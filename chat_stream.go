package deepseek

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// StreamChatCompletionMessage represents a single message in a chat completion stream.
type StreamChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionStream is an interface for receiving streaming chat completion responses.
type ChatCompletionStream interface {
	Recv() (*StreamChatCompletionResponse, error)
	Close() error
}

// chatCompletionStream implements the ChatCompletionStream interface.
type chatCompletionStream struct {
	ctx    context.Context    // Context for cancellation.
	cancel context.CancelFunc // Cancel function for the context.
	resp   *http.Response     // HTTP response from the API call.
	reader *bufio.Reader      // Reader for the response body.
}

// StreamOptions provides options for streaming chat completion responses.
type StreamOptions struct {
	IncludeUsage bool `json:"include_usage"` // Whether to include usage information in the stream. The API returns the usage sometimes even if this is set to false.
}

// StreamUsage represents token usage statistics for a streaming chat completion response. You will get {0 0 0} up until the last stream delta.
type StreamUsage struct {
	PromptTokens            int                     `json:"prompt_tokens"`            // Number of tokens in the prompt.
	CompletionTokens        int                     `json:"completion_tokens"`        // Number of tokens in the completion.
	TotalTokens             int                     `json:"total_tokens"`             // Total number of tokens used.
	PromptCacheHitTokens    int                     `json:"prompt_cache_hit_tokens"`  // Hit Cached Tokens.
	PromptCacheMissTokens   int                     `json:"prompt_cache_miss_tokens"` // Unhit Cached Tokens.
	CompletionTokensDetails CompletionTokensDetails `json:"completion_tokens_details"`
}

type CompletionTokensDetails struct {
	ReasoningTokens int `json:"reasoning_tokens"` //Number of tokens generated by inference model
}

// StreamDelta represents a delta in the chat completion stream.
type StreamDelta struct {
	Role             string     `json:"role,omitempty"`              // Role of the message.
	Content          string     `json:"content"`                     // Content of the message.
	ReasoningContent string     `json:"reasoning_content,omitempty"` // Reasoning content of the message.
	ToolCalls        []ToolCall `json:"tool_calls,omitempty"`        // Optional tool calls related to the message.
}

// Custom unmarshal to support both "reasoning_content" and "reasoning" fields.
// This is necessary because the API may return either field, and we want to prefer "reasoning_content" if present.
func (s *StreamDelta) UnmarshalJSON(data []byte) error {
	type Alias StreamDelta
	aux := &struct {
		Reasoning *string `json:"reasoning,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}
	// Prefer reasoning_content, but fallback to reasoning if present
	if s.ReasoningContent == "" && aux.Reasoning != nil {
		s.ReasoningContent = *aux.Reasoning
	}
	return nil
}

// StreamChoices represents a choice in the chat completion stream.
type StreamChoices struct {
	Index        int         `json:"index"` // Index of the choice.
	Delta        StreamDelta // Delta information for the choice.
	FinishReason string      `json:"finish_reason,omitempty"` // Reason for finishing the generation.
	Logprobs     any         `json:"logprobs,omitempty"`      // Log probabilities for the generated tokens.
}

// StreamChatCompletionResponse represents a single response from a streaming chat completion API call.
type StreamChatCompletionResponse struct {
	ID      string          `json:"id"`              // ID of the response.
	Object  string          `json:"object"`          // Type of object.
	Created int64           `json:"created"`         // Creation timestamp.
	Model   string          `json:"model"`           // Model used.
	Choices []StreamChoices `json:"choices"`         // Choices generated.
	Usage   *StreamUsage    `json:"usage,omitempty"` // Usage statistics (optional).
	Title   string          `json:"title,omitempty"` // Title of the response (optional).
}

// StreamChatCompletionRequest represents the request body for a streaming chat completion API call.
type StreamChatCompletionRequest struct {
	Stream           bool                    `json:"stream,omitempty"`            //Comments: Defaults to true, since it's "STREAM"
	StreamOptions    StreamOptions           `json:"stream_options,omitempty"`    // Optional: Stream options for the request.
	Model            string                  `json:"model"`                       // Required: Model ID, e.g., "deepseek-chat"
	Messages         []ChatCompletionMessage `json:"messages"`                    // Required: List of messages
	FrequencyPenalty float32                 `json:"frequency_penalty,omitempty"` // Optional: Frequency penalty, >= -2 and <= 2
	MaxTokens        int                     `json:"max_tokens,omitempty"`        // Optional: Maximum tokens, > 1
	PresencePenalty  float32                 `json:"presence_penalty,omitempty"`  // Optional: Presence penalty, >= -2 and <= 2
	Temperature      float32                 `json:"temperature,omitempty"`       // Optional: Sampling temperature, <= 2
	TopP             float32                 `json:"top_p,omitempty"`             // Optional: Nucleus sampling parameter, <= 1
	ResponseFormat   *ResponseFormat         `json:"response_format,omitempty"`   // Optional: Custom response format: just don't try, it breaks rn ;)
	Stop             []string                `json:"stop,omitempty"`              // Optional: Stop signals
	Tools            []Tool                  `json:"tools,omitempty"`             // Optional: List of tools
	LogProbs         bool                    `json:"logprobs,omitempty"`          // Optional: Enable log probabilities
	TopLogProbs      int                     `json:"top_logprobs,omitempty"`      // Optional: Number of top tokens with log probabilities, <= 20
	EnableThinking   bool                    `json:"enable_thinking,omitempty"`   // Optional: Enable thinking (for qwen3 api)
}

// Recv receives the next response from the stream.
func (s *chatCompletionStream) Recv() (*StreamChatCompletionResponse, error) {
	reader := s.reader
	for {
		line, err := reader.ReadString('\n') // Read until newline
		if err != nil {
			if err == io.EOF {
				return nil, io.EOF
			}
			return nil, fmt.Errorf("error reading stream: %w", err)
		}

		line = strings.TrimSpace(line)
		if line == "data: [DONE]" {
			return nil, io.EOF // End of stream
		}
		if len(line) > 6 && line[:6] == "data: " {
			trimmed := line[6:] // Trim the "data: " prefix
			var response StreamChatCompletionResponse
			if err := json.Unmarshal([]byte(trimmed), &response); err != nil {
				return nil, fmt.Errorf("unmarshal error: %w, raw data: %s", err, trimmed)
			}
			if response.Usage == nil {
				response.Usage = &StreamUsage{}
			}
			return &response, nil
		}
	}
}

// Close terminates the stream.
func (s *chatCompletionStream) Close() error {
	s.cancel()
	err := s.resp.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to close response body: %w", err)
	}
	return nil
}

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
	PromptTokens     int `json:"prompt_tokens"`     // Number of tokens in the prompt.
	CompletionTokens int `json:"completion_tokens"` // Number of tokens in the completion.
	TotalTokens      int `json:"total_tokens"`      // Total number of tokens used.
}

// StreamDelta represents a delta in the chat completion stream.
type StreamDelta struct {
	Role             string `json:"role,omitempty"`              // Role of the message.
	Content          string `json:"content"`                     // Content of the message.
	ReasoningContent string `json:"reasoning_content,omitempty"` // Reasoning content of the message.
}

// StreamChoices represents a choice in the chat completion stream.
type StreamChoices struct {
	Index        int         `json:"index"` // Index of the choice.
	Delta        StreamDelta // Delta information for the choice.
	FinishReason string      `json:"finish_reason,omitempty"` // Reason for finishing the generation.
	Logprobs     Logprobs    `json:"logprobs,omitempty"`      // Log probabilities for the generated tokens.
}

// StreamChatCompletionResponse represents a single response from a streaming chat completion API call.
type StreamChatCompletionResponse struct {
	ID      string          `json:"id"`              // ID of the response.
	Object  string          `json:"object"`          // Type of object.
	Created int64           `json:"created"`         // Creation timestamp.
	Model   string          `json:"model"`           // Model used.
	Choices []StreamChoices `json:"choices"`         // Choices generated.
	Usage   *StreamUsage    `json:"usage,omitempty"` // Usage statistics (optional).
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
}

type StreamChatCompletionRequestOption func(*StreamChatCompletionRequest)

func NewDefaultStreamChatCompletionRequest(model string, messages []ChatCompletionMessage,
	opts ...StreamChatCompletionRequestOption) *StreamChatCompletionRequest {

	locChatCompletionReq := &StreamChatCompletionRequest{
		Stream:           true,
		Model:            model,
		Messages:         messages,
		MaxTokens:        ConstChatMaxTokensDefault,
		Temperature:      ConstChatTemperatureDefault,
		TopP:             ConstChatTopPDefault,
		PresencePenalty:  ConstChatPresencePenaltyDefault,
		FrequencyPenalty: ConstChatFrequencyPenaltyDefault,
	}
	for _, o := range opts {
		o(locChatCompletionReq)
	}
	return locChatCompletionReq
}

func WithStreamChatUpdate(update bool) StreamChatCompletionRequestOption {
	return func(req *StreamChatCompletionRequest) {
		req.Stream = update
	}
}

func WithStreamChatModel(model string, messages []ChatCompletionMessage) StreamChatCompletionRequestOption {
	return func(r *StreamChatCompletionRequest) {
		r.Model = model
		r.Messages = messages
	}
}

func WithStreamChatMaxTokens(maxTokens int) StreamChatCompletionRequestOption {
	return func(r *StreamChatCompletionRequest) {
		r.MaxTokens = maxTokens
	}
}

func WithStreamChatTemperature(temperature float32) StreamChatCompletionRequestOption {
	return func(r *StreamChatCompletionRequest) {
		r.Temperature = temperature
	}
}

func WithStreamChatTopP(topP float32) StreamChatCompletionRequestOption {
	return func(r *StreamChatCompletionRequest) {
		r.TopP = topP
	}
}

func WithStreamChatPresencePenalty(presencePenalty float32) StreamChatCompletionRequestOption {
	return func(r *StreamChatCompletionRequest) {
		r.PresencePenalty = presencePenalty
	}
}

func WithStreamChatFrequencyPenalty(frequencyPenalty float32) StreamChatCompletionRequestOption {
	return func(r *StreamChatCompletionRequest) {
		r.FrequencyPenalty = frequencyPenalty
	}
}

func WithStreamChatStop(stop []string) StreamChatCompletionRequestOption {
	return func(r *StreamChatCompletionRequest) {
		r.Stop = stop
	}
}

func WithStreamChatLogprobs(logprobs bool) StreamChatCompletionRequestOption {
	return func(r *StreamChatCompletionRequest) {
		r.LogProbs = logprobs
	}
}

func WithCStreamhatTopLogProbs(topLogProbs int) StreamChatCompletionRequestOption {
	return func(r *StreamChatCompletionRequest) {
		r.TopLogProbs = topLogProbs
	}
}

func CheckStreamChatCompletionRequest(request *StreamChatCompletionRequest) error {
	if request == nil {
		return fmt.Errorf("request cannot be nil")
	}
	if request.Model == "" {
		return fmt.Errorf("model cannot be empty")
	}
	if len(request.Messages) == 0 {
		return fmt.Errorf("messages cannot be empty")
	}
	if request.MaxTokens <= ConstChatMaxTokensMin {
		return fmt.Errorf("max tokens must be > 1")
	}
	if request.MaxTokens > ConstChatMaxTokensMax {
		return fmt.Errorf("max tokens must be <= 4000")
	}
	if request.FrequencyPenalty < ConstChatFrequencyPenaltyMin {
		return fmt.Errorf("frequency penalty must be >= %v", ConstChatFrequencyPenaltyMin)
	}
	if request.FrequencyPenalty > ConstChatFrequencyPenaltyMax {
		return fmt.Errorf("frequency penalty must be < =%v", ConstChatFrequencyPenaltyMax)
	}
	if request.TopLogProbs > ConstChatTopLogprobsMax {
		return fmt.Errorf("logprobs must be <= %v", ConstChatTopLogprobsMax)
	}
	if request.Temperature > ConstChatTemperatureMax {
		return fmt.Errorf("temperature must be <= %v", ConstChatTemperatureMax)
	}
	if request.TopP > ConstChatTopPMax {
		return fmt.Errorf("top p must be <= %v", ConstChatTopPMax)
	}

	return nil
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

package deepseek

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ---------------------------------------------------------------------------
// Content block types
// ---------------------------------------------------------------------------

// AnthropicTextBlock represents a text content block in an Anthropic message.
type AnthropicTextBlock struct {
	Type string `json:"type"` // "text"
	Text string `json:"text"`
}

// AnthropicThinkingBlock represents a thinking content block.
type AnthropicThinkingBlock struct {
	Type      string `json:"type"`      // "thinking"
	Thinking  string `json:"thinking"`
	Signature string `json:"signature,omitempty"`
}

// AnthropicToolUseInput represents the typed input arguments for a tool_use block.
// Use this when constructing tool_result responses with structured data.
type AnthropicToolUseInput map[string]interface{}

// AnthropicToolUseBlock represents a tool use content block from the model.
type AnthropicToolUseBlock struct {
	Type  string               `json:"type"`  // "tool_use"
	ID    string               `json:"id"`    // Unique identifier for this tool use
	Name  string               `json:"name"`  // The name of the tool being called
	Input AnthropicToolUseInput `json:"input"` // The input arguments for the tool
}

// AnthropicToolResultBlock represents a tool result content block from the user.
type AnthropicToolResultBlock struct {
	Type        string      `json:"type"`         // "tool_result"
	ToolUseID   string      `json:"tool_use_id"`  // The ID of the tool use this result is for
	Content     interface{} `json:"content"`       // The result content (string or []content block)
	IsError     bool        `json:"is_error,omitempty"` // Whether this is an error result
}

// ---------------------------------------------------------------------------
// Tool types
// ---------------------------------------------------------------------------

// AnthropicTool defines a tool that the model may use.
type AnthropicTool struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	InputSchema interface{} `json:"input_schema"`
}

// AnthropicToolChoice controls which tool (if any) the model should use.
type AnthropicToolChoice struct {
	Type string `json:"type"`          // "auto", "any", "tool", or "none"
	Name string `json:"name,omitempty"` // Required when type is "tool"
}

// ---------------------------------------------------------------------------
// Configuration types
// ---------------------------------------------------------------------------

// AnthropicThinkingConfig configures thinking/reasoning for the model.
type AnthropicThinkingConfig struct {
	Type         string `json:"type"`                    // "enabled"
	BudgetTokens int    `json:"budget_tokens,omitempty"` // Ignored by DeepSeek but included for spec compatibility
}

// AnthropicOutputConfig configures output generation parameters.
type AnthropicOutputConfig struct {
	Effort string `json:"effort,omitempty"` // e.g. "high", "low"
}

// AnthropicMetadata holds optional metadata for the request.
type AnthropicMetadata struct {
	UserID string `json:"user_id,omitempty"`
}

// ---------------------------------------------------------------------------
// Message types
// ---------------------------------------------------------------------------

// AnthropicMessage represents a single message in the conversation.
type AnthropicMessage struct {
	Role    string        `json:"role"`    // "user" or "assistant"
	Content []interface{} `json:"content"` // Array of content blocks
}

// ---------------------------------------------------------------------------
// Request / Response
// ---------------------------------------------------------------------------

// AnthropicRequest represents a request to the Anthropic-compatible API.
type AnthropicRequest struct {
	Model         string                  `json:"model"`
	MaxTokens     int                     `json:"max_tokens"`
	StopSequences []string                `json:"stop_sequences,omitempty"`
	Stream        bool                    `json:"stream,omitempty"`
	System        string                  `json:"system,omitempty"`
	Temperature   float32                 `json:"temperature,omitempty"`
	TopP          float32                 `json:"top_p,omitempty"`
	Thinking      *AnthropicThinkingConfig `json:"thinking,omitempty"`
	OutputConfig  *AnthropicOutputConfig   `json:"output_config,omitempty"`
	Tools         []AnthropicTool          `json:"tools,omitempty"`
	ToolChoice    *AnthropicToolChoice     `json:"tool_choice,omitempty"`
	Metadata      *AnthropicMetadata       `json:"metadata,omitempty"`
	Messages      []AnthropicMessage       `json:"messages"`
}

// AnthropicResponse represents a response from the Anthropic-compatible API.
type AnthropicResponse struct {
	ID           string                `json:"id"`
	Type         string                `json:"type"`
	Role         string                `json:"role"`
	Content      []interface{}         `json:"content"`
	Model        string                `json:"model"`
	StopReason   string                `json:"stop_reason,omitempty"`
	StopSequence string                `json:"stop_sequence,omitempty"`
	Usage        AnthropicUsage        `json:"usage"`
}

// AnthropicUsage represents token usage in an Anthropic response.
type AnthropicUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// ---------------------------------------------------------------------------
// Streaming types
// ---------------------------------------------------------------------------

// AnthropicStream is an interface for receiving streaming Anthropic message responses.
type AnthropicStream interface {
	Recv() (*AnthropicResponse, error)
	Close() error
}

// anthropicStream implements the AnthropicStream interface.
type anthropicStream struct {
	ctx    context.Context
	cancel context.CancelFunc
	resp   *http.Response
	reader *bufio.Reader
	// Accumulated state across SSE events
	response *AnthropicResponse
}

// Recv reads the next SSE event and returns the accumulated response so far.
// On message_stop, io.EOF is returned.
func (s *anthropicStream) Recv() (*AnthropicResponse, error) {
	reader := s.reader
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil, io.EOF
			}
			return nil, fmt.Errorf("error reading stream: %w", err)
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Standard SSE data format: "data: {...}"
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")

		// Parse the SSE event
		var event struct {
			Type    string          `json:"type"`
			Message *AnthropicResponse `json:"message,omitempty"`
		}
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			return nil, fmt.Errorf("error parsing SSE event: %w", err)
		}

		switch event.Type {
		case "message_start":
			// Start of a new message
			s.response = event.Message
		case "message_delta":
			// Delta updates to the message (e.g., stop_reason, usage)
			var delta struct {
				Delta struct {
					StopReason   string        `json:"stop_reason"`
					StopSequence string        `json:"stop_sequence"`
				} `json:"delta"`
				Usage AnthropicUsage `json:"usage"`
			}
			if err := json.Unmarshal([]byte(data), &delta); err != nil {
				return nil, fmt.Errorf("error parsing message_delta: %w", err)
			}
			if s.response != nil {
				if delta.Delta.StopReason != "" {
					s.response.StopReason = delta.Delta.StopReason
				}
				if delta.Delta.StopSequence != "" {
					s.response.StopSequence = delta.Delta.StopSequence
				}
				if delta.Usage.InputTokens > 0 || delta.Usage.OutputTokens > 0 {
					s.response.Usage = delta.Usage
				}
			}
		case "content_block_start":
			// New content block started
			var block struct {
				Index       int         `json:"index"`
				ContentBlock interface{} `json:"content_block"`
			}
			if err := json.Unmarshal([]byte(data), &block); err != nil {
				return nil, fmt.Errorf("error parsing content_block_start: %w", err)
			}
			if s.response != nil {
				// content_block will be properly typed by the JSON unmarshaler
				s.response.Content = append(s.response.Content, block.ContentBlock)
			}
		case "content_block_delta":
			// Delta to an existing content block
			var delta struct {
				Index int `json:"index"`
				Delta struct {
					Type  string `json:"type"`
					Text  string `json:"text"`
					PartialJSON string `json:"partial_json,omitempty"`
				} `json:"delta"`
			}
			if err := json.Unmarshal([]byte(data), &delta); err != nil {
				return nil, fmt.Errorf("error parsing content_block_delta: %w", err)
			}
			if s.response != nil && delta.Index < len(s.response.Content) {
				// Append text to the content block at the given index
				if block, ok := s.response.Content[delta.Index].(map[string]interface{}); ok {
					if delta.Delta.Type == "text_delta" {
						block["text"] = block["text"].(string) + delta.Delta.Text
					}
				}
			}
		case "message_stop":
			// Message is complete
			return nil, io.EOF
		}

		if s.response != nil {
			resp := *s.response
			return &resp, nil
		}
	}
}

// Close terminates the stream.
func (s *anthropicStream) Close() error {
	s.cancel()
	return s.resp.Body.Close()
}

// ---------------------------------------------------------------------------
// Client
// ---------------------------------------------------------------------------

// AnthropicClient wraps the existing Client but uses the Anthropic-compatible base URL.
type AnthropicClient struct {
	Client
}

// NewAnthropicClient creates a new client for the DeepSeek Anthropic-compatible API endpoint.
// It uses BaseURL "https://api.deepseek.com/anthropic" and path "v1/messages".
func NewAnthropicClient(authToken string) *AnthropicClient {
	return &AnthropicClient{
		Client: Client{
			AuthToken: authToken,
			BaseURL:   "https://api.deepseek.com/",
			Path:      "anthropic/v1/messages",
		},
	}
}

// MapClaudeModelToDeepSeek maps Claude model names to DeepSeek models.
// - claude-opus* -> deepseek-v4-pro
// - claude-haiku* or claude-sonnet* -> deepseek-v4-flash
// Unrecognized names are returned unchanged.
func MapClaudeModelToDeepSeek(model string) string {
	if strings.HasPrefix(model, "claude-opus") {
		return "deepseek-v4-pro"
	}
	if strings.HasPrefix(model, "claude-haiku") || strings.HasPrefix(model, "claude-sonnet") {
		return "deepseek-v4-flash"
	}
	return model
}

// CreateAnthropicMessage sends a message to the Anthropic-compatible API and returns the response.
func (c *AnthropicClient) CreateAnthropicMessage(
	ctx context.Context,
	request *AnthropicRequest,
) (*AnthropicResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	// Map Claude model names to DeepSeek models
	request.Model = MapClaudeModelToDeepSeek(request.Model)

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	url := fmt.Sprintf("%s%s", c.BaseURL, c.Path)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("x-api-key", c.AuthToken)
	req.Header.Set("Content-Type", "application/json")

	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= 400 {
		// Reuse existing error handling for API errors
		return nil, HandleAPIError(resp)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var anthropicResp AnthropicResponse
	if err := json.Unmarshal(respBody, &anthropicResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &anthropicResp, nil
}

// CreateAnthropicMessageStream sends a streaming message request to the Anthropic-compatible API.
func (c *AnthropicClient) CreateAnthropicMessageStream(
	ctx context.Context,
	request *AnthropicRequest,
) (AnthropicStream, error) {
	if request == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	request.Stream = true

	// Map Claude model names to DeepSeek models
	request.Model = MapClaudeModelToDeepSeek(request.Model)

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	url := fmt.Sprintf("%s%s", c.BaseURL, c.Path)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("x-api-key", c.AuthToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("cache-control", "no-cache")

	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	if resp.StatusCode >= 400 {
		defer func() {
			_ = resp.Body.Close()
		}()
		return nil, HandleAPIError(resp)
	}

	ctx, cancel := context.WithCancel(ctx)
	stream := &anthropicStream{
		ctx:    ctx,
		cancel: cancel,
		resp:   resp,
		reader: bufio.NewReader(resp.Body),
	}
	return stream, nil
}

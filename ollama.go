package deepseek

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	utils "github.com/cohesion-org/deepseek-go/utils"
	api "github.com/ollama/ollama/api"
)

// IsOllamaRunning checks if the Ollama server is running by sending a GET request to the API
// endpoint. It returns true if the server is running and responds with a 200 OK status.
// http://localhost:11434/api/tags is the endpoint to check if the server is running.
func IsOllamaRunning() bool {
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Get("http://localhost:11434/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// OllamaStreamResponse represents the response format from Ollama when streaming chat completions.
type OllamaStreamResponse struct {
	Model              string      `json:"model"`                          // Name of the model used for the response
	CreatedAt          string      `json:"created_at"`                     // Timestamp when the response was created
	Message            api.Message `json:"message"`                        // Message content and role from Ollama's API
	Done               bool        `json:"done"`                           // Indicates if the stream is finished
	DoneReason         string      `json:"done_reason,omitempty"`          // Optional reason for completion (e.g., "stop", "length")
	TotalDuration      int64       `json:"total_duration,omitempty"`       // Total time taken for the request (nanoseconds)
	LoadDuration       int64       `json:"load_duration,omitempty"`        // Time spent loading the model (nanoseconds)
	PromptEvalCount    int64       `json:"prompt_eval_count,omitempty"`    // Number of tokens evaluated in the prompt
	PromptEvalDuration int64       `json:"prompt_eval_duration,omitempty"` // Time spent evaluating the prompt (nanoseconds)
	EvalCount          int64       `json:"eval_count,omitempty"`           // Number of tokens generated in the response
	EvalDuration       int64       `json:"eval_duration,omitempty"`        // Time spent generating the response (nanoseconds)
}

// ollamaCompletionStream implements the ChatCompletionStream interface for Ollama.
type ollamaCompletionStream struct {
	ctx    context.Context    // Context for stream cancellation
	cancel context.CancelFunc // Function to cancel the context
	resp   *http.Response     // HTTP response from the Ollama API
	reader *bufio.Reader      // Buffered reader for streaming response
}

// convertToOllamaMessages converts deepseek messages to ollama format
func convertToOllamaMessages(messages []ChatCompletionMessage) []api.Message {
	converted := make([]api.Message, len(messages))
	for i, msg := range messages {
		converted[i] = api.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	return converted
}

// convertToOllamaMessages converts deepseek messages to ollama format
func convertToOllamaMessagesWithImage(messages []ChatCompletionMessageWithImage) (err error, _ []api.Message) {
	converted := make([]api.Message, len(messages))
	for i, msg := range messages {
		switch content := msg.Content.(type) {
		case []ContentItem:
			// Join all text items into a single string
			var text []string
			var imageData []api.ImageData

			for _, item := range content {
				if item.Type == "text" {
					text = append(text, item.Text)
				} else if item.Type == "image_url" && item.Image != nil {
					if urlStr, ok := item.Image.URL.(string); ok {
						// Extract base64 data after the comma
						parts := strings.Split(urlStr, ",")
						if len(parts) != 2 {
							return fmt.Errorf("invalid image URL format"), nil
						}
						base64Data, err := base64.StdEncoding.DecodeString(parts[1])
						if err != nil {
							return fmt.Errorf("error decoding to base64 %w", err), nil
						}
						imageData = append(imageData, base64Data)
					}
				}
			}

			converted[i] = api.Message{
				Role:    msg.Role,
				Content: strings.Join(text, "\n"),
				Images:  imageData,
			}

		case string:
			converted[i] = api.Message{
				Role:    msg.Role,
				Content: content,
			}
		default:
			return fmt.Errorf("unsupported content type: %T", content), nil
		}
	}
	return nil, converted
}

// convertToDeepseekResponse converts ollama response to deepseek format
func convertToDeepseekResponse(response api.ChatResponse) *ChatCompletionResponse {
	return &ChatCompletionResponse{
		Model:   response.Model,
		Created: response.CreatedAt.Unix(),
		Choices: []Choice{
			{
				Message: Message{
					Role:    response.Message.Role,
					Content: response.Message.Content,
				},
				FinishReason: response.DoneReason,
			},
		},
		Usage: Usage{
			TotalTokens: response.PromptEvalCount + response.EvalCount,
		},
	}
}

// CreateOllamaChatCompletion sends a chat completion request to the Ollama API
// Note from maintainer: This is a wrapper around the Ollama API. It is not a direct implementation of deepseek-go.
func CreateOllamaChatCompletion(req *ChatCompletionRequest) (ChatCompletionResponse, error) {
	if !IsOllamaRunning() {
		return ChatCompletionResponse{}, fmt.Errorf("Ollama server is not running")
	}

	if req == nil {
		return ChatCompletionResponse{}, fmt.Errorf("request cannot be nil")
	}

	client, err := api.ClientFromEnvironment()
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("failed to create client: %w", err)
	}

	var lastResponse api.ChatResponse
	response := func(response api.ChatResponse) error {
		lastResponse = response
		return nil
	}

	stream := false
	err = client.Chat(context.Background(), &api.ChatRequest{
		Model:    req.Model,
		Messages: convertToOllamaMessages(req.Messages),
		Stream:   &stream,
	}, response)

	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("error sending request: %w", err)
	}

	convertedResponse := convertToDeepseekResponse(lastResponse)
	return *convertedResponse, nil
}

// CreateOllamaChatCompletionStream sends a chat completion request with stream = true and returns the delta
func CreateOllamaChatCompletionStream(
	ctx context.Context,
	request *StreamChatCompletionRequest,
) (*ollamaCompletionStream, error) {
	if !IsOllamaRunning() {
		return &ollamaCompletionStream{}, fmt.Errorf("Ollama server is not running")
	}
	if request == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	c := Client{
		BaseURL: "http://localhost:11434",
	}
	var s bool = true
	// Convert messages to Ollama format
	ollamaRequest := &api.ChatRequest{
		Model:    request.Model,
		Messages: convertToOllamaMessages(request.Messages),
		Stream:   &s,
	}

	req, err := utils.NewRequestBuilder(c.AuthToken).
		SetBaseURL(c.BaseURL).
		SetPath("/api/chat/").
		SetBodyFromStruct(ollamaRequest).
		Build(ctx)

	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}

	resp, err := HandleSendChatCompletionRequest(c, req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, HandleAPIError(resp)
	}

	ctx, cancel := context.WithCancel(ctx)
	stream := &ollamaCompletionStream{
		ctx:    ctx,
		cancel: cancel,
		resp:   resp,
		reader: bufio.NewReader(resp.Body),
	}
	return stream, nil
}

// Recv receives the next response from the Ollama stream
func (s *ollamaCompletionStream) Recv() (*StreamChatCompletionResponse, error) {
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

		var ollamaResp OllamaStreamResponse
		if err := json.Unmarshal([]byte(line), &ollamaResp); err != nil {
			return nil, fmt.Errorf("unmarshal error: %w, raw data: %s", err, line)
		}

		// Convert Ollama response to StreamChatCompletionResponse format
		response := &StreamChatCompletionResponse{
			Model: ollamaResp.Model,
			Choices: []StreamChoices{
				{
					Index: 0,
					Delta: StreamDelta{
						Content: ollamaResp.Message.Content,
						Role:    ollamaResp.Message.Role,
					},
					FinishReason: ollamaResp.DoneReason,
				},
			},
		}

		if ollamaResp.Done && ollamaResp.Message.Content == "" {
			return nil, io.EOF
		}

		return response, nil
	}
}

// Close terminates the Ollama stream
func (s *ollamaCompletionStream) Close() error {
	s.cancel()
	err := s.resp.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to close response body: %w", err)
	}
	return nil
}

// CreateOllamaChatCompletionWithImage sends a chat completion request with image to the Ollama API
// Note from maintainer: This is a wrapper around the Ollama API. It is not a direct implementation of deepseek-go.
func CreateOllamaChatCompletionWithImage(req *ChatCompletionRequestWithImage) (ChatCompletionResponse, error) {
	if !IsOllamaRunning() {
		return ChatCompletionResponse{}, fmt.Errorf("Ollama server is not running")
	}

	if req == nil {
		return ChatCompletionResponse{}, fmt.Errorf("request cannot be nil")
	}

	client, err := api.ClientFromEnvironment()
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("failed to create client: %w", err)
	}

	var lastResponse api.ChatResponse
	response := func(response api.ChatResponse) error {
		lastResponse = response
		return nil
	}

	stream := false
	err, messages := convertToOllamaMessagesWithImage(req.Messages)
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("error converting messages: %w", err)
	}
	err = client.Chat(context.Background(), &api.ChatRequest{
		Model:    req.Model,
		Messages: messages,
		Stream:   &stream,
	}, response)

	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("error sending request: %w", err)
	}

	convertedResponse := convertToDeepseekResponse(lastResponse)
	return *convertedResponse, nil
}

// CreateOllamaChatCompletionStream sends a chat completion request with stream = true and returns the delta
func CreateOllamaChatCompletionStreamWithImage(
	ctx context.Context,
	request *StreamChatCompletionRequestWithImage,
) (*ollamaCompletionStream, error) {

	if !IsOllamaRunning() {
		return &ollamaCompletionStream{}, fmt.Errorf("Ollama server is not running")
	}

	if request == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	err, messages := convertToOllamaMessagesWithImage(request.Messages)
	if err != nil {
		return nil, fmt.Errorf("error converting messages: %w", err)
	}
	var s bool = true
	ollamaRequest := &api.ChatRequest{
		Model:    request.Model,
		Messages: messages,
		Stream:   &s,
	}
	c := Client{
		BaseURL: "http://localhost:11434",
	}
	req, err := utils.NewRequestBuilder(c.AuthToken).
		SetBaseURL(c.BaseURL).
		SetPath("/api/chat/").
		SetBodyFromStruct(ollamaRequest).
		Build(ctx)

	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}

	resp, err := HandleSendChatCompletionRequest(c, req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, HandleAPIError(resp)
	}

	ctx, cancel := context.WithCancel(ctx)
	stream := &ollamaCompletionStream{
		ctx:    ctx,
		cancel: cancel,
		resp:   resp,
		reader: bufio.NewReader(resp.Body),
	}
	return stream, nil
}

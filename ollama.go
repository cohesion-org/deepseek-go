package deepseek

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	utils "github.com/cohesion-org/deepseek-go/utils"
	api "github.com/ollama/ollama/api"
)

// OllamaStreamResponse represents the response format from Ollama
type OllamaStreamResponse struct {
	Model              string      `json:"model"`
	CreatedAt          string      `json:"created_at"`
	Message            api.Message `json:"message"`
	Done               bool        `json:"done"`
	DoneReason         string      `json:"done_reason,omitempty"`
	TotalDuration      int64       `json:"total_duration,omitempty"`
	LoadDuration       int64       `json:"load_duration,omitempty"`
	PromptEvalCount    int64       `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int64       `json:"prompt_eval_duration,omitempty"`
	EvalCount          int64       `json:"eval_count,omitempty"`
	EvalDuration       int64       `json:"eval_duration,omitempty"`
}

// ollamaCompletionStream implements the ChatCompletionStream interface for Ollama
type ollamaCompletionStream struct {
	ctx    context.Context
	cancel context.CancelFunc
	resp   *http.Response
	reader *bufio.Reader
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
		return ChatCompletionResponse{}, err
	}

	convertedResponse := convertToDeepseekResponse(lastResponse)
	return *convertedResponse, nil
}

// CreateOllamaChatCompletionStream sends a chat completion request with stream = true and returns the delta
func CreateOllamaChatCompletionStream(
	ctx context.Context,
	request *StreamChatCompletionRequest,
) (*ollamaCompletionStream, error) {
	if request == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	c := Client{
		BaseURL: "http://localhost:11434",
	}
	request.Stream = true

	req, err := utils.NewRequestBuilder(c.AuthToken).
		SetBaseURL(c.BaseURL).
		SetPath("/api/chat/").
		SetBodyFromStruct(request).
		Build(ctx)

	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}

	resp, err := HandleSendChatCompletionRequest(c, req)
	if err != nil {
		return nil, err
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

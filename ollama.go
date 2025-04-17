package deepseek

import (
	"context"
	"fmt"

	api "github.com/ollama/ollama/api"
)

func CreateOllamaCompletion(req *ChatCompletionRequest) (ChatCompletionResponse, error) {
	if req == nil {
		return ChatCompletionResponse{}, fmt.Errorf("request cannot be nil")
	}
	if len(req.Messages) == 0 {
		return ChatCompletionResponse{}, fmt.Errorf("messages cannot be empty")
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

	messages := make([]api.Message, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = api.Message{
			Role:    msg.Role,
			Content: msg.Content,
			// Image is officially not supported by deepseek yet, what should I do?
			// Tools calls can be added tho, but with a slight modification
		}
	}

	stream := false
	err = client.Chat(context.Background(), &api.ChatRequest{
		Model:    req.Model,
		Messages: messages,
		Stream:   &stream,
	}, response)

	if err != nil {
		return ChatCompletionResponse{}, err
	}

	// Handle the response from Ollma to Deepseek Style
	convertedResponse := &ChatCompletionResponse{
		Model:   lastResponse.Model,
		Created: lastResponse.CreatedAt.Unix(),
		Choices: []Choice{
			{
				Message: Message{
					Role:    lastResponse.Message.Role,
					Content: lastResponse.Message.Content,
				},
				FinishReason: lastResponse.DoneReason,
			},
		},
		Usage: Usage{
			TotalTokens: lastResponse.PromptEvalCount + lastResponse.EvalCount,
		},
	}

	return *convertedResponse, nil
}

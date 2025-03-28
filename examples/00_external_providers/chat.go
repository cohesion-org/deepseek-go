package deepseek_examples

import (
	"context"
	"fmt"
	"log"
	"os"

	deepseek "github.com/cohesion-org/deepseek-go"
	constants "github.com/cohesion-org/deepseek-go/constants"
)

// ExternalProviders demonstrates how to use the Deepseek client with external providers.
// This is not the library you should be using for this but if you do decide to use other models, they can
// be accessed by changing the baseURL to the desired provider and replacing model with the desired model name.
// For instance, using a mistral model from OpenRouter would require changing the baseURL to "https://openrouter.ai/api/v1/"
// and the model to "mistralai/mistral-7b-instruct"

func ExternalProviders() {

	// Azure
	baseURL := "https://models.inference.ai.azure.com/"

	// OpenRouter
	// baseURL := "https://openrouter.ai/api/v1/"

	// Set up the Deepseek client
	client := deepseek.NewClient(os.Getenv("PROVIDER_API_KEY"), baseURL)

	// Create a chat completion request
	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.AzureDeepSeekR1,
		// Model: deepseek.OpenRouterDeepSeekR1,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: constants.ChatMessageRoleUser, Content: "Which is the tallest mountain in the world?"},
		},
	}

	// Send the request and handle the response
	ctx := context.Background()
	response, _, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Print the response
	fmt.Println("Response:", response.Choices[0].Message.Content)
}

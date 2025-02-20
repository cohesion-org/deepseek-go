package deepseek_examples

import (
	"context"
	"fmt"
	"log"
	"os"

	deepseek "github.com/cohesion-org/deepseek-go"
)

// EstimateTokens demonstrates how to estimate the tokens used for a request.
func EstimateTokens() {
	client := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))
	request := deepseek.NewDefaultChatCompletionRequest(
		deepseek.DeepSeekChat,
		[]deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleUser, Content: "The text to evaluate the time is: Who is the greatest singer in the world?"},
		},
	)
	ctx := context.Background()

	tokens := deepseek.EstimateTokensFromMessages(request)
	fmt.Println("Estimated tokens for the request is: ", tokens.EstimatedTokens)
	response, err := client.CreateChatCompletion(ctx, request)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Println("Response:", response.Choices[0].Message.Content, "\nActual Tokens Used:", response.Usage.PromptTokens)
}

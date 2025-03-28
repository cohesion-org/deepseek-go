package deepseek_examples

import (
	"context"
	"fmt"
	"log"
	"os"

	deepseek "github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/constants"
)

func Chat() {
	client := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))
	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: constants.ChatMessageRoleUser, Content: "Which is the tallest mountain in the world?"},
			{Role: constants.ChatMessageRoleSystem, Content: "Answer every question using slang"},
		},
	}
	ctx := context.Background()
	response, _, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Println("Response:", response.Choices[0].Message.Content)
}

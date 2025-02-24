package deepseek_examples

import (
	"context"
	"fmt"
	"log"
	"os"

	deepseek "github.com/cohesion-org/deepseek-go"
)

func Chat() {
	client := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))
	request := deepseek.NewDefaultChatCompletionRequest(
		deepseek.DeepSeekChat,
		[]deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleUser, Content: "Which is the tallest mountain in the world?"},
			{Role: deepseek.ChatMessageRoleSystem, Content: "Answer every question using slang"},
		},
	)
	ctx := context.Background()
	response, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Println("Response:", response.Choices[0].Message.Content)
}

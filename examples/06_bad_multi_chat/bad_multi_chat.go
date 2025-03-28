//this example is considered bad because it manually puts in Role and Content to ChatCompletionMessage. Visit #6 and #7 for more information. It does work tho!

package deepseek_examples

import (
	"context"
	"fmt"
	"log"

	deepseek "github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/constants"
)

// BadMultiChat shows an outdated way of creating a multi-chat completion
func BadMultiChat() {
	client := deepseek.NewClient("DEEPSEEK_API_KEY")
	ctx := context.Background()

	messages := []deepseek.ChatCompletionMessage{{
		Role:    constants.ChatMessageRoleUser,
		Content: "What's the highest mountain in the world? One word response only.",
	}}

	response1, _, err := client.CreateChatCompletion(ctx, &deepseek.ChatCompletionRequest{
		Model:    deepseek.DeepSeekChat,
		Messages: messages,
	})
	if err != nil {
		log.Fatalf("Round 1 failed: %v", err)
	}

	messages = append(messages, deepseek.ChatCompletionMessage{
		Role:    response1.Choices[0].Message.Role,
		Content: response1.Choices[0].Message.Content,
	})

	fmt.Printf("Messages after Round 1: %+v\n", messages)

	messages = append(messages, deepseek.ChatCompletionMessage{
		Role:    constants.ChatMessageRoleUser,
		Content: "What is the second?",
	})

	response2, _, err := client.CreateChatCompletion(ctx, &deepseek.ChatCompletionRequest{
		Model:    deepseek.DeepSeekChat,
		Messages: messages,
	})
	if err != nil {
		log.Fatalf("Round 2 failed: %v", err)
	}

	fmt.Printf("Final messages: %+v\n", append(messages, deepseek.ChatCompletionMessage{
		Role:    response2.Choices[0].Message.Role,
		Content: response2.Choices[0].Message.Content,
	}))
}

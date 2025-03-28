package deepseek_examples

import (
	"context"
	"log"

	deepseek "github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/constants"
)

func MultiChat() {
	client := deepseek.NewClient("DEEPSEEK_API_KEY")
	ctx := context.Background()

	messages := []deepseek.ChatCompletionMessage{{
		Role:    constants.ChatMessageRoleUser,
		Content: "Who is the president of the United States? One word response only.",
	}}

	// Round 1: First API call
	response1, _, err := client.CreateChatCompletion(ctx, &deepseek.ChatCompletionRequest{
		Model:    deepseek.DeepSeekChat,
		Messages: messages,
	})
	if err != nil {
		log.Fatalf("Round 1 failed: %v", err)
	}

	response1Message, err := deepseek.MapMessageToChatCompletionMessage(response1.Choices[0].Message)
	if err != nil {
		log.Fatalf("Mapping to message failed: %v", err)
	}
	messages = append(messages, response1Message)

	log.Printf("The messages after response 1 are: %v", messages)
	// Round 2: Second API call
	messages = append(messages, deepseek.ChatCompletionMessage{
		Role:    constants.ChatMessageRoleUser,
		Content: "Who was the one in the previous term.",
	})

	response2, _, err := client.CreateChatCompletion(ctx, &deepseek.ChatCompletionRequest{
		Model:    deepseek.DeepSeekChat,
		Messages: messages,
	})
	if err != nil {
		log.Fatalf("Round 2 failed: %v", err)
	}

	response2Message, err := deepseek.MapMessageToChatCompletionMessage(response2.Choices[0].Message)
	if err != nil {
		log.Fatalf("Mapping to message failed: %v", err)
	}
	messages = append(messages, response2Message)
	log.Printf("The messages after response 1 are: %v", messages)

}

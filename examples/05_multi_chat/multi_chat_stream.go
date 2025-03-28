package deepseek_examples

import (
	"context"
	"errors"
	"io"
	"log"

	deepseek "github.com/cohesion-org/deepseek-go"
	constants "github.com/cohesion-org/deepseek-go/constants"
)

// MultiChatStream demonstrates how to use the ChatStream API for multi-turn chat completion.
func MultiChatStream() {
	client := deepseek.NewClient("DEEPSEEK_API_KEY")
	ctx := context.Background()

	messages := []deepseek.ChatCompletionMessage{{
		Role:    constants.ChatMessageRoleUser,
		Content: "What's the highest mountain in the world? One word response only.",
	}}

	// First round of conversation
	response1, err := streamChatCompletion(ctx, client, messages)
	if err != nil {
		log.Fatalf("Round 1 failed: %v", err)
	}

	// can't use mappers here to append to messages because the response is a stream
	messages = append(messages, deepseek.ChatCompletionMessage{
		Role:    constants.ChatMessageRoleAssistant,
		Content: response1,
	})

	log.Printf("Messages after Round 1: %+v", messages)

	// Second round of conversation
	messages = append(messages, deepseek.ChatCompletionMessage{
		Role:    constants.ChatMessageRoleUser,
		Content: "What is the second?",
	})

	response2, err := streamChatCompletion(ctx, client, messages)
	if err != nil {
		log.Fatalf("Round 2 failed: %v", err)
	}

	messages = append(messages, deepseek.ChatCompletionMessage{
		Role:    constants.ChatMessageRoleAssistant,
		Content: response2,
	})

	log.Printf("Final messages: %+v", messages)
}

// Helper function to handle streaming chat completion. Just returns the final message for this example.
func streamChatCompletion(ctx context.Context, client *deepseek.Client, messages []deepseek.ChatCompletionMessage) (string, error) {
	request := &deepseek.StreamChatCompletionRequest{
		Model:       deepseek.DeepSeekChat,
		Messages:    messages,
		Stream:      true,
		Temperature: 1.5,
	}

	stream, _, err := client.CreateChatCompletionStream(ctx, request)
	if err != nil {
		return "", err
	}
	defer stream.Close()

	var fullMessage string
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			log.Println("Stream finished")
			break
		}
		if err != nil {
			return "", err
		}
		for _, choice := range response.Choices {
			fullMessage += choice.Delta.Content // Accumulate chunk content
		}
	}
	return fullMessage, nil
}

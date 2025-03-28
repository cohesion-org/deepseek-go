package deepseek_examples

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	deepseek "github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/constants"
)

// Streaming function demonstrates how to use the streaming feature of the DeepSeek API for chat completion.
func Streaming() {
	client := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))
	request := &deepseek.StreamChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: constants.ChatMessageRoleUser, Content: "Just testing if the streaming feature is working or not!"},
		},
		Stream: true,
	}
	ctx := context.Background()

	stream, _, err := client.CreateChatCompletionStream(ctx, request)
	if err != nil {
		log.Fatalf("ChatCompletionStream error: %v", err)
	}
	var fullMessage string
	defer stream.Close()
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			break
		}
		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			break
		}
		for _, choice := range response.Choices {
			fullMessage += choice.Delta.Content // Accumulate chunk content
			log.Println(choice.Delta.Content)
		}
	}
	log.Println("The full message is: ", fullMessage)
}

// StreamingWithReasoningContent function demonstrates how to use the streaming feature of the DeepSeek API for chat completion with reasoning content (R1 Model).
func StreamingWithReasoningContent() {
	client := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))
	request := &deepseek.StreamChatCompletionRequest{
		Model: deepseek.DeepSeekReasoner,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: constants.ChatMessageRoleUser, Content: "Hello, how are you?"},
		},
		Stream: true,
	}
	ctx := context.Background()

	stream, _, err := client.CreateChatCompletionStream(ctx, request)
	if err != nil {
		log.Fatalf("ChatCompletionStream error: %v", err)
	}

	var fullMessage string
	var fullReasoning string
	defer stream.Close()
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			break
		}
		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			break
		}
		for _, choice := range response.Choices {
			fullMessage += choice.Delta.Content            // Accumulate chunk content
			fullReasoning += choice.Delta.ReasoningContent // Accumulate chunk reasoning content
			if choice.Delta.ReasoningContent != "" {
				log.Println("Reasoning: ", choice.Delta.ReasoningContent)
			}
			if choice.Delta.Content != "" {
				log.Println("Content:", choice.Delta.Content)
			}
		}
		if streamUsage := response.Usage; streamUsage != nil && streamUsage.TotalTokens > 0 {
			log.Printf("Prompt tokens: %d, Completion tokens: %d, Total tokens: %d",
				streamUsage.PromptTokens, streamUsage.CompletionTokens, streamUsage.TotalTokens)
		}
	}
	log.Println("Full message: ", fullMessage)
	log.Println("\nFull reasoning: ", fullReasoning)
}

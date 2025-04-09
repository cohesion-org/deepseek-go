package deepseek_examples

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	deepseek "github.com/cohesion-org/deepseek-go"
)

// ExampleChatWithImage demonstrates how to use the OpenRouter API with an image.
func ExampleChatWithImage() {
	request := &deepseek.ChatCompletionRequestWithImage{
		Model: "google/gemini-2.0-flash-001",
		Messages: []deepseek.ChatCompletionMessageWithImage{
			deepseek.NewImageMessage(
				deepseek.ChatMessageRoleUser,
				"How would you name this file in snake case? Only return the name and extension.",
				"https://raw.githubusercontent.com/Vein05/nomnom/refs/heads/main/nomnom.png",
			),
		},
	}
	client := deepseek.NewClient(os.Getenv("OPENROUTER_API_KEY"),
		"https://openrouter.ai/api/v1/")

	response, err := client.CreateChatCompletionWithImage(context.Background(), request)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println("Response:", response.Choices[0].Message.Content)
}

// ExampleStreamWithImage demonstrates how to use the OpenRouter API with a streaming image.
func ExampleStreamWithImage() {
	client := deepseek.NewClient(os.Getenv("OPENROUTER_API_KEY"),
		"https://openrouter.ai/api/v1/")
	request := &deepseek.StreamChatCompletionRequestWithImage{
		Model: "google/gemini-2.0-flash-001",
		Messages: []deepseek.ChatCompletionMessageWithImage{
			deepseek.NewImageMessage(
				deepseek.ChatMessageRoleUser,
				"How would you name this file in snake case? Only return the name and extension.",
				"https://raw.githubusercontent.com/Vein05/nomnom/refs/heads/main/nomnom.png",
			),
		},
		Stream: true,
	}
	ctx := context.Background()
	stream, err := client.CreateChatCompletionStreamWithImage(ctx, request)
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
			fullMessage += choice.Delta.Content
			log.Println(choice.Delta.Content)
		}
	}
}

// ExampleChatWithImageBase64 demonstrates how to use the OpenRouter API with a base64 image.
func ExampleChatWithImageBase64() {
	base64content, err := deepseek.ImageToBase64("https://raw.githubusercontent.com/Vein05/nomnom/refs/heads/main/nomnom.png")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	request := &deepseek.ChatCompletionRequestWithImage{
		Model: "google/gemini-2.0-flash-001",
		Messages: []deepseek.ChatCompletionMessageWithImage{
			deepseek.NewImageMessage(
				deepseek.ChatMessageRoleUser,
				"How would you name this file in snake case? Only return the name and extension.",
				base64content,
			),
		},
	}
	client := deepseek.NewClient(os.Getenv("OPENROUTER_API_KEY"),
		"https://openrouter.ai/api/v1/")

	response, err := client.CreateChatCompletionWithImage(context.Background(), request)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println("Response:", response.Choices[0].Message.Content)
}

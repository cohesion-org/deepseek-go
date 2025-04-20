package deepseek_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"testing"

	"github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/constants"
)

func TestChatCompletionWithOllama(t *testing.T) {
	if !deepseek.IsOllamaRunning() {
		t.Skip("Ollama server is not running, skipping test")
	}

	req := &deepseek.ChatCompletionRequest{
		Model: "llava:latest",
		Messages: []deepseek.ChatCompletionMessage{
			{Role: constants.ChatMessageRoleUser, Content: "Hello how are you?"}},
	}
	res, err := deepseek.CreateOllamaChatCompletion(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.Choices[0].Message.Content == "" {
		t.Fatal("Expected non-empty response")
	}

	t.Logf("The reponse is: %v", res)

}

func TestChatCompletionWithOllamaErrors(t *testing.T) {
	if !deepseek.IsOllamaRunning() {
		t.Skip("Ollama server is not running, skipping test")
	}

	tests := []struct {
		name    string
		req     *deepseek.ChatCompletionRequest
		wantErr bool
	}{
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
		},
		{
			name: "empty messages",
			req: &deepseek.ChatCompletionRequest{
				Model:    "llava:latest",
				Messages: []deepseek.ChatCompletionMessage{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := deepseek.CreateOllamaChatCompletion(tt.req)
			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

// TestChatCompletionWithOllamaStream tests the streaming chat completion with Ollama
func TestChatCompletionWithOllamaStream(t *testing.T) {
	if !deepseek.IsOllamaRunning() {
		t.Skip("Ollama server is not running, skipping test")
	}

	ctx := context.Background()
	req := &deepseek.StreamChatCompletionRequest{
		Stream: true,
		Model:  "llava:latest",
		Messages: []deepseek.ChatCompletionMessage{
			{Role: constants.ChatMessageRoleUser, Content: "Whoa are you?"}},
	}

	res, err := deepseek.CreateOllamaChatCompletionStream(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Close()
	fullMessage := ""

	for {
		response, err := res.Recv()
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
	if fullMessage == "" {
		t.Fatal("Expected non-empty response")
	}

	t.Logf("The reponse is: %s", fullMessage)

}

// TestChatCompletionWithOllamaStreamWithImage tests the streaming chat completion with image
func TestChatCompletionWithOllamaWithImage(t *testing.T) {
	if !deepseek.IsOllamaRunning() {
		t.Skip("Ollama server is not running, skipping test")
	}
	imgData, err := deepseek.ImageToBase64("internal/images/deepseek-go-big.png")
	if err != nil {
		t.Fatal(err)
	}
	req := &deepseek.ChatCompletionRequestWithImage{
		Model: "llava:latest",
		Messages: []deepseek.ChatCompletionMessageWithImage{
			deepseek.NewImageMessage("user", "What is this?", imgData),
		},
	}
	fmt.Printf("Request: %v\n", req)

	res, err := deepseek.CreateOllamaChatCompletionWithImage(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.Choices[0].Message.Content == "" {
		t.Fatal("Expected non-empty response")
	}

	t.Logf("The reponse content is: %s", res.Choices[0].Message.Content)

}

func TestChatCompletionWithOllamaStreamWithImage(t *testing.T) {
	if !deepseek.IsOllamaRunning() {
		t.Skip("Ollama server is not running, skipping test")
	}
	imgData, err := deepseek.ImageToBase64("internal/images/deepseek-go-big.png")
	if err != nil {
		t.Fatal(err)
	}
	req := &deepseek.StreamChatCompletionRequestWithImage{
		Model: "llava:latest",
		Messages: []deepseek.ChatCompletionMessageWithImage{
			deepseek.NewImageMessage("user", "What is this?", imgData),
		},
	}
	ctx := context.Background()
	req.Stream = true
	fmt.Printf("Request: %v\n", req)

	res, err := deepseek.CreateOllamaChatCompletionStreamWithImage(ctx, req)
	if err != nil {
		t.Fatalf("The error is %s", err)
	}

	defer res.Close()
	fullMessage := ""

	for {
		response, err := res.Recv()
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
	if fullMessage == "" {
		t.Fatal("Expected non-empty response")
	}

	t.Logf("The reponse is: %s", fullMessage)

}

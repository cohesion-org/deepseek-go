package deepseek_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/constants"
)

func isOllamaRunning() bool {
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Get("http://localhost:11434/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
func TestChatCompletionWithOllama(t *testing.T) {
	if !isOllamaRunning() {
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
	if !isOllamaRunning() {
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
	if !isOllamaRunning() {
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

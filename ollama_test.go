package deepseek_test

import (
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
	res, err := deepseek.CreateOllamaCompletion(req)
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
				Model:    "llama2:latest",
				Messages: []deepseek.ChatCompletionMessage{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := deepseek.CreateOllamaCompletion(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateOllamaCompletion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package deepseek_examples

import (
	"context"
	"fmt"
	"log"
	"os"

	deepseek "github.com/cohesion-org/deepseek-go"
)

// Anthropic demonstrates the Anthropic-compatible API using DeepSeek's endpoint
// at https://api.deepseek.com/anthropic. Use DeepSeek model names directly
// (deepseek-v4-pro, deepseek-v4-flash) as shown below.
//
// For Claude Code / Claude Desktop compatibility, pass Claude model names and they
// are auto-mapped: claude-opus* -> deepseek-v4-pro, claude-haiku*/claude-sonnet* ->
// deepseek-v4-flash. Unsupported model names also auto-map to deepseek-v4-flash.
//
// Features supported: text, thinking, tool_use, tool_result content blocks.
// Not supported (by DeepSeek): images, documents, cache_control, MCP, code execution.
func Anthropic() {
	key := os.Getenv("DEEPSEEK_API_KEY")
	if key == "" {
		log.Fatal("DEEPSEEK_API_KEY environment variable not set")
	}

	client := deepseek.NewAnthropicClient(key)

	// Non-streaming example with text + thinking
	req := &deepseek.AnthropicRequest{
		Model:     deepseek.DeepSeekV4Pro,
		MaxTokens: 256,
		System:    "You are a helpful assistant. Answer concisely.",
		Thinking: &deepseek.AnthropicThinkingConfig{
			Type: "enabled",
		},
		Messages: []deepseek.AnthropicMessage{
			{
				Role: "user",
				Content: []interface{}{
					deepseek.AnthropicTextBlock{
						Type: "text",
						Text: "Explain the P vs NP problem in one paragraph.",
					},
				},
			},
		},
	}

	ctx := context.Background()
	resp, err := client.CreateAnthropicMessage(ctx, req)
	if err != nil {
		log.Fatalf("anthropic error: %v", err)
	}

	fmt.Printf("Model: %s, Stop reason: %s\n", resp.Model, resp.StopReason)
	for _, block := range resp.Content {
		switch b := block.(type) {
		case map[string]interface{}:
			if b["type"] == "text" {
				fmt.Printf("Text: %s\n", b["text"])
			}
			if b["type"] == "thinking" {
				fmt.Printf("Thinking: %s\n", b["thinking"])
			}
		}
	}
	fmt.Printf("Usage: %d input, %d output tokens\n", resp.Usage.InputTokens, resp.Usage.OutputTokens)

	// Streaming example
	streamReq := &deepseek.AnthropicRequest{
		Model:     deepseek.DeepSeekV4Flash,
		MaxTokens: 128,
		Messages: []deepseek.AnthropicMessage{
			{
				Role: "user",
				Content: []interface{}{
					deepseek.AnthropicTextBlock{
						Type: "text",
						Text: "Say hello in exactly three words.",
					},
				},
			},
		},
	}

	stream, err := client.CreateAnthropicMessageStream(ctx, streamReq)
	if err != nil {
		log.Fatalf("anthropic stream error: %v", err)
	}
	defer stream.Close()

	fmt.Println("\n=== Streaming response ===")
	for {
		event, err := stream.Recv()
		if err != nil {
			break // io.EOF or other error ends the stream
		}
		if event != nil {
			for _, block := range event.Content {
				if b, ok := block.(map[string]interface{}); ok && b["type"] == "text" {
					fmt.Print(b["text"])
				}
			}
		}
	}
	fmt.Println()
}

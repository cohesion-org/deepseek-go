package deepseek_examples

import (
	"context"
	"fmt"
	"log"
	"os"

	deepseek "github.com/cohesion-org/deepseek-go"
)

// ReasoningEffort demonstrates how to control reasoning depth with the
// reasoning_effort parameter on V4 models. Valid values are "high" (default)
// and "max" (deepest reasoning). "low" and "medium" are silently mapped
// to "high"; "xhigh" is silently mapped to "max".
//
// Disabled parameters in thinking mode: temperature, top_p, presence_penalty,
// frequency_penalty (accepted but have no effect).
func ReasoningEffort() {
	key := os.Getenv("DEEPSEEK_API_KEY")
	if key == "" {
		log.Fatal("DEEPSEEK_API_KEY environment variable not set")
	}

	client := deepseek.NewClient(key)

	// Example 1: "max" effort forces the deepest chain-of-thought reasoning.
	// Use this for complex tasks where accuracy matters most.
	reqMax := &deepseek.ChatCompletionRequest{
		Model:           deepseek.DeepSeekV4Pro,
		ReasoningEffort: "max",
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleUser, Content: "Solve the Riemann hypothesis."},
		},
	}

	ctx := context.Background()
	resp, err := client.CreateChatCompletion(ctx, reqMax)
	if err != nil {
		log.Fatalf("max effort error: %v", err)
	}
	fmt.Println("=== Reasoning (max) ===")
	fmt.Println(resp.Choices[0].Message.ReasoningContent)
	fmt.Println("=== Answer (max) ===")
	fmt.Println(resp.Choices[0].Message.Content)

	// Example 2: "high" effort (default) for quicker responses.
	reqHigh := &deepseek.ChatCompletionRequest{
		Model:           deepseek.DeepSeekV4Pro,
		ReasoningEffort: "high",
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleUser, Content: "What is 2+2?"},
		},
	}

	resp, err = client.CreateChatCompletion(ctx, reqHigh)
	if err != nil {
		log.Fatalf("high effort error: %v", err)
	}
	fmt.Println("=== Reasoning (high) ===")
	fmt.Println(resp.Choices[0].Message.ReasoningContent)
	fmt.Println("=== Answer (high) ===")
	fmt.Println(resp.Choices[0].Message.Content)
}

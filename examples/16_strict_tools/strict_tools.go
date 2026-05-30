package deepseek_examples

import (
	"context"
	"fmt"
	"log"
	"os"

	deepseek "github.com/cohesion-org/deepseek-go"
)

// StrictTools demonstrates strict function calling mode (beta). When Function.Strict
// is true, the SDK automatically routes requests through the /beta endpoint and the
// model strictly adheres to the provided JSON Schema.
//
// Strict mode constraints:
//   - All object properties must be listed in "required"
//   - Every object must set additionalProperties: false
//   - String fields: minLength/maxLength are NOT supported; use pattern/format/const instead
//   - Array fields: minItems/maxItems are NOT supported
//   - Supported string formats: email, hostname, ipv4, ipv6, uuid
func StrictTools() {
	key := os.Getenv("DEEPSEEK_API_KEY")
	if key == "" {
		log.Fatal("DEEPSEEK_API_KEY environment variable not set")
	}

	client := deepseek.NewClient(key)

	getWeather := deepseek.Tool{
		Type: "function",
		Function: deepseek.Function{
			Name:        "get_weather",
			Description: "Get weather of a location, the user should supply a location first",
			Strict:      true,
			Parameters: &deepseek.FunctionParameters{
				Type: "object",
				Properties: map[string]interface{}{
					"location": map[string]interface{}{
						"type":        "string",
						"description": "The city and state to get the weather for",
					},
				},
				Required: []string{"location"},
			},
		},
	}

	ctx := context.Background()
	req := &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekV4Flash,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleUser, Content: "What's the weather in Tokyo?"},
		},
		Tools: []deepseek.Tool{getWeather},
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		log.Fatalf("strict tools error: %v", err)
	}

	msg := resp.Choices[0].Message
	fmt.Println("Tool calls:", msg.ToolCalls)

	if len(msg.ToolCalls) > 0 {
		fmt.Printf("Function called: %s(%s)\n",
			msg.ToolCalls[0].Function.Name,
			msg.ToolCalls[0].Function.Arguments,
		)
	}
}

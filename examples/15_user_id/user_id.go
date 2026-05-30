package deepseek_examples

import (
	"context"
	"fmt"
	"log"
	"os"

	deepseek "github.com/cohesion-org/deepseek-go"
)

// UserID demonstrates how to use the user_id parameter for end-user isolation.
// user_id enables three isolation functions:
//  1. Content safety isolation — distinguishes end-users for safety handling
//  2. KVCache isolation — separates cache data for privacy between users
//  3. Scheduling isolation — per-user_id concurrency caps (with increased quotas)
//
// Constraints: must match [a-zA-Z0-9\-_]+ , max 512 characters.
// Do not include user privacy information (emails, names, PII) in user_id.
func UserID() {
	key := os.Getenv("DEEPSEEK_API_KEY")
	if key == "" {
		log.Fatal("DEEPSEEK_API_KEY environment variable not set")
	}

	client := deepseek.NewClient(key)
	ctx := context.Background()

	users := []struct{ id, question string }{
		{"user-42", "What is the capital of France?"},
		{"user-7", "Explain quantum computing in one sentence."},
		{"user-99", "How far is the moon?"},
	}

	for _, u := range users {
		req := &deepseek.ChatCompletionRequest{
			Model:  deepseek.DeepSeekV4Flash,
			UserID: u.id,
			Messages: []deepseek.ChatCompletionMessage{
				{Role: deepseek.ChatMessageRoleUser, Content: u.question},
			},
		}

		resp, err := client.CreateChatCompletion(ctx, req)
		if err != nil {
			log.Fatalf("user %s error: %v", u.id, err)
		}
		fmt.Printf("User %s: %s\n", u.id, resp.Choices[0].Message.Content)
	}
}

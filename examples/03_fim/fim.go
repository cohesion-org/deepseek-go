package deepseek_examples

import (
	"context"
	"fmt"
	"log"
	"os"

	deepseek "github.com/cohesion-org/deepseek-go"
)

// FIM demonstrates how to use the FIM API for FIM completion.
func FIM() {
	client := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))
	request := deepseek.NewDefaultFIMCompletionRequest(deepseek.DeepSeekChat, "def add(a, b):")
	ctx := context.Background()
	response, err := client.CreateFIMCompletion(ctx, request)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println("\n", response.Choices[0].Text)
}

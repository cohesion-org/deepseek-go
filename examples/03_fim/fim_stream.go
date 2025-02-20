package deepseek_examples

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	deepseek "github.com/cohesion-org/deepseek-go"
)

// FIMStream demonstrates how to use the FIMStream API for FIM completion.
func FIMStream() {
	client := deepseek.NewClient("DEEPSEEK_API_KEY")
	request := deepseek.NewDefaultFIMStreamCompletionRequest(deepseek.DeepSeekChat, "def add(a, b): ")
	ctx := context.Background()
	stream, err := client.CreateFIMStreamCompletion(ctx, request)
	if err != nil {
		log.Fatalf("FIMCompletionStream error: %v", err)
	}
	defer stream.FIMClose()
	for {
		response, err := stream.FIMRecv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			break
		}
		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			break
		}
		log.Printf("%v", response.Choices[0].Text)
	}
}

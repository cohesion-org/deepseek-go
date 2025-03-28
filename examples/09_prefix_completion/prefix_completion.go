// please read https://api-docs.deepseek.com/guides/chat_prefix_completion for more information

package deepseek_examples

import (
	"context"
	"fmt"
	"log"
	"os"

	deepseek "github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/constants"
)

// ChatPrefix demonstrates how to use the Chat API for Chat completion with a prefix.
func ChatPrefix() {
	client := deepseek.NewClient(
		os.Getenv("DEEPSEEK_API_KEY"),
		"https://api.deepseek.com/beta/") // Use the beta endpoint

	ctx := context.Background()

	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: constants.ChatMessageRoleUser, Content: "Please write quick sort code"},
			{Role: constants.ChatMessageRoleAssistant, Content: "```python\n", Prefix: true},
		},
		Stop: []string{"```"}, // Stop the prefix when the assistant sends the closing triple backticks
	}
	response, _, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println(response.Choices[0].Message.Content)

}

// ChatPrefixWithJsonMode demonstrates how to use the Chat API for Chat completion with a prefix and JSON mode.
func ChatPrefixWithJsonMode() {
	// Book represents a book in a library
	type Book struct {
		ISBN            string `json:"isbn"`
		Title           string `json:"title"`
		Author          string `json:"author"`
		Genre           string `json:"genre"`
		PublicationYear int    `json:"publication_year"`
		Available       bool   `json:"available"`
	}

	type Books struct {
		Books []Book `json:"books"`
	}
	// Creating a new client using OpenRouter; you can use your own API key and endpoint.
	client := deepseek.NewClient(
		os.Getenv("DEEPSEEK_API_KEY"),
		"https://api.deepseek.com/beta/") // Use the beta endpoint

	ctx := context.Background()

	prompt := `Provide book details in JSON format. Generate 10 JSON objects.
	Please provide the JSON in the following format: { "books": [...] }
	Example: {"isbn": "978-0321765723", "title": "The Lord of the Rings", "author": "J.R.R. Tolkien", "genre": "Fantasy", "publication_year": 1954, "available": true}`

	resp, _, err := client.CreateChatCompletion(ctx, &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: constants.ChatMessageRoleUser, Content: prompt},
			{Role: constants.ChatMessageRoleAssistant, Content: "```json\n", Prefix: true},
		},
		Stop:     []string{"```"}, // Stop the prefix when the assistant sends the closing triple backticks
		JSONMode: true,
	})
	if err != nil {
		log.Fatalf("Failed to create chat completion: %v", err)
	}
	if resp == nil || len(resp.Choices) == 0 {
		log.Fatal("No response or choices found")
	}

	log.Printf("Response: %s", resp.Choices[0].Message.Content)

	extractor := deepseek.NewJSONExtractor(nil)
	var books Books
	if err := extractor.ExtractJSON(resp, &books); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n\nExtracted Books: %+v\n\n", books)

	// Basic validation to check if we got some books
	if len(books.Books) == 0 {
		log.Print("No books were extracted from the JSON response")
	} else {
		fmt.Println("Successfully extracted", len(books.Books), "books.")
	}

}

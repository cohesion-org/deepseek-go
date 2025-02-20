package deepseek_examples

import (
	"context"
	"fmt"
	"log"
	"os"

	deepseek "github.com/cohesion-org/deepseek-go"
)

// JsonMode is a client specific feature that demonstrates how to use the JSON mode for chat completion.
func JsonMode() {
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
		os.Getenv("OPENROUTER_API_KEY"),
		"https://openrouter.ai/api/v1/",
	)
	ctx := context.Background()

	prompt := `Provide book details in JSON format. Generate 10 JSON objects. 
	Please provide the JSON in the following format: { "books": [...] }
	Example: {"isbn": "978-0321765723", "title": "The Lord of the Rings", "author": "J.R.R. Tolkien", "genre": "Fantasy", "publication_year": 1954, "available": true}`

	resp, err := client.CreateChatCompletion(ctx, deepseek.NewDefaultChatCompletionRequest(
		"mistralai/codestral-2501", // Or another suitable model
		[]deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleUser, Content: prompt},
		},
		deepseek.WithChatJSONMode(true),
	))
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

// JsonModeWithSchema is a client specific feature that demonstrates how to use the JSON mode for chat completion with a schema.
func JsonModeWithSchema() {
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

	client := deepseek.NewClient(
		os.Getenv("OPENROUTER_API_KEY"),
		"https://openrouter.ai/api/v1/",
	)
	ctx := context.Background()

	prompt := `Provide book details in JSON format. Generate 10 JSON objects. 
Please provide the JSON in the following format: { "books": [...] }
Example: {"isbn": "978-0321765723", "title": "The Lord of the Rings", "author": "J.R.R. Tolkien", "genre": "Fantasy", "publication_year": 1954, "available": true}`

	resp, err := client.CreateChatCompletion(ctx, deepseek.NewDefaultChatCompletionRequest(
		"mistralai/codestral-2501", // Or another suitable model
		[]deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleUser, Content: prompt},
		},
		deepseek.WithChatJSONMode(true),
	))
	if err != nil {
		log.Fatalf("Failed to create chat completion: %v", err)
	}
	if resp == nil || len(resp.Choices) == 0 {
		log.Fatal("No response or choices found")
	}

	log.Printf("Response: %s", resp.Choices[0].Message.Content)

	// Define the schema (optional, but highly recommended)
	schema := `{
	"type": "object",
	"properties": {
		"books": {
			"type": "array",
			"items": {
				"type": "object",
				"properties": {
					"isbn": {"type": "string"},
					"title": {"type": "string"},
					"author": {"type": "string"},
					"genre": {"type": "string"},
					"publication_year": {"type": "integer"},
					"available": {"type": "boolean"}
				},
				"required": ["isbn", "title", "author", "genre", "publication_year", "available"]
			}
		}
	},
	"required": ["books"]
}`

	extractor := deepseek.NewJSONExtractor([]byte(schema)) // Pass the schema to the extractor
	var books Books
	if err := extractor.ExtractJSON(resp, &books); err != nil {
		log.Fatalf("JSON Extraction Error: %v", err)
	}

	fmt.Printf("\n\nExtracted Books: %+v\n\n", books)

	if len(books.Books) == 0 {
		log.Print("No books were extracted from the JSON response")
	} else {
		fmt.Println("Successfully extracted", len(books.Books), "books.")
	}

}

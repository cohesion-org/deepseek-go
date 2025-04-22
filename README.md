# Deepseek-Go

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE) 
[![Go Report Card](https://goreportcard.com/badge/github.com/cohesion-org/deepseek-go)](https://goreportcard.com/report/github.com/cohesion-org/deepseek-go)

Deepseek-Go is a Go-based API client for the [Deepseek](https://deepseek.com) platform. It provides a clean and type-safe interface to interact with Deepseek's AI features, including chat completions with streaming, token usage tracking, and more.


## Installation

```sh
go get github.com/cohesion-org/deepseek-go
```
deepseek-go currently uses `go 1.24.0`

## Features

- **Chat Completion**: Easily send chat messages and receive responses from Deepseek's AI models. It also supports streaming.
- **Modular Design**: The library is structured into reusable components for building, sending, and handling requests and responses.
- **External Providers**: Deepseek-go also supports external providers like OpenRouter, Azure, and even Ollama. 
- **MIT License**: Open-source and free for both personal and commercial use.

The recent gain in popularity and cybersecurity issues Deepseek has seen makes for many problems while using the API. Please refer to the [status](https://status.deepseek.com/) page for the current status.

## Getting Started

Here's a quick example of how to use the library:

### Prerequisites

Before using the library, ensure you have:
- A valid Deepseek API key.
- Go installed on your system.

### Supported Models

- **deepseek-chat**  
  A versatile model designed for conversational tasks. <br/>
  Usage: `Model: deepseek.DeepSeekChat`

- **deepseek-reasoner**  
  A specialized model for reasoning-based tasks.  
  Usage: `Model: deepseek.DeepSeekReasoner`. <br/>
  **Note:** The [reasoner](https://api-docs.deepseek.com/guides/reasoning_model) requires unique conditions. Please refer to this issue [#8](https://github.com/cohesion-org/deepseek-go/issues/8). 

### External Providers
- **Azure DeepSeekR1**  
	Same as `deepseek-reasoner`, but provided by Azure. <br/>
	Usage: `Model: deepseek.AzureDeepSeekR1`

- **OpenRouter DeepSeek1** <br/>
	Same as `deepseek-reasoner`, but provided by OpenRouter. <br/>
  	Usage: `Model: deepseek.OpenRouterR1`

- **Ollama Support** <br/>
	Please read [Ollama Support](#ollama) for more info about this!


<details open>
<summary> Chat </summary>

### Example for chatting with deepseek

Even more examples are avilable [here](/examples/README.md)

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	deepseek "github.com/cohesion-org/deepseek-go"
	
)

func main() {
	// Set up the Deepseek client
	client := deepseek.NewClient("") // Empty API key triggers env lookup for "DEEPSEEK_API_KEY"

	// Create a chat completion request
	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleSystem, Content: "Answer every question using slang."},
			{Role: deepseek.ChatMessageRoleUser, Content: "Which is the tallest mountain in the world?"},
		},
	}

	// Send the request and handle the response
	ctx := context.Background()
	response, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Print the response
	fmt.Println("Response:", response.Choices[0].Message.Content)
}
```
</details>

## More Examples:

<details>
<summary> Using external providers such as Azure or OpenRouter. </summary>

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	deepseek "github.com/cohesion-org/deepseek-go"
)

func main() {

	// Azure
	baseURL := "https://models.inference.ai.azure.com/"

	// OpenRouter
	// baseURL := "https://openrouter.ai/api/v1/"

	// Set up the Deepseek client
    client := deepseek.NewClient(os.Getenv("PROVIDER_API_KEY"), baseURL)

	// Create a chat completion request
	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.AzureDeepSeekR1,
		// Model: deepseek.OpenRouterDeepSeekR1,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleUser, Content: "Which is the tallest mountain in the world?"},
		},
	}

	// Send the request and handle the response
	ctx := context.Background()
	response, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Print the response
	fmt.Println("Response:", response.Choices[0].Message.Content)
}
```

Note: If you wish to use other providers that are not supported by us, you can simply extend the baseURL(as shown above), and pass the name of your model as a string to `Model` while creating the `ChatCompletionRequest`. This will work as long as the provider follows the same API structure as Azure or OpenRouter.


</details>

<details >
	<summary> Sending other params like Temp, Stop </summary>
	<strong> You just need to extend the ChatCompletionMessage with the supported parameters. </strong>

```go
	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleUser, Content: "What is the meaning of deepseek"},
			{Role: deepseek.ChatMessageRoleSystem, Content: "Answer every question using slang"},
		},
		Temperature: 1.0,
		Stop:        []string{"yo", "hello"},
		ResponseFormat: &deepseek.ResponseFormat{
			Type: "text",
		},
	}
```

</details>

<details >
	<summary> Multi-Conversation with Deepseek. </summary>

```go
package deepseek_examples

import (
	"context"
	"log"

	deepseek "github.com/cohesion-org/deepseek-go"
)

func MultiChat() {
	client := deepseek.NewClient("DEEPSEEK_API_KEY")
	ctx := context.Background()

	messages := []deepseek.ChatCompletionMessage{{
		Role:    deepseek.ChatMessageRoleUser,
		Content: "Who is the president of the United States? One word response only.",
	}}

	// Round 1: First API call
	response1, err := client.CreateChatCompletion(ctx, &deepseek.ChatCompletionRequest{
		Model:    deepseek.DeepSeekChat,
		Messages: messages,
	})
	if err != nil {
		log.Fatalf("Round 1 failed: %v", err)
	}

	response1Message, err := deepseek.MapMessageToChatCompletionMessage(response1.Choices[0].Message)
	if err != nil {
		log.Fatalf("Mapping to message failed: %v", err)
	}
	messages = append(messages, response1Message)

	log.Printf("The messages after response 1 are: %v", messages)
	// Round 2: Second API call
	messages = append(messages, deepseek.ChatCompletionMessage{
		Role:    deepseek.ChatMessageRoleUser,
		Content: "Who was the one in the previous term.",
	})

	response2, err := client.CreateChatCompletion(ctx, &deepseek.ChatCompletionRequest{
		Model:    deepseek.DeepSeekChat,
		Messages: messages,
	})
	if err != nil {
		log.Fatalf("Round 2 failed: %v", err)
	}

	response2Message, err := deepseek.MapMessageToChatCompletionMessage(response2.Choices[0].Message)
	if err != nil {
		log.Fatalf("Mapping to message failed: %v", err)
	}
	messages = append(messages, response2Message)
	log.Printf("The messages after response 1 are: %v", messages)

}

```

</details>

<details>
<summary> Chat with Streaming </summary>

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	deepseek "github.com/cohesion-org/deepseek-go"
)

func main() {
	client := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))
	request := &deepseek.StreamChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleUser, Content: "Just testing if the streaming feature is working or not!"},
		},
		Stream: true,
	}
	ctx := context.Background()

	stream, err := client.CreateChatCompletionStream(ctx, request)
	if err != nil {
		log.Fatalf("ChatCompletionStream error: %v", err)
	}
	var fullMessage string
	defer stream.Close()
	for {
		response, err := stream.Recv()
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
	log.Println("The full message is: ", fullMessage)
}
```
</details>

<details>
<summary> Get the balance(s) of the user. </summary>

```go
package main

import (
	"context"
	"log"
	"os"

	deepseek "github.com/cohesion-org/deepseek-go"
)

func main() {
	client := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))
	ctx := context.Background()
	balance, err := deepseek.GetBalance(client, ctx)
	if err != nil {
		log.Fatalf("Error getting balance: %v", err)
	}

	if balance == nil {
		log.Fatalf("Balance is nil")
	}

	if len(balance.BalanceInfos) == 0 {
		log.Fatalf("No balance information returned")
	}
	log.Printf("%+v\n", balance)
}
```
</details>

<details>
<summary> Get the list of All the models the API supports right now. This is different from what deepseek-go might support. </summary>

```go
func ListModels() {
	client := deepseek.NewClient("DEEPSEEK_API_KEY")
	ctx := context.Background()
	models, err := deepseek.ListAllModels(client, ctx)
	if err != nil {
		t.Fatalf("Error listing models: %v", err)
	}
	fmt.Printf("\n%+v\n", models)
}
```
</details>

<details> 
<summary> Get the estimated tokens for the request. </summary>

This is adpated from [the  Deepseek's estimation](https://api-docs.deepseek.com/quick_start/token_usage).

```go
func Estimation() {
	client := deepseek.NewClient("DEEPSEEK_API_KEY"))
	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleSystem, Content: "Just respond with the time it might take you to complete this request."},
			{Role: deepseek.ChatMessageRoleUser, Content: "The text to evaluate the time is: Who is the greatest singer in the world?"},
		},
	}
	ctx := context.Background()

	tokens := deepseek.EstimateTokensFromMessages(request)
	fmt.Println("Estimated tokens for the request is: ", tokens.EstimatedTokens)
	response, err := client.CreateChatCompletion(ctx, request)

	if err != nil {
		log.Fatalf("error: %v", err)
	}
	
	fmt.Println("Response:", response.Choices[0].Message.Content, "\nActual Tokens Used:", response.Usage.PromptTokens)
}

```

</details>

<details> 

<summary> JSON mode for JSON extraction</summary>

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	deepseek "github.com/cohesion-org/deepseek-go"
)

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

	resp, err := client.CreateChatCompletion(ctx, &deepseek.ChatCompletionRequest{
		Model: "mistralai/codestral-2501", // Or another suitable model
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleUser, Content: prompt},
		},
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
```
You can see more examples inside the examples folder.

</details>

<details> <summary> Add more settings to your client with NewClientWithOptions </summary>

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/cohesion-org/deepseek-go"
)

func main() {
    client, err := deepseek.NewClientWithOptions("your-api-key",
        deepseek.WithBaseURL("https://custom-api.com/"),
        deepseek.WithTimeout(10*time.Second),
    )
    if err != nil {
        log.Fatalf("Error creating client: %v", err)
    }

    fmt.Printf("Client initialized with BaseURL: %s and Timeout: %v\n", client.BaseURL, client.Timeout)
}
 ```
See the examples folder for more information.
</details>

<details> 
<summary> FIM Mode(Beta) </summary>

In FIM (Fill In the Middle) completion, users can provide a prefix and a suffix (optional), and the model will complete the content in between. FIM is commonly used for content completion、code completion.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	deepseek "github.com/cohesion-org/deepseek-go"
)

func FIM() {
	client := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))
	request := &deepseek.FIMCompletionRequest{
		Model:  deepseek.DeepSeekChat,
		Prompt: "def add(a, b):",
	}
	ctx := context.Background()
	response, err := client.CreateFIMCompletion(ctx, request)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println("\n", response.Choices[0].Text)
}

```

</details>

<details> 
<summary> Chat Prefix Completion (Beta)</summary>
The chat prefix completion follows the [Chat Completion API](https://api-docs.deepseek.com/guides/chat_prefix_completion), where users provide an assistant's prefix message for the model to complete the rest of the message.

```go

package main

import (
	"context"
	"fmt"
	"log"

	deepseek "github.com/cohesion-org/deepseek-go"
)

func ChatPrefix() {
	client := deepseek.NewClient(
		DEEPSEEK_API_KEY,
		"https://api.deepseek.com/beta/") // Use the beta endpoint

	ctx := context.Background()

	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleUser, Content: "Please write quick sort code"},
			{Role: deepseek.ChatMessageRoleAssistant, Content: "```python", Prefix: true},
		},
		Stop: []string{"```"}, // Stop the prefix when the assistant sends the closing triple backticks
	}
	response, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println(response.Choices[0].Message.Content)

}

```
</details>

<details>
<summary> Using external providers with image support (OpenRouter) </summary>

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    deepseek "github.com/cohesion-org/deepseek-go"
)

func main() {
    // Create request with image URL
    request := &deepseek.ChatCompletionRequestWithImage{
        Model: "google/gemini-2.0-flash-001",
        Messages: []deepseek.ChatCompletionMessageWithImage{
            deepseek.NewImageMessage(
                deepseek.ChatMessageRoleUser,
                "Describe this image",
                "https://example.com/path/to/image.jpg",
            ),
        },
    }

    // Initialize client with OpenRouter
    client := deepseek.NewClient(
        os.Getenv("OPENROUTER_API_KEY"),
        "https://openrouter.ai/api/v1/",
    )

    // Send request and get response
    response, err := client.CreateChatCompletionWithImage(context.Background(), request)
    if err != nil {
        log.Fatalf("error: %v", err)
    }

    fmt.Println("Response:", response.Choices[0].Message.Content)
}
```

For more advanced examples including streaming and base64 image support, see [OpenRouter Images Examples](/examples/13_openrouter_images/openrouter_images.go).

</details>

---
## Getting a Deepseek Key

To use the Deepseek API, you need an API key. You can obtain one by signing up on the [Deepseek website](https://platform.deepseek.com/api_keys)

## Ollama

Deepseek-go supports the usage of Ollama, but because of Ollama not following OpenAI policy, there are some extra types you need to be aware about. This is still an experimental feature so please understand that. 

You can find all information about it at [Ollama Docs](/examples/ollama.md). 

---


## Running Tests

### Setup

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

2. Add your DeepSeek API key to `.env`:
   ```
   TEST_DEEPSEEK_API_KEY=your_api_key_here
   ```

3. (Optional) Configure test timeout:
   ```
   # Default is 30s, increase for slower connections
   TEST_TIMEOUT=1m
   ```

### Test Organization

The tests are organized into several files and folders:

### Main Package
- `client_test.go`: Client configuration and error handling
- `chat_test.go`: Chat completion functionality 
- `chat_stream_test.go`: Chat streaming functionality
- `models_test.go`: Model listing and retrieval
- `balance_test.go`: Account balance operations
- `tokens_test.go`: Token estimation utilities
- `json_test.go`: JSON mode for extraction
- `fim_test.go`: Tests for the FIM beta implementation
<!-- - `errors_test.go`: Tests the error handler -->
- `requestHandler_test.go`: Tests for the request handler
- `responseHandler_test.go`: Tests for the response handler

### Utils Package
- `utils/requestBuilder_test.go`: Tests for the request builder

### Running Tests

1. Run all tests (requires API key):
   ```bash
   go test -v ./...
   ```

2. Run tests in short mode (skips API calls):
   ```bash
   go test -v -short ./...
   ```

3. Run tests with race detection:
   ```bash
   go test -v -race ./...
   ```

4. Run tests with coverage:
   ```bash
   go test -v -coverprofile=coverage.txt -covermode=atomic ./...
   ```

   View coverage in browser:
   ```bash
   go tool cover -html=coverage.txt
   ```

5. Run specific test:
   ```bash
   # Example: Run only chat completion tests
   go test -v -run TestCreateChatCompletion ./...
   ```
## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Credits
- **`chat.go` Inspiration**: Adapted from [sashabaranov/go-openai](https://github.com/sashabaranov/go-openai/tree/master).
- **`json.go` Inspiration**: Thanks a lot for [Mr. Peter](https://github.com/peterlodri92).

---

## Images

<div style="display:flex; justify-content: space-between; margin:20px;">
  <img src="internal/images/deepseek-go-big.png" alt="Deepseek Go Big Logo" style="border-radius:2%;"
  height=250px>
  <img src="internal/images/deepseek-go.png" alt="Deepseek Go Logo" style="scale: 90%; border-radius:100%"
  height=250px>

</div>

Feel free to contribute, open issues, or submit PRs to help improve Deepseek-Go! Let us know if you encounter any issues.

# Deepseek-Go

<p align="center">
  <img src="internal/images/deepseek-go-big.png" alt="deepseek-go" width="400">
</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/cohesion-org/deepseek-go"><img src="https://pkg.go.dev/badge/github.com/cohesion-org/deepseek-go.svg" alt="Go Reference"></a>
  <a href="https://goreportcard.com/report/github.com/cohesion-org/deepseek-go"><img src="https://goreportcard.com/badge/github.com/cohesion-org/deepseek-go" alt="Go Report Card"></a>
  <a href="https://github.com/cohesion-org/deepseek-go/actions/workflows/test.yml"><img src="https://github.com/cohesion-org/deepseek-go/actions/workflows/test.yml/badge.svg" alt="CI"></a>
  <a href="https://github.com/cohesion-org/deepseek-go/releases"><img src="https://img.shields.io/github/v/release/cohesion-org/deepseek-go" alt="Release"></a>
  <a href="https://github.com/cohesion-org/deepseek-go"><img src="https://img.shields.io/github/go-mod/go-version/cohesion-org/deepseek-go" alt="Go Version"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="MIT License"></a>
</p>

Deepseek-Go is a Go-based API client for the [Deepseek](https://deepseek.com) platform. It provides a clean and type-safe interface to interact with Deepseek's AI features, including chat completions with streaming, token usage tracking, and more.


## Installation

[![Go Reference](https://pkg.go.dev/badge/github.com/cohesion-org/deepseek-go.svg)](https://pkg.go.dev/github.com/cohesion-org/deepseek-go)

```sh
go get github.com/cohesion-org/deepseek-go
```
deepseek-go currently uses `go 1.26.0`

## Features

- **Chat Completion**: Send chat messages and receive responses from DeepSeek's V4 models with streaming support.
- **Anthropic API**: Full Anthropic-compatible endpoint support via `NewAnthropicClient` with content blocks, tool use, and streaming.
- **Thinking Mode**: Chain-of-thought reasoning with configurable `reasoning_effort` (`"high"` / `"max"`).
- **Tool Calling**: Function calling with standard and strict mode (beta, with automatic `/beta` routing).
- **FIM Completion**: Fill-in-the-Middle completions for code generation, with streaming.
- **JSON Output**: Structured JSON output with schema extraction via `ResponseFormat`.
- **External Providers**: OpenRouter, Azure, Ollama, and any OpenAI-compatible endpoint.
- **Balance & Models**: Check account balance and list available models.
- **Token Estimation**: Client-side token counting for Chinese and English text.
- **Modular Design**: Reusable components for building, sending, and handling requests and responses.
- **MIT License**: Open-source and free for both personal and commercial use.

For API status and uptime, refer to the [DeepSeek Status](https://status.deepseek.com/) page.

## Contents

- [Installation](#installation)
- [Features](#features)
- [Getting Started](#getting-started)
  - [Supported Models](#supported-models)
  - [External Providers](#external-providers)
- [Examples](#more-examples)
- [Getting a DeepSeek Key](#getting-a-deepseek-key)
- [Ollama Support](#ollama)
- [Running Tests](#running-tests)
- [Who Uses deepseek-go](#who-uses-deepseek-go)
- [License](#license)

## Getting Started

Here's a quick example of how to use the library:

### Prerequisites

Before using the library, ensure you have:
- A valid Deepseek API key.
- Go installed on your system.

### Supported Models

- **deepseek-v4-flash** (current)  
  Flagship model with 1M context, 384K max output. Supports thinking mode, tool calls, JSON output, FIM, and prefix completion. <br/>
  Usage: `Model: deepseek.DeepSeekV4Flash`

- **deepseek-v4-pro** (current)  
  Premium reasoning model with 1M context, 384K max output. Best for complex reasoning and agent tasks. <br/>
  Usage: `Model: deepseek.DeepSeekV4Pro`

- **deepseek-chat** (deprecated, sunset 2026/07/24)  
  Maps to `deepseek-v4-flash` non-thinking mode. <br/>
  Usage: `Model: deepseek.DeepSeekV4Flash` — emits a deprecation warning to stderr.

- **deepseek-reasoner** (deprecated, sunset 2026/07/24)  
  Maps to `deepseek-v4-flash` thinking mode. <br/>
  Usage: `Model: deepseek.DeepSeekReasoner` — emits a deprecation warning to stderr. 

### External Providers
- **Azure DeepSeekR1**  
	DeepSeek R1 provided by Azure. <br/>
	Usage: `Model: deepseek.AzureDeepSeekR1`

- **OpenRouter** <br/>
	OpenRouter provides access to DeepSeek R1 and distill models. <br/>
	Usage: `Model: deepseek.OpenRouterDeepSeekR1` (and other `OpenRouterDeepSeek*` constants)

- **MiniMax** <br/>
	Models: `deepseek.MiniMaxM3`, `deepseek.MiniMaxM2_7` <br/>
	OpenAI-compatible API: `deepseek.MiniMaxBaseURL` (global), `deepseek.MiniMaxCNBaseURL` (China) <br/>
	Anthropic-compatible API: `deepseek.MiniMaxAnthropicAPIBaseURL` (global), `deepseek.MiniMaxCNAnthropicAPIBaseURL` (China) <br/>
	When using `NewAnthropicClient`, set `Path` to `"/v1/messages"` with an Anthropic-compatible API base URL.

- **Ollama Support** <br/>
	Please read [Ollama Support](#ollama) for more info about this!


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
		Model: deepseek.DeepSeekV4Flash,
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

	// MiniMax (global)
	// baseURL := deepseek.MiniMaxBaseURL

	// MiniMax (China)
	// baseURL := deepseek.MiniMaxCNBaseURL

	// Set up the Deepseek client
    client := deepseek.NewClient(os.Getenv("PROVIDER_API_KEY"), baseURL)

	// Create a chat completion request
	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.AzureDeepSeekR1,
		// Model: deepseek.OpenRouterDeepSeekR1,
		// Model: deepseek.MiniMaxM3,
		// Model: deepseek.MiniMaxM2_7,
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
		Model: deepseek.DeepSeekV4Flash,
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
		Model:    deepseek.DeepSeekV4Flash,
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
		Model:    deepseek.DeepSeekV4Flash,
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
		Model: deepseek.DeepSeekV4Flash,
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
		Model: deepseek.DeepSeekV4Flash,
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

**Using local model servers without an API key:**

By default, `NewClientWithOptions` requires an API key. If you're connecting to a local model server (e.g., vLLM, llama.cpp, LocalAI) that doesn't require authentication, use `WithoutAPIKeyValidation()`:

```go
client, err := deepseek.NewClientWithOptions("",
    deepseek.WithBaseURL("<BASE_URL>"),
    deepseek.WithoutAPIKeyValidation(),
)
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
		Model:  deepseek.DeepSeekV4Flash,
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
		Model: deepseek.DeepSeekV4Flash,
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
<details>
<summary> Thinking Mode with reasoning_effort </summary>

DeepSeek V4 models support chain-of-thought reasoning. Control the reasoning depth with `ReasoningEffort`:

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
	client := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))

	request := &deepseek.ChatCompletionRequest{
		Model:           deepseek.DeepSeekV4Pro,
		ReasoningEffort: "max", // "high" (default) or "max"
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleUser, Content: "Explain quantum entanglement."},
		},
	}

	ctx := context.Background()
	response, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Println("Reasoning:", response.Choices[0].Message.ReasoningContent)
	fmt.Println("Answer:", response.Choices[0].Message.Content)
}
```

Note: When thinking mode is active, `temperature`, `top_p`, `presence_penalty`, and `frequency_penalty` are accepted but have no effect.

</details>

<details>
<summary> Tool Calling with Strict Mode </summary>

Standard function calling and beta strict mode (automatically routes to `/beta` when `Strict: true`):

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
	client := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))

	getWeather := deepseek.Tool{
		Type: "function",
		Function: deepseek.Function{
			Name:        "get_weather",
			Description: "Get weather of a location",
			Strict:      true, // triggers /beta endpoint routing
			Parameters: &deepseek.FunctionParameters{
				Type: "object",
				Properties: map[string]interface{}{
					"location": map[string]interface{}{
						"type": "string",
					},
				},
				Required: []string{"location"},
			},
		},
	}

	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekV4Flash,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleUser, Content: "What's the weather in Tokyo?"},
		},
		Tools: []deepseek.Tool{getWeather},
	}

	response, err := client.CreateChatCompletion(context.Background(), request)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Println("Tool calls:", response.Choices[0].Message.ToolCalls)
}
```

</details>

<details>
<summary> Anthropic-Compatible API </summary>

Use `NewAnthropicClient` to talk to DeepSeek's Anthropic-compatible endpoint. Claude model names are automatically mapped to DeepSeek models:

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
	client := deepseek.NewAnthropicClient(os.Getenv("DEEPSEEK_API_KEY"))

	request := &deepseek.AnthropicRequest{
		Model:     "claude-opus-4-7", // auto-mapped to deepseek-v4-pro
		MaxTokens: 1024,
		System:    "You are a helpful assistant.",
		Messages: []deepseek.AnthropicMessage{
			{
				Role: "user",
				Content: []interface{}{
					deepseek.AnthropicTextBlock{Type: "text", Text: "Hello, how are you?"},
				},
			},
		},
	}

	response, err := client.CreateAnthropicMessage(context.Background(), request)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for _, block := range response.Content {
		if text, ok := block.(deepseek.AnthropicTextBlock); ok {
			fmt.Println(text.Text)
		}
	}
}
```

See `examples/` for more Anthropic usage patterns including tool use and streaming.

</details>

<details>
<summary> Multi-tenant isolation with user_id </summary>

Use `UserID` to isolate content safety, KVCache, and scheduling per end-user:

```go
request := &deepseek.ChatCompletionRequest{
	Model:  deepseek.DeepSeekV4Flash,
	UserID: "user-42", // [a-zA-Z0-9-_]+, max 512 chars
	Messages: []deepseek.ChatCompletionMessage{
		{Role: deepseek.ChatMessageRoleUser, Content: "Hello!"},
	},
}
```

Do not include user privacy information (emails, names, PII) in `user_id`.

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

2. Add your DeepSeek API key to `.env` if you plan to run live integration tests:
   ```
   TEST_DEEPSEEK_API_KEY=your_api_key_here
   ```

3. (Optional) Add your OpenRouter key if you plan to run the integration image tests:
   ```
   OPENROUTER_API_KEY=your_openrouter_key_here
   ```

4. (Optional) Configure test timeout:
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

1. Run the default offline test suite:
   ```bash
   go test -v ./...
   ```

   This uses the built-in mock DeepSeek server for the package's API contract tests, so it does not require a DeepSeek API key.

2. Run tests in short mode:
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

5. Run the live integration lane against the real APIs:
   ```bash
   DEEPSEEK_LIVE_TESTS=1 go test -v -tags=integration ./...
   ```

   Notes:
   `TEST_DEEPSEEK_API_KEY` or `DEEPSEEK_API_KEY` is required for the DeepSeek integration tests.
   `OPENROUTER_API_KEY` is also required for the image/OpenRouter integration tests.

6. Use the Makefile shortcuts:
   ```bash
   make test
   make test-short
   make test-race
   make test-integration
   ```

7. Run a specific test:
   ```bash
   # Example: Run only chat completion tests
   go test -v -run TestCreateChatCompletion ./...
   ```
## Who Uses deepseek-go

deepseek-go is trusted by **149+ projects** and counting.

Notable dependents include [Ech0](https://github.com/lin-snow/Ech0) (1,900+ stars) and many others.

[![Used by](https://img.shields.io/badge/used_by-149-blue)](https://github.com/cohesion-org/deepseek-go/network/dependents)

*Using deepseek-go in your project? Open a PR to add your name here!*

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

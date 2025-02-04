# Deepseek-Go

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE) 
[![Go Report Card](https://goreportcard.com/badge/github.com/cohesion-org/deepseek-go)](https://goreportcard.com/report/github.com/cohesion-org/deepseek-go)

Deepseek-Go is a Go-based API wrapper for the [Deepseek](https://deepseek.com) platform. It provides a clean and type-safe interface to interact with Deepseek's AI features, including chat completions with streaming, token usage tracking, and more.

This library is designed for developers building Go applications that require seamless integration with Deepseek AI.

The recent gain in popularity and cybersecurity issues Deepseek has seen makes for many problems while using the API. Please refer to the [status](https://status.deepseek.com/) page for the current status.

## Features

- **Chat Completion**: Easily send chat messages and receive responses from Deepseek's AI models. It also supports streaming.
- **Modular Design**: The library is structured into reusable components for building, sending, and handling requests and responses.
- **MIT License**: Open-source and free for both personal and commercial use.

## Installation

To use Deepseek-Go, ensure you have Go installed, and run:

```sh
go get github.com/cohesion-org/deepseek-go
```

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

- **Azure DeepSeekR1**  
  Same as `deepseek-reasoner`. Its Just provided by Azure . <br/>
  Usage: `Model: deepseek.AzureDeepSeekR1`

<details open>
<summary> Chat </summary>

### Example for chatting with deepseek

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	deepseek "github.com/cohesion-org/deepseek-go"
	constants "github.com/cohesion-org/deepseek-go/constants"
)

func main() {
	// Set up the Deepseek client
    client := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))

	// Create a chat completion request
	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: constants.ChatMessageRoleSystem, Content: "Answer every question using slang."},
			{Role: constants.ChatMessageRoleUser, Content: "Which is the tallest mountain in the world?"},
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

Using Azure:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	deepseek "github.com/cohesion-org/deepseek-go"
	constants "github.com/cohesion-org/deepseek-go/constants"
)

func main() {
	// Set up the Deepseek client
    client := deepseek.NewAzureClient(os.Getenv("YOUR_TOKEN"))

	// Create a chat completion request
	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.AzureDeepSeekR1,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: constants.ChatMessageRoleSystem, Content: "Answer every question using slang."},
			{Role: constants.ChatMessageRoleUser, Content: "Which is the tallest mountain in the world?"},
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

Save this code to a file (e.g., `main.go`), and run it:
```sh
go run main.go
```
</details>

<details >
	<summary> Sending other params like Temp, Stop </summary>
	<strong> You just need to extend the ChatCompletionMessage with the supported parameters. </strong>

```go
	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: constants.ChatMessageRoleUser, Content: "What is the meaning of deepseek"},
			{Role: constants.ChatMessageRoleSystem, Content: "Answer every question using slang"},
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
	"github.com/cohesion-org/deepseek-go/constants"
)

func MultiChat() {
	client := deepseek.NewClient("DEEPSEEK_API_KEY")
	ctx := context.Background()

	messages := []deepseek.ChatCompletionMessage{{
		Role:    constants.ChatMessageRoleUser,
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
		Role:    constants.ChatMessageRoleUser,
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
	constants "github.com/cohesion-org/deepseek-go/constants"
)

func main() {
	client := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))
	request := &deepseek.StreamChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: constants.ChatMessageRoleUser, Content: "Just testing if the streaming feature is working or not!"},
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

---
## Getting a Deepseek Key

To use the Deepseek API, you need an API key. You can obtain one by signing up on the [Deepseek website](https://platform.deepseek.com/api_keys)

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
<!-- - `client_test.go`: Client configuration and error handling -->
- `chat_test.go`: Chat completion functionality 
- `chat_stream_test.go`: Chat streaming functionality
- `models_test.go`: Model listing and retrieval
- `balance_test.go`: Account balance operations
- `tokens_test.go`: Token estimation utilities
<!-- - `errors_test.go`: Tests the error handler -->
### Handlers Package
- `handlers/requestHandler_test.go`: Tests for the request handler
- `handlers/responseHandler_test.go`: Tests for the response handler

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

---

Feel free to contribute, open issues, or submit PRs to help improve Deepseek-Go! Let us know if you encounter any issues.

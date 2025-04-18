# Using Ollama with deepseek-go
This example demonstrates how to use Ollama with deepseek-go to generate chat completions. Ollama's API functionality is extended and standardized to match the OpenAI format.

## Prerequisites

- [Ollama](https://ollama.ai/) installed and running locally
- A compatible model downloaded (e.g., llama2, codellama, mistral)

## Usage Example

### Chat

```go
package main

import (
    "fmt"
    "github.com/cohesion-org/deepseek-go"
    "github.com/cohesion-org/deepseek-go/constants"
)

func main() {
    req := &deepseek.ChatCompletionRequest{
        Model: "llama2:latest",
        Messages: []deepseek.ChatCompletionMessage{
            {
                Role: constants.ChatMessageRoleUser,
                Content: "What is the capital of France?",
            },
        },
    }
    
    resp, err := deepseek.CreateOllamaCompletion(req)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
}
```

### Stream

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "io"
    "github.com/cohesion-org/deepseek-go"
    "github.com/cohesion-org/deepseek-go/constants"
)

func main() {
    ctx := context.Background()
    messages := []deepseek.ChatCompletionMessage{
        {
            Role:    constants.ChatMessageRoleUser,
            Content: "What is artificial intelligence?",
        },
    }

    req := &deepseek.ChatCompletionRequest{
        Stream: true,
        Model:  "llava:latest",
        Messages: messages,
    }

    stream, err := deepseek.CreateOllamaChatCompletionStream(ctx, req)
    if err != nil {
        fmt.Printf("Error creating stream: %v\n", err)
        return
    }
    defer stream.Close()

    for {
        response, err := stream.Recv()
        if errors.Is(err, io.EOF) {
            break
        }
        if err != nil {
            fmt.Printf("Stream error: %v\n", err)
            return
        }

        for _, choice := range response.Choices {
            fmt.Print(choice.Delta.Content)
        }
    }
}
```

## Features

✅ **Supported**
- Basic chat completion
- Multiple messages in conversation
- Usage statistics
- Model selection

❌ **Current Limitations**
- No image handling
- No function and tool calling

## Roadmap

- [✅] Streaming support
- [ ] Image handling
- [ ] Enhanced function calling
- [ ] Improved error handling
- [ ] Additional configuration options

## Requirements

- Ollama installed and running
- Proper environment configuration
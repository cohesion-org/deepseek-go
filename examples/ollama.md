# Using Ollama with deepseek-go
This example demonstrates how to use Ollama with deepseek-go to generate chat completions. Ollama's API functionality is extended and standardized to match the OpenAI format.

## Prerequisites

- [Ollama](https://ollama.ai/) installed and running locally
- A compatible model downloaded (e.g., llama2, codellama, mistral)

## Usage Example

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

## Features

✅ **Supported**
- Basic chat completion
- Multiple messages in conversation
- Usage statistics
- Model selection

❌ **Current Limitations**
- No streaming support
- No image handling
- Limited function calling
- Response format differences

## Roadmap

- [ ] Streaming support
- [ ] Image handling
- [ ] Enhanced function calling
- [ ] Improved error handling
- [ ] Additional configuration options

## Requirements

- Ollama installed and running
- Proper environment configuration
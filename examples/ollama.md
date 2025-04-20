# Using Ollama with deepseek-go
This example demonstrates how to use Ollama with deepseek-go to generate chat completions. Ollama's API functionality is extended and standardized to match the OpenAI format.

## Prerequisites

- [Ollama](https://ollama.ai/) installed and running locally
- A compatible model downloaded (e.g., llama2, codellama, mistral)

## Usage Example

- [Chat](#chat)
- [Streaming](#stream)
- [Image](#image)
- [Stream with Image](#stream-with-image)

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
    
    resp, err := deepseek.CreateOllamaChatCompletion(req)
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
    req := &deepseek.StreamChatCompletionRequest{
        Model: "llava:latest",
        Messages: []deepseek.ChatCompletionMessage{
            {
                Role:    constants.ChatMessageRoleUser,
                Content: "What is artificial intelligence?",
            },
        },
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

### Image 

```go
package main

import (
    "fmt"
    "github.com/cohesion-org/deepseek-go"
)

func main() {
    // Convert image to base64
    imgData, err := deepseek.ImageToBase64("path/to/your/image.png")
    if err != nil {
        fmt.Printf("Error converting image: %v\n", err)
        return
    }

    // Create request with image
    req := &deepseek.ChatCompletionRequestWithImage{
        Model: "llava:latest",
        Messages: []deepseek.ChatCompletionMessageWithImage{
            deepseek.NewImageMessage("user", "What is this image about?", imgData),
        },
    }

    // Send request and get response
    resp, err := deepseek.CreateOllamaChatCompletionWithImage(req)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
}
```

### Stream with Image

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "io"
    "github.com/cohesion-org/deepseek-go"
)

func main() {
    // Convert image to base64
    imgData, err := deepseek.ImageToBase64("path/to/your/image.png")
    if err != nil {
        fmt.Printf("Error converting image: %v\n", err)
        return
    }

    // Create request with image
    req := &deepseek.StreamChatCompletionRequestWithImage{
        Model: "llava:latest",
        Messages: []deepseek.ChatCompletionMessageWithImage{
            deepseek.NewImageMessage("user", "What is this image about?", imgData),
        },
    }
    req.Stream = true

    // Create stream
    ctx := context.Background()
    stream, err := deepseek.CreateOllamaChatCompletionStreamWithImage(ctx, req)
    if err != nil {
        fmt.Printf("Error creating stream: %v\n", err)
        return
    }
    defer stream.Close()

    // Read from stream
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
- [✅] Image handling
- [ ] Enhanced function calling
- [ ] Improved error handling
- [ ] Additional configuration options

## Requirements

- Ollama installed and running
- Proper environment configuration
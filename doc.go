// Package deepseek is a Go-based API client for the Deepseek platform, providing a clean and type-safe interface
// to interact with Deepseek's AI features and compatible providers.
//
// Features:
//
// - Chat Completion: Send chat messages and receive responses from Deepseek's AI models with support for both
//   regular and streaming responses.
//
// - Model Support: Access to official Deepseek models (DeepSeekChat, DeepSeekCoder, DeepSeekReasoner) and
//   external providers like OpenRouter, Azure, and Ollama.
//
// - Advanced Features: Support for Fill-In-the-Middle (FIM) completions, JSON mode for structured outputs,
//   function calling, and image processing capabilities.
//
// - Token Management: Track token usage and estimate token counts for requests.
//
// - Flexible Configuration: Customize client behavior with options for base URLs, timeouts, and HTTP clients.
//
// - Balance Tracking: Check account balance and usage information.
//
// The package is designed with a modular architecture, separating request building, sending, and response handling
// into reusable components. This makes it easy to extend and maintain while providing a simple interface for users.
//
// Basic usage example:
//
//	client := deepseek.NewClient("YOUR_API_KEY")
//	request := &deepseek.ChatCompletionRequest{
//		Model: deepseek.DeepSeekChat,
//		Messages: []deepseek.ChatCompletionMessage{
//			{Role: deepseek.ChatMessageRoleUser, Content: "Hello, how are you?"},
//		},
//	}
//	response, err := client.CreateChatCompletion(context.Background(), request)
//
// For more examples and detailed documentation, visit the project repository:
// https://github.com/cohesion-org/deepseek-go
package deepseek
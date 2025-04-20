# Examples

This directory contains various examples demonstrating different capabilities of the deepseek-go library.

| #  | Example Name                                   | Description |
|----|----------------------------------------------|-------------|
| 0  | **[Using External Providers](00_external_providers/chat.go)** | Supports external providers through `baseURL` extension. Missing model constants can be reported via issues or pull requests. |
| 1  | **[Basic Chat Example](01_chat/chat.go)**  | Demonstrates basic chat functionality. |
| 2  | **[Chat with Streaming](02_chat_stream/chat_stream.go)** | Implements streaming chat responses, including `ReasoningContent` with R1. |
| 3  | **[Fill-in-Middle (FIM)](03_fim/fim.go)** | Example of fill-in-middle completion with streaming support. |
| 4  | **[JSON Mode](04_json_mode/json_mode.go)** | Demonstrates JSON mode for structured responses. This is a client-specific feature. |
| 5  | **[Multi-Chat](05_multi_chat/multi_chat.go)** | Example of handling multiple concurrent chat sessions. |
| 6  | **[Bad Multi-Chat](06_bad_multi_chat/bad_multi_chat.go)** | Demonstrates incorrect handling of multiple chats (for educational purposes). |
| 7  | **[Balance Example](07_balance/balance.go)** | Shows balance-related functionality. |
| 8  | **[Client with Options](08_newClientWithOptions/newClientWithOptions.go)** | Demonstrates creating a client with custom options. |
| 9  | **[Prefix Completion](09_prefix_completion/prefix_completion.go)** | Example of prefix-based completion. |
| 10 | **[Token Usage Estimation](10_token_usage/token_usage.go)** | Demonstrates how to estimate and track token usage for requests (based on Deepseek’s documentation). |
| 11 | **[List Supported Models](11_list_models/list_models.go)** | Shows how to list all supported models through the Deepseek API. |
| 12 | **[Function Calling](12_function_calling/function_calling.go)** | Demonstrates function calling capabilities. |
| 13 | **[OpenRouter Images](13_openrouter_images/openrouter_images.go)** | Example usage with OpenRouter images. |

# Ollama
Ollama docs are located [here](./ollama.md).
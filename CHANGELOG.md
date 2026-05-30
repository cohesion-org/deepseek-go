# Changelog

All notable changes to this project will be documented in this file.

## [1.4.0] - 2026-05-30

### Added
- Anthropic-compatible API support via AnthropicClient
- V4 model constants: DeepSeekV4Flash, DeepSeekV4Pro
- reasoning_effort parameter on all request structs
- user_id parameter on all request structs
- Strict function calling with automatic /beta routing
- json_schema type in ResponseFormat
- CompletionTokensDetails in non-streaming Usage
- Runtime deprecation warnings for old model names
- New examples: reasoning_effort, user_id, strict_tools, anthropic

### Changed
- Go version bumped to 1.26.0
- Ollama dependency updated from v0.6.5 to v0.24.0
- FIM float types normalized to float32 for consistency
- All tests, examples, and mocks migrated to V4 model constants

### Deprecated
- DeepSeekChat, DeepSeekCoder, DeepSeekReasoner constants (sunset 2026/07/24)

### Fixed
- tool_choice added to StreamChatCompletionRequest
- 4 Dependabot security alerts resolved (Ollama CVEs)

## [1.3.4] - 2026-03-23

### Added
- Compatibility for DeepSeek December 1, 2025 V3.2 API update
- Support for current thinking request format
- Preserved reasoning_content and tool_calls in tool-call loops

### Changed
- Tests migrated to mock server (no longer require DEEPSEEK_API_KEY)
- CI reworked to run offline mock-backed suite by default
- Live API tests moved to explicit integration lane
- GitHub Actions workflows upgraded for Node 24 transition

## [1.3.3] - 2026-02-12

### Added
- DeepThink support and extraFields configuration
- Improved FunctionCalling to handle user input and tool calls

## [1.3.2] - 2025-05-28

### Added
- enable_thinking tag for Qwen3 API
- doc.go and improved file-level comments for pkg.go.dev

## [1.3.1] - 2025-04-28

### Added
- Environment variable fallback for empty authentication
- Logprobs support
- Reasoning support with OpenRouter

### Changed
- Improved error handling and error messages
- Better error handling in NewClient when API key is missing

## [1.3.0] - 2025-04-20

### Added
- Ollama integration: CreateOllamaChatCompletion and CreateOllamaChatCompletionStream
- New OLLAMA model constants
- Ollama documentation and examples

### Changed
- Go version bumped to 1.24.0
- Added github.com/ollama/ollama dependency

### Removed
- Outdated comments from requestHandler.go

## [1.2.12] - 2025-04-17

### Fixed
- url.ParseRequestURI() incorrectly handling absolute file paths as valid URLs
- Local file paths no longer mistakenly treated as web URLs

## [1.2.1] - 2025-02-14

### Added
- Chat Prefix Completion (Beta) via beta endpoint
- Prefix guidance for assistant messages
- New example demonstrating prefix completion with JSON mode

## [1.2.0] - 2025-02-13

### Breaking
- Renamed Tools struct to Tool for clarity
- Renamed Parameters struct to FunctionParameters

### Added
- ToolChoice feature for flexible tool selection (auto or manual)
- FIM streaming support (FIMRecv, FIMClose)

### Changed
- Significant refactor of FIMStreamCompletionRequest, FIMStreamChoice, FIMStreamCompletionResponse

## [1.1.2] - 2025-02-11

### Added
- FIM Completion (Beta) - CreateFIMCompletion and CreateFIMCompletionStream
- StreamOptions for StreamChatCompletionRequest
- ToolChoice for ChatCompletionRequest
- ToolCall to Message object

### Changed
- Improved error handling when initializing client with different baseURL
- Enhanced documentation and testing coverage

## [1.1.1] - 2025-02-10

### Added
- Client Base URL extension for custom providers (OpenRouter, Azure)
- JSON Mode with automatic JSON extraction and parsing into Go structs
- NewClientWithOptions for flexible client creation with functional options
- Constants for common OpenRouter models

### Changed
- Improved error handling for request sending and processing

## [1.1.0] - 2025-02-10

### Notes
- Version 1.1.0 was released and subsequently retracted due to versioning issues with pkg.dev. The feature set was merged into v1.1.1.

## [1.0.0] - 2025-02-05

### Notes
- Version 1.0.0 was released and subsequently retracted due to versioning issues with pkg.dev. Users should upgrade to v1.1.1 or later.

## [0.2.1] - 2025-02-05

### Fixed
- Retracted premature versions v1.1.0 and v1.0.1 from the package

## [0.2.0] - 2025-02-04

### Added
- Base URL extension for external providers (Azure, OpenRouter, and others)
- Timeout configuration via DEEPSEEK_TIMEOUT environment variable
- Token estimation function (EstimateTokensFromMessages)

### Changed
- Improved error handling for HTML error responses from external servers
- Expanded testing coverage across most files

## [0.1.1] - 2025-01-30

### Fixed
- Streaming API 404 error caused by Deepseek API changes to stream endpoint

## [0.1.0] - 2025-01-28

### Added
- Initial release of deepseek-go SDK
- Mappers utility (mappers.go) for multi-round conversations
- MapMessageToChatCompletionMessage function for converting responses to request messages
- Basic chat completion and streaming support

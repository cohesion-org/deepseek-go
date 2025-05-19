package deepseek

import (
	"errors"

	"github.com/cohesion-org/deepseek-go/constants"
)

const (
	// ChatMessageRoleSystem is the role of a system message
	ChatMessageRoleSystem = constants.ChatMessageRoleSystem
	// ChatMessageRoleUser is the role of a user message
	ChatMessageRoleUser = constants.ChatMessageRoleUser
	// ChatMessageRoleAssistant is the role of an assistant message
	ChatMessageRoleAssistant = constants.ChatMessageRoleAssistant
	// ChatMessageRoleTool is the role of a tool message
	ChatMessageRoleTool = constants.ChatMessageRoleTool
)

var (
	// ErrChatCompletionStreamNotSupported is returned when streaming is not supported with the method.
	ErrChatCompletionStreamNotSupported = errors.New("streaming is not supported with this method")
	// ErrChatCompletionRequestNil is returned when the request is nil.
	ErrUnexpectedResponseFormat = errors.New("unexpected response format")
)

// ChatCompletionMessage represents a single message in a chat completion conversation.
type ChatCompletionMessage struct {
	Role             string     `json:"role"`                        // The role of the message sender, e.g., "user", "assistant", "system".
	Content          string     `json:"content"`                     // The content of the message.
	Prefix           bool       `json:"prefix,omitempty"`            // The prefix of the message (optional) for Chat Prefix Completion [Beta Feature].
	ReasoningContent string     `json:"reasoning_content,omitempty"` // The reasoning content of the message (optional) when using the reasoner model with Chat Prefix Completion. When using this feature, the Prefix parameter must be set to true.
	ToolCallID       string     `json:"tool_call_id,omitempty"`      // Tool call that this message is responding to.
	ToolCalls        []ToolCall `json:"tool_calls,omitempty"`        // Optional tool calls.

}

// FunctionParameters defines the parameters for a function.
type FunctionParameters struct {
	Type       string                 `json:"type"`                 // The type of the parameters, e.g., "object" (required).
	Properties map[string]interface{} `json:"properties,omitempty"` // The properties of the parameters (optional).
	Required   []string               `json:"required,omitempty"`   // A list of required parameter names (optional).
}

// Function defines the structure of a function tool.
type Function struct {
	Name        string              `json:"name"`                 // The name of the function (required).
	Description string              `json:"description"`          // A description of the function (required).
	Parameters  *FunctionParameters `json:"parameters,omitempty"` // The parameters of the function (optional).
}

// Tool defines the structure for a tool.
type Tool struct {
	Type     string   `json:"type"`     // The type of the tool, e.g., "function" (required).
	Function Function `json:"function"` // The function details (required).
}

// ToolChoice defines the structure for a tool choice.
type ToolChoice struct {
	Type     string             `json:"type"`               // The type of the tool, e.g., "function" (required).
	Function ToolChoiceFunction `json:"function,omitempty"` // The function details (optional, but required if type is "function").
}

// ToolChoiceFunction defines the function details within ToolChoice.
type ToolChoiceFunction struct {
	Name string `json:"name"` // The name of the function to call (required).
}

// ResponseFormat defines the structure for the response format.
type ResponseFormat struct {
	Type string `json:"type"` // The desired response format, either "text" or "json_object".
}

// ChatCompletionRequest defines the structure for a chat completion request.
type ChatCompletionRequest struct {
	Model            string                  `json:"model"`                       // The ID of the model to use (required).
	Messages         []ChatCompletionMessage `json:"messages"`                    // A list of messages comprising the conversation (required).
	FrequencyPenalty float32                 `json:"frequency_penalty,omitempty"` // Penalty for new tokens based on their frequency in the text so far (optional).
	MaxTokens        int                     `json:"max_tokens,omitempty"`        // The maximum number of tokens to generate in the chat completion (optional).
	PresencePenalty  float32                 `json:"presence_penalty,omitempty"`  // Penalty for new tokens based on their presence in the text so far (optional).
	Temperature      float32                 `json:"temperature,omitempty"`       // The sampling temperature, between 0 and 2 (optional).
	TopP             float32                 `json:"top_p,omitempty"`             // The nucleus sampling parameter, between 0 and 1 (optional).
	ResponseFormat   *ResponseFormat         `json:"response_format,omitempty"`   // The desired response format (optional).
	Stop             []string                `json:"stop,omitempty"`              // A list of sequences where the model should stop generating further tokens (optional).
	Tools            []Tool                  `json:"tools,omitempty"`             // A list of tools the model may use (optional).
	ToolChoice       interface{}             `json:"tool_choice,omitempty"`       // Controls which (if any) tool is called by the model (optional).
	LogProbs         bool                    `json:"logprobs,omitempty"`          // Whether to return log probabilities of the most likely tokens (optional).
	TopLogProbs      int                     `json:"top_logprobs,omitempty"`      // The number of top most likely tokens to return log probabilities for (optional).
	JSONMode         bool                    `json:"json,omitempty"`              // [deepseek-go feature] Optional: Enable JSON mode. If you're using the JSON mode, please mention "json" anywhere in your prompt, and also include the JSON schema in the request.
	EnableThinking   bool                    `json:"enable_thinking,omitempty"`   // Optional: Enable thinking (for qwen3 api)
}

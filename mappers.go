package deepseek

import (
	"errors"
	"fmt"

	"github.com/cohesion-org/deepseek-go/constants"
)

var validRoles = map[string]bool{
	constants.ChatMessageRoleUser:      true,
	constants.ChatMessageRoleAssistant: true,
	constants.ChatMessageRoleSystem:    true,
}

// MapMessageToChatCompletionMessage maps a Message to a ChatCompletionMessage
func MapMessageToChatCompletionMessage(m Message) (ChatCompletionMessage, error) {
	if m.Role == "" {
		return ChatCompletionMessage{}, errors.New("message role cannot be empty")
	}

	if m.Content == "" && m.ReasoningContent == "" && len(m.ToolCalls) == 0 {
		return ChatCompletionMessage{}, errors.New("message content, reasoning content, and tool calls cannot all be empty")
	}
	if !validRoles[m.Role] {
		return ChatCompletionMessage{}, fmt.Errorf("invalid role: %s. Valid roles are can be found in official deepseek documentation", m.Role)
	}

	return ChatCompletionMessage{
		Role:             m.Role,
		Content:          m.Content,
		ReasoningContent: m.ReasoningContent,
		ToolCalls:        m.ToolCalls,
	}, nil
}

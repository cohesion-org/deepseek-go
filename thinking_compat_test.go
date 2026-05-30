package deepseek_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cohesion-org/deepseek-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateChatCompletion_EnableThinkingUsesTopLevelThinking(t *testing.T) {
	var captured map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			_ = r.Body.Close()
		}()

		err := json.NewDecoder(r.Body).Decode(&captured)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(`{"id":"chatcmpl-1","object":"chat.completion","created":1,"model":"deepseek-chat","choices":[{"index":0,"message":{"role":"assistant","content":"ok"},"finish_reason":"stop"}],"usage":{"prompt_tokens":1,"completion_tokens":1,"total_tokens":2}}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	client, err := deepseek.NewClientWithOptions("token", deepseek.WithBaseURL(server.URL+"/"))
	require.NoError(t, err)

	resp, err := client.CreateChatCompletion(context.Background(), &deepseek.ChatCompletionRequest{
		Model:          deepseek.DeepSeekV4Flash,
		EnableThinking: true,
		ExtraFields: map[string]interface{}{
			"foo": "bar",
		},
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleUser, Content: "hello"},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.NotContains(t, captured, "enable_thinking")
	assert.NotContains(t, captured, "extra_fields")
	assert.Equal(t, "bar", captured["foo"])

	thinking, ok := captured["thinking"].(map[string]interface{})
	require.True(t, ok, "expected top-level thinking object")
	assert.Equal(t, "enabled", thinking["type"])
}

func TestCreateChatCompletionStream_EnableThinkingUsesTopLevelThinking(t *testing.T) {
	var captured map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			_ = r.Body.Close()
		}()

		err := json.NewDecoder(r.Body).Decode(&captured)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "text/event-stream")
		_, err = io.WriteString(w, "data: {\"id\":\"chatcmpl-1\",\"object\":\"chat.completion.chunk\",\"created\":1,\"model\":\"deepseek-chat\",\"choices\":[{\"index\":0,\"delta\":{\"content\":\"ok\"},\"finish_reason\":\"stop\"}]}\n\n")
		require.NoError(t, err)
		_, err = io.WriteString(w, "data: [DONE]\n\n")
		require.NoError(t, err)
	}))
	defer server.Close()

	client, err := deepseek.NewClientWithOptions("token", deepseek.WithBaseURL(server.URL+"/"))
	require.NoError(t, err)

	stream, err := client.CreateChatCompletionStream(context.Background(), &deepseek.StreamChatCompletionRequest{
		Model:          deepseek.DeepSeekV4Flash,
		EnableThinking: true,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleUser, Content: "hello"},
		},
	})
	require.NoError(t, err)
	defer func() {
		_ = stream.Close()
	}()

	chunk, err := stream.Recv()
	require.NoError(t, err)
	require.NotNil(t, chunk)

	thinking, ok := captured["thinking"].(map[string]interface{})
	require.True(t, ok, "expected top-level thinking object")
	assert.Equal(t, "enabled", thinking["type"])
	assert.NotContains(t, captured, "enable_thinking")
	assert.NotContains(t, captured, "extra_fields")
}

func TestMapMessageToChatCompletionMessage_PreservesReasoningAndToolCalls(t *testing.T) {
	message := deepseek.Message{
		Role:             deepseek.ChatMessageRoleAssistant,
		Content:          "",
		ReasoningContent: "Need to call a tool first.",
		ToolCalls: []deepseek.ToolCall{
			{
				ID:   "call_123",
				Type: "function",
				Function: deepseek.ToolCallFunction{
					Name:      "get_time",
					Arguments: "{}",
				},
			},
		},
	}

	mapped, err := deepseek.MapMessageToChatCompletionMessage(message)
	require.NoError(t, err)

	assert.Equal(t, message.Role, mapped.Role)
	assert.Equal(t, message.Content, mapped.Content)
	assert.Equal(t, message.ReasoningContent, mapped.ReasoningContent)
	assert.Equal(t, message.ToolCalls, mapped.ToolCalls)
}

func TestCreateChatCompletion_ToolContinuationPreservesReasoningContent(t *testing.T) {
	var captured map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			_ = r.Body.Close()
		}()

		err := json.NewDecoder(r.Body).Decode(&captured)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(`{"id":"chatcmpl-2","object":"chat.completion","created":2,"model":"deepseek-chat","choices":[{"index":0,"message":{"role":"assistant","content":"final answer"},"finish_reason":"stop"}],"usage":{"prompt_tokens":2,"completion_tokens":1,"total_tokens":3}}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	client, err := deepseek.NewClientWithOptions("token", deepseek.WithBaseURL(server.URL+"/"))
	require.NoError(t, err)

	assistantMessage, err := deepseek.MapMessageToChatCompletionMessage(deepseek.Message{
		Role:             deepseek.ChatMessageRoleAssistant,
		Content:          "",
		ReasoningContent: "I should call a tool.",
		ToolCalls: []deepseek.ToolCall{
			{
				ID:   "call_123",
				Type: "function",
				Function: deepseek.ToolCallFunction{
					Name:      "get_time",
					Arguments: "{}",
				},
			},
		},
	})
	require.NoError(t, err)

	_, err = client.CreateChatCompletion(context.Background(), &deepseek.ChatCompletionRequest{
		Model:    deepseek.DeepSeekV4Flash,
		Thinking: &deepseek.ThinkingConfig{Type: "enabled"},
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleUser, Content: "what time is it?"},
			assistantMessage,
			{Role: deepseek.ChatMessageRoleTool, ToolCallID: "call_123", Content: "2026-03-23T12:00:00Z"},
		},
	})
	require.NoError(t, err)

	messages, ok := captured["messages"].([]interface{})
	require.True(t, ok)
	require.Len(t, messages, 3)

	assistantPayload, ok := messages[1].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "assistant", assistantPayload["role"])
	assert.Equal(t, "I should call a tool.", assistantPayload["reasoning_content"])

	toolCalls, ok := assistantPayload["tool_calls"].([]interface{})
	require.True(t, ok)
	require.Len(t, toolCalls, 1)
}

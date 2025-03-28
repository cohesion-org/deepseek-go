package deepseek_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/constants"
	"github.com/cohesion-org/deepseek-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateChatCompletionStream(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	client := deepseek.NewClient(config.APIKey)

	ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
	defer cancel()

	stream, _, err := client.CreateChatCompletionStream(ctx, &deepseek.StreamChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{
				Role:    constants.ChatMessageRoleUser,
				Content: "Tell me a joke about artificial intelligence",
			},
		},
	})
	require.NoError(t, err)
	defer stream.Close()

	var contentBuffer string
	var receivedFinishReason bool

	for {
		resp, err := stream.Recv()
		if err != nil {
			break
		}

		if len(resp.Choices) > 0 {
			contentBuffer += resp.Choices[0].Delta.Content
			if resp.Choices[0].FinishReason != "" {
				receivedFinishReason = true
			}
		}
	}

	assert.True(t, len(contentBuffer) > 0, "should receive streamed content")
	assert.Contains(t, contentBuffer, "AI", "response should be related to AI")
	assert.True(t, receivedFinishReason, "should receive finish reason")
}

func TestStreamingMultiChat(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	client := deepseek.NewClient(config.APIKey)
	ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
	defer cancel()

	var messages []deepseek.ChatCompletionMessage

	t.Run("InitialQuestion", func(t *testing.T) {
		// First question about highest mountain
		messages = []deepseek.ChatCompletionMessage{{
			Role:    constants.ChatMessageRoleUser,
			Content: "What's the highest mountain in the world? Respond with only the name.",
		}}

		response, err := streamChatCompletion(t, ctx, client, messages)
		require.NoError(t, err, "initial streaming should succeed")
		assert.NotEmpty(t, response, "response should not be empty")
		assert.Contains(t, response, "Everest", "should identify Mount Everest")

		// Append assistant response
		messages = append(messages, deepseek.ChatCompletionMessage{
			Role:    constants.ChatMessageRoleAssistant,
			Content: response,
		})
	})

	t.Run("FollowUpQuestion", func(t *testing.T) {
		// Add follow-up question
		messages = append(messages, deepseek.ChatCompletionMessage{
			Role:    constants.ChatMessageRoleUser,
			Content: "What's the second highest? Respond with only the name.",
		})

		response, err := streamChatCompletion(t, ctx, client, messages)
		require.NoError(t, err, "follow-up streaming should succeed")
		assert.NotEmpty(t, response, "response should not be empty")
		assert.Contains(t, response, "K2", "should identify K2 mountain")

		// Validate conversation history
		require.Len(t, messages, 3, "should have 3 messages in history")
		assert.Equal(t, constants.ChatMessageRoleUser, messages[0].Role)
		assert.Equal(t, constants.ChatMessageRoleAssistant, messages[1].Role)
		assert.Equal(t, constants.ChatMessageRoleUser, messages[2].Role)
	})
}

func streamChatCompletion(
	t *testing.T,
	ctx context.Context,
	client *deepseek.Client,
	messages []deepseek.ChatCompletionMessage,
) (string, error) {
	req := &deepseek.StreamChatCompletionRequest{
		Model:    deepseek.DeepSeekChat,
		Messages: messages,
		Stream:   true,
	}

	stream, _, err := client.CreateChatCompletionStream(ctx, req)
	require.NoError(t, err, "should create stream without error")
	defer func() {
		if cerr := stream.Close(); cerr != nil {
			t.Errorf("error closing stream: %v", cerr)
		}
	}()

	var (
		fullMessage  string
		receivedRole bool
	)

	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err, "stream should not error")

		// Validate response structure
		assert.Equal(t, "chat.completion.chunk", resp.Object)
		assert.Equal(t, deepseek.DeepSeekChat, resp.Model)
		assert.NotEmpty(t, resp.ID)
		assert.NotZero(t, resp.Created)

		if len(resp.Choices) > 0 {
			chunk := resp.Choices[0]
			if chunk.Delta.Role != "" {
				receivedRole = true
				assert.Equal(t, constants.ChatMessageRoleAssistant, chunk.Delta.Role)
			}
			fullMessage += chunk.Delta.Content
		}
	}

	assert.True(t, receivedRole, "should receive assistant role")
	assert.NotEmpty(t, fullMessage, "should accumulate message content")
	return fullMessage, nil
}

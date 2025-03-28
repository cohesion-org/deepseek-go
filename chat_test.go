package deepseek_test

import (
	"context"
	"testing"

	"github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/constants"
	"github.com/cohesion-org/deepseek-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateChatCompletion(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	client := deepseek.NewClient(config.APIKey)

	tests := []struct {
		name        string
		req         *deepseek.ChatCompletionRequest
		wantErr     bool
		validateRes func(t *testing.T, res *deepseek.ChatCompletionResponse)
	}{
		{
			name: "basic completion",
			req: &deepseek.ChatCompletionRequest{
				Model: deepseek.DeepSeekChat,
				Messages: []deepseek.ChatCompletionMessage{
					{Role: constants.ChatMessageRoleUser, Content: "Say hello!"},
				},
			},
			wantErr: false,
			validateRes: func(t *testing.T, res *deepseek.ChatCompletionResponse) {
				assert.NotEmpty(t, res.Choices[0].Message.Content)
			},
		},
		{
			name: "empty messages",
			req: &deepseek.ChatCompletionRequest{
				Model:    deepseek.DeepSeekChat,
				Messages: []deepseek.ChatCompletionMessage{},
			},
			wantErr: true,
		},
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
			defer cancel()

			resp, _, err := client.CreateChatCompletion(ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, resp)

			// Validate common response structure
			assert.NotEmpty(t, resp.ID)
			assert.NotEmpty(t, resp.Created)
			assert.Equal(t, "chat.completion", resp.Object)
			assert.Equal(t, tt.req.Model, resp.Model)
			assert.NotEmpty(t, resp.Choices)
			assert.NotNil(t, resp.Usage)

			// Validate specific test case expectations
			if tt.validateRes != nil {
				tt.validateRes(t, resp)
			}
		})
	}
}

func TestMultiChatConversation(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	client := deepseek.NewClient(config.APIKey)
	ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
	defer cancel()

	// Initial message setup
	messages := []deepseek.ChatCompletionMessage{{
		Role:    constants.ChatMessageRoleUser,
		Content: "Who is the current president of the United States? Respond with only the last name.",
	}}

	// First round of conversation
	t.Run("FirstResponse", func(t *testing.T) {
		req := &deepseek.ChatCompletionRequest{
			Model:    deepseek.DeepSeekChat,
			Messages: messages,
		}

		resp, _, err := client.CreateChatCompletion(ctx, req)
		require.NoError(t, err, "initial request should succeed")
		require.NotNil(t, resp, "response should not be nil")
		require.NotEmpty(t, resp.Choices, "response should contain choices")

		// Validate response structure
		assert.Equal(t, "chat.completion", resp.Object)
		assert.Equal(t, deepseek.DeepSeekChat, resp.Model)
		assert.NotZero(t, resp.Created)
		assert.NotZero(t, resp.Usage.TotalTokens)

		// Convert and append response
		responseMessage, err := deepseek.MapMessageToChatCompletionMessage(resp.Choices[0].Message)
		require.NoError(t, err, "message mapping should succeed")
		messages = append(messages, responseMessage)
	})

	// Add follow-up question
	messages = append(messages, deepseek.ChatCompletionMessage{
		Role:    constants.ChatMessageRoleUser,
		Content: "Who was the immediate predecessor? Respond with only the last name.",
	})

	// Second round of conversation
	t.Run("SecondResponse", func(t *testing.T) {
		req := &deepseek.ChatCompletionRequest{
			Model:    deepseek.DeepSeekChat,
			Messages: messages,
		}

		resp, _, err := client.CreateChatCompletion(ctx, req)
		require.NoError(t, err, "follow-up request should succeed")
		require.NotNil(t, resp, "response should not be nil")
		require.NotEmpty(t, resp.Choices, "response should contain choices")

		// Validate conversation continuity
		assert.Greater(t, len(messages), 2, "should have conversation history")

		// Check response content
		content := resp.Choices[0].Message.Content
		assert.NotEmpty(t, content, "response should have content")
		assert.NotEqual(t, messages[1].Content, content, "responses should be different")

		// Validate response structure
		assert.Equal(t, "chat.completion", resp.Object)
		assert.Equal(t, deepseek.DeepSeekChat, resp.Model)
		assert.NotZero(t, resp.Created)
		assert.NotZero(t, resp.Usage.TotalTokens)
	})
}

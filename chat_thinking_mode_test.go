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

func TestThinkingModeCompletion(t *testing.T) {
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
			name: "thinking mode with deepseek-reasoner model",
			req: &deepseek.ChatCompletionRequest{
				Model: deepseek.DeepSeekReasoner,
				Messages: []deepseek.ChatCompletionMessage{
					{Role: constants.ChatMessageRoleUser, Content: "9.11 and 9.8, which is greater?"},
				},
			},
			wantErr: false,
			validateRes: func(t *testing.T, res *deepseek.ChatCompletionResponse) {
				assert.NotEmpty(t, res.Choices[0].Message.Content)
				assert.NotEmpty(t, res.Choices[0].Message.ReasoningContent, "reasoning content should be present in thinking mode")
			},
		},
		{
			name: "thinking mode with deepseek-chat and thinking parameter",
			req: &deepseek.ChatCompletionRequest{
				Model: deepseek.DeepSeekChat,
				Messages: []deepseek.ChatCompletionMessage{
					{Role: constants.ChatMessageRoleUser, Content: "How many Rs are there in the word 'strawberry'?"},
				},
				Thinking: &deepseek.ThinkingMode{Type: "enabled"},
			},
			wantErr: false,
			validateRes: func(t *testing.T, res *deepseek.ChatCompletionResponse) {
				assert.NotEmpty(t, res.Choices[0].Message.Content)
				assert.NotEmpty(t, res.Choices[0].Message.ReasoningContent, "reasoning content should be present when thinking is enabled")
				// Note: API returns deepseek-reasoner when thinking is enabled, even if deepseek-chat was requested
				assert.Equal(t, deepseek.DeepSeekReasoner, res.Model, "model should be deepseek-reasoner when thinking is enabled")
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

			resp, err := client.CreateChatCompletion(ctx, tt.req)
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
			assert.NotEmpty(t, resp.Choices)
			assert.NotNil(t, resp.Usage)

			// Validate specific test case expectations
			if tt.validateRes != nil {
				tt.validateRes(t, resp)
			}
		})
	}
}

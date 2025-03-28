package deepseek_test

import (
	"context"
	"os"
	"testing"

	"github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/constants"
	"github.com/cohesion-org/deepseek-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatPrefixCompletion(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	client := deepseek.NewClient(config.APIKey, "https://api.deepseek.com/beta/")

	reasoningContent, err := os.ReadFile("utils/reasoningContent.txt")
	require.NoError(t, err, "failed to read reasoning content file")

	tests := []struct {
		name        string
		req         *deepseek.ChatCompletionRequest
		wantErr     bool
		validateRes func(t *testing.T, res *deepseek.ChatCompletionResponse)
	}{
		{
			name: "basic prefix completion",
			req: &deepseek.ChatCompletionRequest{
				Model: deepseek.DeepSeekChat,
				Messages: []deepseek.ChatCompletionMessage{
					{Role: constants.ChatMessageRoleUser, Content: "Please write quick sort code"},
					{Role: constants.ChatMessageRoleAssistant, Content: "```python\n", Prefix: true},
				},
			},
			wantErr: false,
			validateRes: func(t *testing.T, res *deepseek.ChatCompletionResponse) {
				assert.NotEmpty(t, res.Choices[0].Message.Content)
			},
		},
		{
			name: "prefix with reasoner and reasoning_content",
			req: &deepseek.ChatCompletionRequest{
				Model: deepseek.DeepSeekReasoner,
				Messages: []deepseek.ChatCompletionMessage{
					{Role: constants.ChatMessageRoleUser, Content: "Please write quick sort code"},
					{Role: constants.ChatMessageRoleAssistant, Content: "```python\n", Prefix: true, ReasoningContent: string(reasoningContent)},
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

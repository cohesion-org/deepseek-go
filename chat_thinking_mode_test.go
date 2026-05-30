package deepseek_test

import (
	"context"
	"testing"

	"github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestThinkingModeCompletion(t *testing.T) {
	server := testutil.NewMockDeepSeekServer(t)
	defer server.Close()
	client := testutil.NewMockClient(t, server)

	tests := []struct {
		name        string
		req         *deepseek.ChatCompletionRequest
		wantErr     bool
		validateRes func(t *testing.T, res *deepseek.ChatCompletionResponse)
	}{
		{
			name: "deep think chat model",
			req: &deepseek.ChatCompletionRequest{
				Model: deepseek.DeepSeekV4Pro,
				Messages: []deepseek.ChatCompletionMessage{
					{Role: deepseek.ChatMessageRoleUser, Content: "Should this AI system be allowed to make final diagnostic decisions without human oversight?"},
				},
			},
			wantErr: false,
			validateRes: func(t *testing.T, res *deepseek.ChatCompletionResponse) {
				assert.NotEmpty(t, res.Choices[0].Message.Content)
				assert.NotEmpty(t, res.Choices[0].Message.ReasoningContent)
			},
		},
		{
			name: "empty messages",
			req: &deepseek.ChatCompletionRequest{
				Model:    deepseek.DeepSeekV4Pro,
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
			resp, err := client.CreateChatCompletion(context.Background(), tt.req)
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

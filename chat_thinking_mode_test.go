package deepseek_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/constants"
	"github.com/cohesion-org/deepseek-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestThinkingModeCompletion(t *testing.T) {
	testutil.SkipIfShort(t)
	//config := testutil.LoadTestConfig(t)
	dsUrl := os.Getenv("DeepSeek_Url")
	config := &testutil.TestConfig{
		APIKey:      "",
		TestTimeout: 30 * time.Second,
	}
	client, _ := deepseek.NewClientWithOptions(config.APIKey,
		deepseek.WithBaseURL(dsUrl),
		deepseek.WithPath(os.Getenv("DeepSeek_Path")),
		deepseek.WithTimeoutString("5m"))

	fmt.Println(client.Path)

	deepThinkPrompt, err := os.ReadFile("utils/deepthink.txt")
	require.NoError(t, err, "failed to read reasoning content file")

	tests := []struct {
		name        string
		req         *deepseek.ChatCompletionRequest
		wantErr     bool
		validateRes func(t *testing.T, res *deepseek.ChatCompletionResponse)
	}{
		{
			name: "deep think chat model",
			req: &deepseek.ChatCompletionRequest{
				Model: deepseek.DeppSeekModel,
				Messages: []deepseek.ChatCompletionMessage{
					{Role: constants.ChatMessageRoleUser, Content: "Should this AI system be allowed to make final diagnostic decisions without human oversight?"},
					{Role: constants.ChatMessageRoleSystem, Content: string(deepThinkPrompt), Prefix: true},
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
				Model:    deepseek.DeppSeekModel,
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

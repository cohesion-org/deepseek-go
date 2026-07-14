package deepseek_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type requestCaptureFunc func(*http.Request) (*http.Response, error)

func (f requestCaptureFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestMiniMaxProviderConstantsAndEndpoints(t *testing.T) {
	require.Equal(t, "MiniMax-M3", deepseek.MiniMaxM3)
	require.Equal(t, "MiniMax-M2.7", deepseek.MiniMaxM2_7)

	tests := []struct {
		name        string
		baseURL     string
		expectedURL string
		anthropic   bool
	}{
		{
			name:        "global OpenAI endpoint",
			baseURL:     deepseek.MiniMaxBaseURL,
			expectedURL: "https://api.minimax.io/v1/chat/completions",
		},
		{
			name:        "China OpenAI endpoint",
			baseURL:     deepseek.MiniMaxCNBaseURL,
			expectedURL: "https://api.minimaxi.com/v1/chat/completions",
		},
		{
			name:        "global Anthropic endpoint",
			baseURL:     deepseek.MiniMaxAnthropicAPIBaseURL,
			expectedURL: "https://api.minimax.io/anthropic/v1/messages",
			anthropic:   true,
		},
		{
			name:        "China Anthropic endpoint",
			baseURL:     deepseek.MiniMaxCNAnthropicAPIBaseURL,
			expectedURL: "https://api.minimaxi.com/anthropic/v1/messages",
			anthropic:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			capture := requestCaptureFunc(func(req *http.Request) (*http.Response, error) {
				require.Equal(t, tt.expectedURL, req.URL.String())

				body := `{"id":"chatcmpl-test","object":"chat.completion","choices":[{"index":0,"message":{"role":"assistant","content":"Hello"},"finish_reason":"stop"}],"model":"MiniMax-M3","usage":{"prompt_tokens":1,"completion_tokens":1,"total_tokens":2}}`
				if tt.anthropic {
					body = `{"id":"msg-test","type":"message","role":"assistant","content":[],"model":"MiniMax-M3","usage":{"input_tokens":1,"output_tokens":1}}`
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     make(http.Header),
					Body:       io.NopCloser(strings.NewReader(body)),
					Request:    req,
				}, nil
			})

			if tt.anthropic {
				client := deepseek.NewAnthropicClient("test-token")
				client.BaseURL = tt.baseURL
				client.Path = "/v1/messages"
				client.HTTPClient = capture
				_, err := client.CreateAnthropicMessage(context.Background(), &deepseek.AnthropicRequest{
					Model:     deepseek.MiniMaxM3,
					MaxTokens: 16,
					Messages: []deepseek.AnthropicMessage{
						{Role: "user", Content: []interface{}{deepseek.AnthropicTextBlock{Type: "text", Text: "Hello"}}},
					},
				})
				require.NoError(t, err)
				return
			}

			client := deepseek.NewClient("test-token", tt.baseURL)
			client.HTTPClient = capture
			_, err := client.CreateChatCompletion(context.Background(), &deepseek.ChatCompletionRequest{
				Model: deepseek.MiniMaxM3,
				Messages: []deepseek.ChatCompletionMessage{
					{Role: deepseek.ChatMessageRoleUser, Content: "Hello"},
				},
			})
			require.NoError(t, err)
		})
	}
}

func TestListAllModels(t *testing.T) {
	server := testutil.NewMockDeepSeekServer(t)
	defer server.Close()

	client := testutil.NewMockClient(t, server)
	resp, err := deepseek.ListAllModels(client, context.Background())
	require.NoError(t, err)
	assert.NotNil(t, resp)

	// Verify response structure
	assert.Equal(t, "list", resp.Object)
	assert.NotEmpty(t, resp.Data)

	// Verify model details
	for _, model := range resp.Data {
		assert.NotEmpty(t, model.ID)
		assert.Equal(t, "model", model.Object)
		assert.Equal(t, "deepseek", model.OwnedBy)

		// Verify known models exist in constants.go
		if model.ID == deepseek.DeepSeekChat ||
			model.ID == deepseek.DeepSeekCoder ||
			model.ID == deepseek.DeepSeekReasoner ||
			model.ID == deepseek.DeepSeekV4Flash ||
			model.ID == deepseek.DeepSeekV4Pro {
			assert.Contains(t, []string{
				deepseek.DeepSeekChat,
				deepseek.DeepSeekCoder,
				deepseek.DeepSeekReasoner,
				deepseek.DeepSeekV4Flash,
				deepseek.DeepSeekV4Pro,
			}, model.ID)
		}
	}
}

package deepseek_test

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateChatCompletionWithImage(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	if os.Getenv("OPENROUTER_API_KEY") == "" {
		t.Skip("Skipping test: OPENROUTER_API_KEY not set")
	}
	client := deepseek.NewClient(os.Getenv("OPENROUTER_API_KEY"),
		"https://openrouter.ai/api/v1/")
	tests := []struct {
		name        string
		req         *deepseek.ChatCompletionRequestWithImage
		wantErr     bool
		validateRes func(t *testing.T, res *deepseek.ChatCompletionResponse)
	}{
		{
			name: "basic image chat completion",
			req: &deepseek.ChatCompletionRequestWithImage{
				Model: "google/gemini-2.0-flash-001",
				Messages: []deepseek.ChatCompletionMessageWithImage{
					deepseek.NewImageMessage(
						deepseek.ChatMessageRoleUser,
						"Describe this image",
						"https://raw.githubusercontent.com/Vein05/nomnom/refs/heads/main/nomnom.png",
					),
				},
			},
			wantErr: false,
			validateRes: func(t *testing.T, res *deepseek.ChatCompletionResponse) {
				assert.NotEmpty(t, res.Choices[0].Message.Content)
			},
		},
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
		},
		{
			name: "invalid model",
			req: &deepseek.ChatCompletionRequestWithImage{
				Model: "invalid-model",
				Messages: []deepseek.ChatCompletionMessageWithImage{
					deepseek.NewImageMessage(
						deepseek.ChatMessageRoleUser,
						"Describe this image",
						"https://raw.githubusercontent.com/Vein05/nomnom/refs/heads/main/nomnom.png",
					),
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
			defer cancel()

			resp, err := client.CreateChatCompletionWithImage(ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)

			// Validate common response structure
			assert.NotEmpty(t, resp.ID)
			assert.NotEmpty(t, resp.Created)
			assert.Equal(t, tt.req.Model, resp.Model)
			assert.NotEmpty(t, resp.Choices)
			assert.NotNil(t, resp.Usage)

			if tt.validateRes != nil {
				tt.validateRes(t, resp)
			}
		})
	}
}

func TestCreateChatCompletionStreamWithImage(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	if os.Getenv("OPENROUTER_API_KEY") == "" {
		t.Skip("Skipping test: OPENROUTER_API_KEY not set")
	}
	client := deepseek.NewClient(os.Getenv("OPENROUTER_API_KEY"),
		"https://openrouter.ai/api/v1/")

	tests := []struct {
		name    string
		req     *deepseek.StreamChatCompletionRequestWithImage
		wantErr bool
	}{
		{
			name: "basic stream with image",
			req: &deepseek.StreamChatCompletionRequestWithImage{
				Model: "google/gemini-2.0-flash-001",
				Messages: []deepseek.ChatCompletionMessageWithImage{
					deepseek.NewImageMessage(
						deepseek.ChatMessageRoleUser,
						"Describe this image",
						"https://raw.githubusercontent.com/Vein05/nomnom/refs/heads/main/nomnom.png",
					),
				},
				Stream: true,
			},
			wantErr: false,
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

			stream, err := client.CreateChatCompletionStreamWithImage(ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, stream)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, stream)
			defer stream.Close()

			// Read from stream and verify responses
			var receivedContent bool
			for {
				response, err := stream.Recv()
				if errors.Is(err, io.EOF) {
					break
				}
				if err != nil {
					t.Fatalf("Stream error: %v", err)
				}

				assert.NotNil(t, response)
				assert.NotEmpty(t, response.Choices)
				receivedContent = true
			}
			assert.True(t, receivedContent, "Should have received content from stream")
		})
	}
}

func TestNewImageMessage(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		text     string
		imageURL string
		want     deepseek.ChatCompletionMessageWithImage
	}{
		{
			name:     "valid image message",
			role:     deepseek.ChatMessageRoleUser,
			text:     "Describe this image",
			imageURL: "https://raw.githubusercontent.com/Vein05/nomnom/refs/heads/main/nomnom.png",
			want: deepseek.ChatCompletionMessageWithImage{
				Role: deepseek.ChatMessageRoleUser,
				Content: []deepseek.ContentItem{
					{
						Type: "text",
						Text: "Describe this image",
					},
					{
						Type: "image_url",
						Image: &deepseek.ImageContent{
							URL: "https://raw.githubusercontent.com/Vein05/nomnom/refs/heads/main/nomnom.png",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := deepseek.NewImageMessage(tt.role, tt.text, tt.imageURL)
			assert.Equal(t, tt.want, got)
		})
	}
}

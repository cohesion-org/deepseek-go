package deepseek_test

import (
	"context"
	"strings"
	"testing"

	"github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateFIMCompletion(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	client := deepseek.NewClient(config.APIKey)

	tests := []struct {
		name        string
		req         *deepseek.FIMCompletionRequest
		wantErr     bool
		validateRes func(t *testing.T, res *deepseek.FIMCompletionResponse)
	}{
		{
			name:    "basic FIM completion",
			req:     deepseek.NewDefaultFIMCompletionRequest(deepseek.DeepSeekChat, "func main() {\n    fmt.Println(\"hel"),
			wantErr: false,
			validateRes: func(t *testing.T, res *deepseek.FIMCompletionResponse) {
				assert.NotEmpty(t, res.Choices[0].Text)
			},
		},
		{
			name:    "empty prompt",
			req:     deepseek.NewDefaultFIMCompletionRequest(deepseek.DeepSeekChat, ""),
			wantErr: true,
		},
		{
			name:    "invalid model",
			req:     deepseek.NewDefaultFIMCompletionRequest("invalid-model-123", "some code"),
			wantErr: true,
		},
		{
			name: "max tokens exceeded",
			req: deepseek.NewDefaultFIMCompletionRequest(deepseek.DeepSeekChat, "long prompt "+strings.Repeat("test ", 1000),
				deepseek.WithMaxTokens(5000)),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
			defer cancel()

			resp, err := client.CreateFIMCompletion(ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, resp)

			if tt.req != nil {
				// Validate common response structure
				assert.NotEmpty(t, resp.ID)
				assert.NotEmpty(t, resp.Created)
				assert.Equal(t, "text_completion", resp.Object)
				assert.Equal(t, tt.req.Model, resp.Model)
				assert.NotEmpty(t, resp.Choices)
				assert.NotNil(t, resp.Usage)
			}
			// Validate specific test case expectations
			if tt.validateRes != nil {
				tt.validateRes(t, resp)
			}
		})
	}
}

func TestCreateFIMCompletionWithParameters(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	client := deepseek.NewClient(config.APIKey)

	tests := []struct {
		name        string
		req         *deepseek.FIMCompletionRequest
		wantErr     bool
		validateRes func(t *testing.T, res *deepseek.FIMCompletionResponse)
	}{
		{
			name: "FIM completion with temperature and top_p",
			req: deepseek.NewDefaultFIMCompletionRequest(deepseek.DeepSeekChat, "func main() {\n    fmt.Println(\"hel",
				deepseek.WithTemperature(0.5), deepseek.WithTopP(0.9)),
			wantErr: false,
			validateRes: func(t *testing.T, res *deepseek.FIMCompletionResponse) {
				assert.NotEmpty(t, res.Choices[0].Text)
			},
		},
		{
			name: "FIM completion with invalid temperature",
			req: deepseek.NewDefaultFIMCompletionRequest(
				deepseek.DeepSeekChat,
				"func main() {\n    fmt.Println(\"hel", deepseek.WithTemperature(2.5)), // Invalid temperature
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
			defer cancel()

			resp, err := client.CreateFIMCompletion(ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, resp)

			if tt.req != nil {
				// Validate common response structure (same as in TestCreateFIMCompletion)
				assert.NotEmpty(t, resp.ID)
				assert.NotEmpty(t, resp.Created)
				assert.Equal(t, "text_completion", resp.Object)
				assert.Equal(t, tt.req.Model, resp.Model)
				assert.NotEmpty(t, resp.Choices)
				assert.NotNil(t, resp.Usage)
			}
			// Validate specific test case expectations
			if tt.validateRes != nil {
				tt.validateRes(t, resp)
			}
		})
	}
}

func TestFIMCompletionResponseStructure(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	client := deepseek.NewClient(config.APIKey)

	req := deepseek.NewDefaultFIMCompletionRequest(deepseek.DeepSeekChat, "func main() {\n    fmt.Println(\"hel")

	ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
	defer cancel()

	resp, err := client.CreateFIMCompletion(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Check overall response structure
	assert.NotEmpty(t, resp.ID)
	assert.NotEmpty(t, resp.Object)
	assert.NotEmpty(t, resp.Created)
	assert.Equal(t, "text_completion", resp.Object)
	assert.Equal(t, req.Model, resp.Model)
	assert.NotEmpty(t, resp.Choices)
	assert.NotNil(t, resp.Usage)

	// Check choices array
	for _, choice := range resp.Choices {
		assert.NotEmpty(t, choice.Text)
		assert.NotNil(t, choice.Index) // Or assert.GreaterOrEqual(t, choice.Index, 0)
		// LogProbs can be nil, so just check that it's present or handle it if you expect values.
		assert.NotEmpty(t, choice.FinishReason) // Check for valid values like "stop", "length"
	}

	// Check usage structure
	assert.GreaterOrEqual(t, resp.Usage.PromptTokens, 0)
	assert.GreaterOrEqual(t, resp.Usage.CompletionTokens, 0)
	assert.GreaterOrEqual(t, resp.Usage.TotalTokens, 0)
	assert.Equal(t, resp.Usage.PromptTokens+resp.Usage.CompletionTokens, resp.Usage.TotalTokens) // Check if the total is correct

}

// The api doesn't return logprobs in the response so this test will fail
// func TestFIMCompletionResponseWithLogProbs(t *testing.T) {
// 	testutil.SkipIfShort(t)
// 	config := testutil.LoadTestConfig(t)
// 	client := deepseek.NewClient(config.APIKey)

// 	req := &deepseek.FIMCompletionRequest{
// 		Model:    deepseek.DeepSeekChat,
// 		Prompt:   "func main() {\n    fmt.Println(\"hel",
// 		Logprobs: 10, // Request log probabilities
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
// 	defer cancel()

// 	resp, err := client.CreateFIMCompletion(ctx, req)
// 	require.NoError(t, err)
// 	require.NotNil(t, resp)

// 	// Check choices array
// 	for _, choice := range resp.Choices {
// 		assert.NotEmpty(t, choice.Text)
// 		assert.NotNil(t, choice.Index)
// 		assert.NotNil(t, choice.Logprobs) // LogProbs should be present

// 		logProbs := choice.Logprobs
// 		log.Printf("LogProbs: %+v\n", logProbs)
// 		assert.NotEmpty(t, logProbs.Content) // Check Content is not empty
// 		for _, contentToken := range logProbs.Content {
// 			assert.NotEmpty(t, contentToken.Token)
// 			// assert.NotZero(t, contentToken.Logprob) // Consider checking Logprob value
// 			if contentToken.Bytes != nil {
// 				assert.NotEmpty(t, contentToken.Bytes) // Check bytes if present
// 			}
// 		}
// 		assert.NotEmpty(t, logProbs.TopLogprobs) // Check TopLogprobs is not empty
// 		for _, topLogprobToken := range logProbs.TopLogprobs {
// 			assert.NotEmpty(t, topLogprobToken.Token)
// 			// assert.NotZero(t, topLogprobToken.Logprob) // Consider checking Logprob value
// 			if topLogprobToken.Bytes != nil {
// 				assert.NotEmpty(t, topLogprobToken.Bytes) // Check bytes if present
// 			}
// 		}

// 		assert.NotEmpty(t, choice.FinishReason)
// 	}
// }

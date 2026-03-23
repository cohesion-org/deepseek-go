//go:build integration

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

func TestIntegrationChatCompletion(t *testing.T) {
	testutil.SkipUnlessLiveAPI(t)
	config := testutil.LoadTestConfig(t)
	client := deepseek.NewClient(config.APIKey)

	ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
	defer cancel()

	resp, err := client.CreateChatCompletion(ctx, &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: constants.ChatMessageRoleUser, Content: "Reply with exactly: integration-ok"},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.Choices)
	assert.NotEmpty(t, resp.Choices[0].Message.Content)
}

func TestIntegrationChatCompletionStream(t *testing.T) {
	testutil.SkipUnlessLiveAPI(t)
	config := testutil.LoadTestConfig(t)
	client := deepseek.NewClient(config.APIKey)

	ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
	defer cancel()

	stream, err := client.CreateChatCompletionStream(ctx, &deepseek.StreamChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: constants.ChatMessageRoleUser, Content: "Say hello in one short sentence."},
		},
	})
	require.NoError(t, err)
	defer stream.Close()

	var got string
	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err)
		if len(resp.Choices) > 0 {
			got += resp.Choices[0].Delta.Content
		}
	}

	assert.NotEmpty(t, got)
}

func TestIntegrationListAllModels(t *testing.T) {
	testutil.SkipUnlessLiveAPI(t)
	config := testutil.LoadTestConfig(t)
	client := deepseek.NewClient(config.APIKey)

	ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
	defer cancel()

	resp, err := deepseek.ListAllModels(client, ctx)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.Data)
}

func TestIntegrationGetBalance(t *testing.T) {
	testutil.SkipUnlessLiveAPI(t)
	config := testutil.LoadTestConfig(t)
	client := deepseek.NewClient(config.APIKey)

	ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
	defer cancel()

	resp, err := deepseek.GetBalance(client, ctx)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotNil(t, resp.BalanceInfos)
}

func TestIntegrationFIMCompletion(t *testing.T) {
	testutil.SkipUnlessLiveAPI(t)
	config := testutil.LoadTestConfig(t)
	client := deepseek.NewClient(config.APIKey)

	ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
	defer cancel()

	resp, err := client.CreateFIMCompletion(ctx, &deepseek.FIMCompletionRequest{
		Model:  deepseek.DeepSeekChat,
		Prompt: "func main() {\n    fmt.Println(\"hel",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.Choices)
	assert.NotEmpty(t, resp.Choices[0].Text)
}

func TestChatTool(t *testing.T) {
	testutil.SkipUnlessLiveAPI(t)
	config := testutil.LoadTestConfig(t)
	client := deepseek.NewClient(config.APIKey)

	ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
	defer cancel()

	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Tools: []deepseek.Tool{
			{
				Type: "function",
				Function: deepseek.Function{
					Name:        "get_weather",
					Description: "Get weather of a location, the user should supply a location first",
					Parameters: &deepseek.FunctionParameters{
						Type: "object",
						Properties: map[string]interface{}{
							"location": map[string]interface{}{
								"type":        "string",
								"description": "The city and state to get the weather for",
							},
						},
						Required: []string{"location"},
					},
				},
			},
		},
		Messages: []deepseek.ChatCompletionMessage{
			{
				Role:    constants.ChatMessageRoleUser,
				Content: "How's the weather in Hangzhou?",
			},
		},
	}

	response, err := client.CreateChatCompletion(ctx, request)
	require.NoError(t, err)
	require.NotEmpty(t, response.Choices)

	assistantMsg := response.Choices[0].Message
	require.NotEmpty(t, assistantMsg.ToolCalls, "expected at least one tool call in the first response")

	request.Messages = append(request.Messages, deepseek.ChatCompletionMessage{
		Role:             deepseek.ChatMessageRoleAssistant,
		Content:          assistantMsg.Content,
		ReasoningContent: assistantMsg.ReasoningContent,
		ToolCalls:        assistantMsg.ToolCalls,
	})

	request.Messages = append(request.Messages, deepseek.ChatCompletionMessage{
		Role:       constants.ChatMessageRoleTool,
		Content:    `{"location":"Hangzhou","temperature":"24℃"}`,
		ToolCallID: assistantMsg.ToolCalls[0].ID,
	})

	finalResponse, err := client.CreateChatCompletion(ctx, request)
	require.NoError(t, err)
	require.NotEmpty(t, finalResponse.Choices)
	assert.NotEmpty(t, finalResponse.Choices[0].Message.Content)
}

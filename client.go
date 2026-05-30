package deepseek

import (
	"bufio"
	"context"
	"fmt"
	"os"

	utils "github.com/cohesion-org/deepseek-go/utils"
)

// CreateChatCompletion sends a chat completion request and returns the generated response.
func (c *Client) CreateChatCompletion(
	ctx context.Context,
	request *ChatCompletionRequest,
) (*ChatCompletionResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	warnDeprecatedModel(request.Model)

	payload, err := normalizeChatCompletionPayload(request)
	if err != nil {
		return nil, fmt.Errorf("error normalizing request payload: %w", err)
	}

	ctx, tcancel, err := getTimeoutContext(ctx, c.Timeout)
	if err != nil {
		return nil, fmt.Errorf("error getting timeout context: %w", err)
	}
	defer tcancel()

	baseURL := c.BaseURL
	if hasStrictTool(request.Tools) {
		baseURL = resolveBetaBaseURL(c.BaseURL)
	}

	req, err := utils.NewRequestBuilder(c.AuthToken).
		SetBaseURL(baseURL).
		SetPath(c.Path).
		SetBodyFromStruct(payload).
		Build(ctx)

	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}
	resp, err := HandleSendChatCompletionRequest(*c, req)

	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= 400 {
		return nil, HandleAPIError(resp)
	}

	updatedResp, err := HandleChatCompletionResponse(resp)

	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return updatedResp, nil
}

// CreateChatCompletionStream sends a chat completion request with stream = true and returns the delta
func (c *Client) CreateChatCompletionStream(
	ctx context.Context,
	request *StreamChatCompletionRequest,
) (ChatCompletionStream, error) {
	if request == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	warnDeprecatedModel(request.Model)

	ctx, _, err := getTimeoutContext(ctx, c.Timeout)
	if err != nil {
		return nil, fmt.Errorf("error getting timeout context: %w", err)
	}

	request.Stream = true
	payload, err := normalizeStreamChatCompletionPayload(request)
	if err != nil {
		return nil, fmt.Errorf("error normalizing request payload: %w", err)
	}

	baseURL := c.BaseURL
	if hasStrictTool(request.Tools) {
		baseURL = resolveBetaBaseURL(c.BaseURL)
	}

	req, err := utils.NewRequestBuilder(c.AuthToken).
		SetBaseURL(baseURL).
		SetPath(c.Path).
		SetBodyFromStruct(payload).
		BuildStream(ctx)

	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}

	resp, err := HandleSendChatCompletionRequest(*c, req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, HandleAPIError(resp)
	}

	ctx, cancel := context.WithCancel(ctx)
	stream := &chatCompletionStream{
		ctx:    ctx,
		cancel: cancel,
		resp:   resp,
		reader: bufio.NewReader(resp.Body),
	}
	return stream, nil
}

// CreateFIMCompletion is a beta feature. It sends a FIM completion request and returns the generated response.
// the base URL is set to "https://api.deepseek.com/beta/"
func (c *Client) CreateFIMCompletion(
	ctx context.Context,
	request *FIMCompletionRequest,
) (*FIMCompletionResponse, error) {
	if request.MaxTokens > 4000 {
		return nil, fmt.Errorf("max tokens must be <= 4000")
	}
	baseURL := resolveBetaBaseURL(c.BaseURL)

	if request == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	warnDeprecatedModel(request.Model)
	req, err := utils.NewRequestBuilder(c.AuthToken).
		SetBaseURL(baseURL).
		SetPath("/completions").
		SetBodyFromStruct(request).
		Build(ctx)
	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}
	resp, err := HandleSendChatCompletionRequest(*c, req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode >= 400 {
		return nil, HandleAPIError(resp)
	}
	updatedResp, err := HandleFIMCompletionRequest(resp)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return updatedResp, nil
}

// CreateFIMStreamCompletion sends a FIM completion request with stream = true and returns the delta
func (c *Client) CreateFIMStreamCompletion(
	ctx context.Context,
	request *FIMStreamCompletionRequest,
) (FIMChatCompletionStream, error) {
	baseURL := resolveBetaBaseURL(c.BaseURL)

	warnDeprecatedModel(request.Model)

	request.Stream = true
	req, err := utils.NewRequestBuilder(c.AuthToken).
		SetBaseURL(baseURL).
		SetPath("/completions"). //Note to maintainer: This is a really bad implementation with manual path insertion. Please create an issue.
		SetBodyFromStruct(request).
		BuildStream(ctx)

	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}

	resp, err := HandleSendChatCompletionRequest(*c, req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, HandleAPIError(resp)
	}

	ctx, cancel := context.WithCancel(ctx)
	stream := &fimCompletionStream{
		ctx:    ctx,
		cancel: cancel,
		resp:   resp,
		reader: bufio.NewReader(resp.Body),
	}
	return stream, nil
}


var deprecatedModels = map[string]string{
	"deepseek-chat":     "deepseek-v4-flash (non-thinking mode)",
	"deepseek-coder":    "deepseek-v4-flash",
	"deepseek-reasoner": "deepseek-v4-flash (thinking mode)",
}

func warnDeprecatedModel(model string) {
	if replacement, ok := deprecatedModels[model]; ok {
		fmt.Fprintf(os.Stderr, "deepseek-go: WARNING model %q is deprecated and will be removed 2026/07/24. Use %s instead.\n", model, replacement)
	}
}

func hasStrictTool(tools []Tool) bool {
	for _, t := range tools {
		if t.Function.Strict {
			return true
		}
	}
	return false
}

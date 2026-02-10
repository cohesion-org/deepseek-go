package deepseek

import (
	"bufio"
	"context"
	"fmt"

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

	ctx, tcancel, err := getTimeoutContext(ctx, c.Timeout)
	if err != nil {
		return nil, fmt.Errorf("error getting timeout context: %w", err)
	}
	defer tcancel()

	req, err := utils.NewRequestBuilder(c.AuthToken).
		SetBaseURL(c.BaseURL).
		SetPath(c.Path).
		SetBodyFromStruct(request).
		Build(ctx)

	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}
	resp, err := HandleSendChatCompletionRequest(*c, req)

	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

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

	ctx, _, err := getTimeoutContext(ctx, c.Timeout)
	if err != nil {
		return nil, fmt.Errorf("error getting timeout context: %w", err)
	}

	request.Stream = true
	req, err := utils.NewRequestBuilder(c.AuthToken).
		SetBaseURL(c.BaseURL).
		SetPath(c.Path).
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
	baseURL := "https://api.deepseek.com/beta/"

	if request == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
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
	defer resp.Body.Close()
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
	baseURL := "https://api.deepseek.com/beta/"

	request.Stream = true
	req, err := utils.NewRequestBuilder(c.AuthToken).
		SetBaseURL(baseURL).
		SetPath("/completions"). //Note to maintianer: This is a really bad implementation with manual path insertion. Please create an issue.
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

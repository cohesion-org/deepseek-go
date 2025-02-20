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
	if err := CheckChatCompletionRequest(request); err != nil {
		return nil, fmt.Errorf("CreateChatCompletion error: %w", err)
	}

	req, err := utils.NewRequestBuilder(c.AuthToken).
		SetBaseURL(c.BaseURL).
		SetPath("chat/completions").
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

	return updatedResp, err
}

// CreateChatCompletionStream sends a chat completion request with stream = true and returns the delta
func (c *Client) CreateChatCompletionStream(
	ctx context.Context,
	request *StreamChatCompletionRequest,
) (ChatCompletionStream, error) {
	if err := CheckStreamChatCompletionRequest(request); err != nil {
		return nil, fmt.Errorf("CreateChatCompletionStream error: %w", err)
	}
	
	request.Stream = true
	req, err := utils.NewRequestBuilder(c.AuthToken).
		SetBaseURL(c.BaseURL).
		SetPath("chat/completions").
		SetBodyFromStruct(request).
		BuildStream(ctx)

	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}

	resp, err := HandleSendChatCompletionRequest(*c, req)
	if err != nil {
		return nil, err
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
	baseURL := "https://api.deepseek.com/beta/"
	if err := CheckFIMCompletionRequest(request); err != nil {
		return nil, fmt.Errorf("CreateFIMCompletion request is invalid, error: %w", err)
	}

	req, err := utils.NewRequestBuilder(c.AuthToken).
		SetBaseURL(baseURL).
		SetPath("completions").
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
	return updatedResp, err
}

// CreateFIMStreamCompletion sends a FIM completion request with stream = true and returns the delta
func (c *Client) CreateFIMStreamCompletion(
	ctx context.Context,
	request *FIMStreamCompletionRequest,
) (FIMChatCompletionStream, error) {
	baseURL := "https://api.deepseek.com/beta/"

	if err := CheckFIMStreamCompletionRequest(request); err != nil {
		return nil, fmt.Errorf("CheckFIMStreamCompletionRequest request is invalid, error: %w", err)
	}

	request.Stream = true
	req, err := utils.NewRequestBuilder(c.AuthToken).
		SetBaseURL(baseURL).
		SetPath("/completions").
		SetBodyFromStruct(request).
		BuildStream(ctx)

	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}

	resp, err := HandleSendChatCompletionRequest(*c, req)
	if err != nil {
		return nil, err
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

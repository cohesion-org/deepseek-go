package deepseek

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// HandleTimeout gets the timeout duration from the DEEPSEEK_TIMEOUT environment variable.
func HandleTimeout() (time.Duration, error) {
	return handleTimeout()
}

// handleTimeout checks the DEEPSEEK_TIMEOUT environment variable and returns the timeout duration.
func handleTimeout() (time.Duration, error) {
	if err := godotenv.Load(); err != nil {
		_ = err
	}

	timeoutStr := os.Getenv("DEEPSEEK_TIMEOUT")
	if timeoutStr == "" {
		return 5 * time.Minute, nil
	}
	duration, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return 0, fmt.Errorf("invalid timeout duration %q: %w", timeoutStr, err)
	}
	return duration, nil
}

// getTimeoutContext creates a context with a timeout.
// If the timeout is less than or equal to 0, it tries to get the timeout from the environment variable.
// If the timeout is greater than 0, it creates a context with that timeout.
// It returns the context, a cancel function, and an error if any.
func getTimeoutContext(ctx context.Context, timeout time.Duration) (
	context.Context,
	context.CancelFunc,
	error,
) {
	if timeout <= 0 {
		// Try to get timeout from environment variable
		var err error
		timeout, err = handleTimeout()
		if err != nil {
			return nil, nil, fmt.Errorf("error getting timeout from environment: %w", err)
		}
	}

	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	} else {
		cancel = func() {}
	}

	return ctx, cancel, nil
}

// HandleSendChatCompletionRequest sends a request to the DeepSeek API and returns the response.
func HandleSendChatCompletionRequest(c Client, req *http.Request) (*http.Response, error) {
	return c.handleRequest(req)
}

// HandleNormalRequest sends a request to the DeepSeek API and returns the response.
func HandleNormalRequest(c Client, req *http.Request) (*http.Response, error) {
	return c.handleRequest(req)
}

// handleRequest sends the HTTP request using the provided HTTP client.
// If no client is provided, it uses the default HTTP client.
func (c *Client) handleRequest(req *http.Request) (*http.Response, error) {
	client := c.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	return resp, nil
}

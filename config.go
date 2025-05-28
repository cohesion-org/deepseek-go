package deepseek

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

// BaseURL is the base URL for the Deepseek API
const BaseURL string = "https://api.deepseek.com/v1"

// HTTPDoer is an interface for the Do method of http.Client
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is the main struct for interacting with the Deepseek API.
type Client struct {
	AuthToken string        // The authentication token for the API
	BaseURL   string        // The base URL for the API
	Timeout   time.Duration // The timeout for the current Client
	Path      string        // The path for the API request. Defaults to "chat/completions"

	HTTPClient HTTPDoer // The HTTP client to send the request and get the response
}

// NewClient creates a new client with an authentication token and an optional custom baseURL.
// If no baseURL is provided, it defaults to "https://api.deepseek.com/".
// You can't set path with this method. If you want to set path, use NewClientWithOptions.
func NewClient(AuthToken string, baseURL ...string) *Client {
	if AuthToken == "" {
		if envKey, ok := os.LookupEnv("DEEPSEEK_API_KEY"); ok && envKey != "" {
			AuthToken = envKey
		} else {
			fmt.Printf("authToken is empty. Please provide a valid token or set the DEEPSEEK_API_KEY environment variable")
			return nil
		}
	}
	// check if this is a valid URL
	if len(baseURL) > 0 {
		_, err := url.ParseRequestURI(baseURL[0])
		if err != nil {
			fmt.Printf("Invalid URL: %s. \nIf you are using options please use NewClientWithOptions", baseURL[0])
			return nil
		}
	}
	url := "https://api.deepseek.com/"
	if len(baseURL) > 0 {
		url = baseURL[0]
	}
	return &Client{
		AuthToken: AuthToken,
		BaseURL:   url,
		Path:      "chat/completions",
	}
}

// Option configures a Client instance
type Option func(*Client) error

// NewClientWithOptions creates a new client with required authentication token and optional configurations.
// Defaults:
// - BaseURL: "https://api.deepseek.com/"
// - Timeout: 5 minutes
func NewClientWithOptions(authToken string, opts ...Option) (*Client, error) {
	// Check for empty auth token and try to use environment variable
	if authToken == "" {
		if envKey, ok := os.LookupEnv("DEEPSEEK_API_KEY"); ok {
			authToken = envKey
		} else {
			return nil, fmt.Errorf("authToken is empty and DEEPSEEK_API_KEY environment variable is not set")
		}
	}
	if authToken == "" {
		return nil, fmt.Errorf("authToken is empty. Please provide a valid token or set the DEEPSEEK_API_KEY environment variable.")
	}

	client := &Client{
		AuthToken: authToken,
		BaseURL:   "https://api.deepseek.com/",
		Timeout:   5 * time.Minute,
		Path:      "chat/completions",
	}

	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return client, nil
}

// WithBaseURL sets the base URL for the API client
func WithBaseURL(url string) Option {
	return func(c *Client) error {
		c.BaseURL = url
		return nil
	}
}

// WithTimeout sets the timeout for API requests
func WithTimeout(d time.Duration) Option {
	return func(c *Client) error {
		if d < 0 {
			return fmt.Errorf("timeout must be a positive duration")
		}
		c.Timeout = d
		return nil
	}
}

// WithTimeoutString parses a duration string and sets the timeout
// Example valid values: "5s", "2m", "1h"
func WithTimeoutString(s string) Option {
	return func(c *Client) error {
		d, err := time.ParseDuration(s)
		if err != nil {
			return fmt.Errorf("invalid timeout duration %q: %w", s, err)
		}
		return WithTimeout(d)(c)
	}
}

// WithPath sets the path for the API request. Defaults to "chat/completions", if not set.
// Example usages would be "/c/chat/" or any http after the baseURL extension
func WithPath(path string) Option {
	if path == "" {
		path = "chat/completions"
	}
	return func(c *Client) error {
		c.Path = path
		return nil
	}
}

// WithHTTPClient sets the http client for the API client.
func WithHTTPClient(httpclient HTTPDoer) Option {
	return func(c *Client) error {
		c.HTTPClient = httpclient
		return nil
	}
}

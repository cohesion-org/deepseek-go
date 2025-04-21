package deepseek_test

import (
	"testing"
	"time"

	"github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestNewCleint(t *testing.T) {
	testutil.SkipIfShort(t)
	//test empty api key
	client := deepseek.NewClient("")
	require.Nil(t, client)

	//test valid api key
	client = deepseek.NewClient("test")
	require.NotNil(t, client)

	//test with a non url as base url
	client = deepseek.NewClient("test", "invalid-url")
	require.Nil(t, client)

	//test with base url
	client = deepseek.NewClient("test", "https://api.deepseek.com/")
	require.NotNil(t, client)
	require.Equal(t, "https://api.deepseek.com/", client.BaseURL)
}

func TestNewClientWithOptions(t *testing.T) {
	tests := []struct {
		name            string
		opts            []deepseek.Option
		expectedURL     string
		expectedTimeout time.Duration
		expectError     bool
	}{
		{
			name:            "default options",
			opts:            nil,
			expectedURL:     "https://api.deepseek.com/",
			expectedTimeout: 5 * time.Minute,
			expectError:     false,
		},
		{
			name:            "custom base URL",
			opts:            []deepseek.Option{deepseek.WithBaseURL("http://test.com")},
			expectedURL:     "http://test.com",
			expectedTimeout: 5 * time.Minute,
			expectError:     false,
		},
		{
			name:            "custom timeout",
			opts:            []deepseek.Option{deepseek.WithTimeout(10 * time.Second)},
			expectedURL:     "https://api.deepseek.com/",
			expectedTimeout: 10 * time.Second,
			expectError:     false,
		},
		{
			name:        "invalid timeout",
			opts:        []deepseek.Option{deepseek.WithTimeout(-1 * time.Second)},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := deepseek.NewClientWithOptions("token", tt.opts...)
			if tt.expectError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if client.BaseURL != tt.expectedURL {
				t.Errorf("expected BaseURL %s, got %s", tt.expectedURL, client.BaseURL)
			}
			if client.Timeout != tt.expectedTimeout {
				t.Errorf("expected Timeout %v, got %v", tt.expectedTimeout, client.Timeout)
			}
		})
	}
}

func TestNewClientWithOptions_EmptyAuthToken(t *testing.T) {
	// Test with empty auth token and no environment variable
	t.Setenv("DEEPSEEK_API_KEY", "") // Clear the env var for testing

	// Should return error when both auth token and env var are empty
	client, err := deepseek.NewClientWithOptions("")
	if err == nil {
		t.Error("Expected error when both auth token and environment variable are empty, got nil")
	}
	if client != nil {
		t.Error("Expected nil client when auth token validation fails")
	}

	// Test with empty auth token but environment variable set
	t.Setenv("DEEPSEEK_API_KEY", "env-api-key")
	client, err = deepseek.NewClientWithOptions("")
	if err != nil {
		t.Errorf("Expected no error when environment variable is set, got: %v", err)
	}
	if client == nil {
		t.Error("Expected non-nil client when environment variable is used")
	}
	if client != nil && client.AuthToken != "env-api-key" {
		t.Errorf("Expected auth token to be 'env-api-key', got '%s'", client.AuthToken)
	}

	// Test with non-empty auth token (should take precedence over env var)
	client, err = deepseek.NewClientWithOptions("direct-api-key")
	if err != nil {
		t.Errorf("Expected no error when auth token is provided directly, got: %v", err)
	}
	if client == nil {
		t.Error("Expected non-nil client when auth token is provided directly")
	}
	if client != nil && client.AuthToken != "direct-api-key" {
		t.Errorf("Expected auth token to be 'direct-api-key', got '%s'", client.AuthToken)
	}
}

func TestWithTimeoutString(t *testing.T) {
	tests := []struct {
		duration  string
		expectErr bool
		expected  time.Duration
	}{
		{"5s", false, 5 * time.Second},
		{"invalid", true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.duration, func(t *testing.T) {
			opt := deepseek.WithTimeoutString(tt.duration)
			client := &deepseek.Client{}
			err := opt(client)
			if tt.expectErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if client.Timeout != tt.expected {
				t.Errorf("expected timeout %v, got %v", tt.expected, client.Timeout)
			}
		})
	}
}

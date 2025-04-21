package deepseek_test

import (
	"os"
	"testing"
	"time"

	"github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestNewCleint(t *testing.T) {
	testutil.SkipIfShort(t)
	// test empty api key with no env fallback
	t.Run("empty api key with no env fallback", func(t *testing.T) {
		os.Unsetenv("DEEPSEEK_API_KEY")
		client := deepseek.NewClient("")
		require.Nil(t, client)
	})

	// test empty api key with env fallback
	t.Run("empty api key with env fallback", func(t *testing.T) {
		os.Setenv("DEEPSEEK_API_KEY", "test")
		defer os.Unsetenv("DEEPSEEK_API_KEY")

		client := deepseek.NewClient("")
		require.NotNil(t, client)
		require.Equal(t, "test", client.AuthToken)
	})

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
	tests := []struct {
		name          string
		envVar        string
		inputToken    string
		wantErr       bool
		wantClientNil bool
		wantAuthToken string
	}{
		{
			name:          "empty token and empty env var",
			envVar:        "",
			inputToken:    "",
			wantErr:       true,
			wantClientNil: true,
		},
		{
			name:          "empty token, env var set",
			envVar:        "env-api-key",
			inputToken:    "",
			wantErr:       false,
			wantClientNil: false,
			wantAuthToken: "env-api-key",
		},
		{
			name:          "direct token takes precedence",
			envVar:        "env-api-key",
			inputToken:    "direct-api-key",
			wantErr:       false,
			wantClientNil: false,
			wantAuthToken: "direct-api-key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("DEEPSEEK_API_KEY", tt.envVar)
			client, err := deepseek.NewClientWithOptions(tt.inputToken)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if client != nil {
					t.Error("expected nil client when auth token validation fails")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.wantClientNil && client != nil {
				t.Error("expected nil client, got non-nil")
			}
			if !tt.wantClientNil && client == nil {
				t.Error("expected non-nil client, got nil")
			}
			if client != nil && client.AuthToken != tt.wantAuthToken {
				t.Errorf("expected auth token to be '%s', got '%s'", tt.wantAuthToken, client.AuthToken)
			}
		})
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

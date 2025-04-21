// Package testutil provides testing utilities for the DeepSeek client.
package testutil

import (
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

// TestConfig holds test configuration loaded from environment
type TestConfig struct {
	APIKey      string
	TestTimeout time.Duration
}

// LoadTestConfig loads test configuration from environment variables
func LoadTestConfig(t *testing.T) *TestConfig {
	t.Helper()

	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// Ignore error since .env file is optional
		_ = err
	}

	defer func() {
		os.Unsetenv("TEST_DEEPSEEK_API_KEY")
		os.Unsetenv("TEST_TIMEOUT")
	}()

	var apiKey string
	apiKey = os.Getenv("TEST_DEEPSEEK_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("DEEPSEEK_API_KEY")
		t.Log("TEST_DEEPSEEK_API_KEY not ofund, trying DEEPSEEK_API_KEY from environment")
	}
	config := &TestConfig{
		APIKey:      apiKey,
		TestTimeout: 30 * time.Second,
	}

	// Override with environment variables if set
	if timeout := os.Getenv("TEST_TIMEOUT"); timeout != "" {
		if d, err := time.ParseDuration(timeout); err == nil {
			config.TestTimeout = d
		}
	}
	// Skip tests if API key is not set
	if config.APIKey == "" {
		t.Skip("Skipping test: neither TEST_DEEPSEEK_API_KEY nor DEEPSEEK_API_KEY is set")
	}

	return config
}

// SkipIfShort skips long-running tests when -short flag is used
func SkipIfShort(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}
}

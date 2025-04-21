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

	var apiKey string
	if os.Getenv("TEST_DEEPSEEK_API_KEY") == "" {
		apiKey = os.Getenv("DEEPSEEK_API_KEY")
	} else {
		apiKey = os.Getenv("TEST_DEEPSEEK_API_KEY")
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
		t.Skip("Skipping test: TEST_DEEPSEEK_API_KEY not set")
	}

	return config
}

// SkipIfShort skips long-running tests when -short flag is used
func SkipIfShort(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}
}

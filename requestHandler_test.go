package deepseek_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/cohesion-org/deepseek-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleSendChatCompletionRequest(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()
		c := deepseek.NewClient("test", ts.URL)

		req, err := http.NewRequest("GET", ts.URL, nil)
		require.NoError(t, err)

		resp, err := deepseek.HandleSendChatCompletionRequest(*c, req)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("request error", func(t *testing.T) {
		req, err := http.NewRequest("GET", "http://invalid-url", nil)
		require.NoError(t, err)
		c := deepseek.NewClient("test", "http://invalid-url")
		resp, err := deepseek.HandleSendChatCompletionRequest(*c, req)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "error sending request:")
	})
}

func TestHandleNormalRequest(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()
		c := deepseek.NewClient("test", ts.URL)
		req, err := http.NewRequest("GET", ts.URL, nil)
		require.NoError(t, err)

		resp, err := deepseek.HandleNormalRequest(*c, req)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("request error", func(t *testing.T) {
		req, err := http.NewRequest("GET", "http://invalid-url", nil)
		require.NoError(t, err)

		c := deepseek.NewClient("test", "http://invalid-url")
		resp, err := deepseek.HandleNormalRequest(*c, req)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "error sending request:")
	})
}

func TestTimeoutConfiguration(t *testing.T) {
	t.Run("chat completion timeout", func(t *testing.T) {
		start := time.Now()
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			time.Sleep(250 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		req, err := http.NewRequest("GET", ts.URL, nil)
		require.NoError(t, err)
		c := deepseek.NewClient("test", ts.URL)
		resp, err := deepseek.HandleSendChatCompletionRequest(*c, req)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()

		// Verify timeout configuration isn't too short
		assert.WithinDuration(t, start, time.Now(), 300*time.Millisecond)
	})

	t.Run("normal request timeout", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			time.Sleep(50 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		req, err := http.NewRequest("GET", ts.URL, nil)
		require.NoError(t, err)

		start := time.Now()
		c := deepseek.NewClient("test", ts.URL)
		resp, err := deepseek.HandleNormalRequest(*c, req)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()

		// Verify timeout allows successful completion
		assert.WithinDuration(t, start, time.Now(), 100*time.Millisecond)
	})
}

func TestErrorTimeout(t *testing.T) {
	t.Run("client timeout preservation", func(t *testing.T) {
		// Create test server that responds slower than client timeout
		ts := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			time.Sleep(250 * time.Millisecond) // Longer than client timeout
		}))
		defer ts.Close()

		req, err := http.NewRequest("GET", ts.URL, nil)
		require.NoError(t, err)

		// Create client with short timeout
		shortClient := &http.Client{
			Timeout: 5 * time.Millisecond,
		}

		_, err = shortClient.Do(req)
		require.Error(t, err)

		assert.True(t, errors.Is(err, context.DeadlineExceeded) ||
			(os.IsTimeout(err)), "should be timeout error")
	})
}

func TestHandleTimeoutApplication(t *testing.T) {
	t.Run("default timeout applied when no env var set", func(t *testing.T) {
		// Ensure no timeout environment variable is set
		_ = os.Unsetenv("DEEPSEEK_TIMEOUT")

		expectedTimeout := 5 * time.Minute

		timeout, err := deepseek.HandleTimeout()
		require.NoError(t, err)
		assert.Equal(t, expectedTimeout, timeout, "Expected default timeout when env var is missing")

		client, err := deepseek.NewClientWithOptions("your-api-key", deepseek.WithTimeout(timeout))
		require.NoError(t, err)
		assert.Equal(t, expectedTimeout, client.Timeout, "Client timeout should match default")
	})

	t.Run("custom timeout from environment variable", func(t *testing.T) {
		_ = os.Setenv("DEEPSEEK_TIMEOUT", "30s")
		defer func() {
			_ = os.Unsetenv("DEEPSEEK_TIMEOUT")
		}()

		expectedTimeout := 30 * time.Second

		timeout, err := deepseek.HandleTimeout()
		require.NoError(t, err)
		assert.Equal(t, expectedTimeout, timeout, "Expected timeout from env variable")

		client, err := deepseek.NewClientWithOptions("your-api-key", deepseek.WithTimeout(timeout))
		require.NoError(t, err)
		assert.Equal(t, expectedTimeout, client.Timeout, "Client timeout should match env var")
	})

	t.Run("invalid timeout format returns error", func(t *testing.T) {
		_ = os.Setenv("DEEPSEEK_TIMEOUT", "invalid")
		defer func() {
			_ = os.Unsetenv("DEEPSEEK_TIMEOUT")
		}()

		timeout, err := deepseek.HandleTimeout()
		assert.Error(t, err, "Expected error for invalid timeout format")
		assert.Zero(t, timeout, "Invalid timeout should return zero duration")
	})
}

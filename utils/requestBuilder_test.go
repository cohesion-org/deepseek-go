package utils_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	utils "github.com/cohesion-org/deepseek-go/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Define a custom type for context keys
type contextKey string

// Create a specific key for testing
const testContextKey contextKey = "test"

func TestAuthedRequestBuilder(t *testing.T) {
	const (
		testToken   = "test-token-123"
		testBaseURL = "https://api.example.com"
		testPath    = "/v1/chat"
	)

	t.Run("NewRequestBuilder", func(t *testing.T) {
		builder := utils.NewRequestBuilder(testToken)
		assert.Equal(t, testToken, builder.AuthToken)
	})

	t.Run("SetBaseURL", func(t *testing.T) {
		builder := utils.NewRequestBuilder(testToken).SetBaseURL(testBaseURL)
		assert.Equal(t, testBaseURL, builder.BaseURL)
	})

	t.Run("SetPath", func(t *testing.T) {
		builder := utils.NewRequestBuilder(testToken).SetPath(testPath)
		assert.Equal(t, testPath, builder.Path)
	})

	t.Run("SetBodyFromStruct", func(t *testing.T) {
		testBody := struct {
			Message string `json:"message"`
		}{Message: "test"}

		t.Run("valid struct", func(t *testing.T) {
			builder := utils.NewRequestBuilder(testToken).SetBodyFromStruct(testBody)
			expected, _ := json.Marshal(testBody)
			assert.Equal(t, expected, builder.Body)
		})

		t.Run("invalid struct", func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("Expected panic for invalid struct")
				}
			}()
			utils.NewRequestBuilder(testToken).SetBodyFromStruct(make(chan int))
		})
	})

	testRequestConstruction := func(method string, buildMethod func(*utils.AuthedRequest, context.Context) (*http.Request, error)) {
		t.Run(fmt.Sprintf("Build%s", method), func(t *testing.T) {
			testBody := struct{ Message string }{Message: "test"}
			expectedBody, _ := json.Marshal(testBody)

			builder := utils.NewRequestBuilder(testToken).
				SetBaseURL(testBaseURL).
				SetPath(testPath).
				SetBodyFromStruct(testBody)

			ctx := context.WithValue(context.Background(), testContextKey, "value")
			req, err := buildMethod(builder, ctx)

			require.NoError(t, err)
			assert.Equal(t, method, req.Method)
			assert.Equal(t, testBaseURL+testPath, req.URL.String())
			assert.Equal(t, "Bearer "+testToken, req.Header.Get("Authorization"))
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))

			// Verify body
			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(req.Body)
			require.NoError(t, err)
			assert.JSONEq(t, string(expectedBody), buf.String())

			// Verify context
			assert.Equal(t, "value", req.Context().Value(testContextKey))
		})
	}

	// Test different build methods
	testRequestConstruction(http.MethodPost, func(b *utils.AuthedRequest, ctx context.Context) (*http.Request, error) {
		return b.Build(ctx)
	})

	testRequestConstruction(http.MethodGet, func(b *utils.AuthedRequest, ctx context.Context) (*http.Request, error) {
		return b.BuildGet(ctx)
	})

	testRequestConstruction(http.MethodPost, func(b *utils.AuthedRequest, ctx context.Context) (*http.Request, error) {
		return b.BuildStream(ctx)
	})

	t.Run("BuildStream", func(t *testing.T) {
		req, err := utils.NewRequestBuilder(testToken).
			SetBaseURL(testBaseURL).
			SetPath(testPath).
			BuildStream(context.Background())

		require.NoError(t, err)
		assert.Equal(t, "no-cache", req.Header.Get("cache-control"))
	})

	t.Run("ErrorConditions", func(t *testing.T) {
		tests := []struct {
			name        string
			baseURL     string
			path        string
			expectedErr string
		}{
			{"MissingBaseURL", "", "/path", "BaseURL or path not set"},
			{"MissingPath", "http://example.com", "", "BaseURL or path not set"},
			{"BothMissing", "", "", "BaseURL or path not set"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				builder := utils.NewRequestBuilder(testToken).
					SetBaseURL(tt.baseURL).
					SetPath(tt.path)

				_, err := builder.Build(context.Background())
				assert.ErrorContains(t, err, tt.expectedErr)

				_, err = builder.BuildGet(context.Background())
				assert.ErrorContains(t, err, tt.expectedErr)

				_, err = builder.BuildStream(context.Background())
				assert.ErrorContains(t, err, tt.expectedErr)
			})
		}
	})

	t.Run("EdgeCases", func(t *testing.T) {
		t.Run("EmptyBody", func(t *testing.T) {
			req, err := utils.NewRequestBuilder(testToken).
				SetBaseURL(testBaseURL).
				SetPath(testPath).
				Build(context.Background())

			require.NoError(t, err)
			assert.IsType(t, http.NoBody, req.Body)
		})

		t.Run("URLConstruction", func(t *testing.T) {
			testCases := []struct {
				baseURL  string
				path     string
				expected string
			}{
				{"http://host", "/path", "http://host/path"},
				{"http://host/", "/path", "http://host//path"},
				{"http://host", "path", "http://hostpath"},
			}

			for _, tc := range testCases {
				t.Run(tc.expected, func(t *testing.T) {
					req, err := utils.NewRequestBuilder(testToken).
						SetBaseURL(tc.baseURL).
						SetPath(tc.path).
						Build(context.Background())

					require.NoError(t, err)
					assert.Equal(t, tc.expected, req.URL.String())
				})
			}
		})
	})
}

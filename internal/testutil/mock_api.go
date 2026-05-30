package testutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	deepseek "github.com/cohesion-org/deepseek-go"
)

type MockDeepSeekServer struct {
	*httptest.Server
}

// NewMockDeepSeekServer creates a local HTTP server that mimics the small subset of DeepSeek endpoints used by tests.
func NewMockDeepSeekServer(t *testing.T) *MockDeepSeekServer {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch cleanPath(r.URL.Path) {
		case "/chat/completions", "/beta/chat/completions":
			handleMockChatCompletions(t, w, r)
		case "/models":
			WriteJSON(t, w, http.StatusOK, map[string]interface{}{
				"object": "list",
				"data": []map[string]string{
					{"id": deepseek.DeepSeekChat, "object": "model", "owned_by": "deepseek"},
					{"id": deepseek.DeepSeekReasoner, "object": "model", "owned_by": "deepseek"},
					{"id": deepseek.DeepSeekCoder, "object": "model", "owned_by": "deepseek"},
					{"id": deepseek.DeepSeekV4Flash, "object": "model", "owned_by": "deepseek"},
					{"id": deepseek.DeepSeekV4Pro, "object": "model", "owned_by": "deepseek"},
				},
			})
		case "/user/balance":
			WriteJSON(t, w, http.StatusOK, map[string]interface{}{
				"is_available": true,
				"balance_infos": []map[string]string{
					{
						"currency":          "USD",
						"total_balance":     "10.00",
						"granted_balance":   "2.00",
						"topped_up_balance": "8.00",
					},
				},
			})
		case "/beta/completions":
			handleMockFIMCompletions(t, w, r)
		default:
			http.NotFound(w, r)
		}
	}))

	return &MockDeepSeekServer{Server: server}
}

// BaseURL returns the mock server base URL for standard API endpoints.
func (s *MockDeepSeekServer) BaseURL() string {
	return s.URL + "/"
}

// BetaBaseURL returns the mock server base URL for beta endpoints.
func (s *MockDeepSeekServer) BetaBaseURL() string {
	return s.URL + "/beta/"
}

// NewMockClient creates a client configured to talk to the mock server's standard API endpoints.
func NewMockClient(t *testing.T, server *MockDeepSeekServer, opts ...deepseek.Option) *deepseek.Client {
	t.Helper()

	allOpts := append([]deepseek.Option{deepseek.WithBaseURL(server.BaseURL())}, opts...)
	client, err := deepseek.NewClientWithOptions("test-token", allOpts...)
	if err != nil {
		t.Fatalf("failed to create mock client: %v", err)
	}

	return client
}

// NewMockBetaClient creates a client configured to talk to the mock server's beta API endpoints.
func NewMockBetaClient(t *testing.T, server *MockDeepSeekServer, opts ...deepseek.Option) *deepseek.Client {
	t.Helper()

	allOpts := append([]deepseek.Option{deepseek.WithBaseURL(server.BetaBaseURL())}, opts...)
	client, err := deepseek.NewClientWithOptions("test-token", allOpts...)
	if err != nil {
		t.Fatalf("failed to create mock beta client: %v", err)
	}

	return client
}

// WriteJSON writes a JSON response and fails the test immediately if encoding fails.
func WriteJSON(t *testing.T, w http.ResponseWriter, status int, body interface{}) {
	t.Helper()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		t.Fatalf("failed to write json response: %v", err)
	}
}

// WriteSSE writes a simple server-sent event stream terminated by a `[DONE]` frame.
func WriteSSE(t *testing.T, w http.ResponseWriter, chunks ...string) {
	t.Helper()

	w.Header().Set("Content-Type", "text/event-stream")
	for _, chunk := range chunks {
		if _, err := fmt.Fprintf(w, "data: %s\n\n", chunk); err != nil {
			t.Fatalf("failed to write sse chunk: %v", err)
		}
	}
	if _, err := fmt.Fprint(w, "data: [DONE]\n\n"); err != nil {
		t.Fatalf("failed to write sse terminator: %v", err)
	}
}

// SkipUnlessLiveAPI skips tests unless the live API opt-in flag is enabled.
func SkipUnlessLiveAPI(t *testing.T) {
	t.Helper()

	if osValue := strings.TrimSpace(strings.ToLower(getenv("DEEPSEEK_LIVE_TESTS"))); osValue != "1" && osValue != "true" && osValue != "yes" {
		t.Skip("skipping live API test; set DEEPSEEK_LIVE_TESTS=1 to enable")
	}
}

type mockChatRequest struct {
	Model    string            `json:"model"`
	Messages []mockChatMessage `json:"messages"`
	Stream   bool              `json:"stream,omitempty"`
	Tools    []deepseek.Tool   `json:"tools,omitempty"`
	Thinking map[string]any    `json:"thinking,omitempty"`
}

type mockChatMessage struct {
	Role             string              `json:"role"`
	Content          string              `json:"content"`
	Prefix           bool                `json:"prefix,omitempty"`
	ReasoningContent string              `json:"reasoning_content,omitempty"`
	ToolCallID       string              `json:"tool_call_id,omitempty"`
	ToolCalls        []deepseek.ToolCall `json:"tool_calls,omitempty"`
}

type mockFIMRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	Temperature float64 `json:"temperature,omitempty"`
}

func handleMockChatCompletions(t *testing.T, w http.ResponseWriter, r *http.Request) {
	t.Helper()

	var req mockChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(t, w, http.StatusBadRequest, map[string]interface{}{
			"error": map[string]string{"message": "invalid request body"},
		})
		return
	}

	if len(req.Messages) == 0 {
		WriteJSON(t, w, http.StatusBadRequest, map[string]interface{}{
			"error": map[string]string{"message": "messages must not be empty"},
		})
		return
	}

	if req.Model == "invalid-model-123" {
		WriteJSON(t, w, http.StatusBadRequest, map[string]interface{}{
			"error": map[string]string{"message": "model not found"},
		})
		return
	}

	if req.Stream {
		writeMockChatStream(t, w, req)
		return
	}

	content, reasoning, toolCalls := buildMockChatResult(req)
	WriteJSON(t, w, http.StatusOK, map[string]interface{}{
		"id":      "chatcmpl-mock",
		"object":  "chat.completion",
		"created": time.Now().Unix(),
		"model":   req.Model,
		"choices": []map[string]interface{}{
			{
				"index": 0,
				"message": map[string]interface{}{
					"role":              deepseek.ChatMessageRoleAssistant,
					"content":           content,
					"reasoning_content": reasoning,
					"tool_calls":        toolCalls,
				},
				"finish_reason": finishReasonForToolCalls(toolCalls),
			},
		},
		"usage": map[string]int{
			"prompt_tokens":     8,
			"completion_tokens": 6,
			"total_tokens":      14,
		},
	})
}

func handleMockFIMCompletions(t *testing.T, w http.ResponseWriter, r *http.Request) {
	t.Helper()

	var req mockFIMRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(t, w, http.StatusBadRequest, map[string]interface{}{
			"error": map[string]string{"message": "invalid request body"},
		})
		return
	}

	switch {
	case req.Prompt == "":
		WriteJSON(t, w, http.StatusBadRequest, map[string]interface{}{
			"error": map[string]string{"message": "prompt must not be empty"},
		})
		return
	case req.Model == "invalid-model-123":
		WriteJSON(t, w, http.StatusBadRequest, map[string]interface{}{
			"error": map[string]string{"message": "model not found"},
		})
		return
	case req.Temperature > 2:
		WriteJSON(t, w, http.StatusBadRequest, map[string]interface{}{
			"error": map[string]string{"message": "temperature out of range"},
		})
		return
	}

	WriteJSON(t, w, http.StatusOK, map[string]interface{}{
		"id":      "fim-mock",
		"object":  "text_completion",
		"created": time.Now().Unix(),
		"model":   req.Model,
		"choices": []map[string]interface{}{
			{
				"text":          "lo\")\n}",
				"index":         0,
				"finish_reason": "stop",
				"logprobs":      map[string]interface{}{"content": []interface{}{}},
			},
		},
		"usage": map[string]int{
			"prompt_tokens":     4,
			"completion_tokens": 3,
			"total_tokens":      7,
		},
	})
}

func writeMockChatStream(t *testing.T, w http.ResponseWriter, req mockChatRequest) {
	t.Helper()

	content, reasoning, toolCalls := buildMockChatResult(req)
	now := time.Now().Unix()

	if len(toolCalls) > 0 {
		arguments := toolCalls[0].Function.Arguments
		mid := len(arguments)
		if mid > 4 {
			mid = len(arguments) / 2
		}

		WriteSSE(t, w,
			fmt.Sprintf(`{"id":"chatcmpl-stream","object":"chat.completion.chunk","created":%d,"model":"%s","choices":[{"index":0,"delta":{"role":"assistant"},"finish_reason":""}]}`, now, req.Model),
			fmt.Sprintf(`{"id":"chatcmpl-stream","object":"chat.completion.chunk","created":%d,"model":"%s","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"%s","type":"function","function":{"name":"%s","arguments":%q}}]},"finish_reason":""}]}`, now, req.Model, toolCalls[0].ID, toolCalls[0].Function.Name, arguments[:mid]),
			fmt.Sprintf(`{"id":"chatcmpl-stream","object":"chat.completion.chunk","created":%d,"model":"%s","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"function":{"arguments":%q}}]},"finish_reason":"tool_calls"}]}`, now, req.Model, arguments[mid:]),
		)
		return
	}

	chunks := []string{
		fmt.Sprintf(`{"id":"chatcmpl-stream","object":"chat.completion.chunk","created":%d,"model":"%s","choices":[{"index":0,"delta":{"role":"assistant"},"finish_reason":""}]}`, now, req.Model),
	}

	if reasoning != "" {
		chunks = append(chunks, fmt.Sprintf(`{"id":"chatcmpl-stream","object":"chat.completion.chunk","created":%d,"model":"%s","choices":[{"index":0,"delta":{"reasoning_content":%q},"finish_reason":""}]}`, now, req.Model, reasoning))
	}

	chunks = append(chunks,
		fmt.Sprintf(`{"id":"chatcmpl-stream","object":"chat.completion.chunk","created":%d,"model":"%s","choices":[{"index":0,"delta":{"content":%q},"finish_reason":""}],"usage":{"prompt_tokens":4,"completion_tokens":2,"total_tokens":6}}`, now, req.Model, content),
		fmt.Sprintf(`{"id":"chatcmpl-stream","object":"chat.completion.chunk","created":%d,"model":"%s","choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}`, now, req.Model),
	)

	WriteSSE(t, w, chunks...)
}

func buildMockChatResult(req mockChatRequest) (content string, reasoning string, toolCalls []deepseek.ToolCall) {
	lastUser := lastMessageContent(req.Messages, deepseek.ChatMessageRoleUser)
	thinkingMode := req.Model == deepseek.DeepSeekReasoner || req.Model == deepseek.DeepSeekV4Pro || req.Thinking != nil
	hasToolMessage := hasRole(req.Messages, deepseek.ChatMessageRoleTool)
	hasPrefix := hasPrefixMessage(req.Messages)

	switch {
	case len(req.Tools) > 0 && !hasToolMessage:
		if thinkingMode {
			reasoning = "I should inspect the request and call the appropriate tool first."
		}
		name := req.Tools[0].Function.Name
		args := `{"location":"Paris","date":"2026-03-24"}`
		if strings.Contains(strings.ToLower(name), "location") {
			args = `{"country":"France","city":"Paris"}`
		}
		return "", reasoning, []deepseek.ToolCall{
			{
				Index: 0,
				ID:    "call_mock_1",
				Type:  "function",
				Function: deepseek.ToolCallFunction{
					Name:      name,
					Arguments: args,
				},
			},
		}
	case len(req.Tools) > 0 && hasToolMessage:
		if thinkingMode {
			reasoning = "I now have the tool result and can answer the user."
		}
		return "The weather is 24C and clear.", reasoning, nil
	case hasPrefix:
		if thinkingMode || hasReasoningInput(req.Messages) {
			return "def quick_sort(arr):\n    return sorted(arr)", "I will continue the requested code block.", nil
		}
		return "def quick_sort(arr):\n    return sorted(arr)", "", nil
	case strings.Contains(strings.ToLower(lastUser), "highest mountain"):
		return "Everest", reasoning, nil
	case strings.Contains(strings.ToLower(lastUser), "second highest"):
		return "K2", reasoning, nil
	case strings.Contains(strings.ToLower(lastUser), "current president"):
		return "CurrentPresident", reasoning, nil
	case strings.Contains(strings.ToLower(lastUser), "immediate predecessor"):
		return "PreviousPresident", reasoning, nil
	case strings.Contains(strings.ToLower(lastUser), "joke"):
		return "AI told a joke about tensors and it still landed.", reasoning, nil
	case thinkingMode:
		return "Human oversight should remain in place for final diagnostic decisions.", "The stakes are high, so oversight helps catch errors and edge cases.", nil
	default:
		return "Hello from the mock DeepSeek server.", "", nil
	}
}

func finishReasonForToolCalls(toolCalls []deepseek.ToolCall) string {
	if len(toolCalls) > 0 {
		return "tool_calls"
	}

	return "stop"
}

func cleanPath(p string) string {
	if p == "" {
		return "/"
	}

	return path.Clean("/" + p)
}

func lastMessageContent(messages []mockChatMessage, role string) string {
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == role {
			return messages[i].Content
		}
	}

	return ""
}

func hasRole(messages []mockChatMessage, role string) bool {
	for _, message := range messages {
		if message.Role == role {
			return true
		}
	}

	return false
}

func hasPrefixMessage(messages []mockChatMessage) bool {
	for _, message := range messages {
		if message.Prefix {
			return true
		}
	}

	return false
}

func hasReasoningInput(messages []mockChatMessage) bool {
	for _, message := range messages {
		if message.ReasoningContent != "" {
			return true
		}
	}

	return false
}

func getenv(key string) string {
	return os.Getenv(key)
}

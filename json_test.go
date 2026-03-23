package deepseek

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJSONExtractor(t *testing.T) {
	schema := json.RawMessage(`{"type": "object"}`)
	extractor := NewJSONExtractor(schema)

	require.NotNil(t, extractor, "NewJSONExtractor returned nil")
	assert.Equal(t, schema, extractor.schema, "Schema not properly set in constructor")
}

func TestExtractJSON(t *testing.T) {
	tests := []struct {
		name        string
		response    *ChatCompletionResponse
		schema      json.RawMessage
		target      interface{}
		expectError bool
		expectedVal interface{}
	}{
		{
			name:        "Nil Response",
			response:    nil,
			schema:      nil,
			target:      &struct{}{},
			expectError: true,
		},
		{
			name: "Empty Choices",
			response: &ChatCompletionResponse{
				Choices: []Choice{},
			},
			schema:      nil,
			target:      &struct{}{},
			expectError: true,
		},
		{
			name: "Valid JSON Response",
			response: &ChatCompletionResponse{
				Choices: []Choice{
					{
						Message: Message{
							Content: `{"name": "test"}`,
						},
					},
				},
			},
			schema: nil,
			target: &struct {
				Name string `json:"name"`
			}{},
			expectError: false,
			expectedVal: struct {
				Name string `json:"name"`
			}{Name: "test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			extractor := NewJSONExtractor(tt.schema)
			err := extractor.ExtractJSON(tt.response, tt.target)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedVal, *tt.target.(*struct {
					Name string `json:"name"`
				}))
			}
		})
	}
}

func TestValidateJSON(t *testing.T) {
	tests := []struct {
		name        string
		schema      json.RawMessage
		data        []byte
		expectError bool
	}{
		{
			name:        "Valid Object Schema",
			schema:      json.RawMessage(`{"type": "object"}`),
			data:        []byte(`{"name": "test"}`),
			expectError: false,
		},
		{
			name:        "Invalid Object Schema",
			schema:      json.RawMessage(`{"type": "object"}`),
			data:        []byte(`["not an object"]`),
			expectError: true,
		},
		{
			name:        "Valid Array Schema",
			schema:      json.RawMessage(`{"type": "array"}`),
			data:        []byte(`[1,2,3]`),
			expectError: false,
		},
		{
			name:        "Invalid Schema",
			schema:      json.RawMessage(`{invalid`),
			data:        []byte(`{}`),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractor := NewJSONExtractor(tt.schema)
			err := extractor.validateJSON(tt.data)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestExtractJSONContent(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "Plain JSON",
			content:  `{"name": "test"}`,
			expected: `{"name": "test"}`,
		},
		{
			name: "JSON in Code Block",
			content: "```json\n" +
				`{"name": "test"}` +
				"\n```",
			expected: `{"name": "test"}`,
		},
		{
			name:     "Invalid JSON",
			content:  "This is not JSON",
			expected: "",
		},
		{
			name:     "JSON Array in Text",
			content:  "Some text before [1,2,3] and after",
			expected: "[1,2,3]",
		},
	}

	extractor := NewJSONExtractor(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractor.extractJSONContent(tt.content)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFindMatchingBrace(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected int
	}{
		{
			name:     "Simple Object",
			content:  `{"name": "test"}`,
			expected: 15,
		},
		{
			name:     "Nested Object",
			content:  `{"outer": {"inner": "value"}}`,
			expected: 28,
		},
		{
			name:     "Invalid Start",
			content:  `["not an object"]`,
			expected: -1,
		},
		{
			name:     "Unclosed Brace",
			content:  `{"unclosed": "brace"`,
			expected: -1,
		},
	}

	extractor := NewJSONExtractor(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractor.findMatchingBrace(tt.content)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFindMatchingBracket(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected int
	}{
		{
			name:     "Simple Array",
			content:  `[1,2,3]`,
			expected: 6,
		},
		{
			name:     "Nested Array",
			content:  `[[1,2],[3,4]]`,
			expected: 12,
		},
		{
			name:     "Invalid Start",
			content:  `{"not": "array"}`,
			expected: -1,
		},
		{
			name:     "Unclosed Bracket",
			content:  `[1,2,3`,
			expected: -1,
		},
	}

	extractor := NewJSONExtractor(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractor.findMatchingBracket(tt.content)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestJSONExtractionEdgeCases(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	extractor := NewJSONExtractor(nil)

	tests := []struct {
		name     string
		content  string
		expected map[string]interface{}
	}{
		{
			name: "JSON with Unicode",
			content: `{
				"name": "José garcía",
				"symbol": "☺"
			}`,
			expected: map[string]interface{}{
				"name":   "José garcía",
				"symbol": "☺",
			},
		},
		{
			name: "JSON with Special Characters",
			content: `{
				"text": "Hello\nWorld\t!",
				"quotes": "\"quoted\""
			}`,
			expected: map[string]interface{}{
				"text":   "Hello\nWorld\t!",
				"quotes": "\"quoted\"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := &ChatCompletionResponse{
				Choices: []Choice{
					{
						Message: Message{
							Content: tt.content,
						},
					},
				},
			}

			var result map[string]interface{}
			err := extractor.ExtractJSON(response, &result)
			require.NoError(err)
			assert.Equal(tt.expected, result)
		})
	}
}

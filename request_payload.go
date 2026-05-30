package deepseek

import "encoding/json"

// ThinkingConfig configures DeepSeek's thinking mode.
type ThinkingConfig struct {
	Type string `json:"type"`
}

func normalizeChatCompletionPayload(request *ChatCompletionRequest) (map[string]interface{}, error) {
	payload := struct {
		Model            string                  `json:"model"`
		Messages         []ChatCompletionMessage `json:"messages"`
		FrequencyPenalty float32                 `json:"frequency_penalty,omitempty"`
		MaxTokens        int                     `json:"max_tokens,omitempty"`
		PresencePenalty  float32                 `json:"presence_penalty,omitempty"`
		Temperature      float32                 `json:"temperature,omitempty"`
		TopP             float32                 `json:"top_p,omitempty"`
		ResponseFormat   *ResponseFormat         `json:"response_format,omitempty"`
		Stop             []string                `json:"stop,omitempty"`
		Tools            []Tool                  `json:"tools,omitempty"`
		ToolChoice       interface{}             `json:"tool_choice,omitempty"`
		LogProbs         bool                    `json:"logprobs,omitempty"`
		TopLogProbs      int                     `json:"top_logprobs,omitempty"`
		ReasoningEffort  string                  `json:"reasoning_effort,omitempty"`
		UserID           string                  `json:"user_id,omitempty"`
		JSONMode         bool                    `json:"json,omitempty"`
		Thinking         *ThinkingConfig         `json:"thinking,omitempty"`
	}{
		Model:            request.Model,
		Messages:         request.Messages,
		FrequencyPenalty: request.FrequencyPenalty,
		MaxTokens:        request.MaxTokens,
		PresencePenalty:  request.PresencePenalty,
		Temperature:      request.Temperature,
		TopP:             request.TopP,
		ResponseFormat:   request.ResponseFormat,
		Stop:             request.Stop,
		Tools:            request.Tools,
		ToolChoice:       request.ToolChoice,
		LogProbs:         request.LogProbs,
		TopLogProbs:      request.TopLogProbs,
		ReasoningEffort:  request.ReasoningEffort,
		UserID:           request.UserID,
		JSONMode:         request.JSONMode,
		Thinking:         request.Thinking,
	}

	return finalizePayload(payload, request.ExtraFields, request.EnableThinking)
}

func normalizeStreamChatCompletionPayload(request *StreamChatCompletionRequest) (map[string]interface{}, error) {
	payload := struct {
		Stream           bool                    `json:"stream,omitempty"`
		StreamOptions    StreamOptions           `json:"stream_options,omitempty"`
		Model            string                  `json:"model"`
		Messages         []ChatCompletionMessage `json:"messages"`
		FrequencyPenalty float32                 `json:"frequency_penalty,omitempty"`
		MaxTokens        int                     `json:"max_tokens,omitempty"`
		PresencePenalty  float32                 `json:"presence_penalty,omitempty"`
		Temperature      float32                 `json:"temperature,omitempty"`
		TopP             float32                 `json:"top_p,omitempty"`
		ResponseFormat   *ResponseFormat         `json:"response_format,omitempty"`
		Stop             []string                `json:"stop,omitempty"`
		Tools            []Tool                  `json:"tools,omitempty"`
		ToolChoice       interface{}             `json:"tool_choice,omitempty"`
		LogProbs         bool                    `json:"logprobs,omitempty"`
		TopLogProbs      int                     `json:"top_logprobs,omitempty"`
		ReasoningEffort  string                  `json:"reasoning_effort,omitempty"`
		UserID           string                  `json:"user_id,omitempty"`
		Thinking         *ThinkingConfig         `json:"thinking,omitempty"`
	}{
		Stream:           request.Stream,
		StreamOptions:    request.StreamOptions,
		Model:            request.Model,
		Messages:         request.Messages,
		FrequencyPenalty: request.FrequencyPenalty,
		MaxTokens:        request.MaxTokens,
		PresencePenalty:  request.PresencePenalty,
		Temperature:      request.Temperature,
		TopP:             request.TopP,
		ResponseFormat:   request.ResponseFormat,
		Stop:             request.Stop,
		Tools:            request.Tools,
		ToolChoice:       request.ToolChoice,
		LogProbs:         request.LogProbs,
		TopLogProbs:      request.TopLogProbs,
		ReasoningEffort:  request.ReasoningEffort,
		UserID:           request.UserID,
		Thinking:         request.Thinking,
	}

	return finalizePayload(payload, request.ExtraFields, request.EnableThinking)
}

func finalizePayload(base interface{}, extraFields map[string]interface{}, enableThinking bool) (map[string]interface{}, error) {
	body, err := json.Marshal(base)
	if err != nil {
		return nil, err
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	for key, value := range extraFields {
		payload[key] = value
	}

	if _, hasThinking := payload["thinking"]; !hasThinking && enableThinking {
		payload["thinking"] = map[string]string{"type": "enabled"}
	}

	return payload, nil
}

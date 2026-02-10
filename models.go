package deepseek

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	utils "github.com/cohesion-org/deepseek-go/utils"
)

// Official DeepSeek Models
const (
	DeepSeekChat     = "deepseek-chat"     // DeepSeekChat is the official model for chat completions
	DeepSeekCoder    = "deepseek-coder"    // DeepSeekCoder has been combined with DeepSeekChat, but you can still use it. Please read: https://api-docs.deepseek.com/updates#version-2024-09-05
	DeepSeekReasoner = "deepseek-reasoner" // DeepSeekReasoner is the official model for reasoning completions
)

// External Models that can be used with the API
const (
	AzureDeepSeekR1                     = "DeepSeek-R1"                            // Azure model for DeepSeek R1
	OpenRouterDeepSeekR1                = "deepseek/deepseek-r1"                   // OpenRouter model for DeepSeek R1
	OpenRouterDeepSeekR1DistillLlama70B = "deepseek/deepseek-r1-distill-llama-70b" // DeepSeek R1 Distill Llama 70B
	OpenRouterDeepSeekR1DistillLlama8B  = "deepseek/deepseek-r1-distill-llama-8b"  // DeepSeek R1 Distill Llama 8B
	OpenRouterDeepSeekR1DistillQwen14B  = "deepseek/deepseek-r1-distill-qwen-14b"  // DeepSeek R1 Distill Qwen 14B
	OpenRouterDeepSeekR1DistillQwen1_5B = "deepseek/deepseek-r1-distill-qwen-1.5b" // DeepSeek R1 Distill Qwen 1.5B
	OpenRouterDeepSeekR1DistillQwen32B  = "deepseek/deepseek-r1-distill-qwen-32b"  // DeepSeek R1 Distill Qwen 32B
)

// Model represents a model that can be used with the API
type Model struct {
	ID      string `json:"id"`       //The id of the model (string)
	Object  string `json:"object"`   //The object of the model (string)
	OwnedBy string `json:"owned_by"` //The owner of the model(usually deepseek)
}

// APIModels represents the response from the API endpoint.
type APIModels struct {
	Object string  `json:"object"` // Object (string)
	Data   []Model `json:"data"`   // List of Models
}

// ListAllModels sends a request to the API to get all available models.
func ListAllModels(c *Client, ctx context.Context) (*APIModels, error) {
	req, err := utils.NewRequestBuilder(c.AuthToken).
		SetBaseURL("https://api.deepseek.com/").
		SetPath("models").
		BuildGet(ctx)

	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}

	resp, err := HandleNormalRequest(*c, req)

	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, HandleAPIError(resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var models APIModels
	if err := json.Unmarshal(body, &models); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}
	return &models, nil
}

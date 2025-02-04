package deepseek

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	handlers "github.com/cohesion-org/deepseek-go/handlers"
	utils "github.com/cohesion-org/deepseek-go/utils"
)

// Create chat and coder models here!

const (
	DeepSeekChat     = "deepseek-chat"
	DeepSeekCoder    = "deepseek-coder"
	DeepSeekReasoner = "deepseek-reasoner"
	AzureDeepSeekR1  = "DeepSeek-R1"
)

type Model struct {
	ID      string `json:"id"`       //The id of the model (string)
	Object  string `json:"object"`   //The object of the model (string)
	OwnedBy string `json:"owned_by"` //The owner of the model(usally deepseek)
}

type APIModels struct {
	Object string  `json:"object"` //Object (string)
	Data   []Model `json:"data"`   // List of Models
}

// Models supported by the API itself
func ListAllModels(c *Client, ctx context.Context) (*APIModels, error) {
	req, err := utils.NewRequestBuilder(c.AuthToken).
		SetBaseURL("https://api.deepseek.com/").
		SetPath("models").
		BuildGet(ctx)

	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}

	resp, err := handlers.HandleNormalRequest(req)

	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, HandleAPIError(resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var models APIModels
	if err := json.Unmarshal(body, &models); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}
	return &models, nil
}

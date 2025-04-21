package deepseek

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	utils "github.com/cohesion-org/deepseek-go/utils"
)

// BalanceInfo represents the balance information for a specific currency.
type BalanceInfo struct {
	Currency        string `json:"currency"`          //The currency of the balance.
	TotalBalance    string `json:"total_balance"`     //The total available balance, including the granted balance and the topped-up balance.
	GrantedBalance  string `json:"granted_balance"`   //The total not expired granted balance.
	ToppedUpBalance string `json:"topped_up_balance"` //The total topped-up balance.
}

// BalanceResponse represents the response from the API endpoint.
type BalanceResponse struct {
	IsAvailable  bool          `json:"is_available"`  //Whether the user's balance is sufficient for API calls.
	BalanceInfos []BalanceInfo `json:"balance_infos"` //List of Balance infos
}

// GetBalance sends a request to the API to get the user's balance.
func GetBalance(c *Client, ctx context.Context) (*BalanceResponse, error) {

	req, err := utils.NewRequestBuilder(c.AuthToken).
		SetBaseURL("https://api.deepseek.com/").
		SetPath("user/balance").
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

	var balance BalanceResponse
	if err := json.Unmarshal(body, &balance); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}
	return &balance, nil
}

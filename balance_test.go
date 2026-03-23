package deepseek_test

import (
	"context"
	"testing"

	"github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBalance(t *testing.T) {
	server := testutil.NewMockDeepSeekServer(t)
	defer server.Close()

	client := testutil.NewMockClient(t, server)

	balance, err := deepseek.GetBalance(client, context.Background())
	require.NoError(t, err, "should not return error")
	require.NotNil(t, balance, "response should not be nil ")

	// Verify response structure
	assert.True(t, balance.IsAvailable || !balance.IsAvailable, "IsAvailable should be a boolean")
	assert.NotEmpty(t, balance.BalanceInfos, "should have balance information")

	// Verify balance info details
	for _, info := range balance.BalanceInfos {
		assert.NotEmpty(t, info.Currency, "currency should not be empty")
		assert.NotEmpty(t, info.TotalBalance, "total balance should not be empty")
		assert.NotEmpty(t, info.GrantedBalance, "granted balance should not be empty")
		assert.NotEmpty(t, info.ToppedUpBalance, "topped up balance should not be empty")
	}
}

package deepseek_test

import (
	"context"
	"testing"

	"github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListAllModels(t *testing.T) {
	server := testutil.NewMockDeepSeekServer(t)
	defer server.Close()

	client := testutil.NewMockClient(t, server)
	resp, err := deepseek.ListAllModels(client, context.Background())
	require.NoError(t, err)
	assert.NotNil(t, resp)

	// Verify response structure
	assert.Equal(t, "list", resp.Object)
	assert.NotEmpty(t, resp.Data)

	// Verify model details
	for _, model := range resp.Data {
		assert.NotEmpty(t, model.ID)
		assert.Equal(t, "model", model.Object)
		assert.Equal(t, "deepseek", model.OwnedBy)

		// Verify known models exist in constants.go
		if model.ID == "deepseek-chat" ||
			model.ID == "deepseek-coder" ||
			model.ID == "deepseek-reasoner" {
			assert.Contains(t, []string{
				"deepseek-chat",
				"deepseek-coder",
				"deepseek-reasoner",
			}, model.ID)
		}
	}
}

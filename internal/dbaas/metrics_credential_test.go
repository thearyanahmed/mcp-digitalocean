package dbaas

import (
	"context"
	"encoding/json"
	"testing"

	"mcp-digitalocean/internal/dbaas/mocks"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMetricCredentialTool_getMetricsCredentials(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("GetMetricsCredentials", mock.Anything).Return(&godo.DatabaseMetricsCredentials{BasicAuthUsername: "user"}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	mt := &MetricCredentialTool{client: client}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{}}}
	res, err := mt.getMetricsCredentials(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "user")
	mockDB.AssertExpectations(t)
}

func TestMetricCredentialTool_updateMetricsCredentials(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("UpdateMetricsCredentials", mock.Anything, mock.Anything).Return(nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	mt := &MetricCredentialTool{client: client}
	creds := godo.DatabaseMetricsCredentials{BasicAuthUsername: "user"}
	credsJSON, _ := json.Marshal(creds)
	args := map[string]interface{}{"credentials_json": string(credsJSON)}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := mt.updateMetricsCredentials(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Metrics credentials updated successfully")
	mockDB.AssertExpectations(t)
	// Error case: missing credentials_json
	args = map[string]interface{}{}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = mt.updateMetricsCredentials(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "credentials_json is required")
}

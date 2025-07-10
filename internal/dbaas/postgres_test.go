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

func TestPostgreSQLTool_getPostgreSQLConfig(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	val := 42
	mockDB.On("GetPostgreSQLConfig", mock.Anything, "cid").Return(&godo.PostgreSQLConfig{BackupHour: &val}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	pt := &PostgreSQLTool{client: client}
	args := map[string]interface{}{"ID": "cid"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := pt.getPostgreSQLConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "42")
	mockDB.AssertExpectations(t)
	// Error case: missing ID
	reqMissing := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{}}}
	res, err = pt.getPostgreSQLConfig(context.Background(), reqMissing)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Cluster ID is required")
	// API error
	mockDB.On("GetPostgreSQLConfig", mock.Anything, "badid").Return(nil, nil, assert.AnError)
	args = map[string]interface{}{"ID": "badid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = pt.getPostgreSQLConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "api error")
}

func TestPostgreSQLTool_updatePostgreSQLConfig(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	val := 99
	mockDB.On("UpdatePostgreSQLConfig", mock.Anything, "cid", mock.Anything).Return(&godo.Response{}, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	pt := &PostgreSQLTool{client: client}
	config := godo.PostgreSQLConfig{BackupHour: &val}
	configBytes, _ := json.Marshal(config)
	args := map[string]interface{}{"ID": "cid", "config_json": string(configBytes)}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := pt.updatePostgreSQLConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "updated successfully")
	mockDB.AssertExpectations(t)
	// Error case: missing ID
	args = map[string]interface{}{"config_json": string(configBytes)}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = pt.updatePostgreSQLConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Cluster ID is required")
	// Error case: missing config_json
	args = map[string]interface{}{"ID": "cid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = pt.updatePostgreSQLConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "config_json is required")
	// Error case: invalid config_json
	args = map[string]interface{}{"ID": "cid", "config_json": "notjson"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = pt.updatePostgreSQLConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Invalid config_json")
	// API error
	mockDB.On("UpdatePostgreSQLConfig", mock.Anything, "badid", mock.Anything).Return(nil, assert.AnError)
	args = map[string]interface{}{"ID": "badid", "config_json": string(configBytes)}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = pt.updatePostgreSQLConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "api error")
}

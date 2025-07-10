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

func TestOpenSearchTool_getOpensearchConfig(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	val := 12345
	mockDB.On("GetOpensearchConfig", mock.Anything, "cid").Return(&godo.OpensearchConfig{HttpMaxContentLengthBytes: &val}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	ot := &OpenSearchTool{client: client}
	args := map[string]interface{}{"ID": "cid"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := ot.getOpensearchConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "12345")
	mockDB.AssertExpectations(t)
	// Error case: missing ID
	reqMissing := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{}}}
	res, err = ot.getOpensearchConfig(context.Background(), reqMissing)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Cluster ID is required")
	// API error
	mockDB.On("GetOpensearchConfig", mock.Anything, "badid").Return(nil, nil, assert.AnError)
	args = map[string]interface{}{"ID": "badid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = ot.getOpensearchConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "api error")
}

func TestOpenSearchTool_updateOpensearchConfig(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	val := 54321
	mockDB.On("UpdateOpensearchConfig", mock.Anything, "cid", mock.Anything).Return(&godo.Response{}, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	ot := &OpenSearchTool{client: client}
	config := godo.OpensearchConfig{HttpMaxContentLengthBytes: &val}
	configBytes, _ := json.Marshal(config)
	args := map[string]interface{}{"ID": "cid", "config_json": string(configBytes)}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := ot.updateOpensearchConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "updated successfully")
	mockDB.AssertExpectations(t)
	// Error case: missing ID
	args = map[string]interface{}{"config_json": string(configBytes)}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = ot.updateOpensearchConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Cluster ID is required")
	// Error case: missing config_json
	args = map[string]interface{}{"ID": "cid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = ot.updateOpensearchConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "config_json is required")
	// Error case: invalid config_json
	args = map[string]interface{}{"ID": "cid", "config_json": "notjson"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = ot.updateOpensearchConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Invalid config_json")
	// API error
	mockDB.On("UpdateOpensearchConfig", mock.Anything, "badid", mock.Anything).Return(nil, assert.AnError)
	args = map[string]interface{}{"ID": "badid", "config_json": string(configBytes)}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = ot.updateOpensearchConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "api error")
}

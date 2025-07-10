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

func TestMongoTool_getMongoDBConfig(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("GetMongoDBConfig", mock.Anything, "cid").Return(&godo.MongoDBConfig{Verbosity: new(int)}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	mt := &MongoTool{client: client}
	args := map[string]interface{}{"ID": "cid"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := mt.getMongoDBConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "verbosity")
	mockDB.AssertExpectations(t)
	// Error case: missing ID
	reqMissing := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{}}}
	res, err = mt.getMongoDBConfig(context.Background(), reqMissing)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Cluster ID is required")
}

func TestMongoTool_updateMongoDBConfig(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("UpdateMongoDBConfig", mock.Anything, "cid", mock.Anything).Return(nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	mt := &MongoTool{client: client}
	cfg := godo.MongoDBConfig{Verbosity: new(int)}
	cfgJSON, _ := json.Marshal(cfg)
	args := map[string]interface{}{"ID": "cid", "config_json": string(cfgJSON)}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := mt.updateMongoDBConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "MongoDB config updated successfully")
	mockDB.AssertExpectations(t)
	// Error case: missing config_json
	args = map[string]interface{}{"ID": "cid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = mt.updateMongoDBConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "config_json is required")
}

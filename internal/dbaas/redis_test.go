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

func TestRedisTool_getRedisConfig(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	val := "volatile-lru"
	mockDB.On("GetRedisConfig", mock.Anything, "cid").Return(&godo.RedisConfig{RedisMaxmemoryPolicy: &val}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	rt := &RedisTool{client: client}
	args := map[string]interface{}{"ID": "cid"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := rt.getRedisConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "volatile-lru")
	mockDB.AssertExpectations(t)
	// Error case: missing ID
	reqMissing := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{}}}
	res, err = rt.getRedisConfig(context.Background(), reqMissing)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Cluster ID is required")
	// API error
	mockDB.On("GetRedisConfig", mock.Anything, "badid").Return(nil, nil, assert.AnError)
	args = map[string]interface{}{"ID": "badid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = rt.getRedisConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "api error")
}

func TestRedisTool_updateRedisConfig(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	val := "allkeys-lru"
	mockDB.On("UpdateRedisConfig", mock.Anything, "cid", mock.Anything).Return(&godo.Response{}, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	rt := &RedisTool{client: client}
	config := godo.RedisConfig{RedisMaxmemoryPolicy: &val}
	configBytes, _ := json.Marshal(config)
	args := map[string]interface{}{"ID": "cid", "config_json": string(configBytes)}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := rt.updateRedisConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "updated successfully")
	mockDB.AssertExpectations(t)
	// Error case: missing ID
	args = map[string]interface{}{"config_json": string(configBytes)}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = rt.updateRedisConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Cluster ID is required")
	// Error case: missing config_json
	args = map[string]interface{}{"ID": "cid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = rt.updateRedisConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "config_json is required")
	// Error case: invalid config_json
	args = map[string]interface{}{"ID": "cid", "config_json": "notjson"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = rt.updateRedisConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Invalid config_json")
	// API error
	mockDB.On("UpdateRedisConfig", mock.Anything, "badid", mock.Anything).Return(nil, assert.AnError)
	args = map[string]interface{}{"ID": "badid", "config_json": string(configBytes)}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = rt.updateRedisConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "api error")
}

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

func TestMysqlTool_getMySQLConfig(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("GetMySQLConfig", mock.Anything, "cid").Return(&godo.MySQLConfig{SQLMode: new(string)}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	mt := &MysqlTool{client: client}
	args := map[string]interface{}{"ID": "cid"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := mt.getMySQLConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "sql_mode")
	mockDB.AssertExpectations(t)
	// Error case: missing ID
	reqMissing := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{}}}
	res, err = mt.getMySQLConfig(context.Background(), reqMissing)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Cluster ID is required")
}

func TestMysqlTool_updateMySQLConfig(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("UpdateMySQLConfig", mock.Anything, "cid", mock.Anything).Return(nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	mt := &MysqlTool{client: client}
	cfg := godo.MySQLConfig{SQLMode: new(string)}
	cfgJSON, _ := json.Marshal(cfg)
	args := map[string]interface{}{"ID": "cid", "config_json": string(cfgJSON)}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := mt.updateMySQLConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "MySQL config updated successfully")
	mockDB.AssertExpectations(t)
	// Error case: missing config_json
	args = map[string]interface{}{"ID": "cid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = mt.updateMySQLConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "config_json is required")
}

func TestMysqlTool_getSQLMode(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("GetSQLMode", mock.Anything, "cid").Return("STRICT_TRANS_TABLES", nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	mt := &MysqlTool{client: client}
	args := map[string]interface{}{"ID": "cid"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := mt.getSQLMode(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "STRICT_TRANS_TABLES")
	mockDB.AssertExpectations(t)
	// Error case: missing ID
	reqMissing := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{}}}
	res, err = mt.getSQLMode(context.Background(), reqMissing)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Cluster ID is required")
}

func TestMysqlTool_setSQLMode(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("SetSQLMode", mock.Anything, "cid", mock.Anything).Return(nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	mt := &MysqlTool{client: client}
	args := map[string]interface{}{"ID": "cid", "modes": "STRICT_TRANS_TABLES"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := mt.setSQLMode(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "SQL mode set successfully")
	mockDB.AssertExpectations(t)
	// Error case: missing modes
	args = map[string]interface{}{"ID": "cid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = mt.setSQLMode(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "SQL modes are required")
}

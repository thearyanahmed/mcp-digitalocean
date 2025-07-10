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

func TestLogSinkTool_createLogsink(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("CreateLogsink", mock.Anything, "cid", mock.Anything).Return(&godo.DatabaseLogsink{ID: "sink1"}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	lt := &LogSinkTool{client: client}
	cfg := godo.DatabaseLogsinkConfig{}
	cfgJSON, _ := json.Marshal(cfg)
	args := map[string]interface{}{"ID": "cid", "sink_name": "sink1", "sink_type": "type1", "config_json": string(cfgJSON)}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := lt.createLogsink(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "sink1")
	mockDB.AssertExpectations(t)
	// Error case: missing sink_name
	args = map[string]interface{}{"ID": "cid", "sink_type": "type1", "config_json": string(cfgJSON)}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = lt.createLogsink(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "sink_name is required")
}

func TestLogSinkTool_getLogsink(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("GetLogsink", mock.Anything, "cid", "sink2").Return(&godo.DatabaseLogsink{ID: "sink2"}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	lt := &LogSinkTool{client: client}
	args := map[string]interface{}{"ID": "cid", "logsink_id": "sink2"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := lt.getLogsink(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "sink2")
	mockDB.AssertExpectations(t)
	// Error case: missing logsink_id
	args = map[string]interface{}{"ID": "cid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = lt.getLogsink(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "logsink_id is required")
}

func TestLogSinkTool_listLogsinks(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("ListLogsinks", mock.Anything, "cid", mock.Anything).Return([]godo.DatabaseLogsink{{ID: "sink3"}}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	lt := &LogSinkTool{client: client}
	args := map[string]interface{}{"ID": "cid"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := lt.listLogsinks(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "sink3")
	mockDB.AssertExpectations(t)
	// Error case: missing ID
	reqMissing := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{}}}
	res, err = lt.listLogsinks(context.Background(), reqMissing)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Cluster ID is required")
}

func TestLogSinkTool_updateLogsink(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("UpdateLogsink", mock.Anything, "cid", "sink4", mock.Anything).Return(nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	lt := &LogSinkTool{client: client}
	cfg := godo.DatabaseLogsinkConfig{}
	cfgJSON, _ := json.Marshal(cfg)
	args := map[string]interface{}{"ID": "cid", "logsink_id": "sink4", "config_json": string(cfgJSON)}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := lt.updateLogsink(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Logsink updated successfully")
	mockDB.AssertExpectations(t)
	// Error case: missing logsink_id
	args = map[string]interface{}{"ID": "cid", "config_json": string(cfgJSON)}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = lt.updateLogsink(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "logsink_id is required")
}

func TestLogSinkTool_deleteLogsink(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("DeleteLogsink", mock.Anything, "cid", "sink5").Return(nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	lt := &LogSinkTool{client: client}
	args := map[string]interface{}{"ID": "cid", "logsink_id": "sink5"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := lt.deleteLogsink(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Logsink deleted successfully")
	mockDB.AssertExpectations(t)
	// Error case: missing logsink_id
	args = map[string]interface{}{"ID": "cid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = lt.deleteLogsink(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "logsink_id is required")
}

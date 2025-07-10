package dbaas

import (
	"context"
	"testing"

	"mcp-digitalocean/internal/dbaas/mocks"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDBTool_listDBs(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("ListDBs", mock.Anything, "cid", mock.Anything).Return([]godo.DatabaseDB{{Name: "db1"}}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	dt := &DBTool{client: client}
	args := map[string]interface{}{"ID": "cid"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := dt.listDBs(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "db1")
	mockDB.AssertExpectations(t)
	// Error case: missing ID
	reqMissing := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{}}}
	res, err = dt.listDBs(context.Background(), reqMissing)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Cluster ID is required")
}

func TestDBTool_createDB(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("CreateDB", mock.Anything, "cid", mock.Anything).Return(&godo.DatabaseDB{Name: "db2"}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	dt := &DBTool{client: client}
	args := map[string]interface{}{"ID": "cid", "name": "db2"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := dt.createDB(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "db2")
	mockDB.AssertExpectations(t)
	// Error case: missing name
	args = map[string]interface{}{"ID": "cid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = dt.createDB(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Database name is required")
}

func TestDBTool_getDB(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("GetDB", mock.Anything, "cid", "db3").Return(&godo.DatabaseDB{Name: "db3"}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	dt := &DBTool{client: client}
	args := map[string]interface{}{"ID": "cid", "name": "db3"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := dt.getDB(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "db3")
	mockDB.AssertExpectations(t)
	// Error case: missing name
	args = map[string]interface{}{"ID": "cid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = dt.getDB(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Database name is required")
}

func TestDBTool_deleteDB(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("DeleteDB", mock.Anything, "cid", "db4").Return(nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	dt := &DBTool{client: client}
	args := map[string]interface{}{"ID": "cid", "name": "db4"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := dt.deleteDB(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Database deleted successfully")
	mockDB.AssertExpectations(t)
	// Error case: missing name
	args = map[string]interface{}{"ID": "cid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = dt.deleteDB(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Database name is required")
}

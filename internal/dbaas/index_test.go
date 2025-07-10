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

func TestIndexTool_listIndexes(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("ListIndexes", mock.Anything, "cid", mock.Anything).Return([]godo.DatabaseIndex{{IndexName: "idx1"}}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	it := &IndexTool{client: client}
	args := map[string]interface{}{"ID": "cid"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := it.listIndexes(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "idx1")
	mockDB.AssertExpectations(t)
	// Error case: missing ID
	reqMissing := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{}}}
	res, err = it.listIndexes(context.Background(), reqMissing)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Cluster ID is required")
}

func TestIndexTool_deleteIndex(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("DeleteIndex", mock.Anything, "cid", "idx1").Return(nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	it := &IndexTool{client: client}
	args := map[string]interface{}{"ID": "cid", "name": "idx1"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := it.deleteIndex(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Index deleted successfully")
	mockDB.AssertExpectations(t)
	// Error case: missing ID
	args = map[string]interface{}{"name": "idx1"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = it.deleteIndex(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Cluster ID is required")
	// Error case: missing name
	args = map[string]interface{}{"ID": "cid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = it.deleteIndex(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Index name is required")
}

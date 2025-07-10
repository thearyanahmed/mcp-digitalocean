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

func TestPoolTool_listPools(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("ListPools", mock.Anything, "cid", (*godo.ListOptions)(nil)).Return([]godo.DatabasePool{{Name: "pool1"}}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	pt := &PoolTool{client: client}
	args := map[string]interface{}{"ID": "cid"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := pt.listPools(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "pool1")
	mockDB.AssertExpectations(t)
	// Error case: missing ID
	reqMissing := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{}}}
	res, err = pt.listPools(context.Background(), reqMissing)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Cluster ID is required")
	// API error
	mockDB.On("ListPools", mock.Anything, "badid", (*godo.ListOptions)(nil)).Return(nil, nil, assert.AnError)
	args = map[string]interface{}{"ID": "badid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = pt.listPools(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "api error")
}

func TestPoolTool_createPool(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("CreatePool", mock.Anything, "cid", mock.Anything).Return(&godo.DatabasePool{Name: "pool1"}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	pt := &PoolTool{client: client}
	args := map[string]interface{}{
		"ID": "cid", "user": "u", "name": "pool1", "database": "db", "mode": "transaction", "size": float64(5),
	}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := pt.createPool(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "pool1")
	mockDB.AssertExpectations(t)
	// Error cases: missing required fields
	errMsgs := map[string]string{
		"ID":       "Cluster ID is required",
		"user":     "User is required",
		"name":     "Pool name is required",
		"database": "Database is required",
		"mode":     "Mode is required",
		"size":     "Size is required and must be a number",
	}
	for _, field := range []string{"ID", "user", "name", "database", "mode", "size"} {
		badArgs := map[string]interface{}{
			"ID": "cid", "user": "u", "name": "pool1", "database": "db", "mode": "transaction", "size": float64(5),
		}
		delete(badArgs, field)
		req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: badArgs}}
		res, err := pt.createPool(context.Background(), req)
		assert.NoError(t, err)
		assert.Contains(t, res.Content[0].(mcp.TextContent).Text, errMsgs[field])
	}
	// API error
	mockDB.On("CreatePool", mock.Anything, "badid", mock.Anything).Return(nil, nil, assert.AnError)
	args = map[string]interface{}{
		"ID": "badid", "user": "u", "name": "pool1", "database": "db", "mode": "transaction", "size": float64(5),
	}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = pt.createPool(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "api error")
}

func TestPoolTool_getPool(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("GetPool", mock.Anything, "cid", "pool1").Return(&godo.DatabasePool{Name: "pool1"}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	pt := &PoolTool{client: client}
	args := map[string]interface{}{"ID": "cid", "name": "pool1"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := pt.getPool(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "pool1")
	mockDB.AssertExpectations(t)
	// Error cases: missing ID or name
	args = map[string]interface{}{"name": "pool1"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = pt.getPool(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Cluster ID is required")
	args = map[string]interface{}{"ID": "cid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = pt.getPool(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Pool name is required")
	// API error
	mockDB.On("GetPool", mock.Anything, "badid", "pool1").Return(nil, nil, assert.AnError)
	args = map[string]interface{}{"ID": "badid", "name": "pool1"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = pt.getPool(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "api error")
}

func TestPoolTool_deletePool(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("DeletePool", mock.Anything, "cid", "pool1").Return(&godo.Response{}, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	pt := &PoolTool{client: client}
	args := map[string]interface{}{"ID": "cid", "name": "pool1"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := pt.deletePool(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "deleted successfully")
	mockDB.AssertExpectations(t)
	// Error cases: missing ID or name
	args = map[string]interface{}{"name": "pool1"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = pt.deletePool(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Cluster ID is required")
	args = map[string]interface{}{"ID": "cid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = pt.deletePool(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Pool name is required")
	// API error
	mockDB.On("DeletePool", mock.Anything, "badid", "pool1").Return(nil, assert.AnError)
	args = map[string]interface{}{"ID": "badid", "name": "pool1"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = pt.deletePool(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "api error")
}

func TestPoolTool_updatePool(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("UpdatePool", mock.Anything, "cid", "pool1", mock.Anything).Return(&godo.Response{}, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	pt := &PoolTool{client: client}
	args := map[string]interface{}{
		"ID": "cid", "name": "pool1", "database": "db", "mode": "transaction", "size": float64(5), "user": "u",
	}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := pt.updatePool(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "updated successfully")
	mockDB.AssertExpectations(t)
	// Error cases: missing required fields
	errMsgs := map[string]string{
		"ID":       "Cluster ID is required",
		"name":     "Pool name is required",
		"database": "Database is required",
		"mode":     "Mode is required",
		"size":     "Size is required and must be a number",
	}
	for _, field := range []string{"ID", "name", "database", "mode", "size"} {
		badArgs := map[string]interface{}{
			"ID": "cid", "name": "pool1", "database": "db", "mode": "transaction", "size": float64(5), "user": "u",
		}
		delete(badArgs, field)
		req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: badArgs}}
		res, err := pt.updatePool(context.Background(), req)
		assert.NoError(t, err)
		assert.Contains(t, res.Content[0].(mcp.TextContent).Text, errMsgs[field])
	}
	// API error
	mockDB.On("UpdatePool", mock.Anything, "badid", "pool1", mock.Anything).Return(nil, assert.AnError)
	args = map[string]interface{}{
		"ID": "badid", "name": "pool1", "database": "db", "mode": "transaction", "size": float64(5), "user": "u",
	}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = pt.updatePool(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "api error")
}

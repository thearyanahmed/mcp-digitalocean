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

func TestReplicaTool_getReplica(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("GetReplica", mock.Anything, "cid", "rep1").Return(&godo.DatabaseReplica{Name: "rep1"}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	rt := &ReplicaTool{client: client}
	args := map[string]interface{}{"ID": "cid", "name": "rep1"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := rt.getReplica(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "rep1")
	mockDB.AssertExpectations(t)
	// Error cases: missing ID or name
	args = map[string]interface{}{"name": "rep1"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = rt.getReplica(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Cluster ID is required")
	args = map[string]interface{}{"ID": "cid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = rt.getReplica(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Replica name is required")
	// API error
	mockDB.On("GetReplica", mock.Anything, "badid", "rep1").Return(nil, nil, assert.AnError)
	args = map[string]interface{}{"ID": "badid", "name": "rep1"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = rt.getReplica(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "api error")
}

func TestReplicaTool_listReplicas(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("ListReplicas", mock.Anything, "cid", (*godo.ListOptions)(nil)).Return([]godo.DatabaseReplica{{Name: "rep1"}}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	rt := &ReplicaTool{client: client}
	args := map[string]interface{}{"ID": "cid"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := rt.listReplicas(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "rep1")
	mockDB.AssertExpectations(t)
	// Error case: missing ID
	reqMissing := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{}}}
	res, err = rt.listReplicas(context.Background(), reqMissing)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Cluster ID is required")
	// API error
	mockDB.On("ListReplicas", mock.Anything, "badid", (*godo.ListOptions)(nil)).Return(nil, nil, assert.AnError)
	args = map[string]interface{}{"ID": "badid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = rt.listReplicas(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "api error")
}

func TestReplicaTool_createReplica(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("CreateReplica", mock.Anything, "cid", mock.Anything).Return(&godo.DatabaseReplica{Name: "rep1"}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	rt := &ReplicaTool{client: client}
	args := map[string]interface{}{
		"ID": "cid", "name": "rep1", "region": "nyc1", "size": "db-s-1vcpu-1gb",
	}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := rt.createReplica(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "rep1")
	mockDB.AssertExpectations(t)
	// Error cases: missing required fields
	errMsgs := map[string]string{
		"ID":     "Cluster ID is required",
		"name":   "Replica name is required",
		"region": "Replica region is required",
		"size":   "Replica size is required",
	}
	for _, field := range []string{"ID", "name", "region", "size"} {
		badArgs := map[string]interface{}{
			"ID": "cid", "name": "rep1", "region": "nyc1", "size": "db-s-1vcpu-1gb",
		}
		delete(badArgs, field)
		req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: badArgs}}
		res, err := rt.createReplica(context.Background(), req)
		assert.NoError(t, err)
		assert.Contains(t, res.Content[0].(mcp.TextContent).Text, errMsgs[field])
	}
	// API error
	mockDB.On("CreateReplica", mock.Anything, "badid", mock.Anything).Return(nil, nil, assert.AnError)
	args = map[string]interface{}{
		"ID": "badid", "name": "rep1", "region": "nyc1", "size": "db-s-1vcpu-1gb",
	}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = rt.createReplica(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "api error")
}

func TestReplicaTool_deleteReplica(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("DeleteReplica", mock.Anything, "cid", "rep1").Return(&godo.Response{}, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	rt := &ReplicaTool{client: client}
	args := map[string]interface{}{"ID": "cid", "name": "rep1"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := rt.deleteReplica(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "deleted successfully")
	mockDB.AssertExpectations(t)
	// Error cases: missing ID or name
	args = map[string]interface{}{"name": "rep1"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = rt.deleteReplica(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Cluster ID is required")
	args = map[string]interface{}{"ID": "cid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = rt.deleteReplica(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Replica name is required")
	// API error
	mockDB.On("DeleteReplica", mock.Anything, "badid", "rep1").Return(nil, assert.AnError)
	args = map[string]interface{}{"ID": "badid", "name": "rep1"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = rt.deleteReplica(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "api error")
}

func TestReplicaTool_promoteReplicaToPrimary(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("PromoteReplicaToPrimary", mock.Anything, "cid", "rep1").Return(&godo.Response{}, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	rt := &ReplicaTool{client: client}
	args := map[string]interface{}{"ID": "cid", "name": "rep1"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := rt.promoteReplicaToPrimary(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "promoted to primary successfully")
	mockDB.AssertExpectations(t)
	// Error cases: missing ID or name
	args = map[string]interface{}{"name": "rep1"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = rt.promoteReplicaToPrimary(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Cluster ID is required")
	args = map[string]interface{}{"ID": "cid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = rt.promoteReplicaToPrimary(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Replica name is required")
	// API error
	mockDB.On("PromoteReplicaToPrimary", mock.Anything, "badid", "rep1").Return(nil, assert.AnError)
	args = map[string]interface{}{"ID": "badid", "name": "rep1"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = rt.promoteReplicaToPrimary(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "api error")
}

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

func TestMigrationTool_startOnlineMigration(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("StartOnlineMigration", mock.Anything, "cid", mock.Anything).Return(&godo.DatabaseOnlineMigrationStatus{ID: "mig1"}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	mt := &MigrationTool{client: client}
	source := godo.DatabaseOnlineMigrationConfig{Host: "host"}
	sourceJSON, _ := json.Marshal(source)
	args := map[string]interface{}{"ID": "cid", "source_json": string(sourceJSON)}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := mt.startOnlineMigration(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "mig1")
	mockDB.AssertExpectations(t)
	// Error case: missing source_json
	args = map[string]interface{}{"ID": "cid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = mt.startOnlineMigration(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "source_json is required")
}

func TestMigrationTool_stopOnlineMigration(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("StopOnlineMigration", mock.Anything, "cid", "mig2").Return(nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	mt := &MigrationTool{client: client}
	args := map[string]interface{}{"ID": "cid", "migration_id": "mig2"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := mt.stopOnlineMigration(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Online migration stopped successfully")
	mockDB.AssertExpectations(t)
	// Error case: missing migration_id
	args = map[string]interface{}{"ID": "cid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = mt.stopOnlineMigration(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "migration_id is required")
}

func TestMigrationTool_getOnlineMigrationStatus(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("GetOnlineMigrationStatus", mock.Anything, "cid").Return(&godo.DatabaseOnlineMigrationStatus{ID: "mig3"}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	mt := &MigrationTool{client: client}
	args := map[string]interface{}{"ID": "cid"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := mt.getOnlineMigrationStatus(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "mig3")
	mockDB.AssertExpectations(t)
	// Error case: missing ID
	reqMissing := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{}}}
	res, err = mt.getOnlineMigrationStatus(context.Background(), reqMissing)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Cluster ID is required")
}

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

func TestFirewallTool_getFirewallRules(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("GetFirewallRules", mock.Anything, "cid").Return([]godo.DatabaseFirewallRule{{UUID: "rule1"}}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	ft := &FirewallTool{client: client}
	args := map[string]interface{}{"ID": "cid"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := ft.getFirewallRules(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "rule1")
	mockDB.AssertExpectations(t)
	// Error case: missing ID
	reqMissing := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{}}}
	res, err = ft.getFirewallRules(context.Background(), reqMissing)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Cluster ID is required")
}

func TestFirewallTool_updateFirewallRules(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("UpdateFirewallRules", mock.Anything, "cid", mock.Anything).Return(nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	ft := &FirewallTool{client: client}
	rules := []*godo.DatabaseFirewallRule{{UUID: "rule2"}}
	rulesJSON, _ := json.Marshal(rules)
	args := map[string]interface{}{"ID": "cid", "rules_json": string(rulesJSON)}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := ft.updateFirewallRules(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Firewall rules updated successfully")
	mockDB.AssertExpectations(t)
	// Error case: missing rules_json
	args = map[string]interface{}{"ID": "cid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = ft.updateFirewallRules(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "rules_json is required")
}

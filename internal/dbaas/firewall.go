package dbaas

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *ClusterTool) getFirewallRules(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	rules, _, err := s.client.Databases.GetFirewallRules(ctx, id)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonRules, err := json.MarshalIndent(rules, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonRules)), nil
}

func (s *ClusterTool) updateFirewallRules(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	rulesStr, ok := args["rules_json"].(string)
	if !ok || rulesStr == "" {
		return mcp.NewToolResultError("rules_json is required (JSON array of firewall rules)"), nil
	}
	var rules []*godo.DatabaseFirewallRule
	err := json.Unmarshal([]byte(rulesStr), &rules)
	if err != nil {
		return mcp.NewToolResultError("Invalid rules_json: " + err.Error()), nil
	}
	updateReq := &godo.DatabaseUpdateFirewallRulesRequest{Rules: rules}
	_, err = s.client.Databases.UpdateFirewallRules(ctx, id, updateReq)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Firewall rules updated successfully"), nil
}

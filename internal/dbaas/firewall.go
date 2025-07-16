package dbaas

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type FirewallTool struct {
	client *godo.Client
}

func NewFirewallTool(client *godo.Client) *FirewallTool {
	return &FirewallTool{
		client: client,
	}
}

func (s *FirewallTool) getFirewallRules(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

func (s *FirewallTool) updateFirewallRules(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
func (s *FirewallTool) Tools() []server.ServerTool {
	return []server.ServerTool{
		{
			Handler: s.getFirewallRules,
			Tool: mcp.NewTool("digitalocean-dbaascluster-get-firewall-rules",
				mcp.WithDescription("Get the firewall rules for a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
			),
		},
		{
			Handler: s.updateFirewallRules,
			Tool: mcp.NewTool("digitalocean-dbaascluster-update-firewall-rules",
				mcp.WithDescription("Update the firewall rules for a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("rules_json", mcp.Required(), mcp.Description("JSON array of firewall rules to set")),
			),
		},
	}
}

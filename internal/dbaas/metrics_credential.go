package dbaas

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *ClusterTool) getMetricsCredentials(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	creds, _, err := s.client.Databases.GetMetricsCredentials(ctx)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonCreds, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonCreds)), nil
}

func (s *ClusterTool) updateMetricsCredentials(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	credsStr, ok := args["credentials_json"].(string)
	if !ok || credsStr == "" {
		return mcp.NewToolResultError("credentials_json is required (JSON for DatabaseMetricsCredentials)"), nil
	}
	var creds godo.DatabaseMetricsCredentials
	err := json.Unmarshal([]byte(credsStr), &creds)
	if err != nil {
		return mcp.NewToolResultError("Invalid credentials_json: " + err.Error()), nil
	}
	updateReq := &godo.DatabaseUpdateMetricsCredentialsRequest{Credentials: &creds}
	_, err = s.client.Databases.UpdateMetricsCredentials(ctx, updateReq)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Metrics credentials updated successfully"), nil
}

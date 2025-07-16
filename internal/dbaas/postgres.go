package dbaas

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type PostgreSQLTool struct {
	client *godo.Client
}

func NewPostgreSQLTool(client *godo.Client) *PostgreSQLTool {
	return &PostgreSQLTool{
		client: client,
	}
}

func (s *PostgreSQLTool) getPostgreSQLConfig(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster id is required"), nil
	}
	cfg, _, err := s.client.Databases.GetPostgreSQLConfig(ctx, id)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonCfg, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonCfg)), nil
}

func (s *PostgreSQLTool) updatePostgreSQLConfig(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster id is required"), nil
	}
	configStr, ok := args["config_json"].(string)
	if !ok || configStr == "" {
		return mcp.NewToolResultError("config_json is required (JSON for PostgreSQLConfig)"), nil
	}
	var config godo.PostgreSQLConfig
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		return mcp.NewToolResultError("Invalid config_json: " + err.Error()), nil
	}
	_, err = s.client.Databases.UpdatePostgreSQLConfig(ctx, id, &config)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("PostgreSQL config updated successfully"), nil
}

func (s *PostgreSQLTool) Tools() []server.ServerTool {
	return []server.ServerTool{
		{
			Handler: s.getPostgreSQLConfig,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-get-postgresql-config",
				mcp.WithDescription("Get the PostgreSQL config for a cluster by its id"),
				mcp.WithString("id", mcp.Required(), mcp.Description("The cluster UUID")),
			),
		},
		{
			Handler: s.updatePostgreSQLConfig,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-update-postgresql-config",
				mcp.WithDescription("Update the PostgreSQL config for a cluster by its id. Accepts a JSON string for the config."),
				mcp.WithString("id", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("config_json", mcp.Required(), mcp.Description("JSON for the PostgreSQLConfig to set")),
			),
		},
	}
}

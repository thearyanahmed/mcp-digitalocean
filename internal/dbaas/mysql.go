package dbaas

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MysqlTool struct {
	client *godo.Client
}

func NewMysqlTool(client *godo.Client) *MysqlTool {
	return &MysqlTool{
		client: client,
	}
}

func (s *MysqlTool) getMySQLConfig(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	cfg, _, err := s.client.Databases.GetMySQLConfig(ctx, id)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonCfg, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonCfg)), nil
}

func (s *MysqlTool) updateMySQLConfig(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	configStr, ok := args["config_json"].(string)
	if !ok || configStr == "" {
		return mcp.NewToolResultError("config_json is required (JSON for MySQLConfig)"), nil
	}
	var config godo.MySQLConfig
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		return mcp.NewToolResultError("Invalid config_json: " + err.Error()), nil
	}
	_, err = s.client.Databases.UpdateMySQLConfig(ctx, id, &config)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("MySQL config updated successfully"), nil
}

func (s *MysqlTool) getSQLMode(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	mode, _, err := s.client.Databases.GetSQLMode(ctx, id)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText(mode), nil
}

func (s *MysqlTool) setSQLMode(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	modesStr, ok := args["modes"].(string)
	if !ok || modesStr == "" {
		return mcp.NewToolResultError("SQL modes are required (comma-separated)"), nil
	}
	modes := []string{}
	for _, m := range strings.Split(modesStr, ",") {
		m = strings.TrimSpace(m)
		if m != "" {
			modes = append(modes, m)
		}
	}
	_, err := s.client.Databases.SetSQLMode(ctx, id, modes...)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("SQL mode set successfully"), nil
}

func (s *MysqlTool) Tools() []server.ServerTool {
	return []server.ServerTool{
		{
			Handler: s.getMySQLConfig,
			Tool: mcp.NewTool("digitalocean-dbaascluster-get-mysql-config",
				mcp.WithDescription("Get the MySQL config for a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
			),
		},
		{
			Handler: s.updateMySQLConfig,
			Tool: mcp.NewTool("digitalocean-dbaascluster-update-mysql-config",
				mcp.WithDescription("Update the MySQL config for a cluster by its ID. Accepts a JSON string for the config."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("config_json", mcp.Required(), mcp.Description("JSON for the MySQLConfig to set")),
			),
		},
		{
			Handler: s.getSQLMode,
			Tool: mcp.NewTool("digitalocean-dbaascluster-get-sql-mode",
				mcp.WithDescription("Get the SQL mode for a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
			),
		},
		{
			Handler: s.setSQLMode,
			Tool: mcp.NewTool("digitalocean-dbaascluster-set-sql-mode",
				mcp.WithDescription("Set the SQL mode for a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("modes", mcp.Required(), mcp.Description("Comma-separated SQL modes to set")),
			),
		},
	}
}

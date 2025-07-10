package dbaas

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MigrationTool struct {
	client *godo.Client
}

func NewMigrationTool(client *godo.Client) *MigrationTool {
	return &MigrationTool{
		client: client,
	}
}

func (s *MigrationTool) startOnlineMigration(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	sourceStr, ok := args["source_json"].(string)
	if !ok || sourceStr == "" {
		return mcp.NewToolResultError("source_json is required (JSON for DatabaseOnlineMigrationConfig)"), nil
	}
	var source godo.DatabaseOnlineMigrationConfig
	err := json.Unmarshal([]byte(sourceStr), &source)
	if err != nil {
		return mcp.NewToolResultError("Invalid source_json: " + err.Error()), nil
	}
	disableSSL := false
	if dssl, ok := args["disable_ssl"].(string); ok && dssl != "" {
		if b, err := strconv.ParseBool(dssl); err == nil {
			disableSSL = b
		}
	}
	var ignoreDBs []string
	if ignoreStr, ok := args["ignore_dbs"].(string); ok && ignoreStr != "" {
		for _, db := range strings.Split(ignoreStr, ",") {
			db = strings.TrimSpace(db)
			if db != "" {
				ignoreDBs = append(ignoreDBs, db)
			}
		}
	}
	startReq := &godo.DatabaseStartOnlineMigrationRequest{
		Source:     &source,
		DisableSSL: disableSSL,
		IgnoreDBs:  ignoreDBs,
	}
	status, _, err := s.client.Databases.StartOnlineMigration(ctx, id, startReq)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonStatus, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonStatus)), nil
}

func (s *MigrationTool) stopOnlineMigration(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	migrationID, ok := args["migration_id"].(string)
	if !ok || migrationID == "" {
		return mcp.NewToolResultError("migration_id is required"), nil
	}
	_, err := s.client.Databases.StopOnlineMigration(ctx, id, migrationID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Online migration stopped successfully"), nil
}

func (s *MigrationTool) getOnlineMigrationStatus(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	status, _, err := s.client.Databases.GetOnlineMigrationStatus(ctx, id)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonStatus, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonStatus)), nil
}

func (s *MigrationTool) Tools() []server.ServerTool {
	return []server.ServerTool{
		{
			Handler: s.startOnlineMigration,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-start-online-migration",
				mcp.WithDescription("Start an online migration for a database cluster by its ID. Accepts source_json (DatabaseOnlineMigrationConfig as JSON, required), disable_ssl (optional, bool as string), and ignore_dbs (optional, comma-separated)."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("source_json", mcp.Required(), mcp.Description("DatabaseOnlineMigrationConfig as JSON (required)")),
				mcp.WithString("disable_ssl", mcp.Description("Disable SSL for migration (optional, bool as string)")),
				mcp.WithString("ignore_dbs", mcp.Description("Comma-separated list of DBs to ignore (optional)")),
			),
		},
		{
			Handler: s.stopOnlineMigration,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-stop-online-migration",
				mcp.WithDescription("Stop an online migration for a database cluster by its ID and migration_id."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("migration_id", mcp.Required(), mcp.Description("The migration ID to stop")),
			),
		},
		{
			Handler: s.getOnlineMigrationStatus,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-get-online-migration-status",
				mcp.WithDescription("Get the online migration status for a database cluster by its ID."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
			),
		},
	}
}

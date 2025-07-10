package dbaas

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *ClusterTool) startOnlineMigration(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

func (s *ClusterTool) stopOnlineMigration(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

func (s *ClusterTool) getOnlineMigrationStatus(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

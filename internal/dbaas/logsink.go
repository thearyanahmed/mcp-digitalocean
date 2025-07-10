package dbaas

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type LogSinkTool struct {
	client *godo.Client
}

func NewLogSinkTool(client *godo.Client) *LogSinkTool {
	return &LogSinkTool{
		client: client,
	}
}

func (s *LogSinkTool) createLogsink(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	name, ok := args["sink_name"].(string)
	if !ok || name == "" {
		return mcp.NewToolResultError("sink_name is required"), nil
	}
	typeStr, ok := args["sink_type"].(string)
	if !ok || typeStr == "" {
		return mcp.NewToolResultError("sink_type is required"), nil
	}
	configStr, ok := args["config_json"].(string)
	if !ok || configStr == "" {
		return mcp.NewToolResultError("config_json is required (JSON for DatabaseLogsinkConfig)"), nil
	}
	var config godo.DatabaseLogsinkConfig
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		return mcp.NewToolResultError("Invalid config_json: " + err.Error()), nil
	}
	createReq := &godo.DatabaseCreateLogsinkRequest{
		Name:   name,
		Type:   typeStr,
		Config: &config,
	}
	logsink, _, err := s.client.Databases.CreateLogsink(ctx, id, createReq)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonLogsink, err := json.MarshalIndent(logsink, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonLogsink)), nil
}

func (s *LogSinkTool) getLogsink(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	logsinkID, ok := args["logsink_id"].(string)
	if !ok || logsinkID == "" {
		return mcp.NewToolResultError("logsink_id is required"), nil
	}
	logsink, _, err := s.client.Databases.GetLogsink(ctx, id, logsinkID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonLogsink, err := json.MarshalIndent(logsink, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonLogsink)), nil
}

func (s *LogSinkTool) listLogsinks(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	// Optional pagination
	page := 0
	if pStr, ok := args["page"].(string); ok && pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil {
			page = p
		}
	}
	perPage := 0
	if ppStr, ok := args["per_page"].(string); ok && ppStr != "" {
		if pp, err := strconv.Atoi(ppStr); err == nil {
			perPage = pp
		}
	}
	var opts *godo.ListOptions
	if page > 0 || perPage > 0 {
		opts = &godo.ListOptions{Page: page, PerPage: perPage}
	}
	logsinks, _, err := s.client.Databases.ListLogsinks(ctx, id, opts)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonLogsinks, err := json.MarshalIndent(logsinks, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonLogsinks)), nil
}

func (s *LogSinkTool) updateLogsink(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	logsinkID, ok := args["logsink_id"].(string)
	if !ok || logsinkID == "" {
		return mcp.NewToolResultError("logsink_id is required"), nil
	}
	configStr, ok := args["config_json"].(string)
	if !ok || configStr == "" {
		return mcp.NewToolResultError("config_json is required (JSON for DatabaseLogsinkConfig)"), nil
	}
	var config godo.DatabaseLogsinkConfig
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		return mcp.NewToolResultError("Invalid config_json: " + err.Error()), nil
	}
	updateReq := &godo.DatabaseUpdateLogsinkRequest{Config: &config}
	_, err = s.client.Databases.UpdateLogsink(ctx, id, logsinkID, updateReq)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Logsink updated successfully"), nil
}

func (s *LogSinkTool) deleteLogsink(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	logsinkID, ok := args["logsink_id"].(string)
	if !ok || logsinkID == "" {
		return mcp.NewToolResultError("logsink_id is required"), nil
	}
	_, err := s.client.Databases.DeleteLogsink(ctx, id, logsinkID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Logsink deleted successfully"), nil
}

func (s *LogSinkTool) Tools() []server.ServerTool {
	return []server.ServerTool{
		{
			Handler: s.createLogsink,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-create-logsink",
				mcp.WithDescription("Create a logsink for a database cluster by its ID. Accepts sink_name, sink_type, and config_json (DatabaseLogsinkConfig as JSON, all required)."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("sink_name", mcp.Required(), mcp.Description("The logsink name to create")),
				mcp.WithString("sink_type", mcp.Required(), mcp.Description("The logsink type (e.g., opensearch, datadog, logtail, papertrail)")),
				mcp.WithString("config_json", mcp.Required(), mcp.Description("DatabaseLogsinkConfig as JSON (required)")),
			),
		},
		{
			Handler: s.getLogsink,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-get-logsink",
				mcp.WithDescription("Get a logsink for a database cluster by its ID and logsink_id."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("logsink_id", mcp.Required(), mcp.Description("The logsink ID to get")),
			),
		},
		{
			Handler: s.listLogsinks,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-list-logsinks",
				mcp.WithDescription("List logsinks for a database cluster by its ID. Supports pagination: page, per_page (optional, integer as string)."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("page", mcp.Description("Page number for pagination (optional, integer as string)")),
				mcp.WithString("per_page", mcp.Description("Number of results per page (optional, integer as string)")),
			),
		},
		{
			Handler: s.updateLogsink,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-update-logsink",
				mcp.WithDescription("Update a logsink for a database cluster by its ID and logsink_id. Accepts config_json (DatabaseLogsinkConfig as JSON, required)."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("logsink_id", mcp.Required(), mcp.Description("The logsink ID to update")),
				mcp.WithString("config_json", mcp.Required(), mcp.Description("DatabaseLogsinkConfig as JSON (required)")),
			),
		},
		{
			Handler: s.deleteLogsink,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-delete-logsink",
				mcp.WithDescription("Delete a logsink for a database cluster by its ID and logsink_id."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("logsink_id", mcp.Required(), mcp.Description("The logsink ID to delete")),
			),
		},
	}
}

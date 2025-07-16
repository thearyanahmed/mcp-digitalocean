package dbaas

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type OpenSearchTool struct {
	client *godo.Client
}

func NewOpenSearchTool(client *godo.Client) *OpenSearchTool {
	return &OpenSearchTool{
		client: client,
	}
}

func (s *OpenSearchTool) getOpensearchConfig(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster id is required"), nil
	}
	cfg, _, err := s.client.Databases.GetOpensearchConfig(ctx, id)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonCfg, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonCfg)), nil
}

func (s *OpenSearchTool) updateOpensearchConfig(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster id is required"), nil
	}
	configStr, ok := args["config_json"].(string)
	if !ok || configStr == "" {
		return mcp.NewToolResultError("config_json is required (JSON for OpensearchConfig)"), nil
	}
	var config godo.OpensearchConfig
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		return mcp.NewToolResultError("Invalid config_json: " + err.Error()), nil
	}
	_, err = s.client.Databases.UpdateOpensearchConfig(ctx, id, &config)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Opensearch config updated successfully"), nil
}

func (s *OpenSearchTool) Tools() []server.ServerTool {
	return []server.ServerTool{
		{
			Handler: s.getOpensearchConfig,
			Tool: mcp.NewTool("digitalocean-dbaascluster-get-opensearch-config",
				mcp.WithDescription("Get the Opensearch config for a cluster by its id"),
				mcp.WithString("id", mcp.Required(), mcp.Description("The cluster UUID")),
			),
		},
		{
			Handler: s.updateOpensearchConfig,
			Tool: mcp.NewTool("digitalocean-dbaascluster-update-opensearch-config",
				mcp.WithDescription("Update the Opensearch config for a cluster by its id. Accepts a JSON string for the config."),
				mcp.WithString("id", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("config_json", mcp.Required(), mcp.Description("JSON for the OpensearchConfig to set")),
			),
		},
	}
}

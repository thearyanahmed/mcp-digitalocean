package dbaas

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type RedisTool struct {
	client *godo.Client
}

func NewRedisTool(client *godo.Client) *RedisTool {
	return &RedisTool{
		client: client,
	}
}

func (s *RedisTool) getRedisConfig(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster id is required"), nil
	}
	cfg, _, err := s.client.Databases.GetRedisConfig(ctx, id)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonCfg, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonCfg)), nil
}

func (s *RedisTool) updateRedisConfig(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster id is required"), nil
	}
	configStr, ok := args["config_json"].(string)
	if !ok || configStr == "" {
		return mcp.NewToolResultError("config_json is required (JSON for RedisConfig)"), nil
	}
	var config godo.RedisConfig
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		return mcp.NewToolResultError("Invalid config_json: " + err.Error()), nil
	}
	_, err = s.client.Databases.UpdateRedisConfig(ctx, id, &config)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Redis config updated successfully"), nil
}

func (s *RedisTool) Tools() []server.ServerTool {
	return []server.ServerTool{
		{
			Handler: s.getRedisConfig,
			Tool: mcp.NewTool("digitalocean-databases-cluster-get-redis-config",
				mcp.WithDescription("Get the Redis config for a cluster by its id"),
				mcp.WithString("id", mcp.Required(), mcp.Description("The cluster UUID")),
			),
		},
		{
			Handler: s.updateRedisConfig,
			Tool: mcp.NewTool("digitalocean-databases-cluster-update-redis-config",
				mcp.WithDescription("Update the Redis config for a cluster by its id. Accepts a JSON string for the config."),
				mcp.WithString("id", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("config_json", mcp.Required(), mcp.Description("JSON for the RedisConfig to set")),
			),
		},
	}
}

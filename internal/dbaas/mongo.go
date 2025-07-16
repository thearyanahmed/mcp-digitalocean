package dbaas

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MongoTool struct {
	client *godo.Client
}

func NewMongoTool(client *godo.Client) *MongoTool {
	return &MongoTool{
		client: client,
	}
}

func (s *MongoTool) getMongoDBConfig(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster id is required"), nil
	}
	cfg, _, err := s.client.Databases.GetMongoDBConfig(ctx, id)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonCfg, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonCfg)), nil
}

func (s *MongoTool) updateMongoDBConfig(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster id is required"), nil
	}
	configStr, ok := args["config_json"].(string)
	if !ok || configStr == "" {
		return mcp.NewToolResultError("config_json is required (JSON for MongoDBConfig)"), nil
	}
	var config godo.MongoDBConfig
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		return mcp.NewToolResultError("Invalid config_json: " + err.Error()), nil
	}
	_, err = s.client.Databases.UpdateMongoDBConfig(ctx, id, &config)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("MongoDB config updated successfully"), nil
}

func (s *MongoTool) Tools() []server.ServerTool {
	return []server.ServerTool{
		{
			Handler: s.getMongoDBConfig,
			Tool: mcp.NewTool("digitalocean-databases-cluster-get-mongodb-config",
				mcp.WithDescription("Get the MongoDB config for a cluster by its id"),
				mcp.WithString("id", mcp.Required(), mcp.Description("The cluster UUID")),
			),
		},
		{
			Handler: s.updateMongoDBConfig,
			Tool: mcp.NewTool("digitalocean-databases-cluster-update-mongodb-config",
				mcp.WithDescription("Update the MongoDB config for a cluster by its id. Accepts a JSON string for the config."),
				mcp.WithString("id", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("config_json", mcp.Required(), mcp.Description("JSON for the MongoDBConfig to set")),
			),
		},
	}
}

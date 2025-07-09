package dbaas

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ClusterTool struct {
	client *godo.Client
}

func NewClusterTool(client *godo.Client) *ClusterTool {
	return &ClusterTool{
		client: client,
	}
}

func (s *ClusterTool) listCluster(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clusters, _, err := s.client.Databases.List(ctx, nil)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	jsonKey, err := json.MarshalIndent(clusters, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	return mcp.NewToolResultText(string(jsonKey)), nil
}

func (s *ClusterTool) Tools() []server.ServerTool {
	return []server.ServerTool{

		{
			Handler: s.listCluster,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-list",
				mcp.WithDescription("Get list of  Cluster"),
			),
		},
	}
}

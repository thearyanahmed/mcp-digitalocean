package marketplace

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type OneClickTool struct {
	client *godo.Client
}

func NewOneClickTool(client *godo.Client) *OneClickTool {
	return &OneClickTool{
		client: client,
	}
}

func (o *OneClickTool) listOneClickApps(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()

	// Type parameter is optional, defaults to "droplet"
	oneClickType := "droplet"
	if typeArg, ok := args["type"]; ok {
		if typeStr, ok := typeArg.(string); ok && typeStr != "" {
			oneClickType = typeStr
		}
	}

	apps, _, err := o.client.OneClick.List(ctx, oneClickType)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list 1-click apps: %v", err)), nil
	}

	result, err := json.Marshal(map[string]interface{}{
		"apps": apps,
		"type": oneClickType,
	})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(result)), nil
}

func (o *OneClickTool) installKubernetesApps(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()

	// Get cluster UUID
	clusterUUIDArg, ok := args["cluster_uuid"]
	if !ok {
		return mcp.NewToolResultError("cluster_uuid parameter is required"), nil
	}

	clusterUUID, ok := clusterUUIDArg.(string)
	if !ok {
		return mcp.NewToolResultError("cluster_uuid must be a string"), nil
	}

	if clusterUUID == "" {
		return mcp.NewToolResultError("cluster_uuid cannot be empty"), nil
	}

	// Get app slugs
	slugsArg, ok := args["app_slugs"]
	if !ok {
		return mcp.NewToolResultError("app_slugs parameter is required"), nil
	}

	slugsInterface, ok := slugsArg.([]interface{})
	if !ok {
		return mcp.NewToolResultError("app_slugs must be an array"), nil
	}

	slugs := make([]string, len(slugsInterface))
	for i, slug := range slugsInterface {
		slugStr, ok := slug.(string)
		if !ok {
			return mcp.NewToolResultError("all app_slugs must be strings"), nil
		}
		slugs[i] = slugStr
	}

	if len(slugs) == 0 {
		return mcp.NewToolResultError("app_slugs cannot be empty"), nil
	}

	installRequest := &godo.InstallKubernetesAppsRequest{
		Slugs:       slugs,
		ClusterUUID: clusterUUID,
	}

	response, _, err := o.client.OneClick.InstallKubernetes(ctx, installRequest)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to install Kubernetes apps: %v", err)), nil
	}

	result, err := json.Marshal(response)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(result)), nil
}

func (o *OneClickTool) Tools() []server.ServerTool {
	return []server.ServerTool{
		{
			Handler: o.listOneClickApps,
			Tool: mcp.NewTool("digitalocean-1-click-list",
				mcp.WithDescription("List available 1-click applications from the DigitalOcean marketplace"),
				mcp.WithString("type", mcp.Description("Type of 1-click apps to list (e.g., 'droplet', 'kubernetes'). Defaults to 'droplet'")),
			),
		},
		{
			Handler: o.installKubernetesApps,
			Tool: mcp.NewTool("digitalocean-1-click-install-kubernetes",
				mcp.WithDescription("Install 1-click applications on a Kubernetes cluster"),
				mcp.WithString("cluster_uuid", mcp.Required(), mcp.Description("UUID of the Kubernetes cluster to install apps on")),
				mcp.WithArray("app_slugs", mcp.Required(), mcp.Description("Array of app slugs to install")),
			),
		},
	}
}

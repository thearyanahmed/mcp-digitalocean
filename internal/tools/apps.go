package tools

import (
	"context"
	"encoding/json"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AppTool provides app management tools
type AppTool struct {
	client *godo.Client
}

// NewAppTool creates a new app tool
func NewAppTool(client *godo.Client) *AppTool {
	return &AppTool{
		client: client,
	}
}

// CreateApp creates a new app
func (a *AppTool) CreateApp(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Read individual parameters
	name := req.Params.Arguments["Name"].(string)
	region := req.Params.Arguments["Region"].(string)

	// Build the AppSpec manually
	spec := godo.AppSpec{
		Name:       name,
		Region:     region,
	}

	// Create the app
	app, _, err := a.client.Apps.Create(ctx, &godo.AppCreateRequest{Spec: &spec})
	if err != nil {
		return nil, err
	}

	// Serialize to JSON
	jsonApp, err := json.MarshalIndent(app, "", "  ")
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(jsonApp)), nil
}

// DeleteApp deletes an app
func (a *AppTool) DeleteApp(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Read the app ID
	appID := req.Params.Arguments["ID"].(string)

	// Delete the app
	_, err := a.client.Apps.Delete(ctx, appID)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText("App deleted successfully"), nil
}

// Tools returns a list of tool functions
func (a *AppTool) Tools() []server.ServerTool {
	return []server.ServerTool{
		{
			Handler: a.CreateApp,
			Tool: mcp.NewTool("digitalocean-app-create",
				mcp.WithDescription("Create a new app"),
				mcp.WithString("Name", mcp.Required(), mcp.Description("Name of the app")),
				mcp.WithString("Region", mcp.Required(), mcp.Description("Region of the app")),
				mcp.WithString("Tier", mcp.Required(), mcp.Description("Tier of the app")),
			),
		},
		{
			Handler: a.DeleteApp,
			Tool: mcp.NewTool("digitalocean-app-delete",
				mcp.WithDescription("Delete an app"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("ID of the app to delete")),
			),
		},
	}
}
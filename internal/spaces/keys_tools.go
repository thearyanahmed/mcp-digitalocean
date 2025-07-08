package spaces

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// SpacesKeysTool provides Spaces keys management tools
type SpacesKeysTool struct {
	client *godo.Client
}

// NewSpacesKeysTool creates a new Spaces keys tool
func NewSpacesKeysTool(client *godo.Client) *SpacesKeysTool {
	return &SpacesKeysTool{
		client: client,
	}
}

// createSpacesKey creates a new Spaces key
func (s *SpacesKeysTool) createSpacesKey(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	name := args["Name"].(string)

	createRequest := &godo.SpacesKeyCreateRequest{
		Name: name,
		Grants: []*godo.Grant{
			{
				Bucket:     "",
				Permission: godo.SpacesKeyFullAccess,
			},
		},
	}

	key, _, err := s.client.SpacesKeys.Create(ctx, createRequest)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonKey, err := json.MarshalIndent(key, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	return mcp.NewToolResultText(string(jsonKey)), nil
}

// updateSpacesKey updates an existing Spaces key
func (s *SpacesKeysTool) updateSpacesKey(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	keyID := args["ID"].(string)
	name := args["Name"].(string)

	updateRequest := &godo.SpacesKeyUpdateRequest{
		Name: name,
	}

	key, _, err := s.client.SpacesKeys.Update(ctx, keyID, updateRequest)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonKey, err := json.MarshalIndent(key, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	return mcp.NewToolResultText(string(jsonKey)), nil
}

// deleteSpacesKey deletes a Spaces key
func (s *SpacesKeysTool) deleteSpacesKey(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	keyID := req.GetArguments()["ID"].(string)

	_, err := s.client.SpacesKeys.Delete(ctx, keyID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	return mcp.NewToolResultText("Spaces key deleted successfully"), nil
}

// Tools returns a list of tool functions
func (s *SpacesKeysTool) Tools() []server.ServerTool {
	return []server.ServerTool{
		{
			Handler: s.createSpacesKey,
			Tool: mcp.NewTool("digitalocean-spaces-key-create",
				mcp.WithDescription("Create a new Spaces key"),
				mcp.WithString("Name", mcp.Required(), mcp.Description("Name for the Spaces key")),
			),
		},
		{
			Handler: s.updateSpacesKey,
			Tool: mcp.NewTool("digitalocean-spaces-key-update",
				mcp.WithDescription("Update an existing Spaces key"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("ID of the Spaces key to update")),
				mcp.WithString("Name", mcp.Required(), mcp.Description("New name for the Spaces key")),
			),
		},
		{
			Handler: s.deleteSpacesKey,
			Tool: mcp.NewTool("digitalocean-spaces-key-delete",
				mcp.WithDescription("Delete a Spaces key"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("ID of the Spaces key to delete")),
			),
		},
	}
}

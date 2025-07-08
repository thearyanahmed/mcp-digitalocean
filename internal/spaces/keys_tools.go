package spaces

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type KeysTool struct {
	client *godo.Client
}

func NewSpacesKeysTool(client *godo.Client) *KeysTool {
	return &KeysTool{
		client: client,
	}
}

func (s *KeysTool) createSpacesKey(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	nameArg, ok := args["Name"]
	if !ok {
		return mcp.NewToolResultError("Name parameter is required"), nil
	}

	name, ok := nameArg.(string)
	if !ok {
		return mcp.NewToolResultError("Name must be a string"), nil
	}

	if name == "" {
		return mcp.NewToolResultError("Name cannot be empty"), nil
	}

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

func (s *KeysTool) updateSpacesKey(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()

	accessKeyArg, ok := args["AccessKey"]
	if !ok {
		return mcp.NewToolResultError("AccessKey parameter is required"), nil
	}

	accessKey, ok := accessKeyArg.(string)
	if !ok {
		return mcp.NewToolResultError("AccessKey must be a string"), nil
	}

	if accessKey == "" {
		return mcp.NewToolResultError("AccessKey cannot be empty"), nil
	}

	nameArg, ok := args["Name"]
	if !ok {
		return mcp.NewToolResultError("Name parameter is required"), nil
	}

	name, ok := nameArg.(string)
	if !ok {
		return mcp.NewToolResultError("Name must be a string"), nil
	}

	if name == "" {
		return mcp.NewToolResultError("Name cannot be empty"), nil
	}

	updateRequest := &godo.SpacesKeyUpdateRequest{
		Name: name,
	}

	key, _, err := s.client.SpacesKeys.Update(ctx, accessKey, updateRequest)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonKey, err := json.MarshalIndent(key, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	return mcp.NewToolResultText(string(jsonKey)), nil
}

func (s *KeysTool) deleteSpacesKey(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()

	accessKeyArg, ok := args["AccessKey"]
	if !ok {
		return mcp.NewToolResultError("AccessKey parameter is required"), nil
	}

	accessKey, ok := accessKeyArg.(string)
	if !ok {
		return mcp.NewToolResultError("AccessKey must be a string"), nil
	}

	if accessKey == "" {
		return mcp.NewToolResultError("AccessKey cannot be empty"), nil
	}

	_, err := s.client.SpacesKeys.Delete(ctx, accessKey)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Spaces key deleted successfully"), nil
}

func (s *KeysTool) listSpacesKeys(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	keys, _, err := s.client.SpacesKeys.List(ctx, &godo.ListOptions{})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	jsonKeys, err := json.MarshalIndent(keys, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	return mcp.NewToolResultText(string(jsonKeys)), nil
}

func (s *KeysTool) getSpacesKey(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()

	accessKeyArg, ok := args["AccessKey"]
	if !ok {
		return mcp.NewToolResultError("AccessKey parameter is required"), nil
	}

	accessKey, ok := accessKeyArg.(string)
	if !ok {
		return mcp.NewToolResultError("AccessKey must be a string"), nil
	}

	if accessKey == "" {
		return mcp.NewToolResultError("AccessKey cannot be empty"), nil
	}

	key, _, err := s.client.SpacesKeys.Get(ctx, accessKey)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	jsonKey, err := json.MarshalIndent(key, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	return mcp.NewToolResultText(string(jsonKey)), nil
}

func (s *KeysTool) Tools() []server.ServerTool {
	return []server.ServerTool{
		{
			Handler: s.listSpacesKeys,
			Tool: mcp.NewTool("digitalocean-spaces-key-list",
				mcp.WithDescription("List all Spaces keys"),
			),
		},
		{
			Handler: s.getSpacesKey,
			Tool: mcp.NewTool("digitalocean-spaces-key-get",
				mcp.WithDescription("Get a specific Spaces key"),
				mcp.WithString("AccessKey", mcp.Required(), mcp.Description("Access Key of the Spaces key to retrieve")),
			),
		},
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
				mcp.WithString("AccessKey", mcp.Required(), mcp.Description("Access Key of the Spaces key to update")),
				mcp.WithString("Name", mcp.Required(), mcp.Description("New name for the Spaces key")),
			),
		},
		{
			Handler: s.deleteSpacesKey,
			Tool: mcp.NewTool("digitalocean-spaces-key-delete",
				mcp.WithDescription("Delete a Spaces key"),
				mcp.WithString("AccessKey", mcp.Required(), mcp.Description("Access Key of the Spaces key to delete")),
			),
		},
	}
}

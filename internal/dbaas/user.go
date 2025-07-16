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

type UserTool struct {
	client *godo.Client
}

func NewUserTool(client *godo.Client) *UserTool {
	return &UserTool{
		client: client,
	}
}

func (s *UserTool) getUser(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster id is required"), nil
	}
	user, ok := args["user"].(string)
	if !ok || user == "" {
		return mcp.NewToolResultError("User name is required"), nil
	}

	dbUser, _, err := s.client.Databases.GetUser(ctx, id, user)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonUser, err := json.MarshalIndent(dbUser, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonUser)), nil
}

func (s *UserTool) listUsers(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster id is required"), nil
	}

	// Optional pagination
	page := 0
	if pStr, ok := args["page"].(string); ok && pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil {
			page = p
		}
	}
	perPage := 0
	if pp, ok := args["per_page"].(int); ok {
		perPage = pp
	}

	var opts *godo.ListOptions
	if page > 0 || perPage > 0 {
		opts = &godo.ListOptions{Page: page, PerPage: perPage}
	}

	users, _, err := s.client.Databases.ListUsers(ctx, id, opts)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonUsers, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonUsers)), nil
}

func (s *UserTool) createUser(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster id is required"), nil
	}
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return mcp.NewToolResultError("User name is required"), nil
	}

	createReq := &godo.DatabaseCreateUserRequest{Name: name}

	if plugin, ok := args["mysql_auth_plugin"].(string); ok && plugin != "" {
		createReq.MySQLSettings = &godo.DatabaseMySQLUserSettings{AuthPlugin: plugin}
	}

	if settingsStr, ok := args["settings_json"].(string); ok && settingsStr != "" {
		var settings godo.DatabaseUserSettings
		err := json.Unmarshal([]byte(settingsStr), &settings)
		if err != nil {
			return mcp.NewToolResultError("Invalid settings_json: " + err.Error()), nil
		}
		createReq.Settings = &settings
	}

	dbUser, _, err := s.client.Databases.CreateUser(ctx, id, createReq)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonUser, err := json.MarshalIndent(dbUser, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonUser)), nil
}

func (s *UserTool) updateUser(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster id is required"), nil
	}
	user, ok := args["user"].(string)
	if !ok || user == "" {
		return mcp.NewToolResultError("User name is required"), nil
	}

	updateReq := &godo.DatabaseUpdateUserRequest{}

	if settingsStr, ok := args["settings_json"].(string); ok && settingsStr != "" {
		var settings godo.DatabaseUserSettings
		err := json.Unmarshal([]byte(settingsStr), &settings)
		if err != nil {
			return mcp.NewToolResultError("Invalid settings_json: " + err.Error()), nil
		}
		updateReq.Settings = &settings
	}

	dbUser, _, err := s.client.Databases.UpdateUser(ctx, id, user, updateReq)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonUser, err := json.MarshalIndent(dbUser, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonUser)), nil
}

func (s *UserTool) deleteUser(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster id is required"), nil
	}
	user, ok := args["user"].(string)
	if !ok || user == "" {
		return mcp.NewToolResultError("User name is required"), nil
	}

	_, err := s.client.Databases.DeleteUser(ctx, id, user)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("User deleted successfully"), nil
}
func (s *UserTool) Tools() []server.ServerTool {
	return []server.ServerTool{
		{
			Handler: s.getUser,
			Tool: mcp.NewTool("digitalocean-dbaascluster-get-user",
				mcp.WithDescription("Get a database user by cluster id and user name"),
				mcp.WithString("id", mcp.Required(), mcp.Description("The cluster id (UUID)")),
				mcp.WithString("user", mcp.Required(), mcp.Description("The user name")),
			),
		},
		{
			Handler: s.listUsers,
			Tool: mcp.NewTool("digitalocean-dbaascluster-list-users",
				mcp.WithDescription("List database users for a cluster by its id"),
				mcp.WithString("id", mcp.Required(), mcp.Description("The cluster id (UUID)")),
				mcp.WithString("page", mcp.Description("Page number for pagination (optional)")),
				mcp.WithNumber("per_page", mcp.Description("Number of results per page (optional)")),
			),
		},
		{
			Handler: s.createUser,
			Tool: mcp.NewTool("digitalocean-dbaascluster-create-user",
				mcp.WithDescription("Create a database user for a cluster by its id"),
				mcp.WithString("id", mcp.Required(), mcp.Description("The cluster id (UUID)")),
				mcp.WithString("name", mcp.Required(), mcp.Description("The user name")),
				mcp.WithString("mysql_auth_plugin", mcp.Description("MySQL auth plugin (optional, e.g., mysql_native_password)")),
				mcp.WithString("settings_json", mcp.Description("Raw JSON for DatabaseUserSettings (optional)")),
			),
		},
		{
			Handler: s.updateUser,
			Tool: mcp.NewTool("digitalocean-dbaascluster-update-user",
				mcp.WithDescription("Update a database user for a cluster by its id and user name"),
				mcp.WithString("id", mcp.Required(), mcp.Description("The cluster id (UUID)")),
				mcp.WithString("user", mcp.Required(), mcp.Description("The user name")),
				mcp.WithString("settings_json", mcp.Description("Raw JSON for DatabaseUserSettings (optional)")),
			),
		},
		{
			Handler: s.deleteUser,
			Tool: mcp.NewTool("digitalocean-dbaascluster-delete-user",
				mcp.WithDescription("Delete a database user by cluster id and user name"),
				mcp.WithString("id", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("user", mcp.Required(), mcp.Description("The user name to delete")),
			),
		}}
}

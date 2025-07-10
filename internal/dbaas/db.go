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

type DBTool struct {
	client *godo.Client
}

func NewDBTool(client *godo.Client) *DBTool {
	return &DBTool{
		client: client,
	}
}

func (s *DBTool) listDBs(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	dbs, _, err := s.client.Databases.ListDBs(ctx, id, opts)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonDBs, err := json.MarshalIndent(dbs, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonDBs)), nil
}

func (s *DBTool) createDB(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return mcp.NewToolResultError("Database name is required"), nil
	}

	createReq := &godo.DatabaseCreateDBRequest{Name: name}
	db, _, err := s.client.Databases.CreateDB(ctx, id, createReq)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonDB, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonDB)), nil
}

func (s *DBTool) getDB(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return mcp.NewToolResultError("Database name is required"), nil
	}
	db, _, err := s.client.Databases.GetDB(ctx, id, name)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonDB, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonDB)), nil
}

func (s *DBTool) deleteDB(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return mcp.NewToolResultError("Database name is required"), nil
	}
	_, err := s.client.Databases.DeleteDB(ctx, id, name)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Database deleted successfully"), nil
}

func (s *DBTool) Tools() []server.ServerTool {
	return []server.ServerTool{
		{
			Handler: s.listDBs,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-list-dbs",
				mcp.WithDescription("List databases for a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("page", mcp.Description("Page number for pagination")),
				mcp.WithString("per_page", mcp.Description("Number of results per page")),
			),
		},
		{
			Handler: s.createDB,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-create-db",
				mcp.WithDescription("Create a database for a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("name", mcp.Required(), mcp.Description("The database name to create")),
			),
		},
		{
			Handler: s.getDB,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-get-db",
				mcp.WithDescription("Get a database for a cluster by its ID and database name"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("name", mcp.Required(), mcp.Description("The database name to get")),
			),
		},
		{
			Handler: s.deleteDB,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-delete-db",
				mcp.WithDescription("Delete a database for a cluster by its ID and database name"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("name", mcp.Required(), mcp.Description("The database name to delete")),
			),
		},
	}
}

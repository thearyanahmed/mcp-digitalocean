package dbaas

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *ClusterTool) listPools(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	pools, _, err := s.client.Databases.ListPools(ctx, id, opts)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonPools, err := json.MarshalIndent(pools, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonPools)), nil
}

func (s *ClusterTool) createPool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	user, ok := args["user"].(string)
	if !ok || user == "" {
		return mcp.NewToolResultError("User is required"), nil
	}
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return mcp.NewToolResultError("Pool name is required"), nil
	}
	database, ok := args["database"].(string)
	if !ok || database == "" {
		return mcp.NewToolResultError("Database is required"), nil
	}
	mode, ok := args["mode"].(string)
	if !ok || mode == "" {
		return mcp.NewToolResultError("Mode is required"), nil
	}
	sizeF, ok := args["size"].(float64)
	if !ok {
		return mcp.NewToolResultError("Size is required and must be a number"), nil
	}
	size := int(sizeF)

	createReq := &godo.DatabaseCreatePoolRequest{
		User:     user,
		Name:     name,
		Database: database,
		Mode:     mode,
		Size:     size,
	}
	pool, _, err := s.client.Databases.CreatePool(ctx, id, createReq)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonPool, err := json.MarshalIndent(pool, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonPool)), nil
}

func (s *ClusterTool) getPool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return mcp.NewToolResultError("Pool name is required"), nil
	}
	pool, _, err := s.client.Databases.GetPool(ctx, id, name)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonPool, err := json.MarshalIndent(pool, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonPool)), nil
}

func (s *ClusterTool) deletePool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return mcp.NewToolResultError("Pool name is required"), nil
	}
	_, err := s.client.Databases.DeletePool(ctx, id, name)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Pool deleted successfully"), nil
}

func (s *ClusterTool) updatePool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return mcp.NewToolResultError("Pool name is required"), nil
	}
	database, ok := args["database"].(string)
	if !ok || database == "" {
		return mcp.NewToolResultError("Database is required"), nil
	}
	mode, ok := args["mode"].(string)
	if !ok || mode == "" {
		return mcp.NewToolResultError("Mode is required"), nil
	}
	sizeF, ok := args["size"].(float64)
	if !ok {
		return mcp.NewToolResultError("Size is required and must be a number"), nil
	}
	size := int(sizeF)
	user, _ := args["user"].(string)

	updateReq := &godo.DatabaseUpdatePoolRequest{
		User:     user,
		Database: database,
		Mode:     mode,
		Size:     size,
	}
	_, err := s.client.Databases.UpdatePool(ctx, id, name, updateReq)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Pool updated successfully"), nil
}

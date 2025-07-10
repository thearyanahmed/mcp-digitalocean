package dbaas

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type IndexTool struct {
	client *godo.Client
}

func NewIndexTool(client *godo.Client) *IndexTool {
	return &IndexTool{
		client: client,
	}
}

func (s *IndexTool) listIndexes(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}

	opts := &godo.ListOptions{}
	if pStr, ok := args["page"].(string); ok && pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil {
			opts.Page = p
		}
	}
	if ppStr, ok := args["per_page"].(string); ok && ppStr != "" {
		if pp, err := strconv.Atoi(ppStr); err == nil {
			opts.PerPage = pp
		}
	}
	if wpStr, ok := args["with_projects"].(string); ok && wpStr != "" {
		if wp, err := strconv.ParseBool(wpStr); err == nil {
			opts.WithProjects = wp
		}
	}
	if odStr, ok := args["only_deployed"].(string); ok && odStr != "" {
		if od, err := strconv.ParseBool(odStr); err == nil {
			opts.Deployed = od
		}
	}
	if poStr, ok := args["public_only"].(string); ok && poStr != "" {
		if po, err := strconv.ParseBool(poStr); err == nil {
			opts.PublicOnly = po
		}
	}
	if ucStr, ok := args["usecases"].(string); ok && ucStr != "" {
		ucList := []string{}
		for _, u := range strings.Split(ucStr, ",") {
			u = strings.TrimSpace(u)
			if u != "" {
				ucList = append(ucList, u)
			}
		}
		opts.Usecases = ucList
	}

	indexes, _, err := s.client.Databases.ListIndexes(ctx, id, opts)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonIndexes, err := json.MarshalIndent(indexes, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonIndexes)), nil
}

func (s *IndexTool) deleteIndex(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return mcp.NewToolResultError("Index name is required"), nil
	}
	_, err := s.client.Databases.DeleteIndex(ctx, id, name)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Index deleted successfully"), nil
}

func (s *IndexTool) Tools() []server.ServerTool {
	return []server.ServerTool{
		{
			Handler: s.listIndexes,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-list-indexes",
				mcp.WithDescription("List indexes for a cluster by its ID. Supports all ListOptions: page, per_page, with_projects, only_deployed, public_only, usecases (comma-separated)."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("page", mcp.Description("Page number for pagination (optional, integer as string)")),
				mcp.WithString("per_page", mcp.Description("Number of results per page (optional, integer as string)")),
				mcp.WithString("with_projects", mcp.Description("Whether to include project_id fields (optional, bool as string)")),
				mcp.WithString("only_deployed", mcp.Description("Only list deployed agents (optional, bool as string)")),
				mcp.WithString("public_only", mcp.Description("Include only public models (optional, bool as string)")),
				mcp.WithString("usecases", mcp.Description("Comma-separated usecases to filter (optional)")),
			),
		},
		{
			Handler: s.deleteIndex,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-delete-index",
				mcp.WithDescription("Delete an index for a cluster by its ID and index name."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("name", mcp.Required(), mcp.Description("The index name to delete")),
			),
		},
	}
}

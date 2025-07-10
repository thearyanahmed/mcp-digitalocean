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

type ReplicaTool struct {
	client *godo.Client
}

func NewReplicaTool(client *godo.Client) *ReplicaTool {
	return &ReplicaTool{
		client: client,
	}
}

func (s *ReplicaTool) getReplica(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return mcp.NewToolResultError("Replica name is required"), nil
	}
	replica, _, err := s.client.Databases.GetReplica(ctx, id, name)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonReplica, err := json.MarshalIndent(replica, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonReplica)), nil
}

func (s *ReplicaTool) listReplicas(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	replicas, _, err := s.client.Databases.ListReplicas(ctx, id, opts)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonReplicas, err := json.MarshalIndent(replicas, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonReplicas)), nil
}

func (s *ReplicaTool) createReplica(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return mcp.NewToolResultError("Replica name is required"), nil
	}
	region, ok := args["region"].(string)
	if !ok || region == "" {
		return mcp.NewToolResultError("Replica region is required"), nil
	}
	size, ok := args["size"].(string)
	if !ok || size == "" {
		return mcp.NewToolResultError("Replica size is required"), nil
	}
	privateNetworkUUID, _ := args["private_network_uuid"].(string)
	tags := []string{}
	if tagsRaw, ok := args["tags"].(string); ok && tagsRaw != "" {
		for _, t := range strings.Split(tagsRaw, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				tags = append(tags, t)
			}
		}
	}
	storageSizeMib := uint64(0)
	if ssm, ok := args["storage_size_mib"].(float64); ok {
		storageSizeMib = uint64(ssm)
	}

	createReq := &godo.DatabaseCreateReplicaRequest{
		Name:               name,
		Region:             region,
		Size:               size,
		PrivateNetworkUUID: privateNetworkUUID,
		Tags:               tags,
		StorageSizeMib:     storageSizeMib,
	}

	replica, _, err := s.client.Databases.CreateReplica(ctx, id, createReq)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonReplica, err := json.MarshalIndent(replica, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonReplica)), nil
}

func (s *ReplicaTool) deleteReplica(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return mcp.NewToolResultError("Replica name is required"), nil
	}
	_, err := s.client.Databases.DeleteReplica(ctx, id, name)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Replica deleted successfully"), nil
}

func (s *ReplicaTool) promoteReplicaToPrimary(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return mcp.NewToolResultError("Replica name is required"), nil
	}
	_, err := s.client.Databases.PromoteReplicaToPrimary(ctx, id, name)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Replica promoted to primary successfully"), nil
}
func (s *ReplicaTool) Tools() []server.ServerTool {
	return []server.ServerTool{
		{
			Handler: s.getReplica,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-get-replica",
				mcp.WithDescription("Get a replica for a cluster by its ID and replica name"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("name", mcp.Required(), mcp.Description("The replica name to get")),
			),
		},
		{
			Handler: s.listReplicas,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-list-replicas",
				mcp.WithDescription("List replicas for a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("page", mcp.Description("Page number for pagination")),
				mcp.WithString("per_page", mcp.Description("Number of results per page")),
			),
		},
		{
			Handler: s.createReplica,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-create-replica",
				mcp.WithDescription("Create a replica for a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("name", mcp.Required(), mcp.Description("The replica name to create")),
				mcp.WithString("region", mcp.Required(), mcp.Description("The region for the replica")),
				mcp.WithString("size", mcp.Required(), mcp.Description("The size slug for the replica")),
				mcp.WithString("private_network_uuid", mcp.Description("The private network UUID (optional)")),
				mcp.WithString("tags", mcp.Description("Comma-separated tags to apply to the replica (optional)")),
				mcp.WithNumber("storage_size_mib", mcp.Description("The storage size in MiB (optional)")),
			),
		},
		{
			Handler: s.deleteReplica,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-delete-replica",
				mcp.WithDescription("Delete a replica for a cluster by its ID and replica name"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("name", mcp.Required(), mcp.Description("The replica name to delete")),
			),
		},
		{
			Handler: s.promoteReplicaToPrimary,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-promote-replica",
				mcp.WithDescription("Promote a replica to primary for a cluster by its ID and replica name"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("name", mcp.Required(), mcp.Description("The replica name to promote")),
			),
		},
	}
}

package dbaas

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *ClusterTool) getReplica(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

func (s *ClusterTool) listReplicas(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

func (s *ClusterTool) createReplica(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

func (s *ClusterTool) deleteReplica(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

func (s *ClusterTool) promoteReplicaToPrimary(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

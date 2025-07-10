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

func (s *ClusterTool) getKafkaConfig(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	cfg, _, err := s.client.Databases.GetKafkaConfig(ctx, id)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonCfg, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonCfg)), nil
}

func (s *ClusterTool) updateKafkaConfig(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	configStr, ok := args["config_json"].(string)
	if !ok || configStr == "" {
		return mcp.NewToolResultError("config_json is required (JSON for KafkaConfig)"), nil
	}
	var config godo.KafkaConfig
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		return mcp.NewToolResultError("Invalid config_json: " + err.Error()), nil
	}
	_, err = s.client.Databases.UpdateKafkaConfig(ctx, id, &config)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Kafka config updated successfully"), nil
}

func (s *ClusterTool) listTopics(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	topics, _, err := s.client.Databases.ListTopics(ctx, id, opts)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonTopics, err := json.MarshalIndent(topics, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonTopics)), nil
}

func (s *ClusterTool) createTopic(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return mcp.NewToolResultError("Topic name is required"), nil
	}

	var partitionCount *uint32
	if pcStr, ok := args["partition_count"].(string); ok && pcStr != "" {
		if pc, err := strconv.ParseUint(pcStr, 10, 32); err == nil {
			pc32 := uint32(pc)
			partitionCount = &pc32
		}
	}
	var replicationFactor *uint32
	if rfStr, ok := args["replication_factor"].(string); ok && rfStr != "" {
		if rf, err := strconv.ParseUint(rfStr, 10, 32); err == nil {
			rf32 := uint32(rf)
			replicationFactor = &rf32
		}
	}

	var topicConfig *godo.TopicConfig
	if cfgStr, ok := args["config_json"].(string); ok && cfgStr != "" {
		var cfg godo.TopicConfig
		err := json.Unmarshal([]byte(cfgStr), &cfg)
		if err != nil {
			return mcp.NewToolResultError("Invalid config_json: " + err.Error()), nil
		}
		topicConfig = &cfg
	}

	createReq := &godo.DatabaseCreateTopicRequest{
		Name:              name,
		PartitionCount:    partitionCount,
		ReplicationFactor: replicationFactor,
		Config:            topicConfig,
	}
	topic, _, err := s.client.Databases.CreateTopic(ctx, id, createReq)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonTopic, err := json.MarshalIndent(topic, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonTopic)), nil
}

func (s *ClusterTool) getTopic(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return mcp.NewToolResultError("Topic name is required"), nil
	}
	topic, _, err := s.client.Databases.GetTopic(ctx, id, name)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonTopic, err := json.MarshalIndent(topic, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonTopic)), nil
}

func (s *ClusterTool) deleteTopic(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return mcp.NewToolResultError("Topic name is required"), nil
	}
	_, err := s.client.Databases.DeleteTopic(ctx, id, name)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Topic deleted successfully"), nil
}

func (s *ClusterTool) updateTopic(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return mcp.NewToolResultError("Topic name is required"), nil
	}

	var partitionCount *uint32
	if pcStr, ok := args["partition_count"].(string); ok && pcStr != "" {
		if pc, err := strconv.ParseUint(pcStr, 10, 32); err == nil {
			pc32 := uint32(pc)
			partitionCount = &pc32
		}
	}
	var replicationFactor *uint32
	if rfStr, ok := args["replication_factor"].(string); ok && rfStr != "" {
		if rf, err := strconv.ParseUint(rfStr, 10, 32); err == nil {
			rf32 := uint32(rf)
			replicationFactor = &rf32
		}
	}

	var topicConfig *godo.TopicConfig
	if cfgStr, ok := args["config_json"].(string); ok && cfgStr != "" {
		var cfg godo.TopicConfig
		err := json.Unmarshal([]byte(cfgStr), &cfg)
		if err != nil {
			return mcp.NewToolResultError("Invalid config_json: " + err.Error()), nil
		}
		topicConfig = &cfg
	}

	updateReq := &godo.DatabaseUpdateTopicRequest{
		PartitionCount:    partitionCount,
		ReplicationFactor: replicationFactor,
		Config:            topicConfig,
	}
	_, err := s.client.Databases.UpdateTopic(ctx, id, name, updateReq)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Topic updated successfully"), nil
}

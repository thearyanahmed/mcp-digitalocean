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

type KafkaTool struct {
	client *godo.Client
}

func NewKafkaTool(client *godo.Client) *KafkaTool {
	return &KafkaTool{
		client: client,
	}
}

func (s *KafkaTool) getKafkaConfig(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster id is required"), nil
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

func (s *KafkaTool) updateKafkaConfig(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster id is required"), nil
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

func (s *KafkaTool) listTopics(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster id is required"), nil
	}

	opts := &godo.ListOptions{}
	if pStr, ok := args["page"].(string); ok && pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil {
			opts.Page = p
		}
	}

	if pp, ok := args["per_page"].(int); ok {
		opts.PerPage = pp
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

func (s *KafkaTool) createTopic(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster id is required"), nil
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

func (s *KafkaTool) getTopic(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster id is required"), nil
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

func (s *KafkaTool) deleteTopic(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster id is required"), nil
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

func (s *KafkaTool) updateTopic(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster id is required"), nil
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

func (s *KafkaTool) Tools() []server.ServerTool {
	return []server.ServerTool{
		{
			Handler: s.listTopics,
			Tool: mcp.NewTool("digitalocean-databases-cluster-list-topics",
				mcp.WithDescription("List topics for a database cluster by its id (Kafka clusters). Supports all ListOptions: page, per_page, with_projects, only_deployed, public_only, usecases (comma-separated)."),
				mcp.WithString("id", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("page", mcp.Description("Page number for pagination (optional, integer as string)")),
				mcp.WithNumber("per_page", mcp.Description("Number of results per page (optional, integer)")),
				mcp.WithString("with_projects", mcp.Description("Whether to include project_id fields (optional, bool as string)")),
				mcp.WithString("only_deployed", mcp.Description("Only list deployed agents (optional, bool as string)")),
				mcp.WithString("public_only", mcp.Description("Include only public models (optional, bool as string)")),
				mcp.WithString("usecases", mcp.Description("Comma-separated usecases to filter (optional)")),
			),
		},
		{
			Handler: s.createTopic,
			Tool: mcp.NewTool("digitalocean-databases-cluster-create-topic",
				mcp.WithDescription("Create a topic for a Kafka database cluster by its id. Accepts name (required), partition_count, replication_factor, and config_json (TopicConfig as JSON, all optional)."),
				mcp.WithString("id", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("name", mcp.Required(), mcp.Description("The topic name to create")),
				mcp.WithString("partition_count", mcp.Description("Number of partitions (optional, integer as string)")),
				mcp.WithString("replication_factor", mcp.Description("Replication factor (optional, integer as string)")),
				mcp.WithString("config_json", mcp.Description("TopicConfig as JSON (optional)")),
			),
		},
		{
			Handler: s.getTopic,
			Tool: mcp.NewTool("digitalocean-databases-cluster-get-topic",
				mcp.WithDescription("Get a topic for a Kafka database cluster by its id and topic name."),
				mcp.WithString("id", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("name", mcp.Required(), mcp.Description("The topic name to get")),
			),
		},
		{
			Handler: s.deleteTopic,
			Tool: mcp.NewTool("digitalocean-databases-cluster-delete-topic",
				mcp.WithDescription("Delete a topic for a Kafka database cluster by its id and topic name."),
				mcp.WithString("id", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("name", mcp.Required(), mcp.Description("The topic name to delete")),
			),
		},
		{
			Handler: s.updateTopic,
			Tool: mcp.NewTool("digitalocean-databases-cluster-update-topic",
				mcp.WithDescription("Update a topic for a Kafka database cluster by its id and topic name. Accepts partition_count, replication_factor, and config_json (TopicConfig as JSON, all optional)."),
				mcp.WithString("id", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("name", mcp.Required(), mcp.Description("The topic name to update")),
				mcp.WithString("partition_count", mcp.Description("Number of partitions (optional, integer as string)")),
				mcp.WithString("replication_factor", mcp.Description("Replication factor (optional, integer as string)")),
				mcp.WithString("config_json", mcp.Description("TopicConfig as JSON (optional)")),
			),
		},
		{
			Handler: s.getKafkaConfig,
			Tool: mcp.NewTool("digitalocean-databases-cluster-get-kafka-config",
				mcp.WithDescription("Get the Kafka config for a cluster by its id"),
				mcp.WithString("id", mcp.Required(), mcp.Description("The cluster UUID")),
			),
		},
		{
			Handler: s.updateKafkaConfig,
			Tool: mcp.NewTool("digitalocean-databases-cluster-update-kafka-config",
				mcp.WithDescription("Update the Kafka config for a cluster by its id. Accepts a JSON string for the config."),
				mcp.WithString("id", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("config_json", mcp.Required(), mcp.Description("JSON for the KafkaConfig to set")),
			),
		},
	}
}

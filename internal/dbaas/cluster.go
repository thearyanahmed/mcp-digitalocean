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

type ClusterTool struct {
	client *godo.Client
}

func NewClusterTool(client *godo.Client) *ClusterTool {
	return &ClusterTool{
		client: client,
	}
}

func (s *ClusterTool) listCluster(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()

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

	clusters, _, err := s.client.Databases.List(ctx, opts)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonClusters, err := json.MarshalIndent(clusters, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonClusters)), nil
}

func (s *ClusterTool) getCluster(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, ok := req.GetArguments()["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	cluster, _, err := s.client.Databases.Get(ctx, id)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonCluster, err := json.MarshalIndent(cluster, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonCluster)), nil
}

func (s *ClusterTool) createCluster(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()

	name, _ := args["name"].(string)
	engine, _ := args["engine"].(string)
	version, _ := args["version"].(string)
	region, _ := args["region"].(string)
	size, _ := args["size"].(string)
	numNodes, _ := args["num_nodes"].(float64) // JSON numbers are float64

	tags := []string{}
	if tagsRaw, ok := args["tags"].(string); ok && tagsRaw != "" {
		for _, t := range strings.Split(tagsRaw, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				tags = append(tags, t)
			}
		}
	}

	createReq := &godo.DatabaseCreateRequest{
		Name:       name,
		EngineSlug: engine,
		Version:    version,
		Region:     region,
		SizeSlug:   size,
		NumNodes:   int(numNodes),
		Tags:       tags,
	}

	cluster, _, err := s.client.Databases.Create(ctx, createReq)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonCluster, err := json.MarshalIndent(cluster, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonCluster)), nil
}

func (s *ClusterTool) deleteCluster(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, ok := req.GetArguments()["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	_, err := s.client.Databases.Delete(ctx, id)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Cluster deleted successfully"), nil
}

func (s *ClusterTool) resizeCluster(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}

	size, _ := args["size"].(string)
	numNodes := 0
	if n, ok := args["num_nodes"].(float64); ok {
		numNodes = int(n)
	}
	storageSizeMib := uint64(0)
	if ssm, ok := args["storage_size_mib"].(float64); ok {
		storageSizeMib = uint64(ssm)
	}

	resizeReq := &godo.DatabaseResizeRequest{}
	if size != "" {
		resizeReq.SizeSlug = size
	}
	if numNodes > 0 {
		resizeReq.NumNodes = numNodes
	}
	if storageSizeMib > 0 {
		resizeReq.StorageSizeMib = storageSizeMib
	}

	_, err := s.client.Databases.Resize(ctx, id, resizeReq)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Cluster resize initiated successfully"), nil
}

func (s *ClusterTool) getCA(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, ok := req.GetArguments()["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	ca, _, err := s.client.Databases.GetCA(ctx, id)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonCA, err := json.MarshalIndent(ca, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonCA)), nil
}

func (s *ClusterTool) listBackups(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	backups, _, err := s.client.Databases.ListBackups(ctx, id, opts)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonBackups, err := json.MarshalIndent(backups, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonBackups)), nil
}

func (s *ClusterTool) listOptions(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	options, _, err := s.client.Databases.ListOptions(ctx)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonOptions, err := json.MarshalIndent(options, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonOptions)), nil
}

func (s *ClusterTool) upgradeMajorVersion(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	version, ok := args["version"].(string)
	if !ok || version == "" {
		return mcp.NewToolResultError("Target version is required"), nil
	}
	upgradeReq := &godo.UpgradeVersionRequest{Version: version}
	_, err := s.client.Databases.UpgradeMajorVersion(ctx, id, upgradeReq)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Major version upgrade initiated successfully"), nil
}

func (s *ClusterTool) startOnlineMigration(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	sourceStr, ok := args["source_json"].(string)
	if !ok || sourceStr == "" {
		return mcp.NewToolResultError("source_json is required (JSON for DatabaseOnlineMigrationConfig)"), nil
	}
	var source godo.DatabaseOnlineMigrationConfig
	err := json.Unmarshal([]byte(sourceStr), &source)
	if err != nil {
		return mcp.NewToolResultError("Invalid source_json: " + err.Error()), nil
	}
	disableSSL := false
	if dssl, ok := args["disable_ssl"].(string); ok && dssl != "" {
		if b, err := strconv.ParseBool(dssl); err == nil {
			disableSSL = b
		}
	}
	var ignoreDBs []string
	if ignoreStr, ok := args["ignore_dbs"].(string); ok && ignoreStr != "" {
		for _, db := range strings.Split(ignoreStr, ",") {
			db = strings.TrimSpace(db)
			if db != "" {
				ignoreDBs = append(ignoreDBs, db)
			}
		}
	}
	startReq := &godo.DatabaseStartOnlineMigrationRequest{
		Source:     &source,
		DisableSSL: disableSSL,
		IgnoreDBs:  ignoreDBs,
	}
	status, _, err := s.client.Databases.StartOnlineMigration(ctx, id, startReq)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonStatus, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonStatus)), nil
}

func (s *ClusterTool) stopOnlineMigration(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	migrationID, ok := args["migration_id"].(string)
	if !ok || migrationID == "" {
		return mcp.NewToolResultError("migration_id is required"), nil
	}
	_, err := s.client.Databases.StopOnlineMigration(ctx, id, migrationID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Online migration stopped successfully"), nil
}

func (s *ClusterTool) getOnlineMigrationStatus(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	status, _, err := s.client.Databases.GetOnlineMigrationStatus(ctx, id)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonStatus, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonStatus)), nil
}

func (s *ClusterTool) Tools() []server.ServerTool {
	return []server.ServerTool{

		{
			Handler: s.listCluster,
			Tool: mcp.NewTool("do-dbaas-cluster-list",
				mcp.WithDescription("Get list of  Cluster"),
				mcp.WithString("page", mcp.Description("Page number for pagination (optional, integer as string)")),
				mcp.WithString("per_page", mcp.Description("Number of results per page (optional, integer as string)")),
			),
		},
		{
			Handler: s.getCluster,
			Tool: mcp.NewTool("do-dbaas-cluster-get",
				mcp.WithDescription("Get a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The ID of the cluster to retrieve")),
			),
		},
		{
			Handler: s.getCA,
			Tool: mcp.NewTool("do-dbaas-cluster-get-ca",
				mcp.WithDescription("Get the CA certificate for a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The ID of the cluster to retrieve the CA for")),
			),
		},
		{
			Handler: s.createCluster,
			Tool: mcp.NewTool("do-dbaas-cluster-create",
				mcp.WithDescription("Create a new database cluster"),
				mcp.WithString("name", mcp.Required(), mcp.Description("The name of the cluster")),
				mcp.WithString("engine", mcp.Required(), mcp.Description("The engine slug (e.g., valkey, pg, mysql, etc.)")),
				mcp.WithString("version", mcp.Required(), mcp.Description("The version of the engine")),
				mcp.WithString("region", mcp.Required(), mcp.Description("The region slug (e.g., nyc1)")),
				mcp.WithString("size", mcp.Required(), mcp.Description("The size slug (e.g., db-s-2vcpu-4gb)")),
				mcp.WithNumber("num_nodes", mcp.Required(), mcp.Description("The number of nodes")),
				mcp.WithString("tags", mcp.Description("Comma-separated tags to apply to the cluster")),
			),
		},
		{
			Handler: s.deleteCluster,
			Tool: mcp.NewTool("do-dbaas-cluster-delete",
				mcp.WithDescription("Delete a database cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The ID of the cluster to delete")),
			),
		},
		{
			Handler: s.resizeCluster,
			Tool: mcp.NewTool("do-dbaas-cluster-resize",
				mcp.WithDescription("Resize a database cluster by its ID. At least one of size, num_nodes, or storage_size_mib must be provided."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The ID of the cluster to resize")),
				mcp.WithString("size", mcp.Description("The new size slug (e.g., db-s-2vcpu-4gb)")),
				mcp.WithNumber("num_nodes", mcp.Description("The new number of nodes")),
				mcp.WithNumber("storage_size_mib", mcp.Description("The new storage size in MiB")),
			),
		},
		{
			Handler: s.listBackups,
			Tool: mcp.NewTool("do-dbaas-cluster-list-backups",
				mcp.WithDescription("List backups for a database cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The ID of the cluster to list backups for")),
				mcp.WithString("page", mcp.Description("Page number for pagination (optional, integer as string)")),
				mcp.WithString("per_page", mcp.Description("Number of results per page (optional, integer as string)")),
			),
		},
		{
			Handler: s.listOptions,
			Tool: mcp.NewTool("do-dbaas-cluster-list-options",
				mcp.WithDescription("List available database options (engines, versions, sizes, regions, etc) for DigitalOcean managed databases."),
			),
		},
		{
			Handler: s.upgradeMajorVersion,
			Tool: mcp.NewTool("do-dbaas-cluster-upgrade-major-version",
				mcp.WithDescription("Upgrade the major version of a database cluster by its ID. Requires the target version."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("version", mcp.Required(), mcp.Description("The target major version to upgrade to (e.g., 15 for PostgreSQL)")),
			),
		},
		{
			Handler: s.startOnlineMigration,
			Tool: mcp.NewTool("do-dbaas-cluster-start-online-migration",
				mcp.WithDescription("Start an online migration for a database cluster by its ID. Accepts source_json (DatabaseOnlineMigrationConfig as JSON, required), disable_ssl (optional, bool as string), and ignore_dbs (optional, comma-separated)."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("source_json", mcp.Required(), mcp.Description("DatabaseOnlineMigrationConfig as JSON (required)")),
				mcp.WithString("disable_ssl", mcp.Description("Disable SSL for migration (optional, bool as string)")),
				mcp.WithString("ignore_dbs", mcp.Description("Comma-separated list of DBs to ignore (optional)")),
			),
		},
		{
			Handler: s.stopOnlineMigration,
			Tool: mcp.NewTool("do-dbaas-cluster-stop-online-migration",
				mcp.WithDescription("Stop an online migration for a database cluster by its ID and migration_id."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("migration_id", mcp.Required(), mcp.Description("The migration ID to stop")),
			),
		},
		{
			Handler: s.getOnlineMigrationStatus,
			Tool: mcp.NewTool("do-dbaas-cluster-get-online-migration-status",
				mcp.WithDescription("Get the online migration status for a database cluster by its ID."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
			),
		},
	}
}

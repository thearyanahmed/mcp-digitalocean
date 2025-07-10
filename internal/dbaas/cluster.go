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

func (s *ClusterTool) migrateCluster(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	region, ok := args["region"].(string)
	if !ok || region == "" {
		return mcp.NewToolResultError("Target region is required"), nil
	}
	privateNetworkUUID, _ := args["private_network_uuid"].(string)

	migrateReq := &godo.DatabaseMigrateRequest{
		Region:             region,
		PrivateNetworkUUID: privateNetworkUUID,
	}

	_, err := s.client.Databases.Migrate(ctx, id, migrateReq)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Cluster migration initiated successfully"), nil
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

func (s *ClusterTool) updateMaintenance(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	day, _ := args["day"].(string)
	hour, _ := args["hour"].(string)

	if day == "" && hour == "" {
		return mcp.NewToolResultError("At least one of 'day' or 'hour' must be provided"), nil
	}

	maintenanceReq := &godo.DatabaseUpdateMaintenanceRequest{
		Day:  day,
		Hour: hour,
	}

	_, err := s.client.Databases.UpdateMaintenance(ctx, id, maintenanceReq)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Maintenance window updated successfully"), nil
}

func (s *ClusterTool) installUpdate(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, ok := req.GetArguments()["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	_, err := s.client.Databases.InstallUpdate(ctx, id)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Update installation triggered successfully"), nil
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

func (s *ClusterTool) resetUserAuth(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	user, ok := args["user"].(string)
	if !ok || user == "" {
		return mcp.NewToolResultError("User name is required"), nil
	}

	resetReq := &godo.DatabaseResetUserAuthRequest{}

	if plugin, ok := args["mysql_auth_plugin"].(string); ok && plugin != "" {
		resetReq.MySQLSettings = &godo.DatabaseMySQLUserSettings{AuthPlugin: plugin}
	}
	if settingsStr, ok := args["settings_json"].(string); ok && settingsStr != "" {
		var settings godo.DatabaseUserSettings
		err := json.Unmarshal([]byte(settingsStr), &settings)
		if err != nil {
			return mcp.NewToolResultError("Invalid settings_json: " + err.Error()), nil
		}
		resetReq.Settings = &settings
	}

	updatedUser, _, err := s.client.Databases.ResetUserAuth(ctx, id, user, resetReq)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonUser, err := json.MarshalIndent(updatedUser, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonUser)), nil
}

func (s *ClusterTool) getEvictionPolicy(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	policy, _, err := s.client.Databases.GetEvictionPolicy(ctx, id)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText(policy), nil
}

func (s *ClusterTool) setEvictionPolicy(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	policy, ok := args["policy"].(string)
	if !ok || policy == "" {
		return mcp.NewToolResultError("Eviction policy is required"), nil
	}
	_, err := s.client.Databases.SetEvictionPolicy(ctx, id, policy)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Eviction policy set successfully"), nil
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

func (s *ClusterTool) listDatabaseEvents(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	events, _, err := s.client.Databases.ListDatabaseEvents(ctx, id, opts)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonEvents, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonEvents)), nil
}

func (s *ClusterTool) Tools() []server.ServerTool {
	return []server.ServerTool{

		{
			Handler: s.listCluster,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-list",
				mcp.WithDescription("Get list of  Cluster"),
				mcp.WithString("page", mcp.Description("Page number for pagination (optional, integer as string)")),
				mcp.WithString("per_page", mcp.Description("Number of results per page (optional, integer as string)")),
			),
		},
		{
			Handler: s.getCluster,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-get",
				mcp.WithDescription("Get a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The ID of the cluster to retrieve")),
			),
		},
		{
			Handler: s.getCA,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-get-ca",
				mcp.WithDescription("Get the CA certificate for a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The ID of the cluster to retrieve the CA for")),
			),
		},
		{
			Handler: s.createCluster,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-create",
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
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-delete",
				mcp.WithDescription("Delete a database cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The ID of the cluster to delete")),
			),
		},
		{
			Handler: s.resizeCluster,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-resize",
				mcp.WithDescription("Resize a database cluster by its ID. At least one of size, num_nodes, or storage_size_mib must be provided."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The ID of the cluster to resize")),
				mcp.WithString("size", mcp.Description("The new size slug (e.g., db-s-2vcpu-4gb)")),
				mcp.WithNumber("num_nodes", mcp.Description("The new number of nodes")),
				mcp.WithNumber("storage_size_mib", mcp.Description("The new storage size in MiB")),
			),
		},
		{
			Handler: s.migrateCluster,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-migrate",
				mcp.WithDescription("Migrate a database cluster to a new region. Requires region; private_network_uuid is optional."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The ID of the cluster to migrate")),
				mcp.WithString("region", mcp.Required(), mcp.Description("The target region slug (e.g., nyc1)")),
				mcp.WithString("private_network_uuid", mcp.Description("The private network UUID to use in the target region (optional)")),
			),
		},
		{
			Handler: s.updateMaintenance,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-update-maintenance",
				mcp.WithDescription("Update the maintenance window for a database cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The ID of the cluster to update")),
				mcp.WithString("day", mcp.Description("The day of the week for maintenance (e.g., monday)")),
				mcp.WithString("hour", mcp.Description("The hour (in UTC, 24h format) for maintenance (e.g., 14:00)")),
			),
		},
		{
			Handler: s.installUpdate,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-install-update",
				mcp.WithDescription("Trigger installation of updates for a database cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The ID of the cluster to update")),
			),
		},
		{
			Handler: s.listBackups,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-list-backups",
				mcp.WithDescription("List backups for a database cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The ID of the cluster to list backups for")),
				mcp.WithString("page", mcp.Description("Page number for pagination (optional, integer as string)")),
				mcp.WithString("per_page", mcp.Description("Number of results per page (optional, integer as string)")),
			),
		},
		{
			Handler: s.resetUserAuth,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-reset-user-auth",
				mcp.WithDescription("Reset a database user's authentication for a cluster by its ID and user name"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("user", mcp.Required(), mcp.Description("The user name to reset")),
				mcp.WithString("mysql_auth_plugin", mcp.Description("MySQL auth plugin (e.g., mysql_native_password)")),
				mcp.WithString("settings_json", mcp.Description("Raw JSON for advanced DatabaseUserSettings (Kafka, OpenSearch, MongoDB, etc)")),
			),
		},
		{
			Handler: s.getEvictionPolicy,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-get-eviction-policy",
				mcp.WithDescription("Get the eviction policy for a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
			),
		},
		{
			Handler: s.setEvictionPolicy,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-set-eviction-policy",
				mcp.WithDescription("Set the eviction policy for a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("policy", mcp.Required(), mcp.Description("The eviction policy to set")),
			),
		},
		{
			Handler: s.listOptions,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-list-options",
				mcp.WithDescription("List available database options (engines, versions, sizes, regions, etc) for DigitalOcean managed databases."),
			),
		},
		{
			Handler: s.upgradeMajorVersion,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-upgrade-major-version",
				mcp.WithDescription("Upgrade the major version of a database cluster by its ID. Requires the target version."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("version", mcp.Required(), mcp.Description("The target major version to upgrade to (e.g., 15 for PostgreSQL)")),
			),
		},
		{
			Handler: s.listDatabaseEvents,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-list-database-events",
				mcp.WithDescription("List database events for a cluster by its ID. Supports all ListOptions: page, per_page, with_projects, only_deployed, public_only, usecases (comma-separated)."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("page", mcp.Description("Page number for pagination (optional, integer as string)")),
				mcp.WithString("per_page", mcp.Description("Number of results per page (optional, integer as string)")),
				mcp.WithString("with_projects", mcp.Description("Whether to include project_id fields (optional, bool as string)")),
				mcp.WithString("only_deployed", mcp.Description("Only list deployed agents (optional, bool as string)")),
				mcp.WithString("public_only", mcp.Description("Include only public models (optional, bool as string)")),
				mcp.WithString("usecases", mcp.Description("Comma-separated usecases to filter (optional)")),
			),
		},
	}
}

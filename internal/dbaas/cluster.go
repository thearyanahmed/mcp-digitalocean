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

func (s *ClusterTool) getUser(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
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

func (s *ClusterTool) listUsers(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

func (s *ClusterTool) createUser(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
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

func (s *ClusterTool) updateUser(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
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

func (s *ClusterTool) deleteUser(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
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

func (s *ClusterTool) listDBs(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

func (s *ClusterTool) createDB(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

func (s *ClusterTool) getDB(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

func (s *ClusterTool) deleteDB(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

func (s *ClusterTool) getSQLMode(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	mode, _, err := s.client.Databases.GetSQLMode(ctx, id)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText(mode), nil
}

func (s *ClusterTool) setSQLMode(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	modesStr, ok := args["modes"].(string)
	if !ok || modesStr == "" {
		return mcp.NewToolResultError("SQL modes are required (comma-separated)"), nil
	}
	modes := []string{}
	for _, m := range strings.Split(modesStr, ",") {
		m = strings.TrimSpace(m)
		if m != "" {
			modes = append(modes, m)
		}
	}
	_, err := s.client.Databases.SetSQLMode(ctx, id, modes...)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("SQL mode set successfully"), nil
}

func (s *ClusterTool) getFirewallRules(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	rules, _, err := s.client.Databases.GetFirewallRules(ctx, id)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonRules, err := json.MarshalIndent(rules, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonRules)), nil
}

func (s *ClusterTool) updateFirewallRules(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	rulesStr, ok := args["rules_json"].(string)
	if !ok || rulesStr == "" {
		return mcp.NewToolResultError("rules_json is required (JSON array of firewall rules)"), nil
	}
	var rules []*godo.DatabaseFirewallRule
	err := json.Unmarshal([]byte(rulesStr), &rules)
	if err != nil {
		return mcp.NewToolResultError("Invalid rules_json: " + err.Error()), nil
	}
	updateReq := &godo.DatabaseUpdateFirewallRulesRequest{Rules: rules}
	_, err = s.client.Databases.UpdateFirewallRules(ctx, id, updateReq)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Firewall rules updated successfully"), nil
}

func (s *ClusterTool) getPostgreSQLConfig(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	cfg, _, err := s.client.Databases.GetPostgreSQLConfig(ctx, id)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonCfg, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonCfg)), nil
}

func (s *ClusterTool) updatePostgreSQLConfig(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	configStr, ok := args["config_json"].(string)
	if !ok || configStr == "" {
		return mcp.NewToolResultError("config_json is required (JSON for PostgreSQLConfig)"), nil
	}
	var config godo.PostgreSQLConfig
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		return mcp.NewToolResultError("Invalid config_json: " + err.Error()), nil
	}
	_, err = s.client.Databases.UpdatePostgreSQLConfig(ctx, id, &config)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("PostgreSQL config updated successfully"), nil
}

func (s *ClusterTool) getRedisConfig(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	cfg, _, err := s.client.Databases.GetRedisConfig(ctx, id)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonCfg, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonCfg)), nil
}

func (s *ClusterTool) updateRedisConfig(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	configStr, ok := args["config_json"].(string)
	if !ok || configStr == "" {
		return mcp.NewToolResultError("config_json is required (JSON for RedisConfig)"), nil
	}
	var config godo.RedisConfig
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		return mcp.NewToolResultError("Invalid config_json: " + err.Error()), nil
	}
	_, err = s.client.Databases.UpdateRedisConfig(ctx, id, &config)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Redis config updated successfully"), nil
}

func (s *ClusterTool) getMySQLConfig(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	cfg, _, err := s.client.Databases.GetMySQLConfig(ctx, id)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonCfg, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonCfg)), nil
}

func (s *ClusterTool) updateMySQLConfig(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	configStr, ok := args["config_json"].(string)
	if !ok || configStr == "" {
		return mcp.NewToolResultError("config_json is required (JSON for MySQLConfig)"), nil
	}
	var config godo.MySQLConfig
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		return mcp.NewToolResultError("Invalid config_json: " + err.Error()), nil
	}
	_, err = s.client.Databases.UpdateMySQLConfig(ctx, id, &config)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("MySQL config updated successfully"), nil
}

func (s *ClusterTool) getMongoDBConfig(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	cfg, _, err := s.client.Databases.GetMongoDBConfig(ctx, id)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonCfg, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonCfg)), nil
}

func (s *ClusterTool) updateMongoDBConfig(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	configStr, ok := args["config_json"].(string)
	if !ok || configStr == "" {
		return mcp.NewToolResultError("config_json is required (JSON for MongoDBConfig)"), nil
	}
	var config godo.MongoDBConfig
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		return mcp.NewToolResultError("Invalid config_json: " + err.Error()), nil
	}
	_, err = s.client.Databases.UpdateMongoDBConfig(ctx, id, &config)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("MongoDB config updated successfully"), nil
}

func (s *ClusterTool) getOpensearchConfig(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	cfg, _, err := s.client.Databases.GetOpensearchConfig(ctx, id)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonCfg, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonCfg)), nil
}

func (s *ClusterTool) updateOpensearchConfig(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	configStr, ok := args["config_json"].(string)
	if !ok || configStr == "" {
		return mcp.NewToolResultError("config_json is required (JSON for OpensearchConfig)"), nil
	}
	var config godo.OpensearchConfig
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		return mcp.NewToolResultError("Invalid config_json: " + err.Error()), nil
	}
	_, err = s.client.Databases.UpdateOpensearchConfig(ctx, id, &config)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Opensearch config updated successfully"), nil
}

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

func (s *ClusterTool) getMetricsCredentials(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	creds, _, err := s.client.Databases.GetMetricsCredentials(ctx)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonCreds, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonCreds)), nil
}

func (s *ClusterTool) updateMetricsCredentials(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	credsStr, ok := args["credentials_json"].(string)
	if !ok || credsStr == "" {
		return mcp.NewToolResultError("credentials_json is required (JSON for DatabaseMetricsCredentials)"), nil
	}
	var creds godo.DatabaseMetricsCredentials
	err := json.Unmarshal([]byte(credsStr), &creds)
	if err != nil {
		return mcp.NewToolResultError("Invalid credentials_json: " + err.Error()), nil
	}
	updateReq := &godo.DatabaseUpdateMetricsCredentialsRequest{Credentials: &creds}
	_, err = s.client.Databases.UpdateMetricsCredentials(ctx, updateReq)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Metrics credentials updated successfully"), nil
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

func (s *ClusterTool) listIndexes(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

func (s *ClusterTool) deleteIndex(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

func (s *ClusterTool) createLogsink(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	name, ok := args["sink_name"].(string)
	if !ok || name == "" {
		return mcp.NewToolResultError("sink_name is required"), nil
	}
	typeStr, ok := args["sink_type"].(string)
	if !ok || typeStr == "" {
		return mcp.NewToolResultError("sink_type is required"), nil
	}
	configStr, ok := args["config_json"].(string)
	if !ok || configStr == "" {
		return mcp.NewToolResultError("config_json is required (JSON for DatabaseLogsinkConfig)"), nil
	}
	var config godo.DatabaseLogsinkConfig
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		return mcp.NewToolResultError("Invalid config_json: " + err.Error()), nil
	}
	createReq := &godo.DatabaseCreateLogsinkRequest{
		Name:   name,
		Type:   typeStr,
		Config: &config,
	}
	logsink, _, err := s.client.Databases.CreateLogsink(ctx, id, createReq)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonLogsink, err := json.MarshalIndent(logsink, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonLogsink)), nil
}

func (s *ClusterTool) getLogsink(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	logsinkID, ok := args["logsink_id"].(string)
	if !ok || logsinkID == "" {
		return mcp.NewToolResultError("logsink_id is required"), nil
	}
	logsink, _, err := s.client.Databases.GetLogsink(ctx, id, logsinkID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonLogsink, err := json.MarshalIndent(logsink, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonLogsink)), nil
}

func (s *ClusterTool) listLogsinks(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
	logsinks, _, err := s.client.Databases.ListLogsinks(ctx, id, opts)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonLogsinks, err := json.MarshalIndent(logsinks, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonLogsinks)), nil
}

func (s *ClusterTool) updateLogsink(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	logsinkID, ok := args["logsink_id"].(string)
	if !ok || logsinkID == "" {
		return mcp.NewToolResultError("logsink_id is required"), nil
	}
	configStr, ok := args["config_json"].(string)
	if !ok || configStr == "" {
		return mcp.NewToolResultError("config_json is required (JSON for DatabaseLogsinkConfig)"), nil
	}
	var config godo.DatabaseLogsinkConfig
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		return mcp.NewToolResultError("Invalid config_json: " + err.Error()), nil
	}
	updateReq := &godo.DatabaseUpdateLogsinkRequest{Config: &config}
	_, err = s.client.Databases.UpdateLogsink(ctx, id, logsinkID, updateReq)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Logsink updated successfully"), nil
}

func (s *ClusterTool) deleteLogsink(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, ok := args["ID"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Cluster ID is required"), nil
	}
	logsinkID, ok := args["logsink_id"].(string)
	if !ok || logsinkID == "" {
		return mcp.NewToolResultError("logsink_id is required"), nil
	}
	_, err := s.client.Databases.DeleteLogsink(ctx, id, logsinkID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Logsink deleted successfully"), nil
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
			Handler: s.getUser,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-get-user",
				mcp.WithDescription("Get a database user by cluster ID and user name"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster ID (UUID)")),
				mcp.WithString("user", mcp.Required(), mcp.Description("The user name")),
			),
		},
		{
			Handler: s.listUsers,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-list-users",
				mcp.WithDescription("List database users for a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster ID (UUID)")),
				mcp.WithString("page", mcp.Description("Page number for pagination (optional)")),
				mcp.WithString("per_page", mcp.Description("Number of results per page (optional)")),
			),
		},
		{
			Handler: s.createUser,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-create-user",
				mcp.WithDescription("Create a database user for a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster ID (UUID)")),
				mcp.WithString("name", mcp.Required(), mcp.Description("The user name")),
				mcp.WithString("mysql_auth_plugin", mcp.Description("MySQL auth plugin (optional, e.g., mysql_native_password)")),
				mcp.WithString("settings_json", mcp.Description("Raw JSON for DatabaseUserSettings (optional)")),
			),
		},
		{
			Handler: s.updateUser,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-update-user",
				mcp.WithDescription("Update a database user for a cluster by its ID and user name"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster ID (UUID)")),
				mcp.WithString("user", mcp.Required(), mcp.Description("The user name")),
				mcp.WithString("settings_json", mcp.Description("Raw JSON for DatabaseUserSettings (optional)")),
			),
		},
		{
			Handler: s.deleteUser,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-delete-user",
				mcp.WithDescription("Delete a database user by cluster ID and user name"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("user", mcp.Required(), mcp.Description("The user name to delete")),
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
		{
			Handler: s.listPools,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-list-pools",
				mcp.WithDescription("List connection pools for a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("page", mcp.Description("Page number for pagination")),
				mcp.WithString("per_page", mcp.Description("Number of results per page")),
			),
		},
		{
			Handler: s.createPool,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-create-pool",
				mcp.WithDescription("Create a connection pool for a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("user", mcp.Required(), mcp.Description("The user for the pool")),
				mcp.WithString("name", mcp.Required(), mcp.Description("The pool name")),
				mcp.WithString("database", mcp.Required(), mcp.Description("The database for the pool")),
				mcp.WithString("mode", mcp.Required(), mcp.Description("The pool mode")),
				mcp.WithNumber("size", mcp.Required(), mcp.Description("The pool size (number of connections)")),
			),
		},
		{
			Handler: s.getPool,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-get-pool",
				mcp.WithDescription("Get a connection pool for a cluster by its ID and pool name"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("name", mcp.Required(), mcp.Description("The pool name to get")),
			),
		},
		{
			Handler: s.deletePool,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-delete-pool",
				mcp.WithDescription("Delete a connection pool for a cluster by its ID and pool name"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("name", mcp.Required(), mcp.Description("The pool name to delete")),
			),
		},
		{
			Handler: s.updatePool,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-update-pool",
				mcp.WithDescription("Update a connection pool for a cluster by its ID and pool name"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("name", mcp.Required(), mcp.Description("The pool name to update")),
				mcp.WithString("user", mcp.Description("The user for the pool (optional)")),
				mcp.WithString("database", mcp.Required(), mcp.Description("The database for the pool")),
				mcp.WithString("mode", mcp.Required(), mcp.Description("The pool mode")),
				mcp.WithNumber("size", mcp.Required(), mcp.Description("The pool size (number of connections)")),
			),
		},
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
			Handler: s.getSQLMode,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-get-sql-mode",
				mcp.WithDescription("Get the SQL mode for a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
			),
		},
		{
			Handler: s.setSQLMode,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-set-sql-mode",
				mcp.WithDescription("Set the SQL mode for a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("modes", mcp.Required(), mcp.Description("Comma-separated SQL modes to set")),
			),
		},
		{
			Handler: s.getFirewallRules,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-get-firewall-rules",
				mcp.WithDescription("Get the firewall rules for a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
			),
		},
		{
			Handler: s.updateFirewallRules,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-update-firewall-rules",
				mcp.WithDescription("Update the firewall rules for a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("rules_json", mcp.Required(), mcp.Description("JSON array of firewall rules to set")),
			),
		},
		{
			Handler: s.getPostgreSQLConfig,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-get-postgresql-config",
				mcp.WithDescription("Get the PostgreSQL config for a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
			),
		},
		{
			Handler: s.updatePostgreSQLConfig,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-update-postgresql-config",
				mcp.WithDescription("Update the PostgreSQL config for a cluster by its ID. Accepts a JSON string for the config."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("config_json", mcp.Required(), mcp.Description("JSON for the PostgreSQLConfig to set")),
			),
		},
		{
			Handler: s.getRedisConfig,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-get-redis-config",
				mcp.WithDescription("Get the Redis config for a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
			),
		},
		{
			Handler: s.updateRedisConfig,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-update-redis-config",
				mcp.WithDescription("Update the Redis config for a cluster by its ID. Accepts a JSON string for the config."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("config_json", mcp.Required(), mcp.Description("JSON for the RedisConfig to set")),
			),
		},
		{
			Handler: s.getMySQLConfig,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-get-mysql-config",
				mcp.WithDescription("Get the MySQL config for a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
			),
		},
		{
			Handler: s.updateMySQLConfig,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-update-mysql-config",
				mcp.WithDescription("Update the MySQL config for a cluster by its ID. Accepts a JSON string for the config."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("config_json", mcp.Required(), mcp.Description("JSON for the MySQLConfig to set")),
			),
		},
		{
			Handler: s.getMongoDBConfig,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-get-mongodb-config",
				mcp.WithDescription("Get the MongoDB config for a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
			),
		},
		{
			Handler: s.updateMongoDBConfig,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-update-mongodb-config",
				mcp.WithDescription("Update the MongoDB config for a cluster by its ID. Accepts a JSON string for the config."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("config_json", mcp.Required(), mcp.Description("JSON for the MongoDBConfig to set")),
			),
		},
		{
			Handler: s.getOpensearchConfig,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-get-opensearch-config",
				mcp.WithDescription("Get the Opensearch config for a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
			),
		},
		{
			Handler: s.updateOpensearchConfig,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-update-opensearch-config",
				mcp.WithDescription("Update the Opensearch config for a cluster by its ID. Accepts a JSON string for the config."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("config_json", mcp.Required(), mcp.Description("JSON for the OpensearchConfig to set")),
			),
		},
		{
			Handler: s.getKafkaConfig,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-get-kafka-config",
				mcp.WithDescription("Get the Kafka config for a cluster by its ID"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
			),
		},
		{
			Handler: s.updateKafkaConfig,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-update-kafka-config",
				mcp.WithDescription("Update the Kafka config for a cluster by its ID. Accepts a JSON string for the config."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("config_json", mcp.Required(), mcp.Description("JSON for the KafkaConfig to set")),
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
			Handler: s.listTopics,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-list-topics",
				mcp.WithDescription("List topics for a database cluster by its ID (Kafka clusters). Supports all ListOptions: page, per_page, with_projects, only_deployed, public_only, usecases (comma-separated)."),
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
			Handler: s.createTopic,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-create-topic",
				mcp.WithDescription("Create a topic for a Kafka database cluster by its ID. Accepts name (required), partition_count, replication_factor, and config_json (TopicConfig as JSON, all optional)."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("name", mcp.Required(), mcp.Description("The topic name to create")),
				mcp.WithString("partition_count", mcp.Description("Number of partitions (optional, integer as string)")),
				mcp.WithString("replication_factor", mcp.Description("Replication factor (optional, integer as string)")),
				mcp.WithString("config_json", mcp.Description("TopicConfig as JSON (optional)")),
			),
		},
		{
			Handler: s.getTopic,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-get-topic",
				mcp.WithDescription("Get a topic for a Kafka database cluster by its ID and topic name."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("name", mcp.Required(), mcp.Description("The topic name to get")),
			),
		},
		{
			Handler: s.deleteTopic,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-delete-topic",
				mcp.WithDescription("Delete a topic for a Kafka database cluster by its ID and topic name."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("name", mcp.Required(), mcp.Description("The topic name to delete")),
			),
		},
		{
			Handler: s.updateTopic,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-update-topic",
				mcp.WithDescription("Update a topic for a Kafka database cluster by its ID and topic name. Accepts partition_count, replication_factor, and config_json (TopicConfig as JSON, all optional)."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("name", mcp.Required(), mcp.Description("The topic name to update")),
				mcp.WithString("partition_count", mcp.Description("Number of partitions (optional, integer as string)")),
				mcp.WithString("replication_factor", mcp.Description("Replication factor (optional, integer as string)")),
				mcp.WithString("config_json", mcp.Description("TopicConfig as JSON (optional)")),
			),
		},
		{
			Handler: s.getMetricsCredentials,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-get-metrics-credentials",
				mcp.WithDescription("Get metrics credentials for DigitalOcean managed databases (no arguments required)."),
			),
		},
		{
			Handler: s.updateMetricsCredentials,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-update-metrics-credentials",
				mcp.WithDescription("Update metrics credentials for DigitalOcean managed databases. Accepts credentials_json (JSON for DatabaseMetricsCredentials)."),
				mcp.WithString("credentials_json", mcp.Required(), mcp.Description("JSON for the DatabaseMetricsCredentials to set")),
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
		{
			Handler: s.createLogsink,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-create-logsink",
				mcp.WithDescription("Create a logsink for a database cluster by its ID. Accepts sink_name, sink_type, and config_json (DatabaseLogsinkConfig as JSON, all required)."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("sink_name", mcp.Required(), mcp.Description("The logsink name to create")),
				mcp.WithString("sink_type", mcp.Required(), mcp.Description("The logsink type (e.g., opensearch, datadog, logtail, papertrail)")),
				mcp.WithString("config_json", mcp.Required(), mcp.Description("DatabaseLogsinkConfig as JSON (required)")),
			),
		},
		{
			Handler: s.getLogsink,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-get-logsink",
				mcp.WithDescription("Get a logsink for a database cluster by its ID and logsink_id."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("logsink_id", mcp.Required(), mcp.Description("The logsink ID to get")),
			),
		},
		{
			Handler: s.listLogsinks,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-list-logsinks",
				mcp.WithDescription("List logsinks for a database cluster by its ID. Supports pagination: page, per_page (optional, integer as string)."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("page", mcp.Description("Page number for pagination (optional, integer as string)")),
				mcp.WithString("per_page", mcp.Description("Number of results per page (optional, integer as string)")),
			),
		},
		{
			Handler: s.updateLogsink,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-update-logsink",
				mcp.WithDescription("Update a logsink for a database cluster by its ID and logsink_id. Accepts config_json (DatabaseLogsinkConfig as JSON, required)."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("logsink_id", mcp.Required(), mcp.Description("The logsink ID to update")),
				mcp.WithString("config_json", mcp.Required(), mcp.Description("DatabaseLogsinkConfig as JSON (required)")),
			),
		},
		{
			Handler: s.deleteLogsink,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-delete-logsink",
				mcp.WithDescription("Delete a logsink for a database cluster by its ID and logsink_id."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("logsink_id", mcp.Required(), mcp.Description("The logsink ID to delete")),
			),
		},
		{
			Handler: s.startOnlineMigration,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-start-online-migration",
				mcp.WithDescription("Start an online migration for a database cluster by its ID. Accepts source_json (DatabaseOnlineMigrationConfig as JSON, required), disable_ssl (optional, bool as string), and ignore_dbs (optional, comma-separated)."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("source_json", mcp.Required(), mcp.Description("DatabaseOnlineMigrationConfig as JSON (required)")),
				mcp.WithString("disable_ssl", mcp.Description("Disable SSL for migration (optional, bool as string)")),
				mcp.WithString("ignore_dbs", mcp.Description("Comma-separated list of DBs to ignore (optional)")),
			),
		},
		{
			Handler: s.stopOnlineMigration,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-stop-online-migration",
				mcp.WithDescription("Stop an online migration for a database cluster by its ID and migration_id."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
				mcp.WithString("migration_id", mcp.Required(), mcp.Description("The migration ID to stop")),
			),
		},
		{
			Handler: s.getOnlineMigrationStatus,
			Tool: mcp.NewTool("digitalocean-dbaas-cluster-get-online-migration-status",
				mcp.WithDescription("Get the online migration status for a database cluster by its ID."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("The cluster UUID")),
			),
		},
	}
}

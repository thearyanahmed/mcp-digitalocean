package droplet

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// DropletTool provides droplet management tools
type DropletTool struct {
	client *godo.Client
}

// NewDropletTool creates a new droplet tool
func NewDropletTool(client *godo.Client) *DropletTool {
	return &DropletTool{
		client: client,
	}
}

// CreateDroplet creates a new droplet
func (d *DropletTool) createDroplet(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	dropletName := args["Name"].(string)
	size := args["Size"].(string)
	imageID := args["ImageID"].(float64)
	region := args["Region"].(string)
	backup, _ := args["Backup"].(bool)         // Defaults to false
	monitoring, _ := args["Monitoring"].(bool) // Defaults to false
	// Create the droplet
	dropletCreateRequest := &godo.DropletCreateRequest{
		Name:       dropletName,
		Size:       size,
		Image:      godo.DropletCreateImage{ID: int(imageID)},
		Region:     region,
		Backups:    backup,
		Monitoring: monitoring,
	}
	droplet, _, err := d.client.Droplets.Create(ctx, dropletCreateRequest)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("droplet create", err), nil
	}
	jsonDroplet, err := json.MarshalIndent(droplet, "", "  ")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("json marshal", err), nil
	}
	return mcp.NewToolResultText(string(jsonDroplet)), nil
}

// deleteDroplet deletes a droplet
func (d *DropletTool) deleteDroplet(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dropletID := req.GetArguments()["ID"].(float64)
	_, err := d.client.Droplets.Delete(ctx, int(dropletID))
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	return mcp.NewToolResultText("Droplet deleted successfully"), nil
}

// getDropletNeighbors gets a droplet's neighbors
func (d *DropletTool) getDropletNeighbors(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dropletID := req.GetArguments()["ID"].(float64)
	neighbors, _, err := d.client.Droplets.Neighbors(ctx, int(dropletID))
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonNeighbors, err := json.MarshalIndent(neighbors, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	return mcp.NewToolResultText(string(jsonNeighbors)), nil
}

// enablePrivateNetworking enables private networking on a droplet
func (d *DropletTool) enablePrivateNetworking(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dropletID := req.GetArguments()["ID"].(float64)
	action, _, err := d.client.DropletActions.EnablePrivateNetworking(ctx, int(dropletID))
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonAction, err := json.MarshalIndent(action, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	return mcp.NewToolResultText(string(jsonAction)), nil
}

// getDropletKernels gets available kernels for a droplet
func (d *DropletTool) getDropletKernels(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dropletID := req.GetArguments()["ID"].(float64)

	// Use list options to get all kernels
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 100,
	}

	kernels, _, err := d.client.Droplets.Kernels(ctx, int(dropletID), opt)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonKernels, err := json.MarshalIndent(kernels, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	return mcp.NewToolResultText(string(jsonKernels)), nil
}

// rebootDroplet reboots a droplet
func (d *DropletTool) rebootDroplet(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dropletID := req.GetArguments()["ID"].(float64)
	action, _, err := d.client.DropletActions.Reboot(ctx, int(dropletID))
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonAction, err := json.MarshalIndent(action, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	return mcp.NewToolResultText(string(jsonAction)), nil
}

// passwordResetDroplet resets the password for a droplet
func (d *DropletTool) passwordResetDroplet(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dropletID := req.GetArguments()["ID"].(float64)
	action, _, err := d.client.DropletActions.PasswordReset(ctx, int(dropletID))
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonAction, err := json.MarshalIndent(action, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	return mcp.NewToolResultText(string(jsonAction)), nil
}

// RebuildByImageSlugDroplet rebuilds a droplet using an image slug
func (d *DropletTool) rebuildByImageSlugDroplet(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dropletID := req.GetArguments()["ID"].(float64)
	imageSlug := req.GetArguments()["ImageSlug"].(string)
	action, _, err := d.client.DropletActions.RebuildByImageSlug(ctx, int(dropletID), imageSlug)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonAction, err := json.MarshalIndent(action, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	return mcp.NewToolResultText(string(jsonAction)), nil
}

// powerCycleByTag power cycles droplets by tag
func (d *DropletTool) powerCycleByTag(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tag := req.GetArguments()["Tag"].(string)
	actions, _, err := d.client.DropletActions.PowerCycleByTag(ctx, tag)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonActions, err := json.MarshalIndent(actions, "", "  ")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("marshal error", err), nil
	}

	return mcp.NewToolResultText(string(jsonActions)), nil
}

// powerOnByTag powers on droplets by tag
func (d *DropletTool) powerOnByTag(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tag := req.GetArguments()["Tag"].(string)
	actions, _, err := d.client.DropletActions.PowerOnByTag(ctx, tag)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonActions, err := json.MarshalIndent(actions, "", "  ")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("marshal error", err), nil
	}

	return mcp.NewToolResultText(string(jsonActions)), nil
}

// powerOffByTag powers off droplets by tag
func (d *DropletTool) powerOffByTag(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tag := req.GetArguments()["Tag"].(string)
	actions, _, err := d.client.DropletActions.PowerOffByTag(ctx, tag)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonActions, err := json.MarshalIndent(actions, "", "  ")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("marshal error", err), nil
	}

	return mcp.NewToolResultText(string(jsonActions)), nil
}

// shutdownByTag shuts down droplets by tag
func (d *DropletTool) shutdownByTag(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tag := req.GetArguments()["Tag"].(string)
	actions, _, err := d.client.DropletActions.ShutdownByTag(ctx, tag)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonActions, err := json.MarshalIndent(actions, "", "  ")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("marshal error", err), nil
	}

	return mcp.NewToolResultText(string(jsonActions)), nil
}

// enableBackupsByTag enables backups on droplets by tag
func (d *DropletTool) enableBackupsByTag(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tag := req.GetArguments()["Tag"].(string)
	actions, _, err := d.client.DropletActions.EnableBackupsByTag(ctx, tag)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonActions, err := json.MarshalIndent(actions, "", "  ")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("marshal error", err), nil
	}

	return mcp.NewToolResultText(string(jsonActions)), nil
}

// disableBackupsByTag disables backups on droplets by tag
func (d *DropletTool) disableBackupsByTag(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tag := req.GetArguments()["Tag"].(string)
	actions, _, err := d.client.DropletActions.DisableBackupsByTag(ctx, tag)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonActions, err := json.MarshalIndent(actions, "", "  ")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("marshal error", err), nil
	}

	return mcp.NewToolResultText(string(jsonActions)), nil
}

// snapshotByTag takes a snapshot of droplets by tag
func (d *DropletTool) snapshotByTag(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tag := req.GetArguments()["Tag"].(string)
	name := req.GetArguments()["Name"].(string)
	actions, _, err := d.client.DropletActions.SnapshotByTag(ctx, tag, name)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonActions, err := json.MarshalIndent(actions, "", "  ")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("marshal error", err), nil
	}

	return mcp.NewToolResultText(string(jsonActions)), nil
}

// enableIPv6ByTag enables IPv6 on droplets by tag
func (d *DropletTool) enableIPv6ByTag(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tag := req.GetArguments()["Tag"].(string)
	actions, _, err := d.client.DropletActions.EnableIPv6ByTag(ctx, tag)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonActions, err := json.MarshalIndent(actions, "", "  ")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("marshal error", err), nil
	}

	return mcp.NewToolResultText(string(jsonActions)), nil
}

// enablePrivateNetworkingByTag enables private networking on droplets by tag
func (d *DropletTool) enablePrivateNetworkingByTag(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tag := req.GetArguments()["Tag"].(string)
	actions, _, err := d.client.DropletActions.EnablePrivateNetworkingByTag(ctx, tag)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonActions, err := json.MarshalIndent(actions, "", "  ")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("marshal error", err), nil
	}

	return mcp.NewToolResultText(string(jsonActions)), nil
}

// Tools returns a list of tool functions
func (d *DropletTool) getDropletByID(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, ok := req.GetArguments()["ID"].(float64)
	if !ok {
		return mcp.NewToolResultError("Droplet ID is required"), nil
	}
	droplet, _, err := d.client.Droplets.Get(ctx, int(id))
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonData, err := json.MarshalIndent(droplet, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonData)), nil
}

func (d *DropletTool) getDropletActionByID(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dropletID, ok := req.GetArguments()["DropletID"].(float64)
	if !ok {
		return mcp.NewToolResultError("DropletID is required"), nil
	}
	actionID, ok := req.GetArguments()["ActionID"].(float64)
	if !ok {
		return mcp.NewToolResultError("ActionID is required"), nil
	}
	action, _, err := d.client.DropletActions.Get(ctx, int(dropletID), int(actionID))
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}
	jsonData, err := json.MarshalIndent(action, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return mcp.NewToolResultText(string(jsonData)), nil
}

func (d *DropletTool) Tools() []server.ServerTool {
	tools := []server.ServerTool{
		{
			Handler: d.createDroplet,
			Tool: mcp.NewTool("digitalocean-droplet-create",
				mcp.WithDescription("Create a new droplet"),
				mcp.WithString("Name", mcp.Required(), mcp.Description("Name of the droplet")),
				mcp.WithString("Size", mcp.Required(), mcp.Description("Slug of the droplet size (e.g., s-1vcpu-1gb)")),
				mcp.WithNumber("ImageID", mcp.Required(), mcp.Description("ID of the image to use")),
				mcp.WithString("Region", mcp.Required(), mcp.Description("Slug of the region (e.g., nyc3)")),
				mcp.WithBoolean("Backup", mcp.DefaultBool(false), mcp.Description("Whether to enable backups")),
				mcp.WithBoolean("Monitoring", mcp.DefaultBool(false), mcp.Description("Whether to enable monitoring")),
			),
		},
		{
			Handler: d.deleteDroplet,
			Tool: mcp.NewTool("digitalocean-droplet-delete",
				mcp.WithDescription("Delete a droplet"),
				mcp.WithNumber("ID", mcp.Required(), mcp.Description("ID of the droplet to delete")),
			),
		},
		{
			Handler: d.enablePrivateNetworking,
			Tool: mcp.NewTool("digitalocean-droplet-enable-private-net",
				mcp.WithDescription("Enable private networking on a droplet"),
				mcp.WithNumber("ID", mcp.Required(), mcp.Description("ID of the droplet")),
			),
		},
		{
			Handler: d.getDropletKernels,
			Tool: mcp.NewTool("digitalocean-droplet-get-kernels",
				mcp.WithDescription("Get available kernels for a droplet"),
				mcp.WithNumber("ID", mcp.Required(), mcp.Description("ID of the droplet")),
			),
		},
		{
			Handler: d.rebootDroplet,
			Tool: mcp.NewTool("digitalocean-droplet-reboot",
				mcp.WithDescription("Reboot a droplet"),
				mcp.WithNumber("ID", mcp.Required(), mcp.Description("ID of the droplet to reboot")),
			),
		},
		{
			Handler: d.passwordResetDroplet,
			Tool: mcp.NewTool("digitalocean-droplet-password-reset",
				mcp.WithDescription("Reset password for a droplet"),
				mcp.WithNumber("ID", mcp.Required(), mcp.Description("ID of the droplet")),
			),
		},
		{
			Handler: d.rebuildByImageSlugDroplet,
			Tool: mcp.NewTool("digitalocean-droplet-rebuild-by-slug",
				mcp.WithDescription("Rebuild a droplet using an image slug"),
				mcp.WithNumber("ID", mcp.Required(), mcp.Description("ID of the droplet to rebuild")),
				mcp.WithString("ImageSlug", mcp.Required(), mcp.Description("Slug of the image to rebuild from")),
			),
		},
		{
			Handler: d.powerCycleByTag,
			Tool: mcp.NewTool("digitalocean-droplet-power-cycle-by-tag",
				mcp.WithDescription("Power cycle droplets by tag"),
				mcp.WithString("Tag", mcp.Required(), mcp.Description("Tag of the droplets to power cycle")),
			),
		},
		{
			Handler: d.powerOnByTag,
			Tool: mcp.NewTool("digitalocean-droplet-power-on-by-tag",
				mcp.WithDescription("Power on droplets by tag"),
				mcp.WithString("Tag", mcp.Required(), mcp.Description("Tag of the droplets to power on")),
			),
		},
		{
			Handler: d.powerOffByTag,
			Tool: mcp.NewTool("digitalocean-droplet-power-off-by-tag",
				mcp.WithDescription("Power off droplets by tag"),
				mcp.WithString("Tag", mcp.Required(), mcp.Description("Tag of the droplets to power off")),
			),
		},
		{
			Handler: d.shutdownByTag,
			Tool: mcp.NewTool("digitalocean-droplet-shutdown-by-tag",
				mcp.WithDescription("Shutdown droplets by tag"),
				mcp.WithString("Tag", mcp.Required(), mcp.Description("Tag of the droplets to shutdown")),
			),
		},
		{
			Handler: d.enableBackupsByTag,
			Tool: mcp.NewTool("digitalocean-droplet-enable-backups-by-tag",
				mcp.WithDescription("Enable backups on droplets by tag"),
				mcp.WithString("Tag", mcp.Required(), mcp.Description("Tag of the droplets")),
			),
		},
		{
			Handler: d.disableBackupsByTag,
			Tool: mcp.NewTool("digitalocean-droplet-disable-backups-by-tag",
				mcp.WithDescription("Disable backups on droplets by tag"),
				mcp.WithString("Tag", mcp.Required(), mcp.Description("Tag of the droplets")),
			),
		},
		{
			Handler: d.snapshotByTag,
			Tool: mcp.NewTool("digitalocean-droplet-snapshot-by-tag",
				mcp.WithDescription("Take a snapshot of droplets by tag"),
				mcp.WithString("Tag", mcp.Required(), mcp.Description("Tag of the droplets")),
				mcp.WithString("Name", mcp.Required(), mcp.Description("Name for the snapshot")),
			),
		},
		{
			Handler: d.enableIPv6ByTag,
			Tool: mcp.NewTool("digitalocean-droplet-enable-ipv6-by-tag",
				mcp.WithDescription("Enable IPv6 on droplets by tag"),
				mcp.WithString("Tag", mcp.Required(), mcp.Description("Tag of the droplets")),
			),
		},
		{
			Handler: d.enablePrivateNetworkingByTag,
			Tool: mcp.NewTool("digitalocean-droplet-enable-private-net-by-tag",
				mcp.WithDescription("Enable private networking on droplets by tag"),
				mcp.WithString("Tag", mcp.Required(), mcp.Description("Tag of the droplets")),
			),
		},
		{
			Handler: d.getDropletByID,
			Tool: mcp.NewTool("digitalocean-droplet-get",
				mcp.WithDescription("Get a droplet by its ID"),
				mcp.WithNumber("ID", mcp.Required(), mcp.Description("Droplet ID")),
			),
		},
		{
			Handler: d.getDropletActionByID,
			Tool: mcp.NewTool("digitalocean-droplet-action-get",
				mcp.WithDescription("Get a droplet action by droplet ID and action ID"),
				mcp.WithNumber("DropletID", mcp.Required(), mcp.Description("Droplet ID")),
				mcp.WithNumber("ActionID", mcp.Required(), mcp.Description("Action ID")),
			),
		},
	}
	return tools
}

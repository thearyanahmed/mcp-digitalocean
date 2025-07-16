// Package droplet provides tools for managing droplet actions
package droplet

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// DropletActionsTool provides tools for droplet actions
type DropletActionsTool struct {
	client *godo.Client
}

// NewDropletActionsTool creates a new droplet actions tool
func NewDropletActionsTool(client *godo.Client) *DropletActionsTool {
	return &DropletActionsTool{
		client: client,
	}
}

// powerCycleDroplet power cycles a droplet
func (d *DropletActionsTool) powerCycleDroplet(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dropletID := req.GetArguments()["ID"].(float64)
	action, _, err := d.client.DropletActions.PowerCycle(ctx, int(dropletID))
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonAction, err := json.MarshalIndent(action, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	return mcp.NewToolResultText(string(jsonAction)), nil
}

// powerOnDroplet powers on a droplet
func (d *DropletActionsTool) powerOnDroplet(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dropletID := req.GetArguments()["ID"].(float64)
	action, _, err := d.client.DropletActions.PowerOn(ctx, int(dropletID))
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonAction, err := json.MarshalIndent(action, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	return mcp.NewToolResultText(string(jsonAction)), nil
}

// powerOffDroplet powers off a droplet
func (d *DropletActionsTool) powerOffDroplet(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dropletID := req.GetArguments()["ID"].(float64)
	action, _, err := d.client.DropletActions.PowerOff(ctx, int(dropletID))
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonAction, err := json.MarshalIndent(action, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	return mcp.NewToolResultText(string(jsonAction)), nil
}

// shutdownDroplet shuts down a droplet
func (d *DropletActionsTool) shutdownDroplet(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dropletID := req.GetArguments()["ID"].(float64)
	action, _, err := d.client.DropletActions.Shutdown(ctx, int(dropletID))
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonAction, err := json.MarshalIndent(action, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	return mcp.NewToolResultText(string(jsonAction)), nil
}

// restoreDroplet restores a droplet to a backup image
func (d *DropletActionsTool) restoreDroplet(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dropletID := req.GetArguments()["ID"].(float64)
	imageID := req.GetArguments()["ImageID"].(float64)
	action, _, err := d.client.DropletActions.Restore(ctx, int(dropletID), int(imageID))
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonAction, err := json.MarshalIndent(action, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	return mcp.NewToolResultText(string(jsonAction)), nil
}

// resizeDroplet resizes a droplet
func (d *DropletActionsTool) resizeDroplet(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dropletID := req.GetArguments()["ID"].(float64)
	size := req.GetArguments()["Size"].(string)
	resizeDisk, _ := req.GetArguments()["ResizeDisk"].(bool) // Defaults to false
	action, _, err := d.client.DropletActions.Resize(ctx, int(dropletID), size, resizeDisk)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonAction, err := json.MarshalIndent(action, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	return mcp.NewToolResultText(string(jsonAction)), nil
}

// rebuildDroplet rebuilds a droplet using a provided image
func (d *DropletActionsTool) rebuildDroplet(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dropletID := req.GetArguments()["ID"].(float64)
	imageID := req.GetArguments()["ImageID"].(float64)
	action, _, err := d.client.DropletActions.RebuildByImageID(ctx, int(dropletID), int(imageID))
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonAction, err := json.MarshalIndent(action, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	return mcp.NewToolResultText(string(jsonAction)), nil
}

// renameDroplet renames a droplet
func (d *DropletActionsTool) renameDroplet(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dropletID := req.GetArguments()["ID"].(float64)
	name := req.GetArguments()["Name"].(string)
	action, _, err := d.client.DropletActions.Rename(ctx, int(dropletID), name)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonAction, err := json.MarshalIndent(action, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	return mcp.NewToolResultText(string(jsonAction)), nil
}

// changeKernel changes a droplet's kernel
func (d *DropletActionsTool) changeKernel(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dropletID := req.GetArguments()["ID"].(float64)
	kernelID := req.GetArguments()["KernelID"].(float64)
	action, _, err := d.client.DropletActions.ChangeKernel(ctx, int(dropletID), int(kernelID))
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonAction, err := json.MarshalIndent(action, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	return mcp.NewToolResultText(string(jsonAction)), nil
}

// enableIPv6 enables IPv6 on a droplet
func (d *DropletActionsTool) enableIPv6(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dropletID := req.GetArguments()["ID"].(float64)
	action, _, err := d.client.DropletActions.EnableIPv6(ctx, int(dropletID))
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonAction, err := json.MarshalIndent(action, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	return mcp.NewToolResultText(string(jsonAction)), nil
}

// enableBackups enables backups on a droplet
func (d *DropletActionsTool) enableBackups(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dropletID := req.GetArguments()["ID"].(float64)
	action, _, err := d.client.DropletActions.EnableBackups(ctx, int(dropletID))
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonAction, err := json.MarshalIndent(action, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	return mcp.NewToolResultText(string(jsonAction)), nil
}

// disableBackups disables backups on a droplet
func (d *DropletActionsTool) disableBackups(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dropletID := req.GetArguments()["ID"].(float64)
	action, _, err := d.client.DropletActions.DisableBackups(ctx, int(dropletID))
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonAction, err := json.MarshalIndent(action, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	return mcp.NewToolResultText(string(jsonAction)), nil
}

// snapshotDroplet creates a snapshot of a droplet
func (d *DropletActionsTool) snapshotDroplet(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dropletID := req.GetArguments()["ID"].(float64)
	name := req.GetArguments()["Name"].(string)
	action, _, err := d.client.DropletActions.Snapshot(ctx, int(dropletID), name)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("api error", err), nil
	}

	jsonAction, err := json.MarshalIndent(action, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	return mcp.NewToolResultText(string(jsonAction)), nil
}

// Tools returns a list of tool functions
func (d *DropletActionsTool) Tools() []server.ServerTool {
	tools := []server.ServerTool{
		{
			Handler: d.powerCycleDroplet,
			Tool: mcp.NewTool("digitalocean-droplet-action-power-cycle",
				mcp.WithDescription("Power cycle a droplet"),
				mcp.WithNumber("ID", mcp.Required(), mcp.Description("ID of the droplet to power cycle")),
			),
		},
		{
			Handler: d.powerOnDroplet,
			Tool: mcp.NewTool("digitalocean-droplet-action-power-on",
				mcp.WithDescription("Power on a droplet"),
				mcp.WithNumber("ID", mcp.Required(), mcp.Description("ID of the droplet to power on")),
			),
		},
		{
			Handler: d.powerOffDroplet,
			Tool: mcp.NewTool("digitalocean-droplet-action-power-off",
				mcp.WithDescription("Power off a droplet"),
				mcp.WithNumber("ID", mcp.Required(), mcp.Description("ID of the droplet to power off")),
			),
		},
		{
			Handler: d.shutdownDroplet,
			Tool: mcp.NewTool("digitalocean-droplet-action-shutdown",
				mcp.WithDescription("Shutdown a droplet"),
				mcp.WithNumber("ID", mcp.Required(), mcp.Description("ID of the droplet to shutdown")),
			),
		},
		{
			Handler: d.restoreDroplet,
			Tool: mcp.NewTool("digitalocean-droplet-action-restore",
				mcp.WithDescription("Restore a droplet from a backup/snapshot"),
				mcp.WithNumber("ID", mcp.Required(), mcp.Description("ID of the droplet to restore")),
				mcp.WithNumber("ImageID", mcp.Required(), mcp.Description("ID of the backup/snapshot image")),
			),
		},
		{
			Handler: d.resizeDroplet,
			Tool: mcp.NewTool("digitalocean-droplet-action-resize",
				mcp.WithDescription("Resize a droplet"),
				mcp.WithNumber("ID", mcp.Required(), mcp.Description("ID of the droplet to resize")),
				mcp.WithString("Size", mcp.Required(), mcp.Description("Slug of the new size (e.g., s-1vcpu-1gb)")),
				mcp.WithBoolean("ResizeDisk", mcp.DefaultBool(false), mcp.Description("Whether to resize the disk")),
			),
		},
		{
			Handler: d.rebuildDroplet,
			Tool: mcp.NewTool("digitalocean-droplet-action-rebuild",
				mcp.WithDescription("Rebuild a droplet from an image"),
				mcp.WithNumber("ID", mcp.Required(), mcp.Description("ID of the droplet to rebuild")),
				mcp.WithNumber("ImageID", mcp.Required(), mcp.Description("ID of the image to rebuild from")),
			),
		},
		{
			Handler: d.renameDroplet,
			Tool: mcp.NewTool("digitalocean-droplet-action-rename",
				mcp.WithDescription("Rename a droplet"),
				mcp.WithNumber("ID", mcp.Required(), mcp.Description("ID of the droplet to rename")),
				mcp.WithString("Name", mcp.Required(), mcp.Description("New name for the droplet")),
			),
		},
		{
			Handler: d.changeKernel,
			Tool: mcp.NewTool("digitalocean-droplet-action-change-kernel",
				mcp.WithDescription("Change a droplet's kernel"),
				mcp.WithNumber("ID", mcp.Required(), mcp.Description("ID of the droplet")),
				mcp.WithNumber("KernelID", mcp.Required(), mcp.Description("ID of the kernel to switch to")),
			),
		},
		{
			Handler: d.enableIPv6,
			Tool: mcp.NewTool("digitalocean-droplet-action-enable-ipv6",
				mcp.WithDescription("Enable IPv6 on a droplet"),
				mcp.WithNumber("ID", mcp.Required(), mcp.Description("ID of the droplet")),
			),
		},
		{
			Handler: d.enableBackups,
			Tool: mcp.NewTool("digitalocean-droplet-action-enable-backups",
				mcp.WithDescription("Enable backups on a droplet"),
				mcp.WithNumber("ID", mcp.Required(), mcp.Description("ID of the droplet")),
			),
		},
		{
			Handler: d.disableBackups,
			Tool: mcp.NewTool("digitalocean-droplet-action-disable-backups",
				mcp.WithDescription("Disable backups on a droplet"),
				mcp.WithNumber("ID", mcp.Required(), mcp.Description("ID of the droplet")),
			),
		},
		{
			Handler: d.snapshotDroplet,
			Tool: mcp.NewTool("digitalocean-droplet-action-snapshot",
				mcp.WithDescription("Take a snapshot of a droplet"),
				mcp.WithNumber("ID", mcp.Required(), mcp.Description("ID of the droplet")),
				mcp.WithString("Name", mcp.Required(), mcp.Description("Name for the snapshot")),
			),
		},
	}
	return tools
}

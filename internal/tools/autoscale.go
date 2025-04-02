package tools

import (
	"context"
	"encoding/json"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type AutoscaleTool struct {
	client *godo.Client
}

func NewAutoscaleTool(client *godo.Client) *AutoscaleTool {
	return &AutoscaleTool{
		client: client,
	}
}

func (a *AutoscaleTool) CreateAutoscalePool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := req.Params.Arguments["Name"].(string)
	config := req.Params.Arguments["Config"].(map[string]interface{})
	template := req.Params.Arguments["DropletTemplate"].(map[string]interface{})

	createRequest := &godo.DropletAutoscalePoolRequest{
		Name: name,
		Config: &godo.DropletAutoscaleConfiguration{
			MinInstances:            uint64(config["MinInstances"].(float64)),
			MaxInstances:            uint64(config["MaxInstances"].(float64)),
			TargetCPUUtilization:    config["TargetCPUUtilization"].(float64),
			TargetMemoryUtilization: config["TargetMemoryUtilization"].(float64),
			CooldownMinutes:         uint32(config["CooldownMinutes"].(float64)),
		},
		DropletTemplate: &godo.DropletAutoscaleResourceTemplate{
			Size:   template["Size"].(string),
			Image:  template["Image"].(string),
			Region: template["Region"].(string),
		},
	}

	pool, _, err := a.client.DropletAutoscale.Create(ctx, createRequest)
	if err != nil {
		return nil, err
	}

	jsonPool, err := json.MarshalIndent(pool, "", "  ")
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(jsonPool)), nil
}

func (a *AutoscaleTool) DeleteAutoscalePool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id := req.Params.Arguments["ID"].(string)

	_, err := a.client.DropletAutoscale.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText("Autoscale pool deleted successfully"), nil
}

func (a *AutoscaleTool) UpdateAutoscalePool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id := req.Params.Arguments["ID"].(string)
	config := req.Params.Arguments["Config"].(map[string]interface{})
	template := req.Params.Arguments["DropletTemplate"].(map[string]interface{})

	updateRequest := &godo.DropletAutoscalePoolRequest{
		Name: req.Params.Arguments["Name"].(string),
		Config: &godo.DropletAutoscaleConfiguration{
			MinInstances:            uint64(config["MinInstances"].(float64)),
			MaxInstances:            uint64(config["MaxInstances"].(float64)),
			TargetCPUUtilization:    config["TargetCPUUtilization"].(float64),
			TargetMemoryUtilization: config["TargetMemoryUtilization"].(float64),
			CooldownMinutes:         uint32(config["CooldownMinutes"].(float64)),
		},
		DropletTemplate: &godo.DropletAutoscaleResourceTemplate{
			Size:   template["Size"].(string),
			Image:  template["Image"].(string),
			Region: template["Region"].(string),
		},
	}

	pool, _, err := a.client.DropletAutoscale.Update(ctx, id, updateRequest)
	if err != nil {
		return nil, err
	}

	jsonPool, err := json.MarshalIndent(pool, "", "  ")
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(jsonPool)), nil
}

func (a *AutoscaleTool) Tools() []server.ServerTool {
	return []server.ServerTool{
		{
			Handler: a.CreateAutoscalePool,
			Tool: mcp.NewTool("digitalocean-autoscale-create",
				mcp.WithDescription("Create a new autoscale pool"),
				mcp.WithString("Name", mcp.Required(), mcp.Description("Name of the autoscale pool")),
				mcp.WithObject("Config", mcp.Required(), mcp.Description("Configuration for the autoscale pool")),
				mcp.WithObject("DropletTemplate", mcp.Required(), mcp.Description("Droplet template for the autoscale pool")),
			),
		},
		{
			Handler: a.DeleteAutoscalePool,
			Tool: mcp.NewTool("digitalocean-autoscale-delete",
				mcp.WithDescription("Delete an autoscale pool"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("ID of the autoscale pool to delete")),
			),
		},
		{
			Handler: a.UpdateAutoscalePool,
			Tool: mcp.NewTool("digitalocean-autoscale-update",
				mcp.WithDescription("Update an autoscale pool"),
				mcp.WithString("ID", mcp.Required(), mcp.Description("ID of the autoscale pool to update")),
				mcp.WithString("Name", mcp.Required(), mcp.Description("Name of the autoscale pool")),
				mcp.WithObject("Config", mcp.Required(), mcp.Description("Updated configuration for the autoscale pool")),
				mcp.WithObject("DropletTemplate", mcp.Required(), mcp.Description("Updated droplet template for the autoscale pool")),
			),
		},
	}
}

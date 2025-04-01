package registry

import (
	"mcp-digitalocean/internal/resources"
	"mcp-digitalocean/internal/tools"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterTools(s *server.MCPServer, c *godo.Client) {
	s.AddTools(tools.NewDropletTool(c).Tools()...)
	s.AddTools(tools.NewAppTool(c).Tools()...)
}

func RegisterResources(s *server.MCPServer, c *godo.Client) {
	// Register droplet resource and resource templates
	dropletResource := resources.NewDropletMCPResource(c)
	for template, handler := range dropletResource.ResourceTemplates() {
		s.AddResourceTemplate(template, handler)
	}

	// Register sizes resource
	sizesResource := resources.NewSizesMCPResource(c)
	for resource, handler := range sizesResource.Resources() {
		s.AddResource(resource, handler)
	}

	// Register images resources and templates
	imageResource := resources.NewImagesMCPResource(c)
	for resource, handler := range imageResource.Resources() {
		s.AddResource(resource, handler)
	}
	for template, handler := range imageResource.ResourceTemplates() {
		s.AddResourceTemplate(template, handler)
	}

	// Register account resource
	for resource, handler := range resources.NewAccountMCPResource(c).Resources() {
		s.AddResource(resource, handler)
	}

	// Register action resource
	for template, handler := range resources.NewActionMCPResource(c).Resources() {
		s.AddResourceTemplate(template, handler)
	}

	// Register apps resource
	for template, handler := range resources.NewAppMCPResource(c).Resources() {
		s.AddResourceTemplate(template, handler)
	}

	// Register balance resource
	for resource, handler := range resources.NewBalanceMCPResource(c).Resources() {
		s.AddResource(resource, handler)
	}
}

package registry

import (
	"mcp-digitalocean/internal/resources"
	"mcp-digitalocean/internal/tools"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterTools(s *server.MCPServer, c *godo.Client) {
	// Register droplet tools
	s.AddTools(tools.NewDropletTool(c).Tools()...)
}

func RegisterResources(s *server.MCPServer, c *godo.Client) {
	// Register droplet resource
	for template, handler := range resources.NewDropletMCPResource(c).Resources() {
		s.AddResourceTemplate(template, handler)
	}

	// Register sizes resource
	for template, handler := range resources.NewSizesMCPResource(c).Resources() {
		s.AddResourceTemplate(template, handler)
	}
}

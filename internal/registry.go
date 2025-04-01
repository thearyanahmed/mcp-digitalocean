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
	// Register droplet resource
	for template, handler := range resources.NewDropletMCPResource(c).Resources() {
		s.AddResourceTemplate(template, handler)
	}

	// Register sizes resource
	for template, handler := range resources.NewSizesMCPResource(c).Resources() {
		s.AddResourceTemplate(template, handler)
	}

	// Register account resource
	for template, handler := range resources.NewAccountMCPResource(c).Resources() {
		s.AddResourceTemplate(template, handler)
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
	for template, handler := range resources.NewBalanceMCPResource(c).Resources() {
		s.AddResourceTemplate(template, handler)
	}

}

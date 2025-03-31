// filepath: /Users/ashar/Developer/mcp-digitalocean/internal/resources/sizes.go
package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
)

// SizesMCPResource represents a handler for MCP Size resources
type SizesMCPResource struct {
	client *godo.Client
}

// NewSizesMCPResource creates a new Sizes MCP resource handler
func NewSizesMCPResource(client *godo.Client) *SizesMCPResource {
	return &SizesMCPResource{
		client: client,
	}
}

// GetResourceTemplate returns the template for the Sizes MCP resource
func (s *SizesMCPResource) GetResourceTemplate() mcp.ResourceTemplate {
	return mcp.NewResourceTemplate(
		"sizes://all",
		"Droplet Sizes",
		mcp.WithTemplateDescription("Returns all available droplet sizes"),
		mcp.WithTemplateMIMEType("application/json"),
	)
}

// GetSizeResourceTemplate returns the template for a specific Size resource
func (s *SizesMCPResource) GetSizeResourceTemplate() mcp.ResourceTemplate {
	return mcp.NewResourceTemplate(
		"sizes://{slug}",
		"Droplet Size",
		mcp.WithTemplateDescription("Returns information about a specific droplet size"),
		mcp.WithTemplateMIMEType("application/json"),
	)
}

// HandleGetResource handles the Sizes MCP resource requests for all sizes
func (s *SizesMCPResource) HandleGetResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// List all droplet sizes from DigitalOcean API
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200, // Get a large number of sizes at once
	}

	sizes, _, err := s.client.Sizes.List(ctx, opt)
	if err != nil {
		return nil, fmt.Errorf("error fetching droplet sizes: %s", err)
	}

	// Serialize to JSON
	jsonData, err := json.MarshalIndent(sizes, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error serializing sizes: %s", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(jsonData),
		},
	}, nil
}

// HandleGetSizeResource handles requests for specific droplet size
func (s *SizesMCPResource) HandleGetSizeResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// Extract size slug from the URI
	sizeSlug, err := extractSizeSlugFromURI(request.Params.URI)
	if err != nil {
		return nil, fmt.Errorf("invalid size URI: %s", err)
	}

	// List all sizes and find the one with matching slug
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200, // Get a large number of sizes at once
	}

	sizes, _, err := s.client.Sizes.List(ctx, opt)
	if err != nil {
		return nil, fmt.Errorf("error fetching droplet sizes: %s", err)
	}

	var targetSize *godo.Size
	for _, size := range sizes {
		if size.Slug == sizeSlug {
			targetSize = &size
			break
		}
	}

	if targetSize == nil {
		return nil, fmt.Errorf("size with slug '%s' not found", sizeSlug)
	}

	// Serialize to JSON
	jsonData, err := json.MarshalIndent(targetSize, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error serializing size: %s", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(jsonData),
		},
	}, nil
}

// extractSizeSlugFromURI extracts the size slug from the URI
func extractSizeSlugFromURI(uri string) (string, error) {
	// Use regex to extract the slug from the URI format "sizes://{slug}"
	re := regexp.MustCompile(`sizes://([a-zA-Z0-9\-]+)`)
	match := re.FindStringSubmatch(uri)
	if len(match) < 2 {
		return "", fmt.Errorf("could not extract size slug from URI: %s", uri)
	}

	return match[1], nil
}

// Resources returns the available resources for the Sizes MCP resource
func (s *SizesMCPResource) Resources() map[mcp.ResourceTemplate]MCPResourceHandler {
	return map[mcp.ResourceTemplate]MCPResourceHandler{
		s.GetResourceTemplate():     s.HandleGetResource,
		s.GetSizeResourceTemplate(): s.HandleGetSizeResource,
	}
}

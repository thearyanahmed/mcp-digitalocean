package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
)

type AutoscaleMCPResource struct {
	client *godo.Client
}

func NewAutoscaleMCPResource(client *godo.Client) *AutoscaleMCPResource {
	return &AutoscaleMCPResource{
		client: client,
	}
}

func (a *AutoscaleMCPResource) GetResourceTemplate() mcp.ResourceTemplate {
	return mcp.NewResourceTemplate(
		"autoscale://{id}",
		"Autoscale Pool",
		mcp.WithTemplateDescription("Returns autoscale pool information"),
		mcp.WithTemplateMIMEType("application/json"),
	)
}

func (a *AutoscaleMCPResource) HandleGetResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	id, err := extractAutoscaleIDFromURI(request.Params.URI)
	if err != nil {
		return nil, fmt.Errorf("invalid autoscale URI: %s", err)
	}

	pool, _, err := a.client.DropletAutoscale.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching autoscale pool: %s", err)
	}

	jsonData, err := json.MarshalIndent(pool, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error serializing autoscale pool: %s", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(jsonData),
		},
	}, nil
}

func extractAutoscaleIDFromURI(uri string) (string, error) {
	re := regexp.MustCompile(`autoscale://([^/]+)`)
	match := re.FindStringSubmatch(uri)
	if len(match) < 2 {
		return "", fmt.Errorf("could not extract autoscale ID from URI: %s", uri)
	}
	return match[1], nil
}

func (a *AutoscaleMCPResource) ResourceTemplates() map[mcp.ResourceTemplate]MCPResourceHandler {
	return map[mcp.ResourceTemplate]MCPResourceHandler{
		a.GetResourceTemplate(): a.HandleGetResource,
	}
}

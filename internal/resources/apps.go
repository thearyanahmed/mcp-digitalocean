package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
)

// MCPResourceHandler = func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error)

// AppMCPResource represents a handler for MCP App resources
type AppMCPResource struct {
	client *godo.Client
}

// NewAppMCPResource creates a new App MCP resource handler
func NewAppMCPResource(client *godo.Client) *AppMCPResource {
	return &AppMCPResource{
		client: client,
	}
}

// GetResourceTemplate returns the template for the App MCP resource
func (a *AppMCPResource) GetResourceTemplate() mcp.ResourceTemplate {
	return mcp.NewResourceTemplate(
		"apps://{id}",
		"App",
		mcp.WithTemplateDescription("Returns app information"),
		mcp.WithTemplateMIMEType("application/json"),
	)
}

// HandleGetResource handles the App MCP resource requests
func (a *AppMCPResource) HandleGetResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// Extract app ID from the URI
	appID, err := extractAppIDFromURI(request.Params.URI)
	if err != nil {
		return nil, fmt.Errorf("invalid app URI: %s", err)
	}

	// Get app from DigitalOcean API
	app, _, err := a.client.Apps.Get(ctx, appID)
	if err != nil {
		return nil, fmt.Errorf("error fetching app: %s", err)
	}

	// Serialize to JSON
	jsonData, err := json.MarshalIndent(app, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error serializing app: %s", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(jsonData),
		},
	}, nil
}

// extractAppIDFromURI extracts the app ID from the URI
func extractAppIDFromURI(uri string) (string, error) {
	// Use regex to extract the ID from the URI format "apps://{id}"
	re := regexp.MustCompile(`apps://(.+)`)
	match := re.FindStringSubmatch(uri)
	if len(match) < 2 {
		return "", fmt.Errorf("could not extract app ID from URI: %s", uri)
	}

	return match[1], nil
}

// GetDeploymentTemplate returns the template for the App Deployment resource
func (a *AppMCPResource) GetDeploymentTemplate() mcp.ResourceTemplate {
	return mcp.NewResourceTemplate(
		"apps://{id}/deployments/{deployment_id}",
		"App Deployment",
		mcp.WithTemplateDescription("Returns deployment information for an app"),
		mcp.WithTemplateMIMEType("application/json"),
	)
}

// HandleGetDeployment handles the App Deployment resource requests
func (a *AppMCPResource) HandleGetDeployment(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// Extract app ID and deployment ID from the URI
	uriParts := regexp.MustCompile(`apps://(.+)/deployments/(.+)`).FindStringSubmatch(request.Params.URI)
	if len(uriParts) < 3 {
		return nil, fmt.Errorf("invalid deployment URI: %s", request.Params.URI)
	}
	appID, deploymentID := uriParts[1], uriParts[2]

	// Get deployment from DigitalOcean API
	deployment, _, err := a.client.Apps.GetDeployment(ctx, appID, deploymentID)
	if err != nil {
		return nil, fmt.Errorf("error fetching deployment: %s", err)
	}

	// Serialize to JSON
	jsonData, err := json.MarshalIndent(deployment, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error serializing deployment: %s", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(jsonData),
		},
	}, nil
}

// GetTierTemplate returns the template for the App Tier resource
func (a *AppMCPResource) GetTierTemplate() mcp.ResourceTemplate {
	return mcp.NewResourceTemplate(
		"apps://{id}/tier",
		"App Tier",
		mcp.WithTemplateDescription("Returns tier information for an app"),
		mcp.WithTemplateMIMEType("application/json"),
	)
}

// HandleGetTier handles the App Tier resource requests
func (a *AppMCPResource) HandleGetTier(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// Extract app ID from the URI
	appID, err := extractAppIDFromURI(request.Params.URI)
	if err != nil {
		return nil, fmt.Errorf("invalid tier URI: %s", err)
	}

	// Get tier information from DigitalOcean API
	tierInfo, _, err := a.client.Apps.GetTier(ctx, appID)
	if err != nil {
		return nil, fmt.Errorf("error fetching tier information: %s", err)
	}

	// Serialize to JSON
	jsonData, err := json.MarshalIndent(tierInfo, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error serializing tier information: %s", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(jsonData),
		},
	}, nil
}

// Resources returns the available resources for the App MCP resource
func (a *AppMCPResource) Resources() map[mcp.ResourceTemplate]MCPResourceHandler {
	return map[mcp.ResourceTemplate]MCPResourceHandler{
		a.GetResourceTemplate():        a.HandleGetResource,
		a.GetDeploymentTemplate():      a.HandleGetDeployment,
		a.GetTierTemplate():            a.HandleGetTier,
	}
}

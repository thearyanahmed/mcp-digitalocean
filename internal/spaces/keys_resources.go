package spaces

import (
	"context"
	"encoding/json"
	"fmt"

	"mcp-digitalocean/internal/common"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const KeysURI = "spaces_keys://"

type KeysIPMCPResource struct {
	client *godo.Client
}

func NewKeysIPMCPResource(client *godo.Client) *KeysIPMCPResource {
	return &KeysIPMCPResource{
		client: client,
	}
}

func (r *KeysIPMCPResource) getKeysResourceTemplate() mcp.ResourceTemplate {
	return mcp.NewResourceTemplate(
		KeysURI+"{access_key}",
		"Spaces Keys",
		mcp.WithTemplateDescription("Returns Spaces key information"),
		mcp.WithTemplateMIMEType("application/json"),
	)
}

func (r *KeysIPMCPResource) getKeysListResource() mcp.Resource {
	return mcp.NewResource(
		KeysURI+"all",
		"Spaces Keys List",
		mcp.WithResourceDescription("Returns list of all Spaces keys"),
		mcp.WithMIMEType("application/json"),
	)
}

func (r *KeysIPMCPResource) handleGetKeysResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	keyID, err := common.ExtractStringIDFromURI(request.Params.URI)
	if err != nil {
		return nil, fmt.Errorf("invalid Spaces key URI: %w", err)
	}

	if keyID == "" {
		return nil, fmt.Errorf("AccessKey cannot be empty")
	}

	key, _, err := r.client.SpacesKeys.Get(ctx, keyID)
	if err != nil {
		return nil, fmt.Errorf("error fetching Spaces key: %w", err)
	}

	jsonData, err := json.MarshalIndent(key, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling Spaces key: %w", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(jsonData),
		},
	}, nil
}

func (r *KeysIPMCPResource) handleGetKeysListResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	keys, _, err := r.client.SpacesKeys.List(ctx, &godo.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error fetching Spaces keys: %w", err)
	}

	jsonData, err := json.MarshalIndent(keys, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling Spaces keys: %w", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(jsonData),
		},
	}, nil
}

func (r *KeysIPMCPResource) Resources() map[mcp.Resource]server.ResourceHandlerFunc {
	return map[mcp.Resource]server.ResourceHandlerFunc{
		r.getKeysListResource(): r.handleGetKeysListResource,
	}
}

func (r *KeysIPMCPResource) ResourceTemplates() map[mcp.ResourceTemplate]server.ResourceTemplateHandlerFunc {
	return map[mcp.ResourceTemplate]server.ResourceTemplateHandlerFunc{
		r.getKeysResourceTemplate(): r.handleGetKeysResource,
	}
}

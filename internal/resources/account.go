package resources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
)

type AccountMCPResource struct {
	client *godo.Client
}

func NewAccountMCPResource(client *godo.Client) *AccountMCPResource {
	return &AccountMCPResource{
		client: client,
	}
}

func (a *AccountMCPResource) GetResourceTemplate() mcp.ResourceTemplate {
	return mcp.NewResourceTemplate(
		"account://current",
		"Account",
		mcp.WithTemplateDescription("Returns account information"),
		mcp.WithTemplateMIMEType("application/json"),
	)
}

func (a *AccountMCPResource) HandleGetResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// Get account information from DigitalOcean API
	account, _, err := a.client.Account.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching account: %s", err)
	}

	// Serialize to JSON
	jsonData, err := json.MarshalIndent(account, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error serializing account: %s", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			MIMEType: "application/json",
			Text:     string(jsonData),
		},
	}, nil
}

func (a *AccountMCPResource) Resources() map[mcp.ResourceTemplate]MCPResourceHandler {
	return map[mcp.ResourceTemplate]MCPResourceHandler{
		a.GetResourceTemplate(): a.HandleGetResource,
	}
}

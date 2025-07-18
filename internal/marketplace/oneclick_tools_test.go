package marketplace

import (
	"testing"

	"github.com/digitalocean/godo"
	"github.com/stretchr/testify/assert"
)

func TestNewOneClickTool(t *testing.T) {
	client := &godo.Client{}
	tool := NewOneClickTool(client)

	assert.NotNil(t, tool)
	assert.Equal(t, client, tool.client)
}

func TestOneClickTool_Tools(t *testing.T) {
	client := &godo.Client{}
	tool := NewOneClickTool(client)

	tools := tool.Tools()
	assert.Len(t, tools, 2)

	// Check tool names
	toolNames := make([]string, len(tools))
	for i, tool := range tools {
		toolNames[i] = tool.Tool.Name
	}

	assert.Contains(t, toolNames, "digitalocean-oneclick-list")
	assert.Contains(t, toolNames, "digitalocean-oneclick-install-kubernetes")
}

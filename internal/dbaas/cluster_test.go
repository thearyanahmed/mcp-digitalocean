package dbaas

import (
	"context"
	"mcp-digitalocean/internal/dbaas/mocks"
	"testing"
	"time"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func getText(res *mcp.CallToolResult) string {
	if len(res.Content) > 0 {
		if tc, ok := res.Content[0].(mcp.TextContent); ok {
			return tc.Text
		}
	}
	return ""
}

func TestClusterTool_listCluster(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("List", mock.Anything, (*godo.ListOptions)(nil)).Return([]godo.Database{{Name: "test-db"}}, nil, nil)

	client := &godo.Client{}
	client.Databases = mockDB
	ct := &ClusterTool{client: client}

	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{}}}
	res, err := ct.listCluster(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, getText(res), "test-db")
	mockDB.AssertExpectations(t)
}

func TestClusterTool_getCluster(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	cluster := &godo.Database{Name: "my-cluster"}
	mockDB.On("Get", mock.Anything, "abc").Return(cluster, nil, nil)

	client := &godo.Client{}
	client.Databases = mockDB
	ct := &ClusterTool{client: client}

	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{"ID": "abc"}}}
	res, err := ct.getCluster(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, getText(res), "my-cluster")
	mockDB.AssertExpectations(t)

	// Error case: missing ID
	reqMissing := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{}}}
	res, err = ct.getCluster(context.Background(), reqMissing)
	assert.NoError(t, err)
	assert.Contains(t, getText(res), "Cluster ID is required")
}

func TestClusterTool_createCluster(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	created := &godo.Database{Name: "new-cluster"}
	mockDB.On("Create", mock.Anything, mock.AnythingOfType("*godo.DatabaseCreateRequest")).Return(created, nil, nil)

	client := &godo.Client{}
	client.Databases = mockDB
	ct := &ClusterTool{client: client}

	args := map[string]interface{}{
		"name":      "new-cluster",
		"engine":    "pg",
		"version":   "13",
		"region":    "nyc1",
		"size":      "db-s-1vcpu-1gb",
		"num_nodes": float64(2),
	}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := ct.createCluster(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, getText(res), "new-cluster")
	mockDB.AssertExpectations(t)

	// Error case: API error
	mockDB2 := &mocks.DatabasesService{}
	mockDB2.On("Create", mock.Anything, mock.Anything).Return(nil, nil, assert.AnError)
	client.Databases = mockDB2
	ct = &ClusterTool{client: client}
	res, err = ct.createCluster(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, getText(res), "api error")
}

func TestClusterTool_deleteCluster(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("Delete", mock.Anything, "abc").Return(nil, nil)

	client := &godo.Client{}
	client.Databases = mockDB
	ct := &ClusterTool{client: client}

	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{"ID": "abc"}}}
	res, err := ct.deleteCluster(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, getText(res), "Cluster deleted successfully")
	mockDB.AssertExpectations(t)

	// Error case: missing ID
	reqMissing := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{}}}
	res, err = ct.deleteCluster(context.Background(), reqMissing)
	assert.NoError(t, err)
	assert.Contains(t, getText(res), "Cluster ID is required")
}

func TestClusterTool_resizeCluster(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("Resize", mock.Anything, "abc", mock.AnythingOfType("*godo.DatabaseResizeRequest")).Return(nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	ct := &ClusterTool{client: client}
	args := map[string]interface{}{"ID": "abc", "size": "db-s-2vcpu-4gb", "num_nodes": float64(3)}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := ct.resizeCluster(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, getText(res), "Cluster resize initiated successfully")
	mockDB.AssertExpectations(t)
	// Error case: missing ID
	reqMissing := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{}}}
	res, err = ct.resizeCluster(context.Background(), reqMissing)
	assert.NoError(t, err)
	assert.Contains(t, getText(res), "Cluster ID is required")
}

func TestClusterTool_getCA(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	ca := &godo.DatabaseCA{Certificate: []byte("cert-data")}
	mockDB.On("GetCA", mock.Anything, "abc").Return(ca, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	ct := &ClusterTool{client: client}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{"ID": "abc"}}}
	res, err := ct.getCA(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, getText(res), "Y2VydC1kYXRh")
	mockDB.AssertExpectations(t)
	// Error case: missing ID
	reqMissing := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{}}}
	res, err = ct.getCA(context.Background(), reqMissing)
	assert.NoError(t, err)
	assert.Contains(t, getText(res), "Cluster ID is required")
}

func TestClusterTool_listBackups(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("ListBackups", mock.Anything, "abc", mock.Anything).Return([]godo.DatabaseBackup{{CreatedAt: time.Now()}}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	ct := &ClusterTool{client: client}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{"ID": "abc"}}}
	res, err := ct.listBackups(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, getText(res), "created_at")
	mockDB.AssertExpectations(t)
	// Error case: missing ID
	reqMissing := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{}}}
	res, err = ct.listBackups(context.Background(), reqMissing)
	assert.NoError(t, err)
	assert.Contains(t, getText(res), "Cluster ID is required")
}

func TestClusterTool_listOptions(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("ListOptions", mock.Anything).Return(&godo.DatabaseOptions{}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	ct := &ClusterTool{client: client}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{}}}
	res, err := ct.listOptions(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, getText(res), "pg")
	mockDB.AssertExpectations(t)
}

func TestClusterTool_upgradeMajorVersion(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("UpgradeMajorVersion", mock.Anything, "abc", mock.AnythingOfType("*godo.UpgradeVersionRequest")).Return(nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	ct := &ClusterTool{client: client}
	args := map[string]interface{}{"ID": "abc", "version": "15"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := ct.upgradeMajorVersion(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, getText(res), "Major version upgrade initiated successfully")
	mockDB.AssertExpectations(t)
	// Error case: missing version
	args = map[string]interface{}{"ID": "abc"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = ct.upgradeMajorVersion(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, getText(res), "Target version is required")
}

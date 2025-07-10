package dbaas

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"mcp-digitalocean/internal/dbaas/mocks"

	"github.com/mark3labs/mcp-go/mcp"
)

func newUserToolWithMock() (*UserTool, *mocks.DatabasesService) {
	mockSvc := new(mocks.DatabasesService)
	return &UserTool{client: &godo.Client{Databases: mockSvc}}, mockSvc
}

func getTextContent(res *mcp.CallToolResult) string {
	if res == nil || len(res.Content) == 0 {
		return ""
	}
	if tc, ok := res.Content[0].(mcp.TextContent); ok {
		return tc.Text
	}
	return ""
}

func TestUserTool_getUser(t *testing.T) {
	tool, mockSvc := newUserToolWithMock()
	ctx := context.Background()

	dbUser := &godo.DatabaseUser{Name: "testuser"}
	mockSvc.On("GetUser", ctx, "cid", "testuser").Return(dbUser, nil, nil)

	// Success
	res, err := tool.getUser(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"ID": "cid", "user": "testuser"}}})
	assert.NoError(t, err)
	assert.Contains(t, getTextContent(res), "testuser")

	// Missing Cluster ID
	res, err = tool.getUser(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"user": "testuser"}}})
	assert.NoError(t, err)
	assert.Equal(t, "Cluster ID is required", getTextContent(res))

	// Missing user
	res, err = tool.getUser(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"ID": "cid"}}})
	assert.NoError(t, err)
	assert.Equal(t, "User name is required", getTextContent(res))

	// API error
	errApi := errors.New("api fail")
	mockSvc.On("GetUser", ctx, "cid", "failuser").Return(nil, nil, errApi)
	res, err = tool.getUser(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"ID": "cid", "user": "failuser"}}})
	assert.NoError(t, err)
	assert.Contains(t, getTextContent(res), "api error")
}

func TestUserTool_listUsers(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		tool, mockSvc := newUserToolWithMock()
		users := []godo.DatabaseUser{{Name: "u1"}, {Name: "u2"}}
		mockSvc.On("ListUsers", ctx, "cid", (*godo.ListOptions)(nil)).Return(users, nil, nil)
		res, err := tool.listUsers(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"ID": "cid"}}})
		assert.NoError(t, err)
		assert.Contains(t, getTextContent(res), "u1")
		assert.Contains(t, getTextContent(res), "u2")
	})

	t.Run("pagination", func(t *testing.T) {
		tool, mockSvc := newUserToolWithMock()
		users := []godo.DatabaseUser{{Name: "u1"}, {Name: "u2"}}
		opts := &godo.ListOptions{Page: 2, PerPage: 5}
		mockSvc.On("ListUsers", ctx, "cid", opts).Return(users, nil, nil)
		res, err := tool.listUsers(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"ID": "cid", "page": "2", "per_page": "5"}}})
		assert.NoError(t, err)
		assert.Contains(t, getTextContent(res), "u1")
	})

	t.Run("missing cluster ID", func(t *testing.T) {
		tool, _ := newUserToolWithMock()
		res, err := tool.listUsers(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{}}})
		assert.NoError(t, err)
		assert.Equal(t, "Cluster ID is required", getTextContent(res))
	})

	t.Run("api error", func(t *testing.T) {
		tool, mockSvc := newUserToolWithMock()
		errApi := errors.New("api fail")
		mockSvc.On("ListUsers", ctx, "cid", (*godo.ListOptions)(nil)).Return(nil, nil, errApi)
		res, err := tool.listUsers(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"ID": "cid"}}})
		assert.NoError(t, err)
		assert.Contains(t, getTextContent(res), "api error")
	})
}

func TestUserTool_createUser(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		tool, mockSvc := newUserToolWithMock()
		dbUser := &godo.DatabaseUser{Name: "newuser"}
		mockSvc.On("CreateUser", ctx, "cid", mock.AnythingOfType("*godo.DatabaseCreateUserRequest")).Return(dbUser, nil, nil)
		res, err := tool.createUser(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"ID": "cid", "name": "newuser"}}})
		assert.NoError(t, err)
		assert.Contains(t, getTextContent(res), "newuser")
	})

	t.Run("with mysql_auth_plugin", func(t *testing.T) {
		tool, mockSvc := newUserToolWithMock()
		dbUser := &godo.DatabaseUser{Name: "newuser"}
		mockSvc.On("CreateUser", ctx, "cid", mock.MatchedBy(func(req *godo.DatabaseCreateUserRequest) bool {
			return req.MySQLSettings != nil && req.MySQLSettings.AuthPlugin == "mysql_native_password"
		})).Return(dbUser, nil, nil)
		res, err := tool.createUser(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"ID": "cid", "name": "newuser", "mysql_auth_plugin": "mysql_native_password"}}})
		assert.NoError(t, err)
		assert.Contains(t, getTextContent(res), "newuser")
	})

	t.Run("with settings_json", func(t *testing.T) {
		tool, mockSvc := newUserToolWithMock()
		dbUser := &godo.DatabaseUser{Name: "newuser"}
		settings := godo.DatabaseUserSettings{ACL: []*godo.KafkaACL{{ID: "acl1", Permission: "read"}}}
		settingsBytes, _ := json.Marshal(settings)
		mockSvc.On("CreateUser", ctx, "cid", mock.MatchedBy(func(req *godo.DatabaseCreateUserRequest) bool {
			return req.Settings != nil && len(req.Settings.ACL) > 0 && req.Settings.ACL[0].ID == "acl1"
		})).Return(dbUser, nil, nil)
		res, err := tool.createUser(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"ID": "cid", "name": "newuser", "settings_json": string(settingsBytes)}}})
		assert.NoError(t, err)
		assert.Contains(t, getTextContent(res), "newuser")
	})

	t.Run("invalid settings_json", func(t *testing.T) {
		tool, _ := newUserToolWithMock()
		res, err := tool.createUser(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"ID": "cid", "name": "newuser", "settings_json": "notjson"}}})
		assert.NoError(t, err)
		assert.Contains(t, getTextContent(res), "Invalid settings_json")
	})

	t.Run("missing cluster ID", func(t *testing.T) {
		tool, _ := newUserToolWithMock()
		res, err := tool.createUser(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"name": "newuser"}}})
		assert.NoError(t, err)
		assert.Equal(t, "Cluster ID is required", getTextContent(res))
	})

	t.Run("missing name", func(t *testing.T) {
		tool, _ := newUserToolWithMock()
		res, err := tool.createUser(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"ID": "cid"}}})
		assert.NoError(t, err)
		assert.Equal(t, "User name is required", getTextContent(res))
	})

	t.Run("api error", func(t *testing.T) {
		tool, mockSvc := newUserToolWithMock()
		errApi := errors.New("api fail")
		mockSvc.On("CreateUser", ctx, "cid", mock.Anything).Return(nil, nil, errApi)
		res, err := tool.createUser(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"ID": "cid", "name": "failuser"}}})
		assert.NoError(t, err)
		assert.Contains(t, getTextContent(res), "api error")
	})
}

func TestUserTool_updateUser(t *testing.T) {
	tool, mockSvc := newUserToolWithMock()
	ctx := context.Background()
	dbUser := &godo.DatabaseUser{Name: "updateduser"}
	mockSvc.On("UpdateUser", ctx, "cid", "updateduser", mock.AnythingOfType("*godo.DatabaseUpdateUserRequest")).Return(dbUser, nil, nil)

	// Success
	res, err := tool.updateUser(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"ID": "cid", "user": "updateduser"}}})
	assert.NoError(t, err)
	assert.Contains(t, getTextContent(res), "updateduser")

	// With settings_json (using ACL field)
	settings := godo.DatabaseUserSettings{ACL: []*godo.KafkaACL{{ID: "acl2", Permission: "write"}}}
	settingsBytes, _ := json.Marshal(settings)
	mockSvc.On("UpdateUser", ctx, "cid", "updateduser", mock.MatchedBy(func(req *godo.DatabaseUpdateUserRequest) bool {
		return req.Settings != nil && len(req.Settings.ACL) > 0 && req.Settings.ACL[0].ID == "acl2"
	})).Return(dbUser, nil, nil)
	res, err = tool.updateUser(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"ID": "cid", "user": "updateduser", "settings_json": string(settingsBytes)}}})
	assert.NoError(t, err)
	assert.Contains(t, getTextContent(res), "updateduser")

	// Invalid settings_json
	res, err = tool.updateUser(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"ID": "cid", "user": "updateduser", "settings_json": "notjson"}}})
	assert.NoError(t, err)
	assert.Contains(t, getTextContent(res), "Invalid settings_json")

	// Missing Cluster ID
	res, err = tool.updateUser(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"user": "updateduser"}}})
	assert.NoError(t, err)
	assert.Equal(t, "Cluster ID is required", getTextContent(res))

	// Missing user
	res, err = tool.updateUser(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"ID": "cid"}}})
	assert.NoError(t, err)
	assert.Equal(t, "User name is required", getTextContent(res))

	// API error
	errApi := errors.New("api fail")
	mockSvc.On("UpdateUser", ctx, "cid", "failuser", mock.Anything).Return(nil, nil, errApi)
	res, err = tool.updateUser(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"ID": "cid", "user": "failuser"}}})
	assert.NoError(t, err)
	assert.Contains(t, getTextContent(res), "api error")
}

func TestUserTool_deleteUser(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		tool, mockSvc := newUserToolWithMock()
		mockSvc.On("DeleteUser", ctx, "cid", "deluser").Return((*godo.Response)(nil), nil)
		res, err := tool.deleteUser(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"ID": "cid", "user": "deluser"}}})
		assert.NoError(t, err)
		assert.Equal(t, "User deleted successfully", getTextContent(res))
	})

	t.Run("missing cluster ID", func(t *testing.T) {
		tool, _ := newUserToolWithMock()
		res, err := tool.deleteUser(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"user": "deluser"}}})
		assert.NoError(t, err)
		assert.Equal(t, "Cluster ID is required", getTextContent(res))
	})

	t.Run("missing user", func(t *testing.T) {
		tool, _ := newUserToolWithMock()
		res, err := tool.deleteUser(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"ID": "cid"}}})
		assert.NoError(t, err)
		assert.Equal(t, "User name is required", getTextContent(res))
	})

	t.Run("api error", func(t *testing.T) {
		tool, mockSvc := newUserToolWithMock()
		errApi := errors.New("api fail")
		mockSvc.On("DeleteUser", ctx, "cid", "failuser").Return((*godo.Response)(nil), errApi)
		res, err := tool.deleteUser(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"ID": "cid", "user": "failuser"}}})
		assert.NoError(t, err)
		assert.Contains(t, getTextContent(res), "api error")
	})
}

package spaces

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func setupSpacesKeysToolWithMock(spacesKeys *MockSpacesKeysService) *KeysTool {
	client := &godo.Client{}
	client.SpacesKeys = spacesKeys
	return NewSpacesKeysTool(client)
}

func TestSpacesKeysTool_createSpacesKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testKey := &godo.SpacesKey{
		Name:      "test-key",
		AccessKey: "AKIA123456789",
		Grants: []*godo.Grant{
			{
				Bucket:     "",
				Permission: godo.SpacesKeyFullAccess,
			},
		},
	}

	tests := []struct {
		name        string
		args        map[string]any
		mockSetup   func(*MockSpacesKeysService)
		expectError bool
	}{
		{
			name: "Successful create",
			args: map[string]any{
				"Name": "test-key",
			},
			mockSetup: func(m *MockSpacesKeysService) {
				m.EXPECT().
					Create(gomock.Any(), &godo.SpacesKeyCreateRequest{
						Name: "test-key",
						Grants: []*godo.Grant{
							{
								Bucket:     "",
								Permission: godo.SpacesKeyFullAccess,
							},
						},
					}).
					Return(testKey, nil, nil).
					Times(1)
			},
		},
		{
			name:        "Missing Name parameter",
			args:        map[string]any{},
			expectError: true,
		},
		{
			name: "Invalid Name type",
			args: map[string]any{
				"Name": 123,
			},
			expectError: true,
		},
		{
			name: "Empty Name",
			args: map[string]any{
				"Name": "",
			},
			expectError: true,
		},
		{
			name: "API error",
			args: map[string]any{
				"Name": "fail-key",
			},
			mockSetup: func(m *MockSpacesKeysService) {
				m.EXPECT().
					Create(gomock.Any(), &godo.SpacesKeyCreateRequest{
						Name: "fail-key",
						Grants: []*godo.Grant{
							{
								Bucket:     "",
								Permission: godo.SpacesKeyFullAccess,
							},
						},
					}).
					Return(nil, nil, errors.New("api error")).
					Times(1)
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSpacesKeys := NewMockSpacesKeysService(ctrl)
			if tc.mockSetup != nil {
				tc.mockSetup(mockSpacesKeys)
			}
			tool := setupSpacesKeysToolWithMock(mockSpacesKeys)
			req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: tc.args}}
			resp, err := tool.createSpacesKey(context.Background(), req)
			if tc.expectError {
				require.NotNil(t, resp)
				require.True(t, resp.IsError)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.False(t, resp.IsError)
			var outKey godo.SpacesKey
			require.NoError(t, json.Unmarshal([]byte(resp.Content[0].(mcp.TextContent).Text), &outKey))
			require.Equal(t, testKey.AccessKey, outKey.AccessKey)
			require.Equal(t, testKey.Name, outKey.Name)
		})
	}
}

func TestSpacesKeysTool_updateSpacesKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testKey := &godo.SpacesKey{
		Name:      "updated-key",
		AccessKey: "AKIA123456789",
	}

	tests := []struct {
		name        string
		args        map[string]any
		mockSetup   func(*MockSpacesKeysService)
		expectError bool
	}{
		{
			name: "Successful update",
			args: map[string]any{
				"AccessKey": "AKIA123456789",
				"Name":      "updated-key",
			},
			mockSetup: func(m *MockSpacesKeysService) {
				m.EXPECT().
					Update(gomock.Any(), "AKIA123456789", &godo.SpacesKeyUpdateRequest{
						Name: "updated-key",
					}).
					Return(testKey, nil, nil).
					Times(1)
			},
		},
		{
			name: "Missing AccessKey parameter",
			args: map[string]any{
				"Name": "updated-key",
			},
			expectError: true,
		},
		{
			name: "Missing Name parameter",
			args: map[string]any{
				"AccessKey": "AKIA123456789",
			},
			expectError: true,
		},
		{
			name: "Invalid AccessKey type",
			args: map[string]any{
				"AccessKey": 123,
				"Name":      "updated-key",
			},
			expectError: true,
		},
		{
			name: "Invalid Name type",
			args: map[string]any{
				"AccessKey": "AKIA123456789",
				"Name":      123,
			},
			expectError: true,
		},
		{
			name: "Empty AccessKey",
			args: map[string]any{
				"AccessKey": "",
				"Name":      "updated-key",
			},
			expectError: true,
		},
		{
			name: "Empty Name",
			args: map[string]any{
				"AccessKey": "AKIA123456789",
				"Name":      "",
			},
			expectError: true,
		},
		{
			name: "API error",
			args: map[string]any{
				"AccessKey": "AKIA987654321",
				"Name":      "fail-key",
			},
			mockSetup: func(m *MockSpacesKeysService) {
				m.EXPECT().
					Update(gomock.Any(), "AKIA987654321", &godo.SpacesKeyUpdateRequest{
						Name: "fail-key",
					}).
					Return(nil, nil, errors.New("api error")).
					Times(1)
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSpacesKeys := NewMockSpacesKeysService(ctrl)
			if tc.mockSetup != nil {
				tc.mockSetup(mockSpacesKeys)
			}
			tool := setupSpacesKeysToolWithMock(mockSpacesKeys)
			req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: tc.args}}
			resp, err := tool.updateSpacesKey(context.Background(), req)
			if tc.expectError {
				require.NotNil(t, resp)
				require.True(t, resp.IsError)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.False(t, resp.IsError)
			var outKey godo.SpacesKey
			require.NoError(t, json.Unmarshal([]byte(resp.Content[0].(mcp.TextContent).Text), &outKey))
			require.Equal(t, testKey.AccessKey, outKey.AccessKey)
			require.Equal(t, testKey.Name, outKey.Name)
		})
	}
}

func TestSpacesKeysTool_deleteSpacesKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		args        map[string]any
		mockSetup   func(*MockSpacesKeysService)
		expectError bool
		expectText  string
	}{
		{
			name: "Successful delete",
			args: map[string]any{"AccessKey": "AKIA123456789"},
			mockSetup: func(m *MockSpacesKeysService) {
				m.EXPECT().
					Delete(gomock.Any(), "AKIA123456789").
					Return(&godo.Response{}, nil).
					Times(1)
			},
			expectText: "Spaces key deleted successfully",
		},
		{
			name:        "Missing AccessKey parameter",
			args:        map[string]any{},
			expectError: true,
		},
		{
			name:        "Invalid AccessKey type",
			args:        map[string]any{"AccessKey": 123},
			expectError: true,
		},
		{
			name:        "Empty AccessKey",
			args:        map[string]any{"AccessKey": ""},
			expectError: true,
		},
		{
			name: "API error",
			args: map[string]any{"AccessKey": "AKIA987654321"},
			mockSetup: func(m *MockSpacesKeysService) {
				m.EXPECT().
					Delete(gomock.Any(), "AKIA987654321").
					Return(nil, errors.New("api error")).
					Times(1)
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSpacesKeys := NewMockSpacesKeysService(ctrl)
			if tc.mockSetup != nil {
				tc.mockSetup(mockSpacesKeys)
			}
			tool := setupSpacesKeysToolWithMock(mockSpacesKeys)
			req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: tc.args}}
			resp, err := tool.deleteSpacesKey(context.Background(), req)
			if tc.expectError {
				require.NotNil(t, resp)
				require.True(t, resp.IsError)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.False(t, resp.IsError)
			require.Contains(t, resp.Content[0].(mcp.TextContent).Text, tc.expectText)
		})
	}
}

func TestSpacesKeysTool_listSpacesKeys(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testKeys := []*godo.SpacesKey{
		{
			Name:      "test-key-1",
			AccessKey: "AKIA123456789",
		},
		{
			Name:      "test-key-2",
			AccessKey: "AKIA987654321",
		},
	}

	tests := []struct {
		name        string
		mockSetup   func(*MockSpacesKeysService)
		expectError bool
	}{
		{
			name: "Successful list",
			mockSetup: func(m *MockSpacesKeysService) {
				m.EXPECT().
					List(gomock.Any(), &godo.ListOptions{}).
					Return(testKeys, nil, nil).
					Times(1)
			},
		},
		{
			name: "API error",
			mockSetup: func(m *MockSpacesKeysService) {
				m.EXPECT().
					List(gomock.Any(), &godo.ListOptions{}).
					Return(nil, nil, errors.New("api error")).
					Times(1)
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSpacesKeys := NewMockSpacesKeysService(ctrl)
			if tc.mockSetup != nil {
				tc.mockSetup(mockSpacesKeys)
			}
			tool := setupSpacesKeysToolWithMock(mockSpacesKeys)
			req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{}}}
			resp, err := tool.listSpacesKeys(context.Background(), req)
			if tc.expectError {
				require.NotNil(t, resp)
				require.True(t, resp.IsError)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.False(t, resp.IsError)
			var outKeys []godo.SpacesKey
			require.NoError(t, json.Unmarshal([]byte(resp.Content[0].(mcp.TextContent).Text), &outKeys))
			require.Len(t, outKeys, 2)
			require.Equal(t, testKeys[0].Name, outKeys[0].Name)
			require.Equal(t, testKeys[1].Name, outKeys[1].Name)
		})
	}
}

func TestSpacesKeysTool_getSpacesKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testKey := &godo.SpacesKey{
		Name:      "test-key",
		AccessKey: "AKIA123456789",
	}

	tests := []struct {
		name        string
		args        map[string]any
		mockSetup   func(*MockSpacesKeysService)
		expectError bool
	}{
		{
			name: "Successful get",
			args: map[string]any{
				"AccessKey": "AKIA123456789",
			},
			mockSetup: func(m *MockSpacesKeysService) {
				m.EXPECT().
					Get(gomock.Any(), "AKIA123456789").
					Return(testKey, nil, nil).
					Times(1)
			},
		},
		{
			name:        "Missing AccessKey parameter",
			args:        map[string]any{},
			expectError: true,
		},
		{
			name: "Invalid AccessKey type",
			args: map[string]any{
				"AccessKey": 123,
			},
			expectError: true,
		},
		{
			name: "Empty AccessKey",
			args: map[string]any{
				"AccessKey": "",
			},
			expectError: true,
		},
		{
			name: "API error",
			args: map[string]any{
				"AccessKey": "AKIA987654321",
			},
			mockSetup: func(m *MockSpacesKeysService) {
				m.EXPECT().
					Get(gomock.Any(), "AKIA987654321").
					Return(nil, nil, errors.New("api error")).
					Times(1)
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSpacesKeys := NewMockSpacesKeysService(ctrl)
			if tc.mockSetup != nil {
				tc.mockSetup(mockSpacesKeys)
			}
			tool := setupSpacesKeysToolWithMock(mockSpacesKeys)
			req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: tc.args}}
			resp, err := tool.getSpacesKey(context.Background(), req)
			if tc.expectError {
				require.NotNil(t, resp)
				require.True(t, resp.IsError)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.False(t, resp.IsError)
			var outKey godo.SpacesKey
			require.NoError(t, json.Unmarshal([]byte(resp.Content[0].(mcp.TextContent).Text), &outKey))
			require.Equal(t, testKey.Name, outKey.Name)
			require.Equal(t, testKey.AccessKey, outKey.AccessKey)
		})
	}
}

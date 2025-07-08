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

func setupKeysResourceWithMock(spacesKeys *MockSpacesKeysService) *KeysIPMCPResource {
	client := &godo.Client{}
	client.SpacesKeys = spacesKeys
	return NewKeysIPMCPResource(client)
}

func TestKeysIPMCPResource_handleGetKeysResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testKey := &godo.SpacesKey{
		Name:      "test-key",
		AccessKey: "access-key-123",
	}
	tests := []struct {
		name        string
		uri         string
		mockSetup   func(*MockSpacesKeysService)
		expectError bool
	}{
		{
			name: "Successful get",
			uri:  "spaces_keys://key-123",
			mockSetup: func(m *MockSpacesKeysService) {
				m.EXPECT().
					Get(gomock.Any(), "key-123").
					Return(testKey, nil, nil)
			},
		},
		{
			name: "API error",
			uri:  "spaces_keys://key-456",
			mockSetup: func(m *MockSpacesKeysService) {
				m.EXPECT().
					Get(gomock.Any(), "key-456").
					Return(nil, nil, errors.New("api error")).
					Times(1)
			},
			expectError: true,
		},
		{
			name:        "Invalid URI",
			uri:         "spaces_keys123",
			mockSetup:   nil,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSpacesKeys := NewMockSpacesKeysService(ctrl)
			if tc.mockSetup != nil {
				tc.mockSetup(mockSpacesKeys)
			}
			resource := setupKeysResourceWithMock(mockSpacesKeys)
			req := mcp.ReadResourceRequest{
				Params: mcp.ReadResourceParams{
					URI: tc.uri,
				},
			}
			resp, err := resource.handleGetKeysResource(context.Background(), req)
			if tc.expectError {
				require.Error(t, err)
				require.Nil(t, resp)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Len(t, resp, 1)
			content, ok := resp[0].(mcp.TextResourceContents)
			require.True(t, ok)
			require.Equal(t, tc.uri, content.URI)
			require.Equal(t, "application/json", content.MIMEType)
			var outKey godo.SpacesKey
			require.NoError(t, json.Unmarshal([]byte(content.Text), &outKey))
			require.Equal(t, testKey.AccessKey, outKey.AccessKey)
			require.Equal(t, testKey.Name, outKey.Name)
		})
	}
}

func TestKeysIPMCPResource_handleGetKeysListResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testKeys := []*godo.SpacesKey{
		{
			Name:      "test-key-1",
			AccessKey: "access-key-123",
		},
		{
			Name:      "test-key-2",
			AccessKey: "access-key-456",
		},
	}
	tests := []struct {
		name        string
		uri         string
		mockSetup   func(*MockSpacesKeysService)
		expectError bool
	}{
		{
			name: "Successful list",
			uri:  "spaces_keys://all",
			mockSetup: func(m *MockSpacesKeysService) {
				m.EXPECT().
					List(gomock.Any(), gomock.Any()).
					Return(testKeys, nil, nil)
			},
		},
		{
			name: "API error",
			uri:  "spaces_keys://all",
			mockSetup: func(m *MockSpacesKeysService) {
				m.EXPECT().
					List(gomock.Any(), gomock.Any()).
					Return(nil, nil, errors.New("api error"))
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
			resource := setupKeysResourceWithMock(mockSpacesKeys)
			req := mcp.ReadResourceRequest{
				Params: mcp.ReadResourceParams{
					URI: tc.uri,
				},
			}
			resp, err := resource.handleGetKeysListResource(context.Background(), req)
			if tc.expectError {
				require.Error(t, err)
				require.Nil(t, resp)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Len(t, resp, 1)
			content, ok := resp[0].(mcp.TextResourceContents)
			require.True(t, ok)
			require.Equal(t, tc.uri, content.URI)
			require.Equal(t, "application/json", content.MIMEType)
			var outKeys []godo.SpacesKey
			require.NoError(t, json.Unmarshal([]byte(content.Text), &outKeys))
			require.Len(t, outKeys, len(testKeys))
		})
	}
}

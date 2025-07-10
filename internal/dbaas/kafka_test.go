package dbaas

import (
	"context"
	"encoding/json"
	"testing"

	"mcp-digitalocean/internal/dbaas/mocks"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestKafkaTool_getKafkaConfig(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("GetKafkaConfig", mock.Anything, "cid").Return(&godo.KafkaConfig{}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	kt := &KafkaTool{client: client}
	args := map[string]interface{}{"ID": "cid"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := kt.getKafkaConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "{")
	mockDB.AssertExpectations(t)
	// Error case: missing ID
	reqMissing := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{}}}
	res, err = kt.getKafkaConfig(context.Background(), reqMissing)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Cluster ID is required")
}

func TestKafkaTool_updateKafkaConfig(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("UpdateKafkaConfig", mock.Anything, "cid", mock.Anything).Return(nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	kt := &KafkaTool{client: client}
	cfg := godo.KafkaConfig{}
	cfgJSON, _ := json.Marshal(cfg)
	args := map[string]interface{}{"ID": "cid", "config_json": string(cfgJSON)}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := kt.updateKafkaConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Kafka config updated successfully")
	mockDB.AssertExpectations(t)
	// Error case: missing config_json
	args = map[string]interface{}{"ID": "cid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = kt.updateKafkaConfig(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "config_json is required")
}

func TestKafkaTool_listTopics(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("ListTopics", mock.Anything, "cid", mock.Anything).Return([]godo.DatabaseTopic{{Name: "topic1"}}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	kt := &KafkaTool{client: client}
	args := map[string]interface{}{"ID": "cid"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := kt.listTopics(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "topic1")
	mockDB.AssertExpectations(t)
	// Error case: missing ID
	reqMissing := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{}}}
	res, err = kt.listTopics(context.Background(), reqMissing)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Cluster ID is required")
}

func TestKafkaTool_createTopic(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("CreateTopic", mock.Anything, "cid", mock.Anything).Return(&godo.DatabaseTopic{Name: "topic2"}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	kt := &KafkaTool{client: client}
	args := map[string]interface{}{"ID": "cid", "name": "topic2"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := kt.createTopic(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "topic2")
	mockDB.AssertExpectations(t)
	// Error case: missing name
	args = map[string]interface{}{"ID": "cid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = kt.createTopic(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Topic name is required")
}

func TestKafkaTool_getTopic(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("GetTopic", mock.Anything, "cid", "topic3").Return(&godo.DatabaseTopic{Name: "topic3"}, nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	kt := &KafkaTool{client: client}
	args := map[string]interface{}{"ID": "cid", "name": "topic3"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := kt.getTopic(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "topic3")
	mockDB.AssertExpectations(t)
	// Error case: missing name
	args = map[string]interface{}{"ID": "cid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = kt.getTopic(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Topic name is required")
}

func TestKafkaTool_deleteTopic(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("DeleteTopic", mock.Anything, "cid", "topic4").Return(nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	kt := &KafkaTool{client: client}
	args := map[string]interface{}{"ID": "cid", "name": "topic4"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := kt.deleteTopic(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Topic deleted successfully")
	mockDB.AssertExpectations(t)
	// Error case: missing name
	args = map[string]interface{}{"ID": "cid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = kt.deleteTopic(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Topic name is required")
}

func TestKafkaTool_updateTopic(t *testing.T) {
	mockDB := &mocks.DatabasesService{}
	mockDB.On("UpdateTopic", mock.Anything, "cid", "topic5", mock.Anything).Return(nil, nil)
	client := &godo.Client{}
	client.Databases = mockDB
	kt := &KafkaTool{client: client}
	args := map[string]interface{}{"ID": "cid", "name": "topic5"}
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err := kt.updateTopic(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Topic updated successfully")
	mockDB.AssertExpectations(t)
	// Error case: missing name
	args = map[string]interface{}{"ID": "cid"}
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
	res, err = kt.updateTopic(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Topic name is required")
}

//go:build integration

package testing

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

// TestUptimeCheckLifecycle tests the full lifecycle of an uptime check:
// create -> get -> list -> get state -> update -> delete
func TestUptimeCheckLifecycle(t *testing.T) {
	ctx, c, _, cleanup := setupTest(t)
	defer cleanup()

	checkName := fmt.Sprintf("test-uptime-check-%d", time.Now().Unix())
	target := "https://www.digitalocean.com"

	// create uptime check
	t.Log("creating uptime check...")
	createResp, err := c.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "uptimecheck-create",
			Arguments: map[string]interface{}{
				"Name":    checkName,
				"Type":    "https",
				"Target":  target,
				"Regions": []string{"us_east", "us_west"},
				"Enabled": true,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, createResp)
	require.False(t, createResp.IsError, "uptimecheck-create failed: %v", createResp.Content)

	var uptimeCheck godo.UptimeCheck
	err = json.Unmarshal([]byte(createResp.Content[0].(mcp.TextContent).Text), &uptimeCheck)
	require.NoError(t, err)
	require.NotEmpty(t, uptimeCheck.ID, "uptime check ID should not be empty")
	t.Logf("created uptime check: %s (ID: %s)", uptimeCheck.Name, uptimeCheck.ID)

	// cleanup on test completion
	defer func() {
		t.Logf("deleting uptime check %s...", uptimeCheck.ID)
		deleteResp, err := c.CallTool(ctx, mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "uptimecheck-delete",
				Arguments: map[string]interface{}{
					"ID": uptimeCheck.ID,
				},
			},
		})
		if err != nil {
			t.Logf("failed to delete uptime check: %v", err)
			return
		}
		if deleteResp.IsError {
			t.Logf("uptimecheck-delete returned error: %v", deleteResp.Content)
			return
		}
		t.Logf("deleted uptime check %s", uptimeCheck.ID)
	}()

	// get uptime check
	t.Log("getting uptime check...")
	getResp, err := c.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "uptimecheck-get",
			Arguments: map[string]interface{}{
				"ID": uptimeCheck.ID,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, getResp)
	require.False(t, getResp.IsError, "uptimecheck-get failed: %v", getResp.Content)

	var fetchedCheck godo.UptimeCheck
	err = json.Unmarshal([]byte(getResp.Content[0].(mcp.TextContent).Text), &fetchedCheck)
	require.NoError(t, err)
	require.Equal(t, uptimeCheck.ID, fetchedCheck.ID)
	require.Equal(t, checkName, fetchedCheck.Name)
	t.Logf("fetched uptime check: %s", fetchedCheck.Name)

	// list uptime checks
	t.Log("listing uptime checks...")
	listResp, err := c.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "uptimecheck-list",
			Arguments: map[string]interface{}{
				"Page":    1,
				"PerPage": 50,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.False(t, listResp.IsError, "uptimecheck-list failed: %v", listResp.Content)

	var checks []godo.UptimeCheck
	err = json.Unmarshal([]byte(listResp.Content[0].(mcp.TextContent).Text), &checks)
	require.NoError(t, err)

	found := false
	for _, check := range checks {
		if check.ID == uptimeCheck.ID {
			found = true
			break
		}
	}
	require.True(t, found, "created uptime check not found in list")
	t.Logf("found uptime check in list (total: %d)", len(checks))

	// get uptime check state
	t.Log("getting uptime check state...")
	stateResp, err := c.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "uptimecheck-get-state",
			Arguments: map[string]interface{}{
				"ID": uptimeCheck.ID,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, stateResp)
	require.False(t, stateResp.IsError, "uptimecheck-get-state failed: %v", stateResp.Content)

	var checkState godo.UptimeCheckState
	err = json.Unmarshal([]byte(stateResp.Content[0].(mcp.TextContent).Text), &checkState)
	require.NoError(t, err)
	t.Logf("uptime check state: %+v", checkState)

	// update uptime check
	updatedName := checkName + "-updated"
	t.Log("updating uptime check...")
	updateResp, err := c.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "uptimecheck-update",
			Arguments: map[string]interface{}{
				"ID":      uptimeCheck.ID,
				"Name":    updatedName,
				"Type":    "https",
				"Target":  target,
				"Regions": []string{"us_east", "eu_west"},
				"Enabled": true,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.False(t, updateResp.IsError, "uptimecheck-update failed: %v", updateResp.Content)

	var updatedCheck godo.UptimeCheck
	err = json.Unmarshal([]byte(updateResp.Content[0].(mcp.TextContent).Text), &updatedCheck)
	require.NoError(t, err)
	require.Equal(t, updatedName, updatedCheck.Name)
	t.Logf("updated uptime check name to: %s", updatedCheck.Name)
}

// TestUptimeCheckAlertLifecycle tests the full lifecycle of an uptime check alert:
// create check -> create alert -> get alert -> list alerts -> update alert -> delete alert -> delete check
func TestUptimeCheckAlertLifecycle(t *testing.T) {
	ctx, c, _, cleanup := setupTest(t)
	defer cleanup()

	checkName := fmt.Sprintf("test-uptime-alert-check-%d", time.Now().Unix())
	target := "https://www.digitalocean.com"

	// create uptime check first (required for alerts)
	t.Log("creating uptime check for alert testing...")
	createCheckResp, err := c.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "uptimecheck-create",
			Arguments: map[string]interface{}{
				"Name":    checkName,
				"Type":    "https",
				"Target":  target,
				"Regions": []string{"us_east"},
				"Enabled": true,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, createCheckResp)
	require.False(t, createCheckResp.IsError, "uptimecheck-create failed: %v", createCheckResp.Content)

	var uptimeCheck godo.UptimeCheck
	err = json.Unmarshal([]byte(createCheckResp.Content[0].(mcp.TextContent).Text), &uptimeCheck)
	require.NoError(t, err)
	t.Logf("created uptime check: %s (ID: %s)", uptimeCheck.Name, uptimeCheck.ID)

	defer func() {
		t.Logf("deleting uptime check %s...", uptimeCheck.ID)
		_, _ = c.CallTool(ctx, mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "uptimecheck-delete",
				Arguments: map[string]interface{}{
					"ID": uptimeCheck.ID,
				},
			},
		})
	}()

	// create uptime check alert
	alertName := fmt.Sprintf("test-alert-%d", time.Now().Unix())
	t.Log("creating uptime check alert...")
	createAlertResp, err := c.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "uptimecheck-alert-create",
			Arguments: map[string]interface{}{
				"CheckID":      uptimeCheck.ID,
				"Name":         alertName,
				"Type":         "down",
				"Period":       "2m",
				"Emails":       []string{},
				"SlackDetails": []map[string]string{},
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, createAlertResp)
	require.False(t, createAlertResp.IsError, "uptimecheck-alert-create failed: %v", createAlertResp.Content)

	var alert godo.UptimeAlert
	err = json.Unmarshal([]byte(createAlertResp.Content[0].(mcp.TextContent).Text), &alert)
	require.NoError(t, err)
	require.NotEmpty(t, alert.ID, "alert ID should not be empty")
	t.Logf("created uptime check alert: %s (ID: %s)", alert.Name, alert.ID)

	defer func() {
		t.Logf("deleting uptime check alert %s...", alert.ID)
		_, _ = c.CallTool(ctx, mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "uptimecheck-alert-delete",
				Arguments: map[string]interface{}{
					"CheckID": uptimeCheck.ID,
					"AlertID": alert.ID,
				},
			},
		})
	}()

	// get uptime check alert
	t.Log("getting uptime check alert...")
	getAlertResp, err := c.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "uptimecheck-alert-get",
			Arguments: map[string]interface{}{
				"CheckID": uptimeCheck.ID,
				"AlertID": alert.ID,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, getAlertResp)
	require.False(t, getAlertResp.IsError, "uptimecheck-alert-get failed: %v", getAlertResp.Content)

	var fetchedAlert godo.UptimeAlert
	err = json.Unmarshal([]byte(getAlertResp.Content[0].(mcp.TextContent).Text), &fetchedAlert)
	require.NoError(t, err)
	require.Equal(t, alert.ID, fetchedAlert.ID)
	t.Logf("fetched uptime check alert: %s", fetchedAlert.Name)

	// list uptime check alerts
	t.Log("listing uptime check alerts...")
	listAlertResp, err := c.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "uptimecheck-alert-list",
			Arguments: map[string]interface{}{
				"CheckID": uptimeCheck.ID,
				"Page":    1,
				"PerPage": 50,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, listAlertResp)
	require.False(t, listAlertResp.IsError, "uptimecheck-alert-list failed: %v", listAlertResp.Content)

	var alerts []godo.UptimeAlert
	err = json.Unmarshal([]byte(listAlertResp.Content[0].(mcp.TextContent).Text), &alerts)
	require.NoError(t, err)

	found := false
	for _, a := range alerts {
		if a.ID == alert.ID {
			found = true
			break
		}
	}
	require.True(t, found, "created alert not found in list")
	t.Logf("found alert in list (total: %d)", len(alerts))

	// update uptime check alert
	updatedAlertName := alertName + "-updated"
	t.Log("updating uptime check alert...")
	updateAlertResp, err := c.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "uptimecheck-alert-update",
			Arguments: map[string]interface{}{
				"CheckID":      uptimeCheck.ID,
				"AlertID":      alert.ID,
				"Name":         updatedAlertName,
				"Type":         "down",
				"Period":       "3m",
				"Emails":       []string{},
				"SlackDetails": []map[string]string{},
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, updateAlertResp)
	require.False(t, updateAlertResp.IsError, "uptimecheck-alert-update failed: %v", updateAlertResp.Content)

	var updatedAlert godo.UptimeAlert
	err = json.Unmarshal([]byte(updateAlertResp.Content[0].(mcp.TextContent).Text), &updatedAlert)
	require.NoError(t, err)
	require.Equal(t, updatedAlertName, updatedAlert.Name)
	t.Logf("updated alert name to: %s", updatedAlert.Name)
}

// TestAlertPolicyLifecycle tests the full lifecycle of an alert policy:
// create -> get -> list -> update -> delete
func TestAlertPolicyLifecycle(t *testing.T) {
	ctx, c, _, cleanup := setupTest(t)
	defer cleanup()

	policyDescription := fmt.Sprintf("test-policy-%d", time.Now().Unix())

	// create alert policy
	t.Log("creating alert policy...")
	createResp, err := c.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "alert-policy-create",
			Arguments: map[string]interface{}{
				"Type":        "v1/insights/droplet/cpu",
				"Description": policyDescription,
				"Compare":     "GreaterThan",
				"Value":       80,
				"Window":      "5m",
				"Entities":    []string{},
				"Tags":        []string{"test-tag"},
				"Alerts": map[string]interface{}{
					"Email": []string{},
					"Slack": []interface{}{},
				},
				"Enabled": true,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, createResp)
	require.False(t, createResp.IsError, "alert-policy-create failed: %v", createResp.Content)

	var policy godo.AlertPolicy
	err = json.Unmarshal([]byte(createResp.Content[0].(mcp.TextContent).Text), &policy)
	require.NoError(t, err)
	require.NotEmpty(t, policy.UUID, "policy UUID should not be empty")
	t.Logf("created alert policy: %s (UUID: %s)", policy.Description, policy.UUID)

	defer func() {
		t.Logf("deleting alert policy %s...", policy.UUID)
		deleteResp, err := c.CallTool(ctx, mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "alert-policy-delete",
				Arguments: map[string]interface{}{
					"UUID": policy.UUID,
				},
			},
		})
		if err != nil {
			t.Logf("failed to delete alert policy: %v", err)
			return
		}
		if deleteResp.IsError {
			t.Logf("alert-policy-delete returned error: %v", deleteResp.Content)
			return
		}
		t.Logf("deleted alert policy %s", policy.UUID)
	}()

	// get alert policy
	t.Log("getting alert policy...")
	getResp, err := c.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "alert-policy-get",
			Arguments: map[string]interface{}{
				"UUID": policy.UUID,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, getResp)
	require.False(t, getResp.IsError, "alert-policy-get failed: %v", getResp.Content)

	var fetchedPolicy godo.AlertPolicy
	err = json.Unmarshal([]byte(getResp.Content[0].(mcp.TextContent).Text), &fetchedPolicy)
	require.NoError(t, err)
	require.Equal(t, policy.UUID, fetchedPolicy.UUID)
	t.Logf("fetched alert policy: %s", fetchedPolicy.Description)

	// list alert policies
	t.Log("listing alert policies...")
	listResp, err := c.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "alert-policy-list",
			Arguments: map[string]interface{}{
				"Page":    1,
				"PerPage": 50,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.False(t, listResp.IsError, "alert-policy-list failed: %v", listResp.Content)

	var policies []godo.AlertPolicy
	err = json.Unmarshal([]byte(listResp.Content[0].(mcp.TextContent).Text), &policies)
	require.NoError(t, err)

	found := false
	for _, p := range policies {
		if p.UUID == policy.UUID {
			found = true
			break
		}
	}
	require.True(t, found, "created policy not found in list")
	t.Logf("found policy in list (total: %d)", len(policies))

	// update alert policy
	updatedDescription := policyDescription + "-updated"
	t.Log("updating alert policy...")
	updateResp, err := c.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "alert-policy-update",
			Arguments: map[string]interface{}{
				"UUID":        policy.UUID,
				"Type":        "v1/insights/droplet/cpu",
				"Description": updatedDescription,
				"Compare":     "GreaterThan",
				"Value":       90,
				"Window":      "10m",
				"Entities":    []string{},
				"Tags":        []string{"test-tag", "updated-tag"},
				"Alerts": map[string]interface{}{
					"Email": []string{},
					"Slack": []interface{}{},
				},
				"Enabled": true,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.False(t, updateResp.IsError, "alert-policy-update failed: %v", updateResp.Content)

	var updatedPolicy godo.AlertPolicy
	err = json.Unmarshal([]byte(updateResp.Content[0].(mcp.TextContent).Text), &updatedPolicy)
	require.NoError(t, err)
	require.Equal(t, updatedDescription, updatedPolicy.Description)
	t.Logf("updated policy description to: %s", updatedPolicy.Description)
}

// assertUptimeCheckExists verifies an uptime check exists by ID
func assertUptimeCheckExists(ctx context.Context, t *testing.T, c interface{ CallTool(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) }, checkID string) {
	resp, err := c.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "uptimecheck-get",
			Arguments: map[string]interface{}{
				"ID": checkID,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.False(t, resp.IsError, "uptime check %s not found: %v", checkID, resp.Content)
}

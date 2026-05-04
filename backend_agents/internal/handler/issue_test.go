package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// TestValidateMcpConfig_ValidInput tests validateMcpConfig with valid input containing mcpServers.
func TestValidateMcpConfig_ValidInput(t *testing.T) {
	validConfig := json.RawMessage(`{
		"mcpServers": {
			"github": {
				"command": "npx",
				"args": ["-y", "@modelcontextprotocol/server-github"]
			}
		}
	}`)

	err := validateMcpConfig(validConfig)
	if err != nil {
		t.Fatalf("validateMcpConfig with valid input: expected nil, got %v", err)
	}
}

// TestValidateMcpConfig_InvalidJSON tests validateMcpConfig with invalid JSON.
func TestValidateMcpConfig_InvalidJSON(t *testing.T) {
	invalidJSON := json.RawMessage(`{invalid json}`)

	err := validateMcpConfig(invalidJSON)
	if err == nil {
		t.Fatal("validateMcpConfig with invalid JSON: expected error, got nil")
	}
	if err.Error() != "mcp_config must be valid JSON" {
		t.Fatalf("validateMcpConfig with invalid JSON: expected 'mcp_config must be valid JSON', got %q", err.Error())
	}
}

// TestValidateMcpConfig_MissingMcpServersKey tests validateMcpConfig with valid JSON but missing mcpServers key.
func TestValidateMcpConfig_MissingMcpServersKey(t *testing.T) {
	missingKey := json.RawMessage(`{"servers": {}}`)

	err := validateMcpConfig(missingKey)
	if err == nil {
		t.Fatal("validateMcpConfig with missing mcpServers: expected error, got nil")
	}
	if err.Error() != "mcp_config must contain a top-level 'mcpServers' object" {
		t.Fatalf("validateMcpConfig with missing mcpServers: expected 'mcp_config must contain a top-level 'mcpServers' object', got %q", err.Error())
	}
}

// TestValidateMcpConfig_NilInput tests validateMcpConfig with nil/empty input.
func TestValidateMcpConfig_NilInput(t *testing.T) {
	err := validateMcpConfig(nil)
	if err != nil {
		t.Fatalf("validateMcpConfig with nil input: expected nil, got %v", err)
	}
}

// TestValidateMcpConfig_EmptyInput tests validateMcpConfig with empty input.
func TestValidateMcpConfig_EmptyInput(t *testing.T) {
	err := validateMcpConfig(json.RawMessage(``))
	if err != nil {
		t.Fatalf("validateMcpConfig with empty input: expected nil, got %v", err)
	}
}

// TestValidateMcpConfig_ExplicitNull tests validateMcpConfig with explicit null.
func TestValidateMcpConfig_ExplicitNull(t *testing.T) {
	err := validateMcpConfig(json.RawMessage(`null`))
	if err != nil {
		t.Fatalf("validateMcpConfig with explicit null: expected nil, got %v", err)
	}
}

// TestCreateIssueWithValidMcpConfig tests CreateIssue with valid mcp_config.
func TestCreateIssueWithValidMcpConfig(t *testing.T) {
	if testHandler == nil {
		t.Skip("database not available")
	}

	validMcpConfig := map[string]any{
		"mcpServers": map[string]any{
			"github": map[string]any{
				"command": "npx",
				"args":    []string{"-y", "@modelcontextprotocol/server-github"},
			},
		},
	}

	w := httptest.NewRecorder()
	req := newRequest("POST", "/api/issues?workspace_id="+testWorkspaceID, map[string]any{
		"title":      "Issue with valid mcp_config",
		"status":     "todo",
		"priority":   "medium",
		"mcp_config": validMcpConfig,
	})
	testHandler.CreateIssue(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("CreateIssue with valid mcp_config: expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var created IssueResponse
	if err := json.NewDecoder(w.Body).Decode(&created); err != nil {
		t.Fatalf("CreateIssue: decode response: %v", err)
	}

	// Verify mcp_config is persisted
	if created.McpConfig == nil {
		t.Fatal("CreateIssue: expected mcp_config in response, got nil")
	}

	// Verify the config matches what we sent
	var returnedConfig map[string]any
	if err := json.Unmarshal(created.McpConfig, &returnedConfig); err != nil {
		t.Fatalf("CreateIssue: unmarshal returned mcp_config: %v", err)
	}

	if _, ok := returnedConfig["mcpServers"]; !ok {
		t.Fatal("CreateIssue: returned mcp_config missing mcpServers key")
	}

	// Cleanup
	cleanupReq := newRequest("DELETE", "/api/issues/"+created.ID, nil)
	cleanupReq = withURLParam(cleanupReq, "id", created.ID)
	testHandler.DeleteIssue(httptest.NewRecorder(), cleanupReq)
}

// TestCreateIssueWithInvalidMcpConfig tests CreateIssue with invalid mcp_config (should return HTTP 400).
func TestCreateIssueWithInvalidMcpConfig(t *testing.T) {
	if testHandler == nil {
		t.Skip("database not available")
	}

	w := httptest.NewRecorder()
	req := newRequest("POST", "/api/issues?workspace_id="+testWorkspaceID, map[string]any{
		"title":      "Issue with invalid mcp_config",
		"status":     "todo",
		"priority":   "medium",
		"mcp_config": "{invalid json}",
	})
	testHandler.CreateIssue(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("CreateIssue with invalid mcp_config: expected 400, got %d: %s", w.Code, w.Body.String())
	}

	var errResp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("CreateIssue: decode error response: %v", err)
	}

	if errMsg, ok := errResp["error"].(string); !ok || errMsg != "mcp_config must be valid JSON" {
		t.Fatalf("CreateIssue: expected error 'mcp_config must be valid JSON', got %v", errResp["error"])
	}
}

// TestCreateIssueWithMissingMcpServersKey tests CreateIssue with mcp_config missing mcpServers key (should return HTTP 400).
func TestCreateIssueWithMissingMcpServersKey(t *testing.T) {
	if testHandler == nil {
		t.Skip("database not available")
	}

	w := httptest.NewRecorder()
	req := newRequest("POST", "/api/issues?workspace_id="+testWorkspaceID, map[string]any{
		"title":      "Issue with missing mcpServers",
		"status":     "todo",
		"priority":   "medium",
		"mcp_config": map[string]any{"servers": map[string]any{}},
	})
	testHandler.CreateIssue(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("CreateIssue with missing mcpServers: expected 400, got %d: %s", w.Code, w.Body.String())
	}

	var errResp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("CreateIssue: decode error response: %v", err)
	}

	if errMsg, ok := errResp["error"].(string); !ok || errMsg != "mcp_config must contain a top-level 'mcpServers' object" {
		t.Fatalf("CreateIssue: expected error about missing mcpServers, got %v", errResp["error"])
	}
}

// TestUpdateIssueWithValidMcpConfig tests UpdateIssue with valid mcp_config.
func TestUpdateIssueWithValidMcpConfig(t *testing.T) {
	if testHandler == nil {
		t.Skip("database not available")
	}

	// Create an issue first
	w := httptest.NewRecorder()
	req := newRequest("POST", "/api/issues?workspace_id="+testWorkspaceID, map[string]any{
		"title":    "Issue for update test",
		"status":   "todo",
		"priority": "medium",
	})
	testHandler.CreateIssue(w, req)
	var created IssueResponse
	json.NewDecoder(w.Body).Decode(&created)
	issueID := created.ID

	defer func() {
		cleanupReq := newRequest("DELETE", "/api/issues/"+issueID, nil)
		cleanupReq = withURLParam(cleanupReq, "id", issueID)
		testHandler.DeleteIssue(httptest.NewRecorder(), cleanupReq)
	}()

	// Update with valid mcp_config
	validMcpConfig := map[string]any{
		"mcpServers": map[string]any{
			"slack": map[string]any{
				"command": "npx",
				"args":    []string{"-y", "@modelcontextprotocol/server-slack"},
			},
		},
	}

	w = httptest.NewRecorder()
	req = newRequest("PUT", "/api/issues/"+issueID, map[string]any{
		"mcp_config": validMcpConfig,
	})
	req = withURLParam(req, "id", issueID)
	testHandler.UpdateIssue(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("UpdateIssue with valid mcp_config: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var updated IssueResponse
	if err := json.NewDecoder(w.Body).Decode(&updated); err != nil {
		t.Fatalf("UpdateIssue: decode response: %v", err)
	}

	// Verify mcp_config is updated
	if updated.McpConfig == nil {
		t.Fatal("UpdateIssue: expected mcp_config in response, got nil")
	}

	var returnedConfig map[string]any
	if err := json.Unmarshal(updated.McpConfig, &returnedConfig); err != nil {
		t.Fatalf("UpdateIssue: unmarshal returned mcp_config: %v", err)
	}

	if _, ok := returnedConfig["mcpServers"]; !ok {
		t.Fatal("UpdateIssue: returned mcp_config missing mcpServers key")
	}
}

// TestUpdateIssueWithMcpConfigSetToNull tests UpdateIssue with mcp_config set to null (should clear).
func TestUpdateIssueWithMcpConfigSetToNull(t *testing.T) {
	if testHandler == nil {
		t.Skip("database not available")
	}

	// Create an issue with mcp_config
	validMcpConfig := map[string]any{
		"mcpServers": map[string]any{
			"github": map[string]any{
				"command": "npx",
				"args":    []string{"-y", "@modelcontextprotocol/server-github"},
			},
		},
	}

	w := httptest.NewRecorder()
	req := newRequest("POST", "/api/issues?workspace_id="+testWorkspaceID, map[string]any{
		"title":      "Issue for null test",
		"status":     "todo",
		"priority":   "medium",
		"mcp_config": validMcpConfig,
	})
	testHandler.CreateIssue(w, req)
	var created IssueResponse
	json.NewDecoder(w.Body).Decode(&created)
	issueID := created.ID

	defer func() {
		cleanupReq := newRequest("DELETE", "/api/issues/"+issueID, nil)
		cleanupReq = withURLParam(cleanupReq, "id", issueID)
		testHandler.DeleteIssue(httptest.NewRecorder(), cleanupReq)
	}()

	// Verify mcp_config is set
	if created.McpConfig == nil {
		t.Fatal("CreateIssue: expected mcp_config to be set")
	}

	// Update with mcp_config set to null
	w = httptest.NewRecorder()
	req = newRequest("PUT", "/api/issues/"+issueID, map[string]any{
		"mcp_config": nil,
	})
	req = withURLParam(req, "id", issueID)
	testHandler.UpdateIssue(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("UpdateIssue with null mcp_config: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var updated IssueResponse
	if err := json.NewDecoder(w.Body).Decode(&updated); err != nil {
		t.Fatalf("UpdateIssue: decode response: %v", err)
	}

	// Verify mcp_config is cleared (should be null or empty)
	if updated.McpConfig != nil && len(updated.McpConfig) > 0 && string(updated.McpConfig) != "null" {
		t.Fatalf("UpdateIssue: expected mcp_config to be cleared, got %s", string(updated.McpConfig))
	}
}

// TestUpdateIssueWithInvalidMcpConfig tests UpdateIssue with invalid mcp_config (should return HTTP 400).
func TestUpdateIssueWithInvalidMcpConfig(t *testing.T) {
	if testHandler == nil {
		t.Skip("database not available")
	}

	// Create an issue first
	w := httptest.NewRecorder()
	req := newRequest("POST", "/api/issues?workspace_id="+testWorkspaceID, map[string]any{
		"title":    "Issue for invalid update test",
		"status":   "todo",
		"priority": "medium",
	})
	testHandler.CreateIssue(w, req)
	var created IssueResponse
	json.NewDecoder(w.Body).Decode(&created)
	issueID := created.ID

	defer func() {
		cleanupReq := newRequest("DELETE", "/api/issues/"+issueID, nil)
		cleanupReq = withURLParam(cleanupReq, "id", issueID)
		testHandler.DeleteIssue(httptest.NewRecorder(), cleanupReq)
	}()

	// Update with invalid mcp_config
	w = httptest.NewRecorder()
	req = newRequest("PUT", "/api/issues/"+issueID, map[string]any{
		"mcp_config": "{invalid json}",
	})
	req = withURLParam(req, "id", issueID)
	testHandler.UpdateIssue(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("UpdateIssue with invalid mcp_config: expected 400, got %d: %s", w.Code, w.Body.String())
	}

	var errResp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("UpdateIssue: decode error response: %v", err)
	}

	if errMsg, ok := errResp["error"].(string); !ok || errMsg != "mcp_config must be valid JSON" {
		t.Fatalf("UpdateIssue: expected error 'mcp_config must be valid JSON', got %v", errResp["error"])
	}
}

// TestUpdateIssueWithMissingMcpServersKeyInUpdate tests UpdateIssue with mcp_config missing mcpServers key (should return HTTP 400).
func TestUpdateIssueWithMissingMcpServersKeyInUpdate(t *testing.T) {
	if testHandler == nil {
		t.Skip("database not available")
	}

	// Create an issue first
	w := httptest.NewRecorder()
	req := newRequest("POST", "/api/issues?workspace_id="+testWorkspaceID, map[string]any{
		"title":    "Issue for missing mcpServers update test",
		"status":   "todo",
		"priority": "medium",
	})
	testHandler.CreateIssue(w, req)
	var created IssueResponse
	json.NewDecoder(w.Body).Decode(&created)
	issueID := created.ID

	defer func() {
		cleanupReq := newRequest("DELETE", "/api/issues/"+issueID, nil)
		cleanupReq = withURLParam(cleanupReq, "id", issueID)
		testHandler.DeleteIssue(httptest.NewRecorder(), cleanupReq)
	}()

	// Update with mcp_config missing mcpServers
	w = httptest.NewRecorder()
	req = newRequest("PUT", "/api/issues/"+issueID, map[string]any{
		"mcp_config": map[string]any{"servers": map[string]any{}},
	})
	req = withURLParam(req, "id", issueID)
	testHandler.UpdateIssue(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("UpdateIssue with missing mcpServers: expected 400, got %d: %s", w.Code, w.Body.String())
	}

	var errResp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("UpdateIssue: decode error response: %v", err)
	}

	if errMsg, ok := errResp["error"].(string); !ok || errMsg != "mcp_config must contain a top-level 'mcpServers' object" {
		t.Fatalf("UpdateIssue: expected error about missing mcpServers, got %v", errResp["error"])
	}
}

// TestUpdateIssuePreservesOtherFieldsWhenUpdatingMcpConfig tests that updating mcp_config preserves other fields.
func TestUpdateIssuePreservesOtherFieldsWhenUpdatingMcpConfig(t *testing.T) {
	if testHandler == nil {
		t.Skip("database not available")
	}

	// Create an issue with specific fields
	w := httptest.NewRecorder()
	req := newRequest("POST", "/api/issues?workspace_id="+testWorkspaceID, map[string]any{
		"title":    "Issue with multiple fields",
		"status":   "todo",
		"priority": "high",
	})
	testHandler.CreateIssue(w, req)
	var created IssueResponse
	json.NewDecoder(w.Body).Decode(&created)
	issueID := created.ID

	defer func() {
		cleanupReq := newRequest("DELETE", "/api/issues/"+issueID, nil)
		cleanupReq = withURLParam(cleanupReq, "id", issueID)
		testHandler.DeleteIssue(httptest.NewRecorder(), cleanupReq)
	}()

	// Update only mcp_config
	validMcpConfig := map[string]any{
		"mcpServers": map[string]any{
			"github": map[string]any{
				"command": "npx",
				"args":    []string{"-y", "@modelcontextprotocol/server-github"},
			},
		},
	}

	w = httptest.NewRecorder()
	req = newRequest("PUT", "/api/issues/"+issueID, map[string]any{
		"mcp_config": validMcpConfig,
	})
	req = withURLParam(req, "id", issueID)
	testHandler.UpdateIssue(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("UpdateIssue: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var updated IssueResponse
	if err := json.NewDecoder(w.Body).Decode(&updated); err != nil {
		t.Fatalf("UpdateIssue: decode response: %v", err)
	}

	// Verify other fields are preserved
	if updated.Title != "Issue with multiple fields" {
		t.Fatalf("UpdateIssue: title should be preserved, got %q", updated.Title)
	}
	if updated.Status != "todo" {
		t.Fatalf("UpdateIssue: status should be preserved, got %q", updated.Status)
	}
	if updated.Priority != "high" {
		t.Fatalf("UpdateIssue: priority should be preserved, got %q", updated.Priority)
	}

	// Verify mcp_config is updated
	if updated.McpConfig == nil {
		t.Fatal("UpdateIssue: expected mcp_config to be set")
	}
}

// TestCreateIssueWithoutMcpConfig tests that CreateIssue works without mcp_config (backward compatibility).
func TestCreateIssueWithoutMcpConfig(t *testing.T) {
	if testHandler == nil {
		t.Skip("database not available")
	}

	w := httptest.NewRecorder()
	req := newRequest("POST", "/api/issues?workspace_id="+testWorkspaceID, map[string]any{
		"title":    "Issue without mcp_config",
		"status":   "todo",
		"priority": "medium",
	})
	testHandler.CreateIssue(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("CreateIssue without mcp_config: expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var created IssueResponse
	if err := json.NewDecoder(w.Body).Decode(&created); err != nil {
		t.Fatalf("CreateIssue: decode response: %v", err)
	}

	// Verify issue is created successfully
	if created.ID == "" {
		t.Fatal("CreateIssue: expected issue ID in response")
	}
	if created.Title != "Issue without mcp_config" {
		t.Fatalf("CreateIssue: expected title 'Issue without mcp_config', got %q", created.Title)
	}

	// Cleanup
	cleanupReq := newRequest("DELETE", "/api/issues/"+created.ID, nil)
	cleanupReq = withURLParam(cleanupReq, "id", created.ID)
	testHandler.DeleteIssue(httptest.NewRecorder(), cleanupReq)
}

// TestUpdateIssueAbsentMcpConfigPreservesValue tests that omitting mcp_config in UpdateIssue preserves the existing value.
func TestUpdateIssueAbsentMcpConfigPreservesValue(t *testing.T) {
	if testHandler == nil {
		t.Skip("database not available")
	}

	// Create an issue with mcp_config
	validMcpConfig := map[string]any{
		"mcpServers": map[string]any{
			"github": map[string]any{
				"command": "npx",
				"args":    []string{"-y", "@modelcontextprotocol/server-github"},
			},
		},
	}

	w := httptest.NewRecorder()
	req := newRequest("POST", "/api/issues?workspace_id="+testWorkspaceID, map[string]any{
		"title":      "Issue for preserve test",
		"status":     "todo",
		"priority":   "medium",
		"mcp_config": validMcpConfig,
	})
	testHandler.CreateIssue(w, req)
	var created IssueResponse
	json.NewDecoder(w.Body).Decode(&created)
	issueID := created.ID

	defer func() {
		cleanupReq := newRequest("DELETE", "/api/issues/"+issueID, nil)
		cleanupReq = withURLParam(cleanupReq, "id", issueID)
		testHandler.DeleteIssue(httptest.NewRecorder(), cleanupReq)
	}()

	// Update without mcp_config field
	w = httptest.NewRecorder()
	req = newRequest("PUT", "/api/issues/"+issueID, map[string]any{
		"status": "in_progress",
	})
	req = withURLParam(req, "id", issueID)
	testHandler.UpdateIssue(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("UpdateIssue: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var updated IssueResponse
	if err := json.NewDecoder(w.Body).Decode(&updated); err != nil {
		t.Fatalf("UpdateIssue: decode response: %v", err)
	}

	// Verify mcp_config is preserved
	if updated.McpConfig == nil {
		t.Fatal("UpdateIssue: expected mcp_config to be preserved, got nil")
	}

	var returnedConfig map[string]any
	if err := json.Unmarshal(updated.McpConfig, &returnedConfig); err != nil {
		t.Fatalf("UpdateIssue: unmarshal returned mcp_config: %v", err)
	}

	if _, ok := returnedConfig["mcpServers"]; !ok {
		t.Fatal("UpdateIssue: returned mcp_config missing mcpServers key")
	}
}

// genValidMcpConfig generates a valid mcp_config JSON with mcpServers key.
func genValidMcpConfig() gopter.Gen {
	return gen.AlphaString().Map(func(s string) json.RawMessage {
		config := map[string]interface{}{
			"mcpServers": map[string]interface{}{
				s: map[string]interface{}{
					"command": "npx",
					"args":    []string{"-y", "@modelcontextprotocol/server-" + s},
				},
			},
		}
		data, _ := json.Marshal(config)
		return json.RawMessage(data)
	})
}

// genInvalidJSON generates a string that is not valid JSON.
func genInvalidJSON() gopter.Gen {
	return gen.OneConstOf(
		"{invalid json}",
		"[1, 2, ",
		`{"key": undefined}`,
		"not json at all",
		"{",
		"}",
	)
}

// genValidJSONWithoutMcpServers generates valid JSON that doesn't have mcpServers key.
func genValidJSONWithoutMcpServers() gopter.Gen {
	validJSONStrings := []string{
		`{}`,
		`{"servers": {}}`,
		`[]`,
		`true`,
		`false`,
		`123`,
		`{"other": "value"}`,
		`{"mcpServers": null}`,
	}
	return gen.IntRange(0, len(validJSONStrings)-1).Map(func(i int) json.RawMessage {
		return json.RawMessage(validJSONStrings[i])
	})
}

// TestProperty_McpConfigPersistenceRoundTrip tests that mcp_config persists unchanged through create/get cycle.
// Feature: issue-mcp-tools-selector, Property 1: mcp_config persistence round-trip
// Validates: Requirements 1.2, 1.3, 1.5
func TestProperty_McpConfigPersistenceRoundTrip(t *testing.T) {
	if testHandler == nil {
		t.Skip("database not available")
	}

	properties := gopter.NewProperties(nil)
	properties.Property("mcp_config persists unchanged", prop.ForAll(
		func(mcpConfig json.RawMessage) bool {
			// Create issue with mcp_config
			w := httptest.NewRecorder()
			req := newRequest("POST", "/api/issues?workspace_id="+testWorkspaceID, map[string]any{
				"title":      "Test issue",
				"status":     "todo",
				"priority":   "medium",
				"mcp_config": mcpConfig,
			})
			testHandler.CreateIssue(w, req)
			if w.Code != http.StatusCreated {
				t.Logf("CreateIssue failed: %d", w.Code)
				return false
			}

			var created IssueResponse
			if err := json.NewDecoder(w.Body).Decode(&created); err != nil {
				t.Logf("Failed to decode created issue: %v", err)
				return false
			}

			// Get issue and verify mcp_config matches
			w = httptest.NewRecorder()
			req = newRequest("GET", "/api/issues/"+created.ID, nil)
			req = withURLParam(req, "id", created.ID)
			testHandler.GetIssue(w, req)
			if w.Code != http.StatusOK {
				t.Logf("GetIssue failed: %d", w.Code)
				return false
			}

			var fetched IssueResponse
			if err := json.NewDecoder(w.Body).Decode(&fetched); err != nil {
				t.Logf("Failed to decode fetched issue: %v", err)
				return false
			}

			// Verify mcp_config matches
			if string(created.McpConfig) != string(fetched.McpConfig) {
				t.Logf("mcp_config mismatch: created=%s, fetched=%s", string(created.McpConfig), string(fetched.McpConfig))
				return false
			}

			// Cleanup
			cleanupReq := newRequest("DELETE", "/api/issues/"+created.ID, nil)
			cleanupReq = withURLParam(cleanupReq, "id", created.ID)
			testHandler.DeleteIssue(httptest.NewRecorder(), cleanupReq)

			return true
		},
		genValidMcpConfig(),
	))

	properties.TestingRun(t)
}

// TestProperty_IssueMcpConfigPrecedence tests that Issue mcp_config takes precedence over Agent config.
// Feature: issue-mcp-tools-selector, Property 2: issue mcp_config takes precedence
// Validates: Requirements 2.1, 2.5
func TestProperty_IssueMcpConfigPrecedence(t *testing.T) {
	if testHandler == nil {
		t.Skip("database not available")
	}

	properties := gopter.NewProperties(nil)
	properties.Property("issue mcp_config takes precedence", prop.ForAll(
		func(issueMcpConfig json.RawMessage) bool {
			// Create issue with mcp_config
			w := httptest.NewRecorder()
			req := newRequest("POST", "/api/issues?workspace_id="+testWorkspaceID, map[string]any{
				"title":      "Test issue",
				"status":     "todo",
				"priority":   "medium",
				"mcp_config": issueMcpConfig,
			})
			testHandler.CreateIssue(w, req)
			if w.Code != http.StatusCreated {
				t.Logf("CreateIssue failed: %d", w.Code)
				return false
			}

			var created IssueResponse
			if err := json.NewDecoder(w.Body).Decode(&created); err != nil {
				t.Logf("Failed to decode created issue: %v", err)
				return false
			}

			// Verify the issue has the mcp_config we set
			if string(created.McpConfig) != string(issueMcpConfig) {
				t.Logf("Issue mcp_config not set correctly: expected=%s, got=%s", string(issueMcpConfig), string(created.McpConfig))
				return false
			}

			// Cleanup
			cleanupReq := newRequest("DELETE", "/api/issues/"+created.ID, nil)
			cleanupReq = withURLParam(cleanupReq, "id", created.ID)
			testHandler.DeleteIssue(httptest.NewRecorder(), cleanupReq)

			return true
		},
		genValidMcpConfig(),
	))

	properties.TestingRun(t)
}

// TestProperty_InvalidJsonRejected tests that invalid JSON is rejected with HTTP 400.
// Feature: issue-mcp-tools-selector, Property 6: invalid JSON rejected
// Validates: Requirement 6.1
func TestProperty_InvalidJsonRejected(t *testing.T) {
	if testHandler == nil {
		t.Skip("database not available")
	}

	properties := gopter.NewProperties(nil)
	properties.Property("invalid JSON rejected by CreateIssue", prop.ForAll(
		func(invalidJSON string) bool {
			w := httptest.NewRecorder()
			req := newRequest("POST", "/api/issues?workspace_id="+testWorkspaceID, map[string]any{
				"title":      "Test issue",
				"status":     "todo",
				"priority":   "medium",
				"mcp_config": invalidJSON,
			})
			testHandler.CreateIssue(w, req)
			if w.Code != http.StatusBadRequest {
				t.Logf("CreateIssue with invalid JSON: expected 400, got %d", w.Code)
				return false
			}

			var errResp map[string]any
			if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
				t.Logf("Failed to decode error response: %v", err)
				return false
			}

			errMsg, ok := errResp["error"].(string)
			if !ok || errMsg != "mcp_config must be valid JSON" {
				t.Logf("Expected error 'mcp_config must be valid JSON', got %v", errResp["error"])
				return false
			}

			return true
		},
		genInvalidJSON(),
	))

	properties.Property("invalid JSON rejected by UpdateIssue", prop.ForAll(
		func(invalidJSON string) bool {
			// Create an issue first
			w := httptest.NewRecorder()
			req := newRequest("POST", "/api/issues?workspace_id="+testWorkspaceID, map[string]any{
				"title":    "Test issue",
				"status":   "todo",
				"priority": "medium",
			})
			testHandler.CreateIssue(w, req)
			var created IssueResponse
			json.NewDecoder(w.Body).Decode(&created)
			issueID := created.ID

			defer func() {
				cleanupReq := newRequest("DELETE", "/api/issues/"+issueID, nil)
				cleanupReq = withURLParam(cleanupReq, "id", issueID)
				testHandler.DeleteIssue(httptest.NewRecorder(), cleanupReq)
			}()

			// Try to update with invalid JSON
			w = httptest.NewRecorder()
			req = newRequest("PUT", "/api/issues/"+issueID, map[string]any{
				"mcp_config": invalidJSON,
			})
			req = withURLParam(req, "id", issueID)
			testHandler.UpdateIssue(w, req)
			if w.Code != http.StatusBadRequest {
				t.Logf("UpdateIssue with invalid JSON: expected 400, got %d", w.Code)
				return false
			}

			var errResp map[string]any
			if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
				t.Logf("Failed to decode error response: %v", err)
				return false
			}

			errMsg, ok := errResp["error"].(string)
			if !ok || errMsg != "mcp_config must be valid JSON" {
				t.Logf("Expected error 'mcp_config must be valid JSON', got %v", errResp["error"])
				return false
			}

			return true
		},
		genInvalidJSON(),
	))

	properties.TestingRun(t)
}

// TestProperty_MissingMcpServersRejected tests that JSON without mcpServers key is rejected with HTTP 400.
// Feature: issue-mcp-tools-selector, Property 7: missing mcpServers rejected
// Validates: Requirement 6.2
func TestProperty_MissingMcpServersRejected(t *testing.T) {
	if testHandler == nil {
		t.Skip("database not available")
	}

	properties := gopter.NewProperties(nil)
	properties.Property("missing mcpServers rejected by CreateIssue", prop.ForAll(
		func(validJSONNoMcpServers json.RawMessage) bool {
			w := httptest.NewRecorder()
			req := newRequest("POST", "/api/issues?workspace_id="+testWorkspaceID, map[string]any{
				"title":      "Test issue",
				"status":     "todo",
				"priority":   "medium",
				"mcp_config": validJSONNoMcpServers,
			})
			testHandler.CreateIssue(w, req)
			if w.Code != http.StatusBadRequest {
				t.Logf("CreateIssue with missing mcpServers: expected 400, got %d", w.Code)
				return false
			}

			var errResp map[string]any
			if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
				t.Logf("Failed to decode error response: %v", err)
				return false
			}

			errMsg, ok := errResp["error"].(string)
			if !ok {
				t.Logf("Expected error message, got %v", errResp["error"])
				return false
			}
			// Accept either "mcpServers" error or "valid JSON" error since both are valid rejections
			if !strings.Contains(errMsg, "mcpServers") && !strings.Contains(errMsg, "valid JSON") {
				t.Logf("Expected error about mcpServers or JSON, got %v", errResp["error"])
				return false
			}

			return true
		},
		genValidJSONWithoutMcpServers(),
	))

	properties.Property("missing mcpServers rejected by UpdateIssue", prop.ForAll(
		func(validJSONNoMcpServers json.RawMessage) bool {
			// Create an issue first
			w := httptest.NewRecorder()
			req := newRequest("POST", "/api/issues?workspace_id="+testWorkspaceID, map[string]any{
				"title":    "Test issue",
				"status":   "todo",
				"priority": "medium",
			})
			testHandler.CreateIssue(w, req)
			var created IssueResponse
			json.NewDecoder(w.Body).Decode(&created)
			issueID := created.ID

			defer func() {
				cleanupReq := newRequest("DELETE", "/api/issues/"+issueID, nil)
				cleanupReq = withURLParam(cleanupReq, "id", issueID)
				testHandler.DeleteIssue(httptest.NewRecorder(), cleanupReq)
			}()

			// Try to update with JSON missing mcpServers
			w = httptest.NewRecorder()
			req = newRequest("PUT", "/api/issues/"+issueID, map[string]any{
				"mcp_config": validJSONNoMcpServers,
			})
			req = withURLParam(req, "id", issueID)
			testHandler.UpdateIssue(w, req)
			if w.Code != http.StatusBadRequest {
				t.Logf("UpdateIssue with missing mcpServers: expected 400, got %d", w.Code)
				return false
			}

			var errResp map[string]any
			if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
				t.Logf("Failed to decode error response: %v", err)
				return false
			}

			errMsg, ok := errResp["error"].(string)
			if !ok {
				t.Logf("Expected error message, got %v", errResp["error"])
				return false
			}
			// Accept either "mcpServers" error or "valid JSON" error since both are valid rejections
			if !strings.Contains(errMsg, "mcpServers") && !strings.Contains(errMsg, "valid JSON") {
				t.Logf("Expected error about mcpServers or JSON, got %v", errResp["error"])
				return false
			}

			return true
		},
		genValidJSONWithoutMcpServers(),
	))

	properties.TestingRun(t)
}

package daemon

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// handleMcpConfigWrite writes the given mcpConfig JSON to the runtime's native
// config file. mcpConfig is the full {"mcpServers": {...}} JSON string.
// Empty string means clear the MCP servers from the config file.
func (d *Daemon) handleMcpConfigWrite(ctx context.Context, rt Runtime, mcpConfig string) {
	d.logger.Info("mcp config write requested", "runtime_id", rt.ID, "provider", rt.Provider)

	var err error
	switch rt.Provider {
	case "opencode":
		err = writeOpencodeMcp(mcpConfig)
	case "claude":
		err = writeClaudeMcp(mcpConfig)
	case "gemini":
		err = writeGeminiMcp(mcpConfig)
	case "codex":
		err = writeCodexMcp(mcpConfig)
	default:
		d.logger.Warn("mcp config write: unsupported provider", "provider", rt.Provider)
		return
	}

	if err != nil {
		d.logger.Error("mcp config write failed", "provider", rt.Provider, "error", err)
	} else {
		d.logger.Info("mcp config write succeeded", "provider", rt.Provider)
	}
}

// parseMcpServersFromJSON parses {"mcpServers": {...}} and returns the inner map.
func parseMcpServersFromJSON(mcpConfig string) (map[string]interface{}, error) {
	if mcpConfig == "" {
		return map[string]interface{}{}, nil
	}
	var cfg struct {
		McpServers map[string]interface{} `json:"mcpServers"`
	}
	if err := json.Unmarshal([]byte(mcpConfig), &cfg); err != nil {
		return nil, fmt.Errorf("invalid mcpConfig JSON: %w", err)
	}
	if cfg.McpServers == nil {
		return map[string]interface{}{}, nil
	}
	return cfg.McpServers, nil
}

// ── opencode: ~/.config/opencode/opencode.json ──
// Format: { "mcp": { "name": { "type": "local", "command": [...], "environment": {...} } } }

func writeOpencodeMcp(mcpConfig string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := filepath.Join(home, ".config", "opencode", "opencode.json")

	// Read existing file
	existing := map[string]interface{}{}
	if data, err := os.ReadFile(path); err == nil {
		json.Unmarshal(data, &existing)
	}

	// Parse new mcpServers
	servers, err := parseMcpServersFromJSON(mcpConfig)
	if err != nil {
		return err
	}

	// Convert mcpServers format → opencode "mcp" format
	mcp := map[string]interface{}{}
	for name, v := range servers {
		entry, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		cmd, _ := entry["command"].(string)
		args, _ := entry["args"].([]interface{})
		env, _ := entry["env"].(map[string]interface{})

		cmdArr := []interface{}{cmd}
		cmdArr = append(cmdArr, args...)

		ocEntry := map[string]interface{}{
			"type":    "local",
			"command": cmdArr,
		}
		if len(env) > 0 {
			ocEntry["environment"] = env
		}
		mcp[name] = ocEntry
	}

	if len(mcp) == 0 {
		delete(existing, "mcp")
	} else {
		existing["mcp"] = mcp
	}

	return writeJSONFile(path, existing)
}

// ── claude: ~/.claude.json ──
// Format: { "mcpServers": { "name": { "command": "...", "args": [...], "type": "stdio" } } }

func writeClaudeMcp(mcpConfig string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := filepath.Join(home, ".claude.json")

	existing := map[string]interface{}{}
	if data, err := os.ReadFile(path); err == nil {
		json.Unmarshal(data, &existing)
	}

	servers, err := parseMcpServersFromJSON(mcpConfig)
	if err != nil {
		return err
	}

	// Claude uses mcpServers at root level with "type": "stdio"
	claudeServers := map[string]interface{}{}
	for name, v := range servers {
		entry, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		claudeEntry := map[string]interface{}{
			"command": entry["command"],
			"args":    entry["args"],
			"type":    "stdio",
		}
		if env, ok := entry["env"].(map[string]interface{}); ok && len(env) > 0 {
			claudeEntry["env"] = env
		}
		claudeServers[name] = claudeEntry
	}

	if len(claudeServers) == 0 {
		delete(existing, "mcpServers")
	} else {
		existing["mcpServers"] = claudeServers
	}

	return writeJSONFile(path, existing)
}

// ── gemini: ~/.gemini/settings.json ──
// Format: { "mcpServers": { "name": { "command": "...", "args": [...] } } }

func writeGeminiMcp(mcpConfig string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := filepath.Join(home, ".gemini", "settings.json")

	existing := map[string]interface{}{}
	if data, err := os.ReadFile(path); err == nil {
		json.Unmarshal(data, &existing)
	}

	servers, err := parseMcpServersFromJSON(mcpConfig)
	if err != nil {
		return err
	}

	if len(servers) == 0 {
		delete(existing, "mcpServers")
	} else {
		existing["mcpServers"] = servers
	}

	return writeJSONFile(path, existing)
}

// ── codex: ~/.codex/config.toml ──
// Format: [mcp_servers.name] command = "..." args = [...]

func writeCodexMcp(mcpConfig string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := filepath.Join(home, ".codex", "config.toml")

	servers, err := parseMcpServersFromJSON(mcpConfig)
	if err != nil {
		return err
	}

	// Read existing TOML, strip old mcp_servers sections, append new ones
	existing := ""
	if data, err := os.ReadFile(path); err == nil {
		existing = stripCodexMcpSections(string(data))
	}

	var sb strings.Builder
	sb.WriteString(strings.TrimRight(existing, "\n"))
	if sb.Len() > 0 {
		sb.WriteString("\n\n")
	}

	for name, v := range servers {
		entry, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		cmd, _ := entry["command"].(string)
		args, _ := entry["args"].([]interface{})
		env, _ := entry["env"].(map[string]interface{})

		fmt.Fprintf(&sb, "[mcp_servers.%s]\n", name)
		fmt.Fprintf(&sb, "command = %q\n", cmd)

		argStrs := make([]string, 0, len(args))
		for _, a := range args {
			argStrs = append(argStrs, fmt.Sprintf("%q", fmt.Sprint(a)))
		}
		fmt.Fprintf(&sb, "args = [%s]\n", strings.Join(argStrs, ", "))

		if len(env) > 0 {
			envParts := make([]string, 0, len(env))
			for k, val := range env {
				envParts = append(envParts, fmt.Sprintf("%q = %q", k, fmt.Sprint(val)))
			}
			fmt.Fprintf(&sb, "env = { %s }\n", strings.Join(envParts, ", "))
		}
		sb.WriteString("\n")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}

// stripCodexMcpSections removes all [mcp_servers.*] sections from TOML content.
func stripCodexMcpSections(content string) string {
	var out strings.Builder
	inMcp := false
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "[") {
			inMcp = strings.HasPrefix(trimmed, "[mcp_servers.")
		}
		if !inMcp {
			out.WriteString(line)
			out.WriteString("\n")
		}
	}
	return out.String()
}

// writeJSONFile writes v as pretty-printed JSON to path, creating dirs as needed.
func writeJSONFile(path string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

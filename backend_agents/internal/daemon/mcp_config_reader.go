package daemon

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// McpServerInfo holds the parsed info of a single MCP server from a runtime's config file.
type McpServerInfo struct {
	Name    string            `json:"name"`
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env,omitempty"`
}

// ReadRuntimeMcpServers reads the MCP server list from the runtime's native config file.
// Returns nil (no error) if the file doesn't exist or the provider is unsupported.
func ReadRuntimeMcpServers(provider string) ([]McpServerInfo, error) {
	switch provider {
	case "opencode":
		return readOpencodeMcp()
	case "claude":
		return readClaudeMcp()
	case "codex":
		return readCodexMcp()
	case "gemini":
		return readGeminiMcp()
	default:
		return nil, nil
	}
}

// opencode: ~/.config/opencode/opencode.json
// Format: { "mcp": { "name": { "type": "local", "command": ["cmd", "arg1", ...], "environment": {...} } } }
func readOpencodeMcp() ([]McpServerInfo, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, nil
	}
	path := filepath.Join(home, ".config", "opencode", "opencode.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil
	}

	var cfg struct {
		Mcp map[string]struct {
			Command     []string          `json:"command"`
			Environment map[string]string `json:"environment"`
			Enabled     *bool             `json:"enabled"`
		} `json:"mcp"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, nil
	}

	var servers []McpServerInfo
	for name, s := range cfg.Mcp {
		// Skip explicitly disabled servers
		if s.Enabled != nil && !*s.Enabled {
			continue
		}
		if len(s.Command) == 0 {
			continue
		}
		servers = append(servers, McpServerInfo{
			Name:    name,
			Command: s.Command[0],
			Args:    s.Command[1:],
			Env:     s.Environment,
		})
	}
	return servers, nil
}

// claude: ~/.claude.json
func readClaudeMcp() ([]McpServerInfo, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, nil
	}
	return readJsonMcpServers(filepath.Join(home, ".claude.json"))
}

// gemini: ~/.gemini/settings.json
func readGeminiMcp() ([]McpServerInfo, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, nil
	}
	return readJsonMcpServers(filepath.Join(home, ".gemini", "settings.json"))
}

// codex: ~/.codex/config.toml
// Parses [mcp_servers.name] sections without an external TOML library.
func readCodexMcp() ([]McpServerInfo, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, nil
	}
	path := filepath.Join(home, ".codex", "config.toml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil
	}
	return parseCodexToml(string(data)), nil
}

// parseCodexToml parses the [mcp_servers.NAME] sections from a codex config.toml.
func parseCodexToml(content string) []McpServerInfo {
	type entry struct {
		command string
		args    []string
		env     map[string]string
	}
	servers := map[string]*entry{}
	var current string

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Section header: [mcp_servers.NAME]
		if strings.HasPrefix(line, "[mcp_servers.") && strings.HasSuffix(line, "]") {
			current = line[len("[mcp_servers."):]
			current = current[:len(current)-1]
			if _, ok := servers[current]; !ok {
				servers[current] = &entry{env: map[string]string{}}
			}
			continue
		}
		// Left mcp_servers section
		if strings.HasPrefix(line, "[") {
			current = ""
			continue
		}
		if current == "" {
			continue
		}
		e := servers[current]
		if strings.HasPrefix(line, "command") {
			e.command = tomlStringValue(line)
		} else if strings.HasPrefix(line, "args") {
			e.args = tomlStringArray(line)
		} else if strings.HasPrefix(line, "env") {
			e.env = tomlInlineTable(line)
		}
	}

	var result []McpServerInfo
	for name, e := range servers {
		result = append(result, McpServerInfo{
			Name:    name,
			Command: e.command,
			Args:    e.args,
			Env:     e.env,
		})
	}
	return result
}

func tomlStringValue(line string) string {
	idx := strings.Index(line, "=")
	if idx < 0 {
		return ""
	}
	return strings.Trim(strings.TrimSpace(line[idx+1:]), `"`)
}

func tomlStringArray(line string) []string {
	idx := strings.Index(line, "[")
	end := strings.LastIndex(line, "]")
	if idx < 0 || end <= idx {
		return nil
	}
	var result []string
	for _, part := range strings.Split(line[idx+1:end], ",") {
		part = strings.Trim(strings.TrimSpace(part), `"`)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

func tomlInlineTable(line string) map[string]string {
	result := map[string]string{}
	idx := strings.Index(line, "{")
	end := strings.LastIndex(line, "}")
	if idx < 0 || end <= idx {
		return result
	}
	for _, pair := range strings.Split(line[idx+1:end], ",") {
		pair = strings.TrimSpace(pair)
		eqIdx := strings.Index(pair, "=")
		if eqIdx < 0 {
			continue
		}
		k := strings.Trim(strings.TrimSpace(pair[:eqIdx]), `"`)
		v := strings.Trim(strings.TrimSpace(pair[eqIdx+1:]), `"`)
		if k != "" {
			result[k] = v
		}
	}
	return result
}

// readJsonMcpServers reads a JSON file with { "mcpServers": { ... } } format.
func readJsonMcpServers(path string) ([]McpServerInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil
	}
	var cfg struct {
		McpServers map[string]struct {
			Command string            `json:"command"`
			Args    []string          `json:"args"`
			Env     map[string]string `json:"env"`
		} `json:"mcpServers"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, nil
	}
	var servers []McpServerInfo
	for name, s := range cfg.McpServers {
		servers = append(servers, McpServerInfo{
			Name:    name,
			Command: s.Command,
			Args:    s.Args,
			Env:     s.Env,
		})
	}
	return servers, nil
}

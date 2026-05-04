import { describe, it, expect } from "vitest";
import {
  parseMcpServerNames,
  buildMcpConfig,
  isSelectionValidForAgent,
} from "./mcp-picker";

// ============================================================================
// Tests for parseMcpServerNames
// ============================================================================

describe("parseMcpServerNames", () => {
  it("should return empty array for null input", () => {
    expect(parseMcpServerNames(null)).toEqual([]);
  });

  it("should return empty array for empty string", () => {
    expect(parseMcpServerNames("")).toEqual([]);
  });

  it("should return empty array for invalid JSON", () => {
    expect(parseMcpServerNames("not json")).toEqual([]);
    expect(parseMcpServerNames("{invalid}")).toEqual([]);
  });

  it("should return empty array for JSON without mcpServers key", () => {
    expect(parseMcpServerNames("{}")).toEqual([]);
    expect(parseMcpServerNames('{"servers": {}}')).toEqual([]);
    expect(parseMcpServerNames('{"mcpServers": null}')).toEqual([]);
    expect(parseMcpServerNames('{"mcpServers": "not an object"}')).toEqual([]);
  });

  it("should return server names from valid mcp_config", () => {
    const config = JSON.stringify({
      mcpServers: {
        github: { command: "npx", args: ["@mcp/github"] },
        slack: { command: "npx", args: ["@mcp/slack"] },
      },
    });
    const result = parseMcpServerNames(config);
    expect(result).toHaveLength(2);
    expect(result).toContain("github");
    expect(result).toContain("slack");
  });

  it("should handle single server", () => {
    const config = JSON.stringify({
      mcpServers: {
        github: { command: "npx", args: ["@mcp/github"] },
      },
    });
    expect(parseMcpServerNames(config)).toEqual(["github"]);
  });

  it("should handle empty mcpServers object", () => {
    const config = JSON.stringify({ mcpServers: {} });
    expect(parseMcpServerNames(config)).toEqual([]);
  });

  it("should handle multiple servers", () => {
    const config = JSON.stringify({
      mcpServers: {
        github: { command: "npx" },
        slack: { command: "npx" },
        database: { command: "npx" },
        api: { command: "npx" },
      },
    });
    const result = parseMcpServerNames(config);
    expect(result).toHaveLength(4);
    expect(result).toContain("github");
    expect(result).toContain("slack");
    expect(result).toContain("database");
    expect(result).toContain("api");
  });
});

// ============================================================================
// Tests for buildMcpConfig
// ============================================================================

describe("buildMcpConfig", () => {
  const agentConfig = JSON.stringify({
    mcpServers: {
      github: { command: "npx", args: ["@mcp/github"] },
      slack: { command: "npx", args: ["@mcp/slack"] },
      database: { command: "npx", args: ["@mcp/database"] },
    },
  });

  it("should return null for empty selectedServers", () => {
    expect(buildMcpConfig(agentConfig, [])).toBeNull();
  });

  it("should return null for invalid agent config", () => {
    expect(buildMcpConfig("invalid json", ["github"])).toBeNull();
    expect(buildMcpConfig("{}", ["github"])).toBeNull();
    expect(buildMcpConfig('{"servers": {}}', ["github"])).toBeNull();
  });

  it("should return null if no selected servers exist in agent config", () => {
    expect(buildMcpConfig(agentConfig, ["nonexistent"])).toBeNull();
    expect(buildMcpConfig(agentConfig, ["nonexistent1", "nonexistent2"])).toBeNull();
  });

  it("should build config with single selected server", () => {
    const result = buildMcpConfig(agentConfig, ["github"]);
    expect(result).not.toBeNull();
    const parsed = JSON.parse(result!);
    expect(parsed.mcpServers).toHaveProperty("github");
    expect(Object.keys(parsed.mcpServers)).toEqual(["github"]);
  });

  it("should build config with multiple selected servers", () => {
    const result = buildMcpConfig(agentConfig, ["github", "slack"]);
    expect(result).not.toBeNull();
    const parsed = JSON.parse(result!);
    expect(parsed.mcpServers).toHaveProperty("github");
    expect(parsed.mcpServers).toHaveProperty("slack");
    expect(Object.keys(parsed.mcpServers)).toHaveLength(2);
  });

  it("should preserve server configuration from agent config", () => {
    const result = buildMcpConfig(agentConfig, ["github"]);
    const parsed = JSON.parse(result!);
    const agentParsed = JSON.parse(agentConfig);
    expect(parsed.mcpServers.github).toEqual(
      agentParsed.mcpServers.github
    );
  });

  it("should ignore non-existent servers in selection", () => {
    const result = buildMcpConfig(agentConfig, ["github", "nonexistent", "slack"]);
    expect(result).not.toBeNull();
    const parsed = JSON.parse(result!);
    expect(Object.keys(parsed.mcpServers)).toHaveLength(2);
    expect(parsed.mcpServers).toHaveProperty("github");
    expect(parsed.mcpServers).toHaveProperty("slack");
  });

  it("should handle all servers selected", () => {
    const result = buildMcpConfig(agentConfig, ["github", "slack", "database"]);
    expect(result).not.toBeNull();
    const parsed = JSON.parse(result!);
    expect(Object.keys(parsed.mcpServers)).toHaveLength(3);
  });

  it("should return valid JSON", () => {
    const result = buildMcpConfig(agentConfig, ["github"]);
    expect(() => JSON.parse(result!)).not.toThrow();
  });
});

// ============================================================================
// Tests for isSelectionValidForAgent
// ============================================================================

describe("isSelectionValidForAgent", () => {
  const agentConfig1 = JSON.stringify({
    mcpServers: {
      github: { command: "npx" },
      slack: { command: "npx" },
    },
  });

  const agentConfig2 = JSON.stringify({
    mcpServers: {
      github: { command: "npx" },
      database: { command: "npx" },
    },
  });

  const agentConfig3 = JSON.stringify({
    mcpServers: {
      github: { command: "npx" },
      slack: { command: "npx" },
      database: { command: "npx" },
    },
  });

  it("should return true if current config is null", () => {
    expect(isSelectionValidForAgent(null, agentConfig1)).toBe(true);
  });

  it("should return true if new agent config is null", () => {
    const currentConfig = JSON.stringify({
      mcpServers: { github: { command: "npx" } },
    });
    expect(isSelectionValidForAgent(currentConfig, null)).toBe(true);
  });

  it("should return true if both configs are null", () => {
    expect(isSelectionValidForAgent(null, null)).toBe(true);
  });

  it("should return true if all selected servers exist in new agent config", () => {
    const currentConfig = JSON.stringify({
      mcpServers: { github: { command: "npx" } },
    });
    expect(isSelectionValidForAgent(currentConfig, agentConfig1)).toBe(true);
  });

  it("should return false if selected server not in new agent config", () => {
    const currentConfig = JSON.stringify({
      mcpServers: { slack: { command: "npx" } },
    });
    expect(isSelectionValidForAgent(currentConfig, agentConfig2)).toBe(false);
  });

  it("should return false if any selected server missing from new agent config", () => {
    const currentConfig = JSON.stringify({
      mcpServers: {
        github: { command: "npx" },
        slack: { command: "npx" },
      },
    });
    expect(isSelectionValidForAgent(currentConfig, agentConfig2)).toBe(false);
  });

  it("should return true if new agent has superset of selected servers", () => {
    const currentConfig = JSON.stringify({
      mcpServers: { github: { command: "npx" } },
    });
    expect(isSelectionValidForAgent(currentConfig, agentConfig3)).toBe(true);
  });

  it("should return true if current config has empty mcpServers", () => {
    const currentConfig = JSON.stringify({ mcpServers: {} });
    expect(isSelectionValidForAgent(currentConfig, agentConfig1)).toBe(true);
  });

  it("should return true if current config is invalid JSON", () => {
    expect(isSelectionValidForAgent("invalid", agentConfig1)).toBe(true);
  });

  it("should return true if new agent config is invalid JSON", () => {
    const currentConfig = JSON.stringify({
      mcpServers: { github: { command: "npx" } },
    });
    expect(isSelectionValidForAgent(currentConfig, "invalid")).toBe(true);
  });

  it("should handle multiple selected servers correctly", () => {
    const currentConfig = JSON.stringify({
      mcpServers: {
        github: { command: "npx" },
        slack: { command: "npx" },
      },
    });
    expect(isSelectionValidForAgent(currentConfig, agentConfig3)).toBe(true);
    expect(isSelectionValidForAgent(currentConfig, agentConfig2)).toBe(false);
  });
});

// ============================================================================
// Integration Tests
// ============================================================================

describe("Helper functions integration", () => {
  it("should round-trip: parse -> build -> parse", () => {
    const agentConfig = JSON.stringify({
      mcpServers: {
        github: { command: "npx", args: ["@mcp/github"] },
        slack: { command: "npx", args: ["@mcp/slack"] },
      },
    });

    const servers = parseMcpServerNames(agentConfig);
    const built = buildMcpConfig(agentConfig, servers);
    const reparsed = parseMcpServerNames(built);

    expect(reparsed).toEqual(servers);
  });

  it("should handle subset selection round-trip", () => {
    const agentConfig = JSON.stringify({
      mcpServers: {
        github: { command: "npx" },
        slack: { command: "npx" },
        database: { command: "npx" },
      },
    });

    const selected = ["github", "database"];
    const built = buildMcpConfig(agentConfig, selected);
    const reparsed = parseMcpServerNames(built);

    expect(reparsed).toEqual(selected);
  });

  it("should validate selection after agent change", () => {
    const agentConfig1 = JSON.stringify({
      mcpServers: {
        github: { command: "npx" },
        slack: { command: "npx" },
      },
    });

    const agentConfig2 = JSON.stringify({
      mcpServers: {
        github: { command: "npx" },
        database: { command: "npx" },
      },
    });

    // User selects github and slack from agent1
    const selected = ["github", "slack"];
    const issueConfig = buildMcpConfig(agentConfig1, selected);

    // Agent changes to agent2 (has github but not slack)
    const isValid = isSelectionValidForAgent(issueConfig, agentConfig2);
    expect(isValid).toBe(false);

    // If we select only github, it should be valid
    const validSelection = buildMcpConfig(agentConfig1, ["github"]);
    expect(isSelectionValidForAgent(validSelection, agentConfig2)).toBe(true);
  });
});

// ============================================================================
// Property-Based Tests (using fast-check)
// ============================================================================

import * as fc from "fast-check";

// Generator for valid mcp_config with multiple servers
const genValidMcpConfig = () => {
  return fc
    .array(fc.string({ minLength: 1, maxLength: 10 }).filter((s: string) => /^[a-z0-9]+$/.test(s)), {
      minLength: 1,
      maxLength: 5,
    })
    .map((serverNames: string[]) => {
      const mcpServers: Record<string, any> = {};
      serverNames.forEach((name: string) => {
        mcpServers[name] = {
          command: "npx",
          args: [`@mcp/${name}`],
        };
      });
      return JSON.stringify({ mcpServers });
    });
};

describe("Property-Based Tests", () => {
  // Feature: issue-mcp-tools-selector, Property 4: MCP server subset construction
  // Validates: Requirements 3.5, 3.7
  it("Property 4: MCP server subset construction - buildMcpConfig creates config with exactly selected servers", () => {
    fc.assert(
      fc.property(genValidMcpConfig(), (agentConfig) => {
        const serverNames = parseMcpServerNames(agentConfig);
        if (serverNames.length === 0) return true; // Skip if no servers

        // Test with a random subset
        const selectedServers = serverNames.slice(
          0,
          Math.max(1, Math.floor(Math.random() * serverNames.length))
        );
        const result = buildMcpConfig(agentConfig, selectedServers);
        if (result === null) return false;

        const parsed = JSON.parse(result);
        const resultServers = Object.keys(parsed.mcpServers);

        // Verify exactly the selected servers are in the result
        expect(resultServers.sort()).toEqual(selectedServers.sort());
        expect(resultServers).toHaveLength(selectedServers.length);
        return true;
      }),
      { numRuns: 20 }
    );
  });

  // Feature: issue-mcp-tools-selector, Property 5: MCP_Picker displays all agent servers
  // Validates: Requirements 3.3
  it("Property 5: MCP_Picker displays all agent servers - parseMcpServerNames returns all servers", () => {
    fc.assert(
      fc.property(genValidMcpConfig(), (agentConfig) => {
        const parsed = JSON.parse(agentConfig);
        const expectedServers = Object.keys(parsed.mcpServers);

        const result = parseMcpServerNames(agentConfig);

        // Verify all servers are returned
        expect(result.sort()).toEqual(expectedServers.sort());
        expect(result).toHaveLength(expectedServers.length);
        return true;
      }),
      { numRuns: 20 }
    );
  });

  // Feature: issue-mcp-tools-selector, Property 8: Agent change resets invalid selection
  // Validates: Requirements 3.8
  it("Property 8: Agent change resets invalid selection - isSelectionValidForAgent detects incompatible servers", () => {
    fc.assert(
      fc.property(
        fc.tuple(genValidMcpConfig(), genValidMcpConfig()).filter(
          ([config1, config2]) => config1 !== config2
        ),
        ([agentConfig1, agentConfig2]) => {
          const servers1 = parseMcpServerNames(agentConfig1);
          const servers2 = parseMcpServerNames(agentConfig2);

          if (servers1.length === 0 || servers2.length === 0) return true;

          // Test with a random subset from agent1
          const selectedServers = servers1.slice(
            0,
            Math.max(1, Math.floor(Math.random() * servers1.length))
          );

          // Build config from agent1
          const issueConfig = buildMcpConfig(agentConfig1, selectedServers);
          if (issueConfig === null) return true;

          // Check if selection is valid for agent2
          const isValid = isSelectionValidForAgent(issueConfig, agentConfig2);

          // If valid, all selected servers must exist in agent2
          if (isValid) {
            const servers2Names = parseMcpServerNames(agentConfig2);
            selectedServers.forEach((server) => {
              expect(servers2Names).toContain(server);
            });
          } else {
            // If invalid, at least one selected server must not exist in agent2
            const servers2Names = parseMcpServerNames(agentConfig2);
            const hasInvalidServer = selectedServers.some(
              (server) => !servers2Names.includes(server)
            );
            expect(hasInvalidServer).toBe(true);
          }
          return true;
        }
      ),
      { numRuns: 20 }
    );
  });
});

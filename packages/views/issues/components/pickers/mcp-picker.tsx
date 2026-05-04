"use client";

import { useState, useMemo } from "react";
import { AlertCircle } from "lucide-react";
import type { UpdateIssueRequest } from "@multica/core/types";
import { useQuery } from "@tanstack/react-query";
import { useWorkspaceId } from "@multica/core/hooks";
import { agentListOptions } from "@multica/core/workspace/queries";
import {
  PropertyPicker,
  PickerItem,
  PickerSection,
  PickerEmpty,
} from "./property-picker";

// ============================================================================
// Helper Functions (Pure, Testable)
// ============================================================================

/**
 * Parse mcp_config JSON string to extract server names.
 * Returns empty array if config is null, invalid JSON, or missing mcpServers key.
 *
 * @param mcpConfig - JSON string containing mcp_config, or null
 * @returns Array of server names
 */
export function parseMcpServerNames(mcpConfig: string | null): string[] {
  if (!mcpConfig) return [];
  try {
    const parsed = JSON.parse(mcpConfig);
    if (
      typeof parsed === "object" &&
      parsed !== null &&
      "mcpServers" in parsed &&
      typeof parsed.mcpServers === "object" &&
      parsed.mcpServers !== null
    ) {
      return Object.keys(parsed.mcpServers);
    }
    return [];
  } catch {
    return [];
  }
}

/**
 * Build mcp_config JSON containing only selected servers from agent config.
 * Returns null if selectedServers is empty or if agentMcpConfig is invalid.
 *
 * @param agentMcpConfig - Agent's mcp_config JSON string
 * @param selectedServers - Array of server names to include
 * @returns JSON string with selected servers, or null if empty selection
 */
export function buildMcpConfig(
  agentMcpConfig: string,
  selectedServers: string[]
): string | null {
  if (selectedServers.length === 0) return null;

  try {
    const agentConfig = JSON.parse(agentMcpConfig);
    if (
      typeof agentConfig !== "object" ||
      agentConfig === null ||
      !("mcpServers" in agentConfig) ||
      typeof agentConfig.mcpServers !== "object" ||
      agentConfig.mcpServers === null
    ) {
      return null;
    }

    const selectedConfig: Record<string, unknown> = {};
    for (const serverName of selectedServers) {
      if (serverName in agentConfig.mcpServers) {
        selectedConfig[serverName] = agentConfig.mcpServers[serverName];
      }
    }

    if (Object.keys(selectedConfig).length === 0) return null;

    return JSON.stringify({ mcpServers: selectedConfig });
  } catch {
    return null;
  }
}

/**
 * Check if current MCP selection is still valid with a new agent's config.
 * Returns false if any previously selected servers are not in the new agent config.
 *
 * @param currentMcpConfig - Current issue's mcp_config (JSON string or null)
 * @param newAgentMcpConfig - New agent's mcp_config (JSON string or null)
 * @returns true if selection is still valid, false if needs reset
 */
export function isSelectionValidForAgent(
  currentMcpConfig: string | null,
  newAgentMcpConfig: string | null
): boolean {
  if (!currentMcpConfig || !newAgentMcpConfig) return true;

  const currentServers = parseMcpServerNames(currentMcpConfig);
  const newAgentServers = parseMcpServerNames(newAgentMcpConfig);

  // If we can't parse either config, be lenient and return true
  if (currentServers.length === 0 || newAgentServers.length === 0) return true;

  // Check if all current servers exist in new agent config
  return currentServers.every((server) => newAgentServers.includes(server));
}

// ============================================================================
// McpPicker Component
// ============================================================================

export function McpPicker({
  mcpConfig,
  agentMcpConfig,
  onUpdate,
  triggerRender,
  align,
}: {
  mcpConfig: string | null;
  agentMcpConfig: string | null;
  onUpdate: (updates: Partial<UpdateIssueRequest>) => void;
  triggerRender?: React.ReactElement;
  align?: "start" | "center" | "end";
}) {
  const [open, setOpen] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const wsId = useWorkspaceId();
  const { isLoading, isError } = useQuery(agentListOptions(wsId));

  // Get available servers from agent config
  const availableServers = useMemo(
    () => parseMcpServerNames(agentMcpConfig),
    [agentMcpConfig]
  );

  // Get currently selected servers
  const selectedServers = useMemo(
    () => parseMcpServerNames(mcpConfig),
    [mcpConfig]
  );

  // Determine trigger label
  const triggerLabel = useMemo(() => {
    if (selectedServers.length === 0) return "No MCP tools";
    if (selectedServers.length === 1) return selectedServers[0];
    return `${selectedServers.length} tools`;
  }, [selectedServers]);

  const handleServerToggle = (serverName: string) => {
    const newSelected = selectedServers.includes(serverName)
      ? selectedServers.filter((s) => s !== serverName)
      : [...selectedServers, serverName];

    const newConfig = buildMcpConfig(agentMcpConfig || "", newSelected);
    onUpdate({ mcp_config: newConfig });
    setError(null);
  };

  const handleRetry = () => {
    setError(null);
  };

  // Empty state: no agent assigned
  if (!agentMcpConfig) {
    return (
      <PropertyPicker
        open={false}
        onOpenChange={() => {}}
        width="w-48"
        align={align}
        triggerRender={triggerRender}
        trigger={
          <span className="text-muted-foreground">Assign agent first</span>
        }
      >
        <div className="px-2 py-3 text-center text-sm text-muted-foreground">
          Assign an agent to configure MCP tools
        </div>
      </PropertyPicker>
    );
  }

  // Error state: failed to load agents
  if (isError) {
    return (
      <PropertyPicker
        open={open}
        onOpenChange={setOpen}
        width="w-48"
        align={align}
        triggerRender={triggerRender}
        trigger={
          <div className="flex items-center gap-1.5 text-destructive">
            <AlertCircle className="h-3.5 w-3.5" />
            <span>Error</span>
          </div>
        }
      >
        <div className="px-2 py-3 text-center text-sm text-muted-foreground">
          <div className="mb-2">Failed to load MCP tools</div>
          <button
            onClick={handleRetry}
            className="text-xs text-primary hover:underline"
          >
            Retry
          </button>
        </div>
      </PropertyPicker>
    );
  }

  // Loading state
  if (isLoading) {
    return (
      <PropertyPicker
        open={false}
        onOpenChange={() => {}}
        width="w-48"
        align={align}
        triggerRender={triggerRender}
        trigger={<span className="text-muted-foreground">Loading...</span>}
      >
        <div className="px-2 py-3 text-center text-sm text-muted-foreground">
          Loading MCP tools...
        </div>
      </PropertyPicker>
    );
  }

  // No servers available
  if (availableServers.length === 0) {
    return (
      <PropertyPicker
        open={false}
        onOpenChange={() => {}}
        width="w-48"
        align={align}
        triggerRender={triggerRender}
        trigger={
          <span className="text-muted-foreground">No MCP tools available</span>
        }
      >
        <div className="px-2 py-3 text-center text-sm text-muted-foreground">
          Agent has no MCP tools configured
        </div>
      </PropertyPicker>
    );
  }

  return (
    <PropertyPicker
      open={open}
      onOpenChange={setOpen}
      width="w-52"
      align={align}
      triggerRender={triggerRender}
      trigger={
        <span className="truncate">{triggerLabel}</span>
      }
    >
      {error && (
        <div className="px-2 py-2 border-b bg-destructive/10 text-xs text-destructive">
          {error}
          <button
            onClick={handleRetry}
            className="ml-2 underline hover:no-underline"
          >
            Retry
          </button>
        </div>
      )}

      {availableServers.length > 0 && (
        <PickerSection label="MCP Servers">
          {availableServers.map((serverName) => (
            <PickerItem
              key={serverName}
              selected={selectedServers.includes(serverName)}
              onClick={() => handleServerToggle(serverName)}
            >
              <span>{serverName}</span>
            </PickerItem>
          ))}
        </PickerSection>
      )}

      {availableServers.length === 0 && <PickerEmpty />}
    </PropertyPicker>
  );
}

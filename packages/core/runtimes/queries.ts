import { queryOptions } from "@tanstack/react-query";
import { api } from "../api";

export const runtimeKeys = {
  all: (wsId: string) => ["runtimes", wsId] as const,
  list: (wsId: string) => [...runtimeKeys.all(wsId), "list"] as const,
  listMine: (wsId: string) => [...runtimeKeys.all(wsId), "list", "mine"] as const,
  latestVersion: () => ["runtimes", "latestVersion"] as const,
};

export function runtimeListOptions(wsId: string, owner?: "me") {
  return queryOptions({
    queryKey: owner === "me" ? runtimeKeys.listMine(wsId) : runtimeKeys.list(wsId),
    queryFn: () => api.listRuntimes({ workspace_id: wsId, owner }),
  });
}

const GITHUB_RELEASES_URL =
  "https://api.github.com/repos/multica-ai/multica/releases/latest";

export function latestCliVersionOptions() {
  return queryOptions({
    queryKey: runtimeKeys.latestVersion(),
    queryFn: async (): Promise<string | null> => {
      try {
        const resp = await fetch(GITHUB_RELEASES_URL, {
          headers: { Accept: "application/vnd.github+json" },
        });
        if (!resp.ok) return null;
        const data = await resp.json();
        return (data.tag_name as string) ?? null;
      } catch {
        return null;
      }
    },
    staleTime: 10 * 60 * 1000, // 10 minutes
  });
}

/**
 * Build an mcp_config JSON string from a runtime's MCP servers.
 * Merges two sources (manual override wins on conflict):
 *   1. runtime.metadata.mcp_servers — detected from the runtime's config file (heartbeat)
 *   2. runtime.mcp_config — manually added via the UI (stored in DB)
 * Returns null if neither source has any servers.
 */
export function buildMcpConfigFromRuntime(runtime: import("../types").AgentRuntime | undefined): string | null {
  if (!runtime) return null;

  const mcpServers: Record<string, unknown> = {};

  // 1. Detected servers from runtime config file (via heartbeat → metadata)
  const meta = runtime.metadata as Record<string, unknown> | undefined;
  if (meta?.mcp_servers && Array.isArray(meta.mcp_servers)) {
    const servers = meta.mcp_servers as Array<{
      name: string;
      command: string;
      args?: string[];
      env?: Record<string, string>;
    }>;
    for (const s of servers) {
      if (!s.name) continue;
      mcpServers[s.name] = {
        command: s.command,
        args: s.args ?? [],
        ...(s.env && Object.keys(s.env).length > 0 ? { env: s.env } : {}),
      };
    }
  }

  // 2. Manual override from DB (runtime.mcp_config) — merges on top, overrides on conflict
  if (runtime.mcp_config) {
    try {
      const override = JSON.parse(runtime.mcp_config) as { mcpServers?: Record<string, unknown> };
      if (override?.mcpServers && typeof override.mcpServers === "object") {
        Object.assign(mcpServers, override.mcpServers);
      }
    } catch {
      // ignore invalid JSON
    }
  }

  if (Object.keys(mcpServers).length === 0) return null;
  return JSON.stringify({ mcpServers });
}

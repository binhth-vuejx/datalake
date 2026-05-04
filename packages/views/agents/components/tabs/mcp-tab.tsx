"use client";

import { useState } from "react";
import {
  Loader2,
  Save,
  Plus,
  Trash2,
} from "lucide-react";
import type { Agent } from "@multica/core/types";
import { Button } from "@multica/ui/components/ui/button";
import { Input } from "@multica/ui/components/ui/input";
import { Label } from "@multica/ui/components/ui/label";
import { Textarea } from "@multica/ui/components/ui/textarea";
import { toast } from "sonner";

let nextMcpId = 0;

interface McpServerEntry {
  id: number;
  name: string;
  command: string;
  args: string;
  env: string;
}

interface McpConfig {
  mcpServers: Record<string, {
    command: string;
    args: string[];
    env?: Record<string, string>;
  }>;
}

function mcpConfigToEntries(config: McpConfig | null): McpServerEntry[] {
  if (!config || !config.mcpServers) return [];
  return Object.entries(config.mcpServers).map(([name, server]) => ({
    id: nextMcpId++,
    name,
    command: server.command,
    args: server.args.join(" "),
    env: server.env ? JSON.stringify(server.env, null, 2) : "",
  }));
}

function entriesToMcpConfig(entries: McpServerEntry[]): McpConfig {
  const mcpServers: McpConfig["mcpServers"] = {};
  for (const entry of entries) {
    const name = entry.name.trim();
    if (name) {
      const argsArray = entry.args.trim().split(/\s+/).filter(a => a);
      let envObj: Record<string, string> | undefined;
      try {
        if (entry.env.trim()) {
          envObj = JSON.parse(entry.env);
        }
      } catch {
        // Invalid JSON, ignore env
      }
      mcpServers[name] = {
        command: entry.command,
        args: argsArray,
        ...(envObj && { env: envObj }),
      };
    }
  }
  return { mcpServers };
}

export function McpTab({
  agent,
  onSave,
}: {
  agent: Agent;
  onSave: (updates: Partial<Agent>) => Promise<void>;
}) {
  const [mcpEntries, setMcpEntries] = useState<McpServerEntry[]>(() => {
    try {
      const config = agent.mcp_config ? JSON.parse(agent.mcp_config) as McpConfig : null;
      return mcpConfigToEntries(config);
    } catch {
      return [];
    }
  });
  const [saving, setSaving] = useState(false);

  const currentConfig = entriesToMcpConfig(mcpEntries);
  const originalConfig = (() => {
    try {
      return agent.mcp_config ? JSON.parse(agent.mcp_config) as McpConfig : null;
    } catch {
      return null;
    }
  })();
  const dirty =
    JSON.stringify(currentConfig) !== JSON.stringify(originalConfig);

  const addMcpEntry = () => {
    setMcpEntries([
      ...mcpEntries,
      { id: nextMcpId++, name: "", command: "", args: "", env: "" },
    ]);
  };

  const removeMcpEntry = (index: number) => {
    setMcpEntries(mcpEntries.filter((_, i) => i !== index));
  };

  const updateMcpEntry = (
    index: number,
    field: keyof McpServerEntry,
    val: string,
  ) => {
    setMcpEntries(
      mcpEntries.map((entry, i) =>
        i === index ? { ...entry, [field]: val } : entry,
      ),
    );
  };

  const handleSave = async () => {
    const names = mcpEntries.filter((e) => e.name.trim()).map((e) => e.name.trim());
    const uniqueNames = new Set(names);
    if (uniqueNames.size < names.length) {
      toast.error("Duplicate MCP server names");
      return;
    }

    for (const entry of mcpEntries) {
      if (entry.name.trim() && !entry.command.trim()) {
        toast.error(`MCP server "${entry.name}" must have a command`);
        return;
      }
    }

    setSaving(true);
    try {
      await onSave({ mcp_config: JSON.stringify(currentConfig) });
      toast.success("MCP configuration saved");
    } catch {
      toast.error("Failed to save MCP configuration");
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="max-w-2xl space-y-4">
      <div>
        <Label className="text-xs text-muted-foreground">
          MCP Servers
        </Label>
        <p className="text-xs text-muted-foreground mt-0.5">
          Configure Model Context Protocol (MCP) servers that this agent can use.
          Each MCP server provides tools and resources for the agent.
        </p>
      </div>

      <div className="space-y-3">
        {mcpEntries.map((entry, index) => (
          <div key={entry.id} className="border rounded-lg p-3 space-y-2">
            <div className="flex items-center gap-2">
              <Input
                value={entry.name}
                onChange={(e) => updateMcpEntry(index, "name", e.target.value)}
                placeholder="Server name (e.g. rill)"
                className="flex-1 font-mono text-xs"
              />
              <button
                type="button"
                onClick={() => removeMcpEntry(index)}
                className="shrink-0 text-muted-foreground hover:text-destructive"
              >
                <Trash2 className="h-3.5 w-3.5" />
              </button>
            </div>
            <div className="flex items-center gap-2">
              <Input
                value={entry.command}
                onChange={(e) => updateMcpEntry(index, "command", e.target.value)}
                placeholder="Command (e.g. npx)"
                className="flex-1 font-mono text-xs"
              />
            </div>
            <div>
              <Input
                value={entry.args}
                onChange={(e) => updateMcpEntry(index, "args", e.target.value)}
                placeholder="Arguments (space-separated, e.g. mcp-remote http://localhost:9009/mcp)"
                className="font-mono text-xs"
              />
            </div>
            <div>
              <Label className="text-[10px] text-muted-foreground">
                Environment Variables (JSON, optional)
              </Label>
              <Textarea
                value={entry.env}
                onChange={(e) => updateMcpEntry(index, "env", e.target.value)}
                placeholder='{"KEY": "value"}'
                className="font-mono text-xs mt-1 h-16"
              />
            </div>
          </div>
        ))}
      </div>

      <Button
        type="button"
        variant="outline"
        size="sm"
        onClick={addMcpEntry}
        className="w-full gap-1 text-xs"
      >
        <Plus className="h-3 w-3" />
        Add MCP Server
      </Button>

      <Button onClick={handleSave} disabled={!dirty || saving} size="sm">
        {saving ? (
          <Loader2 className="h-3.5 w-3.5 mr-1.5 animate-spin" />
        ) : (
          <Save className="h-3.5 w-3.5 mr-1.5" />
        )}
        Save
      </Button>
    </div>
  );
}

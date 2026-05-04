"use client";

import { useState, useEffect } from "react";
import { Loader2, Save } from "lucide-react";
import type { AgentRuntime } from "@multica/core/types";
import { Button } from "@multica/ui/components/ui/button";
import { Textarea } from "@multica/ui/components/ui/textarea";
import { toast } from "sonner";

const MCP_SUPPORTED_PROVIDERS = ["claude", "codex", "opencode", "gemini"];

interface DetectedMcpServer {
  name: string;
  command: string;
  args?: string[];
  env?: Record<string, string>;
}

function getDetectedServers(runtime: AgentRuntime): DetectedMcpServer[] {
  try {
    const meta = runtime.metadata as Record<string, unknown> | undefined;
    if (!meta?.mcp_servers || !Array.isArray(meta.mcp_servers)) return [];
    return meta.mcp_servers as DetectedMcpServer[];
  } catch {
    return [];
  }
}

function StatusDot({ online }: { online: boolean }) {
  return (
    <span
      className={`inline-flex h-2 w-2 rounded-full shrink-0 ${
        online ? "bg-emerald-500" : "bg-muted-foreground/30"
      }`}
      title={online ? "Active" : "Offline"}
    />
  );
}

export function RuntimeMcpSection({
  runtime,
  canEdit,
  onSave,
  saving,
}: {
  runtime: AgentRuntime;
  canEdit: boolean;
  onSave: (mcpConfig: string | null) => Promise<void>;
  saving: boolean;
}) {
  if (!MCP_SUPPORTED_PROVIDERS.includes(runtime.provider)) return null;
  return (
    <RuntimeMcpSectionInner
      runtime={runtime}
      canEdit={canEdit}
      onSave={onSave}
      saving={saving}
    />
  );
}

function RuntimeMcpSectionInner({
  runtime,
  canEdit,
  onSave,
  saving,
}: {
  runtime: AgentRuntime;
  canEdit: boolean;
  onSave: (mcpConfig: string | null) => Promise<void>;
  saving: boolean;
}) {
  const isOnline = runtime.status === "online";
  const detectedServers = getDetectedServers(runtime);

  // JSON textarea value — initialized from runtime.mcp_config, fallback to detected servers
  const [jsonValue, setJsonValue] = useState<string>(() => {
    // If manual override exists, use it
    if (runtime.mcp_config) {
      try {
        return JSON.stringify(JSON.parse(runtime.mcp_config), null, 2);
      } catch {
        return runtime.mcp_config;
      }
    }
    // Otherwise pre-fill from detected servers (metadata.mcp_servers)
    const detected = getDetectedServers(runtime);
    if (detected.length > 0) {
      const mcpServers: Record<string, unknown> = {};
      for (const s of detected) {
        mcpServers[s.name] = {
          command: s.command,
          args: s.args ?? [],
          ...(s.env && Object.keys(s.env).length > 0 ? { env: s.env } : {}),
        };
      }
      return JSON.stringify({ mcpServers }, null, 2);
    }
    return "";
  });
  const [jsonError, setJsonError] = useState<string | null>(null);

  // Reset when runtime changes
  useEffect(() => {
    if (!runtime.mcp_config) {
      setJsonValue("");
    } else {
      try {
        setJsonValue(JSON.stringify(JSON.parse(runtime.mcp_config), null, 2));
      } catch {
        setJsonValue(runtime.mcp_config);
      }
    }
    setJsonError(null);
  }, [runtime.id, runtime.mcp_config]);

  const originalJson = runtime.mcp_config
    ? (() => { try { return JSON.stringify(JSON.parse(runtime.mcp_config)); } catch { return runtime.mcp_config; } })()
    : "";
  const currentJson = jsonValue.trim()
    ? (() => { try { return JSON.stringify(JSON.parse(jsonValue)); } catch { return null; } })()
    : "";
  // dirty = valid change from original (including clearing)
  const dirty = currentJson !== null
    ? currentJson !== originalJson
    : false;

  const handleChange = (val: string) => {
    setJsonValue(val);
    if (!val.trim()) {
      setJsonError(null);
      return;
    }
    try {
      const parsed = JSON.parse(val);
      if (!parsed.mcpServers || typeof parsed.mcpServers !== "object") {
        setJsonError('Must have a top-level "mcpServers" key');
      } else {
        setJsonError(null);
      }
    } catch {
      setJsonError("Invalid JSON");
    }
  };

  const handleSave = async () => {
    const trimmed = jsonValue.trim();
    if (!trimmed) {
      try {
        await onSave(null);
        toast.success("MCP config cleared");
      } catch {
        toast.error("Failed to save");
      }
      return;
    }
    if (jsonError) {
      toast.error(jsonError);
      return;
    }
    try {
      const compact = JSON.stringify(JSON.parse(trimmed));
      await onSave(compact);
      toast.success("MCP config saved");
    } catch {
      toast.error("Failed to save MCP config");
    }
  };

  return (
    <div className="space-y-4">
      {/* Detected servers from runtime config file */}
      {detectedServers.length > 0 && (
        <div>
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-xs font-medium text-muted-foreground">MCP Servers</h3>
            <span className={`text-[10px] font-medium px-1.5 py-0.5 rounded-full ${
              isOnline
                ? "bg-emerald-500/10 text-emerald-600"
                : "bg-muted text-muted-foreground"
            }`}>
              {isOnline ? "Active" : "Offline"}
            </span>
          </div>
          <div className="space-y-1.5">
            {detectedServers.map((server) => (
              <div
                key={server.name}
                className="flex items-center gap-2.5 rounded-md border px-3 py-2 bg-muted/20"
              >
                <StatusDot online={isOnline} />
                <span className="text-xs font-mono font-medium flex-1 truncate">
                  {server.name}
                </span>
                <span className="text-[10px] text-muted-foreground font-mono truncate max-w-[220px]">
                  {server.command}{server.args?.length ? " " + server.args.join(" ") : ""}
                </span>
              </div>
            ))}
          </div>
          <p className="text-[10px] text-muted-foreground mt-1.5">
            Detected from runtime config file. Updates every heartbeat (~30s).
          </p>
        </div>
      )}

      {/* JSON override editor */}
      {canEdit && (
        <div>
          <h3 className="text-xs font-medium text-muted-foreground mb-2">
            {detectedServers.length > 0 ? "Override MCP Config" : "MCP Config"}
          </h3>
          <Textarea
            value={jsonValue}
            onChange={(e) => handleChange(e.target.value)}
            placeholder={`{\n  "mcpServers": {\n    "my-server": {\n      "command": "npx",\n      "args": ["mcp-remote", "http://localhost:9009/mcp"]\n    }\n  }\n}`}
            className="font-mono text-xs min-h-[160px] resize-y"
            spellCheck={false}
          />
          {jsonError && (
            <p className="text-xs text-destructive mt-1">{jsonError}</p>
          )}
          <div className="flex items-center justify-between mt-2">
            <p className="text-[10px] text-muted-foreground">
              Paste the full <code className="font-mono">mcpServers</code> JSON config. Leave empty to clear.
            </p>
            <Button
              size="sm"
              onClick={handleSave}
              disabled={(!dirty && !!jsonValue.trim()) || saving || !!jsonError}
            >
              {saving ? (
                <Loader2 className="h-3.5 w-3.5 mr-1.5 animate-spin" />
              ) : (
                <Save className="h-3.5 w-3.5 mr-1.5" />
              )}
              Save
            </Button>
          </div>
        </div>
      )}

      {/* Read-only: no servers */}
      {!canEdit && detectedServers.length === 0 && (
        <div>
          <h3 className="text-xs font-medium text-muted-foreground mb-2">MCP Servers</h3>
          <p className="text-xs text-muted-foreground">No MCP servers configured.</p>
        </div>
      )}
    </div>
  );
}

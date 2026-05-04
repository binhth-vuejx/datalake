"use client";

import { Cpu } from "lucide-react";
import { Tooltip, TooltipTrigger, TooltipContent } from "@multica/ui/components/ui/tooltip";

/**
 * McpConfigBadge displays a visual indicator when an issue has MCP tools configured.
 * Shows a small, unobtrusive icon with a tooltip on hover.
 *
 * @param mcpConfig - The issue's mcp_config (JSON string or null)
 * @returns Badge component if mcp_config is not null, otherwise null
 */
export function McpConfigBadge({ mcpConfig }: { mcpConfig: string | null }) {
  if (!mcpConfig) return null;

  return (
    <Tooltip>
      <TooltipTrigger
        render={
          <div className="inline-flex items-center justify-center shrink-0 cursor-default">
            <Cpu className="h-3.5 w-3.5 text-muted-foreground hover:text-foreground transition-colors" />
          </div>
        }
      />
      <TooltipContent side="top" className="text-xs">
        Custom MCP tools configured
      </TooltipContent>
    </Tooltip>
  );
}

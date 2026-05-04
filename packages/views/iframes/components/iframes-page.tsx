"use client";

import { Plus, Frame, Trash2 } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { iframeListOptions } from "@multica/core/iframes";
import { useDeleteIframe } from "@multica/core/iframes/mutations";
import { useWorkspaceId } from "@multica/core/hooks";
import { useWorkspacePaths } from "@multica/core/paths";
import { useModalStore } from "@multica/core/modals";
import { AppLink } from "../../navigation";
import { Skeleton } from "@multica/ui/components/ui/skeleton";
import { Button } from "@multica/ui/components/ui/button";
import type { Iframe } from "@multica/core/types";
import { PageHeader } from "../../layout/page-header";

function formatRelativeDate(date: string): string {
  const diff = Date.now() - new Date(date).getTime();
  const days = Math.floor(diff / (1000 * 60 * 60 * 24));
  if (days < 1) return "Today";
  if (days === 1) return "1d ago";
  if (days < 30) return `${days}d ago`;
  const months = Math.floor(days / 30);
  return `${months}mo ago`;
}

function IframeRow({ iframe }: { iframe: Iframe }) {
  const wsPaths = useWorkspacePaths();
  const deleteIframe = useDeleteIframe();

  const handleDelete = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (confirm("Are you sure you want to delete this iframe?")) {
      deleteIframe.mutate(iframe.id);
    }
  };

  return (
    <div className="group/row flex h-11 items-center gap-2 px-5 text-sm transition-colors hover:bg-accent/40">
      {/* Icon + Name (navigates to detail) */}
      <AppLink
        href={wsPaths.iframeDetail(iframe.id)}
        className="flex min-w-0 flex-1 items-center gap-2"
      >
        <span className="shrink-0 w-[24px] text-center text-base">{iframe.icon || "🖼️"}</span>
        <span className="min-w-0 flex-1 truncate font-medium">{iframe.title}</span>
      </AppLink>

      {/* Description (truncated) */}
      <span className="min-w-0 flex-1 truncate text-muted-foreground">
        {iframe.description || "No description"}
      </span>

      {/* Iframe URL or Script indicator */}
      <span className="w-32 shrink-0 text-xs text-muted-foreground">
        {iframe.iframe_url ? "URL" : iframe.iframe_script ? "Script" : "-"}
      </span>

      {/* Delete button */}
      <button
        type="button"
        onClick={handleDelete}
        className="flex w-8 items-center justify-center shrink-0 rounded hover:bg-destructive/10 hover:text-destructive transition-colors cursor-pointer opacity-0 group-hover/row:opacity-100"
      >
        <Trash2 className="h-3.5 w-3.5" />
      </button>

      {/* Created */}
      <span className="w-20 shrink-0 text-right text-xs text-muted-foreground tabular-nums">
        {formatRelativeDate(iframe.created_at)}
      </span>
    </div>
  );
}

export function IframesPage() {
  const wsId = useWorkspaceId();
  const { data: iframes = [], isLoading } = useQuery(iframeListOptions(wsId)) as { data: Iframe[]; isLoading: boolean };
  const openCreateIframe = () => useModalStore.getState().open("create-iframe");

  return (
    <div className="flex h-full flex-col">
      {/* Header bar */}
      <PageHeader className="justify-between px-5">
        <div className="flex items-center gap-2">
          <Frame className="h-4 w-4 text-muted-foreground" />
          <h1 className="text-sm font-medium">IFRAMEs</h1>
          {!isLoading && iframes.length > 0 && (
            <span className="text-xs text-muted-foreground tabular-nums">{iframes.length}</span>
          )}
        </div>
        <Button size="sm" variant="outline" onClick={openCreateIframe}>
          <Plus className="h-3.5 w-3.5 mr-1" />
          New IFRAME
        </Button>
      </PageHeader>

      {/* Table */}
      <div className="flex-1 overflow-y-auto">
        {isLoading ? (
          <>
            <div className="sticky top-0 z-[1] flex h-8 items-center gap-2 border-b bg-muted/30 px-5">
              <span className="shrink-0 w-[24px]" />
              <Skeleton className="h-3 w-12 flex-1 max-w-[48px]" />
              <Skeleton className="h-3 w-12 flex-1 max-w-[48px]" />
              <Skeleton className="h-3 w-12 shrink-0" />
              <Skeleton className="h-3 w-8 shrink-0" />
              <Skeleton className="h-3 w-12 shrink-0" />
            </div>
            <div className="p-5 pt-1 space-y-1">
              {Array.from({ length: 4 }).map((_, i) => (
                <Skeleton key={i} className="h-11 w-full" />
              ))}
            </div>
          </>
        ) : iframes.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-24 text-muted-foreground">
            <Frame className="h-10 w-10 mb-3 opacity-30" />
            <p className="text-sm">No iframes yet</p>
            <Button size="sm" variant="outline" className="mt-3" onClick={openCreateIframe}>
              Create your first iframe
            </Button>
          </div>
        ) : (
          <>
            {/* Column headers */}
            <div className="sticky top-0 z-[1] flex h-8 items-center gap-2 border-b bg-muted/30 px-5 text-xs font-medium text-muted-foreground">
              {/* Icon spacer + Name */}
              <span className="shrink-0 w-[24px]" />
              <span className="min-w-0 flex-1">Name</span>
              <span className="min-w-0 flex-1">Description</span>
              <span className="w-32 shrink-0">Type</span>
              <span className="w-8 shrink-0" />
              <span className="w-20 text-right shrink-0">Created</span>
            </div>
            {/* Rows */}
            {iframes.map((iframe: Iframe) => (
              <IframeRow key={iframe.id} iframe={iframe} />
            ))}
          </>
        )}
      </div>
    </div>
  );
}

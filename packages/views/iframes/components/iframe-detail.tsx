"use client";

import { ArrowLeft, Edit2, Trash2, Frame, Pin, PinOff } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { iframeDetailOptions } from "@multica/core/iframes";
import { useDeleteIframe } from "@multica/core/iframes/mutations";
import { useWorkspaceId } from "@multica/core/hooks";
import { useWorkspacePaths } from "@multica/core/paths";
import { useModalStore } from "@multica/core/modals";
import { useAuthStore } from "@multica/core/auth";
import { pinListOptions } from "@multica/core/pins";
import { useCreatePin, useDeletePin } from "@multica/core/pins";
import { AppLink, useNavigation } from "../../navigation";
import { Button } from "@multica/ui/components/ui/button";
import { Skeleton } from "@multica/ui/components/ui/skeleton";
import { cn } from "@multica/ui/lib/utils";
import type { Iframe } from "@multica/core/types";

interface IframeDetailProps {
  iframeId?: string;
}

export function IframeDetail({ iframeId: propIframeId }: IframeDetailProps = {}) {
  // Support both Next.js (via URL params) and Astro (via props)
  // In Astro, the iframeId will be passed as a prop
  // In Next.js, we extract it from the pathname
  const nav = useNavigation();
  
  // Extract iframeId from pathname if not provided as prop
  // Pathname format: /[workspaceSlug]/iframes/[id]
  let iframeId: string = propIframeId || "";
  if (!iframeId && typeof window !== "undefined") {
    const pathParts = nav.pathname.split("/");
    const iframesIndex = pathParts.indexOf("iframes");
    if (iframesIndex !== -1 && iframesIndex + 1 < pathParts.length) {
      iframeId = pathParts[iframesIndex + 1] || "";
    }
  }
  const wsId = useWorkspaceId();
  const wsPaths = useWorkspacePaths();
  const openEditModal = () => useModalStore.getState().open("edit-iframe", { iframeId });
  const { user } = useAuthStore();
  const { data: iframe, isLoading } = useQuery(iframeDetailOptions(wsId, iframeId)) as { data: Iframe | undefined; isLoading: boolean };
  console.log("Iframe data:", iframe);
  const deleteIframe = useDeleteIframe();

  const { data: pinnedItems = [] } = useQuery({
    ...pinListOptions(wsId, user?.id ?? ""),
    enabled: !!user,
  });
  const isPinned = pinnedItems.some((p) => p.item_type === "iframe" && p.item_id === iframeId);
  const createPin = useCreatePin();
  const deletePinMut = useDeletePin();

  const handleDelete = () => {
    if (confirm("Are you sure you want to delete this iframe?")) {
      deleteIframe.mutate(iframeId, {
        onSuccess: () => {
          window.location.href = wsPaths.iframes();
        },
      });
    }
  };

  if (isLoading) {
    return (
      <div className="flex h-full flex-col">
        <div className="flex items-center gap-2 px-5 py-4 border-b">
          <Skeleton className="h-4 w-4" />
          <Skeleton className="h-6 w-32" />
        </div>
        <div className="flex-1 p-5">
          <Skeleton className="h-8 w-64 mb-4" />
          <Skeleton className="h-4 w-full mb-2" />
          <Skeleton className="h-4 w-3/4 mb-8" />
          <Skeleton className="h-64 w-full" />
        </div>
      </div>
    );
  }

  if (!iframe) {
    return (
      <div className="flex h-full items-center justify-center text-muted-foreground">
        Iframe not found
      </div>
    );
  }

  return (
    <div className="flex h-full flex-col">
      {/* Header */}
      <div className="flex items-center justify-between px-5 py-4 border-b">
        <div className="flex items-center gap-2">
          <AppLink href={wsPaths.iframes()}>
            <Button variant="ghost" size="sm">
              <ArrowLeft className="h-4 w-4 mr-1" />
              Back
            </Button>
          </AppLink>
          <span className="text-2xl">{iframe.icon || "🖼️"}</span>
          <h1 className="text-lg font-semibold">{iframe.title}</h1>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="ghost"
            size="icon-sm"
            className={cn("text-muted-foreground", isPinned && "text-foreground")}
            title={isPinned ? "Unpin from sidebar" : "Pin to sidebar"}
            onClick={() => {
              if (isPinned) {
                deletePinMut.mutate({ itemType: "iframe", itemId: iframeId });
              } else {
                createPin.mutate({ item_type: "iframe", item_id: iframeId });
              }
            }}
          >
            {isPinned ? <PinOff /> : <Pin />}
          </Button>
          <Button variant="outline" size="sm" onClick={openEditModal}>
            <Edit2 className="h-3.5 w-3.5 mr-1" />
            Edit
          </Button>
          <Button variant="outline" size="sm" onClick={handleDelete}>
            <Trash2 className="h-3.5 w-3.5 mr-1" />
            Delete
          </Button>
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto p-5">
        {/* Iframe Display */}
        <div className="border rounded-lg overflow-hidden bg-background">
          <div className="flex items-center gap-2 px-4 py-2 border-b bg-muted/50">
            <Frame className="h-4 w-4 text-muted-foreground" />
            <span className="text-sm font-medium">{iframe.description?.trim() || "Embedded Content"}</span>
          </div>
          <div className="aspect-video w-full">
            {iframe.iframe_url ? (
              <iframe
                src={iframe.iframe_url}
                className="w-full h-full border-0"
                title={iframe.title}
                sandbox="allow-scripts allow-same-origin allow-popups allow-forms"
              />
            ) : iframe.iframe_script ? (
              <div
                className="w-full h-full"
                dangerouslySetInnerHTML={{ __html: iframe.iframe_script }}
              />
            ) : (
              <div className="flex items-center justify-center h-full text-muted-foreground">
                No iframe content configured
              </div>
            )}
          </div>
        </div>

        {/* Metadata */}
        <div className="mt-6 grid grid-cols-2 gap-4 text-sm">
          <div>
            <span className="text-muted-foreground">Type:</span>
            <span className="ml-2 font-medium">
              {iframe.iframe_url ? "URL" : iframe.iframe_script ? "Script" : "None"}
            </span>
          </div>
          <div>
            <span className="text-muted-foreground">Created:</span>
            <span className="ml-2 font-medium">
              {new Date(iframe.created_at).toLocaleDateString()}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}

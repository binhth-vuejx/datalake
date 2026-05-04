"use client";

import { useCurrentWorkspace } from "./paths/hooks";

/**
 * Returns the current workspace UUID, or "" when called outside a workspace
 * route (e.g. shared issue page, unauthenticated views).
 *
 * Previously threw when no workspace was selected. Changed to return "" so
 * components that call this unconditionally (useActorName, mutations, etc.)
 * don't crash on public/shared pages. Callers that need a guaranteed non-empty
 * wsId should guard with `if (!wsId) return` or use `enabled: !!wsId` in queries.
 */
export function useWorkspaceId(): string {
  const ws = useCurrentWorkspace();
  return ws?.id ?? "";
}

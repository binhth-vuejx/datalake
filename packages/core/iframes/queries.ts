import { queryOptions } from "@tanstack/react-query";
import { api } from "../api";

export const iframeKeys = {
  all: (wsId: string) => ["iframes", wsId] as const,
  list: (wsId: string) => [...iframeKeys.all(wsId), "list"] as const,
  detail: (wsId: string, id: string) =>
    [...iframeKeys.all(wsId), "detail", id] as const,
};

export function iframeListOptions(wsId: string) {
  return queryOptions({
    queryKey: iframeKeys.list(wsId),
    queryFn: () => api.listIframes(),
    select: (data) => data.iframes,
  });
}

export function iframeDetailOptions(wsId: string, id: string) {
  return queryOptions({
    queryKey: iframeKeys.detail(wsId, id),
    queryFn: () => api.getIframe(id),
  });
}

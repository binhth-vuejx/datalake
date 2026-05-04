import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "../api";
import { iframeKeys } from "./queries";
import { useWorkspaceId } from "../hooks";
import type { Iframe, CreateIframeRequest, UpdateIframeRequest, ListIframesResponse } from "../types";

export function useCreateIframe() {
  const qc = useQueryClient();
  const wsId = useWorkspaceId();
  return useMutation({
    mutationFn: (data: CreateIframeRequest) => api.createIframe(data),
    onSuccess: (newIframe) => {
      qc.setQueryData<ListIframesResponse>(iframeKeys.list(wsId), (old) =>
        old && !old.iframes.some((p) => p.id === newIframe.id)
          ? { ...old, iframes: [...old.iframes, newIframe], total: old.total + 1 }
          : old,
      );
    },
    onSettled: () => {
      qc.invalidateQueries({ queryKey: iframeKeys.list(wsId) });
    },
  });
}

export function useUpdateIframe() {
  const qc = useQueryClient();
  const wsId = useWorkspaceId();
  return useMutation({
    mutationFn: ({ id, ...data }: { id: string } & UpdateIframeRequest) =>
      api.updateIframe(id, data),
    onMutate: ({ id, ...data }) => {
      qc.cancelQueries({ queryKey: iframeKeys.list(wsId) });
      const prevList = qc.getQueryData<ListIframesResponse>(iframeKeys.list(wsId));
      const prevDetail = qc.getQueryData<Iframe>(iframeKeys.detail(wsId, id));
      qc.setQueryData<ListIframesResponse>(iframeKeys.list(wsId), (old) =>
        old ? { ...old, iframes: old.iframes.map((p) => (p.id === id ? { ...p, ...data } : p)) } : old,
      );
      qc.setQueryData<Iframe>(iframeKeys.detail(wsId, id), (old) =>
        old ? { ...old, ...data } : old,
      );
      return { prevList, prevDetail, id };
    },
    onError: (_err, _vars, ctx) => {
      if (ctx?.prevList) qc.setQueryData(iframeKeys.list(wsId), ctx.prevList);
      if (ctx?.prevDetail) qc.setQueryData(iframeKeys.detail(wsId, ctx.id), ctx.prevDetail);
    },
    onSettled: (_data, _err, vars) => {
      qc.invalidateQueries({ queryKey: iframeKeys.detail(wsId, vars.id) });
      qc.invalidateQueries({ queryKey: iframeKeys.list(wsId) });
    },
  });
}

export function useDeleteIframe() {
  const qc = useQueryClient();
  const wsId = useWorkspaceId();
  return useMutation({
    mutationFn: (id: string) => api.deleteIframe(id),
    onMutate: async (id) => {
      await qc.cancelQueries({ queryKey: iframeKeys.list(wsId) });
      const prevList = qc.getQueryData<ListIframesResponse>(iframeKeys.list(wsId));
      qc.setQueryData<ListIframesResponse>(iframeKeys.list(wsId), (old) =>
        old ? { ...old, iframes: old.iframes.filter((p) => p.id !== id), total: old.total - 1 } : old,
      );
      qc.removeQueries({ queryKey: iframeKeys.detail(wsId, id) });
      return { prevList };
    },
    onError: (_err, _id, ctx) => {
      if (ctx?.prevList) qc.setQueryData(iframeKeys.list(wsId), ctx.prevList);
    },
    onSettled: () => {
      qc.invalidateQueries({ queryKey: iframeKeys.list(wsId) });
    },
  });
}

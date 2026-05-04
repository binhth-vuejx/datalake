import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "../api";
import { runtimeKeys } from "./queries";

export function useDeleteRuntime(wsId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (runtimeId: string) => api.deleteRuntime(runtimeId),
    onSettled: () => {
      qc.invalidateQueries({ queryKey: runtimeKeys.all(wsId) });
    },
  });
}

export function useUpdateRuntimeMcpConfig(wsId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ runtimeId, mcpConfig }: { runtimeId: string; mcpConfig: string | null }) =>
      api.updateRuntimeMcpConfig(runtimeId, mcpConfig),
    onSettled: () => {
      qc.invalidateQueries({ queryKey: runtimeKeys.all(wsId) });
    },
  });
}

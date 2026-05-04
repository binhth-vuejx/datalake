/**
 * Sales Domain Mutations
 * TanStack Query mutations for sales CRUD operations
 */

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "../api";
import { saleKeys } from "./queries";
import type { CreateSaleRequest, UpdateSaleRequest, Sale } from "./types";

/**
 * Hook to create a new sale
 * Includes optimistic updates and cache invalidation
 */
export function useCreateSale() {
  const qc = useQueryClient();

  return useMutation({
    mutationFn: async (data: CreateSaleRequest) => {
      const response = await api.createSale(data);
      if (!response.success) {
        throw new Error(response.error || "Failed to create sale");
      }
      return response.data;
    },
    onSuccess: (newSale) => {
      if (newSale) {
        // Invalidate the list to refetch with new sale
        qc.invalidateQueries({ queryKey: saleKeys.all() });
      }
    },
  });
}

/**
 * Hook to update an existing sale
 * Includes optimistic updates with rollback on error
 */
export function useUpdateSale() {
  const qc = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, ...data }: { id: number } & UpdateSaleRequest) => {
      const response = await api.updateSale(id, data);
      if (!response.success) {
        throw new Error(response.error || "Failed to update sale");
      }
      return response.data;
    },
    onMutate: async ({ id, ...data }) => {
      // Cancel any outgoing refetches for this sale
      await qc.cancelQueries({ queryKey: saleKeys.detail(id) });

      // Snapshot the previous value
      const previousSale = qc.getQueryData<Sale>(saleKeys.detail(id));

      // Optimistically update the cache
      if (previousSale) {
        qc.setQueryData<Sale>(saleKeys.detail(id), {
          ...previousSale,
          ...data,
          updated_at: new Date().toISOString(),
        });
      }

      return { previousSale };
    },
    onError: (_err, { id }, context) => {
      // Rollback on error
      if (context?.previousSale) {
        qc.setQueryData(saleKeys.detail(id), context.previousSale);
      }
    },
    onSettled: (_data, _err, { id }) => {
      // Invalidate queries to refetch fresh data
      qc.invalidateQueries({ queryKey: saleKeys.detail(id) });
      qc.invalidateQueries({ queryKey: saleKeys.all() });
    },
  });
}

/**
 * Hook to delete a sale
 * Includes cache invalidation
 */
export function useDeleteSale() {
  const qc = useQueryClient();

  return useMutation({
    mutationFn: async (id: number) => {
      const response = await api.deleteSale(id);
      if (!response.success) {
        throw new Error(response.error || "Failed to delete sale");
      }
      return id;
    },
    onSuccess: (deletedId) => {
      // Remove from cache
      qc.removeQueries({ queryKey: saleKeys.detail(deletedId) });
      // Invalidate list to refetch
      qc.invalidateQueries({ queryKey: saleKeys.all() });
    },
  });
}

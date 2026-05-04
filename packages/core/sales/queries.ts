/**
 * Sales Domain Queries
 * TanStack Query hooks for sales data fetching and caching
 */

import { useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "../api";
import type { ListSalesParams } from "./types";

/**
 * Query key factory for sales queries
 * Follows TanStack Query best practices for key organization
 */
export const saleKeys = {
  all: () => ["sales"] as const,
  list: (params?: ListSalesParams) => [...saleKeys.all(), "list", params] as const,
  detail: (id: number) => [...saleKeys.all(), "detail", id] as const,
  byCustomer: (customerId: number, params?: ListSalesParams) =>
    [...saleKeys.all(), "byCustomer", customerId, params] as const,
};

/**
 * Hook to fetch list of sales with pagination
 */
export function useListSales(params?: ListSalesParams) {
  return useQuery({
    queryKey: saleKeys.list(params),
    queryFn: async () => {
      const response = await api.listSales(params);
      if (!response.success) {
        throw new Error(response.error || "Failed to fetch sales");
      }
      return response;
    },
  });
}

/**
 * Hook to fetch a single sale by ID
 */
export function useGetSale(id: number) {
  return useQuery({
    queryKey: saleKeys.detail(id),
    queryFn: async () => {
      return api.getSale(id);
    },
  });
}

/**
 * Hook to fetch sales for a specific customer
 */
export function useGetSalesByCustomer(customerId: number, params?: ListSalesParams) {
  return useQuery({
    queryKey: saleKeys.byCustomer(customerId, params),
    queryFn: async () => {
      const response = await api.getSalesByCustomer(customerId, params);
      if (!response.success) {
        throw new Error(response.error || "Failed to fetch customer sales");
      }
      return response;
    },
  });
}

/**
 * Helper to get the query client for manual cache updates
 */
export function useQueryClientForSales() {
  return useQueryClient();
}

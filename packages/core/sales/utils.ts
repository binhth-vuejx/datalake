/**
 * Sales Domain Utilities
 * Pure functions for sales data manipulation
 */

import type { Sale, SaleStatus } from './types'

/**
 * Sort sales by creation date in descending order (newest first)
 */
export function sortSalesByDate(sales: Sale[]): Sale[] {
  return [...sales].sort(
    (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  )
}

/**
 * Filter sales by status
 */
export function filterSalesByStatus(sales: Sale[], status: SaleStatus): Sale[] {
  return sales.filter((sale) => sale.status === status)
}

/**
 * Calculate total amount across all sales
 */
export function calculateTotalAmount(sales: Sale[]): number {
  return sales.reduce((sum, sale) => sum + sale.amount, 0)
}

/**
 * Group sales by status
 */
export function groupSalesByStatus(sales: Sale[]): Record<SaleStatus, Sale[]> {
  return sales.reduce(
    (acc, sale) => {
      const status = sale.status as SaleStatus
      acc[status] = [...(acc[status] ?? []), sale]
      return acc
    },
    {} as Record<SaleStatus, Sale[]>
  )
}

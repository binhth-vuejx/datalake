/**
 * Sales Domain Types
 * Defines core types and constants for the Sales domain
 */

export type SaleStatus = 'pending' | 'completed' | 'cancelled'

export const STATUS_LABELS: Record<SaleStatus, string> = {
  pending: 'Pending',
  completed: 'Completed',
  cancelled: 'Cancelled',
}

export const STATUS_COLORS: Record<SaleStatus, string> = {
  pending: 'yellow',
  completed: 'green',
  cancelled: 'red',
}

export interface Sale {
  id: number
  customer_id: number
  amount: number
  status: SaleStatus
  created_by: string
  created_at: string
  updated_at: string
}

export interface CreateSaleRequest {
  customer_id: number
  amount: number
}

export interface UpdateSaleRequest {
  amount: number
  status: SaleStatus
}

export interface ListSalesParams {
  limit?: number
  offset?: number
}

export interface ListSalesResponse {
  success: boolean
  data?: Sale[]
  total?: number
  error?: string
}

export interface CreateSaleResponse {
  success: boolean
  data?: Sale
  error?: string
}

export interface UpdateSaleResponse {
  success: boolean
  data?: Sale
  error?: string
}

export interface DeleteSaleResponse {
  success: boolean
  error?: string
}

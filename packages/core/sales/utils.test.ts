/**
 * Sales Domain Utilities Tests
 */

import { describe, it, expect } from "vitest";
import {
  sortSalesByDate,
  filterSalesByStatus,
  calculateTotalAmount,
  groupSalesByStatus,
} from "./utils";
import type { Sale } from "./types";

const mockSales: Sale[] = [
  {
    id: 1,
    customer_id: 100,
    amount: 150.5,
    status: "pending",
    created_by: "user1",
    created_at: "2024-01-15T10:00:00Z",
    updated_at: "2024-01-15T10:00:00Z",
  },
  {
    id: 2,
    customer_id: 101,
    amount: 200.0,
    status: "completed",
    created_by: "user2",
    created_at: "2024-01-16T10:00:00Z",
    updated_at: "2024-01-16T10:00:00Z",
  },
  {
    id: 3,
    customer_id: 102,
    amount: 75.25,
    status: "pending",
    created_by: "user1",
    created_at: "2024-01-14T10:00:00Z",
    updated_at: "2024-01-14T10:00:00Z",
  },
  {
    id: 4,
    customer_id: 103,
    amount: 300.0,
    status: "cancelled",
    created_by: "user3",
    created_at: "2024-01-17T10:00:00Z",
    updated_at: "2024-01-17T10:00:00Z",
  },
];

describe("Sales Utilities", () => {
  describe("sortSalesByDate", () => {
    it("should sort sales by created_at in descending order (newest first)", () => {
      const sorted = sortSalesByDate(mockSales);
      expect(sorted[0].id).toBe(4); // 2024-01-17
      expect(sorted[1].id).toBe(2); // 2024-01-16
      expect(sorted[2].id).toBe(1); // 2024-01-15
      expect(sorted[3].id).toBe(3); // 2024-01-14
    });

    it("should not mutate the original array", () => {
      const original = [...mockSales];
      sortSalesByDate(mockSales);
      expect(mockSales).toEqual(original);
    });
  });

  describe("filterSalesByStatus", () => {
    it("should filter sales by status", () => {
      const pending = filterSalesByStatus(mockSales, "pending");
      expect(pending).toHaveLength(2);
      expect(pending.every((s) => s.status === "pending")).toBe(true);
    });

    it("should return empty array when no sales match status", () => {
      const result = filterSalesByStatus(mockSales, "completed");
      expect(result).toHaveLength(1);
    });
  });

  describe("calculateTotalAmount", () => {
    it("should calculate total amount across all sales", () => {
      const total = calculateTotalAmount(mockSales);
      expect(total).toBe(150.5 + 200.0 + 75.25 + 300.0);
    });

    it("should return 0 for empty array", () => {
      const total = calculateTotalAmount([]);
      expect(total).toBe(0);
    });
  });

  describe("groupSalesByStatus", () => {
    it("should group sales by status", () => {
      const grouped = groupSalesByStatus(mockSales);
      expect(grouped.pending).toHaveLength(2);
      expect(grouped.completed).toHaveLength(1);
      expect(grouped.cancelled).toHaveLength(1);
    });

    it("should have correct sales in each group", () => {
      const grouped = groupSalesByStatus(mockSales);
      expect(grouped.pending.map((s) => s.id)).toEqual([1, 3]);
      expect(grouped.completed.map((s) => s.id)).toEqual([2]);
      expect(grouped.cancelled.map((s) => s.id)).toEqual([4]);
    });
  });
});

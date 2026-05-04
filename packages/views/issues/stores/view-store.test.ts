import { describe, it, expect } from "vitest";
import * as fc from "fast-check";
import { createStore } from "zustand/vanilla";
import { viewStoreSlice, type IssueViewState } from "@multica/core/issues/stores/view-store";

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

/** Create a fresh vanilla store backed by viewStoreSlice (no persist). */
function makeStore() {
  return createStore<IssueViewState>()(viewStoreSlice);
}

// ---------------------------------------------------------------------------
// Unit tests — baseline sanity checks
// ---------------------------------------------------------------------------

describe("IssueViewState – expandedParentIds", () => {
  it("initializes expandedParentIds as an empty Set", () => {
    const store = makeStore();
    expect(store.getState().expandedParentIds.size).toBe(0);
  });

  it("toggleExpandedParent adds an id when absent", () => {
    const store = makeStore();
    store.getState().toggleExpandedParent("parent-1");
    expect(store.getState().expandedParentIds.has("parent-1")).toBe(true);
  });

  it("toggleExpandedParent removes an id when present", () => {
    const store = makeStore();
    store.getState().toggleExpandedParent("parent-1");
    store.getState().toggleExpandedParent("parent-1");
    expect(store.getState().expandedParentIds.has("parent-1")).toBe(false);
  });
});

// ---------------------------------------------------------------------------
// Property-based tests (task 3.1)
// ---------------------------------------------------------------------------

describe("IssueViewState – property-based tests", () => {
  /**
   * Property 4: Toggle expand/collapse is a round-trip.
   *
   * For any parent issue id, calling toggleExpandedParent(id) twice SHALL
   * leave expandedParentIds in the same state as before either call was made.
   *
   * Validates: Requirements 3.3, 3.5
   */
  it("Property 4: toggleExpandedParent twice is a round-trip", () => {
    fc.assert(
      fc.property(
        // Arbitrary: a string id and an arbitrary initial set of expanded ids
        fc.string({ minLength: 1 }),
        fc.array(fc.string({ minLength: 1 }), { minLength: 0, maxLength: 10 }),
        (id, initialIds) => {
          const store = makeStore();

          // Seed the store with the initial expanded ids
          for (const seed of initialIds) {
            store.getState().toggleExpandedParent(seed);
          }

          // Snapshot the state before the round-trip
          const before = new Set(store.getState().expandedParentIds);

          // Perform the round-trip: toggle twice
          store.getState().toggleExpandedParent(id);
          store.getState().toggleExpandedParent(id);

          // State must be identical to before
          const after = store.getState().expandedParentIds;

          if (before.size !== after.size) return false;
          for (const entry of before) {
            if (!after.has(entry)) return false;
          }
          return true;
        },
      ),
      { numRuns: 100 },
    );
  });

  /**
   * Property 5: Toggle adds then removes.
   *
   * For any parent issue id not currently in expandedParentIds, calling
   * toggleExpandedParent(id) SHALL add it; calling it again SHALL remove it.
   *
   * Validates: Requirements 3.2, 3.3, 3.5
   */
  it("Property 5: toggleExpandedParent adds when absent, removes when present", () => {
    fc.assert(
      fc.property(
        fc.string({ minLength: 1 }),
        (id) => {
          const store = makeStore();

          // Precondition: id is not in the set (fresh store, empty set)
          expect(store.getState().expandedParentIds.has(id)).toBe(false);

          // First toggle: must add the id
          store.getState().toggleExpandedParent(id);
          if (!store.getState().expandedParentIds.has(id)) return false;

          // Second toggle: must remove the id
          store.getState().toggleExpandedParent(id);
          if (store.getState().expandedParentIds.has(id)) return false;

          return true;
        },
      ),
      { numRuns: 100 },
    );
  });
});

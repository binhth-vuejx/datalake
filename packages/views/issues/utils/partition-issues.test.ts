import { describe, it, expect } from "vitest";
import * as fc from "fast-check";
import type { Issue } from "@multica/core/types";
import { partitionIssues } from "./partition-issues";

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function makeIssue(overrides: Partial<Issue> = {}): Issue {
  return {
    id: "i-1",
    workspace_id: "ws-1",
    number: 1,
    identifier: "MUL-1",
    title: "Test",
    description: null,
    status: "todo",
    priority: "medium",
    assignee_type: null,
    assignee_id: null,
    creator_type: "member",
    creator_id: "u-1",
    parent_issue_id: null,
    project_id: null,
    position: 0,
    due_date: null,
    mcp_config: null,
    created_at: "2025-01-01T00:00:00Z",
    updated_at: "2025-01-01T00:00:00Z",
    ...overrides,
  };
}

// ---------------------------------------------------------------------------
// Unit tests (task 2.3)
// ---------------------------------------------------------------------------

describe("partitionIssues", () => {
  it("returns all issues as top-level when none have parents", () => {
    const issues = [
      makeIssue({ id: "a", parent_issue_id: null }),
      makeIssue({ id: "b", parent_issue_id: null }),
      makeIssue({ id: "c", parent_issue_id: null }),
    ];
    const { topLevelIssues, subIssuesByParent } = partitionIssues(issues);
    expect(topLevelIssues).toHaveLength(3);
    expect(subIssuesByParent.size).toBe(0);
  });

  it("treats orphaned sub-issues (parent not in list) as top-level", () => {
    const issues = [
      makeIssue({ id: "a", parent_issue_id: null }),
      makeIssue({ id: "b", parent_issue_id: "missing-parent" }),
    ];
    const { topLevelIssues, subIssuesByParent } = partitionIssues(issues);
    expect(topLevelIssues.map((i) => i.id)).toEqual(["a", "b"]);
    expect(subIssuesByParent.size).toBe(0);
  });

  it("handles 3-level hierarchy: grandchild goes under parent (not grandparent)", () => {
    // A (top-level) → B (sub of A) → C (sub of B)
    const a = makeIssue({ id: "A", parent_issue_id: null });
    const b = makeIssue({ id: "B", parent_issue_id: "A" });
    const c = makeIssue({ id: "C", parent_issue_id: "B" });
    const { topLevelIssues, subIssuesByParent } = partitionIssues([a, b, c]);

    // A is top-level; B is a sub-issue of A; C is a sub-issue of B
    expect(topLevelIssues.map((i) => i.id)).toEqual(["A"]);
    expect(subIssuesByParent.get("A")?.map((i) => i.id)).toEqual(["B"]);
    expect(subIssuesByParent.get("B")?.map((i) => i.id)).toEqual(["C"]);
  });

  it("returns empty results for an empty array", () => {
    const { topLevelIssues, subIssuesByParent } = partitionIssues([]);
    expect(topLevelIssues).toHaveLength(0);
    expect(subIssuesByParent.size).toBe(0);
  });

  it("correctly partitions a mixed list", () => {
    const parent = makeIssue({ id: "p1", parent_issue_id: null });
    const child1 = makeIssue({ id: "c1", parent_issue_id: "p1" });
    const child2 = makeIssue({ id: "c2", parent_issue_id: "p1" });
    const standalone = makeIssue({ id: "s1", parent_issue_id: null });

    const { topLevelIssues, subIssuesByParent } = partitionIssues([parent, child1, child2, standalone]);

    expect(topLevelIssues.map((i) => i.id)).toEqual(["p1", "s1"]);
    expect(subIssuesByParent.get("p1")?.map((i) => i.id)).toEqual(["c1", "c2"]);
  });
});

// ---------------------------------------------------------------------------
// Property-based tests (task 2.4)
// ---------------------------------------------------------------------------

/**
 * Arbitrary that generates a non-empty pool of unique IDs, then builds
 * Issue objects where each issue's parent_issue_id is either null or one
 * of the IDs in the pool (allowing self-references and cross-references).
 */
// issueArrayArb is kept for reference but simpleIssueArrayArb is used in tests
// const issueArrayArb = ...

/**
 * Simpler arbitrary: array of issues with unique sequential IDs and
 * random parent_issue_id drawn from the same id pool or null.
 */
const simpleIssueArrayArb = fc
  .integer({ min: 0, max: 20 })
  .chain((n) => {
    const ids = Array.from({ length: n }, (_, i) => `issue-${i}`);
    if (n === 0) return fc.constant([] as Issue[]);
    return fc
      .array(fc.option(fc.constantFrom(...ids), { nil: null }), {
        minLength: n,
        maxLength: n,
      })
      .map((parentIds) =>
        ids.map((id, i) =>
          makeIssue({
            id,
            // avoid self-reference to keep semantics clean
            parent_issue_id: parentIds[i] === id ? null : parentIds[i],
          }),
        ),
      );
  });

describe("partitionIssues – property-based tests", () => {
  /**
   * Property 1: Every issue in topLevelIssues has a parent_issue_id that is
   * null or not present in the id set of the input array.
   *
   * Validates: Requirements 1.1, 1.2
   */
  it("Property 1: top-level issues have no parent in the same list", () => {
    fc.assert(
      fc.property(simpleIssueArrayArb, (issues) => {
        const idSet = new Set(issues.map((i) => i.id));
        const { topLevelIssues } = partitionIssues(issues);
        for (const issue of topLevelIssues) {
          const parentPresent =
            issue.parent_issue_id !== null && idSet.has(issue.parent_issue_id);
          if (parentPresent) return false;
        }
        return true;
      }),
      { numRuns: 100 },
    );
  });

  /**
   * Property 2: Every issue whose parent_issue_id matches another issue's id
   * appears in subIssuesByParent under that parent id, and NOT in topLevelIssues.
   *
   * Validates: Requirements 1.1, 2.1
   */
  it("Property 2: sub-issue map is complete and consistent", () => {
    fc.assert(
      fc.property(simpleIssueArrayArb, (issues) => {
        const idSet = new Set(issues.map((i) => i.id));
        const { topLevelIssues, subIssuesByParent } = partitionIssues(issues);
        const topLevelIds = new Set(topLevelIssues.map((i) => i.id));

        for (const issue of issues) {
          if (issue.parent_issue_id && idSet.has(issue.parent_issue_id)) {
            // Must be in subIssuesByParent under its parent
            const siblings = subIssuesByParent.get(issue.parent_issue_id) ?? [];
            if (!siblings.some((s) => s.id === issue.id)) return false;
            // Must NOT be in topLevelIssues
            if (topLevelIds.has(issue.id)) return false;
          }
        }
        return true;
      }),
      { numRuns: 100 },
    );
  });

  /**
   * Property 3: The union of topLevelIssues and all values in subIssuesByParent
   * equals the original array — no issues dropped or duplicated.
   *
   * Validates: Requirements 1.1, 1.3
   */
  it("Property 3: partition is a lossless split (no drops, no duplicates)", () => {
    fc.assert(
      fc.property(simpleIssueArrayArb, (issues) => {
        const { topLevelIssues, subIssuesByParent } = partitionIssues(issues);

        const allPartitioned: Issue[] = [
          ...topLevelIssues,
          ...Array.from(subIssuesByParent.values()).flat(),
        ];

        // Same count
        if (allPartitioned.length !== issues.length) return false;

        // Same ids (no duplicates, no drops)
        const originalIds = issues.map((i) => i.id).sort();
        const partitionedIds = allPartitioned.map((i) => i.id).sort();
        return JSON.stringify(originalIds) === JSON.stringify(partitionedIds);
      }),
      { numRuns: 100 },
    );
  });
});

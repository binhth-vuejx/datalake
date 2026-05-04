import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, cleanup } from "@testing-library/react";
import * as fc from "fast-check";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { Issue } from "@multica/core/types";

// ---------------------------------------------------------------------------
// Mocks
// ---------------------------------------------------------------------------

vi.mock("@multica/core/hooks", () => ({
  useWorkspaceId: () => "ws-1",
}));

vi.mock("@multica/core/paths", async () => {
  const actual = await vi.importActual<typeof import("@multica/core/paths")>(
    "@multica/core/paths",
  );
  return {
    ...actual,
    useCurrentWorkspace: () => ({ id: "ws-1", name: "Test WS", slug: "test" }),
    useWorkspacePaths: () => actual.paths.workspace("test"),
  };
});

vi.mock("../../navigation", () => ({
  AppLink: ({ children, href, ...props }: any) => (
    <a href={href} {...props}>
      {children}
    </a>
  ),
}));

vi.mock("@multica/core/api", () => ({
  api: {
    listProjects: () => Promise.resolve([]),
  },
  getApi: () => ({ listProjects: () => Promise.resolve([]) }),
  setApiInstance: vi.fn(),
}));

vi.mock("@multica/core/projects/queries", () => ({
  projectListOptions: () => ({ queryKey: ["projects"], queryFn: () => Promise.resolve([]) }),
}));

vi.mock("@multica/core/issues/config", () => ({
  PRIORITY_CONFIG: {
    urgent: { label: "Urgent", bars: 4, color: "text-destructive" },
    high: { label: "High", bars: 3, color: "text-warning" },
    medium: { label: "Medium", bars: 2, color: "text-warning" },
    low: { label: "Low", bars: 1, color: "text-info" },
    none: { label: "No priority", bars: 0, color: "text-muted-foreground" },
  },
}));

// Selection store — mutable so tests can inspect calls
const mockToggle = vi.fn();
const mockSelectedIds = new Set<string>();

vi.mock("@multica/core/issues/stores/selection-store", () => ({
  useIssueSelectionStore: (selector?: any) => {
    const state = { selectedIds: mockSelectedIds, toggle: mockToggle };
    return selector ? selector(state) : state;
  },
}));

// View store — mutable so tests can control expandedParentIds
const mockToggleExpandedParent = vi.fn();
let mockExpandedParentIds = new Set<string>();

const mockCardProperties = {
  priority: true,
  description: true,
  assignee: false,
  dueDate: false,
  project: false,
  childProgress: true,
};

vi.mock("@multica/core/issues/stores/view-store-context", () => ({
  useViewStore: (selector?: any) => {
    const state = {
      expandedParentIds: mockExpandedParentIds,
      toggleExpandedParent: mockToggleExpandedParent,
      cardProperties: mockCardProperties,
    };
    return selector ? selector(state) : state;
  },
}));

// ---------------------------------------------------------------------------
// Test data helpers
// ---------------------------------------------------------------------------

const issueBase = {
  workspace_id: "ws-1",
  project_id: null,
  position: 0,
  assignee_type: null,
  assignee_id: null,
  creator_type: "member" as const,
  creator_id: "user-1",
  due_date: null,
  mcp_config: null,
  created_at: "2026-01-01T00:00:00Z",
  updated_at: "2026-01-01T00:00:00Z",
  description: null,
};

function makeIssue(overrides: Partial<Issue> & Pick<Issue, "id" | "identifier" | "title">): Issue {
  return {
    ...issueBase,
    number: 1,
    status: "todo",
    priority: "medium",
    parent_issue_id: null,
    ...overrides,
  };
}

const parentIssue = makeIssue({ id: "parent-1", identifier: "TES-1", title: "Parent Issue" });
const subIssue1 = makeIssue({ id: "sub-1", identifier: "TES-2", title: "Sub Issue 1", parent_issue_id: "parent-1" });
const subIssue2 = makeIssue({ id: "sub-2", identifier: "TES-3", title: "Sub Issue 2", parent_issue_id: "parent-1" });

// ---------------------------------------------------------------------------
// Import component under test (after mocks)
// ---------------------------------------------------------------------------

import { TreeRow } from "./tree-row";

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function renderTreeRow(props: Partial<Parameters<typeof TreeRow>[0]> = {}) {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={qc}>
      <TreeRow
        issue={parentIssue}
        subIssues={[]}
        childProgressMap={new Map()}
        {...props}
      />
    </QueryClientProvider>,
  );
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

describe("TreeRow", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockExpandedParentIds = new Set<string>();
    mockSelectedIds.clear();
  });

  // 4.3 — chevron present when sub-issues exist
  it("renders chevron button when subIssues.length > 0", () => {
    renderTreeRow({ subIssues: [subIssue1] });
    const chevron = screen.getByRole("button", { name: /expand sub-issues for TES-1/i });
    expect(chevron).toBeInTheDocument();
  });

  // 4.3 — chevron absent when no sub-issues
  it("does not render chevron button when subIssues is empty", () => {
    renderTreeRow({ subIssues: [] });
    expect(screen.queryByRole("button", { name: /expand sub-issues/i })).not.toBeInTheDocument();
  });

  // 4.4 — count badge value
  it("renders count badge with correct sub-issue count", () => {
    renderTreeRow({ subIssues: [subIssue1, subIssue2] });
    const badge = screen.getByTestId("sub-issue-count");
    expect(badge).toHaveTextContent("2");
  });

  // 4.4 — count badge absent when no sub-issues
  it("does not render count badge when subIssues is empty", () => {
    renderTreeRow({ subIssues: [] });
    expect(screen.queryByTestId("sub-issue-count")).not.toBeInTheDocument();
  });

  // 4.3 — clicking chevron calls toggleExpandedParent
  it("calls toggleExpandedParent with issue.id when chevron is clicked", () => {
    renderTreeRow({ subIssues: [subIssue1] });
    const chevron = screen.getByRole("button", { name: /expand sub-issues for TES-1/i });
    fireEvent.click(chevron);
    expect(mockToggleExpandedParent).toHaveBeenCalledWith("parent-1");
  });

  // 4.5 — Enter key triggers toggle
  it("calls toggleExpandedParent when Enter is pressed on chevron", () => {
    renderTreeRow({ subIssues: [subIssue1] });
    const chevron = screen.getByRole("button", { name: /expand sub-issues for TES-1/i });
    fireEvent.keyDown(chevron, { key: "Enter" });
    expect(mockToggleExpandedParent).toHaveBeenCalledWith("parent-1");
  });

  // 4.5 — Space key triggers toggle
  it("calls toggleExpandedParent when Space is pressed on chevron", () => {
    renderTreeRow({ subIssues: [subIssue1] });
    const chevron = screen.getByRole("button", { name: /expand sub-issues for TES-1/i });
    fireEvent.keyDown(chevron, { key: " " });
    expect(mockToggleExpandedParent).toHaveBeenCalledWith("parent-1");
  });

  // 4.5 — other keys do not trigger toggle
  it("does not call toggleExpandedParent for other keys", () => {
    renderTreeRow({ subIssues: [subIssue1] });
    const chevron = screen.getByRole("button", { name: /expand sub-issues for TES-1/i });
    fireEvent.keyDown(chevron, { key: "Tab" });
    expect(mockToggleExpandedParent).not.toHaveBeenCalled();
  });

  // 4.6 — sub-rows hidden when collapsed
  it("does not render sub-issue rows when parent is collapsed", () => {
    mockExpandedParentIds = new Set<string>(); // collapsed
    renderTreeRow({ subIssues: [subIssue1, subIssue2] });
    expect(screen.queryByText("Sub Issue 1")).not.toBeInTheDocument();
    expect(screen.queryByText("Sub Issue 2")).not.toBeInTheDocument();
  });

  // 4.6 — sub-rows rendered when expanded
  it("renders sub-issue rows when parent is expanded", () => {
    mockExpandedParentIds = new Set(["parent-1"]);
    renderTreeRow({ subIssues: [subIssue1, subIssue2] });
    expect(screen.getByText("Sub Issue 1")).toBeInTheDocument();
    expect(screen.getByText("Sub Issue 2")).toBeInTheDocument();
  });

  // 4.6 — indentation class pl-8 applied to sub-issue container
  it("wraps sub-issue rows in a pl-8 container when expanded", () => {
    mockExpandedParentIds = new Set(["parent-1"]);
    const { container } = renderTreeRow({ subIssues: [subIssue1] });
    const indentedContainer = container.querySelector(".pl-8");
    expect(indentedContainer).toBeInTheDocument();
  });

  // 4.7 — childProgressMap passed to sub-issue ListRow
  it("passes childProgress from childProgressMap to each sub-issue row", () => {
    mockExpandedParentIds = new Set(["parent-1"]);
    const childProgressMap = new Map([["sub-1", { done: 2, total: 5 }]]);
    renderTreeRow({ subIssues: [subIssue1], childProgressMap });
    // The progress badge renders "done/total" text
    expect(screen.getByText("2/5")).toBeInTheDocument();
  });

  // aria-expanded reflects expand state (collapsed)
  it("sets aria-expanded=false on chevron when collapsed", () => {
    mockExpandedParentIds = new Set<string>();
    renderTreeRow({ subIssues: [subIssue1] });
    const chevron = screen.getByRole("button", { name: /expand sub-issues for TES-1/i });
    expect(chevron).toHaveAttribute("aria-expanded", "false");
  });

  // aria-expanded reflects expand state (expanded)
  it("sets aria-expanded=true on chevron when expanded", () => {
    mockExpandedParentIds = new Set(["parent-1"]);
    renderTreeRow({ subIssues: [subIssue1] });
    const chevron = screen.getByRole("button", { name: /expand sub-issues for TES-1/i });
    expect(chevron).toHaveAttribute("aria-expanded", "true");
  });

  // aria-label contains issue identifier
  it("chevron aria-label contains the issue identifier", () => {
    renderTreeRow({ subIssues: [subIssue1] });
    expect(
      screen.getByRole("button", { name: "Expand sub-issues for TES-1" }),
    ).toBeInTheDocument();
  });

  // 4.8 — sub-issue checkbox selection
  it("clicking sub-issue checkbox calls toggle with sub-issue id", () => {
    mockExpandedParentIds = new Set(["parent-1"]);
    renderTreeRow({ subIssues: [subIssue1] });
    // The checkbox is hidden by default (group-hover), but we can query by role
    const checkboxes = screen.getAllByRole("checkbox");
    // First checkbox belongs to parent row, second to sub-issue row
    const subCheckbox = checkboxes[1]!;
    fireEvent.click(subCheckbox);
    expect(mockToggle).toHaveBeenCalledWith("sub-1");
  });

  // parent row is always rendered
  it("always renders the parent issue title", () => {
    renderTreeRow({ subIssues: [] });
    expect(screen.getByText("Parent Issue")).toBeInTheDocument();
  });

  // ---------------------------------------------------------------------------
  // Task 7.2 — Lazy fetch / render gate (Requirement 6.2 & 6.3)
  // ---------------------------------------------------------------------------

  // 6.2 — collapsed parent: sub-issue rows are NOT rendered (no unnecessary work)
  it("does not render sub-issue rows when parent is collapsed (lazy render gate)", () => {
    // expandedParentIds does NOT contain "parent-1" → collapsed
    mockExpandedParentIds = new Set<string>();
    renderTreeRow({ subIssues: [subIssue1, subIssue2] });
    // Sub-issue titles must be absent — confirming the render is gated on isExpanded
    expect(screen.queryByText("Sub Issue 1")).not.toBeInTheDocument();
    expect(screen.queryByText("Sub Issue 2")).not.toBeInTheDocument();
  });

  // 6.3 — expanded parent: sub-issue rows ARE rendered
  it("renders sub-issue rows when parent is expanded (lazy render gate)", () => {
    // expandedParentIds contains "parent-1" → expanded
    mockExpandedParentIds = new Set(["parent-1"]);
    renderTreeRow({ subIssues: [subIssue1, subIssue2] });
    expect(screen.getByText("Sub Issue 1")).toBeInTheDocument();
    expect(screen.getByText("Sub Issue 2")).toBeInTheDocument();
  });

  // 6.2 — toggling from expanded back to collapsed hides sub-issue rows
  it("hides sub-issue rows after parent is collapsed again", () => {
    // Start expanded
    mockExpandedParentIds = new Set(["parent-1"]);
    const { rerender } = renderTreeRow({ subIssues: [subIssue1] });
    expect(screen.getByText("Sub Issue 1")).toBeInTheDocument();

    // Simulate collapse by updating the store state
    mockExpandedParentIds = new Set<string>();
    const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    rerender(
      <QueryClientProvider client={qc}>
        <TreeRow issue={parentIssue} subIssues={[subIssue1]} childProgressMap={new Map()} />
      </QueryClientProvider>,
    );
    expect(screen.queryByText("Sub Issue 1")).not.toBeInTheDocument();
  });
});

// ---------------------------------------------------------------------------
// Property-based tests
// ---------------------------------------------------------------------------

describe("TreeRow — property-based tests", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockExpandedParentIds = new Set<string>();
    mockSelectedIds.clear();
  });

  /**
   * Property 6: Sub-issue count badge equals sub-issue array length
   * Validates: Requirements 5.1, 5.3
   */
  it("Property 6: count badge equals subIssues.length for any non-empty sub-issues array", () => {
    fc.assert(
      fc.property(
        // Generate 1–20 sub-issues (non-empty so badge is rendered)
        fc.array(
          fc.record({
            id: fc.uuid(),
            identifier: fc.string({ minLength: 1, maxLength: 10 }),
            title: fc.string({ minLength: 1, maxLength: 50 }),
          }),
          { minLength: 1, maxLength: 20 },
        ),
        (rawSubIssues) => {
          const subIssues = rawSubIssues.map((s) =>
            makeIssue({ id: s.id, identifier: s.identifier, title: s.title, parent_issue_id: "parent-1" }),
          );

          renderTreeRow({ subIssues });

          const badge = screen.getByTestId("sub-issue-count");
          expect(badge).toHaveTextContent(String(subIssues.length));

          cleanup();
        },
      ),
      { numRuns: 100 },
    );
  });

  /**
   * Property 7: aria-expanded reflects expand state
   * Validates: Requirements 7.3
   */
  it("Property 7: aria-expanded is 'true' iff issue.id is in expandedParentIds", () => {
    fc.assert(
      fc.property(
        fc.boolean(),
        (isExpanded) => {
          mockExpandedParentIds = isExpanded ? new Set(["parent-1"]) : new Set<string>();

          renderTreeRow({ subIssues: [subIssue1] });

          const chevron = screen.getByRole("button", { name: /expand sub-issues for TES-1/i });
          expect(chevron).toHaveAttribute("aria-expanded", String(isExpanded));

          cleanup();
        },
      ),
      { numRuns: 100 },
    );
  });

  /**
   * Property 8: aria-label contains issue identifier
   * Validates: Requirements 7.4
   */
  it("Property 8: chevron aria-label contains the issue identifier", () => {
    fc.assert(
      fc.property(
        // Generate identifier strings that are non-empty and contain no regex-special chars
        fc.stringMatching(/^[A-Z]{1,5}-[0-9]{1,4}$/),
        (identifier) => {
          const issue = makeIssue({ id: "parent-1", identifier, title: "Some Issue" });

          renderTreeRow({ issue, subIssues: [subIssue1] });

          const chevron = screen.getByRole("button", { name: new RegExp(identifier, "i") });
          expect(chevron.getAttribute("aria-label")).toContain(identifier);

          cleanup();
        },
      ),
      { numRuns: 100 },
    );
  });
});

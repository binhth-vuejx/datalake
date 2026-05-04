"use client";

import { useEffect } from "react";
import { ChevronRight } from "lucide-react";
import type { Issue } from "@multica/core/types";
import { useViewStore } from "@multica/core/issues/stores/view-store-context";
import { ListRow, type ChildProgress } from "./list-row";

export interface TreeRowProps {
  issue: Issue;
  subIssues: Issue[];
  childProgress?: ChildProgress;
  childProgressMap: Map<string, ChildProgress>;
}

export function TreeRow({ issue, subIssues, childProgress, childProgressMap }: TreeRowProps) {
  const expandedParentIds = useViewStore((s) => s.expandedParentIds);
  const toggleExpandedParent = useViewStore((s) => s.toggleExpandedParent);

  const hasSubIssues = subIssues.length > 0;
  const isExpanded = expandedParentIds.has(issue.id);

  // Auto-expand on first mount if this issue has sub-issues
  useEffect(() => {
    if (hasSubIssues && !expandedParentIds.has(issue.id)) {
      toggleExpandedParent(issue.id);
    }
    // Only run on mount (issue.id is stable)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [issue.id, hasSubIssues]);

  function handleChevronClick() {
    toggleExpandedParent(issue.id);
  }

  function handleChevronKeyDown(event: React.KeyboardEvent<HTMLButtonElement>) {
    if (event.key === "Enter" || event.key === " ") {
      event.preventDefault();
      toggleExpandedParent(issue.id);
    }
  }

  return (
    <div>
      <div className="relative">
        <ListRow issue={issue} childProgress={childProgress} />
        {hasSubIssues && (
          // Position after identifier column:
          // px-4 (16) + checkbox w-4 (16) + gap-2 (8) + AppLink gap-2 (8) + identifier w-16 (64) = 112px
          <div className="absolute right-0 top-0 flex h-9 items-center gap-1 pr-[3rem]">
            <button
              type="button"
              aria-expanded={isExpanded}
              aria-label={`Expand sub-issues for ${issue.identifier}`}
              onClick={handleChevronClick}
              onKeyDown={handleChevronKeyDown}
              className="flex items-center justify-center rounded p-0.5 text-muted-foreground transition-colors hover:bg-accent hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            >
              <ChevronRight
                className={`size-3 transition-transform ${isExpanded ? "rotate-90" : ""}`}
              />
            </button>
            <span
              data-testid="sub-issue-count"
              className="inline-flex items-center rounded-full bg-muted/60 px-1.5 py-0.5 text-[11px] font-medium tabular-nums text-muted-foreground"
            >
              {subIssues.length}
            </span>
          </div>
        )}
      </div>
      {hasSubIssues && isExpanded && (
        <div className="pl-8">
          {subIssues.map((sub) => (
            <ListRow key={sub.id} issue={sub} childProgress={childProgressMap.get(sub.id)} />
          ))}
        </div>
      )}
    </div>
  );
}

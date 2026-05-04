import type { Issue } from "@multica/core/types";

export interface PartitionResult {
  topLevelIssues: Issue[];
  subIssuesByParent: Map<string, Issue[]>;
}

/**
 * O(n) partition of issues into top-level issues and a sub-issues map.
 *
 * An issue is considered top-level if its parent_issue_id is null or
 * its parent is not present in the same array.
 */
export function partitionIssues(issues: Issue[]): PartitionResult {
  const issueIdSet = new Set(issues.map((i) => i.id));
  const subIssuesByParent = new Map<string, Issue[]>();
  const topLevelIssues: Issue[] = [];

  for (const issue of issues) {
    if (issue.parent_issue_id && issueIdSet.has(issue.parent_issue_id)) {
      const siblings = subIssuesByParent.get(issue.parent_issue_id) ?? [];
      siblings.push(issue);
      subIssuesByParent.set(issue.parent_issue_id, siblings);
    } else {
      topLevelIssues.push(issue);
    }
  }

  return { topLevelIssues, subIssuesByParent };
}

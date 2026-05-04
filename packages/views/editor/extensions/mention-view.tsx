"use client";

import { NodeViewWrapper } from "@tiptap/react";
import type { NodeViewProps } from "@tiptap/react";
import { useQuery } from "@tanstack/react-query";
import { issueListOptions, issueDetailOptions } from "@multica/core/issues/queries";
import { useCurrentWorkspace, useWorkspacePathsSafe } from "@multica/core/paths";
import { useNavigation } from "../../navigation";
import { StatusIcon } from "../../issues/components/status-icon";

export function MentionView({ node }: NodeViewProps) {
  const { type, id, label } = node.attrs;

  if (type === "issue") {
    return (
      <NodeViewWrapper as="span" className="inline">
        <IssueMention issueId={id} fallbackLabel={label} />
      </NodeViewWrapper>
    );
  }

  return (
    <NodeViewWrapper as="span" className="inline">
      <span className="mention">@{label ?? id}</span>
    </NodeViewWrapper>
  );
}

function IssueMention({
  issueId,
  fallbackLabel,
}: {
  issueId: string;
  fallbackLabel?: string;
}) {
  const wsId = useCurrentWorkspace()?.id ?? "";
  const p = useWorkspacePathsSafe();
  const issuePath = p?.issueDetail(issueId) ?? "#";

  const { data: issues = [] } = useQuery({
    ...issueListOptions(wsId),
    enabled: !!wsId,
  });
  const { push, openInNewTab } = useNavigation();
  const listIssue = issues.find((i) => i.id === issueId);

  const { data: detailIssue } = useQuery({
    ...issueDetailOptions(wsId, issueId),
    enabled: !!wsId && !listIssue,
  });

  const issue = listIssue ?? detailIssue;

  const tabTitle = issue ? `${issue.identifier}: ${issue.title}` : undefined;
  const handleClick = (e: React.MouseEvent) => {
    if (!p) return; // no navigation in shared view
    e.preventDefault();
    e.stopPropagation();
    if (e.metaKey || e.ctrlKey || e.shiftKey) {
      if (openInNewTab) {
        openInNewTab(issuePath, tabTitle);
      }
      return;
    }
    push(issuePath);
  };

  const cardClass =
    "issue-mention inline-flex items-center gap-1.5 rounded-md border mx-0.5 px-2 py-0.5 text-xs hover:bg-accent transition-colors cursor-pointer max-w-72";

  if (!issue) {
    return (
      <a href={issuePath} onClick={handleClick} className={cardClass}>
        <span className="font-medium text-muted-foreground">
          {fallbackLabel ?? issueId.slice(0, 8)}
        </span>
      </a>
    );
  }

  return (
    <a href={issuePath} onClick={handleClick} className={cardClass}>
      <StatusIcon status={issue.status} className="h-3.5 w-3.5 shrink-0" />
      <span className="font-medium text-muted-foreground shrink-0">{issue.identifier}</span>
      <span className="text-foreground truncate">{issue.title}</span>
    </a>
  );
}

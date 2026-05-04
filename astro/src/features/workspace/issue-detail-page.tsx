import { AstroProviders } from "@/components/astro-providers";
import { IssueDetail } from "@multica/views/issues/components";
import { DashboardLayout } from "@multica/views/layout";
import { useEffect } from "react";
import { setCurrentWorkspace } from "@multica/core/platform/workspace-storage";

interface IssueDetailPageProps {
  issueId: string;
  workspaceSlug?: string;
}

export default function IssueDetailPageWrapper({ issueId, workspaceSlug }: IssueDetailPageProps) {
  useEffect(() => {
    if (workspaceSlug) setCurrentWorkspace(workspaceSlug, null);
  }, [workspaceSlug]);

  return (
    <AstroProviders workspaceSlug={workspaceSlug}>
      <DashboardLayout>
        <IssueDetail issueId={issueId} />
      </DashboardLayout>
    </AstroProviders>
  );
}

import { AstroProviders } from "@/components/astro-providers";
import { IssuesPage } from "@multica/views/issues/components";
import { DashboardLayout } from "@multica/views/layout";
import { useEffect } from "react";
import { setCurrentWorkspace } from "@multica/core/platform/workspace-storage";

interface IssuesPageWrapperProps {
  workspaceSlug?: string;
}

export default function IssuesPageWrapper({ workspaceSlug }: IssuesPageWrapperProps) {
  useEffect(() => {
    if (workspaceSlug) setCurrentWorkspace(workspaceSlug, null);
  }, [workspaceSlug]);

  return (
    <AstroProviders workspaceSlug={workspaceSlug}>
      <DashboardLayout>
        <IssuesPage />
      </DashboardLayout>
    </AstroProviders>
  );
}

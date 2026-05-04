import { AstroProviders } from "@/components/astro-providers";
import { MyIssuesPage } from "@multica/views/my-issues";
import { DashboardLayout } from "@multica/views/layout";
import { useEffect } from "react";
import { setCurrentWorkspace } from "@multica/core/platform/workspace-storage";

interface MyIssuesPageWrapperProps {
  workspaceSlug?: string;
}

export default function MyIssuesPageWrapper({ workspaceSlug }: MyIssuesPageWrapperProps) {
  useEffect(() => {
    if (workspaceSlug) setCurrentWorkspace(workspaceSlug, null);
  }, [workspaceSlug]);

  return (
    <AstroProviders workspaceSlug={workspaceSlug}>
      <DashboardLayout>
        <MyIssuesPage />
      </DashboardLayout>
    </AstroProviders>
  );
}

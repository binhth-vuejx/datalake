import { AstroProviders } from "@/components/astro-providers";
import { AgentsPage } from "@multica/views/agents";
import { DashboardLayout } from "@multica/views/layout";
import { useEffect } from "react";
import { setCurrentWorkspace } from "@multica/core/platform/workspace-storage";

interface AgentsPageWrapperProps {
  workspaceSlug?: string;
}

export default function AgentsPageWrapper({ workspaceSlug }: AgentsPageWrapperProps) {
  useEffect(() => {
    if (workspaceSlug) setCurrentWorkspace(workspaceSlug, null);
  }, [workspaceSlug]);

  return (
    <AstroProviders workspaceSlug={workspaceSlug}>
      <DashboardLayout>
        <AgentsPage />
      </DashboardLayout>
    </AstroProviders>
  );
}

import { AstroProviders } from "@/components/astro-providers";
import { RuntimesPage } from "@multica/views/runtimes";
import { DashboardLayout } from "@multica/views/layout";
import { useEffect } from "react";
import { setCurrentWorkspace } from "@multica/core/platform/workspace-storage";

interface RuntimesPageWrapperProps {
  workspaceSlug?: string;
}

export default function RuntimesPageWrapper({ workspaceSlug }: RuntimesPageWrapperProps) {
  useEffect(() => {
    if (workspaceSlug) setCurrentWorkspace(workspaceSlug, null);
  }, [workspaceSlug]);

  return (
    <AstroProviders workspaceSlug={workspaceSlug}>
      <DashboardLayout>
        <RuntimesPage />
      </DashboardLayout>
    </AstroProviders>
  );
}

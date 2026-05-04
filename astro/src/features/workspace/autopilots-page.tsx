import { AstroProviders } from "@/components/astro-providers";
import { AutopilotsPage } from "@multica/views/autopilots/components";
import { DashboardLayout } from "@multica/views/layout";
import { useEffect } from "react";
import { setCurrentWorkspace } from "@multica/core/platform/workspace-storage";

interface AutopilotsPageWrapperProps {
  workspaceSlug?: string;
}

export default function AutopilotsPageWrapper({ workspaceSlug }: AutopilotsPageWrapperProps) {
  useEffect(() => {
    if (workspaceSlug) setCurrentWorkspace(workspaceSlug, null);
  }, [workspaceSlug]);

  return (
    <AstroProviders workspaceSlug={workspaceSlug}>
      <DashboardLayout>
        <AutopilotsPage />
      </DashboardLayout>
    </AstroProviders>
  );
}

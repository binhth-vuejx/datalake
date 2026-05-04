import { AstroProviders } from "@/components/astro-providers";
import { AutopilotDetailPage } from "@multica/views/autopilots/components";
import { DashboardLayout } from "@multica/views/layout";
import { useEffect } from "react";
import { setCurrentWorkspace } from "@multica/core/platform/workspace-storage";

interface AutopilotDetailPageProps {
  autopilotId: string;
  workspaceSlug?: string;
}

export default function AutopilotDetailPageWrapper({ autopilotId, workspaceSlug }: AutopilotDetailPageProps) {
  useEffect(() => {
    if (workspaceSlug) setCurrentWorkspace(workspaceSlug, null);
  }, [workspaceSlug]);

  return (
    <AstroProviders workspaceSlug={workspaceSlug}>
      <DashboardLayout>
        <AutopilotDetailPage autopilotId={autopilotId} />
      </DashboardLayout>
    </AstroProviders>
  );
}

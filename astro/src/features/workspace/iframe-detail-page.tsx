import { AstroProviders } from "@/components/astro-providers";
import { IframeDetail } from "@multica/views/iframes/components";
import { DashboardLayout } from "@multica/views/layout";
import { useEffect } from "react";
import { setCurrentWorkspace } from "@multica/core/platform/workspace-storage";

interface IframeDetailPageProps {
  workspaceSlug?: string;
}

export default function IframeDetailPageWrapper({ workspaceSlug }: IframeDetailPageProps) {
  useEffect(() => {
    if (workspaceSlug) setCurrentWorkspace(workspaceSlug, null);
  }, [workspaceSlug]);

  return (
    <AstroProviders workspaceSlug={workspaceSlug}>
      <DashboardLayout>
        <IframeDetail />
      </DashboardLayout>
    </AstroProviders>
  );
}

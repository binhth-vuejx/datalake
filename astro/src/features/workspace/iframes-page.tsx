import { AstroProviders } from "@/components/astro-providers";
import { IframesPage } from "@multica/views/iframes/components";
import { DashboardLayout } from "@multica/views/layout";
import { useEffect } from "react";
import { setCurrentWorkspace } from "@multica/core/platform/workspace-storage";

interface IframesPageWrapperProps {
  workspaceSlug?: string;
}

export default function IframesPageWrapper({ workspaceSlug }: IframesPageWrapperProps) {
  useEffect(() => {
    if (workspaceSlug) setCurrentWorkspace(workspaceSlug, null);
  }, [workspaceSlug]);

  return (
    <AstroProviders workspaceSlug={workspaceSlug}>
      <DashboardLayout>
        <IframesPage />
      </DashboardLayout>
    </AstroProviders>
  );
}

import { AstroProviders } from "@/components/astro-providers";
import { SettingsPage } from "@multica/views/settings";
import { DashboardLayout } from "@multica/views/layout";
import { useEffect } from "react";
import { setCurrentWorkspace } from "@multica/core/platform/workspace-storage";

interface SettingsPageWrapperProps {
  workspaceSlug?: string;
}

export default function SettingsPageWrapper({ workspaceSlug }: SettingsPageWrapperProps) {
  useEffect(() => {
    if (workspaceSlug) setCurrentWorkspace(workspaceSlug, null);
  }, [workspaceSlug]);

  return (
    <AstroProviders workspaceSlug={workspaceSlug}>
      <DashboardLayout>
        <SettingsPage />
      </DashboardLayout>
    </AstroProviders>
  );
}

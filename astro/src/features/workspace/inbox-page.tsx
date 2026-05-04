import { AstroProviders } from "@/components/astro-providers";
import { InboxPage } from "@multica/views/inbox";
import { DashboardLayout } from "@multica/views/layout";
import { useEffect } from "react";
import { setCurrentWorkspace } from "@multica/core/platform/workspace-storage";

interface InboxPageWrapperProps {
  workspaceSlug?: string;
}

export default function InboxPageWrapper({ workspaceSlug }: InboxPageWrapperProps) {
  useEffect(() => {
    if (workspaceSlug) setCurrentWorkspace(workspaceSlug, null);
  }, [workspaceSlug]);

  return (
    <AstroProviders workspaceSlug={workspaceSlug}>
      <DashboardLayout>
        <InboxPage />
      </DashboardLayout>
    </AstroProviders>
  );
}

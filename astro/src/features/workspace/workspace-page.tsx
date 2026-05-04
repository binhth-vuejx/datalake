import { AstroProviders } from "@/components/astro-providers";
import { DashboardLayout } from "@multica/views/layout";
import { useEffect } from "react";
import { setCurrentWorkspace } from "@multica/core/platform/workspace-storage";

interface WorkspacePageProps {
  workspaceSlug: string;
}

export default function WorkspacePage({ workspaceSlug }: WorkspacePageProps) {
  useEffect(() => {
    if (workspaceSlug) {
      setCurrentWorkspace(workspaceSlug, null);
    }
  }, [workspaceSlug]);

  return (
    <AstroProviders workspaceSlug={workspaceSlug}>
      <DashboardLayout>
        <div className="flex items-center justify-center h-full">
          <p className="text-muted-foreground">Workspace: {workspaceSlug}</p>
        </div>
      </DashboardLayout>
    </AstroProviders>
  );
}

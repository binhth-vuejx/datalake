import { AstroProviders } from "@/components/astro-providers";
import { ProjectsPage } from "@multica/views/projects/components";
import { DashboardLayout } from "@multica/views/layout";
import { useEffect } from "react";
import { setCurrentWorkspace } from "@multica/core/platform/workspace-storage";

interface ProjectsPageWrapperProps {
  workspaceSlug?: string;
}

export default function ProjectsPageWrapper({ workspaceSlug }: ProjectsPageWrapperProps) {
  useEffect(() => {
    if (workspaceSlug) setCurrentWorkspace(workspaceSlug, null);
  }, [workspaceSlug]);

  return (
    <AstroProviders workspaceSlug={workspaceSlug}>
      <DashboardLayout>
        <ProjectsPage />
      </DashboardLayout>
    </AstroProviders>
  );
}

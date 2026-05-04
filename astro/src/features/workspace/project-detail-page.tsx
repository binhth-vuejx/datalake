import { AstroProviders } from "@/components/astro-providers";
import { ProjectDetail } from "@multica/views/projects/components";
import { DashboardLayout } from "@multica/views/layout";
import { useEffect } from "react";
import { setCurrentWorkspace } from "@multica/core/platform/workspace-storage";

interface ProjectDetailPageProps {
  projectId: string;
  workspaceSlug?: string;
}

export default function ProjectDetailPageWrapper({ projectId, workspaceSlug }: ProjectDetailPageProps) {
  useEffect(() => {
    if (workspaceSlug) setCurrentWorkspace(workspaceSlug, null);
  }, [workspaceSlug]);

  return (
    <AstroProviders workspaceSlug={workspaceSlug}>
      <DashboardLayout>
        <ProjectDetail projectId={projectId} />
      </DashboardLayout>
    </AstroProviders>
  );
}

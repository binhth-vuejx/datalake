import { AstroProviders } from "@/components/astro-providers";
import { SkillsPage } from "@multica/views/skills";
import { DashboardLayout } from "@multica/views/layout";
import { useEffect } from "react";
import { setCurrentWorkspace } from "@multica/core/platform/workspace-storage";

interface SkillsPageWrapperProps {
  workspaceSlug?: string;
}

export default function SkillsPageWrapper({ workspaceSlug }: SkillsPageWrapperProps) {
  useEffect(() => {
    if (workspaceSlug) setCurrentWorkspace(workspaceSlug, null);
  }, [workspaceSlug]);

  return (
    <AstroProviders workspaceSlug={workspaceSlug}>
      <DashboardLayout>
        <SkillsPage />
      </DashboardLayout>
    </AstroProviders>
  );
}

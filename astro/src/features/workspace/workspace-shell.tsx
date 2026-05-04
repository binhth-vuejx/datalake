import { AstroProviders } from "@/components/astro-providers";
import { WorkspaceRouter } from "@/components/workspace-router";

interface WorkspaceShellProps {
  workspaceSlug?: string;
}

export default function WorkspaceShell({ workspaceSlug }: WorkspaceShellProps) {
  return (
    <AstroProviders workspaceSlug={workspaceSlug}>
      <WorkspaceRouter />
    </AstroProviders>
  );
}

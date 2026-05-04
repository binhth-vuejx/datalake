import { useState, useEffect } from "react";
import { DashboardLayout } from "@multica/views/layout";
import { IssuesPage } from "@multica/views/issues/components";
import { IssueDetail } from "@multica/views/issues/components";
import { AgentsPage } from "@multica/views/agents";
import { AutopilotsPage, AutopilotDetailPage } from "@multica/views/autopilots/components";
import { IframesPage, IframeDetail } from "@multica/views/iframes/components";
import { InboxPage } from "@multica/views/inbox";
import { MyIssuesPage } from "@multica/views/my-issues";
import { ProjectsPage, ProjectDetail } from "@multica/views/projects/components";
import { RuntimesPage } from "@multica/views/runtimes";
import { SettingsPage } from "@multica/views/settings";
import { SkillsPage } from "@multica/views/skills";

function parseRoute(pathname: string) {
  // /agents/{workspaceSlug}/{section}/{id?}
  const parts = pathname.split("/").filter(Boolean);
  // parts[0] = "agents", parts[1] = workspaceSlug, parts[2] = section, parts[3] = id
  const isAgentsRoute = parts[0] === "agents";
  const slug = isAgentsRoute ? (parts[1] ?? null) : (parts[0] ?? null);
  const section = isAgentsRoute ? (parts[2] ?? null) : (parts[1] ?? null);
  const id = isAgentsRoute ? (parts[3] ?? null) : (parts[2] ?? null);
  return { slug, section, id };
}

function PageContent({ pathname }: { pathname: string }) {
  const { section, id } = parseRoute(pathname);

  if (section === "issues" && id) return <IssueDetail issueId={id} />;
  if (section === "issues") return <IssuesPage />;
  if (section === "agents") return <AgentsPage />;
  if (section === "autopilots" && id) return <AutopilotDetailPage autopilotId={id} />;
  if (section === "autopilots") return <AutopilotsPage />;
  if (section === "iframes" && id) return <IframeDetail />;
  if (section === "iframes") return <IframesPage />;
  if (section === "inbox") return <InboxPage />;
  if (section === "my-issues") return <MyIssuesPage />;
  if (section === "projects" && id) return <ProjectDetail projectId={id} />;
  if (section === "projects") return <ProjectsPage />;
  if (section === "runtimes") return <RuntimesPage />;
  if (section === "settings") return <SettingsPage />;
  if (section === "skills") return <SkillsPage />;

  // Default: issues page
  return <IssuesPage />;
}

export function WorkspaceRouter() {
  const [pathname, setPathname] = useState(() =>
    typeof window !== "undefined" ? window.location.pathname : "/"
  );

  useEffect(() => {
    const update = () => setPathname(window.location.pathname);

    window.addEventListener("popstate", update);

    const origPush = history.pushState.bind(history);
    const origReplace = history.replaceState.bind(history);

    history.pushState = (...args) => { origPush(...args); update(); };
    history.replaceState = (...args) => { origReplace(...args); update(); };

    return () => {
      window.removeEventListener("popstate", update);
      history.pushState = origPush;
      history.replaceState = origReplace;
    };
  }, []);

  return (
    <DashboardLayout>
      <PageContent pathname={pathname} />
    </DashboardLayout>
  );
}

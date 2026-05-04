import { CoreProvider } from "@multica/core/platform";
import { AstroNavigationProvider } from "@/platform/navigation";
import { WorkspaceSlugProvider } from "@multica/core/paths/hooks";
import { ThemeProvider } from "@multica/ui/components/common/theme-provider";
import { Toaster } from "@multica/ui/components/ui/sonner";
import { setCurrentWorkspace } from "@multica/core/platform/workspace-storage";
import { useEffect, useState } from "react";

function deriveWsUrl(): string {
  if (typeof window === "undefined") return "";
  
  // Connect to WebSocket through nginx proxy
  // nginx listens on 3000 and forwards to backend on 8080
  const proto = window.location.protocol === "https:" ? "wss:" : "ws:";
  return `${proto}//${window.location.host}/ws`;
}

/** Read workspace slug exclusively from the last_workspace_slug cookie.
 * This is the only source of truth — URL slug is not used.
 */
function getSlugFromCookie(): string | null {
  if (typeof document === "undefined") return null;
  const match = document.cookie.match(/(?:^|;\s*)last_workspace_slug=([^;]+)/);
  return match ? decodeURIComponent(match[1]) : null;
}

interface AstroProvidersProps {
  children: React.ReactNode;
  workspaceSlug?: string;
}

export function AstroProviders({ children, workspaceSlug: initialSlug }: AstroProvidersProps) {
  const [slug, setSlug] = useState<string | null>(() => {
    if (typeof window === "undefined") return initialSlug ?? null;
    return getSlugFromCookie();
  });

  // Update slug when cookie changes (e.g. user navigates to different workspace)
  useEffect(() => {
    const updateSlug = () => {
      const newSlug = getSlugFromCookie();
      setSlug(newSlug);
      if (newSlug) setCurrentWorkspace(newSlug, null);
    };

    // Initial set
    updateSlug();

    window.addEventListener("popstate", updateSlug);
    const origPush = history.pushState.bind(history);
    const origReplace = history.replaceState.bind(history);

    history.pushState = (...args) => {
      origPush(...args);
      updateSlug();
    };
    history.replaceState = (...args) => {
      origReplace(...args);
      updateSlug();
    };

    return () => {
      window.removeEventListener("popstate", updateSlug);
      history.pushState = origPush;
      history.replaceState = origReplace;
    };
  }, []);

  return (
    <ThemeProvider forcedTheme="light" enableSystem={false}>
      <CoreProvider
        apiBaseUrl=""
        wsUrl={deriveWsUrl()}
        cookieAuth={true}
      >
        <WorkspaceSlugProvider slug={slug}>
          <AstroNavigationProvider>
            {children}
          </AstroNavigationProvider>
        </WorkspaceSlugProvider>
      </CoreProvider>
      <Toaster />
    </ThemeProvider>
  );
}


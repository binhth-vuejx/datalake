import { Suspense, useEffect, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { sanitizeNextUrl, useAuthStore } from "@multica/core/auth";
import { workspaceKeys } from "@multica/core/workspace/queries";
import {
  paths,
  resolvePostAuthDestination,
  useHasOnboarded,
} from "@multica/core/paths";
import { api } from "@multica/core/api";
import type { Workspace } from "@multica/core/types";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from "@multica/ui/components/ui/card";
import { Button } from "@multica/ui/components/ui/button";
import { Loader2 } from "lucide-react";
import { LoginPage, validateCliCallback } from "@multica/views/auth";
import { AstroProviders } from "@/components/astro-providers";
import { useNavigation } from "@multica/views/navigation";

const googleClientId = import.meta.env.PUBLIC_GOOGLE_CLIENT_ID;

function LoginPageContent() {
  const nav = useNavigation();
  const qc = useQueryClient();
  const user = useAuthStore((s) => s.user);
  const isLoading = useAuthStore((s) => s.isLoading);

  // Extract query params from URL
  const searchParams = new URLSearchParams(typeof window !== "undefined" ? window.location.search : "");
  const cliCallbackRaw = searchParams.get("cli_callback");
  const cliState = searchParams.get("cli_state") || "";
  const platform = searchParams.get("platform");
  const isDesktopHandoff = platform === "desktop" && !cliCallbackRaw;
  const nextUrl = sanitizeNextUrl(searchParams.get("next"));

  const [desktopToken, setDesktopToken] = useState<string | null>(null);
  const [desktopError, setDesktopError] = useState("");
  const hasOnboarded = useHasOnboarded();

  // Already authenticated — honor ?next= or fall back to first workspace
  useEffect(() => {
    if (isLoading || !user || cliCallbackRaw) return;
    if (isDesktopHandoff) {
      api
        .issueCliToken()
        .then(({ token }) => {
          setDesktopToken(token);
          window.location.href = `multica://auth/callback?token=${encodeURIComponent(token)}`;
        })
        .catch((err) => {
          setDesktopError(
            err instanceof Error ? err.message : "Failed to prepare Desktop sign-in",
          );
        });
      return;
    }
    if (!hasOnboarded) {
      nav.push(paths.onboarding());
      return;
    }
    if (nextUrl) {
      nav.push(nextUrl);
      return;
    }
    const list = qc.getQueryData<Workspace[]>(workspaceKeys.list()) ?? [];
    nav.push(resolvePostAuthDestination(list, hasOnboarded));
  }, [isLoading, user, nav, nextUrl, cliCallbackRaw, isDesktopHandoff, hasOnboarded, qc]);

  const handleSuccess = () => {
    const currentUser = useAuthStore.getState().user;
    const onboarded = currentUser?.onboarded_at != null;
    if (!onboarded) {
      nav.push(paths.onboarding());
      return;
    }
    if (nextUrl) {
      nav.push(nextUrl);
      return;
    }
    const list = qc.getQueryData<Workspace[]>(workspaceKeys.list()) ?? [];
    nav.push(resolvePostAuthDestination(list, onboarded));
  };

  const googleState = [
    platform === "desktop" ? "platform:desktop" : "",
    nextUrl ? `next:${nextUrl}` : "",
  ]
    .filter(Boolean)
    .join(",") || undefined;

  if (isDesktopHandoff && user) {
    if (desktopError) {
      return (
        <div className="flex min-h-screen items-center justify-center">
          <Card className="w-full max-w-sm">
            <CardHeader className="text-center">
              <CardTitle className="text-2xl">Sign-in Failed</CardTitle>
              <CardDescription>{desktopError}</CardDescription>
            </CardHeader>
          </Card>
        </div>
      );
    }
    return (
      <div className="flex min-h-screen items-center justify-center">
        <Card className="w-full max-w-sm">
          <CardHeader className="text-center">
            <CardTitle className="text-2xl">Opening Multica</CardTitle>
            <CardDescription>
              {desktopToken
                ? "You should see a prompt to open the Multica desktop app. If nothing happens, click the button below."
                : "Preparing Desktop sign-in..."}
            </CardDescription>
          </CardHeader>
          <CardContent className="flex justify-center">
            {desktopToken ? (
              <Button
                variant="outline"
                onClick={() => {
                  window.location.href = `multica://auth/callback?token=${encodeURIComponent(desktopToken)}`;
                }}
              >
                Open Multica Desktop
              </Button>
            ) : (
              <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
            )}
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <LoginPage
      onSuccess={handleSuccess}
      google={
        googleClientId
          ? {
              clientId: googleClientId,
              redirectUri: `${typeof window !== "undefined" ? window.location.origin : ""}/auth/callback`,
              state: googleState,
            }
          : undefined
      }
      cliCallback={
        cliCallbackRaw && validateCliCallback(cliCallbackRaw)
          ? { url: cliCallbackRaw, state: cliState }
          : undefined
      }
    />
  );
}

export default function LoginPageWrapper() {
  return (
    <AstroProviders>
      <Suspense fallback={null}>
        <LoginPageContent />
      </Suspense>
    </AstroProviders>
  );
}

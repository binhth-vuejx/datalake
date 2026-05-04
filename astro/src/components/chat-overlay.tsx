"use client";

import { Component, type ReactNode } from "react";
import { AstroProviders } from "@/components/astro-providers";
import { ChatFab } from "@multica/views/chat";
import { ChatWindow } from "@multica/views/chat";
import { TooltipProvider } from "@multica/ui/components/ui/tooltip";

interface ChatOverlayProps {
  workspaceSlug?: string;
}

/**
 * Catches any React error inside the chat overlay so a crash here
 * never propagates to the rest of the page.
 */
class ChatErrorBoundary extends Component<
  { children: ReactNode },
  { hasError: boolean }
> {
  state = { hasError: false };

  static getDerivedStateFromError() {
    return { hasError: true };
  }

  componentDidCatch(error: unknown, info: unknown) {
    console.error("[ChatOverlay] error boundary caught:", error, info);
  }

  render() {
    if (this.state.hasError) {
      // Minimal fallback FAB so the user can still reload/interact
      return (
        <button
          style={{
            position: "fixed",
            bottom: "8px",
            right: "8px",
            zIndex: 99999,
            width: "40px",
            height: "40px",
            borderRadius: "50%",
            background: "#f87171",
            border: "none",
            cursor: "pointer",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            boxShadow: "0 2px 8px rgba(0,0,0,0.25)",
            pointerEvents: "auto",
          }}
          title="Chat error — click to reload"
          onClick={() => window.location.reload()}
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="white" strokeWidth="2.5">
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
          </svg>
        </button>
      );
    }
    return this.props.children;
  }
}

/**
 * Standalone chat overlay — renders ChatFab + ChatWindow with full providers.
 * Injected into every page via base.astro as an isolated Astro island
 * (client:only="react") so it survives crashes in other islands.
 *
 * Wrapped in ChatErrorBoundary so even an internal crash shows a fallback
 * instead of a blank corner.
 */
export default function ChatOverlay({ workspaceSlug }: ChatOverlayProps) {
  return (
    <ChatErrorBoundary>
      <AstroProviders workspaceSlug={workspaceSlug}>
        <TooltipProvider>
          {/*
            Full-screen fixed container acts as the positioning parent for
            ChatWindow (absolute) and ChatFab (absolute bottom-2 right-2).
            pointer-events: none on the container so it doesn't block page clicks,
            but ChatWindow/ChatFab re-enable pointer-events on themselves.
          */}
          <div
            className="fixed inset-0 pointer-events-none"
            style={{ zIndex: 99990 }}
          >
            <ChatWindow />
            <ChatFab />
          </div>
        </TooltipProvider>
      </AstroProviders>
    </ChatErrorBoundary>
  );
}

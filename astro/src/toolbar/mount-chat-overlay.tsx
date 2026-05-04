/**
 * mount-chat-overlay.tsx
 *
 * Mount ChatOverlay (React) vào một DOM node độc lập ngoài Astro island tree.
 * Được inject vào mọi trang qua nginx sub_filter.
 * Đảm bảo widget luôn tồn tại dù trang có lỗi frontend ở bất kỳ đâu.
 */

import { Component, type ReactNode } from "react";
import { createRoot } from "react-dom/client";
import ChatOverlay from "@/components/chat-overlay";
// Import global CSS so Tailwind classes work even on Astro error pages
import "@/styles/globals.css";

/** Bắt lỗi React bên trong widget — không để crash lan ra ngoài */
class ChatErrorBoundary extends Component<
  { children: ReactNode },
  { hasError: boolean }
> {
  state = { hasError: false };

  static getDerivedStateFromError() {
    return { hasError: true };
  }

  componentDidCatch(error: unknown) {
    console.error("[ChatOverlay] React error caught by boundary:", error);
  }

  render() {
    if (this.state.hasError) {
      return (
        <button
          style={{
            position: "fixed",
            bottom: "8px",
            right: "8px",
            zIndex: 49,
            width: "40px",
            height: "40px",
            borderRadius: "50%",
            background: "#f87171",
            border: "none",
            cursor: "pointer",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            boxShadow: "0 2px 8px rgba(0,0,0,0.2)",
            pointerEvents: "auto",
          }}
          title="Chat widget error — click to reload"
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

function mountChatOverlay() {
  // Tránh mount 2 lần
  if (document.getElementById("__multica-chat-overlay-root")) return;

  // Container: position fixed, full screen, pointer-events none
  // Dùng inline style để không phụ thuộc vào Tailwind hay bất kỳ CSS nào của page
  const container = document.createElement("div");
  container.id = "__multica-chat-overlay-root";
  container.style.cssText = [
    "position: fixed",
    "top: 0",
    "left: 0",
    "right: 0",
    "bottom: 0",
    "pointer-events: none",
    "z-index: 99999",
    "overflow: visible",
    "isolation: isolate",
  ].join("; ");
  document.body.appendChild(container);

  const slugMatch = window.location.pathname.match(/^\/([^/]+)/);
  const workspaceSlug = slugMatch?.[1] ?? undefined;

  const root = createRoot(container);
  root.render(
    <ChatErrorBoundary>
      <ChatOverlay workspaceSlug={workspaceSlug} />
    </ChatErrorBoundary>
  );
}

if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", mountChatOverlay);
} else {
  mountChatOverlay();
}

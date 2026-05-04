/**
 * Content Transition Wrapper
 * 
 * Wraps the main content area with a view-transition-name
 * so only the content animates, not the sidebar.
 */

import { ReactNode } from "react";

interface ContentTransitionWrapperProps {
  children: ReactNode;
}

export function ContentTransitionWrapper({
  children,
}: ContentTransitionWrapperProps) {
  return (
    <div
      style={{
        viewTransitionName: "main-content",
      } as React.CSSProperties}
      className="flex-1 overflow-auto"
    >
      {children}
    </div>
  );
}

"use client";

import { forwardRef } from "react";
import { useNavigation } from "./context";

interface AppLinkProps extends React.AnchorHTMLAttributes<HTMLAnchorElement> {
  href: string;
}

export const AppLink = forwardRef<HTMLAnchorElement, AppLinkProps>(
  function AppLink({ href, children, onClick, ...props }, ref) {
    const { openInNewTab } = useNavigation();

    const handleClick = (e: React.MouseEvent<HTMLAnchorElement>) => {
      // Handle modifier keys (Cmd/Ctrl/Shift) for opening in new tab
      if (e.metaKey || e.ctrlKey || e.shiftKey) {
        if (openInNewTab) {
          e.preventDefault();
          openInNewTab(href);
        }
        return;
      }

      // Call custom onClick handler if provided
      onClick?.(e);

      // Let the browser handle the click naturally
      // The view-transitions.ts script will intercept it and apply transitions
      // Don't prevent default - let it propagate
    };

    return (
      <a ref={ref} href={href} onClick={handleClick} {...props}>
        {children}
      </a>
    );
  },
);

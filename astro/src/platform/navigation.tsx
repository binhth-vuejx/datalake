import { useEffect, useState } from "react";
import {
  NavigationProvider,
  type NavigationAdapter,
} from "@multica/views/navigation";

function getLocation() {
  return {
    pathname: window.location.pathname,
    searchParams: new URLSearchParams(window.location.search),
  };
}

/**
 * Astro View Transitions Navigation Provider
 * 
 * Works with Astro's View Transitions by allowing natural <a> tag navigation.
 * The AppLink component lets the browser handle clicks naturally so Astro
 * can intercept them and apply View Transitions.
 */
export function AstroNavigationProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const [location, setLocation] = useState(getLocation);

  useEffect(() => {
    // Update location on popstate (back/forward buttons)
    const onPopState = () => setLocation(getLocation());
    window.addEventListener("popstate", onPopState);

    // Update location after Astro transitions
    const handleAfterSwap = () => {
      setLocation(getLocation());
    };
    document.addEventListener("astro:after-swap", handleAfterSwap);

    return () => {
      window.removeEventListener("popstate", onPopState);
      document.removeEventListener("astro:after-swap", handleAfterSwap);
    };
  }, []);

  const adapter: NavigationAdapter = {
    // For Astro, we rely on natural <a> tag navigation
    // The AppLink component allows clicks to propagate naturally
    push: (path: string) => {
      window.location.href = path;
    },
    replace: (path: string) => {
      window.location.replace(path);
    },
    back: () => history.back(),
    pathname: location.pathname,
    searchParams: location.searchParams,
  };

  return <NavigationProvider value={adapter}>{children}</NavigationProvider>;
}

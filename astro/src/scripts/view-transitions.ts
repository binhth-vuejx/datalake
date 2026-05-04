/**
 * View Transitions Handler
 * 
 * Manages Astro View Transitions lifecycle and intercepts navigation
 * to apply smooth transitions between pages - only for content, not sidebar.
 */

let isTransitioning = false;

/**
 * Perform a navigation with View Transitions
 */
async function navigateWithTransition(url: string, replace: boolean = false) {
  if (isTransitioning) return;
  isTransitioning = true;

  try {
    // Check if View Transitions API is available
    if (!("startViewTransition" in document)) {
      // Fallback: full page reload
      if (replace) {
        window.location.replace(url);
      } else {
        window.location.href = url;
      }
      return;
    }

    // Fetch the new page
    const response = await fetch(url);
    if (!response.ok) throw new Error(`Failed to fetch ${url}`);

    const html = await response.text();

    // Parse the new document
    const parser = new DOMParser();
    const newDoc = parser.parseFromString(html, "text/html");

    // Start view transition
    (document as any).startViewTransition(async () => {
      // Update history
      if (replace) {
        history.replaceState(null, "", url);
      } else {
        history.pushState(null, "", url);
      }

      // Find the main content area in both old and new documents
      const oldContent = document.querySelector('[data-astro-content]');
      const newContent = newDoc.querySelector('[data-astro-content]');

      if (oldContent && newContent) {
        // Only replace the content area, keep sidebar intact
        oldContent.replaceWith(newContent);
      } else {
        // Fallback: replace entire body if content wrapper not found
        document.head.replaceChildren(...newDoc.head.childNodes);
        document.body.replaceChildren(...newDoc.body.childNodes);
      }

      // Dispatch custom event for Astro
      document.dispatchEvent(new Event("astro:after-swap"));
    });
  } catch (error) {
    console.error("Navigation error:", error);
    // Fallback to full page reload on error
    if (replace) {
      window.location.replace(url);
    } else {
      window.location.href = url;
    }
  } finally {
    isTransitioning = false;
  }
}

/**
 * Initialize View Transitions event listeners
 */
export function initViewTransitions() {
  // Intercept all link clicks
  document.addEventListener("click", (e: Event) => {
    const target = (e.target as HTMLElement).closest("a");
    if (!target) return;

    const href = target.getAttribute("href");
    if (!href || href.startsWith("#") || href.startsWith("http")) return;

    // Check for modifier keys
    const event = e as MouseEvent;
    if (event.metaKey || event.ctrlKey || event.shiftKey) return;

    // Prevent default and use our transition handler
    e.preventDefault();
    navigateWithTransition(href, false);
  });

  // Listen for Astro's navigation events
  document.addEventListener("astro:before-preparation", () => {
    isTransitioning = true;
  });

  document.addEventListener("astro:after-swap", () => {
    isTransitioning = false;
  });

  document.addEventListener("astro:page-load", () => {
    isTransitioning = false;
  });
}

/**
 * Check if a transition is currently in progress
 */
export function isViewTransitioning(): boolean {
  return isTransitioning;
}

// Initialize on script load
if (typeof document !== "undefined") {
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", initViewTransitions);
  } else {
    initViewTransitions();
  }
}


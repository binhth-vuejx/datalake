import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";
import { render, act } from "@testing-library/react";
import React from "react";
import * as fc from "fast-check";
import { AstroNavigationProvider } from "./navigation";

// Create a context to capture the adapter value passed to NavigationProvider
const AdapterContext = React.createContext<any>(null);

// Mock the NavigationProvider from @multica/views/navigation
vi.mock("@multica/views/navigation", () => ({
  NavigationProvider: ({ children, value }: any) => {
    return (
      <AdapterContext.Provider value={value}>
        {children}
      </AdapterContext.Provider>
    );
  },
}));

describe("AstroNavigationProvider - Property-Based Tests", () => {
  beforeEach(() => {
    // Reset window.location to a known state
    delete (window as any).location;
    window.location = new URL("http://localhost/") as any;
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  /**
   * Property 2: AstroNavigationProvider push/replace correctness
   * **Validates: Requirements 4.2, 4.4, 4.5**
   *
   * For any valid path string (starts with `/`), `push(path)` SHALL result in
   * `adapter.pathname === path` — minimum 100 iterations with fast-check
   *
   * Tag: Feature: web-react-astro-integration, Property 2: AstroNavigationProvider push/replace correctness
   */
  it(
    "Property 2: For any valid path string, push(path) SHALL result in adapter.pathname === path",
    () => {
      fc.assert(
        fc.property(
          // Generate valid path strings: start with "/" followed by any alphanumeric, hyphens, underscores, slashes
          fc.stringMatching(/^\/[a-zA-Z0-9\-_/]*$/),
          (path: string) => {
            // Reset window.location for each iteration
            delete (window as any).location;
            window.location = new URL("http://localhost/") as any;

            // Mock history.pushState to update window.location.pathname
            const originalPushState = window.history.pushState;
            vi.spyOn(window.history, "pushState").mockImplementation(
              (state: any, title: string, url: string) => {
                // Update window.location.pathname to simulate browser behavior
                delete (window as any).location;
                window.location = new URL(`http://localhost${url}`) as any;
              }
            );

            let capturedAdapter: any;

            const TestComponent = () => {
              const adapter = React.useContext(AdapterContext);
              React.useEffect(() => {
                if (adapter) {
                  capturedAdapter = adapter;
                }
              }, [adapter]);
              return null;
            };

            render(
              <AstroNavigationProvider>
                <TestComponent />
              </AstroNavigationProvider>
            );

            // Ensure adapter is captured
            expect(capturedAdapter).toBeDefined();

            // Call push with the generated path
            act(() => {
              capturedAdapter.push(path);
            });

            // Verify that the adapter.pathname was updated
            // Since we mocked pushState to update window.location, the adapter should reflect this
            expect(capturedAdapter.pathname).toBe(path);

            // Restore original pushState
            window.history.pushState = originalPushState;
          }
        ),
        { numRuns: 100 }
      );
    },
    { timeout: 30000 }
  );

  /**
   * Property 2 (replace variant): AstroNavigationProvider replace correctness
   * **Validates: Requirements 4.2, 4.5**
   *
   * For any valid path string, `replace(path)` SHALL result in
   * `adapter.pathname === path` — minimum 100 iterations with fast-check
   *
   * Tag: Feature: web-react-astro-integration, Property 2: AstroNavigationProvider push/replace correctness
   */
  it(
    "Property 2 (replace variant): For any valid path string, replace(path) SHALL result in adapter.pathname === path",
    () => {
      fc.assert(
        fc.property(
          // Generate valid path strings: start with "/" followed by any alphanumeric, hyphens, underscores, slashes
          fc.stringMatching(/^\/[a-zA-Z0-9\-_/]*$/),
          (path: string) => {
            // Reset window.location for each iteration
            delete (window as any).location;
            window.location = new URL("http://localhost/") as any;

            // Mock history.replaceState to update window.location.pathname
            const originalReplaceState = window.history.replaceState;
            vi.spyOn(window.history, "replaceState").mockImplementation(
              (state: any, title: string, url: string) => {
                // Update window.location.pathname to simulate browser behavior
                delete (window as any).location;
                window.location = new URL(`http://localhost${url}`) as any;
              }
            );

            let capturedAdapter: any;

            const TestComponent = () => {
              const adapter = React.useContext(AdapterContext);
              React.useEffect(() => {
                if (adapter) {
                  capturedAdapter = adapter;
                }
              }, [adapter]);
              return null;
            };

            render(
              <AstroNavigationProvider>
                <TestComponent />
              </AstroNavigationProvider>
            );

            // Ensure adapter is captured
            expect(capturedAdapter).toBeDefined();

            // Call replace with the generated path
            act(() => {
              capturedAdapter.replace(path);
            });

            // Verify that the adapter.pathname was updated
            // Since we mocked replaceState to update window.location, the adapter should reflect this
            expect(capturedAdapter.pathname).toBe(path);

            // Restore original replaceState
            window.history.replaceState = originalReplaceState;
          }
        ),
        { numRuns: 100 }
      );
    },
    { timeout: 30000 }
  );
});

import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";
import { render, waitFor, act } from "@testing-library/react";
import React from "react";
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

describe("AstroNavigationProvider", () => {
  beforeEach(() => {
    // Reset window.location to a known state
    delete (window as any).location;
    window.location = new URL("http://localhost/initial-path") as any;
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  it("should mount component and verify initial pathname equals window.location.pathname", () => {
    window.location = new URL("http://localhost/test-initial") as any;

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

    expect(capturedAdapter).toBeDefined();
    expect(capturedAdapter.pathname).toBe("/test-initial");
  });

  it("should call push with a path and verify history.pushState is called", async () => {
    const pushStateSpy = vi.spyOn(window.history, "pushState");
    window.location = new URL("http://localhost/start") as any;

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

    // Call push
    await act(async () => {
      capturedAdapter.push("/test-path");
    });

    // Verify pushState was called with correct arguments
    expect(pushStateSpy).toHaveBeenCalledWith(null, "", "/test-path");

    pushStateSpy.mockRestore();
  });

  it("should call replace with a path and verify history.replaceState is called", async () => {
    const replaceStateSpy = vi.spyOn(window.history, "replaceState");
    window.location = new URL("http://localhost/start") as any;

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

    // Call replace
    await act(async () => {
      capturedAdapter.replace("/replaced-path");
    });

    // Verify replaceState was called with correct arguments
    expect(replaceStateSpy).toHaveBeenCalledWith(null, "", "/replaced-path");

    replaceStateSpy.mockRestore();
  });

  it("should handle popstate event and re-render with updated pathname", async () => {
    window.location = new URL("http://localhost/page1") as any;

    let renderCount = 0;
    const TestComponent = () => {
      const adapter = React.useContext(AdapterContext);
      renderCount++;
      return <div data-testid="pathname-display">{adapter?.pathname}</div>;
    };

    const { getByTestId } = render(
      <AstroNavigationProvider>
        <TestComponent />
      </AstroNavigationProvider>
    );

    // Initial pathname
    expect(getByTestId("pathname-display").textContent).toBe("/page1");
    const initialRenderCount = renderCount;

    // Simulate popstate event (back button)
    // In jsdom, we need to manually update window.location before dispatching popstate
    await act(async () => {
      window.location = new URL("http://localhost/page2") as any;
      window.dispatchEvent(new PopStateEvent("popstate"));
    });

    // Verify component re-rendered
    expect(renderCount).toBeGreaterThan(initialRenderCount);

    // Verify component re-renders with updated pathname
    await waitFor(() => {
      expect(getByTestId("pathname-display").textContent).toBe("/page2");
    });
  });

  it("should call history.back() when back() is called", () => {
    const backSpy = vi.spyOn(window.history, "back");

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

    // Call back
    act(() => {
      capturedAdapter.back();
    });

    // Verify history.back() was called
    expect(backSpy).toHaveBeenCalled();

    backSpy.mockRestore();
  });

  it("should provide searchParams from window.location.search", () => {
    window.location = new URL("http://localhost/page?foo=bar&baz=qux") as any;

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

    expect(capturedAdapter.searchParams).toBeInstanceOf(URLSearchParams);
    expect(capturedAdapter.searchParams.get("foo")).toBe("bar");
    expect(capturedAdapter.searchParams.get("baz")).toBe("qux");
  });
});

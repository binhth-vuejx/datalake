import { describe, it, expect } from "vitest";
import { paths, isGlobalPath } from "./paths";

describe("paths.workspace(slug)", () => {
  const ws = paths.workspace("acme");

  it("builds dashboard paths with /agents/ prefix", () => {
    expect(ws.issues()).toBe("/agents/acme/issues");
    expect(ws.issueDetail("abc-123")).toBe("/agents/acme/issues/abc-123");
    expect(ws.projects()).toBe("/agents/acme/projects");
    expect(ws.projectDetail("p1")).toBe("/agents/acme/projects/p1");
    expect(ws.autopilots()).toBe("/agents/acme/autopilots");
    expect(ws.autopilotDetail("a1")).toBe("/agents/acme/autopilots/a1");
    expect(ws.agents()).toBe("/agents/acme/agents");
    expect(ws.inbox()).toBe("/agents/acme/inbox");
    expect(ws.myIssues()).toBe("/agents/acme/my-issues");
    expect(ws.runtimes()).toBe("/agents/acme/runtimes");
    expect(ws.skills()).toBe("/agents/acme/skills");
    expect(ws.settings()).toBe("/agents/acme/settings");
  });

  it("URL-encodes special characters in ids", () => {
    expect(ws.issueDetail("id with space")).toBe("/agents/acme/issues/id%20with%20space");
  });
});

describe("paths (global)", () => {
  it("builds global paths without slug", () => {
    expect(paths.login()).toBe("/login");
    expect(paths.newWorkspace()).toBe("/workspaces/new");
    expect(paths.invite("inv-1")).toBe("/invite/inv-1");
    expect(paths.authCallback()).toBe("/auth/callback");
  });
});

describe("isGlobalPath", () => {
  it("returns true for pre-workspace routes", () => {
    expect(isGlobalPath("/login")).toBe(true);
    expect(isGlobalPath("/workspaces/new")).toBe(true);
    expect(isGlobalPath("/invite/abc")).toBe(true);
    expect(isGlobalPath("/auth/callback")).toBe(true);
  });

  it("returns false for workspace-scoped paths", () => {
    expect(isGlobalPath("/acme/issues")).toBe(false);
    expect(isGlobalPath("/")).toBe(false);
  });
});

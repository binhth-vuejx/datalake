<script lang="ts">
  let { workspaceSlug } = $props<{ workspaceSlug: string }>();

  let stats = $state([
    { label: "Total Users", value: 128, icon: "👥" },
    { label: "Active Sessions", value: 42, icon: "🟢" },
    { label: "Issues Open", value: 17, icon: "🐛" },
    { label: "Deployments", value: 5, icon: "🚀" },
  ]);
</script>

<div class="admin-layout">
  <header>
    <h1>⚙️ Admin Dashboard</h1>
    <p>Workspace: <strong>{workspaceSlug}</strong></p>
    <nav>
      <a href="/admin/{workspaceSlug}">Dashboard</a>
      <a href="/admin/{workspaceSlug}/users">Users</a>
      <a href="/agents/{workspaceSlug}">→ Agent App</a>
    </nav>
  </header>

  <main>
    <div class="stats-grid">
      {#each stats as stat}
        <div class="stat-card">
          <span class="icon">{stat.icon}</span>
          <span class="value">{stat.value}</span>
          <span class="label">{stat.label}</span>
        </div>
      {/each}
    </div>

    <div class="section">
      <h2>Recent Activity</h2>
      <ul>
        <li>User john@example.com logged in — 2 min ago</li>
        <li>Issue #42 created — 5 min ago</li>
        <li>Runtime "prod-1" restarted — 12 min ago</li>
      </ul>
    </div>
  </main>
</div>

<style>
  .admin-layout { font-family: system-ui; padding: 2rem; background: #0f172a; color: #e2e8f0; min-height: 100vh; }
  header { border-bottom: 1px solid #334155; padding-bottom: 1rem; margin-bottom: 2rem; }
  h1 { color: #38bdf8; margin: 0 0 0.5rem; }
  nav { display: flex; gap: 1rem; margin-top: 1rem; }
  nav a { color: #38bdf8; text-decoration: none; padding: 0.25rem 0.75rem; border: 1px solid #38bdf8; border-radius: 4px; }
  nav a:hover { background: #38bdf8; color: #0f172a; }
  .stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(150px, 1fr)); gap: 1rem; margin-bottom: 2rem; }
  .stat-card { background: #1e293b; border-radius: 8px; padding: 1.5rem; text-align: center; }
  .icon { font-size: 2rem; display: block; }
  .value { font-size: 2rem; font-weight: bold; color: #38bdf8; display: block; }
  .label { color: #94a3b8; font-size: 0.875rem; }
  .section { background: #1e293b; border-radius: 8px; padding: 1.5rem; }
  h2 { color: #38bdf8; margin-top: 0; }
  ul { list-style: none; padding: 0; }
  li { padding: 0.5rem 0; border-bottom: 1px solid #334155; color: #94a3b8; }
</style>

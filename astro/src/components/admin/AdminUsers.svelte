<script lang="ts">
  let { workspaceSlug } = $props<{ workspaceSlug: string }>();

  let users = $state([
    { id: 1, name: "Alice Johnson", email: "alice@example.com", role: "Admin", status: "active" },
    { id: 2, name: "Bob Smith", email: "bob@example.com", role: "Member", status: "active" },
    { id: 3, name: "Carol White", email: "carol@example.com", role: "Member", status: "inactive" },
    { id: 4, name: "Dave Brown", email: "dave@example.com", role: "Viewer", status: "active" },
  ]);

  let search = $state("");
  let filtered = $derived(users.filter(u =>
    u.name.toLowerCase().includes(search.toLowerCase()) ||
    u.email.toLowerCase().includes(search.toLowerCase())
  ));
</script>

<div class="users-page">
  <header>
    <h1>👥 User Management</h1>
    <p>Workspace: <strong>{workspaceSlug}</strong></p>
    <nav>
      <a href="/admin/{workspaceSlug}">← Dashboard</a>
    </nav>
  </header>

  <main>
    <div class="toolbar">
      <input bind:value={search} placeholder="Search users..." />
      <span>{filtered.length} users</span>
    </div>

    <table>
      <thead>
        <tr><th>Name</th><th>Email</th><th>Role</th><th>Status</th></tr>
      </thead>
      <tbody>
        {#each filtered as user}
          <tr>
            <td>{user.name}</td>
            <td>{user.email}</td>
            <td><span class="badge role">{user.role}</span></td>
            <td><span class="badge {user.status}">{user.status}</span></td>
          </tr>
        {/each}
      </tbody>
    </table>
  </main>
</div>

<style>
  .users-page { font-family: system-ui; padding: 2rem; background: #0f172a; color: #e2e8f0; min-height: 100vh; }
  header { border-bottom: 1px solid #334155; padding-bottom: 1rem; margin-bottom: 2rem; }
  h1 { color: #38bdf8; margin: 0 0 0.5rem; }
  nav a { color: #38bdf8; text-decoration: none; }
  .toolbar { display: flex; gap: 1rem; align-items: center; margin-bottom: 1rem; }
  input { background: #1e293b; border: 1px solid #334155; color: #e2e8f0; padding: 0.5rem 1rem; border-radius: 4px; flex: 1; }
  table { width: 100%; border-collapse: collapse; background: #1e293b; border-radius: 8px; overflow: hidden; }
  th { background: #334155; padding: 0.75rem 1rem; text-align: left; color: #94a3b8; font-size: 0.875rem; }
  td { padding: 0.75rem 1rem; border-bottom: 1px solid #334155; }
  .badge { padding: 0.2rem 0.6rem; border-radius: 9999px; font-size: 0.75rem; font-weight: 600; }
  .badge.role { background: #1d4ed8; color: #bfdbfe; }
  .badge.active { background: #166534; color: #bbf7d0; }
  .badge.inactive { background: #7f1d1d; color: #fecaca; }
</style>

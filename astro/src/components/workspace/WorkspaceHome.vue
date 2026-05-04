<template>
  <div class="workspace-home">
    <header>
      <h1>🏠 {{ workspaceSlug }}</h1>
      <p>Welcome to your workspace</p>
      <nav>
        <a :href="`/${workspaceSlug}/dashboard`">Dashboard</a>
        <a :href="`/agents/${workspaceSlug}`">→ Agent App</a>
        <a :href="`/admin/${workspaceSlug}`">→ Admin</a>
      </nav>
    </header>

    <main>
      <div class="quick-links">
        <div v-for="link in quickLinks" :key="link.href" class="link-card">
          <span class="icon">{{ link.icon }}</span>
          <a :href="link.href">{{ link.label }}</a>
          <p>{{ link.desc }}</p>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
defineProps<{ workspaceSlug: string }>();

const quickLinks = [
  { href: "#", icon: "📋", label: "Issues", desc: "Track and manage issues" },
  { href: "#", icon: "🤖", label: "Agents", desc: "AI coding agents" },
  { href: "#", icon: "📦", label: "Projects", desc: "Organize your work" },
  { href: "#", icon: "📬", label: "Inbox", desc: "Notifications & updates" },
];
</script>

<style scoped>
.workspace-home { font-family: system-ui; padding: 2rem; background: #0f172a; color: #e2e8f0; min-height: 100vh; }
header { border-bottom: 1px solid #334155; padding-bottom: 1rem; margin-bottom: 2rem; }
h1 { color: #a78bfa; margin: 0 0 0.5rem; }
nav { display: flex; gap: 1rem; margin-top: 1rem; }
nav a { color: #a78bfa; text-decoration: none; padding: 0.25rem 0.75rem; border: 1px solid #a78bfa; border-radius: 4px; }
nav a:hover { background: #a78bfa; color: #0f172a; }
.quick-links { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 1rem; }
.link-card { background: #1e293b; border-radius: 8px; padding: 1.5rem; }
.icon { font-size: 2rem; display: block; margin-bottom: 0.5rem; }
a { color: #a78bfa; font-weight: 600; font-size: 1.1rem; }
p { color: #94a3b8; margin: 0.5rem 0 0; font-size: 0.875rem; }
</style>

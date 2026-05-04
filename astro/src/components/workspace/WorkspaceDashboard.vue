<template>
  <div class="dashboard">
    <header>
      <h1>📊 Dashboard</h1>
      <p>Workspace: <strong>{{ workspaceSlug }}</strong></p>
      <nav>
        <a :href="`/${workspaceSlug}`">← Home</a>
        <a :href="`/agents/${workspaceSlug}`">Agent App</a>
      </nav>
    </header>

    <main>
      <div class="metrics">
        <div v-for="m in metrics" :key="m.label" class="metric">
          <span class="value">{{ m.value }}</span>
          <span class="label">{{ m.label }}</span>
          <span :class="['trend', m.trend > 0 ? 'up' : m.trend < 0 ? 'down' : 'neutral']">
            {{ m.trend > 0 ? '↑' : m.trend < 0 ? '↓' : '—' }}
            {{ m.trend !== 0 ? Math.abs(m.trend) + '%' : '' }}
          </span>
        </div>
      </div>

      <div class="activity">
        <h2>Recent Activity</h2>
        <div v-for="item in activity" :key="item.id" class="activity-item">
          <span class="dot"></span>
          <div>
            <strong>{{ item.title }}</strong>
            <span class="time">{{ item.time }}</span>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
defineProps<{ workspaceSlug: string }>();

const metrics = ref([
  { label: "Open Issues", value: 24, trend: -5 },
  { label: "Completed", value: 87, trend: 12 },
  { label: "In Progress", value: 11, trend: 3 },
  { label: "Agents Active", value: 3, trend: 0 },
]);

const activity = ref([
  { id: 1, title: "Issue #42 resolved by Agent", time: "2 min ago" },
  { id: 2, title: "New runtime 'dev-local' connected", time: "15 min ago" },
  { id: 3, title: "Project 'Backend API' updated", time: "1 hour ago" },
]);
</script>

<style scoped>
.dashboard { font-family: system-ui; padding: 2rem; background: #0f172a; color: #e2e8f0; min-height: 100vh; }
header { border-bottom: 1px solid #334155; padding-bottom: 1rem; margin-bottom: 2rem; }
h1 { color: #a78bfa; margin: 0 0 0.5rem; }
nav { display: flex; gap: 1rem; margin-top: 1rem; }
nav a { color: #a78bfa; text-decoration: none; padding: 0.25rem 0.75rem; border: 1px solid #a78bfa; border-radius: 4px; }
.metrics { display: grid; grid-template-columns: repeat(auto-fit, minmax(150px, 1fr)); gap: 1rem; margin-bottom: 2rem; }
.metric { background: #1e293b; border-radius: 8px; padding: 1.5rem; text-align: center; }
.value { font-size: 2.5rem; font-weight: bold; color: #a78bfa; display: block; }
.label { color: #94a3b8; font-size: 0.875rem; display: block; }
.trend { font-size: 0.75rem; display: block; margin-top: 0.25rem; }
.trend.up { color: #4ade80; }
.trend.down { color: #f87171; }
.trend.neutral { color: #64748b; }
.activity { background: #1e293b; border-radius: 8px; padding: 1.5rem; }
h2 { color: #a78bfa; margin-top: 0; }
.activity-item { display: flex; gap: 1rem; align-items: flex-start; padding: 0.75rem 0; border-bottom: 1px solid #334155; }
.dot { width: 8px; height: 8px; background: #a78bfa; border-radius: 50%; margin-top: 6px; flex-shrink: 0; }
.time { color: #64748b; font-size: 0.75rem; margin-left: 0.5rem; }
</style>

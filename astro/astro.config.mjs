import { defineConfig } from "astro/config";
import vue from "@astrojs/vue";
import svelte from "@astrojs/svelte";
import react from "@astrojs/react";
import tailwindcss from "@tailwindcss/vite";
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

// ── Setup monorepo root path ────────────────────────────────────────────────
const __dirname = path.dirname(fileURLToPath(import.meta.url));
const monorepoRoot = path.resolve(__dirname, "..");

// ── Backend URL for proxy ───────────────────────────────────────────────────
const BACKEND_URL = process.env.BACKEND_URL ?? "http://localhost:3000";

// ── Vite plugin: inject persistent widgets (Svelte) vào /@vite/client ────────
// + inject chat overlay FAB vào mọi HTML response (kể cả Astro error pages)
function persistentWidgetPlugin() {
  const widgetPath = path.resolve("src/toolbar/persistent-widget.js");
  const chatWidgetPath = path.resolve("src/toolbar/chat-widget.js");

  return {
    name: "persistent-widget",
    apply: "serve",
    enforce: "pre",

    configureServer(server) {
      let cachedClientCode = null;

      async function getProcessedClient() {
        if (cachedClientCode) return cachedClientCode;
        const result = await server.transformRequest("/@vite/client");
        if (result?.code) {
          cachedClientCode = result.code;
          return cachedClientCode;
        }
        return null;
      }

      // Intercept /@vite/client để inject persistent-widget + chat-widget (Svelte)
      server.middlewares.use("/@vite/client", async (_req, res, next) => {
        try {
          const clientCode = await getProcessedClient();
          if (!clientCode) return next();

          const widgetCode = fs.readFileSync(widgetPath, "utf-8");
          const chatWidgetCode = fs.existsSync(chatWidgetPath)
            ? fs.readFileSync(chatWidgetPath, "utf-8")
            : "";

          const combined = clientCode
            + "\n\n/* ── persistent-widget ── */\n" + widgetCode
            + (chatWidgetCode ? "\n\n/* ── chat-widget ── */\n" + chatWidgetCode : "");

          res.setHeader("Content-Type", "application/javascript");
          res.setHeader("Cache-Control", "no-cache");
          res.end(combined);
        } catch (e) {
          next(e);
        }
      });

      server.watcher.on("change", (file) => {
        if (file.includes("persistent-widget") || file.includes("chat-widget")) {
          cachedClientCode = null;
        }
      });
    },
  };
}

export default defineConfig({
  integrations: [
    vue(),
    svelte({
      compilerOptions: {
        hydratable: true,
      },
    }),
    react(),
  ],
  vite: {
    plugins: [
      tailwindcss(),
      persistentWidgetPlugin(),
    ],
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "src"),
        "@multica/core": path.resolve(monorepoRoot, "packages/core"),
        "@multica/ui": path.resolve(monorepoRoot, "packages/ui"),
        "@multica/views": path.resolve(monorepoRoot, "packages/views"),
      },
      dedupe: ["react", "react-dom", "@tanstack/react-query", "zustand"],
    },
    optimizeDeps: {
      include: ["react", "react-dom", "@tanstack/react-query", "zustand"],
    },
    ssr: {
      noExternal: ["@multica/core", "@multica/ui", "@multica/views"],
    },
  },
  devToolbar: { enabled: false },
  output: "server",
});

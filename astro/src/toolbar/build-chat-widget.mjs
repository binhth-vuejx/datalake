/**
 * build-chat-widget.mjs
 *
 * Bundle ChatWidget.svelte thành một file JS self-contained
 * để nhúng vào chat-widget.js (inject qua /@vite/client).
 *
 * Chạy: node src/toolbar/build-chat-widget.mjs
 */

import { build } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import path from "node:path";
import fs from "node:fs";
import { fileURLToPath } from "node:url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const outDir = path.resolve(__dirname, "../../.chat-widget-build");

await build({
  configFile: false,
  plugins: [
    svelte({
      compilerOptions: {
        customElement: true,
      },
    }),
  ],
  build: {
    outDir,
    emptyOutDir: true,
    lib: {
      entry: path.resolve(__dirname, "ChatWidget.svelte"),
      name: "ChatWidget",
      fileName: "chat-widget",
      formats: ["iife"],
    },
    rollupOptions: {
      output: {
        inlineDynamicImports: true,
      },
    },
    minify: false,
    sourcemap: false,
  },
  logLevel: "warn",
});

// Đọc output và ghi vào chat-widget.js
const builtFile = path.resolve(outDir, "chat-widget.iife.js");
const builtCode = fs.readFileSync(builtFile, "utf-8");

const output = `/**
 * chat-widget.js — AUTO GENERATED, do not edit manually.
 * Source: src/toolbar/ChatWidget.svelte
 * Rebuild: node src/toolbar/build-chat-widget.mjs
 */

${builtCode}

// Mount vào body
(function mount() {
  function doMount() {
    if (!document.querySelector("chat-widget")) {
      document.body.appendChild(document.createElement("chat-widget"));
    }
  }
  if (document.body) {
    doMount();
  } else {
    document.addEventListener("DOMContentLoaded", doMount);
  }
})();
`;

const outFile = path.resolve(__dirname, "chat-widget.js");
fs.writeFileSync(outFile, output, "utf-8");

console.log("✅ chat-widget.js rebuilt from ChatWidget.svelte");
console.log(`   Output: ${outFile}`);

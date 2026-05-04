/**
 * build-widget.mjs
 *
 * Bundle WidgetContent.svelte thành một file JS self-contained
 * để nhúng vào persistent-widget.js (inject qua /@vite/client).
 *
 * Chạy: node src/toolbar/build-widget.mjs
 */

import { build } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import path from "node:path";
import fs from "node:fs";
import { fileURLToPath } from "node:url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const outDir = path.resolve(__dirname, "../../.widget-build");

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
      entry: path.resolve(__dirname, "WidgetContent.svelte"),
      name: "PersistentWidget",
      fileName: "widget",
      formats: ["iife"], // IIFE = self-contained, không cần import
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

// Đọc output và ghi vào persistent-widget.js
const builtFile = path.resolve(outDir, "widget.iife.js");
const builtCode = fs.readFileSync(builtFile, "utf-8");

// Wrap: chỉ giữ phần define custom element + mount logic
const output = `/**
 * persistent-widget.js — AUTO GENERATED, do not edit manually.
 * Source: src/toolbar/WidgetContent.svelte
 * Rebuild: node src/toolbar/build-widget.mjs
 */

${builtCode}

// Mount vào body
(function mount() {
  if (!document.querySelector("persistent-widget")) {
    document.body.appendChild(document.createElement("persistent-widget"));
  }
  if (!document.body) {
    document.addEventListener("DOMContentLoaded", mount);
  }
})();
`;

const outFile = path.resolve(__dirname, "persistent-widget.js");
fs.writeFileSync(outFile, output, "utf-8");

console.log("✅ persistent-widget.js rebuilt from WidgetContent.svelte");
console.log(`   Output: ${outFile}`);

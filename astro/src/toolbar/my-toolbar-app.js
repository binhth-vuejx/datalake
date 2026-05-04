/**
 * Astro Dev Toolbar — custom inspect app.
 * Hover over any element to highlight it and see DOM info.
 */
export default {
  id: "my-app",
  name: "My App",
  icon: "🛠️",

  init(canvas, eventTarget) {
    // ── State ──────────────────────────────────────────────────────────────
    let inspectActive = false;
    let pinned = false;  // true = đang pin 1 element, hover tạm dừng

    // ── Main window ────────────────────────────────────────────────────────
    const win = document.createElement("astro-dev-toolbar-window");

    const title = document.createElement("h2");
    title.textContent = "My Custom Toolbar App";
    title.style.cssText = "margin:0 0 12px; font-size:16px; font-weight:600;";
    win.appendChild(title);

    // ── Inspect toggle button ──────────────────────────────────────────────
    const inspectBtn = document.createElement("astro-dev-toolbar-button");
    inspectBtn.textContent = "🔍 Start Inspect";
    inspectBtn.setAttribute("button-style", "purple");
    inspectBtn.style.display = "block";
    inspectBtn.style.marginBottom = "8px";

    inspectBtn.addEventListener("click", () => {
      inspectActive = !inspectActive;
      inspectBtn.textContent = inspectActive ? "✋ Stop Inspect" : "🔍 Start Inspect";
      inspectBtn.setAttribute("button-style", inspectActive ? "red" : "purple");

      if (!inspectActive) {
        pinned = false;
        removeHighlight();
      }

      eventTarget.dispatchEvent(
        new CustomEvent("app-notification", {
          detail: {
            state: inspectActive ? "info" : "warning",
            title: inspectActive ? "Inspect ON" : "Inspect OFF",
            description: inspectActive
              ? "Hover over elements to inspect them."
              : "Inspect mode disabled.",
          },
        }),
      );
    });
    win.appendChild(inspectBtn);

    // ── Divider ────────────────────────────────────────────────────────────
    const hr = document.createElement("hr");
    hr.style.cssText = "margin:12px 0; border-color:rgba(255,255,255,0.1);";
    win.appendChild(hr);

    // ── Info panel (shows hovered element info) ────────────────────────────
    const infoPanel = document.createElement("div");
    infoPanel.style.cssText = `
      font-family: monospace;
      font-size: 12px;
      color: #ccc;
      min-height: 60px;
      padding: 8px;
      background: rgba(255,255,255,0.05);
      border-radius: 6px;
      white-space: pre-wrap;
      word-break: break-all;
    `;
    infoPanel.textContent = "Hover over an element to inspect it.";
    win.appendChild(infoPanel);

    // ── Other buttons ──────────────────────────────────────────────────────
    const hr2 = document.createElement("hr");
    hr2.style.cssText = "margin:12px 0; border-color:rgba(255,255,255,0.1);";
    win.appendChild(hr2);

    const clearBtn = document.createElement("astro-dev-toolbar-button");
    clearBtn.textContent = "Clear localStorage";
    clearBtn.setAttribute("button-style", "outline");
    clearBtn.style.marginBottom = "8px";
    clearBtn.addEventListener("click", () => {
      localStorage.clear();
      eventTarget.dispatchEvent(
        new CustomEvent("app-notification", {
          detail: { state: "success", title: "Done", description: "localStorage cleared!" },
        }),
      );
    });
    win.appendChild(clearBtn);

    const logBtn = document.createElement("astro-dev-toolbar-button");
    logBtn.textContent = "Log page info";
    logBtn.addEventListener("click", () => {
      console.info("[Toolbar] URL:", window.location.href);
      console.info("[Toolbar] Islands:", document.querySelectorAll("[data-astro-source-file]").length);
      eventTarget.dispatchEvent(
        new CustomEvent("app-notification", {
          detail: { state: "info", title: "Logged", description: "Check the browser console." },
        }),
      );
    });
    win.appendChild(logBtn);

    canvas.appendChild(win);

    // ── Highlight helpers ──────────────────────────────────────────────────
    let highlight = null;

    function createHighlight() {
      if (highlight) return;
      highlight = document.createElement("div");
      highlight.style.cssText = `
        position: fixed;
        pointer-events: none;
        z-index: 49;
        outline: 2px solid #7c3aed;
        background: rgba(124, 58, 237, 0.1);
        border-radius: 2px;
        transition: all 0.05s;
      `;
      document.body.appendChild(highlight);
    }

    function moveHighlight(el) {
      if (!highlight) return;
      const rect = el.getBoundingClientRect();
      highlight.style.top    = rect.top + "px";
      highlight.style.left   = rect.left + "px";
      highlight.style.width  = rect.width + "px";
      highlight.style.height = rect.height + "px";
      // Pinned = solid border, hover = dashed
      highlight.style.outlineStyle = pinned ? "solid" : "dashed";
      highlight.style.outlineColor = pinned ? "#f59e0b" : "#7c3aed";
      highlight.style.background   = pinned
        ? "rgba(245,158,11,0.1)"
        : "rgba(124,58,237,0.1)";
    }

    function removeHighlight() {
      if (highlight) { highlight.remove(); highlight = null; }
    }

    function getVueComponentFile(el) {
      // Vue 3 stores component instance in various internal properties
      // We need to find the vnode first, then walk up the component tree

      function findVueInstance(node) {
        // Try multiple Vue 3 internal keys
        const vueKeys = [
          "__vueParentComponent",
          "__vnode",
          "_vnode",
          "__vue_app__",
        ];

        for (const key of vueKeys) {
          if (node[key]) {
            // If it's a vnode, get the component from it
            if (node[key].component) return node[key].component;
            // If it's already a component instance
            if (node[key].type) return node[key];
          }
        }

        // Check all keys that start with __v (Vue internals)
        for (const key of Object.keys(node)) {
          if (key.startsWith("__v") && node[key]?.component) {
            return node[key].component;
          }
          if (key.startsWith("__v") && node[key]?.type?.__file) {
            return node[key];
          }
        }

        return null;
      }

      // Walk up DOM tree, at each node try to find Vue instance
      let node = el;
      while (node && node !== document.body) {
        let vInst = findVueInstance(node);

        // Walk up Vue component tree from this instance
        let depth = 0;
        while (vInst && depth < 20) {
          const file = vInst.type?.__file;
          if (file) {
            // Found a component with __file — return it
            return file.replace(/.*\/src\//, "src/");
          }
          // Go to parent component
          vInst = vInst.parent;
          depth++;
        }

        node = node.parentElement;
      }
      return null;
    }

    function getSvelteComponentFile(el) {
      // Svelte stores component info in $$, $$props, etc.
      let node = el;
      while (node && node !== document.body) {
        // Svelte 3/4 internals
        if (node.$$ && node.$$.ctx) {
          // Try to extract file from component definition
          const comp = node.$$;
          // Svelte doesn't expose __file like Vue — use component-url from astro-island
          return null; // fallback to island detection
        }
        node = node.parentElement;
      }
      return null;
    }

    // ── Pre-scan: build a map of astro-island → component-url ─────────────
    // Astro wraps Vue/React islands in <astro-island component-url="...">
    // For static HTML (h1, p, etc.) we fall back to the current page URL
    function buildIslandMap() {
      const map = new Map(); // element → component file
      document.querySelectorAll("astro-island").forEach((island) => {
        const url = island.getAttribute("component-url") || "";
        const file = url.replace(/.*\/src\//, "src/").replace(/\?.*$/, "");
        // Map the island itself and all its descendants
        island.querySelectorAll("*").forEach((child) => {
          map.set(child, file);
        });
        map.set(island, file);
      });
      return map;
    }

    let islandMap = buildIslandMap();

    // Rebuild map when DOM changes (HMR)
    new MutationObserver(() => {
      islandMap = buildIslandMap();
    }).observe(document.body, { childList: true, subtree: true });

    function getAstroSourceFile(el) {
      // 1. Check if element is inside an astro-island (Vue/React component)
      const islandFile = islandMap.get(el);
      if (islandFile) {
        return { file: islandFile, loc: "", type: "island" };
      }

      // 2. Walk up to find astro-island ancestor
      let node = el.parentElement;
      while (node && node !== document.documentElement) {
        if (node.tagName?.toLowerCase() === "astro-island") {
          const url = node.getAttribute("component-url") || "";
          const file = url.replace(/.*\/src\//, "src/").replace(/\?.*$/, "");
          if (file) return { file, loc: "", type: "island" };
        }
        // Check data-astro-source-file (injected by Astro dev toolbar)
        const srcFile = node.getAttribute?.("data-astro-source-file");
        if (srcFile) {
          const loc = node.getAttribute?.("data-astro-source-loc") ?? "";
          return {
            file: srcFile.replace(/.*\/src\//, "src/"),
            loc,
            type: "astro",
          };
        }
        node = node.parentElement;
      }

      // 3. Fallback: static HTML element → belongs to current page
      // Derive page file from URL path
      const path = window.location.pathname;
      const pageName = path === "/" ? "index" : path.replace(/^\//, "").replace(/\/$/, "");
      return {
        file: `src/pages/${pageName}.astro`,
        loc: "",
        type: "page",
      };
    }

    function getElementInfo(el) {
      const tag = el.tagName.toLowerCase();
      const id = el.id ? `#${el.id}` : "";
      const classes = el.classList.length
        ? "." + Array.from(el.classList).join(".")
        : "";
      const rect = el.getBoundingClientRect();

      const lines = [
        `Tag:   <${tag}${id}${classes}>`,
        `Size:  ${Math.round(rect.width)}×${Math.round(rect.height)}px`,
        `Pos:   top ${Math.round(rect.top)}  left ${Math.round(rect.left)}`,
      ];

      lines.push(`─────────────────────────`);

      // 1. Page gốc — luôn hiện, derive từ URL
      const path = window.location.pathname;
      const pageName = path === "/" ? "index" : path.replace(/^\//, "").replace(/\/$/, "");
      lines.push(`📄 Page: src/pages/${pageName}.astro`);

      // 2. Astro island (nếu element nằm trong island)
      const astro = getAstroSourceFile(el);
      if (astro && astro.type === "island") {
        lines.push(`🟣 Island: ${astro.file}`);
      }

      // 3. Vue component (nearest)
      const vueFile = getVueComponentFile(el);
      if (vueFile) {
        lines.push(`🟢 Vue: ${vueFile}`);
      }

      // 4. Svelte component (nearest)
      const svelteFile = getSvelteComponentFile(el);
      if (svelteFile) {
        lines.push(`🟠 Svelte: ${svelteFile}`);
      }

      return lines.join("\n");
    }

    // ── Build markdown report từ element ──────────────────────────────────
    function buildMarkdown(el) {
      const tag = el.tagName.toLowerCase();
      const id = el.id ? `#${el.id}` : "";
      const classes = el.classList.length
        ? Array.from(el.classList).join(", ")
        : "_none_";
      const rect = el.getBoundingClientRect();

      const path = window.location.pathname;
      const pageName = path === "/" ? "index" : path.replace(/^\//, "").replace(/\/$/, "");
      const pageFile = `src/pages/${pageName}.astro`;

      const astro = getAstroSourceFile(el);
      const vueFile = getVueComponentFile(el);

      const lines = [
        `## 🐛 Yêu cầu từ Inspect`,
        ``,
        `### Thông tin element`,
        `| Trường | Giá trị |`,
        `|--------|---------|`,
        `| Tag | \`<${tag}${id}>\` |`,
        `| Classes | \`${classes}\` |`,
        `| Size | ${Math.round(rect.width)}×${Math.round(rect.height)}px |`,
        `| Position | top: ${Math.round(rect.top)}, left: ${Math.round(rect.left)} |`,
        `| URL | ${window.location.href} |`,
        ``,
        `### Vị trí trong code`,
        `| Layer | File |`,
        `|-------|------|`,
        `| 📄 Page | \`${pageFile}\` |`,
      ];

      if (astro && astro.type === "island") {
        lines.push(`| 🟣 Island | \`${astro.file}\` |`);
      }
      if (vueFile) {
        lines.push(`| 🟢 Vue Component | \`${vueFile}\` |`);
      }

      lines.push(``);
      lines.push(`### Mô tả yêu cầu`);
      lines.push(`> _(Điền mô tả yêu cầu tại đây)_`);

      return lines.join("\n");
    }

    // ── Mouse event listeners on the main document ─────────────────────────
    document.addEventListener("mouseover", (e) => {
      if (!inspectActive || pinned) return;
      const el = e.target;
      if (!el || el === document.body || el === document.documentElement) return;

      createHighlight();
      moveHighlight(el);
      infoPanel.textContent = getElementInfo(el);
    });

    document.addEventListener("mouseout", () => {
      if (!inspectActive || pinned) return;
      if (highlight) {
        highlight.style.outline = "none";
        highlight.style.background = "transparent";
      }
    });

    // ── Pin/Unpin dùng pointerdown (không intercept click của toolbar buttons) ──
    document.addEventListener("pointerdown", (e) => {
      if (!inspectActive) return;

      const el = e.target;
      if (!el || el === document.body || el === document.documentElement) return;

      // Bỏ qua nếu click vào toolbar overlay (highlight div) hoặc astro toolbar
      if (el === highlight) return;
      const tagName = el.tagName?.toLowerCase() ?? "";
      if (tagName.startsWith("astro-")) return;  // skip astro-dev-toolbar, astro-island, etc.

      e.preventDefault();

      if (pinned) {
        // Unpin → resume hover
        pinned = false;
        infoPanel.textContent = "Hover over an element to inspect it.";
        const oldCreateBtn = win.querySelector("#create-issue-btn");
        if (oldCreateBtn) oldCreateBtn.remove();
        removeHighlight();
        eventTarget.dispatchEvent(
          new CustomEvent("app-notification", {
            detail: { state: "info", title: "Unpinned", description: "Hover to inspect again." },
          }),
        );
      } else {
        // Pin current element
        const pinnedEl = el;
        pinned = true;
        createHighlight();
        moveHighlight(pinnedEl);
        const info = getElementInfo(pinnedEl);
        infoPanel.textContent = "📌 PINNED\n" + info;
        console.info("[Inspect pinned]", pinnedEl.tagName);

        // Xóa nút cũ nếu có
        const oldCreateBtn = win.querySelector("#create-issue-btn");
        if (oldCreateBtn) oldCreateBtn.remove();

        // Tạo nút "Tạo yêu cầu"
        const createBtn = document.createElement("astro-dev-toolbar-button");
        createBtn.id = "create-issue-btn";
        createBtn.textContent = "📝 Tạo yêu cầu";
        createBtn.setAttribute("button-style", "purple");
        createBtn.style.cssText = "display:block; margin-top:8px;";

        // Dùng pointerup để tránh conflict với document pointerdown
        createBtn.addEventListener("pointerdown", (ev) => {
          ev.stopPropagation(); // Ngăn document pointerdown handler
        });

        createBtn.addEventListener("click", async (ev) => {
          ev.stopPropagation();
          console.info("[Toolbar] ⏳ Đang gửi...");
          createBtn.textContent = "⏳ Đang gửi...";
          createBtn.setAttribute("button-style", "outline");

          let markdown;
          try {
            markdown = buildMarkdown(pinnedEl);
            console.info("[Toolbar] Markdown OK, length:", markdown.length);
          } catch (err) {
            console.error("[Toolbar] buildMarkdown error:", err);
            createBtn.textContent = "❌ Lỗi build";
            createBtn.setAttribute("button-style", "red");
            return;
          }

          try {
            const res = await fetch("http://localhost:8080/api/issues", {
              method: "POST",
              headers: { "Content-Type": "application/json" },
              credentials: "include",
              body: JSON.stringify({
                title: `[Inspect] <${pinnedEl.tagName.toLowerCase()}> tại ${window.location.pathname}`,
                description: markdown,
                priority: "medium",
              }),
            });
            console.info("[Toolbar] Response:", res.status);
            if (res.ok) {
              createBtn.textContent = "✅ Đã tạo";
              eventTarget.dispatchEvent(new CustomEvent("app-notification", {
                detail: { state: "success", title: "Tạo yêu cầu thành công", description: "Issue đã được tạo." },
              }));
            } else {
              const body = await res.json().catch(() => ({}));
              throw new Error(body.detail || body.error || `HTTP ${res.status}`);
            }
          } catch (fetchErr) {
            console.error("[Toolbar] Fetch error:", fetchErr);
            createBtn.textContent = "❌ Lỗi — thử lại";
            createBtn.setAttribute("button-style", "red");
            eventTarget.dispatchEvent(new CustomEvent("app-notification", {
              detail: { state: "error", title: "Lỗi", description: String(fetchErr.message) },
            }));
          }
        });

        infoPanel.insertAdjacentElement("afterend", createBtn);

        eventTarget.dispatchEvent(new CustomEvent("app-notification", {
          detail: { state: "success", title: "Pinned", description: "Click again to unpin." },
        }));
      }
    });
  },
};

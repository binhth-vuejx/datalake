<svelte:options customElement="persistent-widget" />

<script>
  // ── State ────────────────────────────────────────────────────────────────
  let isOpen = $state(false);
  let inspectActive = $state(false);
  let pinned = $state(false);
  let infoText = $state("Hover over an element to inspect it.");
  let pinnedEl = $state(null);
  let sendState = $state("idle"); // idle | sending | ok | error

  // ── Highlight overlay (appended to document.body, outside shadow DOM) ───
  let highlight = null;

  function createHighlight() {
    if (highlight) return;
    highlight = document.createElement("div");
    highlight.style.cssText = `
      position: fixed;
      pointer-events: none;
      z-index: 49;
      outline: 2px dashed #7c3aed;
      background: rgba(124,58,237,0.1);
      border-radius: 2px;
      transition: top 0.05s, left 0.05s, width 0.05s, height 0.05s;
    `;
    document.body.appendChild(highlight);
  }

  function moveHighlight(el, isPinned) {
    if (!highlight) return;
    const rect = el.getBoundingClientRect();
    highlight.style.top    = rect.top + "px";
    highlight.style.left   = rect.left + "px";
    highlight.style.width  = rect.width + "px";
    highlight.style.height = rect.height + "px";
    highlight.style.outlineStyle = isPinned ? "solid" : "dashed";
    highlight.style.outlineColor = isPinned ? "#f59e0b" : "#7c3aed";
    highlight.style.background   = isPinned
      ? "rgba(245,158,11,0.1)"
      : "rgba(124,58,237,0.1)";
  }

  function removeHighlight() {
    if (highlight) { highlight.remove(); highlight = null; }
  }

  // ── Source detection helpers ─────────────────────────────────────────────
  function getVueFile(el) {
    function findVueInst(node) {
      const keys = ["__vueParentComponent", "__vnode", "_vnode", "__vue_app__"];
      for (const k of keys) {
        if (node[k]) {
          if (node[k].component) return node[k].component;
          if (node[k].type) return node[k];
        }
      }
      for (const k of Object.keys(node)) {
        if (k.startsWith("__v") && node[k]?.component) return node[k].component;
        if (k.startsWith("__v") && node[k]?.type?.__file) return node[k];
      }
      return null;
    }
    let node = el;
    while (node && node !== document.body) {
      let inst = findVueInst(node);
      let depth = 0;
      while (inst && depth < 20) {
        const file = inst.type?.__file;
        if (file) return file.replace(/.*\/src\//, "src/");
        inst = inst.parent;
        depth++;
      }
      node = node.parentElement;
    }
    return null;
  }

  function buildIslandMap() {
    const map = new Map();
    document.querySelectorAll("astro-island").forEach((island) => {
      const url = island.getAttribute("component-url") || "";
      const file = url.replace(/.*\/src\//, "src/").replace(/\?.*$/, "");
      island.querySelectorAll("*").forEach((c) => map.set(c, file));
      map.set(island, file);
    });
    return map;
  }

  let islandMap = buildIslandMap();
  new MutationObserver(() => { islandMap = buildIslandMap(); })
    .observe(document.body, { childList: true, subtree: true });

  function getAstroSource(el) {
    const islandFile = islandMap.get(el);
    if (islandFile) return { file: islandFile, type: "island" };
    let node = el.parentElement;
    while (node && node !== document.documentElement) {
      if (node.tagName?.toLowerCase() === "astro-island") {
        const url = node.getAttribute("component-url") || "";
        const file = url.replace(/.*\/src\//, "src/").replace(/\?.*$/, "");
        if (file) return { file, type: "island" };
      }
      const srcFile = node.getAttribute?.("data-astro-source-file");
      if (srcFile) return { file: srcFile.replace(/.*\/src\//, "src/"), type: "astro" };
      node = node.parentElement;
    }
    const p = window.location.pathname;
    const name = p === "/" ? "index" : p.replace(/^\//, "").replace(/\/$/, "");
    return { file: `src/pages/${name}.astro`, type: "page" };
  }

  function getElementInfo(el) {
    const tag = el.tagName.toLowerCase();
    const id = el.id ? `#${el.id}` : "";
    const cls = el.classList.length ? "." + Array.from(el.classList).join(".") : "";
    const rect = el.getBoundingClientRect();
    const p = window.location.pathname;
    const name = p === "/" ? "index" : p.replace(/^\//, "").replace(/\/$/, "");

    const lines = [
      `Tag:   <${tag}${id}${cls}>`,
      `Size:  ${Math.round(rect.width)}×${Math.round(rect.height)}px`,
      `Pos:   top ${Math.round(rect.top)}  left ${Math.round(rect.left)}`,
      `─────────────────────────`,
      `📄 Page: src/pages/${name}.astro`,
    ];

    const astro = getAstroSource(el);
    if (astro.type === "island") lines.push(`🟣 Island: ${astro.file}`);

    const vueFile = getVueFile(el);
    if (vueFile) lines.push(`🟢 Vue: ${vueFile}`);

    return lines.join("\n");
  }

  function buildMarkdown(el) {
    const tag = el.tagName.toLowerCase();
    const id = el.id ? `#${el.id}` : "";
    const cls = el.classList.length ? Array.from(el.classList).join(", ") : "_none_";
    const rect = el.getBoundingClientRect();
    const p = window.location.pathname;
    const name = p === "/" ? "index" : p.replace(/^\//, "").replace(/\/$/, "");
    const pageFile = `src/pages/${name}.astro`;
    const astro = getAstroSource(el);
    const vueFile = getVueFile(el);

    const lines = [
      `## 🐛 Yêu cầu từ Inspect`,
      ``,
      `### Thông tin element`,
      `| Trường | Giá trị |`,
      `|--------|---------|`,
      `| Tag | \`<${tag}${id}>\` |`,
      `| Classes | \`${cls}\` |`,
      `| Size | ${Math.round(rect.width)}×${Math.round(rect.height)}px |`,
      `| Position | top: ${Math.round(rect.top)}, left: ${Math.round(rect.left)} |`,
      `| URL | ${window.location.href} |`,
      ``,
      `### Vị trí trong code`,
      `| Layer | File |`,
      `|-------|------|`,
      `| 📄 Page | \`${pageFile}\` |`,
    ];
    if (astro.type === "island") lines.push(`| 🟣 Island | \`${astro.file}\` |`);
    if (vueFile) lines.push(`| 🟢 Vue Component | \`${vueFile}\` |`);
    lines.push(``, `### Mô tả yêu cầu`, `> _(Điền mô tả yêu cầu tại đây)_`);
    return lines.join("\n");
  }

  // ── Inspect toggle ───────────────────────────────────────────────────────
  function toggleInspect() {
    inspectActive = !inspectActive;
    if (!inspectActive) {
      pinned = false;
      pinnedEl = null;
      sendState = "idle";
      infoText = "Hover over an element to inspect it.";
      removeHighlight();
    }
  }

  // ── Send issue ───────────────────────────────────────────────────────────
  async function sendIssue() {
    if (!pinnedEl) return;
    sendState = "sending";
    try {
      const markdown = buildMarkdown(pinnedEl);
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
      if (res.ok) {
        sendState = "ok";
      } else {
        const body = await res.json().catch(() => ({}));
        throw new Error(body.detail || body.error || `HTTP ${res.status}`);
      }
    } catch (e) {
      console.error("[Widget] Send error:", e);
      sendState = "error";
    }
  }

  // ── DOM event listeners (attached to document, not shadow DOM) ───────────
  $effect(() => {
    function onMouseOver(e) {
      if (!inspectActive || pinned) return;
      const el = e.target;
      if (!el || el === document.body || el === document.documentElement) return;
      // Skip our own widget
      if (el.closest?.("persistent-widget")) return;
      createHighlight();
      moveHighlight(el, false);
      infoText = getElementInfo(el);
    }

    function onMouseOut() {
      if (!inspectActive || pinned) return;
      if (highlight) {
        highlight.style.outline = "none";
        highlight.style.background = "transparent";
      }
    }

    function onPointerDown(e) {
      if (!inspectActive) return;
      const el = e.target;
      if (!el || el === document.body || el === document.documentElement) return;
      if (el === highlight) return;
      if (el.tagName?.toLowerCase().startsWith("astro-")) return;
      if (el.closest?.("persistent-widget")) return;

      e.preventDefault();

      if (pinned) {
        pinned = false;
        pinnedEl = null;
        sendState = "idle";
        infoText = "Hover over an element to inspect it.";
        removeHighlight();
      } else {
        pinned = true;
        pinnedEl = el;
        sendState = "idle";
        createHighlight();
        moveHighlight(el, true);
        infoText = "📌 PINNED\n" + getElementInfo(el);
      }
    }

    document.addEventListener("mouseover", onMouseOver);
    document.addEventListener("mouseout", onMouseOut);
    document.addEventListener("pointerdown", onPointerDown);

    return () => {
      document.removeEventListener("mouseover", onMouseOver);
      document.removeEventListener("mouseout", onMouseOut);
      document.removeEventListener("pointerdown", onPointerDown);
      removeHighlight();
    };
  });
</script>

<!-- Widget wrapper — fixed góc dưới phải -->
<div class="widget">
  {#if isOpen}
    <div class="panel">
      <div class="panel-header">
        <span class="title">🔍 Inspect</span>
        <button class="close-btn" onclick={() => isOpen = false} aria-label="Close">✕</button>
      </div>

      <!-- Inspect toggle -->
      <button
        class="inspect-btn"
        class:active={inspectActive}
        onclick={toggleInspect}
      >
        {inspectActive ? "✋ Stop Inspect" : "🔍 Start Inspect"}
      </button>

      <!-- Info output -->
      <pre class="info">{infoText}</pre>

      <!-- Send issue button — chỉ hiện khi đã pin -->
      {#if pinned}
        <button
          class="send-btn"
          class:sending={sendState === "sending"}
          class:ok={sendState === "ok"}
          class:error={sendState === "error"}
          onclick={sendIssue}
          disabled={sendState === "sending" || sendState === "ok"}
        >
          {#if sendState === "idle"}📝 Tạo yêu cầu
          {:else if sendState === "sending"}⏳ Đang gửi...
          {:else if sendState === "ok"}✅ Đã tạo
          {:else}❌ Lỗi — thử lại{/if}
        </button>
      {/if}
    </div>
  {/if}

  <!-- FAB button -->
  <button class="fab" onclick={() => isOpen = !isOpen} aria-label="Toggle widget">
    {isOpen ? "✕" : "�"}
  </button>
</div>

<style>
  .widget {
    position: fixed;
    bottom: 24px;
    right: 24px;
    z-index: 499;
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 8px;
    font-family: system-ui, sans-serif;
  }

  /* ── Panel ── */
  .panel {
    background: #1e1e2e;
    color: #cdd6f4;
    border: 1px solid #45475a;
    border-radius: 12px;
    padding: 14px 16px;
    width: 300px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
    font-size: 13px;
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .panel-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .title {
    font-weight: 600;
    font-size: 14px;
  }

  .close-btn {
    background: none;
    border: none;
    color: #6c7086;
    cursor: pointer;
    font-size: 14px;
    padding: 0;
    line-height: 1;
  }

  .close-btn:hover { color: #cdd6f4; }

  /* ── Inspect button ── */
  .inspect-btn {
    padding: 6px 12px;
    border-radius: 6px;
    border: none;
    font-size: 12px;
    font-weight: 600;
    cursor: pointer;
    background: #7c3aed;
    color: #fff;
    transition: background 0.15s;
    text-align: left;
  }

  .inspect-btn.active {
    background: #dc2626;
  }

  .inspect-btn:hover {
    filter: brightness(1.15);
  }

  /* ── Info pre ── */
  .info {
    margin: 0;
    font-family: monospace;
    font-size: 11px;
    color: #a6adc8;
    background: rgba(255, 255, 255, 0.04);
    border-radius: 6px;
    padding: 8px;
    white-space: pre-wrap;
    word-break: break-all;
    min-height: 56px;
    line-height: 1.6;
  }

  /* ── Send button ── */
  .send-btn {
    padding: 6px 12px;
    border-radius: 6px;
    border: none;
    font-size: 12px;
    font-weight: 600;
    cursor: pointer;
    background: #7c3aed;
    color: #fff;
    transition: background 0.15s;
  }

  .send-btn.sending { background: #6c7086; cursor: default; }
  .send-btn.ok      { background: #16a34a; cursor: default; }
  .send-btn.error   { background: #dc2626; }
  .send-btn:disabled { opacity: 0.8; }

  /* ── FAB ── */
  .fab {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    background: #89b4fa;
    color: #1e1e2e;
    border: none;
    font-size: 20px;
    cursor: pointer;
    box-shadow: 0 4px 16px rgba(137, 180, 250, 0.4);
    transition: transform 0.15s ease, box-shadow 0.15s ease;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .fab:hover {
    transform: scale(1.1);
    box-shadow: 0 6px 20px rgba(137, 180, 250, 0.6);
  }
</style>

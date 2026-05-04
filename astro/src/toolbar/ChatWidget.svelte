<svelte:options customElement="chat-widget" />

<script>
  // ── State ────────────────────────────────────────────────────────────────
  let isOpen = $state(false);
  let isExpanded = $state(false);
  let input = $state("");
  let sending = $state(false);
  let messages = $state([]);
  let selectedAgent = $state(null);
  let agents = $state([]);
  let agentDropdownOpen = $state(false);
  let workspaceSlug = $state(null);

  // ── Derive workspace slug from URL ───────────────────────────────────────
  function getWorkspaceSlug() {
    const parts = window.location.pathname.split("/").filter(Boolean);
    if (parts[0] === "agents" || parts[0] === "admin") return parts[1] ?? null;
    if (parts[0] === "auth") return null;
    return parts[0] ?? null;
  }

  // ── Fetch agents ─────────────────────────────────────────────────────────
  async function fetchAgents() {
    try {
      const res = await fetch("/api/agents", { credentials: "include" });
      if (!res.ok) return;
      const data = await res.json();
      agents = data?.agents ?? data ?? [];
      if (agents.length > 0 && !selectedAgent) {
        selectedAgent = agents[0];
      }
    } catch {}
  }

  // ── Send message ─────────────────────────────────────────────────────────
  async function sendMessage() {
    const text = input.trim();
    if (!text || sending) return;

    const slug = workspaceSlug ?? getWorkspaceSlug();
    if (!slug) {
      messages = [...messages, { role: "system", content: "⚠️ Không tìm thấy workspace. Hãy mở một workspace trước." }];
      return;
    }

    messages = [...messages, { role: "user", content: text }];
    input = "";
    sending = true;

    try {
      // Create issue with agent task
      const body = {
        title: text.length > 80 ? text.slice(0, 80) + "..." : text,
        description: text,
        priority: "medium",
      };
      if (selectedAgent) body.agent_id = selectedAgent.id;

      const res = await fetch("/api/issues", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify(body),
      });

      if (res.ok) {
        const issue = await res.json();
        messages = [...messages, {
          role: "assistant",
          content: `✅ Đã tạo task: **${issue.title}**\n\nAgent sẽ xử lý ngay.`,
          issueId: issue.id,
        }];
      } else {
        const err = await res.json().catch(() => ({}));
        throw new Error(err.error || `HTTP ${res.status}`);
      }
    } catch (e) {
      messages = [...messages, { role: "system", content: `❌ Lỗi: ${e.message}` }];
    } finally {
      sending = false;
    }
  }

  function handleKeydown(e) {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  }

  function toggleOpen() {
    isOpen = !isOpen;
    if (isOpen && agents.length === 0) {
      workspaceSlug = getWorkspaceSlug();
      fetchAgents();
    }
  }

  function clearMessages() {
    messages = [];
  }
</script>

<div class="chat-widget">
  {#if isOpen}
    <div class="panel" class:expanded={isExpanded}>
      <!-- Header -->
      <div class="header">
        <div class="header-left">
          <span class="bot-icon">🤖</span>
          <span class="title">
            {selectedAgent?.name ?? "Ask Multica"}
          </span>
        </div>
        <div class="header-actions">
          <button class="icon-btn" onclick={() => isExpanded = !isExpanded} title={isExpanded ? "Thu nhỏ" : "Mở rộng"}>
            {isExpanded ? "⊡" : "⊞"}
          </button>
          <button class="icon-btn" onclick={clearMessages} title="Xóa lịch sử">🗑</button>
          <button class="icon-btn" onclick={() => isOpen = false} title="Đóng">✕</button>
        </div>
      </div>

      <!-- Messages -->
      <div class="messages">
        {#if messages.length === 0}
          <div class="empty">
            <span class="empty-icon">💬</span>
            <p>Hỏi Multica bất cứ điều gì hoặc giao task cho agent.</p>
          </div>
        {:else}
          {#each messages as msg}
            <div class="message {msg.role}">
              {#if msg.role === "user"}
                <div class="bubble user-bubble">{msg.content}</div>
              {:else if msg.role === "assistant"}
                <div class="bubble assistant-bubble">{msg.content}</div>
              {:else}
                <div class="bubble system-bubble">{msg.content}</div>
              {/if}
            </div>
          {/each}
          {#if sending}
            <div class="message assistant">
              <div class="bubble assistant-bubble typing">
                <span></span><span></span><span></span>
              </div>
            </div>
          {/if}
        {/if}
      </div>

      <!-- Input area -->
      <div class="input-area">
        <textarea
          bind:value={input}
          onkeydown={handleKeydown}
          placeholder="Tell me what to do..."
          rows="2"
          disabled={sending}
        ></textarea>
        <div class="input-footer">
          <!-- Agent selector -->
          <div class="agent-selector">
            <button class="agent-btn" onclick={() => agentDropdownOpen = !agentDropdownOpen}>
              <span class="agent-icon">🤖</span>
              <span class="agent-name">{selectedAgent?.name ?? "Agent"}</span>
              <span class="chevron">▾</span>
            </button>
            {#if agentDropdownOpen && agents.length > 0}
              <div class="agent-dropdown">
                {#each agents as agent}
                  <button
                    class="agent-option"
                    class:selected={selectedAgent?.id === agent.id}
                    onclick={() => { selectedAgent = agent; agentDropdownOpen = false; }}
                  >
                    🤖 {agent.name}
                  </button>
                {/each}
              </div>
            {/if}
          </div>

          <!-- Send button -->
          <button
            class="send-btn"
            onclick={sendMessage}
            disabled={sending || !input.trim()}
            title="Gửi (Enter)"
          >
            ↑
          </button>
        </div>
      </div>
    </div>
  {/if}

  <!-- FAB -->
  <button class="fab" onclick={toggleOpen} aria-label="Chat with Multica">
    {#if isOpen}
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
        <line x1="18" y1="6" x2="6" y2="18"></line>
        <line x1="6" y1="6" x2="18" y2="18"></line>
      </svg>
    {:else}
      <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
      </svg>
    {/if}
  </button>
</div>

<style>
  .chat-widget {
    position: fixed;
    bottom: 24px;
    right: 80px; /* offset from persistent-widget */
    z-index: 49;
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 8px;
    font-family: system-ui, -apple-system, sans-serif;
  }

  /* ── Panel ── */
  .panel {
    background: #fff;
    border: 1px solid #e5e7eb;
    border-radius: 16px;
    width: 340px;
    height: 480px;
    box-shadow: 0 8px 40px rgba(0,0,0,0.15);
    display: flex;
    flex-direction: column;
    overflow: hidden;
    transition: width 0.2s, height 0.2s;
  }

  .panel.expanded {
    width: 480px;
    height: 640px;
  }

  /* ── Header ── */
  .header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 14px;
    border-bottom: 1px solid #f3f4f6;
    background: #fafafa;
    border-radius: 16px 16px 0 0;
  }

  .header-left {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .bot-icon {
    font-size: 18px;
  }

  .title {
    font-weight: 600;
    font-size: 14px;
    color: #111827;
    max-width: 160px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .header-actions {
    display: flex;
    gap: 4px;
  }

  .icon-btn {
    background: none;
    border: none;
    cursor: pointer;
    color: #9ca3af;
    font-size: 14px;
    padding: 4px 6px;
    border-radius: 6px;
    line-height: 1;
    transition: background 0.1s, color 0.1s;
  }

  .icon-btn:hover {
    background: #f3f4f6;
    color: #374151;
  }

  /* ── Messages ── */
  .messages {
    flex: 1;
    overflow-y: auto;
    padding: 12px;
    display: flex;
    flex-direction: column;
    gap: 8px;
    scroll-behavior: smooth;
  }

  .empty {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    color: #9ca3af;
    text-align: center;
    gap: 8px;
  }

  .empty-icon {
    font-size: 32px;
  }

  .empty p {
    font-size: 13px;
    margin: 0;
    line-height: 1.5;
  }

  .message {
    display: flex;
  }

  .message.user {
    justify-content: flex-end;
  }

  .bubble {
    max-width: 85%;
    padding: 8px 12px;
    border-radius: 12px;
    font-size: 13px;
    line-height: 1.5;
    white-space: pre-wrap;
    word-break: break-word;
  }

  .user-bubble {
    background: #7c3aed;
    color: #fff;
    border-bottom-right-radius: 4px;
  }

  .assistant-bubble {
    background: #f3f4f6;
    color: #111827;
    border-bottom-left-radius: 4px;
  }

  .system-bubble {
    background: #fef3c7;
    color: #92400e;
    border-radius: 8px;
    font-size: 12px;
  }

  /* Typing indicator */
  .typing {
    display: flex;
    gap: 4px;
    align-items: center;
    padding: 10px 14px;
  }

  .typing span {
    width: 6px;
    height: 6px;
    background: #9ca3af;
    border-radius: 50%;
    animation: bounce 1.2s infinite;
  }

  .typing span:nth-child(2) { animation-delay: 0.2s; }
  .typing span:nth-child(3) { animation-delay: 0.4s; }

  @keyframes bounce {
    0%, 60%, 100% { transform: translateY(0); }
    30% { transform: translateY(-6px); }
  }

  /* ── Input area ── */
  .input-area {
    border-top: 1px solid #f3f4f6;
    padding: 10px 12px;
    background: #fff;
    border-radius: 0 0 16px 16px;
  }

  textarea {
    width: 100%;
    border: 1px solid #e5e7eb;
    border-radius: 10px;
    padding: 8px 10px;
    font-size: 13px;
    font-family: inherit;
    resize: none;
    outline: none;
    color: #111827;
    background: #f9fafb;
    box-sizing: border-box;
    transition: border-color 0.15s;
    line-height: 1.5;
  }

  textarea:focus {
    border-color: #7c3aed;
    background: #fff;
  }

  textarea::placeholder {
    color: #9ca3af;
  }

  textarea:disabled {
    opacity: 0.6;
  }

  .input-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-top: 8px;
  }

  /* ── Agent selector ── */
  .agent-selector {
    position: relative;
  }

  .agent-btn {
    display: flex;
    align-items: center;
    gap: 4px;
    background: none;
    border: 1px solid #e5e7eb;
    border-radius: 20px;
    padding: 4px 10px;
    font-size: 12px;
    cursor: pointer;
    color: #374151;
    transition: border-color 0.15s;
  }

  .agent-btn:hover {
    border-color: #7c3aed;
    color: #7c3aed;
  }

  .agent-icon { font-size: 14px; }
  .agent-name { max-width: 100px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .chevron { font-size: 10px; color: #9ca3af; }

  .agent-dropdown {
    position: absolute;
    bottom: calc(100% + 4px);
    left: 0;
    background: #fff;
    border: 1px solid #e5e7eb;
    border-radius: 10px;
    box-shadow: 0 4px 16px rgba(0,0,0,0.1);
    min-width: 160px;
    overflow: hidden;
    z-index: 10;
  }

  .agent-option {
    display: block;
    width: 100%;
    text-align: left;
    padding: 8px 12px;
    font-size: 12px;
    background: none;
    border: none;
    cursor: pointer;
    color: #374151;
    transition: background 0.1s;
  }

  .agent-option:hover { background: #f3f4f6; }
  .agent-option.selected { background: #ede9fe; color: #7c3aed; font-weight: 600; }

  /* ── Send button ── */
  .send-btn {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: #7c3aed;
    color: #fff;
    border: none;
    font-size: 16px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background 0.15s, transform 0.1s;
    flex-shrink: 0;
  }

  .send-btn:hover:not(:disabled) {
    background: #6d28d9;
    transform: scale(1.05);
  }

  .send-btn:disabled {
    background: #d1d5db;
    cursor: default;
  }

  /* ── FAB ── */
  .fab {
    width: 52px;
    height: 52px;
    border-radius: 50%;
    background: #fff;
    color: #374151;
    border: 1px solid #e5e7eb;
    cursor: pointer;
    box-shadow: 0 4px 16px rgba(0,0,0,0.12);
    transition: transform 0.15s, box-shadow 0.15s;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .fab:hover {
    transform: scale(1.08);
    box-shadow: 0 6px 20px rgba(0,0,0,0.18);
  }
</style>

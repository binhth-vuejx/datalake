<script>
  let isOpen = $state(false);

  function toggle() {
    isOpen = !isOpen;
  }
</script>

<!-- Floating container, fixed góc dưới phải -->
<div class="floating-wrapper">
  {#if isOpen}
    <div class="floating-panel">
      <p>👋 Widget content ở đây</p>
      <p>Component này hoàn toàn độc lập với Vue.</p>
    </div>
  {/if}

  <button class="floating-btn" onclick={toggle} aria-label="Toggle widget">
    {isOpen ? "✕" : "💬"}
  </button>
</div>

<style>
  .floating-wrapper {
    position: fixed;
    bottom: 24px;
    right: 24px;
    z-index: 9999;
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 8px;
  }

  .floating-panel {
    background: #1e1e2e;
    color: #cdd6f4;
    border: 1px solid #45475a;
    border-radius: 12px;
    padding: 16px 20px;
    width: 240px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
    font-size: 14px;
    line-height: 1.6;
  }

  .floating-btn {
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
  }

  .floating-btn:hover {
    transform: scale(1.1);
    box-shadow: 0 6px 20px rgba(137, 180, 250, 0.6);
  }
</style>

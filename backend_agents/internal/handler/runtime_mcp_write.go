package handler

import (
	"sync"
)

// McpWriteStore holds pending MCP config write requests for daemon runtimes.
// When the user saves MCP config via the UI, the server queues a write request.
// The daemon picks it up on the next heartbeat and writes to the local config file.
type McpWriteStore struct {
	mu      sync.Mutex
	pending map[string]string // runtimeID -> mcpConfig JSON (empty string = clear)
}

func NewMcpWriteStore() *McpWriteStore {
	return &McpWriteStore{
		pending: make(map[string]string),
	}
}

// Queue stores a pending MCP config write for a runtime.
// mcpConfig is the full JSON string (or empty to clear).
func (s *McpWriteStore) Queue(runtimeID, mcpConfig string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pending[runtimeID] = mcpConfig
}

// PopPending returns and removes the pending write for a runtime, if any.
// Returns ("", false) if nothing is pending.
func (s *McpWriteStore) PopPending(runtimeID string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.pending[runtimeID]
	if ok {
		delete(s.pending, runtimeID)
	}
	return v, ok
}

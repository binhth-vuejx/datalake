package commands

import (
	"encoding/json"
	"time"
)

// Command represents a write operation
type Command struct {
	Name      string          `json:"name"`
	Input     json.RawMessage `json:"input"`
	ExecutedBy string          `json:"executed_by"`
	ExecutedAt time.Time       `json:"executed_at"`
}

// CommandResult represents the result of a command execution
type CommandResult struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}

// CommandHandler defines the interface for command handlers
type CommandHandler interface {
	Handle(input json.RawMessage) (json.RawMessage, error)
}

// CommandRegistry manages all available commands
type CommandRegistry struct {
	handlers map[string]CommandHandler
}

// NewCommandRegistry creates a new command registry
func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		handlers: make(map[string]CommandHandler),
	}
}

// Register registers a command handler
func (r *CommandRegistry) Register(name string, handler CommandHandler) {
	r.handlers[name] = handler
}

// Execute executes a command
func (r *CommandRegistry) Execute(name string, input json.RawMessage) (json.RawMessage, error) {
	handler, ok := r.handlers[name]
	if !ok {
		return nil, ErrCommandNotFound
	}
	return handler.Handle(input)
}

// GetHandler returns a handler for a command
func (r *CommandRegistry) GetHandler(name string) (CommandHandler, error) {
	handler, ok := r.handlers[name]
	if !ok {
		return nil, ErrCommandNotFound
	}
	return handler, nil
}

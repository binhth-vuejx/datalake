package commands

import "errors"

var (
	ErrCommandNotFound      = errors.New("command not found")
	ErrInvalidInput         = errors.New("invalid command input")
	ErrExecutionFailed      = errors.New("command execution failed")
	ErrValidationFailed     = errors.New("validation failed")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrInternalServerError  = errors.New("internal server error")
)

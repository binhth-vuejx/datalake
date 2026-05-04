package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/datalake-free/internal/domain/customer"
	"github.com/yourorg/datalake-free/internal/infrastructure/postgres"
)

// UpdateCustomerHandler handles customer updates with audit logging
type UpdateCustomerHandler struct {
	repo        customer.Repository
	auditLogger customer.AuditLogger
	db          *pgxpool.Pool
}

func NewUpdateCustomerHandler(
	repo customer.Repository,
	auditLogger customer.AuditLogger,
	db *pgxpool.Pool,
) *UpdateCustomerHandler {
	return &UpdateCustomerHandler{
		repo:        repo,
		auditLogger: auditLogger,
		db:          db,
	}
}

// Handle updates an existing customer with transaction and audit logging
func (h *UpdateCustomerHandler) Handle(ctx context.Context, cmd customer.UpdateCustomerCommand) (*customer.Customer, error) {
	// Validate input
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Execute within transaction
	var result *customer.Customer
	err := postgres.WithTx(ctx, h.db, func(txCtx context.Context) error {
		// Get existing customer
		oldCustomer, err := h.repo.GetByID(txCtx, cmd.ID)
		if err != nil {
			return fmt.Errorf("get customer: %w", err)
		}
		if oldCustomer == nil {
			return fmt.Errorf("customer not found: %d", cmd.ID)
		}

		// Check if new email is already used by another customer
		if cmd.Email != oldCustomer.Email {
			existing, err := h.repo.GetByEmail(txCtx, cmd.Email)
			if err == nil && existing != nil && existing.ID != cmd.ID {
				return fmt.Errorf("email already exists: %s", cmd.Email)
			}
		}

		// Update customer entity
		updatedCustomer := &customer.Customer{
			ID:        cmd.ID,
			Name:      cmd.Name,
			Email:     cmd.Email,
			Phone:     cmd.Phone,
			CreatedBy: oldCustomer.CreatedBy,
			CreatedAt: oldCustomer.CreatedAt,
			UpdatedAt: time.Now(),
		}

		// Save updated customer
		if err := h.repo.Update(txCtx, updatedCustomer); err != nil {
			return fmt.Errorf("update customer: %w", err)
		}

		// Log audit
		if err := h.auditLogger.LogUpdate(txCtx, oldCustomer, updatedCustomer, cmd.UpdatedBy); err != nil {
			return fmt.Errorf("log audit: %w", err)
		}

		result = updatedCustomer
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

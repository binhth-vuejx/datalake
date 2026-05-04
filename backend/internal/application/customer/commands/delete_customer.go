package commands

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/datalake-free/internal/domain/customer"
	"github.com/yourorg/datalake-free/internal/infrastructure/postgres"
)

// DeleteCustomerHandler handles customer deletion with audit logging
type DeleteCustomerHandler struct {
	repo        customer.Repository
	auditLogger customer.AuditLogger
	db          *pgxpool.Pool
}

func NewDeleteCustomerHandler(
	repo customer.Repository,
	auditLogger customer.AuditLogger,
	db *pgxpool.Pool,
) *DeleteCustomerHandler {
	return &DeleteCustomerHandler{
		repo:        repo,
		auditLogger: auditLogger,
		db:          db,
	}
}

// Handle deletes a customer with transaction and audit logging
func (h *DeleteCustomerHandler) Handle(ctx context.Context, cmd customer.DeleteCustomerCommand) error {
	// Validate input
	if err := cmd.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	// Execute within transaction
	err := postgres.WithTx(ctx, h.db, func(txCtx context.Context) error {
		// Get customer before deletion
		customer, err := h.repo.GetByID(txCtx, cmd.ID)
		if err != nil {
			return fmt.Errorf("get customer: %w", err)
		}
		if customer == nil {
			return fmt.Errorf("customer not found: %d", cmd.ID)
		}

		// Delete customer
		if err := h.repo.Delete(txCtx, cmd.ID); err != nil {
			return fmt.Errorf("delete customer: %w", err)
		}

		// Log audit
		if err := h.auditLogger.LogDelete(txCtx, customer, cmd.DeletedBy); err != nil {
			return fmt.Errorf("log audit: %w", err)
		}

		return nil
	})

	return err
}

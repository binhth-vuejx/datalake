package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/datalake-free/internal/domain/customer"
	"github.com/yourorg/datalake-free/internal/infrastructure/postgres"
)

// CreateCustomerHandler handles customer creation with audit logging
type CreateCustomerHandler struct {
	repo        customer.Repository
	auditLogger customer.AuditLogger
	db          *pgxpool.Pool
}

func NewCreateCustomerHandler(
	repo customer.Repository,
	auditLogger customer.AuditLogger,
	db *pgxpool.Pool,
) *CreateCustomerHandler {
	return &CreateCustomerHandler{
		repo:        repo,
		auditLogger: auditLogger,
		db:          db,
	}
}

// Handle creates a new customer with transaction and audit logging
func (h *CreateCustomerHandler) Handle(ctx context.Context, cmd customer.CreateCustomerCommand) (*customer.Customer, error) {
	// Validate input
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Execute within transaction
	var result *customer.Customer
	err := postgres.WithTx(ctx, h.db, func(txCtx context.Context) error {
		// Create customer entity
		now := time.Now()
		newCustomer := &customer.Customer{
			Name:      cmd.Name,
			Email:     cmd.Email,
			Phone:     cmd.Phone,
			CreatedBy: cmd.CreatedBy,
			CreatedAt: now,
			UpdatedAt: now,
		}

		// Check if email already exists
		existing, err := h.repo.GetByEmail(txCtx, cmd.Email)
		if err == nil && existing != nil {
			return fmt.Errorf("email already exists: %s", cmd.Email)
		}

		// Save customer
		if err := h.repo.Create(txCtx, newCustomer); err != nil {
			return fmt.Errorf("create customer: %w", err)
		}

		// Log audit
		if err := h.auditLogger.LogCreate(txCtx, newCustomer, cmd.CreatedBy); err != nil {
			return fmt.Errorf("log audit: %w", err)
		}

		result = newCustomer
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

package commands

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/datalake-free/internal/domain/sale"
)

type UpdateSaleHandler struct {
	repo        sale.Repository
	auditLogger interface {
		LogUpdateSale(ctx context.Context, oldSale, newSale *sale.Sale, updatedBy string) error
	}
	db *pgxpool.Pool
}

func NewUpdateSaleHandler(repo sale.Repository, auditLogger interface {
	LogUpdateSale(ctx context.Context, oldSale, newSale *sale.Sale, updatedBy string) error
}, db *pgxpool.Pool) *UpdateSaleHandler {
	return &UpdateSaleHandler{
		repo:        repo,
		auditLogger: auditLogger,
		db:          db,
	}
}

func (h *UpdateSaleHandler) Handle(ctx context.Context, cmd sale.UpdateSaleCommand) (*sale.Sale, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	// Retrieve existing sale for audit comparison
	oldSale, err := h.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	// Create updated sale with new values
	newSale := &sale.Sale{
		ID:         oldSale.ID,
		CustomerID: oldSale.CustomerID,
		Amount:     cmd.Amount,
		Status:     cmd.Status,
		CreatedBy:  oldSale.CreatedBy,
		CreatedAt:  oldSale.CreatedAt,
		UpdatedAt:  time.Now(),
	}

	// Persist changes
	if err := h.repo.Update(ctx, newSale); err != nil {
		return nil, err
	}

	// Log audit event with old and new data
	if err := h.auditLogger.LogUpdateSale(ctx, oldSale, newSale, cmd.UpdatedBy); err != nil {
		// Log but don't fail
		_ = err
	}

	return newSale, nil
}

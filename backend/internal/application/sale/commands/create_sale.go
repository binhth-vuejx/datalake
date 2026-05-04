package commands

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/datalake-free/internal/domain/sale"
)

type CreateSaleHandler struct {
	repo        sale.Repository
	auditLogger interface {
		LogCreateSale(ctx context.Context, sale *sale.Sale, createdBy string) error
	}
	db *pgxpool.Pool
}

func NewCreateSaleHandler(repo sale.Repository, auditLogger interface {
	LogCreateSale(ctx context.Context, sale *sale.Sale, createdBy string) error
}, db *pgxpool.Pool) *CreateSaleHandler {
	return &CreateSaleHandler{
		repo:        repo,
		auditLogger: auditLogger,
		db:          db,
	}
}

func (h *CreateSaleHandler) Handle(ctx context.Context, cmd sale.CreateSaleCommand) (*sale.Sale, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	s := &sale.Sale{
		CustomerID: cmd.CustomerID,
		Amount:     cmd.Amount,
		Status:     "pending",
		CreatedBy:  cmd.CreatedBy,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := h.repo.Create(ctx, s); err != nil {
		return nil, err
	}

	if err := h.auditLogger.LogCreateSale(ctx, s, cmd.CreatedBy); err != nil {
		// Log but don't fail
		_ = err
	}

	return s, nil
}

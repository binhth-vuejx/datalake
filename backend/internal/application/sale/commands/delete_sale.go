package commands

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/datalake-free/internal/domain/sale"
)

type DeleteSaleHandler struct {
	repo        sale.Repository
	auditLogger interface {
		LogDeleteSale(ctx context.Context, sale *sale.Sale, deletedBy string) error
	}
	db *pgxpool.Pool
}

func NewDeleteSaleHandler(repo sale.Repository, auditLogger interface {
	LogDeleteSale(ctx context.Context, sale *sale.Sale, deletedBy string) error
}, db *pgxpool.Pool) *DeleteSaleHandler {
	return &DeleteSaleHandler{
		repo:        repo,
		auditLogger: auditLogger,
		db:          db,
	}
}

func (h *DeleteSaleHandler) Handle(ctx context.Context, cmd sale.DeleteSaleCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	existing, err := h.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err := h.repo.Delete(ctx, cmd.ID); err != nil {
		return err
	}

	if err := h.auditLogger.LogDeleteSale(ctx, existing, cmd.DeletedBy); err != nil {
		// Log but don't fail
		_ = err
	}

	return nil
}

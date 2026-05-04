package commands

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/datalake-free/internal/domain/order"
)

type DeleteOrderHandler struct {
	repo        order.Repository
	auditLogger interface {
		LogDeleteOrder(ctx context.Context, order *order.Order, deletedBy string) error
	}
	db *pgxpool.Pool
}

func NewDeleteOrderHandler(repo order.Repository, auditLogger interface {
	LogDeleteOrder(ctx context.Context, order *order.Order, deletedBy string) error
}, db *pgxpool.Pool) *DeleteOrderHandler {
	return &DeleteOrderHandler{
		repo:        repo,
		auditLogger: auditLogger,
		db:          db,
	}
}

func (h *DeleteOrderHandler) Handle(ctx context.Context, cmd order.DeleteOrderCommand) error {
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

	if err := h.auditLogger.LogDeleteOrder(ctx, existing, cmd.DeletedBy); err != nil {
		// Log but don't fail
		_ = err
	}

	return nil
}

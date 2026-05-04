package commands

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/datalake-free/internal/domain/order"
)

type UpdateOrderHandler struct {
	repo        order.Repository
	auditLogger interface {
		LogUpdateOrder(ctx context.Context, oldOrder, newOrder *order.Order, updatedBy string) error
	}
	db *pgxpool.Pool
}

func NewUpdateOrderHandler(repo order.Repository, auditLogger interface {
	LogUpdateOrder(ctx context.Context, oldOrder, newOrder *order.Order, updatedBy string) error
}, db *pgxpool.Pool) *UpdateOrderHandler {
	return &UpdateOrderHandler{
		repo:        repo,
		auditLogger: auditLogger,
		db:          db,
	}
}

func (h *UpdateOrderHandler) Handle(ctx context.Context, cmd order.UpdateOrderCommand) (*order.Order, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	existing, err := h.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	oldOrder := *existing

	existing.TotalAmount = cmd.TotalAmount
	existing.Status = cmd.Status
	existing.UpdatedAt = time.Now()

	if err := h.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	if err := h.auditLogger.LogUpdateOrder(ctx, &oldOrder, existing, cmd.UpdatedBy); err != nil {
		// Log but don't fail
		_ = err
	}

	return existing, nil
}

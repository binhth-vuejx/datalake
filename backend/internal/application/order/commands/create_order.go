package commands

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/datalake-free/internal/domain/order"
)

type CreateOrderHandler struct {
	repo        order.Repository
	auditLogger interface {
		LogCreateOrder(ctx context.Context, order *order.Order, createdBy string) error
	}
	db *pgxpool.Pool
}

func NewCreateOrderHandler(repo order.Repository, auditLogger interface {
	LogCreateOrder(ctx context.Context, order *order.Order, createdBy string) error
}, db *pgxpool.Pool) *CreateOrderHandler {
	return &CreateOrderHandler{
		repo:        repo,
		auditLogger: auditLogger,
		db:          db,
	}
}

func (h *CreateOrderHandler) Handle(ctx context.Context, cmd order.CreateOrderCommand) (*order.Order, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	o := &order.Order{
		CustomerID:  cmd.CustomerID,
		TotalAmount: cmd.TotalAmount,
		Status:      "pending",
		CreatedBy:   cmd.CreatedBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.repo.Create(ctx, o); err != nil {
		return nil, err
	}

	if err := h.auditLogger.LogCreateOrder(ctx, o, cmd.CreatedBy); err != nil {
		// Log but don't fail
		_ = err
	}

	return o, nil
}

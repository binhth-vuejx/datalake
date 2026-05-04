package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/datalake-free/internal/domain/order"
)

// OrderRepository implements order.Repository using PostgreSQL
type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create saves a new order
func (r *OrderRepository) Create(ctx context.Context, o *order.Order) error {
	query := `
		INSERT INTO orders (customer_id, total_amount, status, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := r.db.QueryRow(ctx, query,
		o.CustomerID, o.TotalAmount, o.Status, o.CreatedBy, o.CreatedAt, o.UpdatedAt,
	).Scan(&o.ID)

	if err != nil {
		return fmt.Errorf("insert order: %w", err)
	}

	return nil
}

// Update updates an existing order
func (r *OrderRepository) Update(ctx context.Context, o *order.Order) error {
	query := `
		UPDATE orders
		SET total_amount = $1, status = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := r.db.Exec(ctx, query,
		o.TotalAmount, o.Status, o.UpdatedAt, o.ID,
	)

	if err != nil {
		return fmt.Errorf("update order: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("order not found: %d", o.ID)
	}

	return nil
}

// Delete removes an order
func (r *OrderRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM orders WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete order: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("order not found: %d", id)
	}

	return nil
}

// GetByID retrieves an order by ID
func (r *OrderRepository) GetByID(ctx context.Context, id int64) (*order.Order, error) {
	query := `
		SELECT id, customer_id, total_amount, status, created_by, created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	var o order.Order
	err := r.db.QueryRow(ctx, query, id).Scan(
		&o.ID, &o.CustomerID, &o.TotalAmount, &o.Status, &o.CreatedBy, &o.CreatedAt, &o.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}

	return &o, nil
}

// GetByCustomerID retrieves all orders for a customer
func (r *OrderRepository) GetByCustomerID(ctx context.Context, customerID int64) ([]*order.Order, error) {
	query := `
		SELECT id, customer_id, total_amount, status, created_by, created_at, updated_at
		FROM orders
		WHERE customer_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("query orders by customer: %w", err)
	}
	defer rows.Close()

	var orders []*order.Order
	for rows.Next() {
		var o order.Order
		err := rows.Scan(
			&o.ID, &o.CustomerID, &o.TotalAmount, &o.Status, &o.CreatedBy, &o.CreatedAt, &o.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}
		orders = append(orders, &o)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return orders, nil
}

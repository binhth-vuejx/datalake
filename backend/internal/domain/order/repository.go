package order

import "context"

// Repository interface for persistence
type Repository interface {
	// Create saves a new order
	Create(ctx context.Context, order *Order) error

	// Update updates an existing order
	Update(ctx context.Context, order *Order) error

	// Delete removes an order
	Delete(ctx context.Context, id int64) error

	// GetByID retrieves an order by ID
	GetByID(ctx context.Context, id int64) (*Order, error)

	// GetByCustomerID retrieves all orders for a customer
	GetByCustomerID(ctx context.Context, customerID int64) ([]*Order, error)
}

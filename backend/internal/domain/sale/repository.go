package sale

import "context"

// Repository interface for persistence
type Repository interface {
	// Create saves a new sale
	Create(ctx context.Context, sale *Sale) error

	// Update updates an existing sale
	Update(ctx context.Context, sale *Sale) error

	// Delete removes a sale
	Delete(ctx context.Context, id int64) error

	// GetByID retrieves a sale by ID
	GetByID(ctx context.Context, id int64) (*Sale, error)

	// GetByCustomerID retrieves all sales for a customer
	GetByCustomerID(ctx context.Context, customerID int64) ([]*Sale, error)

	// List retrieves sales with pagination support
	// Returns a slice of sales and the total count
	List(ctx context.Context, limit, offset int) ([]*Sale, int, error)
}

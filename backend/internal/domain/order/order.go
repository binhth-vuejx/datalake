package order

import (
	"errors"
	"time"
)

// Order is the domain entity
type Order struct {
	ID          int64
	CustomerID  int64
	TotalAmount float64
	Status      string
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CreateOrderCommand is the input for creating an order
type CreateOrderCommand struct {
	CustomerID  int64
	TotalAmount float64
	CreatedBy   string
}

// UpdateOrderCommand is the input for updating an order
type UpdateOrderCommand struct {
	ID          int64
	TotalAmount float64
	Status      string
	UpdatedBy   string
}

// DeleteOrderCommand is the input for deleting an order
type DeleteOrderCommand struct {
	ID        int64
	DeletedBy string
}

// Validate business rules for create
func (cmd CreateOrderCommand) Validate() error {
	if cmd.CustomerID == 0 {
		return errors.New("customer_id is required")
	}
	if cmd.TotalAmount <= 0 {
		return errors.New("total_amount must be positive")
	}
	if cmd.CreatedBy == "" {
		return errors.New("created_by is required")
	}
	return nil
}

// Validate business rules for update
func (cmd UpdateOrderCommand) Validate() error {
	if cmd.ID == 0 {
		return errors.New("id is required")
	}
	if cmd.TotalAmount < 0 {
		return errors.New("total_amount cannot be negative")
	}
	if cmd.Status == "" {
		return errors.New("status is required")
	}
	if cmd.UpdatedBy == "" {
		return errors.New("updated_by is required")
	}
	return nil
}

// Validate business rules for delete
func (cmd DeleteOrderCommand) Validate() error {
	if cmd.ID == 0 {
		return errors.New("id is required")
	}
	if cmd.DeletedBy == "" {
		return errors.New("deleted_by is required")
	}
	return nil
}

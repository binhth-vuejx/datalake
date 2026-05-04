package sale

import (
	"errors"
	"time"
)

// Sale is the domain entity
type Sale struct {
	ID         int64
	CustomerID int64
	Amount     float64
	Status     string
	CreatedBy  string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// CreateSaleCommand is the input for creating a sale
type CreateSaleCommand struct {
	CustomerID int64
	Amount     float64
	CreatedBy  string
}

// UpdateSaleCommand is the input for updating a sale
type UpdateSaleCommand struct {
	ID        int64
	Amount    float64
	Status    string
	UpdatedBy string
}

// DeleteSaleCommand is the input for deleting a sale
type DeleteSaleCommand struct {
	ID        int64
	DeletedBy string
}

// Validate business rules for create
func (cmd CreateSaleCommand) Validate() error {
	if cmd.CustomerID <= 0 {
		return errors.New("customer_id must be a positive integer")
	}
	if cmd.Amount <= 0 {
		return errors.New("amount must be positive")
	}
	if cmd.CreatedBy == "" {
		return errors.New("created_by is required")
	}
	if len(cmd.CreatedBy) > 255 {
		return errors.New("created_by must not exceed 255 characters")
	}
	return nil
}

// Validate business rules for update
func (cmd UpdateSaleCommand) Validate() error {
	if cmd.ID <= 0 {
		return errors.New("id must be a positive integer")
	}
	if cmd.Amount <= 0 {
		return errors.New("amount must be positive")
	}
	validStatuses := map[string]bool{
		"pending":   true,
		"completed": true,
		"cancelled": true,
	}
	if !validStatuses[cmd.Status] {
		return errors.New("status must be one of: pending, completed, cancelled")
	}
	if cmd.UpdatedBy == "" {
		return errors.New("updated_by is required")
	}
	if len(cmd.UpdatedBy) > 255 {
		return errors.New("updated_by must not exceed 255 characters")
	}
	return nil
}

// Validate business rules for delete
func (cmd DeleteSaleCommand) Validate() error {
	if cmd.ID <= 0 {
		return errors.New("id must be a positive integer")
	}
	if cmd.DeletedBy == "" {
		return errors.New("deleted_by is required")
	}
	if len(cmd.DeletedBy) > 255 {
		return errors.New("deleted_by must not exceed 255 characters")
	}
	return nil
}

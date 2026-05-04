package customer

import (
	"context"
	"errors"
	"time"
)

// Customer is the domain entity
type Customer struct {
	ID        int64
	Name      string
	Email     string
	Phone     string
	CreatedBy string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateCustomerCommand is the input for creating a customer
type CreateCustomerCommand struct {
	Name      string
	Email     string
	Phone     string
	CreatedBy string
}

// UpdateCustomerCommand is the input for updating a customer
type UpdateCustomerCommand struct {
	ID        int64
	Name      string
	Email     string
	Phone     string
	UpdatedBy string
}

// DeleteCustomerCommand is the input for deleting a customer
type DeleteCustomerCommand struct {
	ID        int64
	DeletedBy string
}

// Validate business rules for create
func (cmd CreateCustomerCommand) Validate() error {
	if cmd.Name == "" {
		return errors.New("name is required")
	}
	if cmd.Email == "" {
		return errors.New("email is required")
	}
	if len(cmd.Email) < 5 || len(cmd.Email) > 255 {
		return errors.New("email must be between 5 and 255 characters")
	}
	if cmd.CreatedBy == "" {
		return errors.New("created_by is required")
	}
	return nil
}

// Validate business rules for update
func (cmd UpdateCustomerCommand) Validate() error {
	if cmd.ID == 0 {
		return errors.New("id is required")
	}
	if cmd.Name == "" {
		return errors.New("name is required")
	}
	if cmd.Email == "" {
		return errors.New("email is required")
	}
	if len(cmd.Email) < 5 || len(cmd.Email) > 255 {
		return errors.New("email must be between 5 and 255 characters")
	}
	if cmd.UpdatedBy == "" {
		return errors.New("updated_by is required")
	}
	return nil
}

// Validate business rules for delete
func (cmd DeleteCustomerCommand) Validate() error {
	if cmd.ID == 0 {
		return errors.New("id is required")
	}
	if cmd.DeletedBy == "" {
		return errors.New("deleted_by is required")
	}
	return nil
}

// Repository interface for persistence
type Repository interface {
	// Create saves a new customer
	Create(ctx context.Context, customer *Customer) error

	// Update updates an existing customer
	Update(ctx context.Context, customer *Customer) error

	// Delete removes a customer
	Delete(ctx context.Context, id int64) error

	// GetByID retrieves a customer by ID
	GetByID(ctx context.Context, id int64) (*Customer, error)

	// GetByEmail retrieves a customer by email
	GetByEmail(ctx context.Context, email string) (*Customer, error)
}

// AuditLogger interface for logging changes
type AuditLogger interface {
	// LogCreate logs a customer creation
	LogCreate(ctx context.Context, customer *Customer, createdBy string) error

	// LogUpdate logs a customer update
	LogUpdate(ctx context.Context, oldCustomer, newCustomer *Customer, updatedBy string) error

	// LogDelete logs a customer deletion
	LogDelete(ctx context.Context, customer *Customer, deletedBy string) error
}

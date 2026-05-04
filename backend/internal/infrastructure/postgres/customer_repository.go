package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/datalake-free/internal/domain/customer"
)

// CustomerRepository implements customer.Repository using PostgreSQL
type CustomerRepository struct {
	db *pgxpool.Pool
}

func NewCustomerRepository(db *pgxpool.Pool) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// Create saves a new customer
func (r *CustomerRepository) Create(ctx context.Context, c *customer.Customer) error {
	query := `
		INSERT INTO customers (name, email, phone, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := r.db.QueryRow(ctx, query,
		c.Name, c.Email, c.Phone, c.CreatedBy, c.CreatedAt, c.UpdatedAt,
	).Scan(&c.ID)

	if err != nil {
		return fmt.Errorf("insert customer: %w", err)
	}

	return nil
}

// Update updates an existing customer
func (r *CustomerRepository) Update(ctx context.Context, c *customer.Customer) error {
	query := `
		UPDATE customers
		SET name = $1, email = $2, phone = $3, updated_at = $4
		WHERE id = $5
	`

	result, err := r.db.Exec(ctx, query,
		c.Name, c.Email, c.Phone, c.UpdatedAt, c.ID,
	)

	if err != nil {
		return fmt.Errorf("update customer: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("customer not found: %d", c.ID)
	}

	return nil
}

// Delete removes a customer
func (r *CustomerRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM customers WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete customer: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("customer not found: %d", id)
	}

	return nil
}

// GetByID retrieves a customer by ID
func (r *CustomerRepository) GetByID(ctx context.Context, id int64) (*customer.Customer, error) {
	query := `
		SELECT id, name, email, phone, created_by, created_at, updated_at
		FROM customers
		WHERE id = $1
	`

	var c customer.Customer
	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.Email, &c.Phone, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("get customer: %w", err)
	}

	return &c, nil
}

// GetByEmail retrieves a customer by email
func (r *CustomerRepository) GetByEmail(ctx context.Context, email string) (*customer.Customer, error) {
	query := `
		SELECT id, name, email, phone, created_by, created_at, updated_at
		FROM customers
		WHERE email = $1
	`

	var c customer.Customer
	err := r.db.QueryRow(ctx, query, email).Scan(
		&c.ID, &c.Name, &c.Email, &c.Phone, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("get customer by email: %w", err)
	}

	return &c, nil
}

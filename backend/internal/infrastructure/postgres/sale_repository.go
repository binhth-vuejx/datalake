package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/datalake-free/internal/domain/sale"
)

// SalesRepository implements sale.Repository using PostgreSQL
type SalesRepository struct {
	db *pgxpool.Pool
}

func NewSalesRepository(db *pgxpool.Pool) *SalesRepository {
	return &SalesRepository{db: db}
}

// Create saves a new sale
func (r *SalesRepository) Create(ctx context.Context, s *sale.Sale) error {
	query := `
		INSERT INTO sales (customer_id, amount, status, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := r.db.QueryRow(ctx, query,
		s.CustomerID, s.Amount, s.Status, s.CreatedBy, s.CreatedAt, s.UpdatedAt,
	).Scan(&s.ID)

	if err != nil {
		return fmt.Errorf("insert sale: %w", err)
	}

	return nil
}

// Update updates an existing sale
func (r *SalesRepository) Update(ctx context.Context, s *sale.Sale) error {
	query := `
		UPDATE sales
		SET amount = $1, status = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := r.db.Exec(ctx, query,
		s.Amount, s.Status, s.UpdatedAt, s.ID,
	)

	if err != nil {
		return fmt.Errorf("update sale: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("sale not found: %d", s.ID)
	}

	return nil
}

// Delete removes a sale
func (r *SalesRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM sales WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete sale: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("sale not found: %d", id)
	}

	return nil
}

// GetByID retrieves a sale by ID
func (r *SalesRepository) GetByID(ctx context.Context, id int64) (*sale.Sale, error) {
	query := `
		SELECT id, customer_id, amount, status, created_by, created_at, updated_at
		FROM sales
		WHERE id = $1
	`

	var s sale.Sale
	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.CustomerID, &s.Amount, &s.Status, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("get sale: %w", err)
	}

	return &s, nil
}

// GetByCustomerID retrieves all sales for a customer
func (r *SalesRepository) GetByCustomerID(ctx context.Context, customerID int64) ([]*sale.Sale, error) {
	query := `
		SELECT id, customer_id, amount, status, created_by, created_at, updated_at
		FROM sales
		WHERE customer_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("query sales by customer: %w", err)
	}
	defer rows.Close()

	var sales []*sale.Sale
	for rows.Next() {
		var s sale.Sale
		err := rows.Scan(
			&s.ID, &s.CustomerID, &s.Amount, &s.Status, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan sale: %w", err)
		}
		sales = append(sales, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return sales, nil
}

// List retrieves sales with pagination support
func (r *SalesRepository) List(ctx context.Context, limit, offset int) ([]*sale.Sale, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM sales`
	var total int
	err := r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count sales: %w", err)
	}

	// Get paginated results
	query := `
		SELECT id, customer_id, amount, status, created_by, created_at, updated_at
		FROM sales
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query sales: %w", err)
	}
	defer rows.Close()

	var sales []*sale.Sale
	for rows.Next() {
		var s sale.Sale
		err := rows.Scan(
			&s.ID, &s.CustomerID, &s.Amount, &s.Status, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan sale: %w", err)
		}
		sales = append(sales, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return sales, total, nil
}

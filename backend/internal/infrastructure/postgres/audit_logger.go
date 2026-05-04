package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/datalake-free/internal/domain/customer"
	"github.com/yourorg/datalake-free/internal/domain/order"
	"github.com/yourorg/datalake-free/internal/domain/sale"
)

// AuditLogger implements customer.AuditLogger using PostgreSQL
type AuditLogger struct {
	db *pgxpool.Pool
}

func NewAuditLogger(db *pgxpool.Pool) *AuditLogger {
	return &AuditLogger{db: db}
}

// LogCreate logs a customer creation
func (l *AuditLogger) LogCreate(ctx context.Context, c *customer.Customer, createdBy string) error {
	newData, err := json.Marshal(map[string]interface{}{
		"id":         c.ID,
		"name":       c.Name,
		"email":      c.Email,
		"phone":      c.Phone,
		"created_by": c.CreatedBy,
		"created_at": c.CreatedAt,
	})
	if err != nil {
		return fmt.Errorf("marshal new data: %w", err)
	}

	query := `
		INSERT INTO audit_log (table_name, record_id, action, new_data, changed_by, changed_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = l.db.Exec(ctx, query,
		"customers", c.ID, "CREATE", newData, createdBy, time.Now(),
	)

	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}

	return nil
}

// LogUpdate logs a customer update
func (l *AuditLogger) LogUpdate(ctx context.Context, oldCustomer, newCustomer *customer.Customer, updatedBy string) error {
	oldData, err := json.Marshal(map[string]interface{}{
		"name":  oldCustomer.Name,
		"email": oldCustomer.Email,
		"phone": oldCustomer.Phone,
	})
	if err != nil {
		return fmt.Errorf("marshal old data: %w", err)
	}

	newData, err := json.Marshal(map[string]interface{}{
		"name":  newCustomer.Name,
		"email": newCustomer.Email,
		"phone": newCustomer.Phone,
	})
	if err != nil {
		return fmt.Errorf("marshal new data: %w", err)
	}

	query := `
		INSERT INTO audit_log (table_name, record_id, action, old_data, new_data, changed_by, changed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = l.db.Exec(ctx, query,
		"customers", newCustomer.ID, "UPDATE", oldData, newData, updatedBy, time.Now(),
	)

	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}

	return nil
}

// LogDelete logs a customer deletion
func (l *AuditLogger) LogDelete(ctx context.Context, c *customer.Customer, deletedBy string) error {
	oldData, err := json.Marshal(map[string]interface{}{
		"id":         c.ID,
		"name":       c.Name,
		"email":      c.Email,
		"phone":      c.Phone,
		"created_by": c.CreatedBy,
		"created_at": c.CreatedAt,
	})
	if err != nil {
		return fmt.Errorf("marshal old data: %w", err)
	}

	query := `
		INSERT INTO audit_log (table_name, record_id, action, old_data, changed_by, changed_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = l.db.Exec(ctx, query,
		"customers", c.ID, "DELETE", oldData, deletedBy, time.Now(),
	)

	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}

	return nil
}


// LogCreate logs an order creation
func (l *AuditLogger) LogCreateOrder(ctx context.Context, o *order.Order, createdBy string) error {
	newData, err := json.Marshal(map[string]interface{}{
		"id":            o.ID,
		"customer_id":   o.CustomerID,
		"total_amount":  o.TotalAmount,
		"status":        o.Status,
		"created_by":    o.CreatedBy,
		"created_at":    o.CreatedAt,
	})
	if err != nil {
		return fmt.Errorf("marshal new data: %w", err)
	}

	query := `
		INSERT INTO audit_log (table_name, record_id, action, new_data, changed_by, changed_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = l.db.Exec(ctx, query,
		"orders", o.ID, "CREATE", newData, createdBy, time.Now(),
	)

	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}

	return nil
}

// LogUpdate logs an order update
func (l *AuditLogger) LogUpdateOrder(ctx context.Context, oldOrder, newOrder *order.Order, updatedBy string) error {
	oldData, err := json.Marshal(map[string]interface{}{
		"total_amount": oldOrder.TotalAmount,
		"status":       oldOrder.Status,
	})
	if err != nil {
		return fmt.Errorf("marshal old data: %w", err)
	}

	newData, err := json.Marshal(map[string]interface{}{
		"total_amount": newOrder.TotalAmount,
		"status":       newOrder.Status,
	})
	if err != nil {
		return fmt.Errorf("marshal new data: %w", err)
	}

	query := `
		INSERT INTO audit_log (table_name, record_id, action, old_data, new_data, changed_by, changed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = l.db.Exec(ctx, query,
		"orders", newOrder.ID, "UPDATE", oldData, newData, updatedBy, time.Now(),
	)

	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}

	return nil
}

// LogDelete logs an order deletion
func (l *AuditLogger) LogDeleteOrder(ctx context.Context, o *order.Order, deletedBy string) error {
	oldData, err := json.Marshal(map[string]interface{}{
		"id":           o.ID,
		"customer_id":  o.CustomerID,
		"total_amount": o.TotalAmount,
		"status":       o.Status,
		"created_by":   o.CreatedBy,
		"created_at":   o.CreatedAt,
	})
	if err != nil {
		return fmt.Errorf("marshal old data: %w", err)
	}

	query := `
		INSERT INTO audit_log (table_name, record_id, action, old_data, changed_by, changed_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = l.db.Exec(ctx, query,
		"orders", o.ID, "DELETE", oldData, deletedBy, time.Now(),
	)

	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}

	return nil
}


// LogCreateSale logs a sale creation
func (l *AuditLogger) LogCreateSale(ctx context.Context, s *sale.Sale, createdBy string) error {
	newData, err := json.Marshal(map[string]interface{}{
		"id":         s.ID,
		"customer_id": s.CustomerID,
		"amount":     s.Amount,
		"status":     s.Status,
		"created_by": s.CreatedBy,
		"created_at": s.CreatedAt,
	})
	if err != nil {
		return fmt.Errorf("marshal new data: %w", err)
	}

	query := `
		INSERT INTO audit_log (table_name, record_id, action, new_data, changed_by, changed_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = l.db.Exec(ctx, query,
		"sales", s.ID, "CREATE", newData, createdBy, time.Now(),
	)

	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}

	return nil
}

// LogUpdateSale logs a sale update
func (l *AuditLogger) LogUpdateSale(ctx context.Context, oldSale, newSale *sale.Sale, updatedBy string) error {
	oldData, err := json.Marshal(map[string]interface{}{
		"amount": oldSale.Amount,
		"status": oldSale.Status,
	})
	if err != nil {
		return fmt.Errorf("marshal old data: %w", err)
	}

	newData, err := json.Marshal(map[string]interface{}{
		"amount": newSale.Amount,
		"status": newSale.Status,
	})
	if err != nil {
		return fmt.Errorf("marshal new data: %w", err)
	}

	query := `
		INSERT INTO audit_log (table_name, record_id, action, old_data, new_data, changed_by, changed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = l.db.Exec(ctx, query,
		"sales", newSale.ID, "UPDATE", oldData, newData, updatedBy, time.Now(),
	)

	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}

	return nil
}

// LogDeleteSale logs a sale deletion
func (l *AuditLogger) LogDeleteSale(ctx context.Context, s *sale.Sale, deletedBy string) error {
	oldData, err := json.Marshal(map[string]interface{}{
		"id":         s.ID,
		"customer_id": s.CustomerID,
		"amount":     s.Amount,
		"status":     s.Status,
		"created_by": s.CreatedBy,
		"created_at": s.CreatedAt,
	})
	if err != nil {
		return fmt.Errorf("marshal old data: %w", err)
	}

	query := `
		INSERT INTO audit_log (table_name, record_id, action, old_data, changed_by, changed_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = l.db.Exec(ctx, query,
		"sales", s.ID, "DELETE", oldData, deletedBy, time.Now(),
	)

	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}

	return nil
}

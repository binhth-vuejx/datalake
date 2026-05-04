package dataset

import "context"

// Repository is the port (interface) — infrastructure implements this.
type Repository interface {
	// Config operations
	LoadConfig(ctx context.Context) ([]Dataset, error)

	// Query execution
	ExecuteQuery(ctx context.Context, dataset Dataset, query Query, params map[string]interface{}) (*QueryResult, error)
}

package rilldata

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/yourorg/datalake-free/internal/domain/dataset"
)

// Repository implements dataset.Repository using YAML + PostgreSQL
type Repository struct {
	configPath string
	db         *pgxpool.Pool
	datasets   []dataset.Dataset // Cached from YAML
}

func NewRepository(configPath string, db *pgxpool.Pool) (*Repository, error) {
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	return &Repository{
		configPath: configPath,
		db:         db,
		datasets:   cfg.ToDomain(),
	}, nil
}

func (r *Repository) LoadConfig(ctx context.Context) ([]dataset.Dataset, error) {
	return r.datasets, nil
}

func (r *Repository) ExecuteQuery(
	ctx context.Context,
	ds dataset.Dataset,
	query dataset.Query,
	params map[string]interface{},
) (*dataset.QueryResult, error) {
	// Validate params
	if err := query.ValidateParams(params); err != nil {
		return nil, err
	}

	// Build SQL args
	args := make([]interface{}, len(query.Params))
	for i, p := range query.Params {
		val, ok := params[p.Name]
		if !ok {
			val = p.Default
		}
		converted, err := convertParam(val, p.Type)
		if err != nil {
			return nil, fmt.Errorf("param %s: %w", p.Name, err)
		}
		args[i] = converted
	}

	// Execute SQL
	rows, err := r.db.Query(ctx, query.SQL, args...)
	if err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}
	defer rows.Close()

	// Scan results
	cols := rows.FieldDescriptions()
	var results []map[string]interface{}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range cols {
			row[string(col.Name)] = values[i]
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scan rows: %w", err)
	}

	return &dataset.QueryResult{
		Data:       results,
		Count:      len(results),
		ExecutedAt: time.Now(),
	}, nil
}

// convertParam converts a parameter value to the expected type
func convertParam(val interface{}, targetType dataset.ParamType) (interface{}, error) {
	if val == nil {
		return nil, nil
	}

	switch targetType {
	case dataset.ParamTypeInt64:
		switch v := val.(type) {
		case int:
			return int64(v), nil
		case int64:
			return v, nil
		case int32:
			return int64(v), nil
		case int16:
			return int64(v), nil
		case int8:
			return int64(v), nil
		case float64:
			return int64(v), nil
		case float32:
			return int64(v), nil
		case string:
			return strconv.ParseInt(v, 10, 64)
		default:
			return nil, fmt.Errorf("cannot convert %T to int64", val)
		}

	case dataset.ParamTypeFloat64:
		switch v := val.(type) {
		case float64:
			return v, nil
		case float32:
			return float64(v), nil
		case int:
			return float64(v), nil
		case int64:
			return float64(v), nil
		case int32:
			return float64(v), nil
		case int16:
			return float64(v), nil
		case int8:
			return float64(v), nil
		case string:
			return strconv.ParseFloat(v, 64)
		default:
			return nil, fmt.Errorf("cannot convert %T to float64", val)
		}

	case dataset.ParamTypeString:
		if s, ok := val.(string); ok {
			return s, nil
		}
		return fmt.Sprintf("%v", val), nil

	case dataset.ParamTypeTimestamp:
		switch v := val.(type) {
		case time.Time:
			return v, nil
		case string:
			// Try parsing common timestamp formats
			formats := []string{
				time.RFC3339,
				time.RFC3339Nano,
				"2006-01-02T15:04:05Z",
				"2006-01-02 15:04:05",
				"2006-01-02",
			}
			for _, format := range formats {
				if t, err := time.Parse(format, v); err == nil {
					return t, nil
				}
			}
			return nil, fmt.Errorf("cannot parse timestamp: %s", v)
		default:
			return nil, fmt.Errorf("cannot convert %T to timestamp", val)
		}

	case dataset.ParamTypeBool:
		switch v := val.(type) {
		case bool:
			return v, nil
		case string:
			return strconv.ParseBool(v)
		case int, int64, int32, int16, int8:
			return v != 0, nil
		default:
			return nil, fmt.Errorf("cannot convert %T to bool", val)
		}

	default:
		return nil, fmt.Errorf("unsupported param type: %s", targetType)
	}
}

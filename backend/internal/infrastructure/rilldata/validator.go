package rilldata

import (
	"fmt"
	"strings"
)

// validateConfig validates the merged configuration
func validateConfig(cfg *Config) error {
	if cfg.Version == "" {
		return fmt.Errorf("version is required")
	}

	if len(cfg.Datasets) == 0 {
		return fmt.Errorf("at least one dataset is required")
	}

	// Check for duplicate dataset names
	names := make(map[string]bool)
	for _, ds := range cfg.Datasets {
		if names[ds.Name] {
			return fmt.Errorf("duplicate dataset name: %s", ds.Name)
		}
		names[ds.Name] = true

		if err := validateDataset(ds); err != nil {
			return fmt.Errorf("dataset %s: %w", ds.Name, err)
		}
	}

	return nil
}

// validateDataset validates a single dataset configuration
func validateDataset(ds DatasetConfig) error {
	if ds.Name == "" {
		return fmt.Errorf("dataset name is required")
	}

	if ds.Table == "" {
		return fmt.Errorf("table is required")
	}

	if len(ds.Queries) == 0 {
		return fmt.Errorf("at least one query is required")
	}

	// Check for duplicate query names
	queryNames := make(map[string]bool)
	for _, q := range ds.Queries {
		if queryNames[q.Name] {
			return fmt.Errorf("duplicate query name: %s", q.Name)
		}
		queryNames[q.Name] = true

		if err := validateQuery(q); err != nil {
			return fmt.Errorf("query %s: %w", q.Name, err)
		}
	}

	return nil
}

// validateQuery validates a single query configuration
func validateQuery(q QueryDef) error {
	if q.Name == "" {
		return fmt.Errorf("query name is required")
	}

	if q.Method == "" {
		return fmt.Errorf("method is required")
	}

	if q.Method != "GET" && q.Method != "POST" {
		return fmt.Errorf("method must be GET or POST, got: %s", q.Method)
	}

	if q.SQL == "" {
		return fmt.Errorf("sql is required")
	}

	// Validate params
	for i, p := range q.Params {
		if err := validateParam(p); err != nil {
			return fmt.Errorf("param %d: %w", i, err)
		}
	}

	// Validate SQL placeholders match params count
	placeholderCount := countSQLPlaceholders(q.SQL)
	if placeholderCount != len(q.Params) {
		return fmt.Errorf("SQL has %d placeholders but %d params defined", placeholderCount, len(q.Params))
	}

	return nil
}

// validateParam validates a single parameter configuration
func validateParam(p ParamDef) error {
	if p.Name == "" {
		return fmt.Errorf("param name is required")
	}

	if p.Type == "" {
		return fmt.Errorf("param type is required")
	}

	validTypes := []string{"int64", "float64", "string", "timestamp", "bool"}
	isValid := false
	for _, t := range validTypes {
		if p.Type == t {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("invalid param type: %s (must be one of: %s)", p.Type, strings.Join(validTypes, ", "))
	}

	return nil
}

// countSQLPlaceholders counts the number of $1, $2, etc. placeholders in SQL
func countSQLPlaceholders(sql string) int {
	count := 0
	for i := 1; ; i++ {
		placeholder := fmt.Sprintf("$%d", i)
		if strings.Contains(sql, placeholder) {
			count++
		} else {
			break
		}
	}
	return count
}

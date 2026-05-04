package dataset

import "time"

// ─── Value Objects ────────────────────────────────────────────────────────────

type ParamType string

const (
	ParamTypeInt64     ParamType = "int64"
	ParamTypeFloat64   ParamType = "float64"
	ParamTypeString    ParamType = "string"
	ParamTypeTimestamp ParamType = "timestamp"
	ParamTypeBool      ParamType = "bool"
)

type QueryMethod string

const (
	QueryMethodGET  QueryMethod = "GET"
	QueryMethodPOST QueryMethod = "POST"
)

// ─── Entities ─────────────────────────────────────────────────────────────────

type Dataset struct {
	Name        string
	Description string
	Table       string
	Columns     []Column
	Queries     []Query
}

type Column struct {
	Name        string
	Type        ParamType
	Description string
}

type Query struct {
	Name        string
	Description string
	Method      QueryMethod
	SQL         string
	Params      []Param
}

type Param struct {
	Name     string
	Type     ParamType
	Default  interface{}
	Required bool
}

// ─── Query Result ─────────────────────────────────────────────────────────────

type QueryResult struct {
	Data       []map[string]interface{}
	Count      int
	Error      string
	ExecutedAt time.Time
}

// ─── Business Rules ───────────────────────────────────────────────────────────

func (q *Query) ValidateParams(params map[string]interface{}) error {
	for _, p := range q.Params {
		val, exists := params[p.Name]
		if !exists {
			if p.Required {
				return ErrMissingRequiredParam{ParamName: p.Name}
			}
			continue
		}
		if err := validateParamType(val, p.Type); err != nil {
			return ErrInvalidParamType{ParamName: p.Name, Expected: p.Type, Got: val}
		}
	}
	return nil
}

func (d *Dataset) FindQuery(name string) (*Query, error) {
	for i := range d.Queries {
		if d.Queries[i].Name == name {
			return &d.Queries[i], nil
		}
	}
	return nil, ErrQueryNotFound{QueryName: name}
}

// validateParamType checks if the value matches the expected type
func validateParamType(val interface{}, expectedType ParamType) error {
	switch expectedType {
	case ParamTypeInt64:
		switch val.(type) {
		case int, int64, int32, int16, int8:
			return nil
		case float64:
			// JSON numbers are float64, check if it's an integer
			if f, ok := val.(float64); ok && f == float64(int64(f)) {
				return nil
			}
		}
	case ParamTypeFloat64:
		switch val.(type) {
		case float64, float32, int, int64, int32, int16, int8:
			return nil
		}
	case ParamTypeString:
		if _, ok := val.(string); ok {
			return nil
		}
	case ParamTypeTimestamp:
		if _, ok := val.(string); ok {
			return nil
		}
		if _, ok := val.(time.Time); ok {
			return nil
		}
	case ParamTypeBool:
		if _, ok := val.(bool); ok {
			return nil
		}
	}
	return ErrInvalidParamType{Expected: expectedType, Got: val}
}

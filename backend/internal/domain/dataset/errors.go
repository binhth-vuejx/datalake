package dataset

import "fmt"

type ErrMissingRequiredParam struct {
	ParamName string
}

func (e ErrMissingRequiredParam) Error() string {
	return fmt.Sprintf("missing required param: %s", e.ParamName)
}

type ErrInvalidParamType struct {
	ParamName string
	Expected  ParamType
	Got       interface{}
}

func (e ErrInvalidParamType) Error() string {
	if e.ParamName != "" {
		return fmt.Sprintf("param %s: expected %s, got %T", e.ParamName, e.Expected, e.Got)
	}
	return fmt.Sprintf("expected %s, got %T", e.Expected, e.Got)
}

type ErrQueryNotFound struct {
	QueryName string
}

func (e ErrQueryNotFound) Error() string {
	return fmt.Sprintf("query not found: %s", e.QueryName)
}

type ErrDatasetNotFound struct {
	DatasetName string
}

func (e ErrDatasetNotFound) Error() string {
	return fmt.Sprintf("dataset not found: %s", e.DatasetName)
}

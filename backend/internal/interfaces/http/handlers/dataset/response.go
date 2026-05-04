package dataset

import (
	datasetqueries "github.com/yourorg/datalake-free/internal/application/dataset/queries"
)

// DatasetResponse is the API representation of a dataset.
type DatasetResponse struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Queries     []QueryResponse `json:"queries"`
}

type QueryResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Method      string `json:"method" enum:"GET,POST"`
}

type ListDatasetsOutput struct {
	Body struct {
		Datasets []DatasetResponse `json:"datasets"`
	}
}

type ExecuteQueryOutput struct {
	Body struct {
		Data  []map[string]interface{} `json:"data"`
		Count int                      `json:"count"`
		Error string                   `json:"error,omitempty"`
	}
}

// toDatasetResponse maps application result → API response.
func toDatasetResponse(info datasetqueries.DatasetInfo) DatasetResponse {
	queries := make([]QueryResponse, len(info.Queries))
	for i, q := range info.Queries {
		queries[i] = QueryResponse{
			Name:        q.Name,
			Description: q.Description,
			Method:      q.Method,
		}
	}
	return DatasetResponse{
		Name:        info.Name,
		Description: info.Description,
		Queries:     queries,
	}
}

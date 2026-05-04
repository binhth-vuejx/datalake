package queries

import (
	"context"

	"github.com/yourorg/datalake-free/internal/domain/dataset"
)

type ListDatasetsQuery struct{}

type ListDatasetsResult struct {
	Datasets []DatasetInfo
}

type DatasetInfo struct {
	Name        string
	Description string
	Queries     []QueryInfo
}

type QueryInfo struct {
	Name        string
	Description string
	Method      string
}

type ListDatasetsHandler struct {
	repo dataset.Repository
}

func NewListDatasetsHandler(repo dataset.Repository) *ListDatasetsHandler {
	return &ListDatasetsHandler{repo: repo}
}

func (h *ListDatasetsHandler) Handle(ctx context.Context, q ListDatasetsQuery) (*ListDatasetsResult, error) {
	datasets, err := h.repo.LoadConfig(ctx)
	if err != nil {
		return nil, err
	}

	result := &ListDatasetsResult{
		Datasets: make([]DatasetInfo, len(datasets)),
	}

	for i, ds := range datasets {
		queries := make([]QueryInfo, len(ds.Queries))
		for j, q := range ds.Queries {
			queries[j] = QueryInfo{
				Name:        q.Name,
				Description: q.Description,
				Method:      string(q.Method),
			}
		}

		result.Datasets[i] = DatasetInfo{
			Name:        ds.Name,
			Description: ds.Description,
			Queries:     queries,
		}
	}

	return result, nil
}

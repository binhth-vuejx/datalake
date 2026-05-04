package queries

import (
	"context"

	"github.com/yourorg/datalake-free/internal/domain/dataset"
)

type ExecuteQueryCommand struct {
	DatasetName string
	QueryName   string
	Params      map[string]interface{}
}

type ExecuteQueryResult struct {
	Data  []map[string]interface{}
	Count int
	Error string
}

type ExecuteQueryHandler struct {
	repo dataset.Repository
}

func NewExecuteQueryHandler(repo dataset.Repository) *ExecuteQueryHandler {
	return &ExecuteQueryHandler{repo: repo}
}

func (h *ExecuteQueryHandler) Handle(ctx context.Context, cmd ExecuteQueryCommand) (*ExecuteQueryResult, error) {
	// Load datasets
	datasets, err := h.repo.LoadConfig(ctx)
	if err != nil {
		return nil, err
	}

	// Find dataset
	var targetDataset *dataset.Dataset
	for i := range datasets {
		if datasets[i].Name == cmd.DatasetName {
			targetDataset = &datasets[i]
			break
		}
	}
	if targetDataset == nil {
		return nil, dataset.ErrDatasetNotFound{DatasetName: cmd.DatasetName}
	}

	// Find query
	query, err := targetDataset.FindQuery(cmd.QueryName)
	if err != nil {
		return nil, err
	}

	// Execute
	result, err := h.repo.ExecuteQuery(ctx, *targetDataset, *query, cmd.Params)
	if err != nil {
		return &ExecuteQueryResult{Error: err.Error()}, nil
	}

	return &ExecuteQueryResult{
		Data:  result.Data,
		Count: result.Count,
	}, nil
}

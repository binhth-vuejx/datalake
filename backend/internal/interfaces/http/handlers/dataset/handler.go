package dataset

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	datasetqueries "github.com/yourorg/datalake-free/internal/application/dataset/queries"
)

// Handler wires application use-cases to huma HTTP handlers.
type Handler struct {
	executeQueryHandler *datasetqueries.ExecuteQueryHandler
	listDatasetsHandler *datasetqueries.ListDatasetsHandler
}

func NewHandler(
	executeQuery *datasetqueries.ExecuteQueryHandler,
	listDatasets *datasetqueries.ListDatasetsHandler,
) *Handler {
	return &Handler{
		executeQueryHandler: executeQuery,
		listDatasetsHandler: listDatasets,
	}
}

// Register mounts all dataset routes onto the huma API.
func Register(api huma.API, h *Handler) {
	// List all datasets
	huma.Register(api, huma.Operation{
		OperationID: "list-datasets",
		Method:      http.MethodGet,
		Path:        "/api/datasets",
		Summary:     "List all available datasets",
		Tags:        []string{"datasets"},
	}, h.ListDatasets)

	// Execute query (GET)
	huma.Register(api, huma.Operation{
		OperationID: "execute-query-get",
		Method:      http.MethodGet,
		Path:        "/api/datasets/{dataset}/query/{query}",
		Summary:     "Execute a dataset query with query parameters",
		Tags:        []string{"datasets"},
	}, h.ExecuteQueryGET)

	// Execute query (POST)
	huma.Register(api, huma.Operation{
		OperationID: "execute-query-post",
		Method:      http.MethodPost,
		Path:        "/api/datasets/{dataset}/query/{query}",
		Summary:     "Execute a dataset query with body parameters",
		Tags:        []string{"datasets"},
	}, h.ExecuteQueryPOST)
}

func (h *Handler) ListDatasets(ctx context.Context, input *ListDatasetsInput) (*ListDatasetsOutput, error) {
	result, err := h.listDatasetsHandler.Handle(ctx, datasetqueries.ListDatasetsQuery{})
	if err != nil {
		return nil, err
	}

	out := &ListDatasetsOutput{}
	for _, ds := range result.Datasets {
		out.Body.Datasets = append(out.Body.Datasets, toDatasetResponse(ds))
	}
	return out, nil
}

func (h *Handler) ExecuteQueryGET(ctx context.Context, input *ExecuteQueryGETInput) (*ExecuteQueryOutput, error) {
	// Convert query params to map[string]interface{}
	params := make(map[string]interface{})
	for k, v := range input.Params {
		params[k] = v
	}

	result, err := h.executeQueryHandler.Handle(ctx, datasetqueries.ExecuteQueryCommand{
		DatasetName: input.Dataset,
		QueryName:   input.Query,
		Params:      params,
	})
	if err != nil {
		return nil, err
	}

	out := &ExecuteQueryOutput{}
	out.Body.Data = result.Data
	out.Body.Count = result.Count
	out.Body.Error = result.Error
	return out, nil
}

func (h *Handler) ExecuteQueryPOST(ctx context.Context, input *ExecuteQueryPOSTInput) (*ExecuteQueryOutput, error) {
	result, err := h.executeQueryHandler.Handle(ctx, datasetqueries.ExecuteQueryCommand{
		DatasetName: input.Dataset,
		QueryName:   input.Query,
		Params:      input.Body,
	})
	if err != nil {
		return nil, err
	}

	out := &ExecuteQueryOutput{}
	out.Body.Data = result.Data
	out.Body.Count = result.Count
	out.Body.Error = result.Error
	return out, nil
}

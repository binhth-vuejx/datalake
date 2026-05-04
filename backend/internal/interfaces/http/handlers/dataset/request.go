package dataset

// ─── Input structs for huma handlers ─────────────────────────────────────────
// These are the API contract — separate from domain entities.

type ListDatasetsInput struct{}

type ExecuteQueryGETInput struct {
	Dataset string            `path:"dataset" doc:"Dataset name"`
	Query   string            `path:"query"   doc:"Query name"`
	Params  map[string]string `query:"*"      doc:"Query parameters"`
}

type ExecuteQueryPOSTInput struct {
	Dataset string                 `path:"dataset" doc:"Dataset name"`
	Query   string                 `path:"query"   doc:"Query name"`
	Body    map[string]interface{} `json:"-"`
}

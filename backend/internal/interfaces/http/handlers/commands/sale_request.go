package commands

// CreateSaleInput is the request for creating a sale
type CreateSaleInput struct {
	Body struct {
		CustomerID int64   `json:"customer_id" doc:"Customer ID" required:"true"`
		Amount     float64 `json:"amount" doc:"Sale amount in USD" required:"true"`
	}
}

// UpdateSaleInput is the request for updating a sale
type UpdateSaleInput struct {
	ID   int64 `path:"id" doc:"Sale ID" required:"true"`
	Body struct {
		Amount float64 `json:"amount" doc:"Sale amount in USD" required:"true"`
		Status string  `json:"status" doc:"Sale status (pending, completed, cancelled)" required:"true"`
	}
}

// DeleteSaleInput is the request for deleting a sale
type DeleteSaleInput struct {
	ID int64 `path:"id" doc:"Sale ID" required:"true"`
}

// ListSalesInput is the request for listing sales with pagination
type ListSalesInput struct {
	Limit  int `query:"limit" doc:"Maximum number of sales to return" default:"50" minimum:"1" maximum:"1000"`
	Offset int `query:"offset" doc:"Number of sales to skip" default:"0" minimum:"0"`
}

// GetSaleByCustomerInput is the request for getting sales by customer ID
type GetSaleByCustomerInput struct {
	CustomerID int64 `path:"customer_id" doc:"Customer ID" required:"true"`
	Limit      int   `query:"limit" doc:"Maximum number of sales to return" default:"50" minimum:"1" maximum:"1000"`
	Offset     int   `query:"offset" doc:"Number of sales to skip" default:"0" minimum:"0"`
}

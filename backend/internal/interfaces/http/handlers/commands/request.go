package commands

// ─── Customer Commands ────────────────────────────────────────────────────

// CreateCustomerInput is the request body for creating a customer
type CreateCustomerInput struct {
	Body struct {
		Name  string `json:"name" doc:"Customer name" required:"true"`
		Email string `json:"email" doc:"Customer email" required:"true"`
		Phone string `json:"phone" doc:"Customer phone number"`
	}
}

// UpdateCustomerInput is the request body for updating a customer
type UpdateCustomerInput struct {
	ID   int64 `path:"id" doc:"Customer ID" required:"true"`
	Body struct {
		Name  string `json:"name" doc:"Customer name"`
		Email string `json:"email" doc:"Customer email"`
		Phone string `json:"phone" doc:"Customer phone number"`
	}
}

// DeleteCustomerInput is the request body for deleting a customer
type DeleteCustomerInput struct {
	ID int64 `path:"id" doc:"Customer ID" required:"true"`
}

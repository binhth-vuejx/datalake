package commands

// CreateOrderInput is the request for creating an order
type CreateOrderInput struct {
	Body struct {
		CustomerID  int64   `json:"customer_id" doc:"Customer ID" required:"true"`
		TotalAmount float64 `json:"total_amount" doc:"Order total amount" required:"true"`
	}
}

// UpdateOrderInput is the request for updating an order
type UpdateOrderInput struct {
	ID   int64 `path:"id" doc:"Order ID" required:"true"`
	Body struct {
		TotalAmount float64 `json:"total_amount" doc:"Order total amount" required:"true"`
		Status      string  `json:"status" doc:"Order status (pending, completed, cancelled)" required:"true"`
	}
}

// DeleteOrderInput is the request for deleting an order
type DeleteOrderInput struct {
	ID int64 `path:"id" doc:"Order ID" required:"true"`
}

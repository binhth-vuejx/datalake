package commands

import "time"

// OrderData represents an order in the system
type OrderData struct {
	ID          int64     `json:"id" doc:"Order ID"`
	CustomerID  int64     `json:"customer_id" doc:"Customer ID"`
	TotalAmount float64   `json:"total_amount" doc:"Order total amount"`
	Status      string    `json:"status" doc:"Order status"`
	CreatedBy   string    `json:"created_by" doc:"User who created the order"`
	CreatedAt   time.Time `json:"created_at" doc:"Creation timestamp"`
	UpdatedAt   time.Time `json:"updated_at" doc:"Last update timestamp"`
}

// CreateOrderOutput is the response for creating an order
type CreateOrderOutput struct {
	Status int `doc:"HTTP status code"`
	Body struct {
		Success bool       `json:"success" doc:"Whether the operation was successful"`
		Data    *OrderData `json:"data,omitempty" doc:"Created order data"`
		Error   string     `json:"error,omitempty" doc:"Error message if operation failed"`
	}
}

// UpdateOrderOutput is the response for updating an order
type UpdateOrderOutput struct {
	Status int `doc:"HTTP status code"`
	Body struct {
		Success bool       `json:"success" doc:"Whether the operation was successful"`
		Data    *OrderData `json:"data,omitempty" doc:"Updated order data"`
		Error   string     `json:"error,omitempty" doc:"Error message if operation failed"`
	}
}

// DeleteOrderOutput is the response for deleting an order
type DeleteOrderOutput struct {
	Status int `doc:"HTTP status code"`
	Body struct {
		Success bool                `json:"success" doc:"Whether the operation was successful"`
		Data    *DeleteMessageData  `json:"data,omitempty" doc:"Deletion result data"`
		Error   string              `json:"error,omitempty" doc:"Error message if operation failed"`
	}
}

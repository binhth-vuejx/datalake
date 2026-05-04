package commands

import "time"

// SaleData represents a sale in the system
type SaleData struct {
	ID         int64     `json:"id" doc:"Sale ID"`
	CustomerID int64     `json:"customer_id" doc:"Customer ID"`
	Amount     float64   `json:"amount" doc:"Sale amount in USD"`
	Status     string    `json:"status" doc:"Sale status"`
	CreatedBy  string    `json:"created_by" doc:"User who created the sale"`
	CreatedAt  time.Time `json:"created_at" doc:"Creation timestamp"`
	UpdatedAt  time.Time `json:"updated_at" doc:"Last update timestamp"`
}

// CreateSaleOutput is the response for creating a sale
type CreateSaleOutput struct {
	Status int `doc:"HTTP status code"`
	Body struct {
		Success bool       `json:"success" doc:"Whether the operation was successful"`
		Data    *SaleData  `json:"data,omitempty" doc:"Created sale data"`
		Error   string     `json:"error,omitempty" doc:"Error message if operation failed"`
	}
}

// UpdateSaleOutput is the response for updating a sale
type UpdateSaleOutput struct {
	Status int `doc:"HTTP status code"`
	Body struct {
		Success bool       `json:"success" doc:"Whether the operation was successful"`
		Data    *SaleData  `json:"data,omitempty" doc:"Updated sale data"`
		Error   string     `json:"error,omitempty" doc:"Error message if operation failed"`
	}
}

// DeleteSaleOutput is the response for deleting a sale
type DeleteSaleOutput struct {
	Status int `doc:"HTTP status code"`
	Body struct {
		Success bool   `json:"success" doc:"Whether the operation was successful"`
		Error   string `json:"error,omitempty" doc:"Error message if operation failed"`
	}
}

// ListSalesOutput is the response for listing sales
type ListSalesOutput struct {
	Status int `doc:"HTTP status code"`
	Body struct {
		Success bool       `json:"success" doc:"Whether the operation was successful"`
		Data    []SaleData `json:"data,omitempty" doc:"List of sales"`
		Total   int        `json:"total" doc:"Total count of sales"`
		Error   string     `json:"error,omitempty" doc:"Error message if operation failed"`
	}
}

// GetSaleByCustomerOutput is the response for getting sales by customer
type GetSaleByCustomerOutput struct {
	Status int `doc:"HTTP status code"`
	Body struct {
		Success bool       `json:"success" doc:"Whether the operation was successful"`
		Data    []SaleData `json:"data,omitempty" doc:"List of sales for the customer"`
		Total   int        `json:"total" doc:"Total count of sales for the customer"`
		Error   string     `json:"error,omitempty" doc:"Error message if operation failed"`
	}
}

// GetSaleByIDOutput is the response for getting a sale by ID
type GetSaleByIDOutput struct {
	Status int `doc:"HTTP status code"`
	Body struct {
		Success bool       `json:"success" doc:"Whether the operation was successful"`
		Data    *SaleData  `json:"data,omitempty" doc:"Sale data"`
		Error   string     `json:"error,omitempty" doc:"Error message if operation failed"`
	}
}

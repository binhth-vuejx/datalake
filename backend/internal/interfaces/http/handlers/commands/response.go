package commands

import "time"

// CustomerData represents a customer in the system
type CustomerData struct {
	ID        int64     `json:"id" doc:"Customer ID"`
	Name      string    `json:"name" doc:"Customer name"`
	Email     string    `json:"email" doc:"Customer email"`
	Phone     string    `json:"phone" doc:"Customer phone number"`
	CreatedBy string    `json:"created_by" doc:"User who created the customer"`
	CreatedAt time.Time `json:"created_at" doc:"Creation timestamp"`
	UpdatedAt time.Time `json:"updated_at" doc:"Last update timestamp"`
}

// DeleteMessageData represents the deletion confirmation message
type DeleteMessageData struct {
	Message string `json:"message" doc:"Deletion confirmation message"`
}

// CreateCustomerOutput is the response for creating a customer
type CreateCustomerOutput struct {
	Status int `doc:"HTTP status code"`
	Body struct {
		Success bool          `json:"success" doc:"Whether the operation was successful"`
		Data    *CustomerData `json:"data,omitempty" doc:"Created customer data"`
		Error   string        `json:"error,omitempty" doc:"Error message if operation failed"`
	}
}

// UpdateCustomerOutput is the response for updating a customer
type UpdateCustomerOutput struct {
	Status int `doc:"HTTP status code"`
	Body struct {
		Success bool          `json:"success" doc:"Whether the operation was successful"`
		Data    *CustomerData `json:"data,omitempty" doc:"Updated customer data"`
		Error   string        `json:"error,omitempty" doc:"Error message if operation failed"`
	}
}

// DeleteCustomerOutput is the response for deleting a customer
type DeleteCustomerOutput struct {
	Status int `doc:"HTTP status code"`
	Body struct {
		Success bool                `json:"success" doc:"Whether the operation was successful"`
		Data    *DeleteMessageData  `json:"data,omitempty" doc:"Deletion result data"`
		Error   string              `json:"error,omitempty" doc:"Error message if operation failed"`
	}
}

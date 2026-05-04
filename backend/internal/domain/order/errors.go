package order

import "errors"

var (
	ErrOrderNotFound           = errors.New("order not found")
	ErrCustomerIDRequired      = errors.New("customer_id is required")
	ErrAmountMustBePositive    = errors.New("amount must be positive")
	ErrAmountCannotBeNegative  = errors.New("amount cannot be negative")
	ErrOrderIDRequired         = errors.New("order_id is required")
	ErrInvalidStatus           = errors.New("invalid order status")
)

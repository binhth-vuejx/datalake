package sale

import (
	"testing"
)

// TestCreateSaleCommandValidation tests validation rules for CreateSaleCommand
func TestCreateSaleCommandValidation(t *testing.T) {
	tests := []struct {
		name    string
		cmd     CreateSaleCommand
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid command",
			cmd: CreateSaleCommand{
				CustomerID: 1,
				Amount:     100.50,
				CreatedBy:  "user1",
			},
			wantErr: false,
		},
		{
			name: "invalid customer_id (zero)",
			cmd: CreateSaleCommand{
				CustomerID: 0,
				Amount:     100.50,
				CreatedBy:  "user1",
			},
			wantErr: true,
			errMsg:  "customer_id must be a positive integer",
		},
		{
			name: "invalid customer_id (negative)",
			cmd: CreateSaleCommand{
				CustomerID: -1,
				Amount:     100.50,
				CreatedBy:  "user1",
			},
			wantErr: true,
			errMsg:  "customer_id must be a positive integer",
		},
		{
			name: "invalid amount (zero)",
			cmd: CreateSaleCommand{
				CustomerID: 1,
				Amount:     0,
				CreatedBy:  "user1",
			},
			wantErr: true,
			errMsg:  "amount must be positive",
		},
		{
			name: "invalid amount (negative)",
			cmd: CreateSaleCommand{
				CustomerID: 1,
				Amount:     -50.00,
				CreatedBy:  "user1",
			},
			wantErr: true,
			errMsg:  "amount must be positive",
		},
		{
			name: "empty created_by",
			cmd: CreateSaleCommand{
				CustomerID: 1,
				Amount:     100.50,
				CreatedBy:  "",
			},
			wantErr: true,
			errMsg:  "created_by is required",
		},
		{
			name: "created_by exceeds max length",
			cmd: CreateSaleCommand{
				CustomerID: 1,
				Amount:     100.50,
				CreatedBy:  string(make([]byte, 256)),
			},
			wantErr: true,
			errMsg:  "created_by must not exceed 255 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

// TestUpdateSaleCommandValidation tests validation rules for UpdateSaleCommand
func TestUpdateSaleCommandValidation(t *testing.T) {
	tests := []struct {
		name    string
		cmd     UpdateSaleCommand
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid command",
			cmd: UpdateSaleCommand{
				ID:        1,
				Amount:    150.75,
				Status:    "completed",
				UpdatedBy: "user1",
			},
			wantErr: false,
		},
		{
			name: "invalid id (zero)",
			cmd: UpdateSaleCommand{
				ID:        0,
				Amount:    150.75,
				Status:    "completed",
				UpdatedBy: "user1",
			},
			wantErr: true,
			errMsg:  "id must be a positive integer",
		},
		{
			name: "invalid amount (zero)",
			cmd: UpdateSaleCommand{
				ID:        1,
				Amount:    0,
				Status:    "completed",
				UpdatedBy: "user1",
			},
			wantErr: true,
			errMsg:  "amount must be positive",
		},
		{
			name: "invalid status",
			cmd: UpdateSaleCommand{
				ID:        1,
				Amount:    150.75,
				Status:    "invalid",
				UpdatedBy: "user1",
			},
			wantErr: true,
			errMsg:  "status must be one of: pending, completed, cancelled",
		},
		{
			name: "valid status pending",
			cmd: UpdateSaleCommand{
				ID:        1,
				Amount:    150.75,
				Status:    "pending",
				UpdatedBy: "user1",
			},
			wantErr: false,
		},
		{
			name: "valid status cancelled",
			cmd: UpdateSaleCommand{
				ID:        1,
				Amount:    150.75,
				Status:    "cancelled",
				UpdatedBy: "user1",
			},
			wantErr: false,
		},
		{
			name: "empty updated_by",
			cmd: UpdateSaleCommand{
				ID:        1,
				Amount:    150.75,
				Status:    "completed",
				UpdatedBy: "",
			},
			wantErr: true,
			errMsg:  "updated_by is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

// TestDeleteSaleCommandValidation tests validation rules for DeleteSaleCommand
func TestDeleteSaleCommandValidation(t *testing.T) {
	tests := []struct {
		name    string
		cmd     DeleteSaleCommand
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid command",
			cmd: DeleteSaleCommand{
				ID:        1,
				DeletedBy: "user1",
			},
			wantErr: false,
		},
		{
			name: "invalid id (zero)",
			cmd: DeleteSaleCommand{
				ID:        0,
				DeletedBy: "user1",
			},
			wantErr: true,
			errMsg:  "id must be a positive integer",
		},
		{
			name: "empty deleted_by",
			cmd: DeleteSaleCommand{
				ID:        1,
				DeletedBy: "",
			},
			wantErr: true,
			errMsg:  "deleted_by is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

// Property-based test: Amount Validation
// Property 1: For any CreateSaleCommand with positive amount, validation should succeed
// For any CreateSaleCommand with amount <= 0, validation should fail
func TestAmountValidationProperty(t *testing.T) {
	// Test positive amounts
	positiveAmounts := []float64{0.01, 1.0, 100.50, 999999.99}
	for _, amount := range positiveAmounts {
		cmd := CreateSaleCommand{
			CustomerID: 1,
			Amount:     amount,
			CreatedBy:  "user1",
		}
		if err := cmd.Validate(); err != nil {
			t.Errorf("Positive amount %v should be valid, got error: %v", amount, err)
		}
	}

	// Test non-positive amounts
	nonPositiveAmounts := []float64{0, -0.01, -100.0}
	for _, amount := range nonPositiveAmounts {
		cmd := CreateSaleCommand{
			CustomerID: 1,
			Amount:     amount,
			CreatedBy:  "user1",
		}
		if err := cmd.Validate(); err == nil {
			t.Errorf("Non-positive amount %v should be invalid", amount)
		}
	}
}

// Property-based test: Status Constraint
// Property 2: For any UpdateSaleCommand with valid status, update should succeed
// For any UpdateSaleCommand with invalid status, validation should fail
func TestStatusConstraintProperty(t *testing.T) {
	validStatuses := []string{"pending", "completed", "cancelled"}
	for _, status := range validStatuses {
		cmd := UpdateSaleCommand{
			ID:        1,
			Amount:    100.0,
			Status:    status,
			UpdatedBy: "user1",
		}
		if err := cmd.Validate(); err != nil {
			t.Errorf("Valid status %q should pass validation, got error: %v", status, err)
		}
	}

	invalidStatuses := []string{"invalid", "PENDING", "Completed", "", "processing"}
	for _, status := range invalidStatuses {
		cmd := UpdateSaleCommand{
			ID:        1,
			Amount:    100.0,
			Status:    status,
			UpdatedBy: "user1",
		}
		if err := cmd.Validate(); err == nil {
			t.Errorf("Invalid status %q should fail validation", status)
		}
	}
}

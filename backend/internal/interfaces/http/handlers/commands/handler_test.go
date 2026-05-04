package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// setupTestDB creates a test database connection
func setupTestDB(t *testing.T) *pgxpool.Pool {
	ctx := context.Background()
	dbURL := "postgres://multica:multica@localhost:5432/multica_test?sslmode=disable"
	
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Skipf("Could not connect to test database: %v", err)
	}

	// Ping to verify connection
	if err := pool.Ping(ctx); err != nil {
		t.Skipf("Could not ping test database: %v", err)
	}

	return pool
}

// setupTestRouter creates a test router with huma API
func setupTestRouter(t *testing.T, db *pgxpool.Pool) (http.Handler, *Handler) {
	router := chi.NewMux()
	config := huma.DefaultConfig("Test API", "1.0.0")
	api := humachi.New(router, config)

	handler := NewHandler(db)
	Register(api, handler)

	return router, handler
}

// TestCreateSaleSuccess tests creating a sale with valid input
func TestCreateSaleSuccess(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router, _ := setupTestRouter(t, db)

	// Create a customer first
	customerPayload := map[string]interface{}{
		"name":  "Test Customer",
		"email": "test@example.com",
	}
	customerBody, _ := json.Marshal(customerPayload)
	customerReq := httptest.NewRequest("POST", "/api/customers", bytes.NewReader(customerBody))
	customerReq.Header.Set("Content-Type", "application/json")
	customerW := httptest.NewRecorder()
	router.ServeHTTP(customerW, customerReq)

	if customerW.Code != http.StatusOK {
		t.Skipf("Failed to create customer: %d", customerW.Code)
	}

	var customerResp map[string]interface{}
	json.Unmarshal(customerW.Body.Bytes(), &customerResp)
	customerID := int64(customerResp["data"].(map[string]interface{})["id"].(float64))

	// Now create a sale
	payload := map[string]interface{}{
		"customer_id": customerID,
		"amount":      150.50,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/sales", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !resp["success"].(bool) {
		t.Errorf("Expected success=true, got %v", resp["success"])
	}

	if resp["data"] == nil {
		t.Errorf("Expected data in response, got nil")
	}
}

// TestCreateSaleInvalidAmount tests creating a sale with invalid amount
func TestCreateSaleInvalidAmount(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router, _ := setupTestRouter(t, db)

	// Create a customer first
	customerPayload := map[string]interface{}{
		"name":  "Test Customer",
		"email": "test@example.com",
	}
	customerBody, _ := json.Marshal(customerPayload)
	customerReq := httptest.NewRequest("POST", "/api/customers", bytes.NewReader(customerBody))
	customerReq.Header.Set("Content-Type", "application/json")
	customerW := httptest.NewRecorder()
	router.ServeHTTP(customerW, customerReq)

	var customerResp map[string]interface{}
	json.Unmarshal(customerW.Body.Bytes(), &customerResp)
	customerID := int64(customerResp["data"].(map[string]interface{})["id"].(float64))

	// Try to create a sale with invalid amount (≤ 0)
	payload := map[string]interface{}{
		"customer_id": customerID,
		"amount":      -50.0,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/sales", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["success"].(bool) {
		t.Errorf("Expected success=false for invalid amount, got true")
	}
}

// TestListSalesSuccess tests listing sales with pagination
func TestListSalesSuccess(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router, _ := setupTestRouter(t, db)

	// List sales
	req := httptest.NewRequest("GET", "/api/sales?limit=50&offset=0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !resp["success"].(bool) {
		t.Errorf("Expected success=true, got %v", resp["success"])
	}

	if resp["data"] == nil {
		t.Errorf("Expected data in response, got nil")
	}

	if resp["total"] == nil {
		t.Errorf("Expected total in response, got nil")
	}
}

// TestListSalesPagination tests pagination parameters
func TestListSalesPagination(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router, _ := setupTestRouter(t, db)

	// Test with limit and offset
	req := httptest.NewRequest("GET", "/api/sales?limit=10&offset=5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["success"].(bool) {
		data := resp["data"].([]interface{})
		if len(data) > 10 {
			t.Errorf("Expected at most 10 items, got %d", len(data))
		}
	}
}

// TestGetSaleByIDSuccess tests getting a sale by ID
func TestGetSaleByIDSuccess(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router, _ := setupTestRouter(t, db)

	// Create a customer and sale first
	customerPayload := map[string]interface{}{
		"name":  "Test Customer",
		"email": "test@example.com",
	}
	customerBody, _ := json.Marshal(customerPayload)
	customerReq := httptest.NewRequest("POST", "/api/customers", bytes.NewReader(customerBody))
	customerReq.Header.Set("Content-Type", "application/json")
	customerW := httptest.NewRecorder()
	router.ServeHTTP(customerW, customerReq)

	var customerResp map[string]interface{}
	json.Unmarshal(customerW.Body.Bytes(), &customerResp)
	customerID := int64(customerResp["data"].(map[string]interface{})["id"].(float64))

	salePayload := map[string]interface{}{
		"customer_id": customerID,
		"amount":      150.50,
	}
	saleBody, _ := json.Marshal(salePayload)
	saleReq := httptest.NewRequest("POST", "/api/sales", bytes.NewReader(saleBody))
	saleReq.Header.Set("Content-Type", "application/json")
	saleW := httptest.NewRecorder()
	router.ServeHTTP(saleW, saleReq)

	var saleResp map[string]interface{}
	json.Unmarshal(saleW.Body.Bytes(), &saleResp)
	saleID := int64(saleResp["data"].(map[string]interface{})["id"].(float64))

	// Get the sale by ID
	req := httptest.NewRequest("GET", "/api/sales/"+string(rune(saleID)), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if !resp["success"].(bool) {
		t.Errorf("Expected success=true, got %v", resp["success"])
	}
}

// TestGetSaleByIDNotFound tests getting a non-existent sale
func TestGetSaleByIDNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router, _ := setupTestRouter(t, db)

	// Try to get a non-existent sale
	req := httptest.NewRequest("GET", "/api/sales/99999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

// TestUpdateSaleSuccess tests updating a sale
func TestUpdateSaleSuccess(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router, _ := setupTestRouter(t, db)

	// Create a customer and sale first
	customerPayload := map[string]interface{}{
		"name":  "Test Customer",
		"email": "test@example.com",
	}
	customerBody, _ := json.Marshal(customerPayload)
	customerReq := httptest.NewRequest("POST", "/api/customers", bytes.NewReader(customerBody))
	customerReq.Header.Set("Content-Type", "application/json")
	customerW := httptest.NewRecorder()
	router.ServeHTTP(customerW, customerReq)

	var customerResp map[string]interface{}
	json.Unmarshal(customerW.Body.Bytes(), &customerResp)
	customerID := int64(customerResp["data"].(map[string]interface{})["id"].(float64))

	salePayload := map[string]interface{}{
		"customer_id": customerID,
		"amount":      150.50,
	}
	saleBody, _ := json.Marshal(salePayload)
	saleReq := httptest.NewRequest("POST", "/api/sales", bytes.NewReader(saleBody))
	saleReq.Header.Set("Content-Type", "application/json")
	saleW := httptest.NewRecorder()
	router.ServeHTTP(saleW, saleReq)

	var saleResp map[string]interface{}
	json.Unmarshal(saleW.Body.Bytes(), &saleResp)
	saleID := int64(saleResp["data"].(map[string]interface{})["id"].(float64))

	// Update the sale
	updatePayload := map[string]interface{}{
		"amount": 200.0,
		"status": "completed",
	}
	updateBody, _ := json.Marshal(updatePayload)
	req := httptest.NewRequest("PUT", "/api/sales/"+string(rune(saleID)), bytes.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if !resp["success"].(bool) {
		t.Errorf("Expected success=true, got %v", resp["success"])
	}
}

// TestUpdateSaleInvalidStatus tests updating a sale with invalid status
func TestUpdateSaleInvalidStatus(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router, _ := setupTestRouter(t, db)

	// Create a customer and sale first
	customerPayload := map[string]interface{}{
		"name":  "Test Customer",
		"email": "test@example.com",
	}
	customerBody, _ := json.Marshal(customerPayload)
	customerReq := httptest.NewRequest("POST", "/api/customers", bytes.NewReader(customerBody))
	customerReq.Header.Set("Content-Type", "application/json")
	customerW := httptest.NewRecorder()
	router.ServeHTTP(customerW, customerReq)

	var customerResp map[string]interface{}
	json.Unmarshal(customerW.Body.Bytes(), &customerResp)
	customerID := int64(customerResp["data"].(map[string]interface{})["id"].(float64))

	salePayload := map[string]interface{}{
		"customer_id": customerID,
		"amount":      150.50,
	}
	saleBody, _ := json.Marshal(salePayload)
	saleReq := httptest.NewRequest("POST", "/api/sales", bytes.NewReader(saleBody))
	saleReq.Header.Set("Content-Type", "application/json")
	saleW := httptest.NewRecorder()
	router.ServeHTTP(saleW, saleReq)

	var saleResp map[string]interface{}
	json.Unmarshal(saleW.Body.Bytes(), &saleResp)
	saleID := int64(saleResp["data"].(map[string]interface{})["id"].(float64))

	// Try to update with invalid status
	updatePayload := map[string]interface{}{
		"amount": 200.0,
		"status": "invalid_status",
	}
	updateBody, _ := json.Marshal(updatePayload)
	req := httptest.NewRequest("PUT", "/api/sales/"+string(rune(saleID)), bytes.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["success"].(bool) {
		t.Errorf("Expected success=false for invalid status, got true")
	}
}

// TestDeleteSaleSuccess tests deleting a sale
func TestDeleteSaleSuccess(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router, _ := setupTestRouter(t, db)

	// Create a customer and sale first
	customerPayload := map[string]interface{}{
		"name":  "Test Customer",
		"email": "test@example.com",
	}
	customerBody, _ := json.Marshal(customerPayload)
	customerReq := httptest.NewRequest("POST", "/api/customers", bytes.NewReader(customerBody))
	customerReq.Header.Set("Content-Type", "application/json")
	customerW := httptest.NewRecorder()
	router.ServeHTTP(customerW, customerReq)

	var customerResp map[string]interface{}
	json.Unmarshal(customerW.Body.Bytes(), &customerResp)
	customerID := int64(customerResp["data"].(map[string]interface{})["id"].(float64))

	salePayload := map[string]interface{}{
		"customer_id": customerID,
		"amount":      150.50,
	}
	saleBody, _ := json.Marshal(salePayload)
	saleReq := httptest.NewRequest("POST", "/api/sales", bytes.NewReader(saleBody))
	saleReq.Header.Set("Content-Type", "application/json")
	saleW := httptest.NewRecorder()
	router.ServeHTTP(saleW, saleReq)

	var saleResp map[string]interface{}
	json.Unmarshal(saleW.Body.Bytes(), &saleResp)
	saleID := int64(saleResp["data"].(map[string]interface{})["id"].(float64))

	// Delete the sale
	req := httptest.NewRequest("DELETE", "/api/sales/"+string(rune(saleID)), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if !resp["success"].(bool) {
		t.Errorf("Expected success=true, got %v", resp["success"])
	}
}

// TestDeleteSaleNotFound tests deleting a non-existent sale
func TestDeleteSaleNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router, _ := setupTestRouter(t, db)

	// Try to delete a non-existent sale
	req := httptest.NewRequest("DELETE", "/api/sales/99999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["success"].(bool) {
		t.Errorf("Expected success=false for non-existent sale, got true")
	}
}

// TestGetSaleByCustomerSuccess tests getting sales by customer ID
func TestGetSaleByCustomerSuccess(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router, _ := setupTestRouter(t, db)

	// Create a customer first
	customerPayload := map[string]interface{}{
		"name":  "Test Customer",
		"email": "test@example.com",
	}
	customerBody, _ := json.Marshal(customerPayload)
	customerReq := httptest.NewRequest("POST", "/api/customers", bytes.NewReader(customerBody))
	customerReq.Header.Set("Content-Type", "application/json")
	customerW := httptest.NewRecorder()
	router.ServeHTTP(customerW, customerReq)

	var customerResp map[string]interface{}
	json.Unmarshal(customerW.Body.Bytes(), &customerResp)
	customerID := int64(customerResp["data"].(map[string]interface{})["id"].(float64))

	// Get sales by customer
	req := httptest.NewRequest("GET", "/api/sales/customer/"+string(rune(customerID))+"?limit=50&offset=0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if !resp["success"].(bool) {
		t.Errorf("Expected success=true, got %v", resp["success"])
	}
}

// TestOpenAPISchemaGeneration tests that OpenAPI schema is generated correctly
func TestOpenAPISchemaGeneration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router, _ := setupTestRouter(t, db)

	// Request OpenAPI schema
	req := httptest.NewRequest("GET", "/openapi.json", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for OpenAPI schema, got %d", w.Code)
	}

	var schema map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &schema); err != nil {
		t.Fatalf("Failed to unmarshal OpenAPI schema: %v", err)
	}

	// Verify schema structure
	if schema["openapi"] == nil {
		t.Errorf("Expected openapi version in schema")
	}

	if schema["paths"] == nil {
		t.Errorf("Expected paths in schema")
	}

	paths := schema["paths"].(map[string]interface{})

	// Verify sales endpoints are documented
	expectedPaths := []string{
		"/api/sales",
		"/api/sales/{id}",
		"/api/sales/customer/{customer_id}",
	}

	for _, path := range expectedPaths {
		if paths[path] == nil {
			t.Errorf("Expected path %s in OpenAPI schema", path)
		}
	}
}

// TestResponseFormat tests that responses follow the expected format
func TestResponseFormat(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router, _ := setupTestRouter(t, db)

	// Create a customer and sale
	customerPayload := map[string]interface{}{
		"name":  "Test Customer",
		"email": "test@example.com",
	}
	customerBody, _ := json.Marshal(customerPayload)
	customerReq := httptest.NewRequest("POST", "/api/customers", bytes.NewReader(customerBody))
	customerReq.Header.Set("Content-Type", "application/json")
	customerW := httptest.NewRecorder()
	router.ServeHTTP(customerW, customerReq)

	var customerResp map[string]interface{}
	json.Unmarshal(customerW.Body.Bytes(), &customerResp)
	customerID := int64(customerResp["data"].(map[string]interface{})["id"].(float64))

	salePayload := map[string]interface{}{
		"customer_id": customerID,
		"amount":      150.50,
	}
	saleBody, _ := json.Marshal(salePayload)
	saleReq := httptest.NewRequest("POST", "/api/sales", bytes.NewReader(saleBody))
	saleReq.Header.Set("Content-Type", "application/json")
	saleW := httptest.NewRecorder()
	router.ServeHTTP(saleW, saleReq)

	var resp map[string]interface{}
	json.Unmarshal(saleW.Body.Bytes(), &resp)

	// Verify response format
	if resp["success"] == nil {
		t.Errorf("Expected 'success' field in response")
	}

	if resp["data"] == nil {
		t.Errorf("Expected 'data' field in response")
	}

	data := resp["data"].(map[string]interface{})

	// Verify SaleData fields
	expectedFields := []string{"id", "customer_id", "amount", "status", "created_by", "created_at", "updated_at"}
	for _, field := range expectedFields {
		if data[field] == nil {
			t.Errorf("Expected field '%s' in SaleData", field)
		}
	}
}

// TestStatusCodes tests that correct HTTP status codes are returned
func TestStatusCodes(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router, _ := setupTestRouter(t, db)

	tests := []struct {
		name           string
		method         string
		path           string
		body           map[string]interface{}
		expectedStatus int
	}{
		{
			name:           "GET /api/sales returns 200",
			method:         "GET",
			path:           "/api/sales",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET /api/sales/99999 returns 404",
			method:         "GET",
			path:           "/api/sales/99999",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			if tt.body != nil {
				body, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(tt.method, tt.path, bytes.NewReader(body))
			if tt.body != nil {
				req.Header.Set("Content-Type", "application/json")
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

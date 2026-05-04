package commands

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	customercmds "github.com/yourorg/datalake-free/internal/application/customer/commands"
	ordercmds "github.com/yourorg/datalake-free/internal/application/order/commands"
	salecmds "github.com/yourorg/datalake-free/internal/application/sale/commands"
	"github.com/yourorg/datalake-free/internal/domain/customer"
	"github.com/yourorg/datalake-free/internal/domain/order"
	"github.com/yourorg/datalake-free/internal/domain/sale"
	"github.com/yourorg/datalake-free/internal/infrastructure/postgres"
)

type Handler struct {
	createCustomerHandler *customercmds.CreateCustomerHandler
	updateCustomerHandler *customercmds.UpdateCustomerHandler
	deleteCustomerHandler *customercmds.DeleteCustomerHandler
	createOrderHandler    *ordercmds.CreateOrderHandler
	updateOrderHandler    *ordercmds.UpdateOrderHandler
	deleteOrderHandler    *ordercmds.DeleteOrderHandler
	createSaleHandler     *salecmds.CreateSaleHandler
	updateSaleHandler     *salecmds.UpdateSaleHandler
	deleteSaleHandler     *salecmds.DeleteSaleHandler
	saleRepo              sale.Repository
}

func NewHandler(db *pgxpool.Pool) *Handler {
	customerRepo := postgres.NewCustomerRepository(db)
	orderRepo := postgres.NewOrderRepository(db)
	saleRepo := postgres.NewSalesRepository(db)
	auditLogger := postgres.NewAuditLogger(db)

	return &Handler{
		createCustomerHandler: customercmds.NewCreateCustomerHandler(customerRepo, auditLogger, db),
		updateCustomerHandler: customercmds.NewUpdateCustomerHandler(customerRepo, auditLogger, db),
		deleteCustomerHandler: customercmds.NewDeleteCustomerHandler(customerRepo, auditLogger, db),
		createOrderHandler:    ordercmds.NewCreateOrderHandler(orderRepo, auditLogger, db),
		updateOrderHandler:    ordercmds.NewUpdateOrderHandler(orderRepo, auditLogger, db),
		deleteOrderHandler:    ordercmds.NewDeleteOrderHandler(orderRepo, auditLogger, db),
		createSaleHandler:     salecmds.NewCreateSaleHandler(saleRepo, auditLogger, db),
		updateSaleHandler:     salecmds.NewUpdateSaleHandler(saleRepo, auditLogger, db),
		deleteSaleHandler:     salecmds.NewDeleteSaleHandler(saleRepo, auditLogger, db),
		saleRepo:              saleRepo,
	}
}

func Register(api huma.API, h *Handler) {
	huma.Register(api, huma.Operation{
		OperationID: "create-customer",
		Method:      "POST",
		Path:        "/api/customers",
		Summary:     "Create a new customer",
		Tags:        []string{"customers"},
	}, h.CreateCustomer)
	
	huma.Register(api, huma.Operation{
		OperationID: "update-customer",
		Method:      "PUT",
		Path:        "/api/customers/{id}",
		Summary:     "Update a customer",
		Tags:        []string{"customers"},
	}, h.UpdateCustomer)
	
	huma.Register(api, huma.Operation{
		OperationID: "delete-customer",
		Method:      "DELETE",
		Path:        "/api/customers/{id}",
		Summary:     "Delete a customer",
		Tags:        []string{"customers"},
	}, h.DeleteCustomer)

	huma.Register(api, huma.Operation{
		OperationID: "create-order",
		Method:      "POST",
		Path:        "/api/orders",
		Summary:     "Create a new order",
		Tags:        []string{"orders"},
	}, h.CreateOrder)

	huma.Register(api, huma.Operation{
		OperationID: "update-order",
		Method:      "PUT",
		Path:        "/api/orders/{id}",
		Summary:     "Update an order",
		Tags:        []string{"orders"},
	}, h.UpdateOrder)

	huma.Register(api, huma.Operation{
		OperationID: "delete-order",
		Method:      "DELETE",
		Path:        "/api/orders/{id}",
		Summary:     "Delete an order",
		Tags:        []string{"orders"},
	}, h.DeleteOrder)

	huma.Register(api, huma.Operation{
		OperationID: "createSale",
		Method:      "POST",
		Path:        "/api/sales",
		Summary:     "Create a new sale",
		Tags:        []string{"sales"},
	}, h.CreateSale)

	huma.Register(api, huma.Operation{
		OperationID: "listSales",
		Method:      "GET",
		Path:        "/api/sales",
		Summary:     "List all sales with pagination",
		Tags:        []string{"sales"},
	}, h.ListSales)

	huma.Register(api, huma.Operation{
		OperationID: "getSaleByID",
		Method:      "GET",
		Path:        "/api/sales/{id}",
		Summary:     "Get a sale by ID",
		Tags:        []string{"sales"},
	}, h.GetSaleByID)

	huma.Register(api, huma.Operation{
		OperationID: "listSalesByCustomer",
		Method:      "GET",
		Path:        "/api/sales/customer/{customer_id}",
		Summary:     "List sales for a customer",
		Tags:        []string{"sales"},
	}, h.GetSaleByCustomer)

	huma.Register(api, huma.Operation{
		OperationID: "updateSale",
		Method:      "PUT",
		Path:        "/api/sales/{id}",
		Summary:     "Update a sale",
		Tags:        []string{"sales"},
	}, h.UpdateSale)

	huma.Register(api, huma.Operation{
		OperationID: "deleteSale",
		Method:      "DELETE",
		Path:        "/api/sales/{id}",
		Summary:     "Delete a sale",
		Tags:        []string{"sales"},
	}, h.DeleteSale)
}

func (h *Handler) CreateCustomer(ctx context.Context, input *CreateCustomerInput) (*CreateCustomerOutput, error) {
	cmd := customer.CreateCustomerCommand{
		Name:      input.Body.Name,
		Email:     input.Body.Email,
		Phone:     input.Body.Phone,
		CreatedBy: "system",
	}

	result, err := h.createCustomerHandler.Handle(ctx, cmd)
	if err != nil {
		output := &CreateCustomerOutput{}
		output.Status = 400
		output.Body.Success = false
		output.Body.Error = err.Error()
		return output, nil
	}

	output := &CreateCustomerOutput{}
	output.Status = 200
	output.Body.Success = true
	output.Body.Data = &CustomerData{
		ID:        result.ID,
		Name:      result.Name,
		Email:     result.Email,
		Phone:     result.Phone,
		CreatedBy: result.CreatedBy,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}
	return output, nil
}

func (h *Handler) UpdateCustomer(ctx context.Context, input *UpdateCustomerInput) (*UpdateCustomerOutput, error) {
	cmd := customer.UpdateCustomerCommand{
		ID:        input.ID,
		Name:      input.Body.Name,
		Email:     input.Body.Email,
		Phone:     input.Body.Phone,
		UpdatedBy: "system",
	}

	result, err := h.updateCustomerHandler.Handle(ctx, cmd)
	if err != nil {
		output := &UpdateCustomerOutput{}
		output.Status = 400
		output.Body.Success = false
		output.Body.Error = err.Error()
		return output, nil
	}

	output := &UpdateCustomerOutput{}
	output.Status = 200
	output.Body.Success = true
	output.Body.Data = &CustomerData{
		ID:        result.ID,
		Name:      result.Name,
		Email:     result.Email,
		Phone:     result.Phone,
		CreatedBy: result.CreatedBy,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}
	return output, nil
}

func (h *Handler) DeleteCustomer(ctx context.Context, input *DeleteCustomerInput) (*DeleteCustomerOutput, error) {
	cmd := customer.DeleteCustomerCommand{
		ID:        input.ID,
		DeletedBy: "system",
	}

	err := h.deleteCustomerHandler.Handle(ctx, cmd)
	if err != nil {
		output := &DeleteCustomerOutput{}
		output.Status = 404
		output.Body.Success = false
		output.Body.Error = err.Error()
		return output, nil
	}

	output := &DeleteCustomerOutput{}
	output.Status = 200
	output.Body.Success = true
	output.Body.Data = &DeleteMessageData{
		Message: fmt.Sprintf("Customer %d deleted successfully", input.ID),
	}
	return output, nil
}

func (h *Handler) CreateOrder(ctx context.Context, input *CreateOrderInput) (*CreateOrderOutput, error) {
	cmd := order.CreateOrderCommand{
		CustomerID:  input.Body.CustomerID,
		TotalAmount: input.Body.TotalAmount,
		CreatedBy:   "system",
	}

	result, err := h.createOrderHandler.Handle(ctx, cmd)
	if err != nil {
		output := &CreateOrderOutput{}
		output.Status = 400
		output.Body.Success = false
		output.Body.Error = err.Error()
		return output, nil
	}

	output := &CreateOrderOutput{}
	output.Status = 200
	output.Body.Success = true
	output.Body.Data = &OrderData{
		ID:          result.ID,
		CustomerID:  result.CustomerID,
		TotalAmount: result.TotalAmount,
		Status:      result.Status,
		CreatedBy:   result.CreatedBy,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}
	return output, nil
}

func (h *Handler) UpdateOrder(ctx context.Context, input *UpdateOrderInput) (*UpdateOrderOutput, error) {
	cmd := order.UpdateOrderCommand{
		ID:          input.ID,
		TotalAmount: input.Body.TotalAmount,
		Status:      input.Body.Status,
		UpdatedBy:   "system",
	}

	result, err := h.updateOrderHandler.Handle(ctx, cmd)
	if err != nil {
		output := &UpdateOrderOutput{}
		output.Status = 400
		output.Body.Success = false
		output.Body.Error = err.Error()
		return output, nil
	}

	output := &UpdateOrderOutput{}
	output.Status = 200
	output.Body.Success = true
	output.Body.Data = &OrderData{
		ID:          result.ID,
		CustomerID:  result.CustomerID,
		TotalAmount: result.TotalAmount,
		Status:      result.Status,
		CreatedBy:   result.CreatedBy,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}
	return output, nil
}

func (h *Handler) DeleteOrder(ctx context.Context, input *DeleteOrderInput) (*DeleteOrderOutput, error) {
	cmd := order.DeleteOrderCommand{
		ID:        input.ID,
		DeletedBy: "system",
	}

	err := h.deleteOrderHandler.Handle(ctx, cmd)
	if err != nil {
		output := &DeleteOrderOutput{}
		output.Status = 404
		output.Body.Success = false
		output.Body.Error = err.Error()
		return output, nil
	}

	output := &DeleteOrderOutput{}
	output.Status = 200
	output.Body.Success = true
	output.Body.Data = &DeleteMessageData{
		Message: fmt.Sprintf("Order %d deleted successfully", input.ID),
	}
	return output, nil
}

func (h *Handler) CreateSale(ctx context.Context, input *CreateSaleInput) (*CreateSaleOutput, error) {
	cmd := sale.CreateSaleCommand{
		CustomerID: input.Body.CustomerID,
		Amount:     input.Body.Amount,
		CreatedBy:  "system",
	}

	result, err := h.createSaleHandler.Handle(ctx, cmd)
	if err != nil {
		output := &CreateSaleOutput{}
		output.Status = 400
		output.Body.Success = false
		output.Body.Error = err.Error()
		return output, nil
	}

	output := &CreateSaleOutput{}
	output.Status = 200
	output.Body.Success = true
	output.Body.Data = &SaleData{
		ID:         result.ID,
		CustomerID: result.CustomerID,
		Amount:     result.Amount,
		Status:     result.Status,
		CreatedBy:  result.CreatedBy,
		CreatedAt:  result.CreatedAt,
		UpdatedAt:  result.UpdatedAt,
	}
	return output, nil
}

func (h *Handler) ListSales(ctx context.Context, input *ListSalesInput) (*ListSalesOutput, error) {
	sales, total, err := h.saleRepo.List(ctx, input.Limit, input.Offset)
	if err != nil {
		output := &ListSalesOutput{}
		output.Status = 500
		output.Body.Success = false
		output.Body.Error = err.Error()
		return output, nil
	}

	output := &ListSalesOutput{}
	output.Status = 200
	output.Body.Success = true
	output.Body.Total = total
	output.Body.Data = make([]SaleData, len(sales))
	for i, s := range sales {
		output.Body.Data[i] = SaleData{
			ID:         s.ID,
			CustomerID: s.CustomerID,
			Amount:     s.Amount,
			Status:     s.Status,
			CreatedBy:  s.CreatedBy,
			CreatedAt:  s.CreatedAt,
			UpdatedAt:  s.UpdatedAt,
		}
	}
	return output, nil
}

func (h *Handler) GetSaleByID(ctx context.Context, input *struct {
	ID int64 `path:"id" doc:"Sale ID" required:"true"`
}) (*GetSaleByIDOutput, error) {
	s, err := h.saleRepo.GetByID(ctx, input.ID)
	if err != nil {
		output := &GetSaleByIDOutput{}
		output.Status = 404
		output.Body.Success = false
		output.Body.Error = "sale not found"
		return output, nil
	}

	output := &GetSaleByIDOutput{}
	output.Status = 200
	output.Body.Success = true
	output.Body.Data = &SaleData{
		ID:         s.ID,
		CustomerID: s.CustomerID,
		Amount:     s.Amount,
		Status:     s.Status,
		CreatedBy:  s.CreatedBy,
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
	}
	return output, nil
}

func (h *Handler) GetSaleByCustomer(ctx context.Context, input *GetSaleByCustomerInput) (*GetSaleByCustomerOutput, error) {
	sales, err := h.saleRepo.GetByCustomerID(ctx, input.CustomerID)
	if err != nil {
		output := &GetSaleByCustomerOutput{}
		output.Status = 500
		output.Body.Success = false
		output.Body.Error = err.Error()
		return output, nil
	}

	// Apply pagination manually
	total := len(sales)
	start := input.Offset
	end := input.Offset + input.Limit
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	paginatedSales := sales[start:end]

	output := &GetSaleByCustomerOutput{}
	output.Status = 200
	output.Body.Success = true
	output.Body.Total = total
	output.Body.Data = make([]SaleData, len(paginatedSales))
	for i, s := range paginatedSales {
		output.Body.Data[i] = SaleData{
			ID:         s.ID,
			CustomerID: s.CustomerID,
			Amount:     s.Amount,
			Status:     s.Status,
			CreatedBy:  s.CreatedBy,
			CreatedAt:  s.CreatedAt,
			UpdatedAt:  s.UpdatedAt,
		}
	}
	return output, nil
}

func (h *Handler) UpdateSale(ctx context.Context, input *UpdateSaleInput) (*UpdateSaleOutput, error) {
	cmd := sale.UpdateSaleCommand{
		ID:        input.ID,
		Amount:    input.Body.Amount,
		Status:    input.Body.Status,
		UpdatedBy: "system",
	}

	result, err := h.updateSaleHandler.Handle(ctx, cmd)
	if err != nil {
		output := &UpdateSaleOutput{}
		output.Status = 400
		output.Body.Success = false
		output.Body.Error = err.Error()
		return output, nil
	}

	output := &UpdateSaleOutput{}
	output.Status = 200
	output.Body.Success = true
	output.Body.Data = &SaleData{
		ID:         result.ID,
		CustomerID: result.CustomerID,
		Amount:     result.Amount,
		Status:     result.Status,
		CreatedBy:  result.CreatedBy,
		CreatedAt:  result.CreatedAt,
		UpdatedAt:  result.UpdatedAt,
	}
	return output, nil
}

func (h *Handler) DeleteSale(ctx context.Context, input *DeleteSaleInput) (*DeleteSaleOutput, error) {
	cmd := sale.DeleteSaleCommand{
		ID:        input.ID,
		DeletedBy: "system",
	}

	err := h.deleteSaleHandler.Handle(ctx, cmd)
	if err != nil {
		output := &DeleteSaleOutput{}
		output.Status = 404
		output.Body.Success = false
		output.Body.Error = err.Error()
		return output, nil
	}

	output := &DeleteSaleOutput{}
	output.Status = 200
	output.Body.Success = true
	return output, nil
}

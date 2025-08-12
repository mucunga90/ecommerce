package api

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
	"github.com/mucunga90/ecommerce/internal"
)

type orderPayload struct {
	CustomerID uuid.UUID `json:"customer_id" validate:"required"`
	Items      []struct {
		ProductID uuid.UUID `json:"product_id" validate:"required"`
		Price     float64   `json:"price" validate:"required,gt=0"`
		Quantity  int       `json:"quantity" validate:"required,min=1"`
	} `json:"items" validate:"required,min=1"`
}

func (p *orderPayload) asDomain() *internal.Order {
	items := make([]internal.OrderItem, len(p.Items))
	for i := range p.Items {
		items[i] = internal.OrderItem{
			ProductID: p.Items[i].ProductID,
			Quantity:  p.Items[i].Quantity,
			UnitPrice: p.Items[i].Price,
		}
	}
	return &internal.Order{
		CustomerID: p.CustomerID,
		Items:      items,
	}
}

type orderHandler struct {
	manager orderManager
}

func NewOrderHandler(manager orderManager) *orderHandler {
	return &orderHandler{manager: manager}
}

func (h *orderHandler) MakeOrder(c echo.Context) error {
	defer func() {
		if r := recover(); r != nil {
			c.Logger().Errorf("Panic in CreateOrder: %v", r)
			_ = c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "internal server error",
			})
		}
	}()

	var payload orderPayload

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
	}

	if len(payload.Items) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "empty order list"})
	}

	if err := h.manager.CreateOrder(payload.asDomain()); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "order created",
	})
}

type orderManager interface {
	CreateOrder(o *internal.Order) error
}

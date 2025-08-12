package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// priceHandler handles price-related requests
type priceHandler struct {
	manager priceManager
}

// NewPriceHandler creates a new price handler
func NewPriceHandler(manager priceManager) *priceHandler {
	return &priceHandler{manager: manager}
}

// ProductPrices handles the request to get product prices for a specific category
func (h *priceHandler) ProductPrices(c echo.Context) error {
	defer func() {
		if r := recover(); r != nil {
			c.Logger().Errorf("Panic in fetching product prices: %v", r)
			_ = c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "internal server error",
			})
		}
	}()

	category := c.QueryParam("category")
	if category == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "category is required"})
	}

	avgPrice, err := h.manager.ProductAveragePrice(category)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"category":      category,
		"average_price": avgPrice,
	})
}

type priceManager interface {
	ProductAveragePrice(category string) (float64, error)
}

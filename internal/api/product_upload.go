package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ProductPayload struct {
	Name       string   `json:"name" validate:"required"`
	Price      float64  `json:"price" validate:"required"`
	Categories []string `json:"categories" validate:"required,min=1"`
}

func (p *ProductPayload) Validate() error {
	if p.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name is required")
	}
	if len(p.Categories) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "at least one category is required")
	}
	return nil
}

type uploadHandler struct {
	manager uploadManager
}

// NewUploadHandler creates a new API handler with the given DB
func NewUploadHandler(manager uploadManager) *uploadHandler {
	return &uploadHandler{manager: manager}
}

func (h *uploadHandler) ProductUpload(c echo.Context) error {
	defer func() {
		if r := recover(); r != nil {
			c.Logger().Errorf("Panic in uploading products: %v", r)
			_ = c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "internal server error",
			})
		}
	}()

	var payload []ProductPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
	}

	if len(payload) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "empty product list"})
	}

	if err := h.manager.ProductUpload(payload); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "products created",
	})
}

type uploadManager interface {
	ProductUpload(items []ProductPayload) error
}

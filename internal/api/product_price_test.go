package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/mucunga90/ecommerce/internal/api/mocks"
	"github.com/stretchr/testify/require"
)

func TestProductPrices_Success(t *testing.T) {
	e := echo.New()
	manager := &mocks.MockpriceManager{}
	handler := NewPriceHandler(manager)

	categoryID := "electronics"
	manager.On("ProductAveragePrice", categoryID).Return(123.45, nil)

	req := httptest.NewRequest(http.MethodGet, "/prices/"+categoryID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("category_id")
	c.SetParamValues(categoryID)

	err := handler.ProductPrices(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.JSONEq(t, `{"category_id":"electronics","average_price":123.45}`, rec.Body.String())
	manager.AssertExpectations(t)
}

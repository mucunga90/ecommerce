package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mucunga90/ecommerce/internal/api/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMakeOrder_Success(t *testing.T) {
	e := echo.New()
	manager := &mocks.MockorderManager{}
	manager.On("CreateOrder", mock.AnythingOfType("*internal.Order")).Return(nil)

	handler := NewOrderHandler(manager)

	payload := orderPayload{
		CustomerID: uuid.New(),
		Items: []struct {
			ProductID uuid.UUID `json:"product_id" validate:"required"`
			Price     float64   `json:"price" validate:"required,gt=0"`
			Quantity  int       `json:"quantity" validate:"required,min=1"`
		}{
			{
				ProductID: uuid.New(),
				Price:     10.5,
				Quantity:  2,
			},
		},
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.MakeOrder(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "order created")
	manager.AssertCalled(t, "CreateOrder", mock.AnythingOfType("*internal.Order"))
	manager.AssertExpectations(t) // <-- Added
}

func TestMakeOrder_ManagerError(t *testing.T) {
	e := echo.New()
	manager := &mocks.MockorderManager{}
	manager.On("CreateOrder", mock.AnythingOfType("*internal.Order")).Return(echo.NewHTTPError(http.StatusInternalServerError, "db error"))

	handler := NewOrderHandler(manager)

	payload := orderPayload{
		CustomerID: uuid.New(),
		Items: []struct {
			ProductID uuid.UUID `json:"product_id" validate:"required"`
			Price     float64   `json:"price" validate:"required,gt=0"`
			Quantity  int       `json:"quantity" validate:"required,min=1"`
		}{
			{
				ProductID: uuid.New(),
				Price:     10.5,
				Quantity:  2,
			},
		},
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.MakeOrder(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.Contains(t, rec.Body.String(), "db error")
	manager.AssertCalled(t, "CreateOrder", mock.AnythingOfType("*internal.Order"))
	manager.AssertExpectations(t) // <-- Added
}

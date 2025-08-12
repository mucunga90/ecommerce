package storage

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mucunga90/ecommerce/internal"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupOrderTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost port=5433 user=test password=test dbname=testdb sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&internal.Customer{}, &internal.Order{}))
	return db
}

func TestCreateOrder(t *testing.T) {
	db := setupOrderTestDB(t)
	s := New(db)

	// Create a customer
	customer := &internal.Customer{
		ID:    uuid.New(),
		Name:  "Test User",
		Email: "test@example.com",
	}
	require.NoError(t, db.Create(customer).Error)

	// Create an order
	order := &internal.Order{
		ID:         uuid.New(),
		CustomerID: customer.ID,
		Total:      99.99,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}
	err := s.CreateOrder(order)
	require.NoError(t, err)

	// Verify order exists
	var found internal.Order
	err = db.First(&found, "id = ?", order.ID).Error
	require.NoError(t, err)
	require.Equal(t, order.CustomerID, found.CustomerID)
	require.Equal(t, order.Total, found.Total)
}

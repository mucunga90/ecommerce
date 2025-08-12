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

func setupProductTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost port=5433 user=test password=test dbname=testdb sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&internal.Category{}, &internal.Product{}))
	return db
}

func TestCreateProductAndGetAveragePriceForCategory(t *testing.T) {
	db := setupProductTestDB(t)
	s := New(db)

	// Create category tree: Electronics -> Phones
	rootCat := &internal.Category{ID: uuid.New(), Name: "Electronics"}
	require.NoError(t, db.Create(rootCat).Error)

	childCat := &internal.Category{ID: uuid.New(), Name: "Phones", ParentID: &rootCat.ID}
	require.NoError(t, db.Create(childCat).Error)

	// Create products in both categories
	p1 := &internal.Product{
		ID:         uuid.New(),
		Name:       "Laptop",
		CategoryID: rootCat.ID,
		Price:      1000.0,
		CreatedAt:  time.Now(),
	}
	p2 := &internal.Product{
		ID:         uuid.New(),
		Name:       "Smartphone",
		CategoryID: childCat.ID,
		Price:      500.0,
		CreatedAt:  time.Now(),
	}
	require.NoError(t, s.CreateProduct(p1))
	require.NoError(t, s.CreateProduct(p2))

	// Average for root should include both products
	avg, err := s.GetAveragePriceForCategory(rootCat.ID)
	require.NoError(t, err)
	require.InDelta(t, 750.0, avg, 0.01)

	// Average for child should only include smartphone
	avgChild, err := s.GetAveragePriceForCategory(childCat.ID)
	require.NoError(t, err)
	require.InDelta(t, 500.0, avgChild, 0.01)
}

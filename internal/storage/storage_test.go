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

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost port=5433 user=test password=test dbname=testdb sslmode=disable"
	var db *gorm.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	require.NoError(t, err)

	// Auto-migrate the schema
	err = db.AutoMigrate(&internal.Category{})
	require.NoError(t, err)

	return db

}

func TestCategoryStorageIntegration(t *testing.T) {
	db := setupTestDB(t)
	s := New(db)

	// 1. Create a root category
	rootCategory := &internal.Category{
		ID:   uuid.New(),
		Name: "Electronics",
	}
	err := s.CreateCategory(rootCategory)
	require.NoError(t, err)

	// 2. Get category by name
	found, err := s.GetCategoryByName("Electronics", nil)
	require.NoError(t, err)
	require.NotNil(t, found)
	require.Equal(t, rootCategory.ID, found.ID)

	// 3. Create a subcategory
	childCategory := &internal.Category{
		ID:       uuid.New(),
		Name:     "Phones",
		ParentID: &rootCategory.ID,
	}
	err = s.CreateCategory(childCategory)
	require.NoError(t, err)

	// 4. Get category tree from root
	tree, err := s.GetCategoryTree(rootCategory.ID)
	require.NoError(t, err)
	require.Len(t, tree, 2) // root + child

	// Verify IDs in tree
	var ids []uuid.UUID
	for _, c := range tree {
		ids = append(ids, c.ID)
	}
	require.Contains(t, ids, rootCategory.ID)
	require.Contains(t, ids, childCategory.ID)
}

package database

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDatabaseConnection(t *testing.T) {
	// Use test database connection string
	dsn := "host=localhost port=5433 user=test password=test dbname=testdb sslmode=disable"
	db, err := New(dsn)
	require.NoError(t, err)
	require.NotNil(t, db)

	// Ping the database to verify connection
	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Ping())
}

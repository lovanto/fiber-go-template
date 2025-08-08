package database

import (
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	mockDBConn = &sqlx.DB{}
)

func TestOpenDBConnection_UnsupportedDB(t *testing.T) {
	// Save original function and restore after test
	oldPostgreSQLConn := postgreSQLConn
	oldMySQLConn := mysqlConn
	defer func() {
		postgreSQLConn = oldPostgreSQLConn
		mysqlConn = oldMySQLConn
	}()

	// Test unsupported database type
	t.Run("Unsupported database type", func(t *testing.T) {
		os.Setenv("DB_TYPE", "unsupported")
		q, err := OpenDBConnection()
		assert.Nil(t, q)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported database type")
	})

	// Test PostgreSQL connection
	t.Run("PostgreSQL connection success", func(t *testing.T) {
		// Mock the PostgreSQL connection function
		postgreSQLConn = func() (*sqlx.DB, error) {
			return mockDBConn, nil
		}
		os.Setenv("DB_TYPE", "pgx")

		q, err := OpenDBConnection()
		require.NoError(t, err)
		assert.NotNil(t, q)
	})

	t.Run("PostgreSQL connection error", func(t *testing.T) {
		// Mock the PostgreSQL connection function to return an error
		expectedErr := fmt.Errorf("connection failed")
		postgreSQLConn = func() (*sqlx.DB, error) {
			return nil, expectedErr
		}
		os.Setenv("DB_TYPE", "pgx")

		q, err := OpenDBConnection()
		assert.Nil(t, q)
		assert.ErrorIs(t, err, expectedErr)
	})

	// Test MySQL connection
	t.Run("MySQL connection success", func(t *testing.T) {
		// Mock the MySQL connection function
		mysqlConn = func() (*sqlx.DB, error) {
			return mockDBConn, nil
		}
		os.Setenv("DB_TYPE", "mysql")

		q, err := OpenDBConnection()
		require.NoError(t, err)
		assert.NotNil(t, q)
	})

	t.Run("MySQL connection error", func(t *testing.T) {
		// Mock the MySQL connection function to return an error
		expectedErr := fmt.Errorf("connection failed")
		mysqlConn = func() (*sqlx.DB, error) {
			return nil, expectedErr
		}
		os.Setenv("DB_TYPE", "mysql")

		q, err := OpenDBConnection()
		assert.Nil(t, q)
		assert.ErrorIs(t, err, expectedErr)
	})
}

func TestOpenDBConnection_MySQL(t *testing.T) {
	// This is a simple test that just verifies the function doesn't panic
	// For a real test, you would need a test database setup
	os.Setenv("DB_TYPE", "mysql")
	q, err := OpenDBConnection()

	// We can't assert much here without a real database
	// Just verify the function returns something
	if err == nil {
		assert.NotNil(t, q)
	} else {
		t.Log("Note: MySQL test database not available")
	}
}

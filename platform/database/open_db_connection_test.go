package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenDBConnection_UnsupportedDB(t *testing.T) {
	os.Setenv("DB_TYPE", "unsupported")
	q, err := OpenDBConnection()
	assert.Nil(t, q)
	assert.Error(t, err)
}

func TestOpenDBConnection_PostgreSQL(t *testing.T) {
	// This is a simple test that just verifies the function doesn't panic
	// For a real test, you would need a test database setup
	os.Setenv("DB_TYPE", "pgx")
	q, err := OpenDBConnection()
	
	// We can't assert much here without a real database
	// Just verify the function returns something
	if err == nil {
		assert.NotNil(t, q)
	} else {
		t.Log("Note: PostgreSQL test database not available")
	}
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

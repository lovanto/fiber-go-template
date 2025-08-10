package database

import (
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var mockDBConn = &sqlx.DB{}

func testUnsupportedDBType(t *testing.T) {
	os.Setenv("DB_TYPE", "unsupported")
	q, err := OpenDBConnection()
	assert.Nil(t, q)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported database type")
}

func testPostgreSQLConnSuccess(t *testing.T) {
	postgreSQLConn = func(builders ...func(string) (string, error)) (*sqlx.DB, error) {
		return mockDBConn, nil
	}
	os.Setenv("DB_TYPE", "pgx")

	q, err := OpenDBConnection()
	require.NoError(t, err)
	assert.NotNil(t, q)
}

func testPostgreSQLConnError(t *testing.T) {
	expectedErr := fmt.Errorf("connection failed")
	postgreSQLConn = func(builders ...func(string) (string, error)) (*sqlx.DB, error) {
		return nil, expectedErr
	}
	os.Setenv("DB_TYPE", "pgx")

	q, err := OpenDBConnection()
	assert.Nil(t, q)
	assert.ErrorIs(t, err, expectedErr)
}

func testMySQLConnSuccess(t *testing.T) {
	mysqlConn = func() (*sqlx.DB, error) {
		return mockDBConn, nil
	}
	os.Setenv("DB_TYPE", "mysql")

	q, err := OpenDBConnection()
	require.NoError(t, err)
	assert.NotNil(t, q)
}

func testMySQLConnError(t *testing.T) {
	expectedErr := fmt.Errorf("connection failed")
	mysqlConn = func() (*sqlx.DB, error) {
		return nil, expectedErr
	}
	os.Setenv("DB_TYPE", "mysql")

	q, err := OpenDBConnection()
	assert.Nil(t, q)
	assert.ErrorIs(t, err, expectedErr)
}

func TestOpenDBConnection(t *testing.T) {
	oldPostgreSQLConn := postgreSQLConn
	oldMySQLConn := mysqlConn
	defer func() {
		postgreSQLConn = oldPostgreSQLConn
		mysqlConn = oldMySQLConn
	}()

	t.Run("Unsupported DB type", testUnsupportedDBType)
	t.Run("PostgreSQL connection success", testPostgreSQLConnSuccess)
	t.Run("PostgreSQL connection error", testPostgreSQLConnError)
	t.Run("MySQL connection success", testMySQLConnSuccess)
	t.Run("MySQL connection error", testMySQLConnError)
}

func TestOpenDBConnection_MySQL_NoPanic(t *testing.T) {
	os.Setenv("DB_TYPE", "mysql")
	q, err := OpenDBConnection()

	if err == nil {
		assert.NotNil(t, q)
	} else {
		t.Log("Note: MySQL test database not available")
	}
}

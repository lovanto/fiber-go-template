package database

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func newMockDB() (*sqlx.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		return nil, nil, err
	}
	return sqlx.NewDb(db, "pgx"), mock, nil
}

func setupCommonEnv() {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "user")
	os.Setenv("DB_PASSWORD", "pass")
	os.Setenv("DB_NAME", "dbname")
	os.Setenv("DB_SSL_MODE", "disable")
	os.Setenv("DB_MAX_CONNECTIONS", "10")
	os.Setenv("DB_MAX_IDLE_CONNECTIONS", "5")
	os.Setenv("DB_MAX_LIFETIME_CONNECTIONS", "300")
}

func testMissingEnvVars(t *testing.T) {
	origHost := os.Getenv("DB_HOST")
	os.Unsetenv("DB_HOST")
	defer os.Setenv("DB_HOST", origHost)

	_, err := PostgreSQLConnection()
	assert.Error(t, err)
}

func testCustomBuilderError(t *testing.T) {
	builder := func(dbType string) (string, error) {
		return "", errors.New("builder fail")
	}
	_, err := PostgreSQLConnection(builder)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error building connection URL")
}

func testConnectionError(t *testing.T) {
	origConnect := sqlxConnectFunc
	defer func() { sqlxConnectFunc = origConnect }()

	sqlxConnectFunc = func(driverName, dataSourceName string) (*sqlx.DB, error) {
		return nil, errors.New("boom connect")
	}

	_, err := PostgreSQLConnection()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected to database")
}

func testPingError(t *testing.T) {
	origConnect := sqlxConnectFunc
	defer func() { sqlxConnectFunc = origConnect }()

	mockDB, mock, _ := newMockDB()
	mock.ExpectPing().WillReturnError(errors.New("ping fail"))

	sqlxConnectFunc = func(driverName, dataSourceName string) (*sqlx.DB, error) {
		return mockDB, nil
	}

	_, err := PostgreSQLConnection()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not sent ping to database")
}

func testSuccessfulConnection(t *testing.T) {
	origConnect := sqlxConnectFunc
	defer func() { sqlxConnectFunc = origConnect }()

	mockDB, mock, _ := newMockDB()
	mock.ExpectPing().WillReturnError(nil) // successful ping

	sqlxConnectFunc = func(driverName, dataSourceName string) (*sqlx.DB, error) {
		return mockDB, nil
	}

	db, err := PostgreSQLConnection()
	assert.NoError(t, err)
	assert.NotNil(t, db)
	time.Sleep(1 * time.Millisecond) // allow SetConnMaxLifetime to be invoked
}

func TestPostgreSQLConnection_FullCoverage(t *testing.T) {
	setupCommonEnv()

	t.Run("missing env vars", testMissingEnvVars)
	t.Run("custom builder error", testCustomBuilderError)
	t.Run("connection error", testConnectionError)
	t.Run("ping error", testPingError)
	t.Run("successful connection", testSuccessfulConnection)
}

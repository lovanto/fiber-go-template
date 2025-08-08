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

// helper to create a mocked *sqlx.DB
func newMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		return nil, nil, err
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	return sqlx.NewDb(db, "pgx"), mock, nil
}

func TestPostgreSQLConnection_FullCoverage(t *testing.T) {
	// common env vars for success cases
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "user")
	os.Setenv("DB_PASSWORD", "pass")
	os.Setenv("DB_NAME", "dbname")
	os.Setenv("DB_SSL_MODE", "disable")
	os.Setenv("DB_MAX_CONNECTIONS", "10")
	os.Setenv("DB_MAX_IDLE_CONNECTIONS", "5")
	os.Setenv("DB_MAX_LIFETIME_CONNECTIONS", "300")

	// 1️⃣ Missing required env vars → ConnectionURLBuilder fails
	t.Run("missing env vars", func(t *testing.T) {
		origHost := os.Getenv("DB_HOST")
		os.Unsetenv("DB_HOST")
		defer os.Setenv("DB_HOST", origHost)

		_, err := PostgreSQLConnection()
		assert.Error(t, err)
	})

	// 2️⃣ Connection error → sqlx.Connect fails
	t.Run("connection error", func(t *testing.T) {
		origConnect := sqlxConnectFunc
		defer func() { sqlxConnectFunc = origConnect }()
		sqlxConnectFunc = func(driverName, dataSourceName string) (*sqlx.DB, error) {
			return nil, errors.New("boom connect")
		}

		_, err := PostgreSQLConnection()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not connected to database")
	})

	// 3️⃣ Ping error after connect
	t.Run("ping error", func(t *testing.T) {
		origConnect := sqlxConnectFunc
		defer func() { sqlxConnectFunc = origConnect }()

		mockDB, mock, _ := newMockDB(t)
		mock.ExpectPing().WillReturnError(errors.New("ping fail"))

		sqlxConnectFunc = func(driverName, dataSourceName string) (*sqlx.DB, error) {
			return mockDB, nil
		}

		_, err := PostgreSQLConnection()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not sent ping to database")
	})

	// 4️⃣ Successful connection
	t.Run("successful connection", func(t *testing.T) {
		origConnect := sqlxConnectFunc
		defer func() { sqlxConnectFunc = origConnect }()

		mockDB, mock, _ := newMockDB(t)
		mock.ExpectPing().WillReturnError(nil) // successful ping

		sqlxConnectFunc = func(driverName, dataSourceName string) (*sqlx.DB, error) {
			return mockDB, nil
		}

		db, err := PostgreSQLConnection()
		assert.NoError(t, err)
		assert.NotNil(t, db)
		time.Sleep(1 * time.Millisecond) // allow SetConnMaxLifetime to be invoked
	})
}

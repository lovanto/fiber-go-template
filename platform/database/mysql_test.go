package database

import (
	"errors"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func setEnvVarsMySQL() {
	os.Setenv("DB_MAX_CONNECTIONS", "10")
	os.Setenv("DB_MAX_IDLE_CONNECTIONS", "5")
	os.Setenv("DB_MAX_LIFETIME_CONNECTIONS", "60000000000") // 1 min in ns
}

func fakeBuilderSuccess(driver string) (string, error) {
	return "mock_dsn", nil
}

func fakeBuilderError(driver string) (string, error) {
	return "", errors.New("builder error")
}

func TestMysqlConnection_Success(t *testing.T) {
	setEnvVarsMySQL()
	dbMock, mock, _ := sqlmock.New(sqlmock.MonitorPingsOption(true))
	mock.ExpectPing()

	sqlxConnect = func(driverName, dataSourceName string) (*sqlx.DB, error) {
		return sqlx.NewDb(dbMock, driverName), nil
	}
	defer func() { sqlxConnect = sqlx.Connect }()

	db, err := MysqlConnection(fakeBuilderSuccess)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMysqlConnection_SuccessCustomBuilder(t *testing.T) {
	setEnvVarsMySQL()
	dbMock, mock, _ := sqlmock.New(sqlmock.MonitorPingsOption(true))
	mock.ExpectPing()

	sqlxConnect = func(driverName, dataSourceName string) (*sqlx.DB, error) {
		return sqlx.NewDb(dbMock, driverName), nil
	}
	defer func() { sqlxConnect = sqlx.Connect }()

	db, err := MysqlConnection(fakeBuilderSuccess)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMysqlConnection_BuilderError(t *testing.T) {
	setEnvVarsMySQL()
	_, err := MysqlConnection(fakeBuilderError)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "builder error")
}

func TestMysqlConnection_ConnectError(t *testing.T) {
	setEnvVarsMySQL()
	sqlxConnect = func(driverName, dataSourceName string) (*sqlx.DB, error) {
		return nil, errors.New("connect fail")
	}
	defer func() { sqlxConnect = sqlx.Connect }()

	_, err := MysqlConnection(fakeBuilderSuccess)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connect fail")
}

func TestMysqlConnection_PingError(t *testing.T) {
	setEnvVarsMySQL()
	dbMock, mock, _ := sqlmock.New(sqlmock.MonitorPingsOption(true))
	mock.ExpectPing().WillReturnError(errors.New("ping fail"))

	sqlxConnect = func(driverName, dataSourceName string) (*sqlx.DB, error) {
		return sqlx.NewDb(dbMock, driverName), nil
	}
	defer func() { sqlxConnect = sqlx.Connect }()

	_, err := MysqlConnection(fakeBuilderSuccess)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ping fail")
}

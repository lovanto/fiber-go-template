package database

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/create-go-app/fiber-go-template/pkg/utils/connection_url_builder"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// allows test override
var sqlxConnect = sqlx.Connect

func MysqlConnection(builders ...func(string) (string, error)) (*sqlx.DB, error) {
	maxConn, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	maxIdleConn, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	maxLifetimeConn, _ := strconv.Atoi(os.Getenv("DB_MAX_LIFETIME_CONNECTIONS"))

	// Use the provided builder or default
	var builder func(string) (string, error) = connection_url_builder.ConnectionURLBuilder
	if len(builders) > 0 && builders[0] != nil {
		builder = builders[0]
	}

	mysqlConnURL, err := builder("mysql")
	if err != nil {
		return nil, fmt.Errorf("error building connection URL: %w", err)
	}

	db, err := sqlxConnect("mysql", mysqlConnURL)
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	db.SetMaxOpenConns(maxConn)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxLifetime(time.Duration(maxLifetimeConn))

	if err := db.Ping(); err != nil {
		defer db.Close()
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}

	return db, nil
}

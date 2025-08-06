package database

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/create-go-app/fiber-go-template/pkg/utils"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

func MysqlConnection() (*sqlx.DB, error) {
	maxConn, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	maxIdleConn, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	maxLifetimeConn, _ := strconv.Atoi(os.Getenv("DB_MAX_LIFETIME_CONNECTIONS"))

	mysqlConnURL, err := utils.ConnectionURLBuilder("mysql")
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Connect("mysql", mysqlConnURL)
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

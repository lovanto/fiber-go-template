package database

import (
	"fmt"
	"os"

	"github.com/create-go-app/fiber-go-template/app/queries"
	"github.com/jmoiron/sqlx"
)

// Queries combines all the query types for our database operations
type Queries struct {
	*queries.UserQueries
	*queries.BookQueries
}

// These function variables allow us to mock the database connections in tests
var (
	postgreSQLConn = PostgreSQLConnection
	mysqlConn      = MysqlConnection
)

func OpenDBConnection() (*Queries, error) {
	var (
		db  *sqlx.DB
		err error
	)

	dbType := os.Getenv("DB_TYPE")

	switch dbType {
	case "pgx":
		db, err = postgreSQLConn()
	case "mysql":
		db, err = mysqlConn()
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	if err != nil {
		return nil, err
	}

	return &Queries{
		UserQueries: &queries.UserQueries{DB: db},
		BookQueries: &queries.BookQueries{DB: db},
	}, nil
}

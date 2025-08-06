package database

import (
	"os"

	"github.com/create-go-app/fiber-go-template/app/queries"
	"github.com/jmoiron/sqlx"
)

type Queries struct {
	*queries.UserQueries
	*queries.BookQueries
}

func OpenDBConnection() (*Queries, error) {
	var (
		db  *sqlx.DB
		err error
	)

	dbType := os.Getenv("DB_TYPE")

	switch dbType {
	case "pgx":
		db, err = PostgreSQLConnection()
	case "mysql":
		db, err = MysqlConnection()
	}

	if err != nil {
		return nil, err
	}

	return &Queries{
		UserQueries: &queries.UserQueries{DB: db},
		BookQueries: &queries.BookQueries{DB: db},
	}, nil
}

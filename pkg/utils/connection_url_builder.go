package utils

import (
	"fmt"
	"os"
)

func ConnectionURLBuilder(n string) (string, error) {
	var url string

	switch n {
	// URL for PostgreSQL connection.
	case "postgres":
		url = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_SSL_MODE"),
		)
	// URL for MySQL connection.
	case "mysql":
		url = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)

	// URL for Redis connection.
	case "redis":
		url = fmt.Sprintf(
			"%s:%s",
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT"),
		)
	// URL for Fiber connection.
	case "fiber":
		url = fmt.Sprintf(
			"%s:%s",
			os.Getenv("SERVER_HOST"),
			os.Getenv("SERVER_PORT"),
		)
	// Default case.
	default:
		return "", fmt.Errorf("connection name '%v' is not supported", n)
	}

	return url, nil
}

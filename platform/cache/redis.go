package cache

import (
	"os"
	"strconv"

	"github.com/create-go-app/fiber-go-template/pkg/utils/connection_url_builder"
	"github.com/redis/go-redis/v9"
)

// URLBuilder is a function type that builds a connection URL for a given service
type URLBuilder func(service string) (string, error)

// DefaultURLBuilder is the default implementation of URLBuilder that uses connection_url_builder
var DefaultURLBuilder = connection_url_builder.ConnectionURLBuilder

// RedisConnection establishes a connection to Redis using the provided URL builder
func RedisConnection() (*redis.Client, error) {
	return NewRedisConnection(DefaultURLBuilder)
}

// NewRedisConnection creates a new Redis connection using the provided URL builder
func NewRedisConnection(urlBuilder URLBuilder) (*redis.Client, error) {
	dbNumber, _ := strconv.Atoi(os.Getenv("REDIS_DB_NUMBER"))
	redisConnURL, err := urlBuilder("redis")
	if err != nil {
		return nil, err
	}

	options := &redis.Options{
		Addr:     redisConnURL,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbNumber,
	}

	return redis.NewClient(options), nil
}

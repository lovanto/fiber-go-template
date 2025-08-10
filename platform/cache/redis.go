package cache

import (
	"os"
	"strconv"

	"github.com/create-go-app/fiber-go-template/pkg/utils/connection_url_builder"
	"github.com/redis/go-redis/v9"
)

type URLBuilder func(service string) (string, error)

var DefaultURLBuilder = connection_url_builder.ConnectionURLBuilder

func RedisConnection() (*redis.Client, error) {
	return NewRedisConnection(DefaultURLBuilder)
}

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

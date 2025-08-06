package cache

import (
	"os"
	"strconv"

	"github.com/create-go-app/fiber-go-template/pkg/utils"

	"github.com/redis/go-redis/v9"
)

func RedisConnection() (*redis.Client, error) {
	dbNumber, _ := strconv.Atoi(os.Getenv("REDIS_DB_NUMBER"))
	redisConnURL, err := utils.ConnectionURLBuilder("redis")
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

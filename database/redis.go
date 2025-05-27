package database

import (
	"github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Use o nome do serviço Docker ("redis")
		Password: "",
		DB:       0,
	})
}

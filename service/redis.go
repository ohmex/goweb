package service

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	*redis.Client
	ctx context.Context
}

func NewRedisClient() *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &RedisClient{client, context.Background()}
}

package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisClient struct {
	*redis.Client
}

func NewRedisClient() *redisClient {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return &redisClient{client}
}

func (c *redisClient) Get(key string) (string, error) {
	return c.Client.Get(context.Background(), key).Result()
}

func (c *redisClient) Set(key string, value interface{}, expiration time.Duration) error {
	return c.Client.Set(context.Background(), key, value, time.Second*10).Err()
}

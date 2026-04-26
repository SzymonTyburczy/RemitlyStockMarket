package repository

import "github.com/redis/go-redis/v9"

// NewRedisClient creates and returns a configured Redis client.
func NewRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{Addr: addr})
}

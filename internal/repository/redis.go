package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisRepo struct {
	client *redis.Client
	ctx    context.Context
}

var RedisRepoInstance *RedisRepo

// NewRedisRepo initializes a new RedisRepo.
func NewRedisRepo(addr string) *RedisRepo {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &RedisRepo{
		client: client,
		ctx:    context.Background(),
	}
}

// CacheURL caches the short code and original URL in Redis.
func (r *RedisRepo) CacheURL(shortCode string, originalURL string, expiration time.Duration) error {
	return r.client.Set(r.ctx, shortCode, originalURL, expiration).Err()
}

// GetCachedURL retrieves the original URL from Redis cache.
func (r *RedisRepo) GetCachedURL(shortCode string) (string, error) {
	val, err := r.client.Get(r.ctx, shortCode).Result()
	if err == redis.Nil {
		return "", nil // Key not found
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

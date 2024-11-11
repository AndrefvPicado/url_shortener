package repository

import "time"

// interface for the URL repository
type PostgresRepository interface {
	SaveURL(shortCode string, originalURL string) error
	GetOriginalURL(shortCode string) (string, error)
}

// interface for the Redis repository
type RedisRepository interface {
	CacheURL(shortCode string, originalURL string, expiration time.Duration) error
	GetCachedURL(shortCode string) (string, error)
}

package services

import (
	"context"
	"errors"
	"time"
	"url-shortener-service/internal/repository"
	"url-shortener-service/internal/utils"
)

const (
	cacheExpiration = 24 * time.Hour
	timeout         = 5 * time.Second
)

type URLProcessResult struct {
	URL string
	Err error
}

func ShortenURL(originalURL string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resultChan := make(chan URLProcessResult, 1)

	go func() {
		shortCode := utils.GenerateShortCode(originalURL)

		// Create error channel for database operations
		dbErrChan := make(chan error, 2)

		// Concurrent database operations
		go func() {
			dbErrChan <- repository.PostgresRepoInstance.SaveURL(shortCode, originalURL)
		}()

		go func() {
			repository.RedisRepoInstance.CacheURL(shortCode, originalURL, cacheExpiration)
			dbErrChan <- nil
		}()

		// Wait for both operations to complete
		for i := 0; i < 2; i++ {
			if err := <-dbErrChan; err != nil {
				resultChan <- URLProcessResult{"", err}
				return
			}
		}

		shortURL := "http://localhost:8080/redirect/" + shortCode
		resultChan <- URLProcessResult{shortURL, nil}
	}()

	// Wait for result or timeout
	select {
	case result := <-resultChan:
		return result.URL, result.Err
	case <-ctx.Done():
		return "", errors.New("operation timed out")
	}
}

func FetchOriginalURL(shortCode string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resultChan := make(chan URLProcessResult, 1)

	go func() {
		// Try Redis first
		redisChan := make(chan URLProcessResult, 1)
		go func() {
			originalURL, err := repository.RedisRepoInstance.GetCachedURL(shortCode)
			redisChan <- URLProcessResult{originalURL, err}
		}()

		// Wait for Redis result or proceed to Postgres
		select {
		case redisResult := <-redisChan:
			if redisResult.Err == nil && redisResult.URL != "" {
				resultChan <- redisResult
				return
			}
		case <-time.After(100 * time.Millisecond):
			// Redis timeout, continue to Postgres
		}

		// Fallback to PostgreSQL
		originalURL, err := repository.PostgresRepoInstance.GetOriginalURL(shortCode)
		if err != nil || originalURL == "" {
			resultChan <- URLProcessResult{"", err}
			return
		}

		// Cache the result asynchronously
		go func() {
			repository.RedisRepoInstance.CacheURL(shortCode, originalURL, cacheExpiration)
		}()

		resultChan <- URLProcessResult{originalURL, nil}
	}()

	// Wait for result or timeout
	select {
	case result := <-resultChan:
		return result.URL, result.Err
	case <-ctx.Done():
		return "", errors.New("operation timed out")
	}
}

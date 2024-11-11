package services

import (
	"time"
	"url-shortener-service/internal/repository"
	"url-shortener-service/internal/utils"
)

const cacheExpiration = 24 * time.Hour

func ShortenURL(originalURL string) (string, error) {
	// generate a short code
	shortCode := utils.GenerateShortCode(originalURL)

	//save to postgres
	if err := repository.PostgresRepoInstance.SaveURL(shortCode, originalURL); err != nil {
		return "", err
	}

	// Cache the result in Redis
	repository.RedisRepoInstance.CacheURL(shortCode, originalURL, cacheExpiration)

	// construct the short URL
	shortURL := "http://localhost:8080/redirect/" + shortCode

	return shortURL, nil
}

func FetchOriginalURL(shortCode string) (string, error) {
	// Check Redis cache first
	originalURL, err := repository.RedisRepoInstance.GetCachedURL(shortCode)
	if err == nil && originalURL != "" {
		return originalURL, nil
	}

	// Fallback to PostgreSQL if not found in Redis
	originalURL, err = repository.PostgresRepoInstance.GetOriginalURL(shortCode)
	if err != nil || originalURL == "" {
		return "", err
	}

	// Cache the URL for future requests
	repository.RedisRepoInstance.CacheURL(shortCode, originalURL, cacheExpiration)
	return originalURL, nil
}

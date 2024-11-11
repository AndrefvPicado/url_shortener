package services

import (
	"errors"
	"url-shortener-service/internal/repository"
	"url-shortener-service/internal/utils"
)

func ShortenURL(originalURL string) (string, error) {
	// generate a short code
	shortCode := utils.GenerateShortCode(originalURL)

	// save the short code and original URL to the database
	if err := repository.SaveURL(shortCode, originalURL); err != nil {
		return "", err
	}

	// construct the short URL
	shortURL := "http://localhost:8080/redirect/" + shortCode

	return shortURL, nil
}

func FetchOriginalURL(shortCode string) (string, error) {
	// fetch the original URL from the database
	originalURL, err := repository.GetOriginalURL(shortCode)
	if err != nil {
		return "", errors.New("original URL not found")
	}

	return originalURL, nil
}

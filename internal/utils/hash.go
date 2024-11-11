package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"strings"
)

// GenerateShortCode generates a unique short code based on the URL.
func GenerateShortCode(originalURL string) string {
	// Create a SHA-256 hash of the original URL
	hash := sha256.Sum256([]byte(originalURL))

	// Encode the first few bytes to base64 and trim to get a shorter code
	shortCode := base64.URLEncoding.EncodeToString(hash[:6])

	// Remove any non-alphanumeric characters (optional)
	return strings.TrimRight(shortCode, "=")
}

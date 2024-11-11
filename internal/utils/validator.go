package utils

import (
	"net/url"
)

// IsValidURL checks if the given URL is valid.
func IsValidURL(input string) bool {
	parsedURL, err := url.ParseRequestURI(input)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	}
	return true
}

package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"url-shortener-service/internal/models"
	"url-shortener-service/internal/services"
)

// CreateShortURLHandler handles URL shortening requests.
func CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// parse the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	var request models.URLRequest
	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, "Invalid request body, non JSON format", http.StatusBadRequest)
		return
	}

	// call the service to create the short URL
	shortURL, err := services.ShortenURL(request.OriginalURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create short URL: %v", err), http.StatusInternalServerError)
		return
	}

	// respond with the short URL
	response := models.URLResponse{ShortURL: shortURL}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	// log the event
	fmt.Println("Short URL created at", time.Now().Format(time.RFC3339))
}

// RedirectHandler handles redirection based on the short code.
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	shortCode := r.URL.Path[len("/redirect/"):]

	// call the service to get the long URL
	if shortCode == "" {
		http.Error(w, "Short code is required", http.StatusBadRequest)
		return
	}

	// call the service to get the original URL
	originalURL, err := services.FetchOriginalURL(shortCode)
	if err != nil {
		http.Error(w, "Short code not found", http.StatusNotFound)
		return
	}

	// redirect to the original URL
	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)

	// log the event
	fmt.Println("Redirecting to original URL")
}

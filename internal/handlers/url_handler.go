package handlers

import (
	"fmt"
	"net/http"
)

func CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: Implement short URL creation logic
	// parse input, call service, and return a short URL
	fmt.Println("Short URL created")
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement redirect logic
	// Extract the short code from the URL path and redirect
	fmt.Println("Redirecting to original URL")
}

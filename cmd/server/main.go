package main

import (
	"log"
	"net/http"
	"url-shortener-service/config"
	"url-shortener-service/internal/handlers"
)

func main() {
	// Load configuration
	config := config.LoadConfig()

	// Initialize http server and routes
	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", handlers.CreateShortURLHandler)
	mux.HandleFunc("/redirect/", handlers.RedirectHandler)

	// Start server
	serverAddress := ":" + config.ServerPort
	log.Printf("Server starting on %s\n", serverAddress)
	if err := http.ListenAndServe(serverAddress, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

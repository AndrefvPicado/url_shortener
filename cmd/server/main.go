package main

import (
	"fmt"
	"log"
	"net/http"
	"url-shortener-service/config"
	"url-shortener-service/internal/handlers"
	"url-shortener-service/internal/repository"
)

func main() {
	// Load configuration
	config := config.LoadConfig()

	// Initialize http server and routes
	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", handlers.CreateShortURLHandler)
	mux.HandleFunc("/redirect/", handlers.RedirectHandler)

	// Initialize dependencies
	postgresRepo, err := repository.NewPostgresRepo(
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			config.DBUser,
			config.DBPass,
			config.DBHost,
			config.DBPort,
			config.DBName,
		))
	if err != nil {
		log.Fatalf("Failed to initialize PostgreSQL: %v", err)
	}
	repository.PostgresRepoInstance = postgresRepo

	repository.RedisRepoInstance = repository.NewRedisRepo("localhost:6379")

	// Start server
	serverAddress := ":" + config.ServerPort
	log.Printf("Server starting on %s\n", serverAddress)
	if err := http.ListenAndServe(serverAddress, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

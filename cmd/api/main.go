package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go_api/internal/app/route"
	serverconfig "go_api/internal/config"
	"go_api/internal/storage"
)

func main() {
	// Load config
	config, err := serverconfig.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Connect to database
	if err := storage.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer storage.Close()

	// Run migrations
	log.Println("Running database migrations...")
	if err := storage.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations completed successfully")

	// Connect to Redis
	redisClient := storage.ConnectRedis()
	if redisClient == nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Set up HTTP server
	mux := http.NewServeMux()

	// Server instance
	serverAddr := fmt.Sprintf(":%s", config.ServerPort)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: route.SetupRoutes(mux, storage.GetDB(), redisClient),
	}

	// Setup graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("Shutting down server...")
		if err := server.Close(); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
		storage.Close()
		os.Exit(0)
	}()

	fmt.Printf("Listening on port %s\n", config.ServerPort)
	fmt.Printf("Serving on http://localhost%s\n", serverAddr)

	// Run server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}

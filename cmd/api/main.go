package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	serverconfig "go_api/internal/config"
	"go_api/internal/database"
	"go_api/internal/handlers"
	"go_api/internal/routes"
)

func main() {
	// Load config
	config, err := serverconfig.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Connect to database
	if err := database.Connect(config); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run migrations
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Set up HTTP server
	mux := http.NewServeMux()

	// Setup handler
	handler := handlers.NewHandler(database.GetDB())

	// Setup routes
	routes.SetupRoutes(mux, handler)

	// Server instance
	serverAddr := fmt.Sprintf(":%s", config.ServerPort)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: mux,
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
		database.Close()
		os.Exit(0)
	}()

	fmt.Printf("Listening on port %s\n", config.ServerPort)
	fmt.Printf("Serving on http://localhost%s\n", serverAddr)

	// Run server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}

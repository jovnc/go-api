package main

import (
	"fmt"
	serverconfig "go_api/configs"
	"go_api/internal/handlers"
	"go_api/internal/routes"
	"log"
	"net/http"
)

func main() {

	// Load config
	config, err := serverconfig.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Set up HTTP server
	mux := http.NewServeMux()

	// Setup handler
	handler := handlers.NewHandler()

	// Setup routes
	routes.SetupRoutes(mux, handler)

	// Server instance
	serverAddr := fmt.Sprintf(":%s", config.ServerPort)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: mux,
	}

	fmt.Printf("Listening on port %s\n", config.ServerPort)
	fmt.Printf("Serving on http://localhost%s\n", serverAddr)

	// Run server
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

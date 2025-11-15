package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go_api/internal/app/handler"
	"go_api/internal/app/route"
	serverconfig "go_api/internal/config"
	"go_api/internal/database"
)

func main() {
	// Parse command-line flags
	migrateOnly := flag.Bool("migrate-only", false, "Run database migrations and exit")
	flag.Parse()

	// Load config
	config, err := serverconfig.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run migrations (if flag is set)
	if *migrateOnly {
		log.Println("Running database migrations...")
		if err := database.Migrate(); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		log.Println("Migrations completed successfully")
		return
	}

	// Connect to Redis
	redisClient := database.ConnectRedis()
	if redisClient == nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Set up HTTP server
	mux := http.NewServeMux()

	// Setup handler
	handler := handler.NewHandler(database.GetDB(), redisClient)

	// Server instance
	serverAddr := fmt.Sprintf(":%s", config.ServerPort)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: route.SetupRoutes(mux, handler),
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

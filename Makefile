.PHONY: run build migrate clean deps fmt stop help

APP_NAME=go_api
BINARY_NAME=gobin
BUILD_DIR=./bin
GO_FILES=$(shell find . -name "*.go" -not -path "./vendor/*")

run:
	@echo "Running $(APP_NAME)..."
	@go run ./cmd/api || true

build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/api

migrate:
	@echo "Migrating database..."
	@go run ./cmd/api -migrate-only || go run -tags migrate ./cmd/api

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

fmt:
	@echo "Formatting code..."
	@go fmt ./...

stop:
	@echo "Stopping $(APP_NAME)..."
	@pkill -f "go run ./cmd/api" || echo "No $(APP_NAME) process found"

help:
	@echo "Usage: make <target>"
	@echo "Targets:"
	@echo "  run - Run the $(APP_NAME)"
	@echo "  build - Build the $(APP_NAME)"
	@echo "  migrate - Run migrations for the database"
	@echo "  clean - Clean up the build directory"
	@echo "  deps - Download dependencies"
	@echo "  fmt - Format the code"
	@echo "  stop - Stop the $(APP_NAME) process"
	@echo "  help - Show this help message"

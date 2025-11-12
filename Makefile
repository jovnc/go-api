.PHONY: dev start build migrate clean deps fmt stop

APP_NAME=go_api
BINARY_NAME=gobin
BUILD_DIR=./bin
GO_FILES=$(shell find . -name "*.go" -not -path "./vendor/*")

build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/api

dev:
	@echo "Running $(APP_NAME) in development mode..."
	@air

start:build
	@echo "Starting $(APP_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

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

.PHONY: dev start build clean deps fmt stop test

APP_NAME=go_api
BINARY_NAME=go_api
BUILD_DIR=./build
GO_FILES=$(shell find . -name "*.go" -not -path "./vendor/*")

build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/api

dev:
	@echo "Running $(APP_NAME) in development mode..."
	@go tool air

start:build
	@echo "Starting $(APP_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

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

test:
	@echo "Running tests..."
	@go test ./tests/... -v

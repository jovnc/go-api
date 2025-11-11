.PHONY: run build migrate clean

# Default binary name
BINARY_NAME=go_api

run:
	go run ./cmd/api

build:
	go build -o bin/$(BINARY_NAME) ./cmd/api

migrate:
	go run ./cmd/api -migrate-only || go run -tags migrate ./cmd/api

clean:
	rm -rf bin/

deps:
	go mod download
	go mod tidy

.PHONY: build run clean test docker-build docker-run install dev

# Build the application
build:
	@echo "Building NotifyPipe..."
	@go build -o bin/notifypipe ./cmd/notifypipe

# Run the application
run:
	@echo "Running NotifyPipe..."
	@go run ./cmd/notifypipe

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf data/
	@rm -f notifypipe

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Install dependencies
install:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t ghcr.io/fatlirmorina/notifypipe:latest .

# Run with Docker Compose
docker-run:
	@echo "Starting NotifyPipe with Docker Compose..."
	@docker-compose up -d

# Stop Docker Compose
docker-stop:
	@echo "Stopping NotifyPipe..."
	@docker-compose down

# Development mode with hot reload
dev:
	@echo "Starting development server..."
	@air || go run ./cmd/notifypipe

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@golangci-lint run || go vet ./...

# Generate Go files
generate:
	@echo "Generating Go files..."
	@go generate ./...

# Show help
help:
	@echo "NotifyPipe Makefile Commands:"
	@echo "  make build       - Build the application"
	@echo "  make run         - Run the application"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make test        - Run tests"
	@echo "  make install     - Install dependencies"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run  - Run with Docker Compose"
	@echo "  make docker-stop - Stop Docker Compose"
	@echo "  make dev         - Run in development mode"
	@echo "  make fmt         - Format code"
	@echo "  make lint        - Lint code"

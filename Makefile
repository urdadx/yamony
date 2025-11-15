# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	
	
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go &
	@pnpm install --prefer-offline --no-fund --prefix ./frontend
	@pnpm run dev --prefix ./frontend
# Create DB container
docker-run:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v
# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                $(shell go env GOPATH)/bin/air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

# Database Migrations
migrate-up:
	@echo "Running migrations..."
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING="postgres://shinobi:password@localhost:5432/yamonydb?sslmode=disable&search_path=public" ~/go/bin/goose -dir internal/database/schema up

migrate-down:
	@echo "Rolling back migrations..."
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING="postgres://shinobi:password@localhost:5432/yamonydb?sslmode=disable&search_path=public" ~/go/bin/goose -dir internal/database/schema down

migrate-status:
	@echo "Migration status..."
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING="postgres://shinobi:password@localhost:5432/yamonydb?sslmode=disable&search_path=public" ~/go/bin/goose -dir internal/database/schema status

.PHONY: all build run test clean watch docker-run docker-down itest migrate-up migrate-down migrate-status

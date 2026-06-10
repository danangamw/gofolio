# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	@CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/go-cms ./cmd/app

# Run the application (with hot-reload via air)
run:
	@if command -v air > /dev/null; then air; else go run ./cmd/app; fi
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
	@rm -rf bin/

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

.PHONY: all build run test clean watch docker-run docker-down itest seed migrate-diff migrate-apply migrate-status migrate-lint

# Seed admin user
seed:
	@go run ./cmd/seed

# ─── DATABASE MIGRATIONS (ATLAS) ─────────────────────────────────────────────
# Atlas reads GORM models via cmd/migrate/loader.go and generates versioned SQL
# files into the migrations/ directory.

# Generate a new migration file from GORM model changes.
# Usage: make migrate-diff name=add_tags_to_blogs
migrate-diff:
	@atlas migrate diff $(name) --env local

# Apply pending migrations to the local database.
migrate-apply:
	@atlas migrate apply --env local --url "$(shell grep DATABASE_URL .env | cut -d= -f2-)"

# Show current migration status.
migrate-status:
	@atlas migrate status --env local --url "$(shell grep DATABASE_URL .env | cut -d= -f2-)"

# Lint migration files for safety issues.
migrate-lint:
	@atlas migrate lint --env local --git-base main

# Apply migrations to PRODUCTION (requires DATABASE_URL env var).
migrate-prod:
	@atlas migrate apply \
	  --dir "file://migrations" \
	  --url "$(DATABASE_URL)"

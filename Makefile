.PHONY: help build run dev test clean migrate-up migrate-down docker-up docker-down seed

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building application..."
	@go build -o bin/api cmd/api/main.go

run: ## Run the application
	@echo "Running application..."
	@go run cmd/api/main.go

dev: ## Run with air (hot reload)
	@echo "Running with hot reload..."
	@air

test: ## Run tests
	@echo "Running tests..."
	@go test -v -cover ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build files
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

deps: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

docker-up: ## Start docker services
	@echo "Starting docker services..."
	@docker-compose up -d

docker-down: ## Stop docker services
	@echo "Stopping docker services..."
	@docker-compose down

docker-logs: ## View docker logs
	@docker-compose logs -f

docker-local-up: ## Start local development services (PostgreSQL + Redis)
	@echo "Starting local development services..."
	@docker-compose -f docker-compose.local.yml up -d
	@echo "Waiting for services to be ready..."
	@timeout /t 5 /nobreak > nul
	@docker-compose -f docker-compose.local.yml ps

docker-local-down: ## Stop local development services
	@echo "Stopping local development services..."
	@docker-compose -f docker-compose.local.yml down

docker-local-logs: ## View local development logs
	@docker-compose -f docker-compose.local.yml logs -f

docker-local-restart: ## Restart local development services
	@docker-compose -f docker-compose.local.yml restart

migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	@go run scripts/migrate.go up

migrate-down: ## Run database migrations down
	@echo "Running migrations down..."
	@go run scripts/migrate.go down

migrate-create: ## Create a new migration (usage: make migrate-create name=create_users)
	@echo "Creating migration..."
	@go run scripts/migrate.go create $(name)

seed: ## Seed the database
	@echo "Seeding database..."
	@go run scripts/seed.go

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger docs..."
	@swag init -g cmd/api/main.go -o docs

install-tools: ## Install development tools
	@echo "Installing tools..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest

.DEFAULT_GOAL := help

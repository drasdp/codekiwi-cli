# ============================================
# CodeKiwi Development Makefile
# ============================================
#
# This Makefile provides convenient commands for developing CodeKiwi.
# It handles building the CLI, Docker images, and managing dev environments.
#
# Quick Start:
#   make help          # Show all available commands
#   make dev-setup     # Initial setup for development
#   make dev-start     # Start development environment
#   make cli-build     # Build the CLI binary
#   make test          # Run tests
# ============================================

.PHONY: help
help: ## Show this help message
	@echo "CodeKiwi Development Commands"
	@echo "=============================="
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ============================================
# CLI Development
# ============================================

.PHONY: cli-build
cli-build: ## Build the CLI binary
	@echo "Building CLI..."
	cd cli-go && go build -o codekiwi cmd/codekiwi/main.go
	@echo "CLI binary built at: cli-go/codekiwi"

.PHONY: cli-install
cli-install: cli-build ## Build and install CLI to /usr/local/bin
	@echo "Installing CLI to /usr/local/bin..."
	sudo cp cli-go/codekiwi /usr/local/bin/codekiwi
	@echo "CLI installed successfully!"
	@codekiwi --version

.PHONY: cli-test
cli-test: cli-build ## Build CLI and run test commands
	@echo "Testing CLI..."
	cd cli-go && ./codekiwi --help
	cd cli-go && ./codekiwi --version

.PHONY: cli-run
cli-run: cli-build ## Build and run CLI with 'start' command
	@echo "Running CLI in development mode..."
	cd cli-go && ./codekiwi start

# ============================================
# Docker Runtime Development
# ============================================

.PHONY: runtime-build
runtime-build: ## Build the runtime Docker image locally
	@echo "Building runtime image..."
	docker build -t drasdp/codekiwi-runtime:dev ./runtime
	@echo "Runtime image built: drasdp/codekiwi-runtime:dev"

.PHONY: runtime-test
runtime-test: runtime-build ## Build and test the runtime image
	@echo "Testing runtime image..."
	docker run --rm drasdp/codekiwi-runtime:dev nginx -t
	@echo "Runtime image test passed!"

# ============================================
# Development Environment (docker-compose.dev.yaml)
# ============================================

.PHONY: dev-start
dev-start: ## Start development environment with docker-compose.dev.yaml
	@echo "Starting development environment..."
	docker compose -f docker-compose.dev.yaml up -d --build
	@echo "Development environment started!"
	@echo "Web UI: http://localhost:8080"

.PHONY: dev-logs
dev-logs: ## View logs from development environment
	docker compose -f docker-compose.dev.yaml logs -f

.PHONY: dev-stop
dev-stop: ## Stop development environment
	@echo "Stopping development environment..."
	docker compose -f docker-compose.dev.yaml down
	@echo "Development environment stopped!"

.PHONY: dev-restart
dev-restart: dev-stop dev-start ## Restart development environment

.PHONY: dev-clean
dev-clean: dev-stop ## Stop and clean development environment (remove volumes)
	@echo "Cleaning development environment..."
	docker compose -f docker-compose.dev.yaml down -v
	@echo "Development environment cleaned!"

# ============================================
# Full Development Workflow
# ============================================

.PHONY: dev-setup
dev-setup: ## Initial setup for development (install dependencies)
	@echo "Setting up development environment..."
	@echo "1. Checking Go installation..."
	@go version || (echo "Go is not installed. Please install Go 1.20+"; exit 1)
	@echo "2. Checking Docker installation..."
	@docker --version || (echo "Docker is not installed. Please install Docker"; exit 1)
	@echo "3. Installing Go dependencies..."
	cd cli-go && go mod download
	@echo "4. Building CLI..."
	@$(MAKE) cli-build
	@echo ""
	@echo "Development setup complete!"
	@echo "Next steps:"
	@echo "  make cli-run       # Test CLI locally"
	@echo "  make dev-start     # Start dev environment"

.PHONY: dev-all
dev-all: cli-build runtime-build ## Build both CLI and runtime image

.PHONY: test
test: cli-test runtime-test ## Run all tests

# ============================================
# Testing with Real Project
# ============================================

.PHONY: test-project-create
test-project-create: ## Create a test project directory
	@echo "Creating test project..."
	mkdir -p test-project
	@echo "Test project created at: test-project/"

.PHONY: test-project-start
test-project-start: cli-build ## Build CLI and start test project
	@echo "Starting test project..."
	cd cli-go && ./codekiwi start ../test-project

.PHONY: test-project-clean
test-project-clean: ## Clean test project
	@echo "Cleaning test project..."
	rm -rf test-project
	@echo "Test project cleaned!"

# ============================================
# Git and Release
# ============================================

.PHONY: git-status
git-status: ## Show git status
	@git status

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -f cli-go/codekiwi
	@echo "Clean complete!"

.PHONY: clean-all
clean-all: clean dev-clean test-project-clean ## Clean everything

# ============================================
# Default target
# ============================================

.DEFAULT_GOAL := help

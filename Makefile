# ============================================
# CodeKiwi Development Makefile
# ============================================
# KEEP IT SIMPLE STUPID
#
# Build:
#   make build dev-cli      # Build CLI binary
#   make build dev-runtime  # Build runtime Docker image (dev tag)
#   make build dev          # Build both
#
# Run (like production 'codekiwi' command):
#   make dev start [path]   # Start CodeKiwi with dev CLI + dev image
#   make dev list           # List instances
#   make dev kill [path]    # Kill instance
#
# ============================================

CLI_BINARY := cli-go/codekiwi
RUNTIME_IMAGE := drasdp/codekiwi-runtime:dev

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show available commands
	@echo "CodeKiwi Development - KISS"
	@echo "==========================="
	@echo ""
	@echo "Build commands:"
	@echo "  make build dev-cli      - Build CLI binary"
	@echo "  make build dev-runtime  - Build runtime image (dev tag)"
	@echo "  make build dev          - Build both CLI + runtime"
	@echo ""
	@echo "Dev commands (like 'codekiwi' in production):"
	@echo "  make dev start [path]   - Start CodeKiwi"
	@echo "  make dev list           - List instances"
	@echo "  make dev kill [path]    - Kill instance"
	@echo ""
	@echo "Other:"
	@echo "  make clean              - Clean build artifacts"

# ============================================
# Build Commands
# ============================================

.PHONY: build
build:
	@echo "Usage: make build [dev-cli|dev-runtime|dev]"

build: dev-cli
.PHONY: dev-cli
dev-cli: ## Build CLI binary
	@echo "Building CLI binary..."
	cd cli-go && go build -o codekiwi cmd/codekiwi/main.go
	@echo "✓ CLI binary: cli-go/codekiwi"

build: dev-runtime
.PHONY: dev-runtime
dev-runtime: ## Build runtime Docker image (dev tag)
	@echo "Building runtime image (dev tag)..."
	docker build -t $(RUNTIME_IMAGE) ./runtime
	@echo "✓ Runtime image: $(RUNTIME_IMAGE)"

build: dev
.PHONY: dev
dev: dev-cli dev-runtime ## Build both CLI and runtime
	@echo "✓ Dev build complete"

# ============================================
# Dev Commands (wrapper for cli-go/codekiwi)
# ============================================

.PHONY: check-cli
check-cli:
	@if [ ! -f $(CLI_BINARY) ]; then \
		echo "Error: CLI binary not found. Run 'make build dev-cli' first."; \
		exit 1; \
	fi

.PHONY: check-runtime
check-runtime:
	@if ! docker images -q $(RUNTIME_IMAGE) | grep -q .; then \
		echo "Error: Runtime image not found. Run 'make build dev-runtime' first."; \
		exit 1; \
	fi

.PHONY: dev
dev: check-cli check-runtime
	@echo "Usage: make dev [start|list|kill] [args]"

dev: start
.PHONY: start
start: check-cli check-runtime ## Start CodeKiwi (usage: make dev start [path])
	@cd cli-go && ./codekiwi start $(filter-out $@,$(MAKECMDGOALS))

dev: list
.PHONY: list
list: check-cli ## List instances
	@cd cli-go && ./codekiwi list

dev: kill
.PHONY: kill
kill: check-cli ## Kill instance (usage: make dev kill [path])
	@cd cli-go && ./codekiwi kill $(filter-out $@,$(MAKECMDGOALS))

# ============================================
# Cleanup
# ============================================

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f $(CLI_BINARY)
	@echo "✓ Clean complete"

# Allow passing arguments to targets
%:
	@:

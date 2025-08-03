.PHONY: test test-race test-coverage test-all lint fmt clean examples help

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Testing targets
test: ## Run basic tests
	go test -v ./...

test-race: ## Run tests with race detection
	go test -race -v ./...

test-coverage: ## Run tests with coverage
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html

test-all: test-race test-coverage ## Run all tests

# Code quality
lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...
	go mod tidy

# Benchmarks
bench: ## Run benchmarks
	go test -bench=. -benchtime=1s ./...

# Examples (for development/testing only)
examples: ## Test example applications (no build artifacts)
	cd examples/basic && go mod tidy && go run main.go &
	sleep 2 && pkill -f "go run main.go" || true
	cd examples/multi-config && go mod tidy && go run main.go &
	sleep 2 && pkill -f "go run main.go" || true

# Cleanup
clean: ## Clean up generated files
	rm -f coverage.out coverage.html
	find examples -name "*.json" -delete
	go clean -cache

# Development
dev: fmt lint test ## Run development checks (format, lint, test)

# CI simulation
ci: fmt lint test-all bench ## Run CI checks locally

# Release preparation
release-check: ci ## Final checks before release
	@echo "✅ Library ready for release!"

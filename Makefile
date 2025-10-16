# ZipCodeReader Makefile
# Build, test, run, and clean tasks for the project

.PHONY: help build run test clean clean-db clean-all dev install lint fmt check deps

# Default target
.DEFAULT_GOAL := help

# Variables
BINARY_NAME=zipcodereader
DB_FILE=zipcodereader.db
TEST_DB_FILE=zipcodereader_test.db
GO=go
GOFLAGS=-v
PORT?=8080

## help: Display this help message
help:
	@echo "ZipCodeReader - Available Make Targets:"
	@echo ""
	@echo "  make build          - Build the application binary"
	@echo "  make run            - Build and run the application (local auth)"
	@echo "  make run-oauth      - Build and run with OAuth2 authentication"
	@echo "  make dev            - Run in development mode with live reload"
	@echo "  make test           - Run all unit tests"
	@echo "  make test-verbose   - Run all tests with verbose output"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make clean          - Remove binary and temporary files"
	@echo "  make clean-db       - Remove database files"
	@echo "  make clean-all      - Remove binary, database, and all generated files"
	@echo "  make install        - Install Go dependencies"
	@echo "  make fmt            - Format Go code"
	@echo "  make lint           - Run linters (requires golangci-lint)"
	@echo "  make check          - Run fmt, lint, and tests"
	@echo "  make deps           - Download and verify dependencies"
	@echo ""

## build: Build the application binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) .
	@echo "✅ Build complete: ./$(BINARY_NAME)"

## run: Build and run the application with local authentication
run: build
	@echo "Starting $(BINARY_NAME) with local authentication on port $(PORT)..."
	@echo "Visit http://localhost:$(PORT)"
	@echo "Press Ctrl+C to stop"
	./$(BINARY_NAME)

## run-oauth: Build and run with OAuth2 authentication
run-oauth: build
	@echo "Starting $(BINARY_NAME) with OAuth2 authentication on port $(PORT)..."
	@echo "Visit http://localhost:$(PORT)"
	@echo "Press Ctrl+C to stop"
	./$(BINARY_NAME) --use_oauth2

## dev: Run in development mode (rebuilds on change - requires air or entr)
dev:
	@if command -v air > /dev/null; then \
		echo "Running with air (hot reload)..."; \
		air; \
	elif command -v entr > /dev/null; then \
		echo "Running with entr (auto-rebuild on file change)..."; \
		find . -name '*.go' | entr -r make run; \
	else \
		echo "🔄 Using built-in watch script for hot reload..."; \
		./scripts/watch.sh; \
	fi

## test: Run all unit tests
test:
	@echo "Running unit tests..."
	$(GO) test ./... -short
	@echo "✅ All tests passed"

## test-verbose: Run all tests with verbose output
test-verbose:
	@echo "Running unit tests (verbose)..."
	$(GO) test ./... -v -short

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test ./... -coverprofile=coverage.out
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

## test-handlers: Run only handler tests
test-handlers:
	@echo "Running handler tests..."
	$(GO) test ./handlers -v

## test-services: Run only service tests
test-services:
	@echo "Running service tests..."
	$(GO) test ./services -v

## test-models: Run only model tests
test-models:
	@echo "Running model tests..."
	$(GO) test ./models -v

## clean: Remove binary and temporary files
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html
	@rm -f /tmp/cookies.txt /tmp/inst_cookies.txt
	@echo "✅ Clean complete"

## clean-db: Remove database files
clean-db:
	@echo "Removing database files..."
	@rm -f $(DB_FILE) $(TEST_DB_FILE)
	@echo "✅ Database files removed"

## clean-all: Remove binary, database, and all generated files
clean-all: clean clean-db
	@echo "Removing all generated files..."
	@rm -rf tmp/
	@find . -name '*.backup' -type f -delete
	@echo "✅ Full clean complete"

## install: Install Go dependencies
install:
	@echo "Installing dependencies..."
	$(GO) mod download
	$(GO) mod verify
	@echo "✅ Dependencies installed"

## deps: Download and verify dependencies
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	@echo "Verifying dependencies..."
	$(GO) mod verify
	@echo "Tidying go.mod..."
	$(GO) mod tidy
	@echo "✅ Dependencies ready"

## fmt: Format Go code
fmt:
	@echo "Formatting Go code..."
	$(GO) fmt ./...
	@echo "✅ Code formatted"

## lint: Run linters (requires golangci-lint)
lint:
	@if command -v golangci-lint > /dev/null; then \
		echo "Running golangci-lint..."; \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install with:"; \
		echo "  brew install golangci-lint"; \
		echo "  or"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi

## check: Run fmt, lint, and tests (pre-commit check)
check: fmt test
	@echo "✅ All checks passed - ready to commit!"

## migrate: Run database migrations (create fresh DB)
migrate: clean-db build
	@echo "Creating fresh database with migrations..."
	@./$(BINARY_NAME) &
	@sleep 2
	@pkill -f $(BINARY_NAME)
	@echo "✅ Database migrated"

## seed: Seed database with test data
seed: migrate
	@echo "Seeding database with test data..."
	@./$(BINARY_NAME) &
	@sleep 2
	@curl -s -X POST http://localhost:8080/local/register \
		-H "Content-Type: application/x-www-form-urlencoded" \
		-d "username=instructor1&email=instructor@test.com&password=password123&password_confirm=password123&role=instructor" > /dev/null
	@curl -s -X POST http://localhost:8080/local/register \
		-H "Content-Type: application/x-www-form-urlencoded" \
		-d "username=student1&email=student1@test.com&password=password123&password_confirm=password123&role=student" > /dev/null
	@curl -s -X POST http://localhost:8080/local/register \
		-H "Content-Type: application/x-www-form-urlencoded" \
		-d "username=student2&email=student2@test.com&password=password123&password_confirm=password123&role=student" > /dev/null
	@pkill -f $(BINARY_NAME)
	@echo "✅ Database seeded with test users:"
	@echo "   - instructor1:password123 (instructor)"
	@echo "   - student1:password123 (student)"
	@echo "   - student2:password123 (student)"

## backup: Backup database file
backup:
	@if [ -f $(DB_FILE) ]; then \
		BACKUP_FILE="$(DB_FILE).backup.$$(date +%Y%m%d_%H%M%S)"; \
		cp $(DB_FILE) $$BACKUP_FILE; \
		echo "✅ Database backed up to $$BACKUP_FILE"; \
	else \
		echo "❌ No database file to backup"; \
		exit 1; \
	fi

## restore: Restore database from latest backup
restore:
	@LATEST_BACKUP=$$(ls -t $(DB_FILE).backup.* 2>/dev/null | head -1); \
	if [ -n "$$LATEST_BACKUP" ]; then \
		cp $$LATEST_BACKUP $(DB_FILE); \
		echo "✅ Database restored from $$LATEST_BACKUP"; \
	else \
		echo "❌ No backup files found"; \
		exit 1; \
	fi

## info: Show project information
info:
	@echo "Project: ZipCodeReader"
	@echo "Binary: $(BINARY_NAME)"
	@echo "Database: $(DB_FILE)"
	@echo "Go Version: $$(go version)"
	@echo ""
	@echo "Project Statistics:"
	@echo "  Go Files: $$(find . -name '*.go' -not -path './vendor/*' | wc -l)"
	@echo "  Lines of Code: $$(find . -name '*.go' -not -path './vendor/*' -not -name '*_test.go' -not -name '*.backup' | xargs wc -l 2>/dev/null | tail -1 | awk '{print $$1}')"
	@echo "  Test Files: $$(find . -name '*_test.go' | wc -l)"
	@echo "  Handlers: $$(ls handlers/*.go 2>/dev/null | grep -v _test | wc -l)"
	@echo "  Services: $$(ls services/*.go 2>/dev/null | grep -v _test | wc -l)"
	@echo "  Models: $$(ls models/*.go 2>/dev/null | grep -v _test | wc -l)"

## docker-build: Build Docker image (future)
docker-build:
	@echo "Docker support coming soon..."

## docker-run: Run in Docker container (future)
docker-run:
	@echo "Docker support coming soon..."

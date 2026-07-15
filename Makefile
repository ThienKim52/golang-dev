.PHONY: all build run test test-unit test-integration clean swagger fmt lint deps help

# Variables
APP_NAME=health-check
CMD_DIR=./cmd/api
BUILD_DIR=./build
BINARY_NAME=$(BUILD_DIR)/$(APP_NAME)
GO_FILES=$(shell find . -name '*.go' -not -path './vendor/*' -not -path './build/*')

# Default target
all: swagger run

## build: Build the application
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BINARY_NAME) $(CMD_DIR)/main.go
	@echo "Build complete: $(BINARY_NAME)"

## run: Run the application
run:
	@echo "Running $(APP_NAME)..."
	@go run $(CMD_DIR)/main.go

## test: Run all tests
test-all:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

test:
	@echo "Running all tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html


## test-unit: Run unit tests only
test-unit:
	@echo "Running unit tests..."
	@go test -v -race ./internal/... -run "^((?!integration).)*$$"

## test-integration: Run integration tests only
test-integration:
	@echo "Running integration tests..."
	@go test -v ./internal/intergration_test/...

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out
	@go clean
	@echo "Clean complete"

## swagger: Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/api/main.go -o docs
	@echo "Swagger documentation generated"

## fmt: Format Go code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

## lint: Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run ./... || echo "golangci-lint not installed. Install it from https://golangci-lint.run/usage/install/"

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

## deps-verify: Verify dependencies
deps-verify:
	@echo "Verifying dependencies..."
	@go mod verify

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME):latest .

## docker-run: Run Docker container
docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 $(APP_NAME):latest

## help: Show this help message
help:
	@echo "Available targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
COVERAGE_THRESHOLD = 80
test-cover:
	go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1
	grep -vE "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@total=$$(go tool cover -func=coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
	  echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
	  exit 1; \
	else \
	  echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi



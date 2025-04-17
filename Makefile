# Variables
BINARY_NAME=pixellock
BINARY_DIR=bin
DOCKER_IMAGE_NAME=pixellock
GO_FILES=$(wildcard *.go)
VERSION=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X main.Version=${VERSION}"

# Colors for terminal output
GREEN=\033[0;32m
NC=\033[0m # No Color

.PHONY: all build clean test coverage docker-build docker-run fmt lint help install-deps run install release dist

# Default target
all: clean build test

# Build the application
build:
	@printf "$(GREEN)Building $(BINARY_NAME)...$(NC)\n"
	@mkdir -p $(BINARY_DIR)
	@go build $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)
	@printf "$(GREEN)Done! Binary created at $(BINARY_DIR)/$(BINARY_NAME)$(NC)\n"

# Clean build artifacts
clean:
	@printf "$(GREEN)Cleaning build artifacts...$(NC)\n"
	@rm -rf $(BINARY_DIR)
	@go clean
	@printf "$(GREEN)Cleaned!$(NC)\n"

# Run all tests
test:
	@printf "$(GREEN)Running tests...$(NC)\n"
	@go test -v ./... -cover

# Generate test coverage report
coverage:
	@printf "$(GREEN)Generating coverage report...$(NC)\n"
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out
	@rm coverage.out

# Build Docker image
docker-build:
	@printf "$(GREEN)Building Docker image...$(NC)\n"
	@docker build -t $(DOCKER_IMAGE_NAME) .

# Run in Docker container
docker-run: docker-build
	@printf "$(GREEN)Running in Docker container...$(NC)\n"
	@docker run -it --rm $(DOCKER_IMAGE_NAME)

# Format code
fmt:
	@printf "$(GREEN)Formatting code...$(NC)\n"
	@go fmt ./...
	@printf "$(GREEN)Code formatted!$(NC)\n"

# Run linter
lint:
	@printf "$(GREEN)Running linter...$(NC)\n"
	@if command -v golangci-lint >/dev/null; then \
		golangci-lint run ./...; \
	else \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run ./...; \
	fi

# Install project dependencies
install-deps:
	@printf "$(GREEN)Installing dependencies...$(NC)\n"
	@go mod download
	@printf "$(GREEN)Dependencies installed!$(NC)\n"

# Run the application
run: build
	@printf "$(GREEN)Running $(BINARY_NAME)...$(NC)\n"
	@./$(BINARY_DIR)/$(BINARY_NAME)

# Install binary to GOPATH/bin
install: build
	@printf "$(GREEN)Installing $(BINARY_NAME)...$(NC)\n"
	@go install ./...

# Release build (optimized)
release: clean
	@printf "$(GREEN)Building release version...$(NC)\n"
	@mkdir -p $(BINARY_DIR)
	@go build -ldflags="-s -w $(LDFLAGS)" -o $(BINARY_DIR)/$(BINARY_NAME)
	@printf "$(GREEN)Release build created at $(BINARY_DIR)/$(BINARY_NAME)$(NC)\n"

# Create distribution packages
dist: release
	@printf "$(GREEN)Creating distribution packages...$(NC)\n"
	@mkdir -p dist
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w $(LDFLAGS)" -o $(BINARY_DIR)/$(BINARY_NAME)-linux-amd64
	@tar -czf dist/$(BINARY_NAME)-linux-amd64.tar.gz -C $(BINARY_DIR) $(BINARY_NAME)-linux-amd64
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w $(LDFLAGS)" -o $(BINARY_DIR)/$(BINARY_NAME)-windows-amd64.exe
	@zip -j dist/$(BINARY_NAME)-windows-amd64.zip $(BINARY_DIR)/$(BINARY_NAME)-windows-amd64.exe
	@printf "$(GREEN)Distribution packages created in dist/$(NC)\n"

# Show help
help:
	@echo "Available targets:"
	@echo "  make              : Build and test the project"
	@echo "  make build        : Build the binary"
	@echo "  make clean        : Remove build artifacts"
	@echo "  make test         : Run tests"
	@echo "  make coverage     : Generate test coverage report"
	@echo "  make docker-build : Build Docker image"
	@echo "  make docker-run   : Run in Docker container"
	@echo "  make fmt          : Format code"
	@echo "  make lint         : Run linter"
	@echo "  make install-deps : Install dependencies"
	@echo "  make run          : Build and run the application"
	@echo "  make install      : Install binary to GOPATH/bin"
	@echo "  make release      : Create optimized release build"
	@echo "  make dist         : Create distribution packages"
	@echo "  make help         : Show this help message"
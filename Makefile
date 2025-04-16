# Go executable name
APP_NAME := pixellock
VERSION := 1.0.0

# Source files and directories
SRCS := main.go internal/pixellock/pixellock.go
PKG_LIST := $(shell go list ./... | grep -v /vendor/)

# Output directory
OUTPUT_DIR := bin
DIST_DIR := dist

# Go compiler flags
GOFLAGS := -v
LDFLAGS := -ldflags "-X main.Version=${VERSION}"

.DEFAULT_GOAL := help

# Build the application
build:
	@mkdir -p $(OUTPUT_DIR)
	go build $(GOFLAGS) $(LDFLAGS) -o $(OUTPUT_DIR)/$(APP_NAME) $(SRCS)
	@echo "Build completed successfully!"

# Run the application
run: build
	@$(OUTPUT_DIR)/$(APP_NAME) $(ARGS)

# Run tests with coverage
test:
	@mkdir -p $(OUTPUT_DIR)
	go test -v -race -coverprofile=$(OUTPUT_DIR)/coverage.out $(PKG_LIST)
	go tool cover -html=$(OUTPUT_DIR)/coverage.out -o $(OUTPUT_DIR)/coverage.html

# Run
# Install dependencies
install-deps:
	go mod tidy
	go mod vendor

# Build Docker image
docker-build:
	docker build -t $(APP_NAME) .

# Run Docker image
docker-run: docker-build
	docker run -it $(APP_NAME) $(DOCKER_ARGS)

# Show help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build         Build the application"
	@echo "  run           Build and run the application"
	@echo "  test          Run tests"
	@echo "  clean         Clean the project"
	@echo "  format        Format the code"
	@echo "  install-deps  Install dependencies using go mod"
	@echo "  docker-build  Build Docker image"
	@echo "  docker-run    Build and run Docker image"
	@echo "  help          Show this help message"

.PHONY: build run test clean format install-deps docker-build docker-run help

.PHONY: build clean test run lint vet fmt help hooks swagger-check

# Build variables
BINARY_NAME := swagger-to-http-file
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(BUILD_DATE)"

# Go variables
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOLINT := golint
GOVET := $(GOCMD) vet
GOFMT := gofmt -s -w
GOMOD := $(GOCMD) mod

# Directories
SRC_DIR := ./cmd/swagger-to-http-file
DIST_DIR := ./dist
BIN_DIR := ./bin

# Default target
.DEFAULT_GOAL := help

# Help text
help:
	@echo "Available targets:"
	@echo "  build        - Build the binary"
	@echo "  install      - Install the binary locally"
	@echo "  clean        - Remove build artifacts"
	@echo "  test         - Run tests"
	@echo "  coverage     - Run tests with coverage"
	@echo "  lint         - Run linter (requires golint)"
	@echo "  vet          - Run go vet"
	@echo "  fmt          - Run go fmt"
	@echo "  tidy         - Run go mod tidy"
	@echo "  release      - Prepare a release (requires goreleaser)"
	@echo "  docker       - Build Docker image"
	@echo "  run          - Run the binary with example input"
	@echo "  hooks        - Install Git hooks"
	@echo "  swagger-check - Check for Swagger changes and update HTTP files"
	@echo ""
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"

# Build the binary
build:
	@echo "Building $(BINARY_NAME) version $(VERSION)..."
	@mkdir -p $(BIN_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME) $(SRC_DIR)
	@echo "Build complete: $(BIN_DIR)/$(BINARY_NAME)"

# Install the binary
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GOCMD) install $(LDFLAGS) $(SRC_DIR)
	@echo "Installation complete"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BIN_DIR) $(DIST_DIR)
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run linter
lint:
	@echo "Running linter..."
	$(GOLINT) ./...

# Run go vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

# Run go fmt
fmt:
	@echo "Running go fmt..."
	$(GOFMT) ./...

# Run go mod tidy
tidy:
	@echo "Running go mod tidy..."
	$(GOMOD) tidy

# Release using goreleaser
release:
	@echo "Creating release with goreleaser..."
	goreleaser release --snapshot --rm-dist

# Build Docker image
docker:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):$(VERSION) .

# Run with example input
run:
	@echo "Running $(BINARY_NAME)..."
	@if [ -f "$(BIN_DIR)/$(BINARY_NAME)" ]; then \
		$(BIN_DIR)/$(BINARY_NAME) -i test/samples/petstore.json -o . -v; \
	else \
		echo "Binary not found. Run 'make build' first."; \
	fi

# Install Git hooks
hooks:
	@echo "Installing Git hooks..."
	@chmod +x scripts/install-hooks.sh
	@./scripts/install-hooks.sh

# Check for swagger changes and update HTTP files
swagger-check:
	@echo "Checking for Swagger file changes..."
	@chmod +x scripts/detect-swagger-changes.sh
	@./scripts/detect-swagger-changes.sh --verbose

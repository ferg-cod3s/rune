# Rune CLI Makefile
.PHONY: build test lint fmt clean dev install help

# Variables
BINARY_NAME=rune
MAIN_PATH=./cmd/rune
BUILD_DIR=./bin
VERSION?=dev
LDFLAGS=-ldflags "-X github.com/ferg-cod3s/rune/internal/commands.version=$(VERSION) -X github.com/ferg-cod3s/rune/internal/telemetry.version=$(VERSION) -X github.com/ferg-cod3s/rune/internal/telemetry.segmentWriteKey=$(RUNE_SEGMENT_WRITE_KEY) -X github.com/ferg-cod3s/rune/internal/telemetry.sentryDSN=$(RUNE_SENTRY_DSN)"

# Default target
all: build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Build for development (with race detection)
dev:
	@echo "Building $(BINARY_NAME) for development..."
	@mkdir -p $(BUILD_DIR)
	go build -race $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run tests in watch mode (requires entr)
test-watch:
	@echo "Running tests in watch mode..."
	find . -name "*.go" | entr -c go test ./...

# Lint the code
lint:
	@echo "Running linter..."
	golangci-lint run

# Format the code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

# Vet the code
vet:
	@echo "Vetting code..."
	go vet ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Install the binary to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	go install $(LDFLAGS) $(MAIN_PATH)

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	go run $(MAIN_PATH) $(ARGS)

# Generate shell completions
completions:
	@echo "Generating shell completions..."
	@mkdir -p completions
	$(BUILD_DIR)/$(BINARY_NAME) completion bash > completions/$(BINARY_NAME).bash
	$(BUILD_DIR)/$(BINARY_NAME) completion zsh > completions/$(BINARY_NAME).zsh
	$(BUILD_DIR)/$(BINARY_NAME) completion fish > completions/$(BINARY_NAME).fish

# Check for security vulnerabilities
security:
	@echo "Checking for security vulnerabilities..."
	govulncheck ./...

# Pre-commit checks (run before committing)
pre-commit: fmt vet lint test

# Help
help:
	@echo "Available targets:"
	@echo "  build        - Build the binary"
	@echo "  dev          - Build for development with race detection"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  test-watch   - Run tests in watch mode"
	@echo "  lint         - Run linter"
	@echo "  fmt          - Format code"
	@echo "  vet          - Vet code"
	@echo "  clean        - Clean build artifacts"
	@echo "  install      - Install binary to GOPATH/bin"
	@echo "  tidy         - Tidy dependencies"
	@echo "  deps         - Download dependencies"
	@echo "  build-all    - Build for multiple platforms"
	@echo "  run          - Run the application (use ARGS=... for arguments)"
	@echo "  completions  - Generate shell completions"
	@echo "  security     - Check for security vulnerabilities"
	@echo "  pre-commit   - Run pre-commit checks"
	@echo "  help         - Show this help"
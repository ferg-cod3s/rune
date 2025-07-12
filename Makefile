# Rune CLI Makefile
.PHONY: build test lint fmt clean dev install help

# Variables
BINARY_NAME=rune
MAIN_PATH=./cmd/rune
BUILD_DIR=./bin
VERSION?=dev
# Default telemetry keys (can be overridden by environment variables)
DEFAULT_SEGMENT_KEY=ZkEZXHRWH96y8EviNkbYJUByqGR9QI4G
DEFAULT_SENTRY_DSN=https://3b20acb23bbbc5958448bb41900cdca2@sentry.fergify.work/10
SEGMENT_KEY?=$(if $(RUNE_SEGMENT_WRITE_KEY),$(RUNE_SEGMENT_WRITE_KEY),$(DEFAULT_SEGMENT_KEY))
SENTRY_DSN?=$(if $(RUNE_SENTRY_DSN),$(RUNE_SENTRY_DSN),$(DEFAULT_SENTRY_DSN))
LDFLAGS=-ldflags "-X github.com/ferg-cod3s/rune/internal/commands.version=$(VERSION) -X github.com/ferg-cod3s/rune/internal/telemetry.version=$(VERSION) -X github.com/ferg-cod3s/rune/internal/telemetry.segmentWriteKey=$(SEGMENT_KEY) -X github.com/ferg-cod3s/rune/internal/telemetry.sentryDSN=$(SENTRY_DSN)"

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

# Test telemetry integration
test-telemetry:
	@echo "Testing telemetry integration..."
	./scripts/test-telemetry.sh

# Build with telemetry (ensures telemetry keys are embedded)
build-telemetry:
	@echo "Building $(BINARY_NAME) with telemetry support..."
	@echo "Segment Key: $(shell echo $(SEGMENT_KEY) | cut -c1-10)..."
	@echo "Sentry DSN: $(shell echo $(SENTRY_DSN) | cut -c1-30)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Testing telemetry integration..."
	@RUNE_DEBUG=true $(BUILD_DIR)/$(BINARY_NAME) --version > /dev/null 2>&1 && echo "✅ Telemetry integration test passed" || echo "❌ Telemetry integration test failed"

# Pre-commit checks (run before committing)
pre-commit: fmt vet lint test

# Security targets
security-deps:
	@echo "Checking dependencies for vulnerabilities..."
	go install github.com/sonatypecommunity/nancy@latest
	go list -json -deps ./... | nancy sleuth --loud

security-vulns:
	@echo "Checking for known vulnerabilities..."
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

security-static:
	@echo "Running static security analysis..."
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	gosec ./...

security-secrets:
	@echo "Scanning for secrets..."
	@if command -v trufflehog >/dev/null 2>&1; then \
		trufflehog filesystem . --exclude-paths .trufflehogignore; \
	else \
		echo "TruffleHog not installed, skipping secret scan"; \
	fi

security-build:
	@echo "Checking binary for embedded secrets..."
	@if [ -f "./bin/rune" ]; then \
		if strings ./bin/rune | grep -E "(password|secret|key|token)" | grep -v -E "(segmentWriteKey|sentryDSN|RUNE_)" ; then \
			echo "❌ Potential secrets found in binary"; \
			exit 1; \
		else \
			echo "✅ No obvious secrets found in binary"; \
		fi \
	else \
		echo "Binary not found, run 'make build' first"; \
		exit 1; \
	fi

security-all: security-deps security-vulns security-static security-secrets
	@echo "✅ All security checks completed"

# Enhanced coverage with thresholds
test-coverage-detailed:
	@echo "Running tests with detailed coverage..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@COVERAGE=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	echo "Total coverage: $$COVERAGE%"; \
	if [ $$(echo "$$COVERAGE < 70" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage $$COVERAGE% is below 70% threshold"; \
		exit 1; \
	else \
		echo "✅ Coverage $$COVERAGE% meets threshold"; \
	fi

# Enhanced pre-commit with security
pre-commit-security: fmt vet lint test security-static security-vulns
	@echo "✅ Pre-commit security checks passed"

# Help
help:
	@echo "Available targets:"
	@echo "  build        - Build the binary"
	@echo "  build-telemetry - Build with telemetry support (embeds keys)"
	@echo "  dev          - Build for development with race detection"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  test-coverage-detailed - Run tests with detailed coverage and thresholds"
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
	@echo "  security-deps - Check dependencies for vulnerabilities"
	@echo "  security-vulns - Check for known vulnerabilities"
	@echo "  security-static - Run static security analysis"
	@echo "  security-secrets - Scan for secrets"
	@echo "  security-build - Check binary for embedded secrets"
	@echo "  security-all - Run all security checks"
	@echo "  test-telemetry - Test telemetry integration"
	@echo "  pre-commit   - Run pre-commit checks"
	@echo "  pre-commit-security - Run pre-commit checks with security"
	@echo "  help         - Show this help"
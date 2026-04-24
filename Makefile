# Makefile for DWD Open Data Downloader (Go Version)

# Binary name
BINARY_NAME := downloader

# Build flags
GO_BUILD_FLAGS := -ldflags "-X main.Version=$(VERSION)"
VERSION := $(shell cat internal/version/version.go | grep "Version = " | sed 's/.*Version = //' | tr -d '"')

# Directories
BIN_DIR := bin
DIST_DIR := dist

# Go parameters
GO ?= go
GOFMT ?= gofmt
GOVET ?= go vet

.PHONY: all build test fmt vet clean install run deps lint tidy

# Default target
all: build

# Download dependencies
deps:
	$(GO) mod download
	$(GO) mod tidy

# Build the binary
build: deps
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/downloader
	@echo "Built: $(BIN_DIR)/$(BINARY_NAME) version $(VERSION)"

# Build for release (multiple platforms)
release: clean
	@mkdir -p $(DIST_DIR)
	# Linux amd64
	GOOS=linux GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/downloader
	# Linux arm64
	GOOS=linux GOARCH=arm64 $(GO) build $(GO_BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/downloader
	# macOS amd64
	GOOS=darwin GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/downloader
	# macOS arm64
	GOOS=darwin GOARCH=arm64 $(GO) build $(GO_BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/downloader
	# Windows amd64
	GOOS=windows GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/downloader
	@echo "Release builds created in $(DIST_DIR)/"

# Install the binary
install: deps
	$(GO) install $(GO_BUILD_FLAGS) ./cmd/downloader
	@echo "Installed to $$GOPATH/bin/$(BINARY_NAME)"

# Run the application
run: deps
	$(GO) run ./cmd/downloader

# Run tests
test:
	$(GO) test -v -race ./...

# Run tests with coverage
coverage:
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format code
fmt:
	$(GOFMT) -s -w .
	@echo "Code formatted"

# Vet code
vet:
	$(GOVET) ./...
	@echo "Vetting complete"

# Run linter (uses golangci-lint if available)
lint:
	$(GO) run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run ./...
	@echo "Linting complete"

# Tidy go.mod
tidy:
	$(GO) mod tidy
	@echo "go.mod tidied"

# Run all checks (fmt, vet, test)
check: fmt vet test
	@echo "All checks passed"

# Clean build artifacts
clean:
	rm -rf $(BIN_DIR) $(DIST_DIR) coverage.out coverage.html
	@echo "Cleaned"

# Show help
help:
	@echo "Available targets:"
	@echo "  all        - Build the binary (default)"
	@echo "  build      - Build the binary"
	@echo "  deps       - Download dependencies"
	@echo "  test       - Run tests"
	@echo "  coverage    - Run tests with coverage report"
	@echo "  fmt        - Format code"
	@echo "  vet        - Run go vet"
	@echo "  lint       - Run linter"
	@echo "  tidy       - Tidy go.mod"
	@echo "  check      - Run fmt, vet, and test"
	@echo "  install    - Install the binary"
	@echo "  run        - Run the application"
	@echo "  release    - Build binaries for all platforms"
	@echo "  clean      - Remove build artifacts"
	@echo "  help       - Show this help message"

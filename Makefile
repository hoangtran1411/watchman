# =============================================================================
# WATCHMEN - Makefile
# =============================================================================
# Build automation for Windows development
# =============================================================================

.PHONY: all build clean test lint fmt help install uninstall run check

# Variables
BINARY_NAME=watchmen.exe
CMD_PATH=./cmd/watchmen
VERSION=$(shell git describe --tags --always --dirty 2>NUL || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>NUL || echo "unknown")
BUILD_DATE=$(shell powershell -Command "Get-Date -Format 'yyyy-MM-ddTHH:mm:ssZ'")
LDFLAGS=-s -w -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildDate=$(BUILD_DATE)

# Default target
all: lint test build

# =============================================================================
# Build
# =============================================================================

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build -ldflags="$(LDFLAGS)" -o $(BINARY_NAME) $(CMD_PATH)
	@echo "Build complete: $(BINARY_NAME)"

## build-dev: Build with debug info
build-dev:
	@echo "Building $(BINARY_NAME) (debug)..."
	go build -o $(BINARY_NAME) $(CMD_PATH)

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	@if exist $(BINARY_NAME) del /F $(BINARY_NAME)
	@if exist coverage.out del /F coverage.out
	@echo "Clean complete"

# =============================================================================
# Testing
# =============================================================================

## test: Run all tests
test:
	@echo "Running tests..."
	go test -v -race -coverprofile coverage.out ./...

## test-short: Run tests without race detector
test-short:
	@echo "Running tests (short)..."
	go test -v -coverprofile coverage.out ./...

## coverage: Show test coverage
coverage: test
	@echo "Coverage report:"
	go tool cover -func=coverage.out

## coverage-html: Open coverage in browser
coverage-html: test
	go tool cover -html=coverage.out

# =============================================================================
# Linting & Formatting
# =============================================================================

## lint: Run golangci-lint
lint:
	@echo "Running linter..."
	golangci-lint run ./...

## lint-fix: Run linter with auto-fix
lint-fix:
	@echo "Running linter with auto-fix..."
	golangci-lint run --fix ./...

## fmt: Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w -local github.com/hoangtran1411/watchman .

## vet: Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# =============================================================================
# Development
# =============================================================================

## run: Run the application
run: build
	./$(BINARY_NAME)

## check: Run check command manually
check: build
	./$(BINARY_NAME) check --output json

## version: Show version
version: build
	./$(BINARY_NAME) version

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

## deps-update: Update all dependencies
deps-update:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# =============================================================================
# Installation
# =============================================================================

## install: Install as Windows Service
install: build
	@echo "Installing Watchmen service..."
	powershell -ExecutionPolicy Bypass -File scripts/install.ps1

## uninstall: Uninstall Windows Service
uninstall:
	@echo "Uninstalling Watchmen service..."
	powershell -ExecutionPolicy Bypass -File scripts/uninstall.ps1

# =============================================================================
# Help
# =============================================================================

## help: Show this help
help:
	@echo.
	@echo Watchmen - SQL Server Agent Job Monitor
	@echo.
	@echo Usage:
	@echo   make [target]
	@echo.
	@echo Targets:
	@echo   build         Build the binary
	@echo   build-dev     Build with debug info
	@echo   clean         Remove build artifacts
	@echo   test          Run all tests
	@echo   test-short    Run tests without race detector
	@echo   coverage      Show test coverage
	@echo   coverage-html Open coverage in browser
	@echo   lint          Run golangci-lint
	@echo   lint-fix      Run linter with auto-fix
	@echo   fmt           Format code
	@echo   vet           Run go vet
	@echo   run           Run the application
	@echo   check         Run check command
	@echo   version       Show version
	@echo   deps          Download dependencies
	@echo   deps-update   Update all dependencies
	@echo   install       Install as Windows Service
	@echo   uninstall     Uninstall Windows Service
	@echo   help          Show this help
	@echo.

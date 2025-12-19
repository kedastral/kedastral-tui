.PHONY: build run clean test fmt help

VERSION ?= dev
LDFLAGS := -X main.version=$(VERSION)

## build: Build the TUI binary
build:
	@echo "Building kedastral-tui $(VERSION)..."
	@mkdir -p bin
	@go build -ldflags "$(LDFLAGS)" -o bin/kedastral-tui .

## run: Run the TUI (use env vars or flags to configure)
run:
	@go run .

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/

## test: Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

## fmt: Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@grep -E '^##' Makefile | sed 's/## //'

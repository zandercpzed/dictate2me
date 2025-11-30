.PHONY: help build test lint clean install run dev

# Variables
BINARY_NAME=dictate2me
BINARY_DAEMON=dictate2me-daemon
BUILD_DIR=bin
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/$(BUILD_DIR)

# Colors for output
CYAN=\033[0;36m
GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m # No Color

## help: Display this help message
help:
	@echo "${CYAN}dictate2me Makefile${NC}"
	@echo ""
	@echo "${GREEN}Available targets:${NC}"
	@grep -E '^## ' Makefile | sed 's/## /  /' | column -t -s ':'

## build: Build the application
build:
	@echo "${CYAN}Building ${BINARY_NAME}...${NC}"
	@mkdir -p "$(GOBIN)"
	@go build -o "$(GOBIN)/$(BINARY_NAME)" ./cmd/$(BINARY_NAME)
	@go build -o "$(GOBIN)/$(BINARY_DAEMON)" ./cmd/$(BINARY_DAEMON)
	@echo "${GREEN}✓ Build complete${NC}"

## test: Run all tests
test:
	@echo "${CYAN}Running tests...${NC}"
	@go test -v -race -coverprofile=coverage.out ./...
	@echo "${GREEN}✓ Tests passed${NC}"

## test-coverage: Run tests with coverage report
test-coverage: test
	@echo "${CYAN}Generating coverage report...${NC}"
	@go tool cover -html=coverage.out -o coverage.html
	@echo "${GREEN}✓ Coverage report generated: coverage.html${NC}"

## lint: Run linters
lint:
	@echo "${CYAN}Running linters...${NC}"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
		echo "${GREEN}✓ Linting complete${NC}"; \
	else \
		echo "${YELLOW}⚠ golangci-lint not installed. Run: brew install golangci-lint${NC}"; \
	fi

## fmt: Format code
fmt:
	@echo "${CYAN}Formatting code...${NC}"
	@go fmt ./...
	@goimports -w .
	@echo "${GREEN}✓ Code formatted${NC}"

## vet: Run go vet
vet:
	@echo "${CYAN}Running go vet...${NC}"
	@go vet ./...
	@echo "${GREEN}✓ Vet complete${NC}"

## clean: Remove build artifacts
clean:
	@echo "${CYAN}Cleaning...${NC}"
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "${GREEN}✓ Clean complete${NC}"

## install: Install the application
install: build
	@echo "${CYAN}Installing...${NC}"
	@go install ./cmd/$(BINARY_NAME)
	@go install ./cmd/$(BINARY_DAEMON)
	@echo "${GREEN}✓ Installed to $(shell go env GOPATH)/bin${NC}"

## run: Run the application
run: build
	@echo "${CYAN}Running ${BINARY_NAME}...${NC}"
	@"$(GOBIN)/$(BINARY_NAME)"

## dev: Run in development mode with hot reload
dev:
	@echo "${CYAN}Starting development mode...${NC}"
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "${YELLOW}⚠ air not installed. Run: go install github.com/air-verse/air@latest${NC}"; \
		$(MAKE) run; \
	fi

## deps: Download dependencies
deps:
	@echo "${CYAN}Downloading dependencies...${NC}"
	@go mod download
	@go mod tidy
	@echo "${GREEN}✓ Dependencies downloaded${NC}"

## setup: Setup development environment
setup:
	@echo "${CYAN}Setting up development environment...${NC}"
	@./scripts/setup-dev.sh
	@echo "${GREEN}✓ Setup complete${NC}"

## check: Run all checks (fmt, vet, lint, test)
check: fmt vet lint test
	@echo "${GREEN}✓ All checks passed${NC}"

## models: Download AI models
models:
	@echo "${CYAN}Downloading AI models...${NC}"
	@./scripts/download-models.sh
	@echo "${GREEN}✓ Models downloaded${NC}"

.DEFAULT_GOAL := help

# ChatApp Makefile
# Basic functionality for building, running, migrating, and managing the chat application

# Variables
APP_NAME := chat-app
BINARY_NAME := chat-app
BUILD_DIR := ./bin

# Go specific
GO := go
GO_VERSION := 1.24.0
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

# Database
DB_HOST := 127.0.0.1
DB_PORT := 3306
DB_USER := chat
DB_PASSWORD := chat
DB_NAME := chat_app

.PHONY: build
build:
	@mkdir -p $(BUILD_DIR)
	@$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

.PHONY: run
run:
	@$(GO) run ./cmd/server

.PHONY: migrate
migrate:
	@echo "Running migrations..."
	@$(GO) run ./cmd/migrate -action=up

.PHONY: migrate-status
migrate-status:
	@echo "Checking migration status..."
	@$(GO) run ./cmd/migrate -action=status

.PHONY: migrate-rollback
migrate-rollback:
	@if [ -z "$(steps)" ]; then echo "Usage: make migrate-rollback steps=1"; exit 1; fi
	@echo "Rolling back $(steps) migration(s)..."
	@$(GO) run ./cmd/migrate -action=down -steps=$(steps)

.PHONY: test
test:
	@$(GO) test -v ./...

.PHONY: fmt
fmt:
	@$(GO) fmt ./...

.PHONY: vet
vet:
	@$(GO) vet ./...

.PHONY: tidy
tidy:
	@$(GO) mod tidy

.PHONY: deps
deps:
	@$(GO) mod download

.PHONY: clean
clean:
	@rm -rf $(BUILD_DIR)
	@$(GO) clean

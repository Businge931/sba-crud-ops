include .envrc

APP = sba-crud-ops
GOBASE = $(shell pwd)
GOBIN = $(GOBASE)/build/bin
LINT_PATH = $(GOBASE)/build/lint
MAIN_APP = $(GOBASE)/cmd
MIGRATIONS_PATH=$(GOBASE)/migrations

# Default database connection details (matches docker-compose.yml)
DB_USER ?= admin
DB_PASSWORD ?= adminpassword
DB_NAME ?= sba_crud_ops
DB_HOST ?= localhost
DB_PORT ?= 5433
DB_ADDR ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable


help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## Fetch required dependencies
	go mod tidy -compat=1.22
	go mod download

build: ## Build the application
	go build -o $(GOBIN)/$(APP) $(MAIN_APP)

run: build ## Build and run program
	cd $(MAIN_APP) && go run .

lint: install-golangci ## Linter for developers
	$(LINT_PATH)/golangci-lint run --timeout=5m -c .golangci.yml

lint-fix:
	$(LINT_PATH)/golangci-lint run --timeout=5m -c .golangci.yml --fix

install-golangci: ## Install the correct version of lint
	@GOBIN=$(LINT_PATH) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.58.1

migrate-up: ## Run the migration up
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up $(filter-out $@,$(MAKECMDGOALS))

migrate-down: ## Run the migration down
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down

migration: ## Create a new migration
	@migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(filter-out $@,$(MAKECMDGOALS))		


docker-up: ## Run the docker-compose
	docker compose up -d

docker-down: ## Stop the docker-compose
	docker compose down

docker-restart: ## Restart the docker-compose
	docker-compose restart

test: ## Run all unit tests
	go test ./internal/...

test-validator: ## Run validator tests
	go test ./internal/core/validator

test-service: ## Run service tests
	go test ./internal/service

test-grpc: ## Run gRPC API tests
	go test ./internal/app/grpc

test-integration: ## Run integration tests for repository (requires test database)
	INTEGRATION_TEST=true go test ./internal/repository/postgres

test-cover: ## Run tests with coverage
	go test ./internal/... -cover

test-coverage: ## Run tests and generate coverage profile
	go test ./internal/... -coverprofile=coverage.out

test-coverage-html: ## Generate HTML coverage report
	go tool cover -html=coverage.out -o coverage.html


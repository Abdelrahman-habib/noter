# Environment file selection based on ENV variable
ENV ?= development
ENV_FILE = ./$(ENV).env

.PHONY: run run-dev run-prod run-test run-all migrate-up migrate-down migrate-reset migrate-status seed-up seed-down seed-reset build build-dev build-prod build-test build-all clean test test-short test-cover lint audit upgrade-tool migrate-test-up migrate-test-down migrate-test-reset seed-test-up seed-test-down seed-test-reset init-test-db teardown-test-db

# Check if environment-specific .env file exists, fallback to .env
ifneq (,$(wildcard $(ENV_FILE)))
    include $(ENV_FILE)
    export
    # Goose rejects an explicitly named env file if it is missing.
    GOOSE_ENV = -env $(ENV_FILE)
else ifneq (,$(wildcard ./.env))
    include .env
    export
endif

run:
	go run ./cmd/web -addr="$(HOST):$(PORT)"

run-dev:
	@$(MAKE) run ENV=development

run-prod:
	@$(MAKE) run ENV=production

run-test:
	@$(MAKE) run ENV=test

run-all:
	make run-dev && make run-prod && make run-test

migrate-up:
	goose -dir db/schema/migrations -table goose_migrations $(GOOSE_ENV) up

migrate-down:
	goose -dir db/schema/migrations -table goose_migrations $(GOOSE_ENV) down

migrate-reset:
	goose -dir db/schema/migrations -table goose_migrations $(GOOSE_ENV) reset

migrate-status:
	goose -dir db/schema/migrations -table goose_migrations $(GOOSE_ENV) status

seed-up:
	@if [ "$(ENVIROMENT)" = "test" ] || [ "$(ENVIROMENT)" = "development" ]; then \
		goose -dir db/schema/seed -no-versioning $(GOOSE_ENV) up; \
	else \
		echo "seed is only allowed in test and development environment"; \
		exit 1; \
	fi
seed-down:
	@if [ "$(ENVIROMENT)" = "test" ] || [ "$(ENVIROMENT)" = "development" ]; then \
		goose -dir db/schema/seed -no-versioning $(GOOSE_ENV) down; \
	else \
		echo "seed is only allowed in test and development environment"; \
		exit 1; \
	fi

seed-reset:
	@if [ "$(ENVIROMENT)" = "test" ] || [ "$(ENVIROMENT)" = "development" ]; then \
		goose -dir db/schema/seed -no-versioning $(GOOSE_ENV) reset; \
	else \
		echo "seed is only allowed in test and development environment"; \
		exit 1; \
	fi

migrate-test-up:
	goose -dir db/schema/migrations postgres '$(TEST_DB_DSN)' up

migrate-test-down:
	goose -dir db/schema/migrations postgres '$(TEST_DB_DSN)' down

migrate-test-reset:
	goose -dir db/schema/migrations postgres '$(TEST_DB_DSN)' reset

seed-test-up:
	goose -dir db/schema/seed -no-versioning postgres '$(TEST_DB_DSN)' up

seed-test-down:
	goose -dir db/schema/seed -no-versioning postgres '$(TEST_DB_DSN)' down

seed-test-reset:
	goose -dir db/schema/seed -no-versioning postgres '$(TEST_DB_DSN)' reset

init-test-db:
	@$(MAKE) migrate-test-up
	@$(MAKE) seed-test-up

teardown-test-db:
	@$(MAKE) seed-test-reset
	@$(MAKE) migrate-test-reset

build:
	go build -o bin/noter ./cmd/web

build-dev:
	@$(MAKE) build ENV=development

build-prod:
	@$(MAKE) build ENV=production

build-test:
	@$(MAKE) build ENV=test

build-all:
	make build-dev && make build-prod && make build-test

clean:
	rm -f bin/noter

test:
	go test ./...

# Skip database integration tests
test-short:
	go test -short ./...

test-cover:
	go test -covermode=atomic -coverprofile=.profile.out ./...
	go tool cover -html=.profile.out

lint:
	go vet ./...

audit:
	go vet ./...
	go tool -modfile=go.tool.mod staticcheck ./...
	go tool -modfile=go.tool.mod govulncheck

add-tool:
	@if [ -z "${tool}" ]; then \
		echo "Tool is required example: make add-tool tool=golang.org/x/vuln/cmd/govulncheck"; \
		exit 1; \
	fi
	@echo "Adding tool: ${tool}"
	go get -tool -modfile=go.tool.mod ${tool}

upgrade-tool:
	@if [ -z "${tool}" ]; then \
		echo "Tool is required example: make upgrade-tool tool=golang.org/x/vuln/cmd/govulncheck"; \
		exit 1; \
	fi
	@echo "Upgrading tool: ${tool}"
	go get -tool -modfile=go.tool.mod ${tool}@latest

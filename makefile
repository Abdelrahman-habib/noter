# Environment file selection based on ENV variable
ENV ?= development
ENV_FILE = .env.$(ENV)

.PHONY: run run-dev run-prod run-test run-all migrate-up migrate-down seed-up seed-down seed-reset build build-dev build-prod build-test build-all clean test lint audit migrate-test-up migrate-test-down migrate-test-reset seed-test-up seed-test-down seed-test-reset init-test-db teardown-test-db migrate-reset

# Check if environment-specific .env file exists, fallback to .env
ifneq (,$(wildcard $(ENV_FILE)))
    include $(ENV_FILE)
    export
else ifneq (,$(wildcard ./.env))
    include .env
    export
endif

run:
	go run ./cmd/web -addr="$(HOST):$(PORT)" -env=$(ENVIROMENT) -dsn='$(DB_DSN)' -tls-cert='$(TLS_CERT)' -tls-key='$(TLS_KEY)' -debug

run-dev:
	@$(MAKE) run ENV=development

run-prod:
	@$(MAKE) run ENV=production

run-test:
	@$(MAKE) run ENV=test

run-all:
	make run-dev && make run-prod && make run-test

migrate-up:
	goose -dir db/schema/migrations up

migrate-down:
	goose -dir db/schema/migrations down

migrate-reset:
	goose -dir db/schema/migrations reset

seed-up:
	@if [ "$(ENVIROMENT)" == "test" ] || [ "$(ENVIROMENT)" == "development" ]; then \
		goose -dir db/schema/seed -no-versioning up; \
	else \
		echo "seed is only allowed in test and development environment"; \
		exit 1; \
	fi
seed-down:
	@if [ "$(ENVIROMENT)" == "test" ] || [ "$(ENVIROMENT)" == "development" ]; then \
		goose -dir db/schema/seed -no-versioning down; \
	else \
		echo "seed is only allowed in test and development environment"; \
		exit 1; \
	fi

seed-reset:
	@if [ "$(ENVIROMENT)" == "test" ] || [ "$(ENVIROMENT)" == "development" ]; then \
		goose -dir db/schema/seed -no-versioning reset; \
	else \
		echo "seed is only allowed in test and development environment"; \
		exit 1; \
	fi

migrate-test-up:
	goose dir db/schema/migrations mysql '$(TEST_DB_DSN)' up

migrate-test-down:
	goose dir db/schema/migrations mysql '$(TEST_DB_DSN)' down

migrate-test-reset:
	goose dir db/schema/migrations mysql '$(TEST_DB_DSN)' reset

seed-test-up:
	goose -dir db/schema/seed -no-versioning mysql '$(TEST_DB_DSN)' up

seed-test-down:
	goose -dir db/schema/seed -no-versioning mysql '$(TEST_DB_DSN)' down

seed-test-reset:
	goose -dir db/schema/seed -no-versioning mysql '$(TEST_DB_DSN)' reset

init-test-db:
	@$(MAKE) migrate-test-up
	@$(MAKE) seed-test-up

teardown-test-db:
	@$(MAKE) seed-test-reset
	@$(MAKE) migrate-test-reset

build:
	go build -o bin/snippetbox ./cmd/web

build-dev:
	@$(MAKE) build ENV=development

build-prod:
	@$(MAKE) build ENV=production

build-test:
	@$(MAKE) build ENV=test

build-all:
	make build-dev && make build-prod && make build-test

clean:
	rm -f bin/snippetbox

test:
	go test ./...

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
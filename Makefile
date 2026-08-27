SHELL := /bin/bash

APP := server
GO ?= go

.PHONY: tidy fmt vet test build run docker-build up down logs

tidy:
	$(GO) mod tidy

fmt:
	gofmt -l -w .

vet:
	$(GO) vet ./...

test:
	$(GO) test ./...

build:
	CGO_ENABLED=0 $(GO) build -o bin/$(APP) ./cmd/server

run: build
	PORT=${PORT:-8000} \
	CLICKHOUSE_DSN=${CLICKHOUSE_DSN:-clickhouse://default:default@localhost:9000/default?dial_timeout=5s\&compress=true} \
	LOG_LEVEL=${LOG_LEVEL:-info} \
	SHUTDOWN_TIMEOUT_SECONDS=${SHUTDOWN_TIMEOUT_SECONDS:-10} \
		./bin/$(APP)

docker-build:
	docker build -t otel-map-server:local .

up:
	docker compose up -d

down:
	docker compose down -v

logs:
	docker compose logs -f clickhouse

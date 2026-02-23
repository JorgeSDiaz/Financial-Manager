.PHONY: run build test lint tidy dev dev-build dev-down help

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  run        Run the API server"
	@echo "  build      Build the API binary to bin/api"
	@echo "  test       Run tests with race detector and coverage"
	@echo "  lint       Run golangci-lint"
	@echo "  tidy       Run go mod tidy"
	@echo "  dev-build  Build Docker Compose services"
	@echo "  dev        Start Docker Compose services"
	@echo "  dev-down   Stop and remove Docker Compose services"

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

test:
	go test $(shell go list ./... | grep -v '/mocks') -v -race -coverprofile=coverage.out
	go tool cover -func=coverage.out

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

dev-build:
	docker compose build

dev:
	docker compose up

dev-down:
	docker compose down --remove-orphans

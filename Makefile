.PHONY: run build test lint tidy docker-run docker-build docker-down help

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  run           Run the API server"
	@echo "  build         Build the API binary to bin/api"
	@echo "  test          Run tests with race detector and coverage"
	@echo "  lint          Run golangci-lint"
	@echo "  tidy          Run go mod tidy"
	@echo "  docker-build  Build Docker Compose services"
	@echo "  docker-run    Start Docker Compose services in detached mode"
	@echo "  docker-down   Stop and remove Docker Compose services"

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

docker-build:
	docker compose build

docker-run:
	docker compose up -d

docker-down:
	docker compose down --remove-orphans

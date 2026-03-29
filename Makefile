.PHONY: run build clean test test-verbose migrate-up migrate-down migrate-status docs run-mcp \
        docker-build docker-up docker-down \
        podman-build podman-up podman-down \
        help

SWAG=$(shell go env GOPATH)/bin/swag
GOOSE=$(shell go env GOPATH)/bin/goose

include .env
export

help:
	@echo ""
	@echo "Development"
	@echo "  run              Run the API server"
	@echo "  run-mcp          Run the MCP SSE server"
	@echo "  build            Compile binary to bin/api"
	@echo "  clean            Remove compiled binaries"
	@echo ""
	@echo "Testing"
	@echo "  test             Run all tests"
	@echo "  test-verbose     Run tests with verbose output"
	@echo ""
	@echo "Database"
	@echo "  migrate-up       Apply pending migrations"
	@echo "  migrate-down     Roll back last migration"
	@echo "  migrate-status   Show migration status"
	@echo ""
	@echo "Docs"
	@echo "  docs             Regenerate OpenAPI spec from annotations"
	@echo ""
	@echo "Docker"
	@echo "  docker-build     Build Docker image"
	@echo "  docker-up        Start app via docker compose"
	@echo "  docker-down      Stop docker compose services"
	@echo ""
	@echo "Podman"
	@echo "  podman-build     Build image with podman"
	@echo "  podman-up        Start app via podman compose"
	@echo "  podman-down      Stop podman compose services"
	@echo ""

run:
	go run ./cmd/api

run-mcp:
	go run ./cmd/mcp

build:
	go build -o bin/api ./cmd/api

clean:
	rm -rf bin/

test:
	go test ./...

test-verbose:
	go test -v ./...

docs:
	$(SWAG) init -g cmd/api/main.go --output docs

migrate-up:
	$(GOOSE) -dir db/migrations postgres "$(MIGRATION_URL)" up

migrate-down:
	$(GOOSE) -dir db/migrations postgres "$(MIGRATION_URL)" down

migrate-status:
	$(GOOSE) -dir db/migrations postgres "$(MIGRATION_URL)" status

docker-build:
	docker build -t hello-world-go .

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

podman-build:
	podman build -t hello-world-go .

podman-up:
	podman compose up --build -d

podman-down:
	podman compose down

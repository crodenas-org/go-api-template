.PHONY: run build clean test test-verbose migrate-up migrate-down migrate-status docs

SWAG=$(shell go env GOPATH)/bin/swag

include .env
export

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

clean:
	rm -rf bin/

test:
	go test ./...

test-verbose:
	go test -v ./...

GOOSE=$(shell go env GOPATH)/bin/goose

docs:
	$(SWAG) init -g cmd/api/main.go --output docs

migrate-up:
	$(GOOSE) -dir db/migrations postgres "$(MIGRATION_URL)" up

migrate-down:
	$(GOOSE) -dir db/migrations postgres "$(MIGRATION_URL)" down

migrate-status:
	$(GOOSE) -dir db/migrations postgres "$(MIGRATION_URL)" status

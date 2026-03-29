.PHONY: run build clean test test-verbose

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

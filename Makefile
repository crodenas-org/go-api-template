.PHONY: run build clean

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

clean:
	rm -rf bin/

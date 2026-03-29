.PHONY: run build clean

run:
	go run main.go

build:
	go build -o bin/api main.go

clean:
	rm -rf bin/

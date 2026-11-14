.PHONY: run test build

run:
	go run ./cmd/api

test:
	go test -v -race ./...

build:
	go build -o bin/api ./cmd/api

.PHONY: all generate build test lint fmt sync-schema clean

all: lint build test

generate: sync-schema
	go generate ./...

sync-schema:
	bash scripts/sync-schema.sh

build:
	go build ./...

test:
	go test ./...

lint:
	go vet ./...

fmt:
	go fmt ./...

clean:
	go clean

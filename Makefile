.PHONY: all generate build test lint fmt sync-schema e2e e2e-list clean

E2E_RUN ?= .
E2E_COUNT ?= 1
E2E_TIMEOUT ?= 12m

all: lint build test

generate: sync-schema
	go generate ./...

sync-schema:
	bash scripts/sync-schema.sh

build:
	go build ./...

test:
	go test ./...

e2e:
	go test -tags=e2e ./e2e -run '$(value E2E_RUN)' -count=$(E2E_COUNT) -timeout=$(E2E_TIMEOUT) -v

e2e-list:
	go test -tags=e2e ./e2e -list .

lint:
	go vet ./...

fmt:
	go fmt ./...

clean:
	go clean

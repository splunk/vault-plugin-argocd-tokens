.DEFAULT_GOAL := all
.PHONY: all
all: build test lint

.PHONY: build
build:
	go install ./...

.PHONY: test
test:
	go test ./... -coverprofile cover.out

.PHONY: lint
lint:
	golangci-lint run --timeout=3m

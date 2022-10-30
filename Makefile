MAJOR_MINOR     := 0.1
BUILD           := $(shell  date -u "+%Y%m%d-%H%M%S")
SHORT_COMMIT    := $(shell git rev-parse --short HEAD 2>/dev/null || echo nocommitinfo)
VERSION         := $(MAJOR_MINOR).$(BUILD).$(SHORT_COMMIT)

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

publish:
	@./publish.sh $(VERSION)

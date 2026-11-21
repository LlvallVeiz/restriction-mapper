# restriction-mapper — Makefile
GO ?= go
BIN := rmap

.PHONY: all build test vet fmt check examples clean

all: check examples

build:
	$(GO) build -o $(BIN) ./cmd/rmap

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...

fmt:
	gofmt -l cmd internal tests

# Full verification: formatting, vet, all tests, then regenerate examples.
check: fmt vet test

examples: build
	./scripts/run-examples.sh

clean:
	rm -f $(BIN)

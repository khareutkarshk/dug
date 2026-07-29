BINARY := dug
CMD := ./cmd/edge

VERSION ?= dev
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GO_VERSION := $(shell go version | awk '{print $$3}')

LDFLAGS := \
	-X github.com/khareutkarshk/dug/internal/version.Version=$(VERSION) \
	-X github.com/khareutkarshk/dug/internal/version.Commit=$(COMMIT) \
	-X github.com/khareutkarshk/dug/internal/version.Date=$(DATE) \
	-X github.com/khareutkarshk/dug/internal/version.Go=$(GO_VERSION)

.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make run        - Run gateway"
	@echo "  make build      - Build binary"
	@echo "  make release    - Build release binary"
	@echo "  make test       - Run tests"
	@echo "  make race       - Run race detector"
	@echo "  make lint       - Run golangci-lint"
	@echo "  make fmt        - Format source"
	@echo "  make tidy       - Tidy dependencies"
	@echo "  make clean      - Remove binaries"

.PHONY: run
run:
	go run $(CMD)

.PHONY: build
build:
	go build -o $(BINARY) $(CMD)

.PHONY: release
release:
	CGO_ENABLED=0 go build \
		-ldflags "$(LDFLAGS) -s -w" \
		-o $(BINARY) \
		$(CMD)

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	go test ./...

.PHONY: race
race:
	go test -race ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: lint
lint:
	golangci-lint run

.PHONY: clean
clean:
	rm -f $(BINARY)
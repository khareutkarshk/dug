BINARY := dug
CMD := ./cmd/dug

VERSION ?= dev
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
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
	@echo "  make run        - Run gateway (dug run)"
	@echo "  make build      - Build binary with version metadata"
	@echo "  make install    - Install dug to GOPATH/bin"
	@echo "  make release    - Build stripped release binary"
	@echo "  make test       - Run tests"
	@echo "  make race       - Run race detector"
	@echo "  make lint       - Run golangci-lint"
	@echo "  make fmt        - Format source"
	@echo "  make tidy       - Tidy dependencies"
	@echo "  make clean      - Remove binaries"
	@echo "  make version    - Print embedded version via dug version"

.PHONY: run
run:
	go run $(CMD) run

.PHONY: build
build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) $(CMD)

.PHONY: install
install:
	go install -ldflags "$(LDFLAGS)" $(CMD)

.PHONY: release
release:
	CGO_ENABLED=0 go build \
		-ldflags "$(LDFLAGS) -s -w" \
		-o $(BINARY) \
		$(CMD)

.PHONY: version
version: build
	./$(BINARY) version

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

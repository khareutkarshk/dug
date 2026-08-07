# ============================
# Stage 1 - Build
# ============================
# BUILDPLATFORM: run the Go toolchain on the builder's native arch.
# TARGETOS/TARGETARCH: cross-compile for the image's target platform (Buildx).
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG VERSION=dev
ARG COMMIT=unknown
ARG DATE=unknown

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
    -trimpath \
    -ldflags="-s -w \
      -X github.com/khareutkarshk/dug/internal/version.Version=${VERSION} \
      -X github.com/khareutkarshk/dug/internal/version.Commit=${COMMIT} \
      -X github.com/khareutkarshk/dug/internal/version.BuildDate=${DATE}" \
    -o dug \
    ./cmd/dug


# ============================
# Stage 2 - Runtime
# ============================
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/dug .
COPY configs/docker.yaml ./configs/docker.yaml

EXPOSE 8080

CMD ["./dug", "run", "-config", "configs/docker.yaml"]

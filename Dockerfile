# ============================
# Stage 1 - Build
# ============================
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
ARG COMMIT=unknown
ARG DATE=unknown

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-s -w \
      -X github.com/khareutkarshk/dug/internal/version.Version=${VERSION} \
      -X github.com/khareutkarshk/dug/internal/version.Commit=${COMMIT} \
      -X github.com/khareutkarshk/dug/internal/version.Date=${DATE}" \
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

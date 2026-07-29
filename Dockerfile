# ============================
# Stage 1 - Build
# ============================
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install certificates and git
RUN apk add --no-cache git ca-certificates

# Copy dependency files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build DUG
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-s -w" \
    -o dug \
    ./cmd/edge


# ============================
# Stage 2 - Runtime
# ============================
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

# Copy binary
COPY --from=builder /app/dug .

# Copy docker configuration
COPY configs/docker.yaml ./configs/docker.yaml

EXPOSE 8080

CMD ["./dug", "-config", "configs/docker.yaml"]
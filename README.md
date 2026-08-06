# DUG (Distributed Unified Gateway)

[![Go Reference](https://pkg.go.dev/badge/github.com/khareutkarshk/dug.svg)](https://pkg.go.dev/github.com/khareutkarshk/dug)
[![Release](https://img.shields.io/github/v/release/khareutkarshk/dug)](https://github.com/khareutkarshk/dug/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/khareutkarshk/dug)](https://github.com/khareutkarshk/dug)
[![License](https://img.shields.io/github/license/khareutkarshk/dug?style=flat-square)](LICENSE)
[![Stars](https://img.shields.io/github/stars/khareutkarshk/dug?style=social)](https://github.com/khareutkarshk/dug/stargazers)

A lightweight, production-inspired API Gateway built from scratch in Go.

> **Build. Learn. Route.**

DUG is an API Gateway designed to understand how modern gateways like **NGINX**, **Traefik**, **Envoy**, and **Kong** work internally by implementing their core building blocks from first principles.

Rather than wrapping existing gateway libraries, DUG focuses on building the networking, reliability, and traffic management components yourself to gain a deeper understanding of production backend infrastructure.

---

# ✨ Features

## 🚦 Traffic Management

- ✅ Reverse Proxy
- ✅ Route-based Routing
- ✅ YAML Configuration
- ✅ Graceful Shutdown
- ✅ WebSocket Proxy
- ✅ Dynamic Configuration Reload

---

## ⚖️ Load Balancing

- ✅ Round Robin
- ✅ Smooth Weighted Round Robin (SWRR)
- ✅ Least Connections

---

## 🛡 Reliability

- ✅ Retry Mechanism
- ✅ Exponential Backoff
- ✅ Active Health Checks
- ✅ Passive Health Checks
- ✅ Circuit Breaker
- ✅ Request Timeouts

---

## 🔒 Security & Traffic Control

- ✅ HTTPS / TLS
- ✅ Per-IP Rate Limiting
- ✅ CORS
- ✅ Request Header Manipulation
- ✅ Response Header Manipulation

---

## 📊 Observability

- ✅ Prometheus Metrics
- ✅ Structured Logging (`slog`)
- ✅ Request IDs

---

## 🔍 Service Discovery

- ✅ Static Service Discovery
- ✅ DNS Service Discovery

---

## ⚙️ Configuration

- ✅ YAML Configuration
- ✅ Hot Configuration Reload
- ✅ Atomic Router Reload

---

# 🏗 Architecture

```text
                     Client
                        │
                 HTTP / HTTPS
                        │
                        ▼
              ┌────────────────────┐
              │        DUG         │
              ├────────────────────┤
              │ Reverse Proxy      │
              │ Router             │
              │ Load Balancer      │
              │ Retry Engine       │
              │ Health Checks      │
              │ Circuit Breaker    │
              │ Rate Limiter       │
              │ Header Middleware  │
              │ Service Discovery  │
              │ Metrics            │
              │ Logging            │
              └─────────┬──────────┘
                        │
        ┌───────────────┼───────────────┐
        ▼               ▼               ▼
    Backend A       Backend B       Backend C
```

---

# 🚀 Getting Started

## Install

### Option 1 — `go install` (recommended)

```bash
go install github.com/khareutkarshk/dug/cmd/dug@latest
```

Ensure `$(go env GOPATH)/bin` is on your `PATH`, then verify:

```bash
dug version
```

### Option 2 — Build from source

```bash
git clone https://github.com/khareutkarshk/dug.git
cd dug
go mod tidy
make build
./dug version
```

### Option 3 — Docker

Local compose:

```bash
docker compose up --build
```

Or pull a published release image from GitHub Container Registry:

```bash
docker pull ghcr.io/khareutkarshk/dug:latest
docker run --rm -p 8080:8080 ghcr.io/khareutkarshk/dug:v0.1.0
```

## Generate local TLS certificates

```bash
mkdir -p certs

openssl req -x509 \
-newkey rsa:4096 \
-keyout certs/server.key \
-out certs/server.crt \
-days 365 \
-nodes
```

## Start the gateway

```bash
dug run -config configs/edge.yaml
```

Or from source without installing:

```bash
go run ./cmd/dug run -config configs/edge.yaml
```

---

# 🧰 CLI

DUG ships as a multi-command CLI:

| Command | Description |
|---|---|
| `dug init` | Bootstrap a new DUG project |
| `dug run` | Start the gateway |
| `dug validate` | Validate a configuration file |
| `dug doctor` | Check config, port availability, TLS files, and upstreams |
| `dug version` | Print build metadata |
| `dug help` | Show usage |

```bash
# Bootstrap a project
dug init my-gateway
cd my-gateway
dug run -config configs/edge.yaml

# Start gateway
dug run -config configs/edge.yaml

# Validate config before deploy
dug validate -config configs/edge.yaml
dug validate -config configs/edge.yaml -quiet   # CI-friendly
dug validate -config configs/edge.yaml -json    # machine-readable

# Local diagnostics (non-zero exit if checks fail)
dug doctor -config configs/edge.yaml
dug doctor -config configs/edge.yaml -timeout 2s -json

# Version (human, short, machine-readable)
dug version
dug version -short
dug version -json
```

### Releases

Pushing a version tag (`v*`) runs `.github/workflows/release.yml`, which:

1. Cross-compiles binaries (linux/darwin/windows) with version ldflags
2. Packages them (`.tar.gz` / `.zip`) and publishes a GitHub Release
3. Builds and pushes a multi-arch image to GHCR (`linux/amd64`, `linux/arm64`)

```bash
# Binary (archive contains ./dug)
curl -LO https://github.com/khareutkarshk/dug/releases/latest/download/dug_v0.1.0_linux_amd64.tar.gz
tar -xzf dug_v0.1.0_linux_amd64.tar.gz
./dug version

# Container
docker pull ghcr.io/khareutkarshk/dug:v0.1.0
docker pull ghcr.io/khareutkarshk/dug:latest
```

To cut a release from `main`:

```bash
git tag v0.1.1
git push origin v0.1.1
```

---

# ⚙️ Configuration

DUG is configured using a YAML file.

Configuration supports:

- Routes
- Multiple upstreams
- Load balancing strategy
- Retry policy
- Health checks
- Timeouts
- HTTPS / TLS
- CORS
- Request & Response Headers
- Rate Limiting
- Service Discovery

Example configuration:

```yaml
server:
  port: 8080

routes:
  - path: /
    upstreams:
      - url: http://localhost:3001
        weight: 1
```

Default configuration:

```text
configs/edge.yaml
```

---

# 📁 Project Structure

```text
.
├── cmd/
│   └── dug/          # CLI entrypoint (run, validate, doctor, version, init)
├── configs/
├── examples/
│   ├── README.md             # Example catalog
│   ├── backend/              # Legacy demo upstream :3001
│   ├── backend2/             # Legacy demo upstream :3002
│   ├── echo/                 # Shared configurable mock upstream
│   ├── basic-reverse-proxy/
│   ├── docker-compose/
│   ├── microservices/
│   ├── nextjs-api/
│   ├── ai-gateway/
│   ├── blue-green/
│   └── websocket/
├── benchmarks/
├── internal/
│   ├── app/
│   ├── cli/
│   ├── config/
│   ├── discovery/
│   ├── httpx/
│   ├── logger/
│   ├── metrics/
│   ├── middleware/
│   ├── proxy/
│   ├── ratelimit/
│   ├── router/
│   ├── server/
│   ├── upstream/
│   └── version/
├── test/
└── README.md
```

See [examples/README.md](examples/README.md) for runnable demos (load balancing, path routing, Next.js, AI gateway, blue/green).

---

# 🧪 Testing

Run all tests

```bash
go test ./...
```

Run with race detector

```bash
go test ./... -race
```

Generate coverage

```bash
go test ./... -cover
```

---

# 📈 Current Status

DUG currently includes:

- Reverse Proxy
- Route-based Routing
- Multiple Load Balancing Algorithms
- Retry Engine
- Circuit Breaker
- Active & Passive Health Checks
- HTTPS
- WebSocket Proxy
- Dynamic Configuration Reload
- DNS Service Discovery
- Rate Limiting
- Header Manipulation
- Prometheus Metrics
- Structured Logging
- Request IDs

The core gateway functionality is complete and stable. The next milestone focuses on improving developer experience, packaging, and documentation.

---

# 🛣 Roadmap

## v1.0 — Developer Experience

### CLI

- [x] `dug init`
- [x] `dug run`
- [x] `dug validate`
- [x] `dug doctor`
- [x] `dug version` (`-short`, `-json`)

### Distribution

- [x] Docker Image
- [x] GitHub Releases (tag `v*` workflow)
- [x] Pre-built binaries
- [x] `go install` support

### Documentation

- [x] Installation Guide
- [ ] Configuration Reference
- [ ] Feature Documentation
- [ ] Example Configurations

### CI/CD

- [x] GitHub Actions
- [x] Go Race Tests
- [x] golangci-lint
- [x] Automated Releases (binaries + GHCR)

---

## Future

- JWT Authentication
- API Keys
- OpenTelemetry
- Response Compression
- Response Caching
- HTTP/2
- HTTP/3
- Kubernetes Service Discovery
- Consul Service Discovery
- Canary Routing
- Blue-Green Deployments
- Plugin System

---

# 💡 Why DUG?

DUG started as a way to learn Go by building production-inspired backend infrastructure instead of another CRUD application.

Every major component is implemented from scratch to understand:

- Reverse proxies
- Load balancing algorithms
- Circuit breakers
- Retry strategies
- Health checking
- Service discovery
- Concurrent programming
- High-performance HTTP servers
- Production networking patterns

The name is inspired by **Dug 🐶**, the lovable golden retriever from Pixar's *Up*.

While DUG began as a learning project, it has evolved into a fully functional API Gateway that continues to grow with production-inspired features.

---

# 🤝 Contributing

Contributions, bug reports, ideas, and discussions are always welcome.

If you'd like to contribute:

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Open a Pull Request

---

# ⭐ Support

If you found DUG useful or interesting, consider giving the repository a **Star**.

It helps others discover the project and motivates future development.

---

## 📄 License

This project is licensed under the MIT License. See the `LICENSE` file for details.
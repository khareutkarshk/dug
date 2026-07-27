# DUG (Distributed Unified Gateway)

[![CI](https://github.com/khareutkarshk/dug/actions/workflows/ci.yml/badge.svg)](https://github.com/khareutkarshk/dug/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/khareutkarshk/dug)](https://goreportcard.com/report/github.com/khareutkarshk/dug)
[![Go Reference](https://pkg.go.dev/badge/github.com/khareutkarshk/dug.svg)](https://pkg.go.dev/github.com/khareutkarshk/dug)
[![Release](https://img.shields.io/github/v/release/khareutkarshk/dug)](https://github.com/khareutkarshk/dug/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/khareutkarshk/dug)](https://github.com/khareutkarshk/dug)
[![License](https://img.shields.io/github/license/khareutkarshk/dug)](LICENSE)
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

## Clone the repository

```bash
git clone https://github.com/khareutkarshk/dug.git

cd dug
```

## Install dependencies

```bash
go mod tidy
```

## Generate local TLS certificates

```bash
mkdir certs

openssl req -x509 \
-newkey rsa:4096 \
-keyout certs/server.key \
-out certs/server.crt \
-days 365 \
-nodes
```

## Start the gateway

```bash
go run ./cmd/edge
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
│   └── edge/
├── configs/
├── examples/
│   ├── backend/
│   ├── backend2/
│   └── websocket/
├── internal/
│   ├── app/
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
│   └── upstream/
├── test/
└── README.md
```

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

- [ ] `dug run`
- [ ] `dug validate`
- [ ] `dug version`

### Distribution

- [ ] Docker Image
- [ ] GitHub Releases
- [ ] Pre-built binaries
- [ ] `go install` support

### Documentation

- [ ] Installation Guide
- [ ] Configuration Reference
- [ ] Feature Documentation
- [ ] Example Configurations

### CI/CD

- [ ] GitHub Actions
- [ ] Go Race Tests
- [ ] golangci-lint
- [ ] Automated Releases

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
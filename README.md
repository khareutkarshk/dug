# DUG (Distributed Unified Gateway)

A lightweight API Gateway built from scratch in Go to understand how production gateways work under the hood.

---

## About

DUG (Distributed Unified Gateway) is an API Gateway built from scratch to learn Go by solving real backend engineering problems instead of building another CRUD application.

The name is inspired by **Dug 🐶**, the lovable golden retriever from Pixar's *Up*—one of my favorite movies.

The goal of this project isn't to compete with gateways like **NGINX**, **Traefik**, **Kong**, or **Envoy**, but to understand the engineering decisions behind them by implementing their core building blocks from scratch.

---

## Features

### Traffic Management
- ✅ Reverse Proxy
- ✅ Route-based Routing
- ✅ YAML Configuration
- ✅ Graceful Shutdown

### Load Balancing
- ✅ Round Robin
- ✅ Smooth Weighted Round Robin (SWRR)
- ✅ Least Connections

### Reliability
- ✅ Retry Mechanism
- ✅ Exponential Backoff
- ✅ Active Health Checks
- ✅ Passive Health Checks
- ✅ Circuit Breaker
- ✅ Request Timeouts

### Security & Traffic Control
- ✅ HTTPS / TLS Support
- ✅ Per-IP Rate Limiting
- ✅ CORS
- ✅ Request Header Manipulation
- ✅ Response Header Manipulation

### Observability
- ✅ Prometheus Metrics
- ✅ Structured Logging (`slog`)
- ✅ Request IDs

### Configuration
- ✅ Hot Configuration Reload
- ✅ Atomic Router Reload

### Service Discovery
- ✅ Static Service Discovery
- ✅ DNS Service Discovery

---

## Architecture

```text
                      Client
                         │
                    HTTP / HTTPS
                         │
                         ▼
                ┌──────────────────┐
                │       DUG        │
                ├──────────────────┤
                │ Reverse Proxy    │
                │ Routing          │
                │ Load Balancer    │
                │ Retry            │
                │ Health Checks    │
                │ Circuit Breaker  │
                │ Rate Limiter     │
                │ Header Middleware│
                │ Service Discovery│
                │ Metrics & Logs   │
                └────────┬─────────┘
                         │
          ┌──────────────┼──────────────┐
          ▼              ▼              ▼
      Backend A      Backend B      Backend C
```

---

## Getting Started

Clone the repository:

```bash
git clone https://github.com/khareutkarshk/dug.git

cd dug
```

Install dependencies:

```bash
go mod tidy
```

Generate development TLS certificates:

```bash
mkdir certs

openssl req -x509 \
-newkey rsa:4096 \
-keyout certs/server.key \
-out certs/server.crt \
-days 365 \
-nodes
```

Start the gateway:

```bash
go run ./cmd/edge
```

---

## Configuration

DUG is configured using a YAML file.

Example features supported by configuration:

- Routes
- Load balancing strategy
- Retry policy
- Timeouts
- TLS
- CORS
- Request & Response Headers
- Rate Limiting
- Service Discovery

Configuration lives in:

```text
configs/edge.yaml
```

---

## Project Structure

```text
cmd/
configs/
internal/
    config/
    discovery/
    health/
    middleware/
    proxy/
    router/
    server/
    upstream/
examples/
```

---

## Roadmap

### v0.1.0
- ✅ HTTPS / TLS
- 🚧 WebSocket Support
- 🚧 Integration Tests
- 🚧 Docker Image
- 🚧 GitHub Actions
- 🚧 Documentation Improvements

### Future
- JWT Authentication
- API Keys
- OpenTelemetry
- Response Caching
- HTTP/2
- HTTP/3
- Kubernetes Service Discovery
- Consul Service Discovery
- Canary Routing
- Blue-Green Deployments

---

## Why this project?

I'm documenting my journey of learning Go by building production-inspired backend infrastructure in public.

Rather than relying on existing libraries, every feature is implemented from first principles to understand:

- how gateways route traffic
- how retries and circuit breakers improve reliability
- how load balancing algorithms work
- how service discovery integrates with routing
- how production-ready networking software is designed

If you have suggestions, feedback, or ideas for improving DUG, I'd love to hear them.

⭐ If you find the project interesting, consider giving it a star!
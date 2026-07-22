
---

## About

DUG is an API Gateway built from scratch to learn Go by solving real backend engineering problems instead of building another CRUD application.

The name is inspired by **Dug 🐶**, the lovable golden retriever from Pixar's *Up*—one of my favorite movies.

Behind the fun name is a serious goal: understanding how production gateways like **NGINX**, **Traefik**, **Kong**, and **Envoy** work by implementing their core building blocks from scratch.

---

## Features

- ✅ Reverse Proxy
- ✅ Route-based Routing
- ✅ YAML Configuration
- ✅ Round Robin Load Balancing
- 🚧 Weighted Load Balancing
- ✅ Active Health Checks
- ✅ Passive Health Checks
- ✅ Retry Mechanism
- ✅ Exponential Backoff
- ✅ Circuit Breaker
- ✅ Per-IP Rate Limiting
- ✅ Prometheus Metrics
- ✅ Structured Logging
- ✅ Graceful Shutdown

---

## Architecture

```text
               Client
                  │
                  ▼
          ┌──────────────┐
          │     DUG      │
          ├──────────────┤
          │ Reverse Proxy│
          │ Load Balancer│
          │ Health Check │
          │ Retry        │
          │ Circuit Breaker
          │ Rate Limiter │
          │ Metrics      │
          └──────┬───────┘
                 │
      ┌──────────┼──────────┐
      ▼          ▼          ▼
  Backend A  Backend B  Backend C
```

---

## Running

```bash
git clone https://github.com/khareutkarshk/dug.git

cd dug

go mod tidy

go run ./cmd/edge
```

---

## Roadmap

- Smooth Weighted Round Robin
- Least Connections
- JWT Authentication
- OpenTelemetry
- Service Discovery
- HTTP/2 & HTTP/3
- Hot Configuration Reload

---

## Why this project?

I'm documenting my journey of learning Go by building an API Gateway in public.

Every feature is implemented to understand **how** it works—not just to use a library.

If you have suggestions or feedback, I'd love to hear them.

⭐ If you find the project interesting, consider giving it a star.
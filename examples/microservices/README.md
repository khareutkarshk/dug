# Microservices Path Routing

DUG routes by URL path to dedicated services.

```text
                   ┌─────────────┐
              /users/───────────▶│ users:3001  │
                   ├─────────────┤
Client ──▶ DUG ─── /orders/────▶│ orders:3002 │
           :8080   ├─────────────┤
              /payments/────────▶│ payments:3003│
                   └─────────────┘
```

Each service is the shared [`../echo`](../echo) binary with a different `SERVICE_NAME`.

## Architecture notes

- Route patterns end with `/` so Go's ServeMux matches the path subtree (`/users/`, `/users/42`, …).
- The full path is forwarded to the upstream (no strip-prefix). Services echo the path they received.
- Adding a service = add a container + a route block in `edge.yaml`.

## Run

```bash
docker compose up --build
```

## Try it

```bash
curl -s http://localhost:8080/users/
```

Expected:

```json
{"message":"hello from users","service":"users","version":"v1","path":"/users/","method":"GET","time":"..."}
```

```bash
curl -s http://localhost:8080/orders/
curl -s http://localhost:8080/payments/checkout
```

Expected `service` fields: `orders`, `payments`. Paths echo what you requested.

Confirm routing headers:

```bash
curl -si http://localhost:8080/orders/ | grep -i x-route
# X-Route: orders
```

## Local (without Docker)

```bash
# Terminals for each service
PORT=3001 SERVICE_NAME=users    go run ./examples/echo
PORT=3002 SERVICE_NAME=orders   go run ./examples/echo
PORT=3003 SERVICE_NAME=payments go run ./examples/echo

# Point edge.yaml upstreams at localhost:3001/3002/3003, then:
dug run -config examples/microservices/edge.yaml
```

## Files

| File | Purpose |
|---|---|
| `edge.yaml` | Path-based routes |
| `docker-compose.yml` | Gateway + three services |
| `README.md` | This guide |

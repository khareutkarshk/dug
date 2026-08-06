# Docker Compose Load Balancing

DUG sits in front of two backends and balances traffic with **least connections**.

```text
                 ┌────────────┐
            ┌───▶│ backend1   │
Client ───▶│DUG │   :3001    │
           │:8080├────────────┤
            └───▶│ backend2   │
                 │   :3002    │
                 └────────────┘
```

## Run

From this directory:

```bash
docker compose up --build
```

First start may take a minute while the Go image compiles the backends.

## Try it

```bash
# Hit the gateway repeatedly — responses alternate across backends
for i in 1 2 3 4 5 6; do
  curl -s http://localhost:8080/hello
  echo
done
```

Expected (order may vary):

```json
{"message":"Hello from backend","service":"backend-service"}
{"message":"Hello from backend","service":"backend-service-2"}
...
```

Inspect which backend answered:

```bash
curl -si http://localhost:8080/hello | grep -i x-backend
# X-Backend: 3001   or   X-Backend: 3002
```

Gateway identity header:

```bash
curl -si http://localhost:8080/hello | grep -i x-powered-by
# X-Powered-By: DUG
```

## Files

| File | Purpose |
|---|---|
| `docker-compose.yml` | DUG + backend1 + backend2 |
| `edge.yaml` | Least-connections dual upstream |
| `README.md` | This guide |

Backends are the existing [`../backend`](../backend) and [`../backend2`](../backend2) packages.

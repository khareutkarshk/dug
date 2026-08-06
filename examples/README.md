# DUG Examples

Runnable, self-contained demos that show how to put DUG in front of real traffic shapes.

| Example | What it shows |
|---|---|
| [basic-reverse-proxy](./basic-reverse-proxy) | Single upstream reverse proxy |
| [docker-compose](./docker-compose) | Dual-backend load balancing in Compose |
| [microservices](./microservices) | Path routing to users / orders / payments |
| [nextjs-api](./nextjs-api) | Browser → DUG → APIs (Next.js) |
| [ai-gateway](./ai-gateway) | `/openai` and `/ollama` provider routes |
| [blue-green](./blue-green) | Weighted traffic shift between versions |

## Shared building blocks

| Path | Role |
|---|---|
| [backend](./backend) | Original demo upstream on `:3001` |
| [backend2](./backend2) | Second demo upstream on `:3002` |
| [echo](./echo) | Configurable JSON upstream (`PORT`, `SERVICE_NAME`, `VERSION`) |
| [websocket](./websocket) | WebSocket upstream demo |

Prefer `echo` (or the existing backends) in new examples instead of copying HTTP servers.

## Quick start pattern

Most Compose examples:

```bash
cd examples/<name>
docker compose up --build
```

Local gateway against a hand-started upstream:

```bash
go run ./examples/backend
dug run -config examples/basic-reverse-proxy/edge.yaml
curl -s http://localhost:8080/hello
```

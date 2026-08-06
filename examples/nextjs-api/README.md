# Next.js + DUG API Gateway

The browser only knows about DUG. Backend instances stay private; the gateway exposes `/api/` with CORS for the Next.js origin.

```text
┌─────────────┐     GET /api/catalog      ┌─────┐
│  Next.js    │ ─────────────────────────▶│ DUG │
│  :3000      │                           │:8080│
└─────────────┘                           └──┬──┘
                                             │
                              ┌──────────────┼──────────────┐
                              ▼                             ▼
                        ┌──────────┐                 ┌──────────┐
                        │ api1:3001│                 │ api2:3002│
                        └──────────┘                 └──────────┘
```

## Why DUG is in the middle

- Single public origin for APIs (`localhost:8080`)
- Load balancing across API replicas
- CORS, headers, and timeouts configured once in `edge.yaml`
- Frontend config is just `NEXT_PUBLIC_API_BASE`

## Run (Docker — recommended)

From this directory:

```bash
docker compose up --build
```

Open [http://localhost:3000](http://localhost:3000) and click **Fetch via DUG**.

## Run (local processes)

```bash
# Terminal 1–2 — APIs (shared echo service)
PORT=3001 SERVICE_NAME=catalog-api VERSION=a go run ./examples/echo
PORT=3002 SERVICE_NAME=catalog-api VERSION=b go run ./examples/echo

# Terminal 3 — gateway (uses localhost upstreams)
dug run -config examples/nextjs-api/edge.yaml

# Terminal 4 — frontend
cd examples/nextjs-api/frontend
npm install
NEXT_PUBLIC_API_BASE=http://localhost:8080 npm run dev
```

## Try it with curl

```bash
curl -si http://localhost:8080/api/catalog | head -n 20
```

Expected body (version `a` or `b`):

```json
{
  "message": "hello from catalog-api",
  "service": "catalog-api",
  "version": "a",
  "path": "/api/catalog",
  "method": "GET",
  "time": "..."
}
```

Headers should include `X-Powered-By: DUG`.

CORS preflight from the browser origin:

```bash
curl -si -X OPTIONS http://localhost:8080/api/catalog \
  -H 'Origin: http://localhost:3000' \
  -H 'Access-Control-Request-Method: GET' | grep -i access-control
```

## Files

| File | Purpose |
|---|---|
| `edge.yaml` | Local gateway config (localhost upstreams) |
| `edge.docker.yaml` | Compose network upstreams |
| `docker-compose.yml` | DUG + APIs + Next.js |
| `frontend/` | Minimal Next.js app |
| `README.md` | This guide |

APIs reuse [`../echo`](../echo).

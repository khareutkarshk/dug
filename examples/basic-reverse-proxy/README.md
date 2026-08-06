# Basic Reverse Proxy

Minimal DUG setup: one gateway, one upstream.

```text
Client
  │
  │  GET /
  ▼
┌─────┐     ┌──────────┐
│ DUG │────▶│ backend  │
│:8080│     │  :3001   │
└─────┘     └──────────┘
```

## Prerequisites

- Go 1.22+
- DUG installed (`go install github.com/khareutkarshk/dug/cmd/dug@latest`) or built from this repo

## Run

From the **repository root**:

```bash
# Terminal 1 — upstream (reuses examples/backend)
go run ./examples/backend

# Terminal 2 — gateway
dug run -config examples/basic-reverse-proxy/edge.yaml
# or: go run ./cmd/dug run -config examples/basic-reverse-proxy/edge.yaml
```

## Try it

```bash
curl -s http://localhost:8080/hello
```

Expected:

```json
{"message":"Hello from backend","service":"backend-service"}
```

Health via the upstream (DUG proxies `/` — backend exposes `/health` directly):

```bash
curl -s http://localhost:3001/health
# OK
```

Through the gateway (path is forwarded as-is):

```bash
curl -si http://localhost:8080/hello | head -n 15
```

Look for `X-Powered-By: DUG` on the response.

## Files

| File | Purpose |
|---|---|
| `edge.yaml` | Minimal single-upstream config |
| `README.md` | This guide |

Upstream code lives in [`../backend`](../backend) so it is not duplicated.

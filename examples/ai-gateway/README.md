# AI Gateway

DUG fronts multiple model providers behind stable paths. Mocks stand in for OpenAI and Ollama so the example runs offline.

```text
                 ┌──────────────────┐
           /openai/────────────────▶│ openai-mock:3001 │
Client ──▶ DUG   ├──────────────────┤
           :8080 │                  │
           /ollama/────────────────▶│ ollama-mock:3002 │
                 └──────────────────┘
```

Swap a mock URL for a real provider endpoint when you move beyond local demos (same route shape).

## Run

```bash
docker compose up --build
```

## Try it

### OpenAI-style path

```bash
curl -s http://localhost:8080/openai/v1/chat/completions \
  -H 'Content-Type: application/json' \
  -d '{"model":"gpt-mock","messages":[{"role":"user","content":"hi"}]}'
```

Expected (trimmed):

```json
{
  "id": "chatcmpl-dug-mock",
  "object": "chat.completion",
  "model": "gpt-mock",
  "choices": [
    {
      "message": {
        "role": "assistant",
        "content": "Hello from the OpenAI mock behind DUG"
      }
    }
  ],
  "service": "openai-mock",
  "path": "/openai/v1/chat/completions"
}
```

### Ollama-style path

```bash
curl -s http://localhost:8080/ollama/api/chat \
  -H 'Content-Type: application/json' \
  -d '{"model":"llama-mock","messages":[{"role":"user","content":"hi"}]}'
```

Expected (trimmed):

```json
{
  "model": "llama-mock",
  "message": {
    "role": "assistant",
    "content": "Hello from the Ollama mock behind DUG"
  },
  "done": true,
  "service": "ollama-mock",
  "path": "/ollama/api/chat"
}
```

Provider header:

```bash
curl -si http://localhost:8080/openai/ | grep -i x-ai-provider
# X-AI-Provider: openai
```

## Files

| File | Purpose |
|---|---|
| `edge.yaml` | `/openai/` and `/ollama/` routes |
| `docker-compose.yml` | Gateway + two mocks |
| `README.md` | This guide |

Mocks are [`../echo`](../echo) with `SERVICE_NAME=openai-mock` / `ollama-mock`.

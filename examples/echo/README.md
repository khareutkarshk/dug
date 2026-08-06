# Echo upstream

Tiny JSON service used by DUG examples so we do not copy HTTP servers around.

```bash
PORT=3001 SERVICE_NAME=users VERSION=v1 go run ./examples/echo
curl -s http://localhost:3001/
```

| Env | Default | Meaning |
|---|---|---|
| `PORT` | `3000` | Listen port |
| `SERVICE_NAME` | `echo` | Name embedded in JSON / special mock mode |
| `VERSION` | `v1` | Version label (blue/green demos) |

Special `SERVICE_NAME` values:

- `openai-mock` — OpenAI-shaped chat completion JSON
- `ollama-mock` — Ollama-shaped chat JSON

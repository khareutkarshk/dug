# Blue / Green Weighted Traffic

Two versions of the same API sit behind DUG. **Smooth weighted round robin** sends most traffic to blue and a slice to green.

```text
                 weight:9
            ┌───────────────▶ blue  (v1) :3001
Client ───▶│DUG
           │:8080
            └───────────────▶ green (v2) :3002
                 weight:1
```

Shift traffic by editing weights in `edge.yaml` (DUG hot-reloads the config file).

## Run

```bash
docker compose up --build
```

## Observe the split (~90% / 10%)

```bash
for i in $(seq 1 20); do
  curl -s http://localhost:8080/ | grep -o '"version":"[^"]*"'
done | sort | uniq -c
```

Expected shape (counts will vary slightly):

```text
 18 "version":"blue-v1"
  2 "version":"green-v2"
```

Single request sample:

```bash
curl -s http://localhost:8080/
```

```json
{"message":"hello from shop-api","service":"shop-api","version":"blue-v1","path":"/","method":"GET","time":"..."}
```

## Shift traffic to green

Edit `edge.yaml`:

```yaml
upstreams:
  - url: http://blue:3001
    weight: 1
  - url: http://green:3002
    weight: 9
```

Save the file — DUG reloads automatically. Re-run the loop; most responses should show `green-v2`.

Cut over fully:

```yaml
upstreams:
  - url: http://green:3002
    weight: 1
```

(Remove the blue upstream when you are ready to drain it.)

## Files

| File | Purpose |
|---|---|
| `edge.yaml` | Weighted blue/green upstreams |
| `docker-compose.yml` | Gateway + blue + green |
| `README.md` | This guide |

Both colors use [`../echo`](../echo) with different `VERSION` values.

# {{.Name}}

Bootstrapped with [DUG](https://github.com/khareutkarshk/dug) (`dug init`).

## Quick start

```bash
cd {{.Name}}
dug run -config configs/edge.yaml
```

Validate configuration:

```bash
dug validate -config configs/edge.yaml
```

Run diagnostics:

```bash
dug doctor -config configs/edge.yaml
```

## Docker

```bash
docker compose up
```

This starts DUG from `ghcr.io/khareutkarshk/dug:latest` and a sample upstream on port `3001`.

## Layout

```text
.
├── configs/
│   └── edge.yaml
├── certs/
├── docker-compose.yml
└── README.md
```

Point `routes[].upstreams` in `configs/edge.yaml` at your real backends when you are ready.

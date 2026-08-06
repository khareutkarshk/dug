# DUG Gateway Benchmarks

Generated: 2026-08-06T05:33:35Z

Deterministic httptest upstreams (instant `200 OK`). Rate limiting effectively disabled.

## Results

| Strategy | Concurrency | Requests | RPS | Avg | P50 | P95 | P99 | Failures |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `smooth_weighted` | 1 | 5 | 1841.1 | 540.1µs | 378.3µs | 1.32ms | 1.32ms | 0 |
| `least_connections` | 1 | 5 | 2047.6 | 486.1µs | 453.4µs | 686.8µs | 686.8µs | 0 |

## Notes

- Strategies: Smooth Weighted Round Robin (`smooth_weighted`), Least Connections (`least_connections`)
- Concurrency levels: 1, 10, 100, 500 (override with `BENCH_CONCURRENCY`)
- Requests per worker: default 100 (override with `BENCH_REQUESTS_PER_WORKER`)
- Run with: `make benchmark`

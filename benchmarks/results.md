# DUG Gateway Benchmarks

Generated: 2026-08-06T05:13:11Z

Deterministic httptest upstreams (instant `200 OK`). Rate limiting effectively disabled.

## Results

| Strategy | Concurrency | Requests | RPS | Avg | P50 | P95 | P99 | Failures |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `smooth_weighted` | 1 | 100 | 2794.5 | 357.7µs | 300.3µs | 632.9µs | 1.05ms | 0 |
| `smooth_weighted` | 10 | 1000 | 5266.7 | 1.86ms | 1.35ms | 5.10ms | 8.92ms | 0 |
| `smooth_weighted` | 100 | 10000 | 12881.3 | 7.41ms | 6.23ms | 16.77ms | 30.80ms | 0 |
| `smooth_weighted` | 500 | 50000 | 13755.0 | 35.24ms | 32.44ms | 70.32ms | 96.20ms | 0 |
| `least_connections` | 1 | 100 | 3141.0 | 318.2µs | 259.4µs | 559.3µs | 1.10ms | 0 |
| `least_connections` | 10 | 1000 | 10715.7 | 913.2µs | 548.3µs | 2.24ms | 2.70ms | 0 |
| `least_connections` | 100 | 10000 | 12878.5 | 7.45ms | 6.64ms | 17.54ms | 23.40ms | 0 |
| `least_connections` | 500 | 50000 | 14330.5 | 33.79ms | 30.77ms | 70.89ms | 93.60ms | 0 |

## Notes

- Strategies: Smooth Weighted Round Robin (`smooth_weighted`), Least Connections (`least_connections`)
- Concurrency levels: 1, 10, 100, 500 (override with `BENCH_CONCURRENCY`)
- Requests per worker: default 100 (override with `BENCH_REQUESTS_PER_WORKER`)
- Run with: `make benchmark`

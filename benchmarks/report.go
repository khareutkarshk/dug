//go:build benchmark

package benchmarks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// WriteReport writes a markdown performance report to path.
func WriteReport(path string, results []RunResult) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	var b strings.Builder
	b.WriteString("# DUG Gateway Benchmarks\n\n")
	b.WriteString(fmt.Sprintf("Generated: %s\n\n", time.Now().UTC().Format(time.RFC3339)))
	b.WriteString("Deterministic httptest upstreams (instant `200 OK`). Rate limiting effectively disabled.\n\n")

	b.WriteString("## Results\n\n")
	b.WriteString("| Strategy | Concurrency | Requests | RPS | Avg | P50 | P95 | P99 | Failures |\n")
	b.WriteString("|---|---:|---:|---:|---:|---:|---:|---:|---:|\n")

	for _, r := range results {
		b.WriteString(fmt.Sprintf(
			"| `%s` | %d | %d | %.1f | %s | %s | %s | %s | %d |\n",
			r.Strategy,
			r.Concurrency,
			r.TotalRequests,
			r.RequestsPerSec,
			formatLatency(r.AvgLatency),
			formatLatency(r.P50),
			formatLatency(r.P95),
			formatLatency(r.P99),
			r.Failures,
		))
	}

	b.WriteString("\n## Notes\n\n")
	b.WriteString("- Strategies: Smooth Weighted Round Robin (`smooth_weighted`), Least Connections (`least_connections`)\n")
	b.WriteString("- Concurrency levels: 1, 10, 100, 500 (override with `BENCH_CONCURRENCY`)\n")
	b.WriteString("- Requests per worker: default 100 (override with `BENCH_REQUESTS_PER_WORKER`)\n")
	b.WriteString("- Run with: `make benchmark`\n")

	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func formatLatency(d time.Duration) string {
	if d < time.Microsecond {
		return fmt.Sprintf("%dns", d.Nanoseconds())
	}
	if d < time.Millisecond {
		return fmt.Sprintf("%.1fµs", float64(d.Nanoseconds())/1e3)
	}
	if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d.Microseconds())/1e3)
	}
	return d.Round(time.Millisecond).String()
}

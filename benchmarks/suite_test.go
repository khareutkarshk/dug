//go:build benchmark

package benchmarks

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestGatewayBenchmarks(t *testing.T) {
	concurrencies := parseIntList(
		os.Getenv("BENCH_CONCURRENCY"),
		[]int{1, 10, 100, 500},
	)
	requestsPerWorker := parsePositiveInt(
		os.Getenv("BENCH_REQUESTS_PER_WORKER"),
		100,
	)

	results := make([]RunResult, 0, len(SupportedStrategies())*len(concurrencies))

	for _, strategy := range SupportedStrategies() {
		harness, err := StartHarness(strategy, 3)
		if err != nil {
			t.Fatalf("start harness (%s): %v", strategy, err)
		}

		for _, concurrency := range concurrencies {
			t.Logf("running strategy=%s concurrency=%d requests_per_worker=%d",
				strategy, concurrency, requestsPerWorker)

			result := RunLoad(harness.GatewayURL, RunConfig{
				Strategy:          strategy,
				Concurrency:       concurrency,
				RequestsPerWorker: requestsPerWorker,
			})
			results = append(results, result)

			t.Logf(
				"  rps=%.1f avg=%s p50=%s p95=%s p99=%s failures=%d",
				result.RequestsPerSec,
				formatLatency(result.AvgLatency),
				formatLatency(result.P50),
				formatLatency(result.P95),
				formatLatency(result.P99),
				result.Failures,
			)

			if result.Failures > 0 {
				t.Errorf("%s concurrency=%d: %d/%d requests failed",
					strategy, concurrency, result.Failures, result.TotalRequests)
			}
		}

		if err := harness.Close(); err != nil {
			t.Errorf("close harness (%s): %v", strategy, err)
		}
	}

	reportPath := filepath.Join(".", "results.md")
	if err := WriteReport(reportPath, results); err != nil {
		t.Fatalf("write report: %v", err)
	}
	t.Logf("wrote report to %s", reportPath)
}

func parseIntList(raw string, fallback []int) []int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fallback
	}

	parts := strings.Split(raw, ",")
	out := make([]int, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		n, err := strconv.Atoi(part)
		if err != nil || n < 1 {
			continue
		}
		out = append(out, n)
	}
	if len(out) == 0 {
		return fallback
	}
	return out
}

func parsePositiveInt(raw string, fallback int) int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		return fallback
	}
	return n
}

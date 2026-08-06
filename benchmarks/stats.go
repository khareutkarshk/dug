//go:build benchmark

package benchmarks

import (
	"math"
	"sort"
	"time"
)

type latencyStats struct {
	Avg time.Duration
	P50 time.Duration
	P95 time.Duration
	P99 time.Duration
}

func summarize(latencies []time.Duration) latencyStats {
	if len(latencies) == 0 {
		return latencyStats{}
	}

	sorted := append([]time.Duration(nil), latencies...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	var sum time.Duration
	for _, d := range sorted {
		sum += d
	}

	return latencyStats{
		Avg: sum / time.Duration(len(sorted)),
		P50: percentile(sorted, 50),
		P95: percentile(sorted, 95),
		P99: percentile(sorted, 99),
	}
}

// percentile uses the nearest-rank method on a sorted slice.
func percentile(sorted []time.Duration, p float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	if p <= 0 {
		return sorted[0]
	}
	if p >= 100 {
		return sorted[len(sorted)-1]
	}

	rank := int(math.Ceil(p/100*float64(len(sorted)))) - 1
	if rank < 0 {
		rank = 0
	}
	if rank >= len(sorted) {
		rank = len(sorted) - 1
	}
	return sorted[rank]
}

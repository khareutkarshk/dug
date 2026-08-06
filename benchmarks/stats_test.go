//go:build benchmark

package benchmarks

import (
	"testing"
	"time"
)

func TestPercentile(t *testing.T) {
	sorted := []time.Duration{
		1 * time.Millisecond,
		2 * time.Millisecond,
		3 * time.Millisecond,
		4 * time.Millisecond,
		5 * time.Millisecond,
	}

	if got := percentile(sorted, 50); got != 3*time.Millisecond {
		t.Fatalf("p50=%v want 3ms", got)
	}
	if got := percentile(sorted, 100); got != 5*time.Millisecond {
		t.Fatalf("p100=%v want 5ms", got)
	}
}

func TestSummarizeEmpty(t *testing.T) {
	stats := summarize(nil)
	if stats.Avg != 0 || stats.P99 != 0 {
		t.Fatalf("expected zero stats, got %+v", stats)
	}
}

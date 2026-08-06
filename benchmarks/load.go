//go:build benchmark

package benchmarks

import (
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// RunConfig controls a single load run against the gateway.
type RunConfig struct {
	Strategy          string
	Concurrency       int
	RequestsPerWorker int
}

// RunResult holds aggregate metrics for one load run.
type RunResult struct {
	Strategy       string
	Concurrency    int
	TotalRequests  int
	Successes      int
	Failures       int
	Duration       time.Duration
	RequestsPerSec float64
	AvgLatency     time.Duration
	P50            time.Duration
	P95            time.Duration
	P99            time.Duration
}

// RunLoad issues Concurrent * RequestsPerWorker GETs against targetURL and
// records per-request latency.
func RunLoad(targetURL string, cfg RunConfig) RunResult {
	if cfg.Concurrency < 1 {
		cfg.Concurrency = 1
	}
	if cfg.RequestsPerWorker < 1 {
		cfg.RequestsPerWorker = 1
	}

	total := cfg.Concurrency * cfg.RequestsPerWorker
	latencies := make([]time.Duration, total)

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        cfg.Concurrency * 2,
			MaxIdleConnsPerHost: cfg.Concurrency * 2,
			MaxConnsPerHost:     0,
			IdleConnTimeout:     30 * time.Second,
			DisableCompression:  true,
		},
	}

	var (
		wg        sync.WaitGroup
		successes atomic.Int64
		failures  atomic.Int64
		nextSlot  atomic.Int64
	)

	start := time.Now()
	wg.Add(cfg.Concurrency)

	for worker := 0; worker < cfg.Concurrency; worker++ {
		go func() {
			defer wg.Done()

			for i := 0; i < cfg.RequestsPerWorker; i++ {
				reqStart := time.Now()
				ok := doRequest(client, targetURL)
				latency := time.Since(reqStart)

				slot := int(nextSlot.Add(1) - 1)
				latencies[slot] = latency

				if ok {
					successes.Add(1)
				} else {
					failures.Add(1)
				}
			}
		}()
	}

	wg.Wait()
	duration := time.Since(start)

	stats := summarize(latencies)
	rps := 0.0
	if duration > 0 {
		rps = float64(total) / duration.Seconds()
	}

	return RunResult{
		Strategy:       cfg.Strategy,
		Concurrency:    cfg.Concurrency,
		TotalRequests:  total,
		Successes:      int(successes.Load()),
		Failures:       int(failures.Load()),
		Duration:       duration,
		RequestsPerSec: rps,
		AvgLatency:     stats.Avg,
		P50:            stats.P50,
		P95:            stats.P95,
		P99:            stats.P99,
	}
}

func doRequest(client *http.Client, url string) bool {
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

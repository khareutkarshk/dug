package metrics

import "github.com/prometheus/client_golang/prometheus"

var RequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "dug_requests_total",
		Help: "Total number of requests handled by the gateway.",
	}, []string{"method", "path", "status"},
)

var RequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "dug_request_duration_seconds",
		Help:    "Histogram of request durations.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"},
)

var InFlightRequests = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "dug_inflight_requests",
		Help: "Current number of in-flight requests.",
	},
)

func Register() {
	prometheus.MustRegister(
		RequestsTotal,
		RequestDuration,
		InFlightRequests,
	)
}

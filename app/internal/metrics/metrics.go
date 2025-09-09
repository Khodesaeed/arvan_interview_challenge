package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// RequestTotal is a Counter that tracks total HTTP requests.
	RequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"path", "status"},
	)

	// RequestLatency is a Histogram that tracks request latency.
	RequestLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_latency_seconds",
			Help:    "Request latency in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method", "status"},
	)

	// InflightRequests is a Gauge that tracks in-progress requests.
	InflightRequests = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "inprogress_requests",
			Help: "Total number of requests in progress.",
		},
		[]string{"path", "method"},
	)
)

func init() {
	println("metrics package initialized")
}
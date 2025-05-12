package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"sync"
)

var (
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "zlib",
			Name:      "http_request_duration_seconds",
			Help:      "Duration of HTTP requests in seconds",
			Buckets:   prometheus.DefBuckets,
		}, []string{"path", "method", "status"},
	)

	ActiveIPs = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "zlib",
			Name:      "active_ips_total",
			Help:      "Current number of active IPs",
		},
		[]string{"ip"},
	)

	activeIPsMap = &sync.Map{}

	TotalRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "zlib",
			Name:      "total_requests",
			Help:      "Total number of HTTP requests",
		}, []string{"path", "method"},
	)

	ResponseStatus = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "zlib",
			Name:      "response_status",
			Help:      "HTTP response status codes",
		}, []string{"status"},
	)
)

func RecordIPActivity(ip string) {
	activeIPsMap.Store(ip, true)
	ActiveIPs.WithLabelValues(ip).Set(1)
}

func InitializeIPMetric(ip string) {
	ActiveIPs.WithLabelValues(ip).Set(0)
}

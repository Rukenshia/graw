package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ResponseSize = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "response_size",
		Namespace: "graw",
		Help:      "Body response size (bytes) for API requests",

		Buckets: []float64{32, 256, 1024, 2048, 4098},
	}, []string{"path", "status"})

	Requests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "requests",
		Namespace: "graw",
		Help:      "API requests",
	}, []string{"path", "status"})

	RateLimitUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "rate_limit_used",
		Namespace: "graw",
		Help:      "Rate limit used",
	})

	RateLimitRemaining = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "rate_limit_remaining",
		Namespace: "graw",
		Help:      "Rate limit remaining",
	})

	RateLimitReset = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "rate_limit_reset",
		Namespace: "graw",
		Help:      "Rate limit reset",
	})
)

// Register registers all metrics
func Register() {
	prometheus.MustRegister(ResponseSize)
	prometheus.MustRegister(Requests)
	prometheus.MustRegister(RateLimitUsed)
	prometheus.MustRegister(RateLimitRemaining)
	prometheus.MustRegister(RateLimitReset)
}

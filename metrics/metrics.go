package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ResponseSize = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:      "response_size",
		Namespace: "graw",
		Help:      "Body response size (bytes) for API requests",

		Buckets: []float64{32, 256, 1024, 2048, 4098},
	})

	Requests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "requests",
		Namespace: "graw",
		Help:      "API requests",
	}, []string{"status"})
)

// Register registers all metrics
func Register() {
	prometheus.MustRegister(ResponseSize)
	prometheus.MustRegister(Requests)
}

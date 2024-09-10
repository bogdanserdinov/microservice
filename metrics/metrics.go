package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var once sync.Once

type Metrics struct {
	ApiRequests *prometheus.HistogramVec
}

func New(exporter promauto.Factory) *Metrics {
	metricsInstance := new(Metrics)
	once.Do(func() {
		requests := exporter.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "api_requests",
			Buckets: []float64{0.01, 0.03, 0.05, 0.1, 0.2, 0.3, 0.5, 1, 1.5, 2, 2.5, 3},
		}, []string{})

		metricsInstance = &Metrics{
			ApiRequests: requests,
		}
	})

	return metricsInstance
}

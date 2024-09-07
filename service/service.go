package service

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type DB interface {
	CreateDummy(ctx context.Context, dummy Dummy) error
}

type Service struct {
	db DB

	requests *prometheus.HistogramVec
}

func New(db DB, exporter promauto.Factory) *Service {
	requests := exporter.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "api_requests",
		Buckets: []float64{0.05, 0.1, 0.2, 0.3, 0.5, 1, 1.5, 2, 2.5, 3},
	}, []string{})

	return &Service{
		db:       db,
		requests: requests,
	}
}

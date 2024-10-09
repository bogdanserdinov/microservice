package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type DB interface {
	CreateDummy(ctx context.Context, dummy Dummy) error
}

type Service struct {
	tracer trace.Tracer

	db DB
}

func New(tracer trace.Tracer, db DB) *Service {
	return &Service{
		tracer: tracer,
		db:     db,
	}
}

func (service *Service) Create(ctx context.Context, status Status, description string) error {
	ctx, span := service.tracer.Start(ctx, "create_service")
	defer span.End()

	dummy := Dummy{
		ID:          uuid.New(),
		Status:      status,
		Description: description,
		UpdatedAt:   time.Now().UTC(),
		CreatedAt:   time.Now().UTC(),
	}

	return service.db.CreateDummy(ctx, dummy)
}

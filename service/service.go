package service

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type DB interface {
	CreateDummy(ctx context.Context, dummy Dummy) error
}

type Service struct {
	db DB
}

func New(db DB) *Service {
	return &Service{
		db: db,
	}
}

func (service *Service) Create(ctx context.Context, status Status, description string) error {
	dummy := Dummy{
		ID:          uuid.New(),
		Status:      status,
		Description: description,
		UpdatedAt:   time.Now().UTC(),
		CreatedAt:   time.Now().UTC(),
	}

	return service.db.CreateDummy(ctx, dummy)
}

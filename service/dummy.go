package service

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Status string

const (
	StatusPending Status = "pending"
	StatusSuccess Status = "success"
	StatusFailed  Status = "failed"
)

func (s Status) IsValid() bool {
	return s == StatusPending || s == StatusSuccess || s == StatusFailed
}

type Dummy struct {
	ID          uuid.UUID
	Status      Status
	Description string
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

func (service *Service) Create(ctx context.Context, dummy Dummy) error {
	return service.db.CreateDummy(ctx, dummy)
}

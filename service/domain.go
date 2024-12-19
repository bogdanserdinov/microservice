package service

import (
	"time"

	"github.com/google/uuid"
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

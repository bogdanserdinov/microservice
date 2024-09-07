package dummy

import (
	"time"

	"github.com/google/uuid"
)

type Dummy struct {
	ID          uuid.UUID
	Status      string
	Description string
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

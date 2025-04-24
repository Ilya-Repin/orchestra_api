package model

import (
	"github.com/google/uuid"
	"time"
)

type Registration struct {
	ID        int
	UserID    uuid.UUID
	EventID   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

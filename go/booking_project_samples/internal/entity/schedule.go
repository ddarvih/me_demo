package entity

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ID         uuid.UUID
	RoomID     uuid.UUID
	DaysOfWeek []int
	StartTime  string // HH:MM UTC
	EndTime    string // HH:MM UTC
	CreatedAt  time.Time
}

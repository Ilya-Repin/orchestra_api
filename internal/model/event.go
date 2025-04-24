package model

import "time"

type Event struct {
	ID          int
	Title       string
	Description string
	EventType   EventType
	EventDate   time.Time
	Location    Location
	Capacity    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

package entity

import (
	"time"

	"github.com/google/uuid"
)

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

type Event struct {
	EventID    uuid.UUID
	Date       Date
	EventInfo  string
	CreateDate time.Time
}

type NewEvent struct {
	UserID    uuid.UUID `json:"userid"`
	Date      Date      `json:"date"`
	EventInfo string    `json:"info"`
}

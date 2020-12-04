package entity

import (
	"time"
)

// Flight represents an flight record.
type Flight struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`           // flight name
	Number        string    `json:"number"`         // flight number
	Departure     string    `json:"departure"`      // departure
	DepartureTime time.Time `json:"departure_time"` // scheduled date & time
	Destination   string    `json:"destination"`    // destination
	ArrivalTime   time.Time `json:"arrival_time"`   // expected arrival date & time
	Fare          string    `json:"fare"`           // fare
	Duration      string    `json:"duration"`       // flight duration
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

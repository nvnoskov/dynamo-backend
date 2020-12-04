package flight

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hako/durafmt"
	"github.com/nvnoskov/dynamo-backend/internal/entity"
	"github.com/nvnoskov/dynamo-backend/pkg/log"
)

// Service encapsulates usecase logic for flights.
type Service interface {
	Get(ctx context.Context, id string) (Flight, error)
	Query(ctx context.Context, input SearchFlightRequest, offset, limit int) ([]Flight, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, input CreateFlightRequest) (Flight, error)
	Update(ctx context.Context, id string, input UpdateFlightRequest) (Flight, error)
	Delete(ctx context.Context, id string) (Flight, error)
}

// Flight represents the data about an flight.
type Flight struct {
	entity.Flight
}

// CreateFlightRequest represents an flight creation request.
type CreateFlightRequest struct {
	Name          string    `json:"name"`           // flight name
	Number        string    `json:"number"`         // flight number
	Departure     string    `json:"departure"`      // departure
	DepartureTime time.Time `json:"departure_time"` // scheduled date & time
	Destination   string    `json:"destination"`    // destination
	ArrivalTime   time.Time `json:"arrival_time"`   // expected arrival date & time
	Fare          string    `json:"fare"`           // fare
	Duration      string    `json:"duration"`       // flight duration
}

// Validate validates the CreateFlightRequest fields.
func (m CreateFlightRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 128)),
		validation.Field(&m.Number, validation.Required, validation.Length(0, 20)),
		validation.Field(&m.Departure, validation.Required, validation.Length(0, 100)),
		validation.Field(&m.Destination, validation.Required, validation.Length(0, 100)),
		validation.Field(&m.Fare, validation.Required, validation.Length(0, 20)),
		validation.Field(&m.DepartureTime, validation.Required),
		validation.Field(&m.ArrivalTime, validation.Required),
	)
}

// UpdateFlightRequest represents an flight update request.
type UpdateFlightRequest struct {
	Name          string    `json:"name"`           // flight name
	Number        string    `json:"number"`         // flight number
	Departure     string    `json:"departure"`      // departure
	DepartureTime time.Time `json:"departure_time"` // scheduled date & time
	Destination   string    `json:"destination"`    // destination
	ArrivalTime   time.Time `json:"arrival_time"`   // expected arrival date & time
	Fare          string    `json:"fare"`           // fare
	Duration      string    `json:"duration"`       // flight duration
}

// Validate validates the UpdateFlightRequest fields.
func (m UpdateFlightRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 128)),
		validation.Field(&m.Number, validation.Required, validation.Length(0, 20)),
		validation.Field(&m.Departure, validation.Required, validation.Length(0, 100)),
		validation.Field(&m.Destination, validation.Required, validation.Length(0, 100)),
		validation.Field(&m.Fare, validation.Required, validation.Length(0, 20)),
		validation.Field(&m.DepartureTime, validation.Required),
		validation.Field(&m.ArrivalTime, validation.Required),
	)
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new flight service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// Get returns the flight with the specified the flight ID.
func (s service) Get(ctx context.Context, id string) (Flight, error) {
	flight, err := s.repo.Get(ctx, id)
	if err != nil {
		return Flight{}, err
	}
	return Flight{flight}, nil
}

// Create creates a new flight.
func (s service) Create(ctx context.Context, req CreateFlightRequest) (Flight, error) {
	if err := req.Validate(); err != nil {
		return Flight{}, err
	}
	id := entity.GenerateID()
	now := time.Now()

	duration := durafmt.Parse(req.ArrivalTime.Sub(req.DepartureTime)).String() // calculate flight duration
	err := s.repo.Create(ctx, entity.Flight{
		ID:            id,
		Name:          req.Name,
		Number:        req.Number,
		Departure:     req.Departure,
		DepartureTime: req.DepartureTime,
		Destination:   req.Destination,
		ArrivalTime:   req.ArrivalTime,
		Fare:          req.Fare,
		Duration:      duration,
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	if err != nil {
		return Flight{}, err
	}
	return s.Get(ctx, id)
}

// Update updates the flight with the specified ID.
func (s service) Update(ctx context.Context, id string, req UpdateFlightRequest) (Flight, error) {
	if err := req.Validate(); err != nil {
		return Flight{}, err
	}

	flight, err := s.Get(ctx, id)
	if err != nil {
		return flight, err
	}
	duration := durafmt.Parse(req.ArrivalTime.Sub(req.DepartureTime)).String() // calculate flight duration

	flight.Name = req.Name
	flight.Name = req.Name
	flight.Number = req.Number
	flight.Departure = req.Departure
	flight.DepartureTime = req.DepartureTime
	flight.Destination = req.Destination
	flight.ArrivalTime = req.ArrivalTime
	flight.Fare = req.Fare
	flight.Duration = duration

	flight.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, flight.Flight); err != nil {
		return flight, err
	}
	return flight, nil
}

// Delete deletes the flight with the specified ID.
func (s service) Delete(ctx context.Context, id string) (Flight, error) {
	flight, err := s.Get(ctx, id)
	if err != nil {
		return Flight{}, err
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		return Flight{}, err
	}
	return flight, nil
}

// Count returns the number of flights.
func (s service) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}

// SearchFlightRequest represents an flight update request.
type SearchFlightRequest struct {
	Name          string `json:"name"`           // flight name
	Departure     string `json:"departure"`      // departure
	DepartureTime string `json:"departure_time"` // scheduled date & time
	Destination   string `json:"destination"`    // destination
}

// Query returns the flights with the specified offset and limit.
func (s service) Query(ctx context.Context, req SearchFlightRequest, offset, limit int) ([]Flight, error) {
	log.New().Infof("%+v", req)

	items, err := s.repo.Query(ctx, req, offset, limit)
	if err != nil {
		return nil, err
	}
	result := []Flight{}
	for _, item := range items {
		result = append(result, Flight{item})
	}
	return result, nil
}

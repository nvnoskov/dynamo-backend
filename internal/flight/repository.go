package flight

import (
	"context"
	"time"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/nvnoskov/dynamo-backend/internal/entity"
	"github.com/nvnoskov/dynamo-backend/pkg/dbcontext"
	"github.com/nvnoskov/dynamo-backend/pkg/log"
)

// Repository encapsulates the logic to access flights from the data source.
type Repository interface {
	// Get returns the flight with the specified flight ID.
	Get(ctx context.Context, id string) (entity.Flight, error)
	// Count returns the number of flights.
	Count(ctx context.Context) (int, error)
	// Query returns the list of flights with the given offset and limit.
	Query(ctx context.Context, req SearchFlightRequest, offset, limit int) ([]entity.Flight, error)
	// Create saves a new flight in the storage.
	Create(ctx context.Context, flight entity.Flight) error
	// Update updates the flight with given ID in the storage.
	Update(ctx context.Context, flight entity.Flight) error
	// Delete removes the flight with given ID from the storage.
	Delete(ctx context.Context, id string) error
}

// repository persists flights in database
type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

// NewRepository creates a new flight repository
func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

// Get reads the flight with the specified ID from the database.
func (r repository) Get(ctx context.Context, id string) (entity.Flight, error) {
	var flight entity.Flight
	err := r.db.With(ctx).Select().Model(id, &flight)
	return flight, err
}

// Create saves a new flight record in the database.
// It returns the ID of the newly inserted flight record.
func (r repository) Create(ctx context.Context, flight entity.Flight) error {
	return r.db.With(ctx).Model(&flight).Insert()
}

// Update saves the changes to an flight in the database.
func (r repository) Update(ctx context.Context, flight entity.Flight) error {
	return r.db.With(ctx).Model(&flight).Update()
}

// Delete deletes an flight with the specified ID from the database.
func (r repository) Delete(ctx context.Context, id string) error {
	flight, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&flight).Delete()
}

// Count returns the number of the flight records in the database.
func (r repository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.With(ctx).Select("COUNT(*)").From("flight").Row(&count)
	return count, err
}

// Query retrieves the flight records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, req SearchFlightRequest, offset, limit int) ([]entity.Flight, error) {
	var flights []entity.Flight
	whereOptions := make(dbx.HashExp)
	if req.Name != "" {
		whereOptions["name"] = req.Name
	}
	if req.Departure != "" {
		whereOptions["departure"] = req.Departure
	}
	if req.DepartureTime != "" {
		departureTime, err := time.Parse("2006-01-02", req.DepartureTime)
		if err == nil {
			whereOptions["departure_time"] = dbx.Between("departure_time", departureTime, departureTime.Add(time.Hour*24)) //departureTime
		}
	}
	if req.Destination != "" {
		whereOptions["destination"] = req.Destination
	}

	err := r.db.With(ctx).
		Select().
		Where(whereOptions).
		OrderBy("id").
		Offset(int64(offset)).
		Limit(int64(limit)).
		All(&flights)
	return flights, err
}

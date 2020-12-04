package flight

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/nvnoskov/dynamo-backend/internal/entity"
	"github.com/nvnoskov/dynamo-backend/internal/test"
	"github.com/nvnoskov/dynamo-backend/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	logger, _ := log.NewForTest()
	db := test.DB(t)
	test.ResetTables(t, db, "flight")
	repo := NewRepository(db, logger)

	ctx := context.Background()

	// initial count
	count, err := repo.Count(ctx)
	assert.Nil(t, err)

	// create
	err = repo.Create(ctx, entity.Flight{
		ID:            "test1",
		Name:          "flight1",
		Number:        "123",
		Departure:     "Minsk",
		DepartureTime: time.Now(),
		Destination:   "Stockholm",
		ArrivalTime:   time.Now(),
		Fare:          "100EUR",
		Duration:      "2 hours",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	})
	assert.Nil(t, err)
	count2, _ := repo.Count(ctx)
	assert.Equal(t, 1, count2-count)

	// get
	flight, err := repo.Get(ctx, "test1")
	assert.Nil(t, err)
	assert.Equal(t, "flight1", flight.Name)
	_, err = repo.Get(ctx, "test0")
	assert.Equal(t, sql.ErrNoRows, err)

	// update
	err = repo.Update(ctx, entity.Flight{
		ID:            "test1",
		Name:          "flight1 updated",
		Departure:     "MOSKOW",
		Destination:   "MINSK",
		Fare:          "200 EUR",
		DepartureTime: time.Now(),
		ArrivalTime:   time.Now().Add(3 * time.Hour),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	})
	assert.Nil(t, err)
	flight, _ = repo.Get(ctx, "test1")
	assert.Equal(t, "flight1 updated", flight.Name)
	assert.Equal(t, "MOSKOW", flight.Departure)
	assert.Equal(t, "MINSK", flight.Destination)
	assert.Equal(t, "200 EUR", flight.Fare)

	// query all
	flights, err := repo.Query(ctx, SearchFlightRequest{}, 0, count2)
	assert.Nil(t, err)
	assert.Equal(t, count2, len(flights))

	// query by name
	flightsByName, err := repo.Query(ctx, SearchFlightRequest{Name: "flight1 updated"}, 0, count2)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(flightsByName))

	// delete
	err = repo.Delete(ctx, "test1")
	assert.Nil(t, err)
	_, err = repo.Get(ctx, "test1")
	assert.Equal(t, sql.ErrNoRows, err)
	err = repo.Delete(ctx, "test1")
	assert.Equal(t, sql.ErrNoRows, err)
}

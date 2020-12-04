package flight

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/nvnoskov/dynamo-backend/internal/entity"
	"github.com/nvnoskov/dynamo-backend/pkg/log"
	"github.com/stretchr/testify/assert"
)

var errCRUD = errors.New("error crud")

func TestCreateFlightRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     CreateFlightRequest
		wantError bool
	}{
		{"success", CreateFlightRequest{
			Name:          "test",
			Number:        "test number",
			Departure:     "MOSKOW",
			Destination:   "MINSK",
			Fare:          "200 EUR",
			DepartureTime: time.Now(),
			ArrivalTime:   time.Now().Add(3 * time.Hour),
		}, false},
		{"required", CreateFlightRequest{Name: ""}, true},
		{"too long", CreateFlightRequest{Name: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.model.Validate()
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func TestUpdateFlightRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     UpdateFlightRequest
		wantError bool
	}{
		{"success", UpdateFlightRequest{
			Name:          "test updated",
			Number:        "test number updated",
			Departure:     "MOSKOW updated",
			Destination:   "MINSK updated",
			Fare:          "200 EUR updated",
			DepartureTime: time.Now(),
			ArrivalTime:   time.Now().Add(3 * time.Hour),
		}, false},
		{"required", UpdateFlightRequest{Name: ""}, true},
		{"too long", UpdateFlightRequest{Name: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.model.Validate()
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func Test_service_CRUD(t *testing.T) {
	logger, _ := log.NewForTest()
	s := NewService(&mockRepository{}, logger)

	ctx := context.Background()

	// initial count
	count, _ := s.Count(ctx)
	assert.Equal(t, 0, count)

	// successful creation
	flight, err := s.Create(ctx, CreateFlightRequest{
		Name:          "test",
		Number:        "test number",
		Departure:     "MOSKOW",
		Destination:   "MINSK",
		Fare:          "200 EUR",
		DepartureTime: time.Now(),
		ArrivalTime:   time.Now().Add(3 * time.Hour),
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, flight.ID)
	id := flight.ID
	assert.Equal(t, "test", flight.Name)
	assert.NotEmpty(t, flight.CreatedAt)
	assert.NotEmpty(t, flight.UpdatedAt)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// validation error in creation
	_, err = s.Create(ctx, CreateFlightRequest{Name: ""})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// unexpected error in creation
	_, err = s.Create(ctx, CreateFlightRequest{
		Name:          "error",
		Number:        "test number",
		Departure:     "MOSKOW",
		Destination:   "MINSK",
		Fare:          "200 EUR",
		DepartureTime: time.Now(),
		ArrivalTime:   time.Now().Add(3 * time.Hour),
	})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	_, _ = s.Create(ctx, CreateFlightRequest{
		Name:          "test2",
		Number:        "test number 2",
		Departure:     "MOSKOW 2",
		Destination:   "MINSK 2",
		Fare:          "200 EUR",
		DepartureTime: time.Now(),
		ArrivalTime:   time.Now().Add(3 * time.Hour),
	})

	// update
	flight, err = s.Update(ctx, id, UpdateFlightRequest{
		Name:          "test updated",
		Number:        "test number",
		Departure:     "MOSKOW",
		Destination:   "MINSK",
		Fare:          "200 EUR",
		DepartureTime: time.Now(),
		ArrivalTime:   time.Now().Add(3 * time.Hour),
	})
	assert.Nil(t, err)
	assert.Equal(t, "test updated", flight.Name)
	assert.Equal(t, "MOSKOW", flight.Departure)
	assert.Equal(t, "MINSK", flight.Destination)
	assert.Equal(t, "200 EUR", flight.Fare)
	assert.Equal(t, "3 hours", flight.Duration)

	_, err = s.Update(ctx, "none", UpdateFlightRequest{Name: "test updated"})
	assert.NotNil(t, err)

	// validation error in update
	_, err = s.Update(ctx, id, UpdateFlightRequest{Name: ""})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 2, count)

	// unexpected error in update
	_, err = s.Update(ctx, id, UpdateFlightRequest{
		Name:          "error",
		Number:        "test number",
		Departure:     "MOSKOW",
		Destination:   "MINSK",
		Fare:          "200 EUR",
		DepartureTime: time.Now(),
		ArrivalTime:   time.Now().Add(3 * time.Hour),
	})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 2, count)

	// get
	_, err = s.Get(ctx, "none")
	assert.NotNil(t, err)
	flight, err = s.Get(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "test updated", flight.Name)
	assert.Equal(t, id, flight.ID)

	// query full
	flightsAll, _ := s.Query(ctx, SearchFlightRequest{}, 0, 0)
	assert.Equal(t, 2, len(flightsAll))

	// delete
	_, err = s.Delete(ctx, "none")
	assert.NotNil(t, err)
	flight, err = s.Delete(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, id, flight.ID)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)
}

type mockRepository struct {
	items []entity.Flight
}

func (m mockRepository) Get(ctx context.Context, id string) (entity.Flight, error) {
	for _, item := range m.items {
		if item.ID == id {
			return item, nil
		}
	}
	return entity.Flight{}, sql.ErrNoRows
}

func (m mockRepository) Count(ctx context.Context) (int, error) {
	return len(m.items), nil
}

func (m mockRepository) Query(ctx context.Context, req SearchFlightRequest, offset, limit int) ([]entity.Flight, error) {
	return m.items, nil
}

func (m *mockRepository) Create(ctx context.Context, flight entity.Flight) error {
	if flight.Name == "error" {
		return errCRUD
	}
	m.items = append(m.items, flight)
	return nil
}

func (m *mockRepository) Update(ctx context.Context, flight entity.Flight) error {
	if flight.Name == "error" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.ID == flight.ID {
			m.items[i] = flight
			break
		}
	}
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id string) error {
	for i, item := range m.items {
		if item.ID == id {
			m.items[i] = m.items[len(m.items)-1]
			m.items = m.items[:len(m.items)-1]
			break
		}
	}
	return nil
}

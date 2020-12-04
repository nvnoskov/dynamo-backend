package auth

import (
	"context"
	"database/sql"
	"testing"

	"github.com/nvnoskov/dynamo-backend/internal/entity"
	"github.com/nvnoskov/dynamo-backend/internal/errors"
	"github.com/nvnoskov/dynamo-backend/pkg/log"
	"github.com/stretchr/testify/assert"
)

func Test_service_Authenticate(t *testing.T) {
	logger, _ := log.NewForTest()
	s := NewService(
		&mockRepository{items: []entity.User{
			{
				ID:       "100",
				Name:     "demo",
				Password: "$2a$10$6gKu8va5UqM48gd/iJdrJOyNMx1GgX6OFymxKuccbbC6nS/LKlu5m",
				Email:    "demo@demo.com",
			},
		}},
		// &mockRepository{}
		"test", 100, logger)
	_, err := s.Login(context.Background(), "unknown", "bad")
	assert.Equal(t, errors.Unauthorized(""), err)
	token, err := s.Login(context.Background(), "demo", "pass")
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
}

func Test_service_GenerateJWT(t *testing.T) {
	logger, _ := log.NewForTest()
	s := service{&mockRepository{}, "test", 100, logger}
	token, err := s.generateJWT(entity.User{
		ID:   "100",
		Name: "demo",
	})
	if assert.Nil(t, err) {
		assert.NotEmpty(t, token)
	}
}

type mockRepository struct {
	items []entity.User
}

func (m mockRepository) Get(ctx context.Context, id string) (entity.User, error) {
	for _, item := range m.items {
		if item.ID == id {
			return item, nil
		}
	}
	return entity.User{}, sql.ErrNoRows
}
func (m mockRepository) GetByUsername(ctx context.Context, username string) (entity.User, error) {
	for _, item := range m.items {
		if item.Name == username {
			return item, nil
		}
	}
	return entity.User{}, sql.ErrNoRows
}

func (m *mockRepository) Create(ctx context.Context, flight entity.User) error {
	m.items = append(m.items, flight)
	return nil
}

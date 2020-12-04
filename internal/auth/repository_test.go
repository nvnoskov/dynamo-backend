package auth

import (
	"context"
	"database/sql"
	"testing"

	"github.com/nvnoskov/dynamo-backend/internal/entity"
	"github.com/nvnoskov/dynamo-backend/internal/test"
	"github.com/nvnoskov/dynamo-backend/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	logger, _ := log.NewForTest()
	db := test.DB(t)
	test.ResetTables(t, db, "user")
	repo := NewRepository(db, logger)

	ctx := context.Background()

	// create
	err := repo.Create(ctx, entity.User{
		ID:       "test1",
		Name:     "user1",
		Password: "123",
		Email:    "user1@mail.com",
	})

	// get
	user, err := repo.Get(ctx, "test1")
	assert.Nil(t, err)
	assert.Equal(t, "user1", user.Name)
	_, err = repo.Get(ctx, "test0")
	assert.Equal(t, sql.ErrNoRows, err)

}

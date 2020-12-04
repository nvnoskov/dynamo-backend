package auth

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/nvnoskov/dynamo-backend/internal/entity"
	"github.com/nvnoskov/dynamo-backend/pkg/dbcontext"
	"github.com/nvnoskov/dynamo-backend/pkg/log"
)

// Repository encapsulates the logic to access users from the data source.
type Repository interface {
	// Get returns the user with the specified user ID.
	Get(ctx context.Context, id string) (entity.User, error)
	GetByUsername(ctx context.Context, username string) (entity.User, error)
	// Create saves a new user in the storage.
	Create(ctx context.Context, user entity.User) error
}

// repository persists users in database
type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

// NewRepository creates a new user repository
func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

// Get reads the user with the specified ID from the database.
func (r repository) Get(ctx context.Context, id string) (entity.User, error) {
	var user entity.User
	err := r.db.With(ctx).Select().Model(id, &user)
	return user, err
}

// Get reads the user with the specified ID from the database.
func (r repository) GetByUsername(ctx context.Context, username string) (entity.User, error) {
	var user entity.User
	err := r.db.With(ctx).Select().Where(dbx.HashExp{"name": username}).One(&user)

	return user, err
}

// Create saves a new user record in the database.
// It returns the ID of the newly inserted user record.
func (r repository) Create(ctx context.Context, user entity.User) error {
	return r.db.With(ctx).Model(&user).Insert()
}

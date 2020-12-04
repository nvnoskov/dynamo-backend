package auth

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/nvnoskov/dynamo-backend/internal/entity"
	"github.com/nvnoskov/dynamo-backend/internal/errors"
	"github.com/nvnoskov/dynamo-backend/pkg/log"
	"golang.org/x/crypto/bcrypt"
)

// Service encapsulates the authentication logic.
type Service interface {
	// authenticate authenticates a user using username and password.
	// It returns a JWT token if authentication succeeds. Otherwise, an error is returned.
	Login(ctx context.Context, username, password string) (string, error)
	Register(ctx context.Context, req RegisterRequest) (User, error)
	Get(ctx context.Context, id string) (User, error)
}

// User represents the data about an user.
type User struct {
	entity.User
}

// Identity represents an authenticated user identity.
type Identity interface {
	// GetID returns the user ID.
	GetID() string
	// GetName returns the user name.
	GetName() string
}

type service struct {
	repo            Repository
	signingKey      string
	tokenExpiration int
	logger          log.Logger
}

// NewService creates a new authentication service.
func NewService(repo Repository, signingKey string, tokenExpiration int, logger log.Logger) Service {
	return service{repo, signingKey, tokenExpiration, logger}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(ctx context.Context, username, password string) (string, error) {
	if identity := s.authenticate(ctx, username, password); identity != nil {
		return s.generateJWT(identity)
	}
	return "", errors.Unauthorized("")
}

// RegisterRequest .
type RegisterRequest struct {
	Name     string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// Validate validates the UpdateFlightRequest fields.
func (m RegisterRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 20)),
		validation.Field(&m.Email, validation.Required, validation.Length(0, 50)),
		validation.Field(&m.Password, validation.Required, validation.Length(0, 50)),
	)
}

// Get returns the flight with the specified the flight ID.
func (s service) Get(ctx context.Context, id string) (User, error) {
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		return User{}, err
	}
	return User{user}, nil
}

// Register create a user
// Otherwise, an error is returned.
func (s service) Register(ctx context.Context, req RegisterRequest) (User, error) {

	if err := req.Validate(); err != nil {
		return User{}, err
	}
	id := entity.GenerateID()

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	err = s.repo.Create(ctx, entity.User{
		ID:       id,
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hash),
	})
	if err != nil {
		return User{}, err
	}
	return s.Get(ctx, id)
}

// authenticate authenticates a user using username and password.
// If username and password are correct, an identity is returned. Otherwise, nil is returned.
func (s service) authenticate(ctx context.Context, username, password string) Identity {
	logger := s.logger.With(ctx, "user", username)

	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		logger.Infof("User not found")
		return nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logger.Infof("authentication failed")
		return nil
	}

	logger.Infof("authentication successful")
	return User{user}

}

// generateJWT generates a JWT that encodes an identity.
func (s service) generateJWT(identity Identity) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   identity.GetID(),
		"name": identity.GetName(),
		"exp":  time.Now().Add(time.Duration(s.tokenExpiration) * time.Hour).Unix(),
	}).SignedString([]byte(s.signingKey))
}

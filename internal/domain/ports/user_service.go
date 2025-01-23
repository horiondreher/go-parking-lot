package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/horiondreher/go-parking-lot/internal/adapters/pgsqlc"
	"github.com/horiondreher/go-parking-lot/internal/domain/errors"
)

type NewUser struct {
	FullName string
	Email    string
	Password string
}

type LoginUser struct {
	Email    string
	Password string
}

type NewUserSession struct {
	RefreshTokenID        uuid.UUID
	Email                 string
	RefreshToken          string
	UserAgent             string
	ClientIP              string
	RefreshTokenExpiresAt time.Time
}

type UserService interface {
	CreateUser(ctx context.Context, newUser NewUser) (pgsqlc.CreateUserRow, *errors.DomainError)
	LoginUser(ctx context.Context, loginUser LoginUser) (pgsqlc.User, *errors.DomainError)
	CreateUserSession(ctx context.Context, newUserSession NewUserSession) (pgsqlc.Session, *errors.DomainError)
	GetUserSession(ctx context.Context, refreshTokenID uuid.UUID) (pgsqlc.Session, *errors.DomainError)
	GetUserByUID(ctx context.Context, userUID string) (pgsqlc.User, *errors.DomainError)
}

package services

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/horiondreher/go-parking-lot/internal/adapters/pgsqlc"
	"github.com/horiondreher/go-parking-lot/internal/domain/errors"
	"github.com/horiondreher/go-parking-lot/internal/domain/ports"
	"github.com/horiondreher/go-parking-lot/internal/utils"
)

type UserManager struct {
	store pgsqlc.Querier
}

func NewUserManager(store pgsqlc.Querier) *UserManager {
	return &UserManager{
		store: store,
	}
}

func (service *UserManager) CreateUser(ctx context.Context, newUser ports.NewUser) (pgsqlc.CreateUserRow, *errors.DomainError) {
	hashedPassword, err := utils.HashPassword(newUser.Password)
	if err != nil {
		return pgsqlc.CreateUserRow{}, errors.NewDomainError(http.StatusInternalServerError, errors.InternalError, err.Error(), err)
	}

	args := pgsqlc.CreateUserParams{
		Email:     newUser.Email,
		Password:  hashedPassword,
		FullName:  newUser.FullName,
		IsStaff:   false,
		IsActive:  true,
		LastLogin: time.Now(),
	}

	user, err := service.store.CreateUser(ctx, args)
	if err != nil {
		return pgsqlc.CreateUserRow{}, errors.MatchPostgresError(err)
	}

	return user, nil
}

func (service *UserManager) LoginUser(ctx context.Context, loginUser ports.LoginUser) (pgsqlc.User, *errors.DomainError) {
	user, err := service.store.GetUser(ctx, loginUser.Email)
	if err != nil {
		return pgsqlc.User{}, errors.MatchPostgresError(err)
	}

	err = utils.CheckPassword(loginUser.Password, user.Password)
	if err != nil {
		return pgsqlc.User{}, errors.MatchHashError(err)
	}

	return user, nil
}

func (service *UserManager) CreateUserSession(ctx context.Context, newUserSession ports.NewUserSession) (pgsqlc.Session, *errors.DomainError) {
	session, err := service.store.CreateSession(ctx, pgsqlc.CreateSessionParams{
		UID:          newUserSession.RefreshTokenID,
		UserEmail:    newUserSession.Email,
		RefreshToken: newUserSession.RefreshToken,
		ExpiresAt:    newUserSession.RefreshTokenExpiresAt,
		UserAgent:    newUserSession.UserAgent,
		ClientIP:     newUserSession.ClientIP,
	})
	if err != nil {
		return pgsqlc.Session{}, errors.MatchPostgresError(err)
	}

	return session, nil
}

func (service *UserManager) GetUserSession(ctx context.Context, refreshTokenID uuid.UUID) (pgsqlc.Session, *errors.DomainError) {
	session, err := service.store.GetSession(ctx, refreshTokenID)
	if err != nil {
		return pgsqlc.Session{}, errors.MatchPostgresError(err)
	}

	return session, nil
}

func (service *UserManager) GetUserByUID(ctx context.Context, userUID string) (pgsqlc.User, *errors.DomainError) {
	parsedUID, err := uuid.Parse(userUID)
	if err != nil {
		return pgsqlc.User{}, errors.NewDomainError(http.StatusInternalServerError, errors.UnexpectedError, err.Error(), err)
	}

	user, err := service.store.GetUserByUID(ctx, parsedUID)
	if err != nil {
		return pgsqlc.User{}, errors.MatchPostgresError(err)
	}

	return user, nil
}

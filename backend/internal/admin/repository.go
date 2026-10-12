package admin

import (
	"context"
	"time"
)

// Repository persists admin users and sessions.
type Repository interface {
	CreateUser(ctx context.Context, user User) error
	FindUserByEmail(ctx context.Context, email string) (User, error)
	FindUserByID(ctx context.Context, id string) (User, error)
	CreateSession(ctx context.Context, session Session) error
	FindSession(ctx context.Context, tokenHash string) (Session, error)
	TouchSession(ctx context.Context, tokenHash string, at time.Time) error
	DeleteSession(ctx context.Context, tokenHash string) error
	DeleteExpiredSessions(ctx context.Context, now time.Time) error
}

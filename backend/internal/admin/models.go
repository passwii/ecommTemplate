package admin

import (
	"errors"
	"time"
)

// ErrNotFound is returned when an admin record does not exist.
var ErrNotFound = errors.New("admin: not found")

// ErrEmailTaken is returned when creating a user with an existing email.
var ErrEmailTaken = errors.New("admin: email already registered")

// User is an internal store administrator.
type User struct {
	ID           string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

// Session is an authenticated admin session. Only the token hash is persisted.
type Session struct {
	TokenHash  string
	UserID     string
	CreatedAt  time.Time
	LastSeenAt time.Time
	ExpiresAt  time.Time
}

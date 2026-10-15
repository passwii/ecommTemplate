package admin

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// DefaultSessionTTL is how long an admin session stays valid.
const DefaultSessionTTL = 12 * time.Hour

// Service authenticates store administrators and manages their sessions.
type Service struct {
	repo       Repository
	now        func() time.Time
	sessionTTL time.Duration
	newToken   func() (string, string, error)
	newID      func() (string, error)
}

// NewService builds the admin authentication service.
func NewService(repo Repository) *Service {
	return &Service{
		repo:       repo,
		now:        time.Now,
		sessionTTL: DefaultSessionTTL,
		newToken:   NewSessionToken,
		newID:      generateUserID,
	}
}

// CreateUser registers an admin account with a hashed password.
func (s *Service) CreateUser(ctx context.Context, email, password string) (User, error) {
	normalized := NormalizeEmail(email)
	if normalized == "" {
		return User{}, errors.New("admin: email is required")
	}
	hash, err := HashPassword(password)
	if err != nil {
		return User{}, err
	}
	id, err := s.newID()
	if err != nil {
		return User{}, err
	}
	user := User{ID: id, Email: normalized, PasswordHash: hash}
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return User{}, err
	}
	return user, nil
}

// Login verifies credentials and creates a session. It returns
// ErrInvalidCredentials for both an unknown email and a wrong password, so
// callers cannot enumerate accounts.
func (s *Service) Login(ctx context.Context, email, password string) (string, User, error) {
	user, err := s.repo.FindUserByEmail(ctx, NormalizeEmail(email))
	if errors.Is(err, ErrNotFound) {
		return "", User{}, ErrInvalidCredentials
	}
	if err != nil {
		return "", User{}, fmt.Errorf("look up admin user: %w", err)
	}
	if !VerifyPassword(user.PasswordHash, password) {
		return "", User{}, ErrInvalidCredentials
	}

	token, tokenHash, err := s.newToken()
	if err != nil {
		return "", User{}, err
	}
	now := s.now().UTC()
	if err := s.repo.CreateSession(ctx, Session{
		TokenHash:  tokenHash,
		UserID:     user.ID,
		CreatedAt:  now,
		LastSeenAt: now,
		ExpiresAt:  now.Add(s.sessionTTL),
	}); err != nil {
		return "", User{}, err
	}
	return token, user, nil
}

// Authenticate resolves a session token to its admin user. It rejects expired
// sessions and refreshes the last-seen timestamp for valid ones.
func (s *Service) Authenticate(ctx context.Context, token string) (User, error) {
	if token == "" {
		return User{}, ErrInvalidCredentials
	}
	tokenHash := HashToken(token)

	session, err := s.repo.FindSession(ctx, tokenHash)
	if errors.Is(err, ErrNotFound) {
		return User{}, ErrInvalidCredentials
	}
	if err != nil {
		return User{}, fmt.Errorf("look up admin session: %w", err)
	}

	now := s.now().UTC()
	if !session.ExpiresAt.After(now) {
		_ = s.repo.DeleteSession(ctx, tokenHash)
		return User{}, ErrInvalidCredentials
	}

	user, err := s.repo.FindUserByID(ctx, session.UserID)
	if errors.Is(err, ErrNotFound) {
		_ = s.repo.DeleteSession(ctx, tokenHash)
		return User{}, ErrInvalidCredentials
	}
	if err != nil {
		return User{}, fmt.Errorf("look up admin user: %w", err)
	}

	if err := s.repo.TouchSession(ctx, tokenHash, now); err != nil {
		// A failed touch must not block an otherwise valid session.
		_ = err
	}
	return user, nil
}

// Logout removes a session.
func (s *Service) Logout(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	return s.repo.DeleteSession(ctx, HashToken(token))
}

func generateUserID() (string, error) {
	token, _, err := NewSessionToken()
	if err != nil {
		return "", err
	}
	return "adm_" + token[:22], nil
}

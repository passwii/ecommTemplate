package admin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// PostgresRepository is the SQL-backed implementation of Repository.
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository builds an admin repository over the given database.
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

const uniqueViolation = "23505"

func (r *PostgresRepository) CreateUser(ctx context.Context, user User) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO admin_users (id, email, password_hash)
		VALUES ($1, $2, $3)`, user.ID, user.Email, user.PasswordHash)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == uniqueViolation {
			return ErrEmailTaken
		}
		return fmt.Errorf("create admin user: %w", err)
	}
	return nil
}

func (r *PostgresRepository) FindUserByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := r.db.QueryRowContext(ctx, `
		SELECT id, email, password_hash, created_at
		FROM admin_users
		WHERE email = $1`, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrNotFound
	}
	if err != nil {
		return User{}, fmt.Errorf("find admin user: %w", err)
	}
	return user, nil
}

func (r *PostgresRepository) FindUserByID(ctx context.Context, id string) (User, error) {
	var user User
	err := r.db.QueryRowContext(ctx, `
		SELECT id, email, password_hash, created_at
		FROM admin_users
		WHERE id = $1`, id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrNotFound
	}
	if err != nil {
		return User{}, fmt.Errorf("find admin user: %w", err)
	}
	return user, nil
}

func (r *PostgresRepository) CreateSession(ctx context.Context, session Session) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO admin_sessions (token_hash, admin_user_id, created_at, last_seen_at, expires_at)
		VALUES ($1, $2, $3, $3, $4)`,
		session.TokenHash, session.UserID, session.CreatedAt, session.ExpiresAt)
	if err != nil {
		return fmt.Errorf("create admin session: %w", err)
	}
	return nil
}

func (r *PostgresRepository) FindSession(ctx context.Context, tokenHash string) (Session, error) {
	var session Session
	err := r.db.QueryRowContext(ctx, `
		SELECT token_hash, admin_user_id, created_at, last_seen_at, expires_at
		FROM admin_sessions
		WHERE token_hash = $1`, tokenHash).
		Scan(&session.TokenHash, &session.UserID, &session.CreatedAt, &session.LastSeenAt, &session.ExpiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Session{}, ErrNotFound
	}
	if err != nil {
		return Session{}, fmt.Errorf("find admin session: %w", err)
	}
	return session, nil
}

func (r *PostgresRepository) TouchSession(ctx context.Context, tokenHash string, at time.Time) error {
	if _, err := r.db.ExecContext(ctx, `
		UPDATE admin_sessions SET last_seen_at = $2 WHERE token_hash = $1`, tokenHash, at); err != nil {
		return fmt.Errorf("touch admin session: %w", err)
	}
	return nil
}

func (r *PostgresRepository) DeleteSession(ctx context.Context, tokenHash string) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM admin_sessions WHERE token_hash = $1`, tokenHash); err != nil {
		return fmt.Errorf("delete admin session: %w", err)
	}
	return nil
}

func (r *PostgresRepository) DeleteExpiredSessions(ctx context.Context, now time.Time) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM admin_sessions WHERE expires_at <= $1`, now); err != nil {
		return fmt.Errorf("delete expired admin sessions: %w", err)
	}
	return nil
}

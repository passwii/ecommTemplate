package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	maxOpenConnections = 10
	maxIdleConnections = 5
	connectionLifetime = 30 * time.Minute
	pingTimeout        = 3 * time.Second
)

// Open creates and verifies a PostgreSQL connection pool without exposing the DSN.
func Open(ctx context.Context, databaseURL string) (*sql.DB, error) {
	parsed, err := url.Parse(databaseURL)
	if err != nil || parsed.Scheme != "postgres" && parsed.Scheme != "postgresql" || parsed.Host == "" || parsed.Path == "/" {
		return nil, errors.New("invalid PostgreSQL connection configuration")
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, errors.New("open PostgreSQL connection pool")
	}
	db.SetMaxOpenConns(maxOpenConnections)
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetConnMaxLifetime(connectionLifetime)

	pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping PostgreSQL: %w", redactConnectionDetails(err, parsed))
	}
	return db, nil
}

func redactConnectionDetails(err error, parsed *url.URL) error {
	message := err.Error()
	if parsed.User != nil {
		if password, ok := parsed.User.Password(); ok && password != "" {
			message = replaceAll(message, password, "[redacted]")
		}
		message = replaceAll(message, parsed.User.Username(), "[redacted]")
	}
	return errors.New(message)
}

func replaceAll(value, old, replacement string) string {
	if old == "" {
		return value
	}
	return strings.ReplaceAll(value, old, replacement)
}

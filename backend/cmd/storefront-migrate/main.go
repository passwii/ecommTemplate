package main

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/passwii/ecommTemple/backend/internal/database"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func main() {
	if len(os.Args) != 2 || os.Args[1] != "up" && os.Args[1] != "status" {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/storefront-migrate <up|status>")
		os.Exit(2)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if err := database.ValidateLocalDatabaseURL(databaseURL); err != nil {
		slog.Error("invalid database configuration", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		slog.Error("open database failed")
		os.Exit(1)
	}
	defer db.Close()
	if err := db.PingContext(ctx); err != nil {
		slog.Error("database unavailable; no migration was run")
		os.Exit(1)
	}

	goose.SetBaseFS(migrationFiles)
	if err := goose.SetDialect("postgres"); err != nil {
		slog.Error("configure migration dialect", "error", err)
		os.Exit(1)
	}

	if err := runMigrationCommand(ctx, db, os.Args[1]); err != nil {
		slog.Error("migration command failed", "error", err)
		os.Exit(1)
	}
}

func runMigrationCommand(ctx context.Context, db *sql.DB, command string) error {
	migrationFS, err := fs.Sub(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("load embedded migrations: %w", err)
	}
	provider, err := goose.NewProvider(goose.DialectPostgres, db, migrationFS)
	if err != nil {
		return fmt.Errorf("initialize migration provider: %w", err)
	}
	defer provider.Close()

	switch command {
	case "up":
		results, err := provider.Up(ctx)
		if err != nil {
			return fmt.Errorf("apply migrations: %w", err)
		}
		for _, result := range results {
			slog.Info("migration applied", "version", result.Source.Version, "name", result.Source.Path)
		}
		if len(results) == 0 {
			slog.Info("database schema already current")
		}
	case "status":
		statuses, err := provider.Status(ctx)
		if err != nil {
			return fmt.Errorf("read migration status: %w", err)
		}
		if len(statuses) == 0 {
			fmt.Println("no migrations available")
			return nil
		}
		for _, status := range statuses {
			fmt.Printf("%-8s %s\n", status.State, status.Source.Path)
		}
	default:
		return errors.New("unsupported migration command")
	}
	return nil
}

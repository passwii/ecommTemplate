package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/passwii/ecommTemple/backend/internal/catalog"
	"github.com/passwii/ecommTemple/backend/internal/database"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if err := database.ValidateLocalDatabaseURL(databaseURL); err != nil {
		slog.Error("invalid database configuration", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		slog.Error("open database failed")
		os.Exit(1)
	}
	defer db.Close()
	if err := db.PingContext(ctx); err != nil {
		slog.Error("database unavailable; nothing was seeded")
		os.Exit(1)
	}

	if err := catalog.Seed(ctx, db, catalog.SampleData()); err != nil {
		slog.Error("seed failed", "error", err)
		os.Exit(1)
	}
	slog.Info("sample catalog seeded", "template", catalog.SampleData().TemplateID)
}

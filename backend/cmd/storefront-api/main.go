package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/passwii/ecommTemple/backend/internal/catalog"
	"github.com/passwii/ecommTemple/backend/internal/config"
	"github.com/passwii/ecommTemple/backend/internal/database"
	"github.com/passwii/ecommTemple/backend/internal/httpapi"
	"github.com/passwii/ecommTemple/backend/internal/inquiry"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.Load()
	if err != nil {
		logger.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := database.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("database unavailable", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	catalogRepo := catalog.NewPostgresRepository(db)
	inquiryRepo := inquiry.NewPostgresRepository(db)
	inquiryService := inquiry.NewService(inquiryRepo, catalogRepo)

	server := &http.Server{
		Addr: cfg.HTTPAddr,
		Handler: httpapi.NewHandler(httpapi.Deps{
			DB:               db,
			Catalog:          catalogRepo,
			Inquiries:        inquiryService,
			Limiter:          inquiry.NewRateLimiter(cfg.InquiryRateLimit, cfg.InquiryRateWindow),
			TrustProxyHeader: cfg.TrustProxyHeader,
		}),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("storefront API listening", "addr", cfg.HTTPAddr)
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server stopped", "error", err)
			os.Exit(1)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("HTTP server shutdown failed", "error", err)
			os.Exit(1)
		}
	}
}

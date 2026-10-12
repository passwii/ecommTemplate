package config

import (
	"strings"
	"testing"
)

func TestLoadRequiresDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("HTTP_ADDR", "")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want missing DATABASE_URL error")
	}
	if strings.Contains(err.Error(), "secret") {
		t.Fatalf("Load() error leaked secret: %v", err)
	}
}

func TestLoadRejectsMalformedDatabaseURLWithoutLeakingCredentials(t *testing.T) {
	const secretURL = "postgres://user:secret%zz@localhost:5432/ecommtemple_dev"
	t.Setenv("DATABASE_URL", secretURL)

	_, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want malformed URL error")
	}
	if strings.Contains(err.Error(), "secret") || strings.Contains(err.Error(), "user") {
		t.Fatalf("Load() error leaked connection details: %v", err)
	}
}

func TestLoadUsesDocumentedHTTPDefault(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/ecommtemple_dev?sslmode=disable")
	t.Setenv("HTTP_ADDR", "")

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got.HTTPAddr != ":8080" {
		t.Errorf("HTTPAddr = %q, want %q", got.HTTPAddr, ":8080")
	}
}

func TestLoadAcceptsValidConfiguration(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/ecommtemple_dev?sslmode=disable")
	t.Setenv("HTTP_ADDR", "127.0.0.1:8090")

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got.HTTPAddr != "127.0.0.1:8090" {
		t.Errorf("HTTPAddr = %q, want explicit address", got.HTTPAddr)
	}
	if got.DatabaseURL == "" {
		t.Error("DatabaseURL is empty")
	}
}

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

func TestLoadUsesInquiryRateLimitDefaults(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/ecommtemple_dev?sslmode=disable")
	t.Setenv("INQUIRY_RATE_LIMIT", "")
	t.Setenv("INQUIRY_RATE_WINDOW", "")

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got.InquiryRateLimit != 5 {
		t.Errorf("InquiryRateLimit = %d, want 5", got.InquiryRateLimit)
	}
	if got.InquiryRateWindow.Minutes() != 10 {
		t.Errorf("InquiryRateWindow = %v, want 10m", got.InquiryRateWindow)
	}
}

func TestLoadDefaultsToNotTrustingProxyHeaders(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/ecommtemple_dev?sslmode=disable")
	t.Setenv("TRUST_PROXY_HEADER", "")

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got.TrustProxyHeader {
		t.Error("TrustProxyHeader = true, want false by default")
	}
}

func TestLoadRejectsInvalidTrustProxyHeader(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/ecommtemple_dev?sslmode=disable")
	t.Setenv("TRUST_PROXY_HEADER", "yes-please")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want invalid boolean error")
	}
}

func TestLoadDefaultsSessionCookieToNotSecure(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/ecommtemple_dev?sslmode=disable")
	t.Setenv("SESSION_COOKIE_SECURE", "")

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got.SessionCookieSecure {
		t.Error("SessionCookieSecure = true, want false by default for local http")
	}
}

func TestLoadRejectsInvalidSessionCookieSecure(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/ecommtemple_dev?sslmode=disable")
	t.Setenv("SESSION_COOKIE_SECURE", "maybe")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want invalid boolean error")
	}
}

func TestLoadRejectsInvalidInquiryRateLimit(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/ecommtemple_dev?sslmode=disable")

	t.Run("non-numeric limit", func(t *testing.T) {
		t.Setenv("INQUIRY_RATE_LIMIT", "many")
		if _, err := Load(); err == nil {
			t.Fatal("Load() error = nil, want invalid limit error")
		}
	})

	t.Run("negative limit", func(t *testing.T) {
		t.Setenv("INQUIRY_RATE_LIMIT", "-1")
		if _, err := Load(); err == nil {
			t.Fatal("Load() error = nil, want invalid limit error")
		}
	})

	t.Run("invalid window", func(t *testing.T) {
		t.Setenv("INQUIRY_RATE_LIMIT", "5")
		t.Setenv("INQUIRY_RATE_WINDOW", "soon")
		if _, err := Load(); err == nil {
			t.Fatal("Load() error = nil, want invalid window error")
		}
	})
}

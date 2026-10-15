package database

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestOpenRejectsMalformedURLWithoutLeakingCredentials(t *testing.T) {
	_, err := Open(context.Background(), "postgres://private-user:secret%zz@localhost:5432/ecommtemple_dev")
	if err == nil {
		t.Fatal("Open() error = nil, want invalid URL error")
	}
	if strings.Contains(err.Error(), "private-user") || strings.Contains(err.Error(), "secret") {
		t.Fatalf("Open() error leaked credentials: %v", err)
	}
}

func TestOpenTimesOutWhenDatabaseIsUnavailable(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := Open(ctx, "postgres://nobody:nobody@127.0.0.1:1/ecommtemple_dev?sslmode=disable&connect_timeout=1")
	if err == nil {
		t.Fatal("Open() error = nil, want unavailable database error")
	}
	if strings.Contains(err.Error(), "nobody") {
		t.Fatalf("Open() error leaked credentials: %v", err)
	}
}

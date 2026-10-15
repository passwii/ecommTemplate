package httpapi

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLivenessReturnsOKWithoutDatabase(t *testing.T) {
	handler := NewHandler(Deps{})
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/health/live", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
}

func TestReadinessReturnsUnavailableWithoutDatabase(t *testing.T) {
	handler := NewHandler(Deps{})
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/health/ready", nil))

	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusServiceUnavailable)
	}
}

type failingDriver struct{}

func (failingDriver) Open(string) (driver.Conn, error) {
	return nil, errors.New("connection refused")
}

func TestReadinessReturnsUnavailableWhenDatabaseCannotPing(t *testing.T) {
	// The registered driver always fails to connect; no real database is used.
	driverName := "health-failing-" + time.Now().Format("150405.000000000")
	sql.Register(driverName, failingDriver{})
	db, err := sql.Open(driverName, "unused")
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Millisecond)

	handler := NewHandler(Deps{DB: db})
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/health/ready", nil))

	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusServiceUnavailable)
	}
}

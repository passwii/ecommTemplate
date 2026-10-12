package httpapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/passwii/ecommTemple/backend/internal/catalog"
)

const readinessTimeout = 2 * time.Second

type healthResponse struct {
	Status string `json:"status"`
}

// Deps holds the runtime dependencies used to build the HTTP handler.
type Deps struct {
	DB      *sql.DB
	Catalog catalog.Repository
}

// NewHandler registers public health and catalog endpoints.
func NewHandler(deps Deps) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, _ *http.Request) {
		writeHealth(w, http.StatusOK, "ok")
	})
	mux.HandleFunc("GET /health/ready", func(w http.ResponseWriter, r *http.Request) {
		if deps.DB == nil {
			writeHealth(w, http.StatusServiceUnavailable, "unavailable")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), readinessTimeout)
		defer cancel()
		if err := deps.DB.PingContext(ctx); err != nil {
			writeHealth(w, http.StatusServiceUnavailable, "unavailable")
			return
		}
		writeHealth(w, http.StatusOK, "ok")
	})
	if deps.Catalog != nil {
		registerCatalogRoutes(mux, deps.Catalog)
	}
	return mux
}

func writeHealth(w http.ResponseWriter, status int, state string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(healthResponse{Status: state})
}

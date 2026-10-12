package httpapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/passwii/ecommTemple/backend/internal/catalog"
	"github.com/passwii/ecommTemple/backend/internal/inquiry"
)

const readinessTimeout = 2 * time.Second

type healthResponse struct {
	Status string `json:"status"`
}

// Deps holds the runtime dependencies used to build the HTTP handler.
type Deps struct {
	DB        *sql.DB
	Catalog   catalog.Repository
	Inquiries *inquiry.Service
	Limiter   *inquiry.RateLimiter
	// TrustProxyHeader keys rate limiting on X-Forwarded-For. Enable only
	// behind a trusted proxy that overwrites the header.
	TrustProxyHeader bool
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
	if deps.Inquiries != nil {
		registerInquiryRoutes(mux, deps.Inquiries, deps.Limiter, deps.TrustProxyHeader)
	}
	return mux
}

func writeHealth(w http.ResponseWriter, status int, state string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(healthResponse{Status: state})
}

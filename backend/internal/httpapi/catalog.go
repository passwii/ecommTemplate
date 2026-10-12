package httpapi

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/passwii/ecommTemple/backend/internal/catalog"
)

const (
	defaultProductLimit = 24
	maxProductLimit     = 100
)

type catalogHandlers struct {
	repo catalog.Repository
}

// registerCatalogRoutes adds public read-only catalog and company endpoints.
func registerCatalogRoutes(mux *http.ServeMux, repo catalog.Repository) {
	handlers := &catalogHandlers{repo: repo}
	mux.HandleFunc("GET /api/v1/collections", handlers.listCollections)
	mux.HandleFunc("GET /api/v1/collections/{slug}", handlers.getCollection)
	mux.HandleFunc("GET /api/v1/products", handlers.listProducts)
	mux.HandleFunc("GET /api/v1/products/{slug}", handlers.getProduct)
	mux.HandleFunc("GET /api/v1/company-content", handlers.getCompanyContent)
}

func (h *catalogHandlers) listCollections(w http.ResponseWriter, r *http.Request) {
	collections, err := h.repo.ListPublishedCollections(r.Context())
	if err != nil {
		writeRepositoryError(w, "list collections", err)
		return
	}
	if collections == nil {
		collections = []catalog.Collection{}
	}
	writeData(w, http.StatusOK, collections)
}

func (h *catalogHandlers) getCollection(w http.ResponseWriter, r *http.Request) {
	collection, err := h.repo.GetPublishedCollection(r.Context(), r.PathValue("slug"))
	if err != nil {
		writeRepositoryError(w, "get collection", err)
		return
	}
	writeData(w, http.StatusOK, collection)
}

func (h *catalogHandlers) listProducts(w http.ResponseWriter, r *http.Request) {
	filter, ok := parseProductFilter(w, r)
	if !ok {
		return
	}
	products, err := h.repo.ListPublishedProducts(r.Context(), filter)
	if err != nil {
		writeRepositoryError(w, "list products", err)
		return
	}
	if products == nil {
		products = []catalog.Product{}
	}
	writeData(w, http.StatusOK, products)
}

func (h *catalogHandlers) getProduct(w http.ResponseWriter, r *http.Request) {
	product, err := h.repo.GetPublishedProduct(r.Context(), r.PathValue("slug"))
	if err != nil {
		writeRepositoryError(w, "get product", err)
		return
	}
	writeData(w, http.StatusOK, product)
}

func (h *catalogHandlers) getCompanyContent(w http.ResponseWriter, r *http.Request) {
	content, err := h.repo.GetCompanyContent(r.Context())
	if err != nil {
		writeRepositoryError(w, "get company content", err)
		return
	}
	writeData(w, http.StatusOK, content)
}

func parseProductFilter(w http.ResponseWriter, r *http.Request) (catalog.Filter, bool) {
	query := r.URL.Query()
	filter := catalog.Filter{
		CollectionSlug: query.Get("collection"),
		Limit:          defaultProductLimit,
	}

	if raw := query.Get("limit"); raw != "" {
		limit, err := strconv.Atoi(raw)
		if err != nil || limit < 1 || limit > maxProductLimit {
			writeFieldErrors(w, http.StatusBadRequest, "invalid_query", "Invalid query parameters.", map[string]string{
				"limit": "Use a whole number between 1 and 100.",
			})
			return catalog.Filter{}, false
		}
		filter.Limit = limit
	}

	if raw := query.Get("featured"); raw != "" {
		featured, err := strconv.ParseBool(raw)
		if err != nil {
			writeFieldErrors(w, http.StatusBadRequest, "invalid_query", "Invalid query parameters.", map[string]string{
				"featured": "Use true or false.",
			})
			return catalog.Filter{}, false
		}
		filter.Featured = &featured
	}

	return filter, true
}

func writeRepositoryError(w http.ResponseWriter, operation string, err error) {
	if errors.Is(err, catalog.ErrNotFound) {
		writeError(w, http.StatusNotFound, "not_found", "The requested resource was not found.")
		return
	}
	slog.Error("catalog read failed", "operation", operation, "error", err)
	writeError(w, http.StatusInternalServerError, "internal_error", "The request could not be completed.")
}

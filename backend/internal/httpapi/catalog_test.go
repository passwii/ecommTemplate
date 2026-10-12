package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/passwii/ecommTemple/backend/internal/catalog"
)

type fakeRepository struct {
	collections []catalog.Collection
	products    []catalog.Product
	company     catalog.CompanyContent
	err         error
	lastFilter  catalog.Filter
}

func (f *fakeRepository) ListPublishedCollections(context.Context) ([]catalog.Collection, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.collections, nil
}

func (f *fakeRepository) GetPublishedCollection(_ context.Context, slug string) (catalog.Collection, error) {
	if f.err != nil {
		return catalog.Collection{}, f.err
	}
	for _, collection := range f.collections {
		if collection.Slug == slug {
			return collection, nil
		}
	}
	return catalog.Collection{}, catalog.ErrNotFound
}

func (f *fakeRepository) ListPublishedProducts(_ context.Context, filter catalog.Filter) ([]catalog.Product, error) {
	f.lastFilter = filter
	if f.err != nil {
		return nil, f.err
	}
	return f.products, nil
}

func (f *fakeRepository) GetPublishedProduct(_ context.Context, slug string) (catalog.Product, error) {
	if f.err != nil {
		return catalog.Product{}, f.err
	}
	for _, product := range f.products {
		if product.Slug == slug {
			return product, nil
		}
	}
	return catalog.Product{}, catalog.ErrNotFound
}

func (f *fakeRepository) GetPublishedProductByID(_ context.Context, id string) (catalog.Product, error) {
	if f.err != nil {
		return catalog.Product{}, f.err
	}
	for _, product := range f.products {
		if product.ID == id {
			return product, nil
		}
	}
	return catalog.Product{}, catalog.ErrNotFound
}

func (f *fakeRepository) GetCompanyContent(context.Context) (catalog.CompanyContent, error) {
	if f.err != nil {
		return catalog.CompanyContent{}, f.err
	}
	return f.company, nil
}

func newTestHandler(repo catalog.Repository) http.Handler {
	return NewHandler(Deps{Catalog: repo})
}

func doRequest(t *testing.T, handler http.Handler, target string) *httptest.ResponseRecorder {
	t.Helper()
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, target, nil))
	return recorder
}

func TestProductListReturnsDataEnvelope(t *testing.T) {
	repo := &fakeRepository{products: []catalog.Product{{ID: "p-1", Slug: "p-1", Name: "Product"}}}
	recorder := doRequest(t, newTestHandler(repo), "/api/v1/products")

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", recorder.Code)
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
	var body struct {
		Data []catalog.Product `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v (%s)", err, recorder.Body.String())
	}
	if len(body.Data) != 1 || body.Data[0].Slug != "p-1" {
		t.Fatalf("data = %+v, want one product p-1", body.Data)
	}
}

func TestProductListRejectsInvalidLimit(t *testing.T) {
	for _, limit := range []string{"0", "-1", "101", "abc"} {
		recorder := doRequest(t, newTestHandler(&fakeRepository{}), "/api/v1/products?limit="+limit)
		if recorder.Code != http.StatusBadRequest {
			t.Errorf("limit=%s status = %d, want 400", limit, recorder.Code)
		}
	}
}

func TestProductListRejectsInvalidFeatured(t *testing.T) {
	recorder := doRequest(t, newTestHandler(&fakeRepository{}), "/api/v1/products?featured=maybe")
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", recorder.Code)
	}
}

func TestProductListAppliesFilters(t *testing.T) {
	repo := &fakeRepository{}
	recorder := doRequest(t, newTestHandler(repo), "/api/v1/products?collection=grills&featured=true&limit=5")
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", recorder.Code)
	}
	if repo.lastFilter.CollectionSlug != "grills" {
		t.Errorf("CollectionSlug = %q, want grills", repo.lastFilter.CollectionSlug)
	}
	if repo.lastFilter.Featured == nil || !*repo.lastFilter.Featured {
		t.Errorf("Featured = %v, want true", repo.lastFilter.Featured)
	}
	if repo.lastFilter.Limit != 5 {
		t.Errorf("Limit = %d, want 5", repo.lastFilter.Limit)
	}
}

func TestNotFoundResponsesUseStableErrorCode(t *testing.T) {
	targets := []string{"/api/v1/products/missing", "/api/v1/collections/missing"}
	for _, target := range targets {
		recorder := doRequest(t, newTestHandler(&fakeRepository{}), target)
		if recorder.Code != http.StatusNotFound {
			t.Errorf("%s status = %d, want 404", target, recorder.Code)
		}
		var body struct {
			Error struct {
				Code string `json:"code"`
			} `json:"error"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Error.Code != "not_found" {
			t.Errorf("%s error code = %q, want not_found", target, body.Error.Code)
		}
	}
}

func TestInternalErrorDoesNotLeakDetails(t *testing.T) {
	repo := &fakeRepository{err: errors.New("pq: relation \"products\" does not exist at /srv/db")}
	recorder := doRequest(t, newTestHandler(repo), "/api/v1/products")

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", recorder.Code)
	}
	body := recorder.Body.String()
	for _, leaked := range []string{"pq:", "relation", "/srv/db"} {
		if strings.Contains(body, leaked) {
			t.Errorf("error body leaked internal detail %q: %s", leaked, body)
		}
	}
}

func TestCompanyContentNotFound(t *testing.T) {
	repo := &fakeRepository{err: catalog.ErrNotFound}
	recorder := doRequest(t, newTestHandler(repo), "/api/v1/company-content")
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", recorder.Code)
	}
}

func TestCollectionListReturnsEmptyArrayWhenNoData(t *testing.T) {
	recorder := doRequest(t, newTestHandler(&fakeRepository{}), "/api/v1/collections")
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", recorder.Code)
	}
	if got := strings.TrimSpace(recorder.Body.String()); got != `{"data":[]}` {
		t.Fatalf("body = %s, want {\"data\":[]}", got)
	}
}

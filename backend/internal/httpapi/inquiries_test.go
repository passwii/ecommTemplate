package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/passwii/ecommTemple/backend/internal/catalog"
	"github.com/passwii/ecommTemple/backend/internal/inquiry"
)

type fakeInquiryRepository struct {
	created     []inquiry.Inquiry
	byRequestID map[string]inquiry.Inquiry
	err         error
}

func (f *fakeInquiryRepository) Create(_ context.Context, value inquiry.Inquiry) (inquiry.Inquiry, error) {
	if f.err != nil {
		return inquiry.Inquiry{}, f.err
	}
	value.CreatedAt = time.Date(2026, 10, 12, 12, 0, 0, 0, time.UTC)
	value.UpdatedAt = value.CreatedAt
	f.created = append(f.created, value)
	if value.RequestID != "" {
		if f.byRequestID == nil {
			f.byRequestID = map[string]inquiry.Inquiry{}
		}
		f.byRequestID[value.RequestID] = value
	}
	return value, nil
}

func (f *fakeInquiryRepository) FindByRequestID(_ context.Context, requestID string) (inquiry.Inquiry, error) {
	if existing, ok := f.byRequestID[requestID]; ok {
		return existing, nil
	}
	return inquiry.Inquiry{}, inquiry.ErrNotFound
}

func (f *fakeInquiryRepository) Get(_ context.Context, id string) (inquiry.Inquiry, error) {
	for _, value := range f.created {
		if value.ID == id {
			return value, nil
		}
	}
	return inquiry.Inquiry{}, inquiry.ErrNotFound
}

func (f *fakeInquiryRepository) List(context.Context, inquiry.ListFilter) ([]inquiry.Inquiry, int, error) {
	return nil, 0, nil
}

func (f *fakeInquiryRepository) Update(context.Context, string, inquiry.UpdateInput, string) (inquiry.Inquiry, error) {
	return inquiry.Inquiry{}, inquiry.ErrNotFound
}

func (f *fakeInquiryRepository) ListActivities(context.Context, string, int) ([]inquiry.Activity, error) {
	return nil, nil
}

func (f *fakeInquiryRepository) AddActivity(_ context.Context, activity inquiry.Activity) (inquiry.Activity, error) {
	return activity, nil
}

func inquiryHandler(repo *fakeInquiryRepository, products catalog.Repository, limit int) http.Handler {
	service := inquiry.NewService(repo, products)
	return NewHandler(Deps{
		Catalog:   products,
		Inquiries: service,
		Limiter:   inquiry.NewRateLimiter(limit, time.Minute),
	})
}

func validInquiryBody() string {
	return `{
		"contact": {"name": "Alex Morgan", "email": "buyer@example.com", "company": "Example Importers", "phone": "+1 555 0100"},
		"destinationCountry": "US",
		"productId": "grill-01",
		"quantity": 120,
		"customizationRequest": "Private-label packaging options",
		"privacyConsent": {"accepted": true, "noticeVersion": "inquiry-privacy-v1"},
		"sourcePath": "/products/ridge-4-burner"
	}`
}

func postInquiry(t *testing.T, handler http.Handler, body string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/inquiries", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return recorder
}

func publishedProducts() *fakeRepository {
	return &fakeRepository{products: []catalog.Product{{ID: "grill-01", Slug: "ridge-4-burner", Model: "Ridge 4-Burner", Name: "Ridge 4-Burner"}}}
}

func TestCreateInquiryStoresSubmissionWithServerSnapshot(t *testing.T) {
	repo := &fakeInquiryRepository{}
	recorder := postInquiry(t, inquiryHandler(repo, publishedProducts(), 10), validInquiryBody())

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201 (%s)", recorder.Code, recorder.Body.String())
	}
	if len(repo.created) != 1 {
		t.Fatalf("stored inquiries = %d, want 1", len(repo.created))
	}
	stored := repo.created[0]
	if stored.Status != inquiry.StatusNew {
		t.Errorf("status = %q, want new", stored.Status)
	}
	if stored.ProductModel != "Ridge 4-Burner" || stored.ProductSlug != "ridge-4-burner" {
		t.Errorf("product snapshot = %+v, want server-generated snapshot", stored)
	}
	if !strings.HasPrefix(stored.ID, "inq_") {
		t.Errorf("id = %q, want inq_ prefix", stored.ID)
	}
	if stored.PrivacyNoticeVersion != "inquiry-privacy-v1" || stored.PrivacyConsentedAt.IsZero() {
		t.Errorf("privacy evidence missing: %+v", stored)
	}
}

func TestCreateInquiryReturnsFieldErrors(t *testing.T) {
	body := `{"contact":{"name":"","email":"nope"},"privacyConsent":{"accepted":false}}`
	recorder := postInquiry(t, inquiryHandler(&fakeInquiryRepository{}, publishedProducts(), 10), body)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", recorder.Code)
	}
	var envelope struct {
		Error struct {
			Code   string            `json:"code"`
			Fields map[string]string `json:"fields"`
		} `json:"error"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if envelope.Error.Code != "validation_error" {
		t.Errorf("code = %q, want validation_error", envelope.Error.Code)
	}
	for _, field := range []string{"contactName", "contactEmail", "privacyConsent"} {
		if _, ok := envelope.Error.Fields[field]; !ok {
			t.Errorf("fields missing %q: %v", field, envelope.Error.Fields)
		}
	}
}

func TestCreateInquiryRejectsUnknownProduct(t *testing.T) {
	body := strings.Replace(validInquiryBody(), `"grill-01"`, `"does-not-exist"`, 1)
	recorder := postInquiry(t, inquiryHandler(&fakeInquiryRepository{}, publishedProducts(), 10), body)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404 (%s)", recorder.Code, recorder.Body.String())
	}
}

func TestCreateInquiryRejectsForgedOperationalFields(t *testing.T) {
	body := strings.Replace(validInquiryBody(), `"sourcePath"`, `"status": "won", "owner": "attacker", "sourcePath"`, 1)
	recorder := postInquiry(t, inquiryHandler(&fakeInquiryRepository{}, publishedProducts(), 10), body)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", recorder.Code)
	}
}

func TestCreateInquiryRejectsHoneypot(t *testing.T) {
	body := strings.Replace(validInquiryBody(), `"sourcePath"`, `"referenceCode": "http://spam.example.com", "sourcePath"`, 1)
	recorder := postInquiry(t, inquiryHandler(&fakeInquiryRepository{}, publishedProducts(), 10), body)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", recorder.Code)
	}
}

func TestCreateInquiryIsIdempotentWithRequestID(t *testing.T) {
	body := strings.Replace(validInquiryBody(), `"sourcePath"`, `"requestId": "req-123", "sourcePath"`, 1)
	repo := &fakeInquiryRepository{}
	handler := inquiryHandler(repo, publishedProducts(), 10)

	first := postInquiry(t, handler, body)
	if first.Code != http.StatusCreated {
		t.Fatalf("first status = %d, want 201", first.Code)
	}
	second := postInquiry(t, handler, body)
	if second.Code != http.StatusOK {
		t.Fatalf("second status = %d, want 200", second.Code)
	}
	if len(repo.created) != 1 {
		t.Fatalf("stored inquiries = %d, want 1 (idempotent replay)", len(repo.created))
	}
}

func TestCreateInquiryRateLimits(t *testing.T) {
	handler := inquiryHandler(&fakeInquiryRepository{}, publishedProducts(), 2)

	for i := 0; i < 2; i++ {
		if code := postInquiry(t, handler, validInquiryBody()).Code; code != http.StatusCreated {
			t.Fatalf("request %d status = %d, want 201", i+1, code)
		}
	}
	recorder := postInquiry(t, handler, validInquiryBody())
	if recorder.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want 429", recorder.Code)
	}
	var envelope struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if envelope.Error.Code != "rate_limited" {
		t.Errorf("code = %q, want rate_limited", envelope.Error.Code)
	}
}

func TestCreateInquiryRequiresJSONContentType(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/v1/inquiries", strings.NewReader(validInquiryBody()))
	request.Header.Set("Content-Type", "text/plain")
	recorder := httptest.NewRecorder()
	inquiryHandler(&fakeInquiryRepository{}, publishedProducts(), 10).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("status = %d, want 415", recorder.Code)
	}
}

func TestCreateInquiryRejectsInvalidJSON(t *testing.T) {
	recorder := postInquiry(t, inquiryHandler(&fakeInquiryRepository{}, publishedProducts(), 10), `{"contact":`)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", recorder.Code)
	}
}

func TestCreateInquiryDoesNotLeakInternalErrors(t *testing.T) {
	repo := &fakeInquiryRepository{err: &failingError{}}
	recorder := postInquiry(t, inquiryHandler(repo, publishedProducts(), 10), validInquiryBody())

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", recorder.Code)
	}
	if body := recorder.Body.String(); strings.Contains(body, "pq:") || strings.Contains(body, "/srv/db") {
		t.Errorf("error body leaked internals: %s", body)
	}
}

type failingError struct{}

func (f *failingError) Error() string { return `pq: relation "inquiries" does not exist at /srv/db` }

func TestClientKeyIgnoresForwardingHeadersByDefault(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/v1/inquiries", nil)
	request.RemoteAddr = "203.0.113.9:54321"
	request.Header.Set("X-Forwarded-For", "198.51.100.7")

	if got := clientKey(request, false); got != "203.0.113.9" {
		t.Errorf("clientKey(trust=false) = %q, want the peer address", got)
	}
}

func TestClientKeyUsesForwardedAddressWhenTrusted(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/v1/inquiries", nil)
	request.RemoteAddr = "10.0.0.1:54321"
	request.Header.Set("X-Forwarded-For", "198.51.100.7, 10.0.0.1")

	if got := clientKey(request, true); got != "198.51.100.7" {
		t.Errorf("clientKey(trust=true) = %q, want the forwarded client address", got)
	}
}

func TestClientKeyFallsBackToPeerWhenTrustedButHeaderMissing(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/v1/inquiries", nil)
	request.RemoteAddr = "10.0.0.1:54321"

	if got := clientKey(request, true); got != "10.0.0.1" {
		t.Errorf("clientKey(trust=true, no header) = %q, want the peer address", got)
	}
}

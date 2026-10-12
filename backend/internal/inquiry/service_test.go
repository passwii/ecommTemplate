package inquiry

import (
	"context"
	"errors"
	"testing"

	"github.com/passwii/ecommTemple/backend/internal/catalog"
)

type memoryRepository struct {
	values      map[string]Inquiry
	byRequestID map[string]Inquiry
	createCalls int
	failNext    error
}

func newMemoryRepository() *memoryRepository {
	return &memoryRepository{values: map[string]Inquiry{}, byRequestID: map[string]Inquiry{}}
}

func (m *memoryRepository) Create(_ context.Context, value Inquiry) (Inquiry, error) {
	m.createCalls++
	if m.failNext != nil {
		err := m.failNext
		m.failNext = nil
		return Inquiry{}, err
	}
	m.values[value.ID] = value
	if value.RequestID != "" {
		m.byRequestID[value.RequestID] = value
	}
	return value, nil
}

func (m *memoryRepository) FindByRequestID(_ context.Context, requestID string) (Inquiry, error) {
	if existing, ok := m.byRequestID[requestID]; ok {
		return existing, nil
	}
	return Inquiry{}, ErrNotFound
}

func (m *memoryRepository) Get(_ context.Context, id string) (Inquiry, error) {
	if existing, ok := m.values[id]; ok {
		return existing, nil
	}
	return Inquiry{}, ErrNotFound
}

func (m *memoryRepository) List(context.Context, ListFilter) ([]Inquiry, int, error) {
	return nil, 0, nil
}

func (m *memoryRepository) Update(context.Context, string, UpdateInput, string) (Inquiry, error) {
	return Inquiry{}, ErrNotFound
}

func (m *memoryRepository) ListActivities(context.Context, string, int) ([]Activity, error) {
	return nil, nil
}

func (m *memoryRepository) AddActivity(_ context.Context, activity Activity) (Activity, error) {
	return activity, nil
}

type stubProducts struct {
	product catalog.Product
	err     error
}

func (s stubProducts) GetPublishedProductByID(context.Context, string) (catalog.Product, error) {
	if s.err != nil {
		return catalog.Product{}, s.err
	}
	return s.product, nil
}

func TestServiceCreateStoresInquiry(t *testing.T) {
	repo := newMemoryRepository()
	products := stubProducts{product: catalog.Product{ID: "grill-01", Slug: "ridge-4-burner", Model: "Ridge 4-Burner", Name: "Ridge 4-Burner"}}
	service := NewService(repo, products)

	result, fieldErrors, err := service.Create(context.Background(), validInput())
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if fieldErrors != nil {
		t.Fatalf("fieldErrors = %v, want nil", fieldErrors)
	}
	if !result.Created {
		t.Error("Created = false, want true")
	}
	if result.Inquiry.Status != StatusNew {
		t.Errorf("status = %q, want new", result.Inquiry.Status)
	}
	if result.Inquiry.ProductModel != "Ridge 4-Burner" {
		t.Errorf("product model = %q, want snapshot from product record", result.Inquiry.ProductModel)
	}
}

func TestServiceCreateIsIdempotent(t *testing.T) {
	repo := newMemoryRepository()
	products := stubProducts{product: catalog.Product{ID: "grill-01", Slug: "ridge-4-burner", Name: "Ridge 4-Burner"}}
	service := NewService(repo, products)

	input := validInput()
	input.RequestID = "req-1"

	first, _, err := service.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("first Create() error = %v", err)
	}
	second, _, err := service.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("second Create() error = %v", err)
	}

	if second.Created {
		t.Error("second Created = true, want false for idempotent replay")
	}
	if second.Inquiry.ID != first.Inquiry.ID {
		t.Errorf("second id = %q, want %q", second.Inquiry.ID, first.Inquiry.ID)
	}
	if repo.createCalls != 1 {
		t.Errorf("createCalls = %d, want 1", repo.createCalls)
	}
}

func TestServiceCreateRejectsUnavailableProduct(t *testing.T) {
	repo := newMemoryRepository()
	service := NewService(repo, stubProducts{err: catalog.ErrNotFound})

	_, _, err := service.Create(context.Background(), validInput())
	if !errors.Is(err, ErrProductNotAvailable) {
		t.Fatalf("Create() error = %v, want ErrProductNotAvailable", err)
	}
	if repo.createCalls != 0 {
		t.Errorf("createCalls = %d, want 0", repo.createCalls)
	}
}

func TestServiceCreateReturnsFieldErrorsWithoutStoring(t *testing.T) {
	repo := newMemoryRepository()
	service := NewService(repo, stubProducts{product: catalog.Product{ID: "grill-01"}})

	input := validInput()
	input.ContactEmail = "invalid"

	_, fieldErrors, err := service.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if _, ok := fieldErrors["contactEmail"]; !ok {
		t.Fatalf("fieldErrors = %v, want contactEmail", fieldErrors)
	}
	if repo.createCalls != 0 {
		t.Errorf("createCalls = %d, want 0", repo.createCalls)
	}
}

func TestServiceCreateRecoversFromConcurrentIdempotentInsert(t *testing.T) {
	repo := newMemoryRepository()
	products := stubProducts{product: catalog.Product{ID: "grill-01", Slug: "ridge-4-burner", Name: "Ridge 4-Burner"}}
	service := NewService(repo, products)

	input := validInput()
	input.RequestID = "req-race"

	// Simulate another request having already stored this idempotency key, so the
	// insert fails with a unique violation.
	existing := Inquiry{ID: "inq_existing", RequestID: "req-race", Status: StatusNew}
	repo.byRequestID["req-race"] = existing
	repo.failNext = errors.New(`pq: duplicate key value violates unique constraint "inquiries_request_id_key"`)

	result, _, err := service.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if result.Created {
		t.Error("Created = true, want false")
	}
	if result.Inquiry.ID != "inq_existing" {
		t.Errorf("id = %q, want inq_existing", result.Inquiry.ID)
	}
}

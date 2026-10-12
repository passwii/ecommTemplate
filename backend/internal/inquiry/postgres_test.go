package inquiry

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// testDB returns a database handle scoped to a fresh schema so tests never
// depend on, or disturb, existing data in the development database.
func testDB(t *testing.T) *sql.DB {
	t.Helper()
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set; skipping inquiry integration test")
	}

	admin, err := sql.Open("pgx", url)
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}

	schema := fmt.Sprintf("inquiry_test_%d", time.Now().UnixNano())
	if _, err := admin.Exec("CREATE SCHEMA " + schema); err != nil {
		admin.Close()
		t.Fatalf("create test schema: %v", err)
	}

	for _, table := range []string{"collections", "products", "inquiries", "inquiry_activities"} {
		if _, err := admin.Exec(fmt.Sprintf("CREATE TABLE %s.%s (LIKE public.%s INCLUDING ALL)", schema, table, table)); err != nil {
			admin.Close()
			t.Fatalf("clone table %s: %v", table, err)
		}
	}

	separator := "?"
	if strings.Contains(url, "?") {
		separator = "&"
	}
	scoped, err := sql.Open("pgx", url+separator+"search_path="+schema)
	if err != nil {
		admin.Close()
		t.Fatalf("open scoped database: %v", err)
	}
	scoped.SetMaxOpenConns(1)

	t.Cleanup(func() {
		scoped.Close()
		_, _ = admin.Exec("DROP SCHEMA " + schema + " CASCADE")
		admin.Close()
	})
	return scoped
}

func seedProduct(t *testing.T, db *sql.DB) {
	t.Helper()
	if _, err := db.Exec(`INSERT INTO collections (slug, name, status, sort_order) VALUES ('freestanding', 'Freestanding grills', 'published', 1)`); err != nil {
		t.Fatalf("insert collection: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO products (id, slug, name, collection_id, status)
		VALUES ('grill-01', 'ridge-4-burner', 'Ridge 4-Burner', (SELECT id FROM collections WHERE slug = 'freestanding'), 'published')`); err != nil {
		t.Fatalf("insert product: %v", err)
	}
}

func sampleInquiry(id, requestID string) Inquiry {
	productID := "grill-01"
	quantity := 120
	return Inquiry{
		ID:                   id,
		RequestID:            requestID,
		ContactName:          "Alex Morgan",
		ContactEmail:         "buyer@example.com",
		ContactCompany:       "Example Importers",
		DestinationCountry:   "US",
		ProductID:            &productID,
		ProductSlug:          "ridge-4-burner",
		ProductModel:         "Ridge 4-Burner",
		ProductName:          "Ridge 4-Burner",
		Quantity:             &quantity,
		CustomizationRequest: "Private-label packaging options",
		PrivacyNoticeVersion: "inquiry-privacy-v1",
		PrivacyConsentedAt:   time.Date(2026, 10, 12, 12, 0, 0, 0, time.UTC),
		SourcePath:           "/products/ridge-4-burner",
		Status:               StatusNew,
	}
}

func TestPostgresRepositoryCreatePersistsInquiryAndActivity(t *testing.T) {
	db := testDB(t)
	seedProduct(t, db)
	repo := NewPostgresRepository(db)

	stored, err := repo.Create(context.Background(), sampleInquiry("inq_1", ""))
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if stored.Status != StatusNew {
		t.Errorf("status = %q, want new", stored.Status)
	}
	if stored.ProductModel != "Ridge 4-Burner" || stored.DestinationCountry != "US" {
		t.Errorf("stored inquiry = %+v", stored)
	}
	if stored.CreatedAt.IsZero() {
		t.Error("CreatedAt is zero")
	}

	var activityCount int
	if err := db.QueryRow(`SELECT count(*) FROM inquiry_activities WHERE inquiry_id = 'inq_1' AND activity_type = 'created'`).Scan(&activityCount); err != nil {
		t.Fatalf("count activities: %v", err)
	}
	if activityCount != 1 {
		t.Errorf("created activities = %d, want 1", activityCount)
	}
}

func TestPostgresRepositoryCreateRollsBackOnActivityFailure(t *testing.T) {
	db := testDB(t)
	seedProduct(t, db)
	repo := NewPostgresRepository(db)

	// A status outside the allowed set fails the CHECK constraint, so the whole
	// transaction must roll back and leave no inquiry behind.
	value := sampleInquiry("inq_bad", "")
	value.Status = Status("bogus")

	if _, err := repo.Create(context.Background(), value); err == nil {
		t.Fatal("Create() error = nil, want constraint violation")
	}

	var count int
	if err := db.QueryRow(`SELECT count(*) FROM inquiries`).Scan(&count); err != nil {
		t.Fatalf("count inquiries: %v", err)
	}
	if count != 0 {
		t.Errorf("inquiries after failed create = %d, want 0 (rolled back)", count)
	}
}

func TestPostgresRepositoryFindByRequestID(t *testing.T) {
	db := testDB(t)
	seedProduct(t, db)
	repo := NewPostgresRepository(db)

	if _, err := repo.Create(context.Background(), sampleInquiry("inq_2", "req-abc")); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	found, err := repo.FindByRequestID(context.Background(), "req-abc")
	if err != nil {
		t.Fatalf("FindByRequestID() error = %v", err)
	}
	if found.ID != "inq_2" {
		t.Errorf("id = %q, want inq_2", found.ID)
	}

	if _, err := repo.FindByRequestID(context.Background(), "missing"); !errors.Is(err, ErrNotFound) {
		t.Errorf("FindByRequestID(missing) error = %v, want ErrNotFound", err)
	}
	if _, err := repo.FindByRequestID(context.Background(), ""); !errors.Is(err, ErrNotFound) {
		t.Errorf("FindByRequestID(empty) error = %v, want ErrNotFound", err)
	}
}

func TestPostgresRepositoryGetReturnsNotFoundForUnknownID(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	if _, err := repo.Get(context.Background(), "inq_missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get() error = %v, want ErrNotFound", err)
	}
}

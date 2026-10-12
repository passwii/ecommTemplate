package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func testRepository(t *testing.T) *PostgresRepository {
	t.Helper()
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set; skipping catalog integration test")
	}
	db, err := sql.Open("pgx", url)
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		db.Close()
		t.Fatalf("BeginTx() error = %v", err)
	}
	t.Cleanup(func() {
		_ = tx.Rollback()
		_ = db.Close()
	})

	// Isolate each test in its own schema so results never depend on whether the
	// real database already contains seeded data.
	schema := fmt.Sprintf("catalog_test_%d", time.Now().UnixNano())
	if _, err := tx.ExecContext(context.Background(), "CREATE SCHEMA "+schema); err != nil {
		t.Fatalf("create test schema: %v", err)
	}
	if _, err := tx.ExecContext(context.Background(), "SET LOCAL search_path TO "+schema); err != nil {
		t.Fatalf("set search_path: %v", err)
	}
	for _, table := range []string{"storefront_metadata", "collections", "products", "company_content"} {
		statement := fmt.Sprintf("CREATE TABLE %s (LIKE public.%s INCLUDING ALL)", table, table)
		if _, err := tx.ExecContext(context.Background(), statement); err != nil {
			t.Fatalf("clone table %s: %v", table, err)
		}
	}
	return NewPostgresRepository(tx)
}

func insertCollection(t *testing.T, repo *PostgresRepository, slug, status string, sortOrder int) int64 {
	t.Helper()
	ctx := context.Background()
	var id int64
	err := repo.db.QueryRowContext(ctx,
		`INSERT INTO collections (slug, name, description, status, sort_order) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		slug, "Name "+slug, "Desc "+slug, status, sortOrder).Scan(&id)
	if err != nil {
		t.Fatalf("insert collection %s: %v", slug, err)
	}
	return id
}

func insertProduct(t *testing.T, repo *PostgresRepository, id, slug string, collectionID int64, status string, featured bool, sortOrder int, resourcesJSON string) {
	t.Helper()
	ctx := context.Background()
	_, err := repo.db.ExecContext(ctx,
		`INSERT INTO products (id, slug, name, collection_id, short_description, description, status, featured, sort_order, images, specifications, applications, resources, content_status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, '[]'::jsonb, '[]'::jsonb, '[]'::jsonb, $10::jsonb, 'sample')`,
		id, slug, "Name "+slug, collectionID, "short", "long", status, featured, sortOrder, resourcesJSON)
	if err != nil {
		t.Fatalf("insert product %s: %v", slug, err)
	}
}

func TestListPublishedCollectionsExcludesEmptyAndDraftCollections(t *testing.T) {
	repo := testRepository(t)
	published := insertCollection(t, repo, "published-in-use", "published", 1)
	insertCollection(t, repo, "published-empty", "published", 2)
	insertCollection(t, repo, "draft-in-use", "draft", 3)
	insertProduct(t, repo, "prod-a", "prod-a", published, "published", false, 1, "[]")
	insertProduct(t, repo, "prod-b", "prod-b", published, "draft", false, 2, "[]")

	got, err := repo.ListPublishedCollections(context.Background())
	if err != nil {
		t.Fatalf("ListPublishedCollections() error = %v", err)
	}
	if len(got) != 1 || got[0].Slug != "published-in-use" {
		t.Fatalf("ListPublishedCollections() = %+v, want only published-in-use", got)
	}
}

func TestGetPublishedCollectionReturnsNotFoundForEmptyOrUnknown(t *testing.T) {
	repo := testRepository(t)
	insertCollection(t, repo, "no-products", "published", 1)

	for _, slug := range []string{"no-products", "missing"} {
		_, err := repo.GetPublishedCollection(context.Background(), slug)
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("GetPublishedCollection(%q) error = %v, want ErrNotFound", slug, err)
		}
	}
}

func TestListPublishedProductsExcludesDraftsAndRespectsFilters(t *testing.T) {
	repo := testRepository(t)
	collection := insertCollection(t, repo, "grills", "published", 1)
	other := insertCollection(t, repo, "compact", "published", 2)
	insertProduct(t, repo, "p-1", "p-1", collection, "published", true, 1, "[]")
	insertProduct(t, repo, "p-2", "p-2", collection, "published", false, 2, "[]")
	insertProduct(t, repo, "p-3", "p-3", other, "published", true, 3, "[]")
	insertProduct(t, repo, "p-4", "p-4", collection, "draft", true, 4, "[]")

	all, err := repo.ListPublishedProducts(context.Background(), Filter{Limit: 50})
	if err != nil {
		t.Fatalf("ListPublishedProducts() error = %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("unfiltered count = %d, want 3 (draft excluded)", len(all))
	}

	byCollection, err := repo.ListPublishedProducts(context.Background(), Filter{CollectionSlug: "grills", Limit: 50})
	if err != nil {
		t.Fatalf("ListPublishedProducts(collection) error = %v", err)
	}
	if len(byCollection) != 2 {
		t.Fatalf("collection filter count = %d, want 2", len(byCollection))
	}

	featured := true
	featuredOnly, err := repo.ListPublishedProducts(context.Background(), Filter{Featured: &featured, Limit: 50})
	if err != nil {
		t.Fatalf("ListPublishedProducts(featured) error = %v", err)
	}
	if len(featuredOnly) != 2 {
		t.Fatalf("featured filter count = %d, want 2", len(featuredOnly))
	}

	limited, err := repo.ListPublishedProducts(context.Background(), Filter{Limit: 1})
	if err != nil {
		t.Fatalf("ListPublishedProducts(limit) error = %v", err)
	}
	if len(limited) != 1 {
		t.Fatalf("limit count = %d, want 1", len(limited))
	}
}

func TestGetPublishedProductReturnsNotFoundForDraftAndUnknown(t *testing.T) {
	repo := testRepository(t)
	collection := insertCollection(t, repo, "grills", "published", 1)
	insertProduct(t, repo, "draft-1", "draft-1", collection, "draft", false, 1, "[]")

	for _, slug := range []string{"draft-1", "missing"} {
		_, err := repo.GetPublishedProduct(context.Background(), slug)
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("GetPublishedProduct(%q) error = %v, want ErrNotFound", slug, err)
		}
	}
}

func TestPublishedProductDropsUnsafeResourceURLs(t *testing.T) {
	repo := testRepository(t)
	collection := insertCollection(t, repo, "grills", "published", 1)
	resources := `[{"title":"ok","url":"/files/spec.pdf"},{"title":"bad","url":"javascript:alert(1)"},{"title":"proto","url":"//evil.example.com/x.pdf"}]`
	insertProduct(t, repo, "p-1", "p-1", collection, "published", false, 1, resources)

	got, err := repo.GetPublishedProduct(context.Background(), "p-1")
	if err != nil {
		t.Fatalf("GetPublishedProduct() error = %v", err)
	}
	if len(got.Resources) != 1 || got.Resources[0].URL != "/files/spec.pdf" {
		t.Fatalf("Resources = %+v, want only the safe /files/spec.pdf entry", got.Resources)
	}
}

func TestGetCompanyContentReturnsNotFoundWhenUnseeded(t *testing.T) {
	repo := testRepository(t)
	_, err := repo.GetCompanyContent(context.Background())
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetCompanyContent() error = %v, want ErrNotFound", err)
	}
}

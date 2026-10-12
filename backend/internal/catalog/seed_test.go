package catalog

import (
	"context"
	"strings"
	"testing"
)

func TestSampleDataIsValid(t *testing.T) {
	if err := ValidateSeedData(SampleData()); err != nil {
		t.Fatalf("ValidateSeedData(SampleData()) error = %v", err)
	}
}

func TestSampleDataDoesNotFabricateFacts(t *testing.T) {
	data := SampleData()
	if data.Company.ContentStatus != "sample" {
		t.Errorf("company content status = %q, want sample", data.Company.ContentStatus)
	}
	if len(data.Company.Credentials) != 0 {
		t.Errorf("sample company content must not present credentials as facts, got %d", len(data.Company.Credentials))
	}
	for _, product := range data.Products {
		if product.ContentStatus != "sample" {
			t.Errorf("product %q content status = %q, want sample", product.Slug, product.ContentStatus)
		}
		if product.MinimumOrderQuantity != "" || product.LeadTime != "" {
			t.Errorf("product %q must not assert MOQ or lead time in sample data", product.Slug)
		}
	}
}

func TestValidateSeedDataRejectsBrokenReferences(t *testing.T) {
	base := SampleData()

	t.Run("duplicate product slug", func(t *testing.T) {
		data := SampleData()
		data.Products[1].Slug = data.Products[0].Slug
		if err := ValidateSeedData(data); err == nil {
			t.Fatal("expected duplicate slug error")
		}
	})

	t.Run("unknown collection reference", func(t *testing.T) {
		data := SampleData()
		data.Products[0].CollectionSlug = "does-not-exist"
		if err := ValidateSeedData(data); err == nil {
			t.Fatal("expected unknown collection error")
		}
	})

	t.Run("unsafe resource url", func(t *testing.T) {
		data := SampleData()
		data.Products[0].Resources = []ProductResource{{Title: "bad", URL: "javascript:alert(1)"}}
		if err := ValidateSeedData(data); err == nil {
			t.Fatal("expected unsafe URL error")
		}
	})

	t.Run("missing template id", func(t *testing.T) {
		data := SampleData()
		data.TemplateID = ""
		if err := ValidateSeedData(data); err == nil {
			t.Fatal("expected missing template id error")
		}
	})

	t.Run("unverified credentials", func(t *testing.T) {
		data := base
		data.Company.ContentStatus = "sample"
		data.Company.Credentials = []CompanyCredential{{Name: "CE", Issuer: "unknown", Details: "unverified"}}
		if err := ValidateSeedData(data); err == nil {
			t.Fatal("expected unverified credential error")
		}
	})
}

func TestSeedIsIdempotent(t *testing.T) {
	repo := testRepository(t)
	ctx := context.Background()

	if err := Seed(ctx, repo.db, SampleData()); err != nil {
		t.Fatalf("first Seed() error = %v", err)
	}
	if err := Seed(ctx, repo.db, SampleData()); err != nil {
		t.Fatalf("second Seed() error = %v", err)
	}

	var productCount int
	if err := repo.db.QueryRowContext(ctx, `SELECT count(*) FROM products`).Scan(&productCount); err != nil {
		t.Fatalf("count products: %v", err)
	}
	if productCount != len(SampleData().Products) {
		t.Fatalf("product count after reseed = %d, want %d", productCount, len(SampleData().Products))
	}

	product, err := repo.GetPublishedProduct(ctx, "ridge-4-burner")
	if err != nil {
		t.Fatalf("GetPublishedProduct() error = %v", err)
	}
	if product.Name != "Ridge 4-Burner" {
		t.Errorf("seeded product name = %q", product.Name)
	}
	if !strings.Contains(product.Collection.Slug, "freestanding") {
		t.Errorf("seeded product collection = %q, want freestanding", product.Collection.Slug)
	}
}

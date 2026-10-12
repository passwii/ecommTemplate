package catalog

import (
	"context"
	"encoding/json"
	"fmt"
)

// SeedCollection is one collection row in the sample dataset.
type SeedCollection struct {
	Slug        string
	Name        string
	Description string
	Status      string
	SortOrder   int
}

// SeedProduct is one product row in the sample dataset.
type SeedProduct struct {
	ID                   string
	Slug                 string
	Name                 string
	Model                string
	SKU                  string
	CollectionSlug       string
	ShortDescription     string
	Description          string
	PriceAmount          *float64
	PriceCurrency        string
	CompareAtPrice       *float64
	Status               string
	ContentStatus        string
	Featured             bool
	SortOrder            int
	Badge                string
	Fuel                 string
	Burners              *int
	CookingArea          string
	Finish               string
	MinimumOrderQuantity string
	LeadTime             string
	Applications         []string
	Images               []ProductImage
	Specifications       []ProductSpecification
	Resources            []ProductResource
}

// SeedData is the full sample dataset for the gas-grill storefront.
type SeedData struct {
	TemplateID  string
	Collections []SeedCollection
	Products    []SeedProduct
	Company     CompanyContent
}

func float64Ptr(value float64) *float64 { return &value }
func intPtr(value int) *int             { return &value }

// SampleData returns the reviewed sample dataset. It mirrors the current
// storefront mock data and is explicitly marked as sample content; it does not
// assert unverified certifications, test results, factory capacity, or MOQ.
func SampleData() SeedData {
	return SeedData{
		TemplateID: "gas-grill",
		Collections: []SeedCollection{
			{Slug: "built-in", Name: "Built-in grills", Description: "A considered centerpiece for an outdoor kitchen that feels like part of home.", Status: "published", SortOrder: 10},
			{Slug: "freestanding", Name: "Freestanding grills", Description: "A complete cooking station, ready to make the patio your favorite room.", Status: "published", SortOrder: 20},
			{Slug: "compact", Name: "Compact grills", Description: "Full-flavor cooking, thoughtfully sized for smaller outdoor spaces.", Status: "published", SortOrder: 30},
		},
		Products: []SeedProduct{
			{
				ID: "grill-01", Slug: "ridge-4-burner", Name: "Ridge 4-Burner", Model: "Ridge 4-Burner",
				CollectionSlug: "freestanding", ShortDescription: "Room to cook for the whole table.",
				Description: "A capable centerpiece for generous weekends outside. Four independently controlled burners give you the range to sear, roast and keep a side warm, all at once.",
				PriceAmount: float64Ptr(2499), PriceCurrency: "USD",
				Status: "published", ContentStatus: "sample", Featured: true, SortOrder: 10,
				Badge: "Bestseller", Fuel: "Natural gas or propane", Burners: intPtr(4), CookingArea: "748 sq in", Finish: "Brushed stainless steel",
				Applications: []string{},
				Images:       []ProductImage{{URL: "/images/grill-ridge.jpg", Alt: "Patio lounge seating around an outdoor fire pit beneath a covered, tree-lined terrace", SortOrder: 0}},
				Specifications: []ProductSpecification{
					{Name: "Burners", Value: "4"},
					{Name: "Cooking area", Value: "748 sq in"},
					{Name: "Fuel", Value: "Natural gas or propane"},
					{Name: "Finish", Value: "Brushed stainless steel"},
				},
				Resources: []ProductResource{},
			},
			{
				ID: "grill-02", Slug: "hearth-3-burner", Name: "Hearth 3-Burner", Model: "Hearth 3-Burner",
				CollectionSlug: "freestanding", ShortDescription: "Everyday cooking, made effortless.",
				Description: "The grill you reach for on a Tuesday and gather around on Saturday. Three responsive burners make reliable heat feel simple.",
				PriceAmount: float64Ptr(1899), PriceCurrency: "USD",
				Status: "published", ContentStatus: "sample", Featured: false, SortOrder: 20,
				Fuel: "Propane", Burners: intPtr(3), CookingArea: "560 sq in", Finish: "Graphite and stainless steel",
				Applications: []string{},
				Images:       []ProductImage{{URL: "/images/grill-hearth.jpg", Alt: "Warm wood-paneled cabin kitchen with a dining area beside large windows", SortOrder: 0}},
				Specifications: []ProductSpecification{
					{Name: "Burners", Value: "3"},
					{Name: "Cooking area", Value: "560 sq in"},
					{Name: "Fuel", Value: "Propane"},
					{Name: "Finish", Value: "Graphite and stainless steel"},
				},
				Resources: []ProductResource{},
			},
			{
				ID: "grill-03", Slug: "summit-built-in", Name: "Summit Built-In", Model: "Summit Built-In",
				CollectionSlug: "built-in", ShortDescription: "Made to belong in your outdoor kitchen.",
				Description: "Clean lines and steady, even heat bring a thoughtful finish to the outdoor kitchen. Designed to sit flush with the space around it.",
				PriceAmount: float64Ptr(3299), PriceCurrency: "USD",
				Status: "published", ContentStatus: "sample", Featured: true, SortOrder: 30,
				Badge: "New", Fuel: "Natural gas or propane", Burners: intPtr(5), CookingArea: "860 sq in", Finish: "Brushed stainless steel",
				Applications: []string{},
				Images:       []ProductImage{{URL: "/images/grill-summit.jpg", Alt: "Two people seated outdoors beside a fire pit beneath string lights", SortOrder: 0}},
				Specifications: []ProductSpecification{
					{Name: "Burners", Value: "5"},
					{Name: "Cooking area", Value: "860 sq in"},
					{Name: "Fuel", Value: "Natural gas or propane"},
					{Name: "Finish", Value: "Brushed stainless steel"},
				},
				Resources: []ProductResource{},
			},
			{
				ID: "grill-04", Slug: "pocket-2-burner", Name: "Pocket 2-Burner", Model: "Pocket 2-Burner",
				CollectionSlug: "compact", ShortDescription: "A smaller footprint. No smaller appetite.",
				Description: "Thoughtfully sized for balconies, courtyards and first homes. Two burners give you direct and indirect heat without taking over the space.",
				PriceAmount: float64Ptr(1099), PriceCurrency: "USD",
				Status: "published", ContentStatus: "sample", Featured: false, SortOrder: 40,
				Fuel: "Propane", Burners: intPtr(2), CookingArea: "360 sq in", Finish: "Matte graphite",
				Applications: []string{},
				Images:       []ProductImage{{URL: "/images/grill-pocket.jpg", Alt: "Person cooking at an outdoor grill in a leafy garden", SortOrder: 0}},
				Specifications: []ProductSpecification{
					{Name: "Burners", Value: "2"},
					{Name: "Cooking area", Value: "360 sq in"},
					{Name: "Fuel", Value: "Propane"},
					{Name: "Finish", Value: "Matte graphite"},
				},
				Resources: []ProductResource{},
			},
		},
		Company: CompanyContent{
			About: CompanyStory{
				Heading:  "About this storefront demo",
				Body:     "This is sample content for the gas-grill storefront template. Replace it with verified company information before publishing to customers.",
				ImageURL: "",
				ImageAlt: "",
			},
			Manufacturing: CompanyManufacturing{
				Heading:        "Manufacturing and quality (sample)",
				Overview:       "Sample placeholder. Add only verified manufacturing information, with evidence, before presenting it as fact.",
				QualityHeading: "Quality process (sample)",
				QualityDetails: "Sample placeholder. Describe your real quality process only when it can be supported by evidence.",
			},
			Credentials:   []CompanyCredential{},
			CaseStudies:   []CompanyCaseStudy{},
			FAQs:          []CompanyFAQ{},
			ContentStatus: "sample",
		},
	}
}

// ValidateSeedData checks referential and safety invariants of the dataset
// before it is written.
func ValidateSeedData(data SeedData) error {
	if data.TemplateID == "" {
		return fmt.Errorf("seed data: template id is required")
	}

	collectionSlugs := make(map[string]bool, len(data.Collections))
	for _, collection := range data.Collections {
		if collection.Slug == "" {
			return fmt.Errorf("seed data: collection slug is required")
		}
		if collectionSlugs[collection.Slug] {
			return fmt.Errorf("seed data: duplicate collection slug %q", collection.Slug)
		}
		collectionSlugs[collection.Slug] = true
		if collection.Status != "published" && collection.Status != "draft" {
			return fmt.Errorf("seed data: collection %q has invalid status %q", collection.Slug, collection.Status)
		}
	}

	productIDs := make(map[string]bool, len(data.Products))
	productSlugs := make(map[string]bool, len(data.Products))
	for _, product := range data.Products {
		if product.ID == "" || product.Slug == "" {
			return fmt.Errorf("seed data: product id and slug are required")
		}
		if productIDs[product.ID] {
			return fmt.Errorf("seed data: duplicate product id %q", product.ID)
		}
		if productSlugs[product.Slug] {
			return fmt.Errorf("seed data: duplicate product slug %q", product.Slug)
		}
		productIDs[product.ID] = true
		productSlugs[product.Slug] = true
		if !collectionSlugs[product.CollectionSlug] {
			return fmt.Errorf("seed data: product %q references unknown collection %q", product.Slug, product.CollectionSlug)
		}
		if product.Status != "published" && product.Status != "draft" {
			return fmt.Errorf("seed data: product %q has invalid status %q", product.Slug, product.Status)
		}
		if product.ContentStatus != "sample" && product.ContentStatus != "verified" {
			return fmt.Errorf("seed data: product %q has invalid content status %q", product.Slug, product.ContentStatus)
		}
		for _, image := range product.Images {
			if !SafeResourceURL(image.URL) {
				return fmt.Errorf("seed data: product %q has unsafe image URL %q", product.Slug, image.URL)
			}
		}
		for _, resource := range product.Resources {
			if !SafeResourceURL(resource.URL) {
				return fmt.Errorf("seed data: product %q has unsafe resource URL %q", product.Slug, resource.URL)
			}
		}
	}

	if data.Company.ContentStatus != "sample" && data.Company.ContentStatus != "verified" {
		return fmt.Errorf("seed data: company content has invalid content status %q", data.Company.ContentStatus)
	}
	if data.Company.About.ImageURL != "" && !SafeResourceURL(data.Company.About.ImageURL) {
		return fmt.Errorf("seed data: company about image URL is unsafe")
	}
	// The sample dataset must not fabricate certifications or customer claims.
	if len(data.Company.Credentials) > 0 && data.Company.ContentStatus != "verified" {
		return fmt.Errorf("seed data: unverified content must not present credentials as facts")
	}
	return nil
}

// Seed writes the sample dataset into the database idempotently. Re-running it
// updates existing rows instead of duplicating them.
func Seed(ctx context.Context, db DBTX, data SeedData) error {
	if err := ValidateSeedData(data); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, `
		INSERT INTO storefront_metadata (id, template_id, updated_at)
		VALUES (1, $1, now())
		ON CONFLICT (id) DO UPDATE SET template_id = EXCLUDED.template_id, updated_at = now()`, data.TemplateID); err != nil {
		return fmt.Errorf("seed storefront metadata: %w", err)
	}

	collectionIDs := make(map[string]int64, len(data.Collections))
	for _, collection := range data.Collections {
		var id int64
		err := db.QueryRowContext(ctx, `
			INSERT INTO collections (slug, name, description, status, sort_order, updated_at)
			VALUES ($1, $2, $3, $4, $5, now())
			ON CONFLICT (slug) DO UPDATE SET
				name = EXCLUDED.name,
				description = EXCLUDED.description,
				status = EXCLUDED.status,
				sort_order = EXCLUDED.sort_order,
				updated_at = now()
			RETURNING id`,
			collection.Slug, collection.Name, collection.Description, collection.Status, collection.SortOrder).Scan(&id)
		if err != nil {
			return fmt.Errorf("seed collection %q: %w", collection.Slug, err)
		}
		collectionIDs[collection.Slug] = id
	}

	for _, product := range data.Products {
		images, err := json.Marshal(product.Images)
		if err != nil {
			return fmt.Errorf("encode images for %q: %w", product.Slug, err)
		}
		specifications, err := json.Marshal(product.Specifications)
		if err != nil {
			return fmt.Errorf("encode specifications for %q: %w", product.Slug, err)
		}
		applications, err := json.Marshal(product.Applications)
		if err != nil {
			return fmt.Errorf("encode applications for %q: %w", product.Slug, err)
		}
		resources, err := json.Marshal(product.Resources)
		if err != nil {
			return fmt.Errorf("encode resources for %q: %w", product.Slug, err)
		}

		_, err = db.ExecContext(ctx, `
			INSERT INTO products (
				id, slug, name, model, sku, collection_id, short_description, description,
				price_amount, price_currency, compare_at_price, status, content_status,
				featured, sort_order, badge, fuel, burners, cooking_area, finish,
				minimum_order_quantity, lead_time, applications, images, specifications, resources, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8,
				$9, $10, $11, $12, $13,
				$14, $15, $16, $17, $18, $19, $20,
				$21, $22, $23::jsonb, $24::jsonb, $25::jsonb, $26::jsonb, now()
			)
			ON CONFLICT (id) DO UPDATE SET
				slug = EXCLUDED.slug,
				name = EXCLUDED.name,
				model = EXCLUDED.model,
				sku = EXCLUDED.sku,
				collection_id = EXCLUDED.collection_id,
				short_description = EXCLUDED.short_description,
				description = EXCLUDED.description,
				price_amount = EXCLUDED.price_amount,
				price_currency = EXCLUDED.price_currency,
				compare_at_price = EXCLUDED.compare_at_price,
				status = EXCLUDED.status,
				content_status = EXCLUDED.content_status,
				featured = EXCLUDED.featured,
				sort_order = EXCLUDED.sort_order,
				badge = EXCLUDED.badge,
				fuel = EXCLUDED.fuel,
				burners = EXCLUDED.burners,
				cooking_area = EXCLUDED.cooking_area,
				finish = EXCLUDED.finish,
				minimum_order_quantity = EXCLUDED.minimum_order_quantity,
				lead_time = EXCLUDED.lead_time,
				applications = EXCLUDED.applications,
				images = EXCLUDED.images,
				specifications = EXCLUDED.specifications,
				resources = EXCLUDED.resources,
				updated_at = now()`,
			product.ID, product.Slug, product.Name, product.Model, product.SKU, collectionIDs[product.CollectionSlug],
			product.ShortDescription, product.Description,
			product.PriceAmount, product.PriceCurrency, product.CompareAtPrice, product.Status, product.ContentStatus,
			product.Featured, product.SortOrder, product.Badge, product.Fuel, product.Burners, product.CookingArea, product.Finish,
			product.MinimumOrderQuantity, product.LeadTime, string(applications), string(images), string(specifications), string(resources))
		if err != nil {
			return fmt.Errorf("seed product %q: %w", product.Slug, err)
		}
	}

	about, err := json.Marshal(data.Company.About)
	if err != nil {
		return fmt.Errorf("encode company about: %w", err)
	}
	manufacturing, err := json.Marshal(data.Company.Manufacturing)
	if err != nil {
		return fmt.Errorf("encode company manufacturing: %w", err)
	}
	credentials, err := json.Marshal(data.Company.Credentials)
	if err != nil {
		return fmt.Errorf("encode company credentials: %w", err)
	}
	caseStudies, err := json.Marshal(data.Company.CaseStudies)
	if err != nil {
		return fmt.Errorf("encode company case studies: %w", err)
	}
	faqs, err := json.Marshal(data.Company.FAQs)
	if err != nil {
		return fmt.Errorf("encode company faqs: %w", err)
	}

	if _, err := db.ExecContext(ctx, `
		INSERT INTO company_content (id, about, manufacturing, credentials, case_studies, faqs, content_status, updated_at)
		VALUES (1, $1::jsonb, $2::jsonb, $3::jsonb, $4::jsonb, $5::jsonb, $6, now())
		ON CONFLICT (id) DO UPDATE SET
			about = EXCLUDED.about,
			manufacturing = EXCLUDED.manufacturing,
			credentials = EXCLUDED.credentials,
			case_studies = EXCLUDED.case_studies,
			faqs = EXCLUDED.faqs,
			content_status = EXCLUDED.content_status,
			updated_at = now()`,
		string(about), string(manufacturing), string(credentials), string(caseStudies), string(faqs), data.Company.ContentStatus); err != nil {
		return fmt.Errorf("seed company content: %w", err)
	}

	return nil
}

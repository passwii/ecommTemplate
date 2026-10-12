package catalog

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	defaultLimit = 24
	maxLimit     = 100
)

// PostgresRepository is the SQL-backed implementation of Repository.
type PostgresRepository struct {
	db DBTX
}

// NewPostgresRepository builds a read repository over the given database handle.
func NewPostgresRepository(db DBTX) *PostgresRepository {
	return &PostgresRepository{db: db}
}

const collectionColumns = `c.id, c.slug, c.name, c.description, c.sort_order`

func (r *PostgresRepository) ListPublishedCollections(ctx context.Context) ([]Collection, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT `+collectionColumns+`
		FROM collections c
		WHERE c.status = 'published'
		  AND EXISTS (
			SELECT 1 FROM products p
			WHERE p.collection_id = c.id AND p.status = 'published'
		  )
		ORDER BY c.sort_order, c.id`)
	if err != nil {
		return nil, fmt.Errorf("list published collections: %w", err)
	}
	defer rows.Close()

	collections := []Collection{}
	for rows.Next() {
		collection, err := scanCollection(rows)
		if err != nil {
			return nil, err
		}
		collections = append(collections, collection)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list published collections: %w", err)
	}
	return collections, nil
}

func (r *PostgresRepository) GetPublishedCollection(ctx context.Context, slug string) (Collection, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT `+collectionColumns+`
		FROM collections c
		WHERE c.status = 'published'
		  AND c.slug = $1
		  AND EXISTS (
			SELECT 1 FROM products p
			WHERE p.collection_id = c.id AND p.status = 'published'
		  )`, slug)

	collection, err := scanCollection(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Collection{}, ErrNotFound
	}
	if err != nil {
		return Collection{}, err
	}
	return collection, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanCollection(row rowScanner) (Collection, error) {
	var collection Collection
	if err := row.Scan(&collection.ID, &collection.Slug, &collection.Name, &collection.Description, &collection.SortOrder); err != nil {
		return Collection{}, err
	}
	return collection, nil
}

const productColumns = `
	p.id, p.slug, p.name, p.model, p.sku,
	c.id, c.slug, c.name, c.description, c.sort_order,
	p.short_description, p.description,
	p.images, p.specifications, p.applications, p.resources,
	p.price_amount::float8, p.price_currency, p.compare_at_price::float8,
	p.badge, p.fuel, p.burners, p.cooking_area, p.finish,
	p.minimum_order_quantity, p.lead_time, p.featured, p.sort_order, p.content_status`

func (r *PostgresRepository) ListPublishedProducts(ctx context.Context, filter Filter) ([]Product, error) {
	limit := clampLimit(filter.Limit)
	collectionSlug := filter.CollectionSlug
	featuredOnly := filter.Featured != nil && *filter.Featured

	rows, err := r.db.QueryContext(ctx, `
		SELECT `+productColumns+`
		FROM products p
		JOIN collections c ON c.id = p.collection_id
		WHERE p.status = 'published'
		  AND ($1 = '' OR c.slug = $1)
		  AND (NOT $2 OR p.featured)
		ORDER BY p.sort_order, p.id
		LIMIT $3`, collectionSlug, featuredOnly, limit)
	if err != nil {
		return nil, fmt.Errorf("list published products: %w", err)
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list published products: %w", err)
	}
	return products, nil
}

func (r *PostgresRepository) GetPublishedProduct(ctx context.Context, slug string) (Product, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT `+productColumns+`
		FROM products p
		JOIN collections c ON c.id = p.collection_id
		WHERE p.status = 'published' AND p.slug = $1`, slug)

	product, err := scanProduct(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Product{}, ErrNotFound
	}
	if err != nil {
		return Product{}, err
	}
	return product, nil
}

func scanProduct(row rowScanner) (Product, error) {
	var (
		product               Product
		imagesJSON            []byte
		specificationsJSON    []byte
		applicationsJSON      []byte
		resourcesJSON         []byte
		priceAmount           sql.NullFloat64
		compareAtPrice        sql.NullFloat64
		burners               sql.NullInt64
		collectionID          int64
		collectionSlug        string
		collectionName        string
		collectionDescription string
		collectionSortOrder   int
	)

	err := row.Scan(
		&product.ID, &product.Slug, &product.Name, &product.Model, &product.SKU,
		&collectionID, &collectionSlug, &collectionName, &collectionDescription, &collectionSortOrder,
		&product.ShortDescription, &product.Description,
		&imagesJSON, &specificationsJSON, &applicationsJSON, &resourcesJSON,
		&priceAmount, &product.PriceCurrency, &compareAtPrice,
		&product.Badge, &product.Fuel, &burners, &product.CookingArea, &product.Finish,
		&product.MinimumOrderQuantity, &product.LeadTime, &product.Featured, &product.SortOrder, &product.ContentStatus,
	)
	if err != nil {
		return Product{}, err
	}

	product.Collection = Collection{
		ID:          collectionID,
		Slug:        collectionSlug,
		Name:        collectionName,
		Description: collectionDescription,
		SortOrder:   collectionSortOrder,
	}

	product.Images = []ProductImage{}
	if err := json.Unmarshal(imagesJSON, &product.Images); err != nil {
		return Product{}, fmt.Errorf("decode product images: %w", err)
	}
	product.Images = safeImages(product.Images)

	product.Specifications = []ProductSpecification{}
	if err := json.Unmarshal(specificationsJSON, &product.Specifications); err != nil {
		return Product{}, fmt.Errorf("decode product specifications: %w", err)
	}

	product.Applications = []string{}
	if err := json.Unmarshal(applicationsJSON, &product.Applications); err != nil {
		return Product{}, fmt.Errorf("decode product applications: %w", err)
	}

	product.Resources = []ProductResource{}
	if err := json.Unmarshal(resourcesJSON, &product.Resources); err != nil {
		return Product{}, fmt.Errorf("decode product resources: %w", err)
	}
	product.Resources = safeResources(product.Resources)

	if priceAmount.Valid {
		value := priceAmount.Float64
		product.PriceAmount = &value
	}
	if compareAtPrice.Valid {
		value := compareAtPrice.Float64
		product.CompareAtPrice = &value
	}
	if burners.Valid {
		value := int(burners.Int64)
		product.Burners = &value
	}
	return product, nil
}

func safeImages(images []ProductImage) []ProductImage {
	safe := make([]ProductImage, 0, len(images))
	for _, image := range images {
		if SafeResourceURL(image.URL) {
			safe = append(safe, image)
		}
	}
	return safe
}

func safeResources(resources []ProductResource) []ProductResource {
	safe := make([]ProductResource, 0, len(resources))
	for _, resource := range resources {
		if SafeResourceURL(resource.URL) {
			safe = append(safe, resource)
		}
	}
	return safe
}

func clampLimit(limit int) int {
	if limit <= 0 {
		return defaultLimit
	}
	if limit > maxLimit {
		return maxLimit
	}
	return limit
}

func (r *PostgresRepository) GetCompanyContent(ctx context.Context) (CompanyContent, error) {
	var (
		content       CompanyContent
		aboutJSON     []byte
		manufactureJS []byte
		credentialsJS []byte
		caseStudiesJS []byte
		faqsJSON      []byte
	)

	row := r.db.QueryRowContext(ctx, `
		SELECT about, manufacturing, credentials, case_studies, faqs, content_status
		FROM company_content
		WHERE id = 1`)

	err := row.Scan(&aboutJSON, &manufactureJS, &credentialsJS, &caseStudiesJS, &faqsJSON, &content.ContentStatus)
	if errors.Is(err, sql.ErrNoRows) {
		return CompanyContent{}, ErrNotFound
	}
	if err != nil {
		return CompanyContent{}, fmt.Errorf("read company content: %w", err)
	}

	if err := json.Unmarshal(aboutJSON, &content.About); err != nil {
		return CompanyContent{}, fmt.Errorf("decode company about: %w", err)
	}
	if err := json.Unmarshal(manufactureJS, &content.Manufacturing); err != nil {
		return CompanyContent{}, fmt.Errorf("decode company manufacturing: %w", err)
	}
	content.Credentials = []CompanyCredential{}
	if err := json.Unmarshal(credentialsJS, &content.Credentials); err != nil {
		return CompanyContent{}, fmt.Errorf("decode company credentials: %w", err)
	}
	content.CaseStudies = []CompanyCaseStudy{}
	if err := json.Unmarshal(caseStudiesJS, &content.CaseStudies); err != nil {
		return CompanyContent{}, fmt.Errorf("decode company case studies: %w", err)
	}
	content.FAQs = []CompanyFAQ{}
	if err := json.Unmarshal(faqsJSON, &content.FAQs); err != nil {
		return CompanyContent{}, fmt.Errorf("decode company faqs: %w", err)
	}
	if content.About.ImageURL != "" && !SafeResourceURL(content.About.ImageURL) {
		content.About.ImageURL = ""
	}
	return content, nil
}

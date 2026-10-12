package catalog

import (
	"context"
	"database/sql"
	"errors"
)

// ErrNotFound is returned when a published record does not exist.
var ErrNotFound = errors.New("catalog: not found")

// DBTX is satisfied by both *sql.DB and *sql.Tx, so the repository can run
// inside a transaction for tests.
type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// Filter narrows a published product list query.
type Filter struct {
	CollectionSlug string
	Featured       *bool
	Limit          int
}

// Repository reads public catalog and company content.
type Repository interface {
	ListPublishedCollections(ctx context.Context) ([]Collection, error)
	GetPublishedCollection(ctx context.Context, slug string) (Collection, error)
	ListPublishedProducts(ctx context.Context, filter Filter) ([]Product, error)
	GetPublishedProduct(ctx context.Context, slug string) (Product, error)
	GetCompanyContent(ctx context.Context) (CompanyContent, error)
}

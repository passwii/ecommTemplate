# Catalog and Company Content API Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Store the `gas-grill` sample catalog and company content in the dedicated PostgreSQL database and expose it through public, read-only, versioned JSON endpoints.

**Architecture:** Add a `catalog` package that owns read-only domain models and a PostgreSQL repository, plus HTTP handlers under `/api/v1`. Seed data is loaded from an explicit, idempotent Go seed command that mirrors the current storefront mock data. Draft records never appear in public responses. Persistence stays behind a repository interface so handlers and future callers do not depend on SQL details.

**Tech Stack:** Go 1.25 minimum / Go 1.27.1 toolchain, `net/http`, `database/sql` with pgx v5, Goose v3 embedded migrations, PostgreSQL 18.

## Global Constraints

- `template_id: gas-grill`; this extends the existing template and adds no new visual style.
- Each store owns an independent API and PostgreSQL database; no shared multi-tenant schema.
- Public read responses only include published records; drafts and unknown slugs return `404` without revealing draft existence.
- Never expose SQL, stack traces, internal file paths, or database credentials in error responses.
- Reuse the dedicated `ecommtemple_dev` local database only; never touch other local databases.
- Catalog seed content is sample data. It must not assert unverified certifications, test results, factory capacity, lead times, or MOQ as facts.
- No checkout, payment processing, or third-party commerce plugins.
- Follow `AGENTS.md` and `TESTING.md`.

---

## API Contract (this phase)

Envelope: `{"data": ...}` for success, `{"error": {"code": "...", "message": "...", "fields": {...}}}` for errors.

- `GET /api/v1/collections` → published collections that have at least one published product.
- `GET /api/v1/collections/{slug}` → one published, non-empty collection or `404`.
- `GET /api/v1/products?collection={slug}&featured={bool}&limit={n}` → bounded published list.
- `GET /api/v1/products/{slug}` → one published product or `404`.
- `GET /api/v1/company-content` → published company content (facts and sample flagged separately).

Error codes: `not_found`, `invalid_query`, `internal_error`.

## File Structure

- `backend/cmd/storefront-migrate/migrations/00002_create_catalog.sql` — catalog and company schema.
- `backend/cmd/storefront-seed/main.go` — idempotent seed command (only writes to the guarded project DB).
- `backend/internal/catalog/models.go` — domain models and DTO mapping.
- `backend/internal/catalog/url.go` — safe resource URL validation mirroring the storefront rule.
- `backend/internal/catalog/repository.go` — repository interface.
- `backend/internal/catalog/postgres.go` — SQL implementation.
- `backend/internal/catalog/seed.go` — seed dataset (sample content, explicitly marked).
- `backend/internal/httpapi/errors.go` — JSON envelope + error helpers.
- `backend/internal/httpapi/catalog.go` — catalog and company handlers.
- Tests alongside each file.

---

## Task 1: Catalog domain models and safe URL rule

**Files:** `backend/internal/catalog/models.go`, `backend/internal/catalog/url.go`, `backend/internal/catalog/url_test.go`, `backend/internal/catalog/models_test.go`

**Interfaces produced:** `Collection`, `Product`, `ProductImage`, `ProductSpecification`, `ProductResource`, `CompanyContent` structs with JSON tags matching the storefront contract; `SafeResourceURL(string) bool`.

- [ ] Write failing tests for `SafeResourceURL` mirroring `isSafeProductResourceUrl` cases (site path ok, `//host` rejected, backslash rejected, absolute path with `..` escaping rejected, http/https ok, other schemes rejected).
- [ ] Run tests, confirm failure.
- [ ] Implement `SafeResourceURL`.
- [ ] Write a failing test that `Product` JSON marshals published status and omits unused optional fields.
- [ ] Implement models with `omitempty` where appropriate.
- [ ] Run tests; commit.

## Task 2: Catalog schema migration

**Files:** `backend/cmd/storefront-migrate/migrations/00002_create_catalog.sql`

Schema:
- `collections(id BIGSERIAL PK, slug UNIQUE, name, description, status CHECK draft|published, sort_order, created_at, updated_at)`
- `products(id TEXT PK, slug UNIQUE, name, model, collection_id FK, short_description, description, price_amount NUMERIC, price_currency TEXT, compare_at_price NUMERIC, status CHECK, content_status CHECK sample|verified, featured BOOL, sort_order INT, badge, fuel, burners INT, cooking_area, finish, minimum_order_quantity, lead_time, applications JSONB, images JSONB, specifications JSONB, resources JSONB, created_at, updated_at)`
- `company_content(id SMALLINT PK CHECK (id=1), about JSONB, manufacturing JSONB, credentials JSONB, case_studies JSONB, faqs JSONB, content_status CHECK, updated_at)`
- Indexes on `products(status, sort_order)` and `products(collection_id, status)`.

- [ ] Author the `-- +goose Up`/`-- +goose Down` migration.
- [ ] Run `storefront-migrate up` and `status`; confirm applied and repeatable.
- [ ] Commit.

## Task 3: PostgreSQL read repository

**Files:** `backend/internal/catalog/repository.go`, `backend/internal/catalog/postgres.go`, `backend/internal/catalog/postgres_test.go`

**Interfaces produced:**
- `type Repository interface { ListPublishedCollections(ctx) ([]Collection, error); GetPublishedCollection(ctx, slug) (Collection, error); ListPublishedProducts(ctx, Filter) ([]Product, error); GetPublishedProduct(ctx, slug) (Product, error); GetCompanyContent(ctx) (CompanyContent, error) }`
- `catalog.ErrNotFound` sentinel.

- [ ] Write tests using a real connection to `ecommtemple_dev` guarded by `DATABASE_URL` env presence; skip when absent.
- [ ] Implement queries with parameterized SQL, `context`, and bounded limits.
- [ ] Verify drafts are excluded and unknown slugs return `ErrNotFound`.
- [ ] Run tests; commit.

## Task 4: HTTP catalog endpoints

**Files:** `backend/internal/httpapi/catalog.go`, `backend/internal/httpapi/errors.go`, tests.

- [ ] Write handler tests using an in-memory fake repository.
- [ ] Implement envelope, error mapping, query validation (`limit` 1..100 default 24, `featured` bool), and handlers.
- [ ] Wire routes into `NewHandler`.
- [ ] Run tests; commit.

## Task 5: Idempotent seed command

**Files:** `backend/cmd/storefront-seed/main.go`, `backend/internal/catalog/seed.go`, tests.

- [ ] Write tests validating seed data integrity (unique slugs, collection references resolve, resource URLs safe, content marked `sample`).
- [ ] Implement seed with upserts keyed on stable IDs/slugs.
- [ ] Guard the command to the local `ecommtemple_dev` database like the migration command.
- [ ] Run `storefront-seed`, then verify API returns seeded data.
- [ ] Update `backend/README.md`, root `README.md`, and `docs/superpowers/plans/TODO.md`.
- [ ] Run backend + frontend checks and `git diff --check`; commit.

## Verification

- `cd backend && gofmt -l ./cmd ./internal` (no output), `go test ./...`, `go vet ./...`.
- Migration status applied; seed repeatable without duplicates.
- HTTP smoke: list/detail/company endpoints return seeded data; unknown slug returns 404; draft excluded.
- No credentials in tracked files; `git diff --check` clean.

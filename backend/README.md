# Storefront API (foundation)

This is the Go API foundation for the `gas-grill` storefront. Each deployed store owns its own API and PostgreSQL database. This foundation currently provides configuration, health endpoints, and explicit SQL migrations; catalog and inquiry APIs are not implemented yet.

## Requirements

- Go 1.25 or newer (the checked-in `toolchain` directive selects Go 1.27.1 when automatic toolchain selection is enabled)
- PostgreSQL 18 for the current local development target (production server version should be selected and documented per deployment)
- A dedicated local database named `ecommtemple_dev` (created and migrated for this workspace)

The project reuses the existing host PostgreSQL service. It does not create a second local PostgreSQL container.

## Local database setup

The database must be created once by a local PostgreSQL role with `CREATEDB` permission. Never point the migration command at a different database. The command enforces the exact `ecommtemple_dev` database name.

Set the connection string in the shell or an ignored `backend/.env` file. Do not commit credentials:

```sh
export DATABASE_URL='postgres://YOUR_PG_USER:YOUR_PASSWORD@localhost:5432/ecommtemple_dev?sslmode=disable'
export HTTP_ADDR='127.0.0.1:8080'
```

For passwordless local socket/trust auth, use the role and connection settings configured on your machine, while keeping the database name `ecommtemple_dev`. Do not put production credentials in this file or in source control.

## Migrations

Migrations are explicit and versioned. They are not run automatically on API startup.

```sh
cd backend
go run ./cmd/storefront-migrate status
go run ./cmd/storefront-migrate up
go run ./cmd/storefront-migrate status
```

The local migration command refuses URLs whose database is not `ecommtemple_dev` on a loopback host. The initial migration records the selected template identifier; the second creates the catalog and company-content schema. Do not use destructive reset/drop commands against local or production databases.

## Seed sample catalog data

The seed command is idempotent and only targets the local `ecommtemple_dev` database. It writes the reviewed sample catalog (collections, products) and clearly-marked sample company content.

```sh
cd backend
go run ./cmd/storefront-seed
```

Sample records are marked `content_status = sample`. They are not verified real-world product or company facts. Do not present sample MOQ, lead time, certification, or factory claims as fact.

## Run API

```sh
cd backend
go run ./cmd/storefront-api
```

Health:

- `GET http://127.0.0.1:8080/health/live` reports process liveness.
- `GET http://127.0.0.1:8080/health/ready` reports database readiness.

Public read-only catalog (published records only; drafts and unknown slugs return `404`):

- `GET /api/v1/collections`
- `GET /api/v1/collections/{slug}`
- `GET /api/v1/products?collection={slug}&featured={true|false}&limit={1..100}`
- `GET /api/v1/products/{slug}`
- `GET /api/v1/company-content`

Responses use a `{"data": ...}` envelope. Errors use `{"error": {"code", "message", "fields?"}}` with codes `not_found`, `invalid_query`, and `internal_error`. Health and error responses never return the database URL or internal error detail.

## Validate

```sh
cd backend
gofmt -w ./cmd ./internal
go test ./...
go vet ./...
```

Catalog integration tests run against `ecommtemple_dev` inside a rolled-back transaction and are skipped when `DATABASE_URL` is unset.

Inquiry intake, admin authentication and management endpoints are future phases and are not available yet.

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

The local migration command refuses URLs whose database path is not `/ecommtemple_dev`. The initial migration records the selected template identifier. Do not use destructive reset/drop commands against local or production databases.

## Run API

```sh
cd backend
go run ./cmd/storefront-api
```

- `GET http://127.0.0.1:8080/health/live` reports process liveness.
- `GET http://127.0.0.1:8080/health/ready` reports database readiness.

Health endpoints return only a generic status; they do not return the database URL or low-level connection error.

## Validate

```sh
cd backend
gofmt -w ./cmd ./internal
go test ./...
go vet ./...
```

Catalog, inquiry intake, admin authentication and management endpoints are future phases and must not be considered available based on this foundation.

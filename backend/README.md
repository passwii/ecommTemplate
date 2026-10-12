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

### Configuration reference

- `DATABASE_URL` (required): PostgreSQL connection URL.
- `HTTP_ADDR` (default `:8080`): listen address.
- `TRUST_PROXY_HEADER` (default `false`): when `true`, inquiry rate limiting keys on the first `X-Forwarded-For` value instead of the transport peer address. Enable **only** when a trusted reverse proxy always overwrites that header — the storefront's `/api/inquiries` route does this. Left off behind a proxy, every visitor shares one rate-limit bucket.
- `INQUIRY_RATE_LIMIT` (default `5`) and `INQUIRY_RATE_WINDOW` (default `10m`): per-client inquiry submission limit.
- `SESSION_COOKIE_SECURE` (default `false`): marks the admin session cookie `Secure`. Enable when the storefront is served over HTTPS.

## Migrations

Migrations are explicit and versioned. They are not run automatically on API startup.

```sh
cd backend
go run ./cmd/storefront-migrate status
go run ./cmd/storefront-migrate up
go run ./cmd/storefront-migrate status
```

The local migration command refuses URLs whose database is not `ecommtemple_dev` on a loopback host. Migrations create the storefront metadata, catalog/company-content schema, and the inquiry schema. Do not use destructive reset/drop commands against local or production databases.

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

## Public inquiry intake

`POST /api/v1/inquiries` accepts visitor submissions and stores them in the store database.

Request body (`application/json`, max 64 KB, unknown fields rejected):

```json
{
  "contact": {"name": "Alex Morgan", "email": "buyer@example.com", "company": "Example Importers", "phone": "+1 555 0100"},
  "destinationCountry": "US",
  "productId": "grill-01",
  "quantity": 120,
  "customizationRequest": "Private-label packaging options",
  "privacyConsent": {"accepted": true, "noticeVersion": "inquiry-privacy-v1"},
  "sourcePath": "/products/ridge-4-burner",
  "requestId": "optional-idempotency-key"
}
```

Behavior:

- `contact.name`, `contact.email`, and accepted privacy consent are required. Company, phone, country, product, quantity and message are optional.
- `destinationCountry` must be a two-letter code when present and is normalized to uppercase.
- `sourcePath` must be a path within this storefront.
- `productId` must reference a published product; the server stores its own product snapshot (slug/model/name) so later product edits do not rewrite history. An unknown or unpublished product returns `404`.
- Operational fields (status, owner, activities) are never accepted from visitors; unknown fields are rejected with `400`.
- `referenceCode` is a honeypot field (deliberately non-semantic so browser autofill does not trip it); filling it rejects the submission.- `requestId` makes retries idempotent: a replay returns the existing inquiry with `200` instead of creating a duplicate.
- Submissions are rate limited per client address (`INQUIRY_RATE_LIMIT` per `INQUIRY_RATE_WINDOW`, default 5 per 10m) and return `429` when exceeded. The limiter is in-memory per instance.
- Errors use codes `validation_error`, `not_found`, `rate_limited`, `unsupported_media_type`, and `internal_error`. Failures never return a fake success and never expose SQL or stack detail.

Success responses:

- `201` on create, `200` on idempotent replay:

```json
{"data": {"id": "inq_...", "status": "received", "receivedAt": "2026-10-12T06:01:20Z", "privacyNotice": "inquiry-privacy-v1"}}
```

New inquiries start in the `new` status with a `created` activity recorded by `system`. Inquiry management (authentication, list/detail, status changes, notes) is a later phase and is not implemented yet.

## Admin sign-in and inquiry management

The Go API owns admin authentication. It stores bcrypt password hashes and
opaque session tokens (only the token hash is persisted). Create the first admin
once, against the local development database:

```sh
cd backend
ADMIN_EMAIL='owner@example.com' ADMIN_PASSWORD='at least 12 characters' go run ./cmd/storefront-admin create-admin
```

For a deployed store the same command must be pointed at that store's database
and requires the explicit opt-in, because it is the only command allowed to leave
the local database:

```sh
ALLOW_REMOTE_DATABASE=true DATABASE_URL='postgres://...' ADMIN_EMAIL='owner@example.com' go run ./cmd/storefront-admin create-admin
```

Do not pass the password as a command-line argument; use `ADMIN_PASSWORD` or pipe
it on stdin. Rotate the password by creating a new admin and removing the old one.

Admin endpoints (all require a valid session except login):

- `POST /api/v1/admin/sessions` — sign in; sets an HttpOnly, SameSite=Lax session cookie.
- `DELETE /api/v1/admin/sessions` — sign out and clear the cookie.
- `GET /api/v1/admin/session` — current admin identity.
- `GET /api/v1/admin/inquiries?status=&q=&limit=&offset=` — list with total count.
- `GET /api/v1/admin/inquiries/{id}` — detail plus activity history.
- `PATCH /api/v1/admin/inquiries/{id}` — status, owner, next action/date, loss reason.
- `POST /api/v1/admin/inquiries/{id}/activities` — append a `note` or `contact` activity.

Operational fields can only be changed here; the public inquiry endpoint never
accepts them. Status changes are recorded as history entries with the acting
admin. Unauthenticated requests receive `401 unauthorized`.

Set `SESSION_COOKIE_SECURE=true` when the storefront is served over HTTPS.

## Validate

```sh
cd backend
gofmt -w ./cmd ./internal
go test ./...
go vet ./...
```

Catalog and inquiry integration tests create an isolated schema in `ecommtemple_dev`, so they neither depend on nor disturb existing data, and are skipped when `DATABASE_URL` is unset.

Admin authentication and inquiry management endpoints are future phases and are not available yet.

# Go API and PostgreSQL Foundation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a small, independently deployable Go API foundation for the `gas-grill` storefront that uses a dedicated per-store PostgreSQL database, explicit migrations, configuration, and health checks.

**Architecture:** Keep the Go service in a dedicated `backend/` module, separate from Next.js. Use the standard library `net/http` server and `database/sql` with a PostgreSQL driver, explicit versioned SQL migrations, and environment-based configuration. Local development may connect to the already-running host PostgreSQL at `localhost:5432`, but must use a new project-specific database name and never mutate any pre-existing database.

**Tech Stack:** Go 1.27.1, Go standard library `net/http`, `database/sql`, PostgreSQL 18.6 local server, `github.com/jackc/pgx/v5` v5.7.6 as the `database/sql` driver, Goose v3.26.0 with embedded SQL migrations.

## Global Constraints

- `template_id: gas-grill`; extend this template and do not introduce a second visitor-facing style.
- Each store has an independent deployment and independent PostgreSQL database; no shared multi-tenant schema.
- Do not inspect, alter, drop, or reuse existing local PostgreSQL databases; create/use a separately named `ecommtemple_dev` database only after confirming it does not already exist and obtaining/using explicit local credentials.
- Do not commit database passwords, connection URLs with credentials, session secrets, or production secrets; use ignored local env files or environment variables.
- Use explicit, reviewable SQL migrations; do not silently auto-alter production schema on server startup.
- Keep Next.js presentation/template composition independent of Go persistence and SQL details.
- No checkout, payment processing, or third-party commerce plugins.
- Follow `AGENTS.md`, `TESTING.md`, and the approved design/API baseline.

---

## Repository and Local Runtime Decisions

- Backend module root: `backend/` (separate `go.mod` from the Next.js root package).
- Implemented layout: `backend/cmd/storefront-api/`, `backend/cmd/storefront-migrate/` with embedded `migrations/`, `backend/internal/config/`, `backend/internal/httpapi/`, `backend/internal/database/`, `backend/.env.example`, `backend/.gitignore`, `backend/README.md`.
- Local PostgreSQL server is 18.6 and listens on `localhost:5432`. Dedicated `ecommtemple_dev` was confirmed absent and created; initial migration was applied there only.
- Local current role was `damon` with CREATEDB capability; no password was required for the successful local connection. Do not expose credentials or use this assumption on another machine.
- Backend dependencies are pinned by `go.mod`/`go.sum`: pgx v5.7.6 and Goose v3.26.0.
- Deployment remains one API + one PostgreSQL database per store. Local development reuses the host PostgreSQL; no Docker PostgreSQL service was added.

## Task 1: Establish isolated backend module and configuration

**Files:**
- Create: `backend/go.mod`
- Create: `backend/cmd/storefront-api/main.go`
- Create: `backend/internal/config/config.go`
- Create: `backend/internal/config/config_test.go`
- Create: `backend/.env.example`
- Create: `backend/.gitignore` (ignore `.env`, local data/export artifacts; do not ignore migration files)
- Create: `backend/README.md`
- Modify: root `.gitignore` only if the nested local secret file needs a root-level rule and the nested ignore is insufficient.

**Interfaces:**
- `config.Load() (Config, error)` reads required listen/database settings without logging secrets.
- `Config` contains `HTTPAddr`, `DatabaseURL`, and explicit request/readiness timeouts; local defaults must be documented.
- Main starts an HTTP server and exits with a non-zero status on invalid configuration or startup failure.

- [ ] **Step 1: Write configuration tests** for missing database URL, valid URL, malformed URL, default HTTP address, and ensure errors never include the secret value.
- [ ] **Step 2: Run configuration tests and confirm the expected failure** because `Load` does not yet exist.
- [ ] **Step 3: Implement minimal config parsing and validation**; do not add a dotenv dependency unless repo execution requires it. Environment variables are enough for the first iteration.
- [ ] **Step 4: Add a minimal entrypoint with server timeouts and clean signal-driven shutdown.**
- [ ] **Step 5: Run `cd backend && gofmt -w ./cmd ./internal && go test ./... && go vet ./...`.** Expected: all tests pass and vet is silent.
- [ ] **Step 6: Document `DATABASE_URL` setup with placeholder-only examples and verify no real secret is committed.**

## Task 2: Add safe PostgreSQL connectivity and health endpoints

**Files:**
- Create: `backend/internal/database/postgres.go`
- Create: `backend/internal/database/postgres_test.go`
- Create: `backend/internal/httpapi/health.go`
- Create: `backend/internal/httpapi/health_test.go`
- Modify: `backend/cmd/storefront-api/main.go`
- Modify: `backend/README.md`

**Interfaces:**
- `database.Open(ctx, databaseURL) (*sql.DB, error)` configures bounded pool settings, pings with context, and never returns/logs credential-bearing connection strings.
- `httpapi.NewHandler(db *sql.DB) http.Handler` registers `GET /health/live` and `GET /health/ready`.
- Liveness checks process responsiveness only; readiness checks PostgreSQL availability with a deadline.

- [ ] **Step 1: Add handler tests** for liveness 200, readiness 200 with a healthy database, and readiness 503 when the database ping fails; assert response bodies contain no DSN/secrets.
- [ ] **Step 2: Run handler tests and observe failure before implementing handlers.**
- [ ] **Step 3: Implement health handlers** with JSON content type, stable status payloads, and bounded DB ping timeout.
- [ ] **Step 4: Add PostgreSQL open/pool configuration** using `database/sql` and parameterized queries for all future data access.
- [ ] **Step 5: Run unit tests without requiring a live DB; run the integration health test against the dedicated project database only.**
- [ ] **Step 6: Run `cd backend && go test ./... && go vet ./...`.**

## Task 3: Add explicit migrations and isolated project database bootstrap

**Files:**
- Create: `backend/migrations/000001_create_schema_migrations.up.sql`
- Create: `backend/migrations/000001_create_schema_migrations.down.sql`
- Create: `backend/internal/migrate/migrate.go`
- Create: `backend/internal/migrate/migrate_test.go`
- Create: `backend/cmd/storefront-migrate/main.go` (if migration tool is an in-repo command; otherwise use selected external tool with pinned version)
- Modify: `backend/README.md`
- Create local-only env file if needed: `backend/.env` (must remain ignored; never display contents)

**Interfaces:**
- Migration command supports explicit `up` and `status` (or documented equivalent); down migration is for local/test recovery only and requires explicit target database confirmation.
- Migration version table lives inside the project database and records applied versions/checksums.
- Migrations run in a transaction where PostgreSQL and selected migration strategy permit.

- [ ] **Step 1: Choose a migration runner** with SQL-first, versioned migration files and minimal dependencies; record choice and pinned version in README/go tooling.
- [x] **Step 2: Confirm `ecommtemple_dev` did not exist using a narrowly scoped catalog query** against the configured local PostgreSQL connection; no unrelated database data was enumerated.
- [x] **Step 3: Create only the dedicated project database using explicit local credentials.** The current role had CREATEDB capability; no global roles or permissions were changed.
- [ ] **Step 4: Add further migration runner tests** against a disposable test database or transaction-safe test schema, including failure rollback. Current implementation includes database-name guard tests and was run against the dedicated project DB.
- [x] **Step 5: Implement embedded Goose migration `up` and `status` commands**; refuse database names other than the project database.
- [x] **Step 6: Apply the initial migration only to `ecommtemple_dev`; verify status; no DROP/reset was run.**
- [x] **Step 7: Run `cd backend && go test ./... && go vet ./...` and migration status against the dedicated local DB.**

## Task 4: Document and verify clean local workflow

**Files:**
- Modify: `backend/README.md`
- Modify: root `README.md`
- Modify: root `TESTING.md` only if backend test commands need to become project-wide policy.
- Modify: `docs/superpowers/plans/2026-10-12-storefront-backend-roadmap.md`
- Modify: `docs/superpowers/plans/2026-10-12-storefront-issues.md`

**Interfaces:**
- Backend README describes prerequisites, env vars, database name, migration commands, run commands, tests, and safety boundaries.
- Root README continues to distinguish current frontend capability from newly delivered API foundation; do not claim catalog API or inquiry persistence until later phases implement them.

- [x] **Step 1: Document exact commands** for local DATABASE_URL configuration, migration status/up, API run, and tests.
- [ ] **Step 2: Verify fresh clone workflow** while reusing only the dedicated `ecommtemple_dev` database.
- [ ] **Step 3: Run root `npm run lint`, `npm run typecheck`, `npm test`, and `npm run build`** because the root README and overall repo behavior are touched.
- [x] **Step 4: Run backend `gofmt`, `go test ./...`, and `go vet ./...`.**
- [ ] **Step 5: Run `git diff --check` and inspect `git status`; confirm no `.env`, DB dump, credential, or database files are tracked.**
- [x] **Step 6: Update the local TODO board for completed, evidenced work.**

## Verification and Stop Conditions

- Never run `dropdb`, `DROP DATABASE`, broad PostgreSQL cleanup, global role changes, or database reset against an unknown/existing database.
- Never run database inventory/dump commands against unrelated databases. Only query whether `ecommtemple_dev` exists and its server version/connection readiness.
- If server credentials or permissions are unavailable, stop before creating/changing the database and report the exact command the user needs to authorize/run.
- Verify API health does not expose database credentials or raw internal errors.
- Expected final evidence: local project DB name, PostgreSQL server version (not credentials), migration status, `go test ./...`, `go vet ./...`, root frontend checks if docs/config impact them, and `git diff --check`.

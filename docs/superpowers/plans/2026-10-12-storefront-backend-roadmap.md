# Storefront Backend and Inquiry Roadmap

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Progress the existing `gas-grill` storefront from static mock data to a per-store Go/PostgreSQL-backed catalog and reliable inquiry intake with lightweight follow-up management.

**Architecture:** Each storefront deployment owns one Go API and one PostgreSQL database. Go owns persistence and API contracts; Next.js owns storefront rendering and template-specific composition. Work is split into independently reviewable plans: (1) backend foundation and catalog API, (2) inquiry intake and persistence, and (3) minimal internal inquiry management. This document is the cross-phase roadmap; each phase should receive a scoped issue and detailed task plan before implementation.

**Tech Stack:** Existing Next.js 15 / React 19 / TypeScript storefront; Go HTTP API; PostgreSQL; migration and query tooling to be selected in Phase 1 design. Keep dependencies minimal and align with repository policy.

## Global Constraints

- `template_id: gas-grill`; extend the selected template, do not add a second storefront visual style.
- Every store has an independent deployment and independent PostgreSQL database; do not build shared multi-tenant storage.
- Next.js owns rendering and presentation; template page composition stays under `src/components/templates/<template-id>/` and `src/app/`.
- Reusable catalog/content contracts and access belong under `src/lib/catalog/`; template identifiers/metadata stay in `src/lib/templates/registry.ts`.
- Mock/seed content must distinguish verifiable facts from sample placeholders; never invent certifications, test evidence, delivery promises, MOQ, or factory claims.
- Current browser-local admin data is not production persistence or publication; do not claim otherwise.
- No checkout, payment processing, or third-party commerce plugins.
- Follow `AGENTS.md` and `TESTING.md`; report only checks actually run.

---

## Scope and Phase Order

### Phase 0 — Confirm API and data contract

**Deliverable:** Reviewed contract for catalog/company content and inquiry intake, plus explicit repository/service boundary decisions. This is a design/issue phase, not a backend implementation task.

**Decisions required before implementation:** Go version; HTTP router (or standard library); SQL access strategy; migration tool; PostgreSQL version; local service orchestration; Next.js runtime data fetching and unavailable-API behavior; admin authentication boundary; attachment storage decision; public form required fields and retention/privacy policy.

**Acceptance:** Product/catalog/company/inquiry/activity entities, JSON shapes, identifiers, validation, lifecycle statuses, error conventions, and local startup expectations are documented. Product records retain current storefront-required data while identifying sample content honestly. Store isolation is demonstrated by deployment/database configuration, not tenant IDs in shared tables.

### Phase 1 — Go/PostgreSQL foundation and catalog API

**Deliverable:** Independently runnable Go service with configuration, migrations, database connectivity, health endpoint and test baseline; database-seeded sample catalog/company content exposed through documented API endpoints; Next.js can consume the API for public catalog/company pages.

**Likely issue boundaries:**
1. Go service + PostgreSQL development foundation.
2. Catalog/company schema, seed data and read API.
3. Next.js catalog/company API integration and README update.

**Acceptance:** Local setup documented and reproducible; migrations are explicit/versioned; API responses match the approved TypeScript-facing content contract; published/draft behavior, unknown slugs, empty collections and API failures have defined behavior; current public routes use API data rather than static mock arrays as runtime truth. Any retained mock fixture is test-only or seed input, not a competing production source.

### Phase 2 — Public inquiry intake and persistence

**Deliverable:** Public form and API create inquiries safely and durably in the store database.

**Core visitor-submitted data:** contact name, contact email, optional company and phone/WhatsApp, destination country/region, selected product ID plus submission-time model/name snapshot, optional quantity, customization/message, privacy notice/consent record, source path and submission time. Attachments are excluded unless a dedicated storage/security decision is approved.

**Operational data:** server-owned status (`new`, `in_progress`, `quoted`, `won`, `lost`, `on_hold`), optional owner, last-contact time, next action and due date, activity/notes, optional loss reason. Public callers cannot set operational fields.

**Acceptance:** Server validates inputs and product association, stores a durable inquiry and source/context, protects against basic abuse with documented rate limiting and anti-spam measures, returns safe errors, and records consent evidence. Detail CTA pre-fills the product. Form provides pending/success/failure states and does not expose internal operational values. Automated tests cover invalid input, unknown product, client-forged status, duplicate/repeated abuse handling and persistence failure.

### Phase 3 — Lightweight internal inquiry management

**Deliverable:** Authenticated, single-store internal UI for a solo operator or small number of users.

**Initial capabilities:** inquiry list, detail view, status updates, activity/notes, optional owner, next action/date, search/filter by a few core fields. No complex teams, approvals, scoring, sales automation or cross-store views.

**Acceptance:** All management endpoints require authentication and authorization on the server; no frontend-only route hiding or localStorage-only persistence. Status changes and notes are auditable activities. Basic login/session security, secret configuration and backup/restore guidance are documented. Small-team user model is sufficient without a complex role matrix.

### Phase 4 — Evidence-led extensions

Only after the earlier phases are in use, open scoped issues for attachment storage, configurable custom operational fields, CSV export/import or reporting based on observed need. Core query/report fields remain typed columns. Custom fields use explicit definitions, types, validation and values; do not introduce a generic metadata engine until concrete use cases justify it.

---

## Current Repository Mapping

- `src/lib/catalog/mock-data.ts`: current `Product` and `Collection` types, static mock rows, and lookup helpers. Phase 1 must map this data to the approved API contract; it should not remain the production runtime source after integration.
- `src/lib/catalog/company-content.ts`: company content contract plus localStorage helper. Its types can inform the shared contract; its browser storage is not a backend substitute.
- `src/lib/catalog/product-admin.ts`: browser-local product editing and validation. Do not treat this as server-backed publishing. Any future migration should be separately scoped because write/admin UX has broader implications than read-only storefront integration.
- `src/lib/templates/registry.ts`: existing stable template registration boundary; preserve it.
- `src/app/`: route-level rendering for home, product, collection and company pages; update data access without spreading template conditionals.
- `src/components/templates/gas-grill/`: template-owned composition; retain layout ownership here.
- `README.md`: currently states no API/database and static data source. Update only when the corresponding capability exists and local instructions have been verified.
- `AGENTS.md`, `TESTING.md`: delivery and validation policies; update only if project conventions change.

Do not plan migrations of `.next`, `node_modules`, `.git`, or `.zcode` local/session artifacts into the product architecture.

---

## Recommended Issue Sequence

1. **Design/contract issue:** contract draft is captured in `docs/superpowers/specs/2026-10-12-storefront-api-contract.md`; local development reuses the existing host PostgreSQL with a dedicated project DB. Agree remaining API/catalog/company/inquiry contracts and technology defaults before implementation. It produces the detailed Phase 1 implementation plan.
2. **Backend foundation issue:** Go service, PostgreSQL local environment, configuration, migrations, health endpoint and tests.
3. **Catalog API issue:** catalog/company schema, reviewed sample seed data, read endpoints and API documentation.
4. **Storefront integration issue:** replace static public catalog/company runtime reads with API-backed data, handle empty/unknown/unavailable states and update README.
5. **Inquiry persistence issue:** database model, public endpoint, validation, anti-abuse, source/consent storage and tests.
6. **Inquiry form issue:** CTA placement, product-prefill, form UX, privacy notice and submission states.
7. **Inquiry management issue:** authentication, list/detail, status/activity/next-action operations and operational documentation.
8. **Evidence-led follow-up issues:** attachments/custom fields/import/export/reporting only after validating need.

Issues should be created only after reviewing current repository issue conventions/templates and confirming the selected scope. Each issue must name `template_id: gas-grill`, target pages/data scope, whether it extends the template, dependencies, acceptance criteria and locally required validation.

## Test and Delivery Checkpoints

- Backend: unit tests for validation/domain rules; integration tests against PostgreSQL migrations and persistence; API tests for auth/public routes, errors and abuse controls.
- Storefront: relevant `npm run lint`, `npm run typecheck`, `npm run build`, catalog tests, and manual desktop/mobile/keyboard checks as appropriate.
- Data quality: verify sample credentials/claims are not presented as facts; verify seed reruns are deterministic/idempotent as designed.
- Operations: verify fresh local bootstrap, migration from prior schema, backup/restore instructions, configuration through environment, and safe failure when database/API is unavailable.
- Before handoff, run `git diff --check`; include actual command outcomes, not assumptions.

## Explicit Non-Goals

- Shared multi-tenant database or public multi-store CRM.
- General-purpose Salesforce-like CRM, configurable arbitrary objects, complex role/field-level security.
- BANT/MEDDIC, lead scoring, automated assignment, outreach sequences or marketing automation in the first release.
- Binary attachment storage directly in PostgreSQL.
- Claiming sample products/company content are verified real-world facts.
- Immediate rewrite of the browser-local product admin into a complete publishing CMS; scope it only after catalog read API is stable.

## Progress Tracker

- [x] Confirm current branch/worktree state and pull `origin/main` with fast-forward-only behavior; repository was already up to date.
- [x] Record the design baseline in `docs/superpowers/specs/2026-10-12-storefront-api-inquiry-design.md`.
- [x] Map current catalog/company/product-admin boundaries from repository source.
- [x] Create this cross-phase roadmap.
- [x] User reviews/approves the design baseline and roadmap.
- [x] Map and draft phase issues with scope, dependencies, and acceptance criteria in `docs/superpowers/plans/2026-10-12-storefront-issues.md`.
- [ ] Publish issues to GitHub (blocked: API returned 403 for issue creation).
- [ ] Confirm Phase 1 technology defaults and detailed API contract.
- [ ] Write a bite-sized implementation plan for the approved implementation issue.
- [ ] Implement and verify phases incrementally.

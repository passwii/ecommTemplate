# Project delivery rules

## Mandatory template selection
Before implementing or delivering a new product, collection/category, or storefront, select and record its template. Every change request or delivery plan must state:

- `template_id` (for the first storefront use `gas-grill`)
- target pages and content scope
- whether the work extends the selected template or requires a new template

Do not silently introduce a second visual style inside an existing template. Template selection is a delivery/configuration decision, not a visitor-facing switch in the first release.

## Mandatory design language
Before changing visitor-facing UI, the storefront must have a design language document at `docs/design/<template-id>.md`, written per [`.agents/design-language-protocol.md`](./.agents/design-language-protocol.md). That document is the authority for the store's colours, type, layout, components, motion, responsive behaviour, and accessibility floor. The current document is [`docs/design/gas-grill.md`](./docs/design/gas-grill.md).

- Research comparable stores before fixing a visual direction, and record observed evidence separately from inference.
- Token values in the document and in the implementation must match; do not keep two sources of truth.
- Each independently deployed store gets its own document. Do not reuse one store's palette as a project-wide default.

## Architecture boundaries

- Next.js owns storefront rendering and presentation. Keep template-specific page composition under `src/components/templates/<template-id>/` and `src/app/`; reusable catalog types/data access belong under `src/lib/catalog/`.
- Register template metadata and stable identifiers in `src/lib/templates/registry.ts`. Do not scatter template-name conditionals through product data or shared components.
- The first release implements only `gas-grill`. Future templates should use the same catalog/content contract and may compose different layouts and styling.
- Go is the planned shared backend codebase. Each customer storefront will have its own deployment and its own PostgreSQL database. Do not design the first version as a shared multi-tenant database.
- Database schema and migrations may be shared as code across deployments. Store-specific content remains in that store's database.
- Current storefront data is mock data. Keep its fields API-friendly so a future Go API can replace the source without binding template components to a database implementation.
- No checkout, payment processing, or third-party commerce plugins are part of this first template.

## Change checklist

1. Name the selected `template_id` and the design language document being followed before changing storefront UI.
2. Keep the selected template's design consistent across its pages.
3. Use only the tokens the design document defines; no one-off colours, radii, or type sizes.
4. Prefer shared content contracts with template-owned composition.
5. Verify desktop and mobile layouts, keyboard focus, empty/unknown routes, and build/type checks.
6. Update this file and the design document if the delivery convention changes.

## Agent contribution and PR review

- Follow the project test and merge policy in [TESTING.md](./TESTING.md).
- For frontend PRs, use the author's reported local test results as the required machine evidence; Vercel deployment checks are not a merge prerequisite for this project.
- The reviewer still checks the changed code, the test declaration, and relevant edge cases before merging. Missing, failed, or contradictory local test results block approval; do not claim a test passed unless it was run.
- Product administration remains a browser-local MVP. Do not describe localStorage persistence as a secure or production-ready backend.

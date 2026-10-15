# Design language protocol

Every storefront project and every independently deployed store must have its own design language document before visitor-facing UI is changed. The document is the authority for that store's visual direction; it is not a generic style guide shared across stores.

This protocol complements [`template-delivery.md`](./template-delivery.md). Template selection answers *which composition contract the store uses*. The design document answers *what that store looks like and why*.

## Where the document lives

- One document per store: `docs/design/<template-id>.md`.
- The first store uses `docs/design/gas-grill.md`.
- A new store with the same template but a different brand gets its own file (`docs/design/<template-id>-<store-slug>.md`). Do not overwrite another store's document.
- When the template is still the only consumer of a document, keep the plain template id. Add the store suffix only once a second store needs the same template.

## Required content

A design language document must state, at minimum:

1. **Scope** — `template_id`, store, target routes, and whether the work extends the current visual world or replaces it.
2. **Audience and goal** — who visits, and what a successful visit looks like. B2C browsing and B2B enquiry lead to different pages, density, and calls to action.
3. **Product context** — the physical qualities of the product being sold, plus honest constraints (mock data, no checkout, no verified claims).
4. **Research record** — comparable stores actually observed, what was seen, and what is inference rather than observation. Unreachable or partially readable sources are recorded as such; never presented as verified findings.
5. **Visual principles** — the few rules that decide contested cases.
6. **Tokens** — colour roles with values, type stack with weights and tracking, spacing scale, radii, elevation.
7. **Layout rules** — grid, container width, section rhythm, breakpoints.
8. **Components** — buttons, cards, forms, navigation, and the states each one must cover (default, hover, focus, active, disabled, loading, empty, error).
9. **Motion rules** — what animates, why, and how it degrades under `prefers-reduced-motion`.
10. **Responsive rules** — the explicit `< 768px` collapse for every multi-column layout.
11. **Accessibility floor** — contrast targets, focus visibility, target sizes, semantic structure.
12. **Content rules** — voice, and what the store must not claim.
13. **Review checklist** — the checks a change is measured against.

## Research method

Research is recorded evidence, not decoration. Before fixing a visual direction:

- Look at comparable stores in the same product category.
- Record what was actually observed: page structure, product-browsing patterns, hierarchy, imagery treatment.
- Separate **observed** from **inferred**. If a source could not be read (login wall, timeout, blocked region), record that limitation instead of describing the site from memory.
- Do not copy another brand's layout, identity, or copy. Extract the pattern, not the expression.

## Relationship to code

- Token values in the document must match the implementation. When they diverge, fix the code or amend the document — do not leave two sources of truth.
- Store-specific composition stays under `src/components/templates/<template-id>/`. Reusable catalog types and data access stay under `src/lib/catalog/`.
- Do not scatter store-name or template-name conditionals through shared components.

## Change checklist

1. Name the `template_id` and the design document being followed.
2. Confirm the change belongs to the document's scope; if it introduces a new visual world, update the document first.
3. Keep tokens consistent — no one-off colours, radii, or type sizes that the document does not define.
4. Cover the component states the document lists, including empty, loading, and error.
5. Verify desktop and mobile, keyboard focus, and reduced-motion behaviour.
6. Run the project checks in [`../TESTING.md`](../TESTING.md) and report actual results.
7. Update the design document when a convention changes.

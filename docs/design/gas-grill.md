# EMBER / FIELD - gas-grill storefront design language

Follows [`../../.agents/design-language-protocol.md`](../../.agents/design-language-protocol.md) and [`../../.agents/template-delivery.md`](../../.agents/template-delivery.md).

## 1. Scope

- **Template:** `template_id: gas-grill` (the only registered template; see `src/lib/templates/registry.ts`).
- **Store:** EMBER / FIELD, a single-brand outdoor cooking storefront.
- **Routes in scope:** `/`, `/collections/[slug]`, `/products/[slug]`, `/company`.
- **Nature of this work:** replaces the incumbent visual world. The previous forest-green and pale-grey palette, Manrope + DM Sans pairing, and square-cornered component styling are retired. Routes, copy, catalog data, and component boundaries are preserved.
- **Out of scope:** `/admin/*` (a separate browser-local MVP surface that keeps its own `admin.css`), checkout, payment, and any B2B wholesale presentation.

## 2. Audience and goal

- **Primary visitor:** overseas consumers researching a gas grill for a home, patio, or balcony.
- **Successful visit:** the visitor understands what the brand sells, compares grills by burners, cooking area, and fuel with confidence, opens a product, and starts an enquiry.
- **Consequence for the design:** product photography and comparable specifications carry the page. This is a considered browse, not a high-urgency checkout funnel. There is no cart, so "ask about this grill" is the terminal action and must stay prominent and unambiguous.

## 3. Product context

Gas grills are large, heavy, stainless-steel and cast-iron objects. Buyers compare burners, cooking area, and fuel. The material reality is metal, heat, and precision, not rustic craft. The visual language follows that reality.

Honest constraints that shape the design:

- Catalog content is mock data (`src/lib/catalog/mock-data.ts`). No verified certifications, awards, reviews, or performance claims exist, and the design must not imply them.
- No checkout or payment is connected.
- Photography is a fixed set of six images in `public/images/`. The design must work with that set and must not introduce fake product renders or invented imagery.
- Stock levels, delivery windows, and discount badges are not modeled in the data and must not appear as if they were.

## 4. Research record

Comparable stores inspected before fixing the direction. Observed and inferred are kept separate.

**Observed**

- Weber (weber.com US, gas grills collection, loaded via redirect to `amer.weber.com/en-us/collections/gas-grills`): a dense collection page that leads with product imagery, filters, and a strong price/compare affordance. Category navigation is persistent and product-first.
- Napoleon (napoleon.com, gas grills): the page committed and rendered its title and document shell, but the automation snapshot returned only the document root, so no reliable structural detail was captured.
- Industry reference surfaced during search: premium outdoor cooking positioning leans on commercial-grade materials, modular integration, and performance engineering rather than lifestyle-only storytelling.

**Inferred, not verified**

- A consumer comparison flow generally rewards visible, aligned specification values over prose, which is why this document promotes burners, cooking area, and fuel into a comparable layer.
- Categories in this store (freestanding, built-in, compact) map naturally onto a browse-by-space decision, which justifies keeping category entry points on the home page.

**Limitations**

- Broil King (broilkingbbq.com) failed to connect (`ERR_CONNECTION_RESET`) and no observation is claimed for it.
- Several other category sites were unreachable or returned no readable structure. They are not cited.
- No competitor layout, identity, copy, or imagery was copied. Only the browsing pattern (product-first collection, comparable specs, persistent category entry) informed the direction.

## 5. Visual principles

1. **The product is the brightest object on the page.** Neutral, slightly cool surfaces; photography supplies the warmth and the colour.
2. **One accent, used sparingly.** A single deep pine green marks actions and focus. It is never used as decoration or as a background wash.
3. **Comparable numbers align.** Burners, cooking area, and price render in tabular figures so they line up column-to-column in a grid.
4. **Metal, not meadow.** Graphite and steel tones replace the previous leaf-green palette. The brand stays green through its accent and its name, not through its surfaces.
5. **Motion confirms, never performs.** Animation acknowledges a hover, a reveal, or a state change. It does not run on loop.
6. **No unverifiable claims.** No rating stars, review counts, certification seals, stock counters, or delivery promises that the data does not support.

## 6. Tokens

Defined once in `src/app/globals.css` as custom properties. Values here and in the stylesheet must match.

### Colour

| Token | Value | Role |
|---|---|---|
| `--ink` | `#16181a` | Headings, primary text |
| `--graphite` | `#24282c` | Dark sections, footer, primary button fill |
| `--slate` | `#5c646b` | Body copy, secondary text |
| `--paper` | `#f4f6f5` | Page background |
| `--surface` | `#ffffff` | Cards, spec blocks, header |
| `--line` | `#dcdfdd` | Dividers, borders, input strokes |
| `--accent` | `#1d5c3f` | Primary actions, link hover, focus ring |
| `--accent-strong` | `#164a32` | Accent hover and pressed state |

Rules:

- One accent. No second hue is introduced anywhere in the storefront, including badges and status text.
- No pure `#000` and no pure `#fff` as text colours. `--surface` is the only near-white fill.
- Dark-on-light text must clear WCAG AA (4.5:1 body, 3:1 large). Light-on-dark uses `#f4f6f5` on `--graphite` or darker.
- The accent is never the sole carrier of meaning; pair it with text or an icon.

### Typography

Loaded with `next/font` and self-hosted. The previous Google Fonts `@import` is removed.

| Role | Family | Weights | Notes |
|---|---|---|---|
| Display / headings | Archivo | 600, 700 | Tracking `-0.04em`, line-height `1.05` |
| Body / UI | Instrument Sans | 400, 500, 600 | Line-height `1.6` for body copy |
| Numerals | Archivo | 600 | `font-variant-numeric: tabular-nums` |

Scale:

| Step | Size | Use |
|---|---|---|
| `--text-display` | `clamp(2.75rem, 6vw, 4.5rem)` | Page `h1` |
| `--text-h2` | `clamp(1.75rem, 3.4vw, 2.75rem)` | Section headings |
| `--text-h3` | `1.0625rem` | Card titles |
| `--text-body` | `0.9375rem` | Paragraphs |
| `--text-meta` | `0.75rem` | Labels, captions |

- Body copy is capped at `65ch`.
- Italic is not used in display type; emphasis uses weight or colour within the same family.
- Type is set in English only. No locale, time, or weather strips.

### Spacing, radii, elevation

| Token | Value |
|---|---|
| `--space-section` | `clamp(4.5rem, 9vw, 8rem)` |
| `--space-block` | `clamp(2rem, 4vw, 3.5rem)` |
| `--radius-card` | `14px` |
| `--radius-control` | `8px` |
| `--container` | `1320px` |
| `--shadow-card` | `0 1px 2px rgb(22 24 26 / 0.04), 0 12px 32px rgb(22 24 26 / 0.06)` |

- One radius system: cards and image frames `14px`, buttons and inputs `8px`. No pill buttons, no square cards.
- Shadows are tinted toward the ink hue, never pure black, and are reserved for cards that sit above the page surface.

## 7. Layout

- 12-column grid inside a `1320px` container, with a gutter of `clamp(1.25rem, 5vw, 3rem)`.
- Section rhythm uses `--space-section`. Related blocks use `--space-block`.
- Home page follows: split hero, product grid, category entry, full-width approach band, closing enquiry band. This gives five visually distinct families and avoids repeated split sections back to back.
- At most one horizontal marquee-style element per page, and only where it carries real content. None is planned for the first release.
- Eyebrow labels are rationed to at most one per three sections. The home page uses them on the hero and the approach band only.
- Long specification lists use grouped cards, not a hairline-under-every-row table.

## 8. Components

### Buttons

| Variant | Fill | Text | Border | Use |
|---|---|---|---|---|
| Primary | `--accent` | `#f4f6f5` | none | The single terminal action per view |
| Secondary | transparent | `--ink` | `1px --line` | Supporting action |
| Text link | none | `--ink` | bottom border on hover | Navigation within content |

- Minimum height `48px`, minimum target `44px` on touch.
- Label fits one line at desktop. Primary labels stay at three words or fewer.
- Active state applies `transform: scale(0.98)`.
- One label per intent across the whole storefront. The incumbent site uses several labels for the same two intents, which this release consolidates:

| Intent | Label kept | Labels retired |
|---|---|---|
| Browse the grill range | **Explore grills** | "Find your grill", "Shop all grills", "Explore the collection", "Shop grills" |
| Start an enquiry | **Get in touch** | "Ask about this grill" |

- "Get in touch" is used for every enquiry entry point, including the product detail action and the company contact block. The product detail button may carry product context in surrounding copy, but its label stays "Get in touch".
- Header, footer, hero, and closing band all use "Explore grills" for the browse intent.

### Product card

- Image frame `14px`, ratio `1 / 1.08` in featured grids and `1 / 1` in collection grids.
- Below the image: collection label (meta), product name (`--text-h3`), one-line summary, price in tabular figures.
- Hover: image scales to `1.03` over `450ms`, and the card lifts `2px`. Nothing animates on a loop.
- Focus: the entire card is a single focusable link target with a visible ring.
- Badge (for example "Bestseller") is text in a `--surface` chip. It is a real field on the product record, never decorative.

### Specification display

- Product detail groups specs into two columns of cards on desktop: **Cooking** (burners, cooking area, fuel) and **Build** (finish, and any supplied spec rows). One column on mobile.
- Each spec card shows the label in `--text-meta` and the value in tabular figures at `--text-h3`.
- Trade fields (applications, minimum order quantity, lead time) render only when present, in a separate "Supply" group. They are secondary and must not dominate a consumer page.

### Forms and inputs

- Label above the input, helper text below, error text below that in an error colour that is not the accent.
- Focus ring is `--accent` at `2px` with `2px` offset, always visible.
- No placeholder used as a label.

### Navigation

- Header: wordmark, category links, one enquiry link. Desktop navigation renders on one line at `1024px` and above, height `72px`.
- Below `768px` the category links collapse behind a menu button with a `44px` target; the enquiry link stays visible.
- Footer: wordmark, one-line brand statement, link group, copyright with the sample-content disclaimer.

### Required states

Every interactive surface covers default, hover, focus, active, and disabled. Every data surface covers loading, empty, and error:

- **Loading:** skeleton blocks that match the final layout shape. The company page keeps its polite live-region text.
- **Empty:** a composed block explaining the absence and offering the next step. The collection empty state and company empty state both qualify.
- **Error:** inline and specific. No toast-only errors.

## 9. Motion

Intensity target: **4 / 10**. Enough to feel considered, never enough to compete with the product.

| Moment | Behaviour | Duration | Easing |
|---|---|---|---|
| Section reveal on scroll | Opacity `0 → 1`, translate `16px → 0` | `520ms` | `cubic-bezier(0.16, 1, 0.3, 1)` |
| Card hover | Image `scale(1.03)`, card `translateY(-2px)` | `450ms` | `cubic-bezier(0.2, 0.7, 0.2, 1)` |
| Button press | `scale(0.98)` | `120ms` | `ease-out` |
| Header on scroll | Border and background settle once | `200ms` | `ease-out` |

Rules:

- Animate `transform` and `opacity` only.
- No `window.addEventListener("scroll")`. Reveals use IntersectionObserver or CSS scroll-driven animation.
- Reveals fire once per element. Repeated re-animation on scroll-back is not wanted.
- Under `prefers-reduced-motion: reduce`, all of the above resolve instantly to their end state, and `scroll-behavior` returns to `auto`.
- No parallax, no scroll hijacking, no looping ambient animation.

## 10. Responsive rules

Every multi-column layout declares its `< 768px` collapse explicitly.

| Section | Desktop | `< 768px` |
|---|---|---|
| Hero | Split, copy left, image right | Stacked, image first, copy below |
| Featured product grid | Asymmetric 3-up | 2-up, third card full width or wrapped |
| Collection grid | 3-up | 2-up |
| Category entry | 3-up offset | 2-up, third full width |
| Approach band | Split image and copy | Stacked |
| Product detail | 2 columns | 1 column, image first |
| Spec cards | 2 columns | 1 column |
| Footer | 4 columns | 2 columns with links wrapping |

- Full-bleed sections use `min-height` with `100dvh` where a viewport-height is intended. `100vh` is not used.
- Touch targets are at least `44px` in both dimensions.

## 11. Accessibility floor

- Text contrast: AAA target for hero copy, AA minimum everywhere else.
- Focus is never removed. The focus ring is visible on every interactive element, including inside cards that contain their own link.
- Heading order is sequential; each page has exactly one `h1`.
- Images carry meaningful `alt` text. Decorative frames use an empty `alt`.
- Icon-only controls carry an accessible name.
- Reduced motion is honoured, per section 9.
- Colour is never the only signal; badges and states pair colour with text.

## 12. Content and voice

- Voice: plain, calm, specific. Short declarative sentences. No superlatives, no "elevate", no "unleash", no fabricated engineering claims.
- Copy states what the product is and how it is used. It does not assert performance the data cannot support.
- Sample-content disclaimers stay visible where the project already places them (product detail note, company storage note, footer).
- The brand name EMBER / FIELD, the price format, and the enquiry email are unchanged.

## 13. Review checklist

- [ ] `template_id: gas-grill` stated, and this document is the one being followed.
- [ ] Only tokens defined in section 6 appear in the changed CSS. No one-off colours, radii, or sizes.
- [ ] One accent hue across all four routes.
- [ ] Primary action per view is unambiguous, one line wide, and readable against its fill.
- [ ] Product comparison values use tabular figures and align across a grid row.
- [ ] Every changed multi-column layout has an explicit `< 768px` rule.
- [ ] Loading, empty, and error states exist for every touched data surface.
- [ ] Keyboard focus is visible and ordered correctly on every interactive element.
- [ ] `prefers-reduced-motion: reduce` collapses every animation to its end state.
- [ ] No unverifiable claim, rating, stock counter, or delivery promise was introduced.
- [ ] Desktop, tablet, and mobile verified in the browser.
- [ ] `npm run lint`, `npm run typecheck`, `npm run build` run, with actual results reported in [`../../TESTING.md`](../../TESTING.md) terms.
- [ ] This document updated when a convention changed.

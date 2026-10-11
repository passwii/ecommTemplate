# Template delivery protocol

Any new product, category/collection, or storefront delivery must explicitly choose a `template_id` first. The initial template is `gas-grill` (outdoor cooking). State the template, target routes, and whether the change extends that template or introduces a separately approved template. Do not implement an implicit or visitor-facing style switch.

See `AGENTS.md` for architecture and data-isolation requirements. Template registry: `src/lib/templates/registry.ts`.

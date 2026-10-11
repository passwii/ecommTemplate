# Test and merge policy

This document defines the local validation expected for changes to the `gas-grill` storefront project. Agent contribution instructions link here from [`AGENTS.md`](./AGENTS.md).

## Environment

- Node.js version supported by the root `package.json` / lockfile.
- npm, with dependencies installed using `npm ci` from the repository root.
- No external service, database, or environment secrets are required for these checks.

## Local validation commands

Run commands from the project root and report their actual outcome in the PR description or review handoff:

```bash
npm ci                 # first setup or after package-lock.json changes
npm run lint           # ESLint checks the project
npm run typecheck      # TypeScript checks without emitting files
npm run build           # Next.js production build
npm run dev             # optional interactive browser / manual QA
```

For a frontend PR that changes application code, the author should report the checks relevant to the change. `npm run lint` and `npm run typecheck` are expected for normal UI changes; run `npm run build` for route, config, dependency, or other build-sensitive changes. If the PR author reports all three passing, the reviewer may rely on that report and perform a proportionate code/edge-case review rather than repeating every command. The reviewer may run checks locally when results are missing, unclear, or inconsistent with the diff.

Before handoff, also run:

```bash
git diff --check
```

## PR merge gate

- This project does **not** require a Vercel deployment check or preview deployment to merge.
- Merge a frontend PR when its local validation results are reported as successful, review finds no blocking correctness/security issue, and the PR is mergeable.
- Do not merge when a required local check failed, results are absent/contradictory, or review identifies an unresolved blocking issue. Ask the author to fix or provide the missing evidence.
- Never state a command passed unless it was actually run and succeeded. Include command names and meaningful output (for example lint errors/warnings and typecheck/build exit status).

## Scope and runtime checks

For storefront/admin changes, review relevant routes and states, including published versus draft products, empty collections, unavailable products, form validation, and desktop/mobile behavior as applicable. The product admin currently persists to browser `localStorage` only; it has no server persistence, authentication, or cross-device sync. Do not use production-readiness claims for this MVP.

# Agent Instructions: QR Flow Headplane Action Verification

## Task

- **Roadmap item:** `ROADMAP.md` — Requested implementation brief: QR approval
  from the registration page.
- **Worktree:**
  `/home/denny/Project/headscale-project/headscale-qr-flow-headplane-action`
- **Branch:** `feature/qr-flow-headplane-action`
- **Superproject baseline:** `49ebce52d3608803573cb2a07ba21c461e4c18da`
- **Headplane gitlink:** `b470d2f87e52f3ab654897fff00732bd021a8495`

## Objective

Verify the exact-version Headplane QR scanner route/action: only an
authenticated principal with machine-write capability can select an owner,
submit a valid pending-registration payload through the configured Headscale
API client, refresh machines, record the existing audit event, and reach the
registered-machine route.

## Read First

1. Root `AGENTS.md`, then `headplane/AGENTS.md`.
2. `headplane/docs/SPECIFICATION.md`, `headplane/docs/ARCHITECTURE.md`, and
   `headplane/docs/development/testing.md`.
3. `headplane/app/routes/machines/scan-qr.tsx`,
   `headplane/app/components/qr-scanner/qr-scanner.tsx`,
   `headplane/app/utils/register-key.ts`, and existing QR tests.

## Scope

### In scope

- Server loader/action authorization for `Capabilities.write_machines`.
- Required owner selection and server-side QR payload validation.
- Configured `api.nodes.register(user, authID)` invocation, live-store refresh,
  `machine.register` audit entry, and successful redirect.
- Scanner permission/error/cleanup/accessibility regression coverage.
- Focused tests and defect-only fixes in the Headplane submodule.

### Out of scope

- A UI redesign, a new QR scanner library, direct requests to `server_url`
  encoded in a QR payload, or browser-hardware test automation.
- Headscale Go source, schema/API/payload-version/dependency changes, root
  documentation, and unrelated Headplane features.
- Modifying another worktree or the original checkout.

## Security Invariants

- The QR `server_url` is validation data only; never redirect or make requests
  to it. Use the configured Headscale API client.
- Require `write_machines` in both loader and action.
- Reject malformed JSON, wrong type/version, missing owner/auth ID/server URL,
  invalid registration key, invalid URL, and expired payload before calling the
  API client.
- The selected Headscale user is mandatory. Headscale remains authoritative for
  cache expiry and auth-ID consumption.

## Expected Files

Tests first, normally limited to the Headplane submodule:

- `headplane/tests/component/scan-qr.test.tsx`
- `headplane/tests/component/qr-scanner.test.tsx`
- A focused existing-pattern test under `headplane/tests/unit/`

Implementation only if a test demonstrates a defect:

- `headplane/app/routes/machines/scan-qr.tsx`
- `headplane/app/components/qr-scanner/qr-scanner.tsx`
- `headplane/app/utils/register-key.ts`

## Validation

Run from `headplane/` with Node 24.2–24.x and PNPM 10.4–10.x when available:

```bash
pnpm exec vitest run --project component \
  tests/component/qr-scanner.test.tsx \
  tests/component/scan-qr.test.tsx
pnpm exec vitest run --project unit tests/unit/utils/register-key.test.ts
pnpm run lint
pnpm run build
pnpm run typecheck
git diff --check
```

If the specified runtime/dependencies are unavailable, report the exact command
and blocker; do not treat it as passing and do not change project tooling.

## Submodule Commit Rules

The submodule is detached at the approved gitlink. If product changes are
needed, create a branch inside `headplane/` named
`feature/qr-flow-headplane-action`, commit there first, then commit only the
updated `headplane` gitlink in the superproject branch. Do not push without
coordinator approval and do not alter `main` or `develop`.

## Completion

Report both commit SHAs when changes are needed, all changed files, exact test
results, and any proven environment limitation. Do not merge, rebase,
force-push, or rewrite history.
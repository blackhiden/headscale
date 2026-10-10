# Agent Instructions: QR Flow Documentation Reconciliation

## Task

- **Roadmap item:** `ROADMAP.md` — Requested implementation brief: QR approval
  from the registration page.
- **Worktree:**
  `/home/denny/Project/headscale-project/headscale-qr-flow-docs-reconciliation`
- **Branch:** `feature/qr-flow-docs-reconciliation`
- **Superproject baseline:** `49ebce52d3608803573cb2a07ba21c461e4c18da`
- **Headplane gitlink:** `b470d2f87e52f3ab654897fff00732bd021a8495`

## Objective

Synchronize tracking documentation with the verified existing QR registration
implementation. This is a documentation-status correction, not a behavior,
architecture, API, schema, or product-scope change.

## Read First

1. `AGENTS.md`, `ROADMAP.md`, and the local `session.md`.
2. `docs/usage/registration.md`, `docs/ref/registration.md`, and
   `docs/api/qr-payload.md`.
3. `ARCHITECTURE.md`, `CHANGELOG.md`, `hscontrol/handlers.go`, and
   `headplane/app/routes/machines/scan-qr.tsx`.

## Scope

### In scope

- Change the stale Phase-17 brief status from Planned to Implemented.
- Update the roadmap current-state summary to state the verified QR flow.
- Keep `session.md` concise with the objective, exact baseline, decisions, and
  validation results.
- Correct another document only when implementation evidence proves it is
  inaccurate.

### Out of scope

- Changing requirements, payload contract, API, schema, security model, or UI.
- Inventing screenshots, browser claims, test results, or new capabilities.
- Updating unrelated roadmap phases or resolving pre-existing documentation
  contradictions outside Phase 17.

## Acceptance Criteria

- The Phase-17 heading and current-state summary no longer contradict the
  completed checklist and source implementation.
- Existing user/API/architecture documentation remains unchanged unless a
  verified discrepancy is found.
- Markdown is formatted and links/references remain valid.
- Only task-scoped documentation is committed; do not merge or rebase.

## Validation

Run from this worktree:

```bash
git diff --check
make fmt-mdformat
git diff --check
```

If `mdformat` is unavailable, report that limitation exactly and avoid changing
project tooling.
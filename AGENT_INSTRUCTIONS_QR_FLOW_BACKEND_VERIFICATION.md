# Agent Instructions: QR Flow Backend Verification

## Task

- **Roadmap item:** `ROADMAP.md` — Requested implementation brief: QR approval
  from the registration page.
- **Worktree:**
  `/home/denny/Project/headscale-project/headscale-qr-flow-backend-verification`
- **Branch:** `feature/qr-flow-backend-verification`
- **Superproject baseline:** `49ebce52d3608803573cb2a07ba21c461e4c18da`
- **Headplane gitlink:** `b470d2f87e52f3ab654897fff00732bd021a8495`

## Objective

Verify the existing Headscale QR registration contract and add only focused
regression coverage or defect fixes proved necessary by tests. The QR flow is
an alternate transport for an existing pending auth ID; it must not alter the
existing CLI registration behavior or bypass Headscale authorization.

## Read First

1. `AGENTS.md`, `ROADMAP.md`, and `session.md`.
2. `hscontrol/handlers.go`, `hscontrol/api/v1/auth.go`,
   `hscontrol/state/state.go`, and `hscontrol/types/common.go`.
3. `hscontrol/qr/`, `hscontrol/handlers_test.go`,
   `hscontrol/state/auth_cache_test.go`, and `hscontrol/apiv1_auth_test.go`.

## Scope

### In scope

- Pending-registration-only QR rendering from `GET /register/{auth_id}`.
- Fixed auth-cache expiry and prevention of reload/reinsertion extension.
- QR payload validation: required fields, supported type/version, expiry, and
  no credentials in the documented payload.
- Headscale registration API behavior for invalid, missing, expired, and
  consumed auth IDs and invalid users.
- Focused Go unit/server-level regression tests and defect-only fixes.

### Out of scope

- Database migrations, schema changes, new APIs, payload-version changes, and
  new dependencies.
- Headplane source changes, documentation changes, integration Docker tests,
  or unrelated registration refactoring.
- Editing another worktree, the original checkout, or the nested Headplane
  checkout.

## Invariants

- Keep `headscale auth register --auth-id <auth_id> --user USERNAME` unchanged.
- Do not emit API keys, passwords, cookies, private keys, or user credentials
  in QR data.
- Generate QR only for a cached pending registration and use its existing
  absolute `ExpiresAt()` value.
- Preserve Headscale as the authority for auth ID validity, ownership,
  expiration, and single-use consumption.
- Keep `Cache-Control: no-store` on the public registration page.

## Expected Files

Tests first, normally limited to:

- `hscontrol/handlers_test.go`
- `hscontrol/apiv1_auth_test.go`
- `hscontrol/state/auth_cache_test.go`
- `hscontrol/qr/payload_test.go`
- `hscontrol/qr/qr_test.go`

Change implementation only if a failing focused test demonstrates a defect:

- `hscontrol/handlers.go`
- `hscontrol/state/state.go`
- `hscontrol/qr/payload.go`

## Acceptance Criteria

- Pending registration renders the unchanged CLI command, a safe QR image, and
  `Cache-Control: no-store`.
- Missing/non-registration/expired cache state cannot produce a valid QR.
- Rendering or reinserting a request cannot extend its QR deadline.
- Malformed, unsupported, expired, invalid-user, and consumed registration
  attempts fail through existing Headscale behavior.
- Focused tests, build, lint, and `git diff --check` pass.

## Validation

Run from this worktree:

```bash
go test ./hscontrol/qr -count=1
go test ./hscontrol/state -run 'Test(AuthCache|SetAuthCacheEntry)' -count=1
go test ./hscontrol -run 'Test(RegisterHandlerRendersQRForPendingRegistration|APIV1AuthRegister|Template)' -count=1
go test ./hscontrol -count=1
make build
make lint
git diff --check
```

Do not run `hi` or Docker integration tests for this task.

## Completion

Commit only product test/fix changes using the project’s Go-style commit
convention. Report the commit SHA, changed files, commands run, results, and
any proven limitations. Do not merge or rebase any branch.
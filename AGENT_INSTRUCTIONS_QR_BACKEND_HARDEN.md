# Agent Instructions: QR Backend Hardening

## Task Overview

**Task:** Backend QR Hardening  
**Roadmap Item:** Phase 17 — Non-OIDC Interactive Registration  
**Branch:** `feature/qr-backend-harden`  
**Worktree:** `/home/denny/Project/headscale-project/headscale-qr-backend-harden`  
**Base Commit:** `f808528f` (main)

## Objective

Harden the backend QR code generation, expiry, and cache logic for non-OIDC interactive registration. Ensure robust edge-case handling and comprehensive test coverage for:
- Expired QR codes
- Reused/invalid auth IDs
- Malformed QR payloads
- Cache expiry enforcement
- Error responses

## Scope

### In Scope
- Review and harden `hscontrol/qr/` package
- Review and harden auth cache logic in `hscontrol/state/`
- Add backend tests for edge cases
- Improve error handling and validation
- Ensure cache expiry cannot be bypassed

### Out of Scope
- Frontend changes (handled by separate task)
- Documentation updates (handled by separate task)
- Integration/E2E tests (handled by separate task)
- Changes to core registration flow outside QR

## Files/Modules Expected to Change

- `hscontrol/qr/*.go`
- `hscontrol/state/auth_cache.go`
- `hscontrol/handlers.go` (QR-related error handling)
- `hscontrol/qr/*_test.go`
- `hscontrol/state/*_test.go` (auth cache tests)
- `hscontrol/handlers_test.go` (registration handler tests)

## Dependencies

None. This task is independent and can proceed immediately.

## Testing Requirements

1. **Unit Tests:**
   - Test QR generation with expired cache entries
   - Test QR generation with invalid/missing auth IDs
   - Test malformed QR payload rejection
   - Test cache reinsertion preserves original expiry

2. **Error Handling Tests:**
   - Test appropriate HTTP error codes for invalid requests
   - Test error messages are clear and actionable

3. **Validation:**
   - Run: `go test ./hscontrol/qr -v -count=1`
   - Run: `go test ./hscontrol/state -run AuthCache -v -count=1`
   - Run: `go test ./hscontrol -run Register -v -count=1`
   - Run: `make lint`
   - Run: `go build ./cmd/headscale`

## Acceptance Criteria

- [ ] All edge cases for QR generation and expiry are handled
- [ ] Cache expiry cannot be bypassed by reload or reinsertion
- [ ] Error responses are appropriate and clear
- [ ] All new tests pass
- [ ] No regression in existing tests
- [ ] Code is linted and builds successfully
- [ ] Changes are committed to `feature/qr-backend-harden` branch

## Constraints

- Preserve existing QR code format and API contract
- Do not modify frontend code
- Do not change database schema
- Follow existing code style and patterns
- Use existing error handling conventions

## Relevant Documentation

- `ROADMAP.md` — Phase 17 description
- `session.md` — Phase 17 completion notes
- `docs/usage/registration.md` — Non-OIDC registration flow
- `docs/ref/registration.md` — Registration API reference
- `AGENTS.md` — Project interaction rules

## Implementation Notes

1. Start by reading the current QR implementation:
   - `hscontrol/qr/qr.go`
   - `hscontrol/state/auth_cache.go`
   - `hscontrol/handlers.go` (register endpoints)

2. Identify edge cases not currently covered by tests

3. Add tests first (TDD approach preferred)

4. Implement hardening and validation

5. Verify all tests pass

6. Commit incrementally with clear messages

## Worktree Setup

```bash
cd /home/denny/Project/headscale-project/headscale-qr-backend-harden
git status  # Verify on feature/qr-backend-harden branch
git log --oneline -n 1  # Should show f808528f
```

## Final Verification

Before considering this task complete:
- [ ] All tests pass
- [ ] Code is linted
- [ ] Build succeeds
- [ ] Changes are committed
- [ ] Branch pushed to remote (if applicable)

## Contact

If you discover scope creep or need architectural decisions outside this task, stop and request clarification.

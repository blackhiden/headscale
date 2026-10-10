# Agent Instructions: QR Integration/E2E Testing

## Task Overview

**Task:** Integration/E2E Testing  
**Roadmap Item:** Phase 17 — Non-OIDC Interactive Registration  
**Branch:** `feature/qr-e2e-tests`  
**Worktree:** `/home/denny/Project/headscale-project/headscale-qr-e2e-tests`  
**Base Commit:** `f808528f` (main)

## Objective

Add comprehensive integration and end-to-end tests for the full non-OIDC interactive registration flow with QR code support. Ensure all scenarios are covered, including edge cases and error conditions.

## Scope

### In Scope
- Add integration tests for full registration flow (CLI + QR)
- Test expiry and single-use enforcement
- Test admin approval via Headplane
- Test error scenarios (expired QR, invalid auth ID, etc.)
- Validate on all supported platforms

### Out of Scope
- Backend implementation changes (handled by separate task)
- Frontend implementation changes (handled by separate task)
- Documentation updates (handled by separate task)

## Files/Modules Expected to Change

- `integration/*_test.go` (new or updated test files)
- Possibly: `hscontrol/servertest/*` (if server-level tests are added)
- Test fixtures and helpers in `integration/`

## Dependencies

**Sequence 3** — This task should wait for backend and frontend hardening to ensure all edge cases are implemented before testing them. Can proceed with reviewing existing tests and planning new ones.

## Testing Requirements

1. **Integration Tests to Add/Extend:**
   - Full non-OIDC registration flow with QR
   - QR expiry enforcement
   - Single-use auth ID validation
   - Admin approval via Headplane scanner
   - Error scenarios (expired, invalid, malformed)
   - Browser/device compatibility (if applicable)

2. **Validation:**
   - Run: `go test ./integration -v -run QR`
   - Run: `go test ./hscontrol/servertest -v` (if server-level tests added)
   - Run: `go run ./cmd/hi doctor` (check integration test environment)
   - Run: `go run ./cmd/hi run "TestQR*"` (if using hi test runner)

## Acceptance Criteria

- [ ] Full registration flow tested end-to-end
- [ ] All edge cases covered
- [ ] Error scenarios tested
- [ ] All tests pass
- [ ] No flakes or race conditions
- [ ] Tests are documented
- [ ] Changes are committed to `feature/qr-e2e-tests` branch

## Constraints

- Follow existing integration test patterns
- Use `IntegrationSkip(t)` for integration tests
- Use `EventuallyWithT` for external calls
- Do not perform state-mutating commands inside `EventuallyWithT`
- Read `cmd/hi/README.md` and `integration/README.md` before running tests
- Clean up test artifacts (logs, containers)

## Relevant Documentation

- `ROADMAP.md` — Phase 17 description
- `session.md` — Phase 17 completion notes
- `cmd/hi/README.md` — Integration test runner
- `integration/README.md` — Test authoring patterns
- `AGENTS.md` — Project interaction rules

## Implementation Notes

1. **Read integration test documentation:**
   - `cmd/hi/README.md`
   - `integration/README.md`

2. **Review existing registration tests:**
   - Look for patterns to follow
   - Identify gaps in coverage

3. **Plan new tests:**
   - Full flow: node registration → QR display → scan → approval → success
   - Expiry: QR expired before scan
   - Invalid: malformed QR payload
   - Single-use: auth ID consumed

4. **Implement tests incrementally:**
   - One scenario per test function
   - Clear test names and documentation
   - Proper setup and teardown

5. **Run and validate:**
   - Run tests multiple times to check for flakes
   - Check test logs for errors
   - Verify cleanup

6. **Commit with clear messages**

## Worktree Setup

```bash
cd /home/denny/Project/headscale-project/headscale-qr-e2e-tests
git status  # Verify on feature/qr-e2e-tests branch
git log --oneline -n 1  # Should show f808528f

# Ensure dependencies
go mod download

# Check integration test environment
go run ./cmd/hi doctor
```

## Test Scenarios to Cover

1. **Happy Path:**
   - Node registers → QR displayed → Admin scans → Node approved

2. **Expiry:**
   - QR displayed → Wait past expiry → Scan fails with clear error

3. **Single-Use:**
   - QR scanned → Auth ID consumed → Second scan fails

4. **Invalid Payload:**
   - Malformed QR data → Scan fails with validation error

5. **Browser/Device:**
   - Test on different browsers (if applicable)

## Final Verification

Before considering this task complete:
- [ ] All tests pass
- [ ] Tests run without flakes
- [ ] Test logs are clean
- [ ] Cleanup is successful
- [ ] Changes are committed
- [ ] Branch pushed to remote (if applicable)

## Contact

If you discover scope creep or need architectural decisions outside this task, stop and request clarification.

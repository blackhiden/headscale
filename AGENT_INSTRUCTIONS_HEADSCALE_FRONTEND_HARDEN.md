# Agent Instructions: Headscale Frontend Template Hardening

## Task Overview

**Task:** Headscale Frontend Template Hardening  
**Roadmap Item:** Phase 17 — Non-OIDC Interactive Registration  
**Branch:** `feature/headscale-frontend-harden`  
**Worktree:** `/home/denny/Project/headscale-project/headscale-headscale-frontend-harden`  
**Base Commit:** `f808528f` (main)

## Objective

Harden the Headscale built-in registration page templates (`/register/{auth_id}`) for non-OIDC interactive registration. Improve error handling, user feedback, QR code display, and edge-case handling in the Go HTML templates.

**Note:** Headplane (the separate React UI) is already complete. This task focuses on Headscale's own built-in registration pages served by the Go backend.

## Scope

### In Scope
- Review and harden registration page templates in `hscontrol/templates/`
- Improve error handling and user feedback in templates
- Ensure QR code display is robust
- Add template tests for edge cases
- Improve accessibility and UX of built-in registration pages

### Out of Scope
- Backend logic changes (handled by separate task)
- Headplane changes (already complete)
- Documentation updates (handled by separate task)
- Integration/E2E tests (handled by separate task)

## Files/Modules Expected to Change

- `hscontrol/templates/register_confirm.go` (OIDC registration page template)
- `hscontrol/templates/auth_web.go` (non-OIDC registration page template)
- `hscontrol/handlers.go` (template rendering logic if needed)
- `hscontrol/templates_consistency_test.go` (template tests)
- Possibly: new test files for template rendering

## Dependencies

None. This task is independent and can proceed immediately.

## Testing Requirements

1. **Template Tests:**
   - Test registration page renders correctly with QR code
   - Test error states (expired, invalid auth ID)
   - Test missing/malformed data handling
   - Test template HTML consistency

2. **Validation:**
   - Run: `go test ./hscontrol -run Template -v -count=1`
   - Run: `go test ./hscontrol -run Register -v -count=1`
   - Run: `make lint`
   - Run: `go build ./cmd/headscale`

3. **Manual Verification:**
   - Start Headscale locally
   - Visit `/register/{auth_id}` page
   - Verify QR code displays correctly
   - Verify CLI command displays correctly
   - Check error states

## Acceptance Criteria

- [ ] Registration page templates handle all edge cases gracefully
- [ ] Error messages are clear and actionable
- [ ] QR code display is robust
- [ ] CLI command display is clear
- [ ] All template tests pass
- [ ] No regression in existing templates
- [ ] Code is linted and builds successfully
- [ ] Changes are committed to `feature/headscale-frontend-harden` branch

## Constraints

- Preserve existing template structure and styling
- Do not modify backend logic outside templates
- Follow existing template patterns in `hscontrol/templates/`
- Maintain backward compatibility
- Keep templates simple and maintainable

## Relevant Documentation

- `ROADMAP.md` — Phase 17 description
- `session.md` — Phase 17 completion notes
- `docs/usage/registration.md` — Non-OIDC registration flow
- `AGENTS.md` — Project interaction rules

## Implementation Notes

1. Start by reading the current template implementation:
   - `hscontrol/templates/auth_web.go` (non-OIDC page)
   - `hscontrol/templates/register_confirm.go` (OIDC page)
   - `hscontrol/handlers.go` (template rendering)

2. Review how templates handle:
   - QR code data
   - Expiry timestamps
   - Error states
   - Missing data

3. Identify edge cases not covered:
   - Expired QR codes
   - Missing auth ID
   - Malformed QR data
   - Long usernames/auth IDs

4. Add tests first (TDD approach preferred)

5. Improve template robustness and error handling

6. Verify all tests pass

7. Test manually by starting Headscale locally

8. Commit incrementally with clear messages

## Worktree Setup

```bash
cd /home/denny/Project/headscale-project/headscale-headscale-frontend-harden
git status  # Verify on feature/headscale-frontend-harden branch
git log --oneline -n 1  # Should show f808528f
```

## Template Edge Cases to Handle

- QR code data missing or malformed
- Expiry timestamp in the past
- Very long auth IDs or server URLs
- Missing CSS or JavaScript resources
- XSS protection for user-provided data
- Proper escaping in templates

## Final Verification

Before considering this task complete:
- [ ] All template tests pass
- [ ] Code is linted
- [ ] Build succeeds
- [ ] Manual testing shows proper rendering
- [ ] Changes are committed
- [ ] Branch pushed to remote (if applicable)

## Contact

If you discover scope creep or need architectural decisions outside this task, stop and request clarification.

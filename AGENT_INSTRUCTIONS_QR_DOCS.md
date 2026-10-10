# Agent Instructions: QR Documentation Update

## Task Overview

**Task:** Documentation Update  
**Roadmap Item:** Phase 17 — Non-OIDC Interactive Registration  
**Branch:** `feature/qr-docs`  
**Worktree:** `/home/denny/Project/headscale-project/headscale-qr-docs`  
**Base Commit:** `61878b4b` (existing worktree)

## Objective

Update user/admin documentation for non-OIDC interactive registration with QR code support. Ensure documentation is complete, accurate, and user-friendly. Add troubleshooting guidance and update diagrams/screenshots as needed.

## Scope

### In Scope
- Review and update `docs/usage/registration.md`
- Review and update `docs/ref/registration.md`
- Add troubleshooting/FAQ section
- Update diagrams/screenshots (if needed)
- Ensure consistency with backend/frontend implementation

### Out of Scope
- Backend changes (handled by separate task)
- Frontend changes (handled by separate task)
- Integration/E2E tests (handled by separate task)

## Files/Modules Expected to Change

- `docs/usage/registration.md`
- `docs/ref/registration.md`
- Possibly: `docs/assets/*` (screenshots, diagrams)
- `ROADMAP.md` (if scope clarifications needed)
- `CHANGELOG.md` (if documentation fixes are significant)

## Dependencies

**Sequence 2** — This task should wait for backend and frontend review to ensure documentation matches implementation. However, it can proceed with review and draft updates.

## Testing Requirements

1. **Documentation Review:**
   - Read updated docs as a new user
   - Verify all steps are clear and complete
   - Check links and references

2. **Validation:**
   - Run: `make fmt` (formats docs with mdformat)
   - Verify markdown syntax
   - Check cross-references

## Acceptance Criteria

- [ ] Documentation is complete and accurate
- [ ] Troubleshooting/FAQ section added
- [ ] All flows (CLI, QR) are documented
- [ ] Diagrams/screenshots updated if needed
- [ ] Docs are formatted correctly
- [ ] Cross-references are valid
- [ ] Changes are committed to `feature/qr-docs` branch

## Constraints

- Follow existing documentation style and structure
- Use existing screenshot/diagram conventions
- Do not invent features not in implementation
- Maintain consistency with other docs

## Relevant Documentation

- `ROADMAP.md` — Phase 17 description
- `session.md` — Phase 17 completion notes
- Existing `docs/usage/registration.md`
- Existing `docs/ref/registration.md`
- `AGENTS.md` — Project interaction rules

## Implementation Notes

1. Start by reading the current documentation:
   - `docs/usage/registration.md`
   - `docs/ref/registration.md`

2. Review backend and frontend implementations to understand what changed

3. Identify gaps, inconsistencies, or unclear sections

4. Draft updates with clear, user-friendly language

5. Add troubleshooting section with common issues:
   - QR code expired
   - Scanner not working
   - Invalid auth ID
   - Camera permission issues

6. Update diagrams/screenshots if needed (or note placeholders)

7. Format and validate

8. Commit with clear messages

## Worktree Setup

```bash
cd /home/denny/Project/headscale-project/headscale-qr-docs
git status  # Verify on feature/qr-docs branch
git log --oneline -n 1  # Should show 61878b4b
```

## Troubleshooting Guidance to Add

Consider adding sections for:
- "QR code shows 'expired' error"
- "Scanner not detecting QR code"
- "Camera permission denied"
- "Registration fails after scan"
- "QR code vs CLI: when to use which"

## Final Verification

Before considering this task complete:
- [ ] Documentation is complete
- [ ] Docs are formatted
- [ ] Cross-references checked
- [ ] Changes are committed
- [ ] Branch pushed to remote (if applicable)

## Contact

If you discover scope creep or need architectural decisions outside this task, stop and request clarification.

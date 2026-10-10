# Agent Instructions: Multi-Instance Dashboard (Task 2 of 4)

**Branch:** `feature/multi-instance-dashboard`  
**Worktree:** `/home/denny/Project/headscale-multi-instance-dashboard`  
**Base Commit:** `08c261007d293f2902339d584bfb5c3df629809d`

## Objective
Add UI and backend support for managing multiple Headscale instances from a single dashboard.

## Expected Files
- `headplane/app/routes/instances/` (new) — instance management UI
- `headplane/app/server/instances.ts` (new) — multi-instance API
- `hscontrol/api/v1/instance.go` (optional) — instance metadata endpoint
- `docs/ref/integration/web-ui.md` (update) — multi-instance setup
- Tests in `headplane/tests/` and `hscontrol/`

## Implementation Steps
1. Navigate: `cd /home/denny/Project/headscale-multi-instance-dashboard`
2. Design multi-instance configuration (how users configure multiple backends)
3. Implement instance switcher UI and instance management page
4. Add backend support for instance configuration storage
5. Implement secure credential handling
6. Add tests (unit + integration)
7. Update documentation
8. Test: `cd headplane && pnpm test && pnpm build` and `cd .. && make test`
9. Commit: `git add -A && git commit -m "feat: add multi-instance dashboard" && git push -u origin feature/multi-instance-dashboard`

## Constraints
- Work only in this worktree
- Secure credential storage required
- Follow existing patterns
- No unrelated refactoring

## Acceptance Criteria
- Users can configure multiple Headscale instances
- UI provides instance switcher
- Credentials stored securely
- Docs cover multi-instance setup
- All tests pass

See `docs/parallel-development.md` for workflow details.

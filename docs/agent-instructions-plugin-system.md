# Agent Instructions: Plugin System (Task 1 of 4)

**Branch:** `feature/plugin-system`  
**Worktree:** `/home/denny/Project/headscale-plugin-system`  
**Base Commit:** `08c261007d293f2902339d584bfb5c3df629809d`

## Objective
Design and implement a plugin API for third-party UI components in Headplane.

## Expected Files
- `headplane/app/plugins/` (new) — plugin system
- `headplane/examples/example-plugin/` (new) — example
- `docs/ref/extending.md` (update) — plugin guide
- `headplane/docs/development/plugins.md` (new)
- Tests in `headplane/tests/`

## Implementation Steps
1. Navigate: `cd /home/denny/Project/headscale-plugin-system`
2. Design plugin API (TypeScript interfaces for plugin contracts)
3. Implement plugin loader and registry
4. Create example plugin
5. Add tests (unit + component)
6. Update documentation
7. Test: `cd headplane && pnpm install && pnpm typecheck && pnpm test && pnpm build`
8. Commit: `git add -A && git commit -m "feat: add plugin system" && git push -u origin feature/plugin-system`

## Constraints
- Work only in this worktree
- Follow existing Headplane patterns
- No unrelated refactoring
- Test independently

## Acceptance Criteria
- Plugin system loads plugins from directory/config
- TypeScript interfaces defined
- Example plugin works
- Docs cover plugin development
- All tests pass

See `docs/parallel-development.md` for workflow details.

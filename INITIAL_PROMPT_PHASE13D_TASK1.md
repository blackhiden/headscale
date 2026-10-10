# Phase 13d Initial Prompt - Task 1: Extend Backend Settings API

You are implementing **Phase 13d Task 1** for the Headscale project.

## Your Task

Extend the Headplane settings API to support updating the local administrator's username, name, and avatar fields, with all changes persisting atomically to `config.yaml`.

## Context

- **Project:** Headscale (self-hosted Tailscale control server) + Headplane (web UI)
- **Phase:** 13d — Restore `/admin/admin/users` Editable User Profile UI
- **Your worktree:** `/home/denny/Project/headscale-project/headscale-extend-settings-api`
- **Your branch:** `feature/extend-settings-api`
- **Base commit:** `8cb0d906`

## What Currently Exists

1. Settings API endpoints:
   - `GET /api/v1/headplane/settings` — retrieves settings
   - `POST /api/v1/headplane/settings` — updates settings (currently theme only)
   - `POST /api/v1/headplane/change-password` — changes password

2. Config structure (`config.yaml`):
   ```yaml
   user:
     username: admin
     password: <bcrypt-hash>
     # Missing: name, avatar
   ```

3. Relevant files:
   - `headplane/app/server/config/config-schema.ts` — config schema
   - `headplane/app/routes/api/headplane-settings.ts` — settings API
   - `headplane/app/server/config/config-loader.ts` — config loading
   - Tests in `headplane/tests/`

## What You Need to Do

1. **Extend config schema** to support `user.name` and `user.avatar` (optional fields)
2. **Extend POST `/api/v1/headplane/settings`** to accept username, name, avatar
3. **Add validation:**
   - Username: non-empty string
   - Avatar: valid URL or empty string
4. **Ensure atomic persistence** to `config.yaml` (no partial writes)
5. **Handle missing fields** gracefully (backward compatibility)
6. **Add/update tests** for new fields and validation
7. **Update API documentation**

## Expected API Contract

### GET `/api/v1/headplane/settings` Response:
```json
{
  "username": "admin",
  "name": "Administrator",
  "avatar": "https://example.com/avatar.png",
  "theme": "dark",
  "apiKey": "..."
}
```

### POST `/api/v1/headplane/settings` Request:
```json
{
  "username": "admin",
  "name": "Administrator",
  "avatar": "https://example.com/avatar.png",
  "theme": "dark"
}
```

## Important Constraints

- Must not break existing settings functionality
- Must maintain backward compatibility with old config files (missing name/avatar)
- Must handle atomic writes (no race conditions or partial writes)
- Follow existing Headplane code conventions (TypeScript, React Router 7, Vite)
- Read relevant documentation before starting:
  - `AGENTS.md` — coding conventions
  - `headplane/docs/ARCHITECTURE.md` — Headplane architecture
  - `ROADMAP.md` Phase 13d section

## Steps to Start

1. **Change to your worktree:**
   ```bash
   cd /home/denny/Project/headscale-project/headscale-extend-settings-api
   ```

2. **Read the agent instructions:**
   - Read `AGENT_INSTRUCTIONS_PHASE13D_TASK1.md` in the main repo

3. **Explore existing code:**
   - Read `headplane/app/server/config/config-schema.ts`
   - Read `headplane/app/routes/api/headplane-settings.ts`
   - Read existing tests

4. **Implement changes:**
   - Update config schema
   - Extend settings API
   - Add validation
   - Update tests
   - Update docs

5. **Validate:**
   - Run tests: `cd headplane && pnpm test`
   - Run typecheck: `cd headplane && pnpm typecheck`
   - Test manually if possible

6. **Commit:**
   ```bash
   git add <files>
   git commit -m "feat: extend settings API for username, name, avatar"
   git push -u origin feature/extend-settings-api
   ```

## Acceptance Criteria

- [ ] Config schema supports `user.name` and `user.avatar`
- [ ] POST endpoint accepts and validates new fields
- [ ] Changes persist atomically to `config.yaml`
- [ ] Backward compatibility maintained
- [ ] Tests pass
- [ ] Documentation updated

## Questions?

If you need clarification, ask the user before proceeding.

---

**Ready to start? Change to your worktree and begin!**

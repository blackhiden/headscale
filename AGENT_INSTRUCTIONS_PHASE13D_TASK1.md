# Phase 13d Task 1: Extend Backend Settings API

**Branch:** `feature/extend-settings-api`  
**Worktree:** `/home/denny/Project/headscale-project/headscale-extend-settings-api`  
**Base Commit:** `8cb0d906` (dev branch)

---

## Objective

Extend the Headplane settings API to support updating the local administrator's username, name, and avatar fields, persisting all changes atomically to `config.yaml`.

---

## Current State

- Existing endpoints:
  - `GET /api/v1/headplane/settings` — retrieves user settings (API key, theme, profile name)
  - `POST /api/v1/headplane/settings` — updates settings (currently theme only)
  - `POST /api/v1/headplane/change-password` — changes password with current password verification

- Current implementation:
  - Settings API in `headplane/app/routes/api/headplane-settings.ts` (Remix/React Router)
  - Headscale backend handlers likely in `hscontrol/` or similar
  - Config schema in `headplane/app/server/config/config-schema.ts`

- Current config structure (`config.yaml`):
  ```yaml
  user:
    username: admin
    password: <bcrypt-hash>
    # Missing: name, avatar fields
  ```

---

## Requirements

1. **Extend config schema** to support `user.name` and `user.avatar` (optional fields)
2. **Extend POST `/api/v1/headplane/settings`** to accept username, name, and avatar
3. **Validate** input fields (username: non-empty, avatar: valid URL or empty)
4. **Persist atomically** to `config.yaml` (no partial writes)
5. **Handle missing fields** gracefully (backward compatibility)
6. **Update API documentation** (OpenAPI, comments, or similar)
7. **Add/update tests** for new fields and validation

---

## Expected Files to Modify/Create

- `headplane/app/server/config/config-schema.ts` — add `name` and `avatar` to user schema
- `headplane/app/routes/api/headplane-settings.ts` — extend POST handler
- `headplane/app/server/config/config-loader.ts` — ensure new fields are loaded
- `headplane/app/server/config/config-writer.ts` (or similar) — ensure new fields are written
- Tests:
  - `headplane/tests/unit/server/config/config-schema.test.ts` (or similar)
  - `headplane/tests/integration/api/headplane-settings.test.ts` (or similar)
- Docs:
  - `docs/ref/api/headplane-settings.md` — update API reference
  - `docs/usage/settings.md` — update user guide

---

## Acceptance Criteria

- [ ] Config schema supports `user.name` and `user.avatar` (optional)
- [ ] `POST /api/v1/headplane/settings` accepts username, name, avatar
- [ ] Input validation works (username non-empty, avatar valid URL or empty)
- [ ] Changes persist atomically to `config.yaml`
- [ ] Missing/partial config handled gracefully
- [ ] Tests pass (unit + integration)
- [ ] API documentation updated

---

## Testing

1. **Unit tests:**
   - Config schema validation for new fields
   - Settings handler validation logic

2. **Integration tests:**
   - Update username, name, avatar via API
   - Verify persistence to `config.yaml`
   - Test missing/partial fields
   - Test validation errors

3. **Manual validation:**
   - Start Headplane with extended config
   - Call API to update fields
   - Inspect `config.yaml` for changes
   - Verify no partial writes on errors

---

## Dependencies

- **None** (this is the first task in the sequence)

---

## Sequence

1. Update config schema
2. Extend settings API handler
3. Add validation logic
4. Add/update tests
5. Update documentation
6. Commit and push

---

## Constraints

- Must not break existing settings API functionality
- Must maintain backward compatibility with old config files
- Must handle atomic writes (no race conditions or partial writes)
- Must follow existing code conventions (see AGENTS.md)

---

## Notes

- This is a **backend-focused** task
- The frontend UI task (Task 2) depends on this API contract
- Review `headplane/app/server/config/` directory structure before starting
- Follow Headplane coding conventions (TypeScript, Vite, React Router 7)
- Read `headplane/docs/` for Headplane-specific architecture

---

## Completion Checklist

- [ ] Config schema updated
- [ ] API handler extended
- [ ] Validation added
- [ ] Tests written and passing
- [ ] Documentation updated
- [ ] Changes committed to `feature/extend-settings-api`
- [ ] Ready for review

# Phase 13d Initial Prompt - Task 2: Restore Frontend User Profile UI

You are implementing **Phase 13d Task 2** for the Headscale project.

## Your Task

Restore the `/admin/admin/users` route in Headplane with a comprehensive user profile edit form that allows the local administrator to update their username, password, name, and avatar.

## Context

- **Project:** Headscale (self-hosted Tailscale control server) + Headplane (web UI)
- **Phase:** 13d — Restore `/admin/admin/users` Editable User Profile UI
- **Your worktree:** `/home/denny/Project/headscale-project/headscale-restore-profile-ui`
- **Your branch:** `feature/restore-profile-ui`
- **Base commit:** `8cb0d906`

## What Currently Exists

1. The `/admin/admin/users` route was **removed in Phase 13c**
2. Current admin UI is at `/admin` but lacks editable profile
3. Settings page has password change functionality at `/settings`
4. API client exists at `headplane/app/server/headscale/api/settings.ts`

## What You Need to Do

1. **Restore `/admin/admin/users` route** with profile edit form
2. **Implement form fields:**
   - Username (text input, required)
   - Name (text input, optional)
   - Avatar (URL input, optional)
   - Password (with current password verification)
3. **Integrate with backend API:**
   - `GET /api/v1/headplane/settings` — load profile
   - `POST /api/v1/headplane/settings` — update username, name, avatar
   - `POST /api/v1/headplane/change-password` — update password
4. **Add UI validation:**
   - Username: non-empty
   - Avatar: valid URL or empty
   - Password: current + new + confirmation
5. **Handle missing config** gracefully (empty fields, no errors)
6. **Update navigation** with link to profile page
7. **Add access control** (admin-only)

## Expected API Contract (from Task 1)

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

### POST `/api/v1/headplane/change-password` Request:
```json
{
  "currentPassword": "...",
  "newPassword": "...",
  "confirmPassword": "..."
}
```

## Important Constraints

- Follow existing Headplane UI/UX patterns
- Reuse existing components where possible
- Do NOT reintroduce multi-user UI elements
- Handle missing config gracefully
- Follow Headplane coding conventions (React Router 7, TypeScript, Vite)
- Read relevant documentation before starting:
  - `AGENTS.md` — coding conventions
  - `headplane/docs/ARCHITECTURE.md` — Headplane architecture
  - `ROADMAP.md` Phase 13d section

## Dependencies

- **Task 1 (Backend API)** should be complete before full integration
- You can start UI scaffolding with stubbed/mocked API if Task 1 is still in progress
- Coordinate with Task 1 developer on API contract

## Steps to Start

1. **Change to your worktree:**
   ```bash
   cd /home/denny/Project/headscale-project/headscale-restore-profile-ui
   ```

2. **Read the agent instructions:**
   - Read `AGENT_INSTRUCTIONS_PHASE13D_TASK2.md` in the main repo

3. **Explore existing code:**
   - Read `headplane/app/routes/admin/` for admin UI patterns
   - Read `headplane/app/routes/settings/` for password change patterns
   - Read `headplane/app/components/` for reusable components
   - Check `headplane/app/server/headscale/api/settings.ts` for API client

4. **Implement changes:**
   - Create route structure: `headplane/app/routes/admin/admin/users/route.tsx`
   - Create profile form component
   - Integrate with API client (stub if Task 1 not ready)
   - Add validation logic
   - Update navigation
   - Add tests

5. **Validate:**
   - Run tests: `cd headplane && pnpm test`
   - Run typecheck: `cd headplane && pnpm typecheck`
   - Run dev server: `cd headplane && pnpm dev`
   - Test manually in browser

6. **Commit:**
   ```bash
   git add <files>
   git commit -m "feat: restore /admin/admin/users profile UI"
   git push -u origin feature/restore-profile-ui
   ```

## Acceptance Criteria

- [ ] `/admin/admin/users` route accessible
- [ ] Profile form displays and edits all fields
- [ ] API integration complete
- [ ] Validation works
- [ ] Missing config handled gracefully
- [ ] Navigation updated
- [ ] Access control enforced
- [ ] Tests pass
- [ ] Manual testing complete

## UI/UX Guidelines

- Use existing Headplane form components and patterns
- Follow existing validation and error display patterns
- Use existing button styles and layouts
- Match existing admin page styling
- Show loading states during API calls
- Display success/error messages clearly

## Questions?

If you need clarification or API contract details from Task 1, ask the user before proceeding.

---

**Ready to start? Change to your worktree and begin!**

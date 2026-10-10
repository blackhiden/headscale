# Phase 13d Task 2: Restore Frontend User Profile UI

**Branch:** `feature/restore-profile-ui`  
**Worktree:** `/home/denny/Project/headscale-project/headscale-restore-profile-ui`  
**Base Commit:** `8cb0d906` (dev branch)

---

## Objective

Restore the `/admin/admin/users` route in Headplane with a comprehensive user profile edit form that allows the local administrator to update their username, password, name, and avatar. All changes must integrate with the extended backend settings API.

---

## Current State

- **Route removed:** `/admin/admin/users` was removed in Phase 13c (see `headplane/SESSION2_COMPLETE.md`)
- **Current admin UI:** `/admin` exists but lacks editable profile
- **Existing components:** Settings page has password change functionality
- **API client:** `headplane/app/server/headscale/api/settings.ts` exists for settings API

---

## Requirements

1. **Restore `/admin/admin/users` route** in Headplane
2. **Implement editable profile form** with fields:
   - Username (text input, required)
   - Name (text input, optional)
   - Avatar (URL input or upload, optional)
   - Password (password input with current password verification)
3. **Integrate with backend API:**
   - `GET /api/v1/headplane/settings` — load current profile
   - `POST /api/v1/headplane/settings` — update username, name, avatar
   - `POST /api/v1/headplane/change-password` — update password
4. **Add UI validation:**
   - Username: non-empty
   - Avatar: valid URL or empty
   - Password: confirm current password + new password + confirmation match
5. **Handle missing/partial config** gracefully (show empty fields, no errors)
6. **Update navigation** to include link to profile page
7. **Add access control** (admin-only)

---

## Expected Files to Modify/Create

- `headplane/app/routes/admin/admin/users/route.tsx` — restore/create profile UI route
- `headplane/app/routes/admin/admin/users/components/` — profile form components
- `headplane/app/server/headscale/api/settings.ts` — extend API client for new fields
- `headplane/app/types/settings.ts` (or similar) — add profile field types
- `headplane/app/layout/header.tsx` — update navigation links if needed
- Tests:
  - `headplane/tests/unit/routes/admin-users.test.ts` — profile form unit tests
  - `headplane/tests/integration/routes/admin-users.test.tsx` — profile UI integration tests

---

## Acceptance Criteria

- [ ] `/admin/admin/users` route restored and accessible
- [ ] Profile form displays current username, name, avatar
- [ ] Username, name, avatar can be edited and saved
- [ ] Password can be changed with current password verification
- [ ] UI validation works (username non-empty, avatar valid URL, password match)
- [ ] Missing/partial config handled gracefully (empty fields, no errors)
- [ ] Navigation updated (link to profile page)
- [ ] Access control enforced (admin-only)
- [ ] Tests pass (unit + integration)
- [ ] UI follows existing Headplane design patterns

---

## Testing

1. **Unit tests:**
   - Profile form component rendering
   - Form validation logic
   - API client methods for new fields

2. **Integration tests:**
   - Load profile data from API
   - Update username, name, avatar via form
   - Change password via form
   - Validation error display
   - Missing/partial config handling

3. **Manual validation:**
   - Navigate to `/admin/admin/users`
   - Edit and save username, name, avatar
   - Change password
   - Verify changes persist (reload page, check API)
   - Test with missing config fields

---

## Dependencies

- **Task 1 (Backend API)** must be complete before full integration
- Can start UI scaffolding with stubbed/mocked API during Task 1 development
- API contract must be agreed upon before final integration

---

## Sequence

1. Restore route structure (`route.tsx`)
2. Create profile form component
3. Integrate with API client (stub/mock if Task 1 not ready)
4. Add validation logic
5. Update navigation
6. Add tests
7. Manual validation
8. Commit and push

---

## Constraints

- Must follow existing Headplane UI/UX patterns
- Must reuse existing components where possible
- Must not reintroduce multi-user UI elements
- Must handle missing config gracefully
- Must follow Headplane coding conventions (React Router 7, TypeScript, Vite)

---

## Notes

- This is a **frontend-focused** task
- Review `headplane/app/routes/admin/` for existing admin UI patterns
- Review `headplane/app/routes/settings/` for password change patterns
- Review `headplane/app/components/` for reusable form components
- Follow Headplane design system and component patterns
- Read `headplane/docs/ARCHITECTURE.md` for frontend architecture

---

## API Contract (from Task 1)

### GET `/api/v1/headplane/settings`
Returns:
```json
{
  "username": "admin",
  "name": "Administrator",
  "avatar": "https://example.com/avatar.png",
  "theme": "dark",
  "apiKey": "..."
}
```

### POST `/api/v1/headplane/settings`
Accepts:
```json
{
  "username": "admin",
  "name": "Administrator",
  "avatar": "https://example.com/avatar.png",
  "theme": "dark"
}
```

### POST `/api/v1/headplane/change-password`
Accepts:
```json
{
  "currentPassword": "...",
  "newPassword": "...",
  "confirmPassword": "..."
}
```

---

## Completion Checklist

- [ ] Route restored
- [ ] Profile form implemented
- [ ] API integration complete
- [ ] Validation added
- [ ] Navigation updated
- [ ] Tests written and passing
- [ ] Manual validation complete
- [ ] Changes committed to `feature/restore-profile-ui`
- [ ] Ready for review

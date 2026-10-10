# Phase 13b — Settings Menu (Single User): Implementation Plan

**Phase:** Phase 13b  
**Status:** Planning  
**Date:** 2026-10-03  
**Base Commit:** `8084bfb8` (dev branch)

---

## Executive Summary

Phase 13b adds a Settings menu to Headplane for the authenticated user (from Phase 13a) to:
- Store and manage their Headscale API key (no re-entry at each login)
- Change their Headplane password
- Select theme preference (dark/light)
- Set optional profile name

This is a **single-user** enhancement. Multi-user support (Phase 13c) is deferred.

---

## Current State

### What Exists (Phase 13a — Complete)
✅ Single password authentication for Headplane  
✅ `POST /api/v1/headplane/login` endpoint  
✅ Session management (24-hour expiry, secure cookies)  
✅ Rate limiting (5 attempts/min/IP)  
✅ Password/API key toggle in login UI  

### What's Missing (Phase 13b — This Plan)
❌ Settings page/route  
❌ API key storage (must re-enter at each login)  
❌ Password change capability  
❌ Theme selection/persistence  
❌ Profile name  

---

## Architecture Decisions

### 1. Settings Storage

**Decision:** SQLite table in Headplane's database (already used for audit log).

**Schema:**
```sql
CREATE TABLE headplane_settings (
  id INTEGER PRIMARY KEY CHECK (id = 1),  -- Single row constraint
  api_key_encrypted TEXT,                 -- AES-256-GCM encrypted
  api_key_nonce TEXT,                     -- Encryption nonce
  theme TEXT DEFAULT 'light',             -- 'light' or 'dark'
  profile_name TEXT,                      -- Optional display name
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Encryption:** API key encrypted with AES-256-GCM using PBKDF2-derived key from session token.

**Rationale:**
- Existing SQLite database (no new infrastructure)
- Single row for single-user scenario
- Encrypted API key at rest
- Simple schema, easy to extend in Phase 13c

### 2. API Endpoints

**New endpoints:**

- `GET /api/v1/headplane/settings` — retrieve settings (decrypts API key)
- `POST /api/v1/headplane/settings` — update settings (encrypts API key)
- `POST /api/v1/headplane/change-password` — change password with verification

**Authentication:** All require valid password session token (from Phase 13a login).

### 3. Theme Implementation

**Decision:** CSS variables + `data-theme` attribute on `<html>` element.

**Implementation:**
- Frontend stores theme in settings
- On load, apply theme immediately to avoid flash
- Theme persisted server-side (survives logout)

### 4. Password Change Flow

**Flow:**
1. User submits current password + new password
2. Backend validates current password (constant-time)
3. Backend updates config password (in-memory, persists to config file)
4. All existing sessions remain valid
5. User sees success message

**Note:** Password stored in config file (`headplane.password`), not database. Change requires write access to config or env var update + restart.

---

## Implementation Tasks

### Part 1: Backend — Settings Storage & Endpoints

**Files to modify:**
- `hscontrol/headplane_settings.go` (new)
- `hscontrol/headplane_settings_test.go` (new)
- `hscontrol/app.go` (register endpoints)
- `hscontrol/types/config.go` (if needed)

**Tasks:**
1. Create `headplane_settings` SQLite table schema
2. Add table creation to Headplane database init
3. Implement settings CRUD (get/update encrypted API key, theme, profile)
4. Implement encryption/decryption for API key (AES-256-GCM)
5. Implement password change endpoint
6. Write unit tests (settings CRUD, encryption, password change)
7. Update OpenAPI spec

**Dependencies:** Phase 13a (password auth, session management).

**Estimated effort:** 1-2 days

---

### Part 2: Backend — Session Enhancement

**Files to modify:**
- `hscontrol/headplane_auth.go` (load API key from settings)
- `hscontrol/headplane_auth_test.go` (update tests)

**Tasks:**
1. Modify session validation to load stored API key if present
2. Inject API key into auth principal so frontend API calls work seamlessly
3. Update tests for new behavior

**Dependencies:** Part 1 (settings storage must exist).

**Estimated effort:** 0.5 days

---

### Part 3: Frontend — Settings Route & Components

**Files to create:**
- `headplane/app/routes/settings/profile/page.tsx` (new)
- `headplane/app/routes/settings/profile/components/` (new dir)

**Files to modify:**
- `headplane/app/routes/settings/overview.tsx` (add "Profile" link)
- `headplane/app/server/headscale/api/index.ts` (add settings methods)
- `headplane/app/server/headscale/api/transport.ts` (if needed)

**Tasks:**
1. Create `/settings/profile` route
2. Create settings sections:
   - Account (change password, session info, logout)
   - Integration (API key input with save/update, masked display)
   - Preferences (theme selector: light/dark)
   - Profile (display name input)
3. Implement forms with validation
4. Add "Settings" link to nav bar (visible after password login)
5. Implement theme application on page load
6. Add success/error feedback for all operations

**Dependencies:** Part 1 (backend endpoints must exist).

**Estimated effort:** 2-3 days

---

### Part 4: Frontend — Theme System

**Files to modify:**
- `headplane/app/root.tsx` (apply theme on load)
- `headplane/app/styles/` (CSS variables for themes)
- `headplane/app/routes/settings/profile/components/theme-selector.tsx` (new)

**Tasks:**
1. Define CSS variables for light/dark themes
2. Apply `data-theme` attribute on `<html>` based on settings
3. Implement theme selector component
4. Ensure theme persists across sessions
5. Test theme switching (no flash, immediate application)

**Dependencies:** Part 1 (settings storage), Part 3 (settings route).

**Estimated effort:** 1 day

---

### Part 5: Documentation

**Files to create/modify:**
- `docs/usage/settings.md` (new)
- `docs/usage/authentication.md` (update with settings reference)
- `README.md` (update features list)
- `CHANGELOG.md` (Phase 13b entries)

**Tasks:**
1. Document settings page usage
2. Document API key storage and security
3. Document password change process
4. Document theme selection
5. Update authentication guide with settings link
6. Add Phase 13b to changelog

**Dependencies:** All implementation parts complete.

**Estimated effort:** 0.5 days

---

### Part 6: Testing

**Files to create:**
- `headplane/tests/unit/settings/` (new dir with tests)
- `headplane/tests/e2e/settings.spec.ts` (new)
- `hscontrol/headplane_settings_test.go` (created in Part 1)

**Tasks:**
1. Backend unit tests (settings CRUD, encryption, password change)
2. Frontend unit tests (components, forms, validation)
3. E2E tests (settings flow: login → save API key → logout → login → API key loaded)
4. Integration tests (theme persistence, password change)
5. Security tests (encryption, CSRF, validation)

**Dependencies:** All implementation parts complete.

**Estimated effort:** 1-2 days

---

## Parallelization Analysis

### Safe for Parallel Development

✅ **Part 1 (Backend Settings)** and **Part 3 (Frontend UI)** can proceed in parallel:
- Frontend can mock backend endpoints initially
- Both teams coordinate on API contract (OpenAPI spec)

✅ **Part 5 (Documentation)** can start early with structure/drafts

### Must Be Sequential

⚠️ **Part 2 (Session Enhancement)** depends on Part 1  
⚠️ **Part 4 (Theme System)** depends on Parts 1 & 3  
⚠️ **Part 6 (Testing)** depends on all implementation parts  

### Recommended Approach

**Option A: Sequential (Single Developer/Agent)**
1. Part 1 (Backend Settings) — 1-2 days
2. Part 2 (Session Enhancement) — 0.5 days
3. Part 3 (Frontend UI) — 2-3 days
4. Part 4 (Theme System) — 1 day
5. Part 6 (Testing) — 1-2 days
6. Part 5 (Documentation) — 0.5 days

**Total: 6-9 days**

**Option B: Parallel (2-3 Developers/Agents)**

**Agent 1 (Backend):**
- Part 1 (Backend Settings)
- Part 2 (Session Enhancement)
- Backend testing

**Agent 2 (Frontend):**
- Part 3 (Frontend UI, mock backend initially)
- Part 4 (Theme System)
- Frontend testing

**Agent 3 (Docs/Integration):**
- Part 5 (Documentation drafts)
- Integration testing after Parts 1-4 complete

**Total: 3-5 days (with coordination)**

---

## Git Worktree Strategy (If Parallel)

If using parallel development:

1. Create worktrees from current `dev` branch (`8084bfb8`)
2. Create feature branches:
   - `feature/settings-backend`
   - `feature/settings-frontend`
   - `feature/settings-docs`
3. Each agent works in isolated worktree
4. Merge order: backend → frontend → docs

**Base commit for all branches:** `8084bfb8`

---

## Acceptance Criteria

### Backend
- [ ] Settings table created in Headplane database
- [ ] API key encrypted at rest with AES-256-GCM
- [ ] `GET /api/v1/headplane/settings` returns settings (decrypted API key)
- [ ] `POST /api/v1/headplane/settings` updates settings (encrypted storage)
- [ ] `POST /api/v1/headplane/change-password` validates and changes password
- [ ] Session loads stored API key automatically
- [ ] All backend unit tests pass

### Frontend
- [ ] `/settings/profile` route exists and is accessible after login
- [ ] Settings page has sections: Account, Integration, Preferences, Profile
- [ ] API key input with masked display and save functionality
- [ ] Password change form with current password verification
- [ ] Theme selector (light/dark) with immediate application
- [ ] Profile name input with save functionality
- [ ] "Settings" link visible in nav after login
- [ ] All forms have validation and error handling
- [ ] All frontend unit and E2E tests pass

### Integration
- [ ] User can log in, save API key, log out, log in again → API key loaded
- [ ] User can change password and log in with new password
- [ ] Theme persists across logout/login
- [ ] Profile name persists across logout/login

### Security
- [ ] API key encrypted at rest (verified in tests)
- [ ] Password change validates current password
- [ ] CSRF protection on all state-changing endpoints
- [ ] No secrets logged or exposed in responses

### Documentation
- [ ] Settings page documented in `docs/usage/settings.md`
- [ ] Authentication guide updated with settings reference
- [ ] Security considerations documented
- [ ] README updated with new features
- [ ] CHANGELOG updated with Phase 13b

---

## Risks & Mitigation

| Risk | Impact | Mitigation |
|------|--------|------------|
| Encryption key derivation from session token is fragile | High | Use PBKDF2 with salt, test key rotation scenarios |
| Password change requires config file write | Medium | Document admin-only capability, consider env var override |
| Theme flash on page load | Low | Apply theme synchronously before first render |
| API key exposed in frontend | High | Mask by default, require explicit "reveal" action |
| Backward compat with Phase 13a | Medium | Ensure API key login still works, settings are optional |

---

## Next Steps

1. **Review and approve this plan**
2. **Choose implementation approach** (sequential vs. parallel)
3. **Create Git worktrees** (if parallel)
4. **Begin Part 1** (Backend Settings)

---

## References

- **ROADMAP.md** — Phase 13b requirements
- **Phase 13a Implementation** — `hscontrol/headplane_auth.go`
- **Headplane Database** — `headplane/app/server/db/schema.ts`
- **Existing Settings Routes** — `headplane/app/routes/settings/`

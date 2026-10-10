# Phase 13b Parallel Development — Coordination Document

**Date:** 2026-10-03  
**Base Commit:** `d95c5eff`  
**Status:** Ready for parallel implementation

---

## Overview

Three agents will work in parallel on Phase 13b (Settings Menu):

| Agent | Task | Worktree | Branch | Instructions |
|-------|------|----------|--------|-------------|
| **Agent 1** | Backend Settings | `/home/denny/Project/headscale-settings-backend` | `feature/settings-backend` | `AGENT_INSTRUCTIONS.md` |
| **Agent 2** | Frontend UI | `/home/denny/Project/headscale-settings-frontend` | `feature/settings-frontend` | `AGENT_INSTRUCTIONS.md` |
| **Agent 3** | Docs & Testing | `/home/denny/Project/headscale-settings-docs` | `feature/settings-docs` | `AGENT_INSTRUCTIONS.md` |

---

## Worktree Structure

```
/home/denny/Project/
├── headscale/                          # Main worktree (dev branch)
├── headscale-settings-backend/        # Agent 1 worktree
├── headscale-settings-frontend/       # Agent 2 worktree
└── headscale-settings-docs/           # Agent 3 worktree
```

---

## Agent Responsibilities

### Agent 1: Backend Settings

**Files to create:**
- `hscontrol/headplane_settings.go` (~250 lines)
- `hscontrol/headplane_settings_test.go` (~200 lines)

**Files to modify:**
- `hscontrol/app.go` (register endpoints)
- `hscontrol/headplane_auth.go` (password change, session enhancement)
- `hscontrol/headplane_auth_test.go` (update tests)

**Deliverables:**
- SQLite `headplane_settings` table
- 3 REST endpoints (GET/POST settings, POST change-password)
- AES-256-GCM encryption for API key
- Unit tests (9+ test cases)
- No regressions

**Estimated:** 1-2 days

---

### Agent 2: Frontend UI

**Files to create:**
- `headplane/app/routes/settings/profile/page.tsx`
- `headplane/app/routes/settings/profile/components/*.tsx` (5-6 components)
- `headplane/tests/unit/settings/*.test.ts`
- `headplane/tests/e2e/settings.spec.ts`

**Files to modify:**
- `headplane/app/server/headscale/api/index.ts` (add methods)
- `headplane/app/server/headscale/api/transport.ts` (implement)
- `headplane/app/root.tsx` (theme application)
- `headplane/app/routes/settings/overview.tsx` (add link)
- Navigation component (add Settings link)

**Deliverables:**
- Settings page with 4 sections (Account, Integration, Preferences, Profile)
- Theme system (light/dark with persistence)
- Form validation and error handling
- Unit and E2E tests
- Responsive, accessible UI

**Estimated:** 2-3 days

---

### Agent 3: Documentation & Integration

**Files to create:**
- `docs/usage/settings.md`
- `docs/ref/api/headplane-settings.md`
- `docs/phase13b-integration-testing.md`

**Files to modify:**
- `docs/usage/authentication.md` (add settings reference)
- `README.md` (update features)
- `CHANGELOG.md` (Phase 13b entries)
- `docs/phase13b-implementation-plan.md` (add results)

**Deliverables:**
- Comprehensive user documentation
- API documentation
- Integration test results
- Screenshots (if needed)
- Known issues documented

**Estimated:** 0.5 days drafting + 1 day testing after agents 1 & 2 complete

---

## Coordination Points

### API Contract (Backend ↔ Frontend)

**Agreed endpoints:**

```
GET /api/v1/headplane/settings
Authorization: Bearer <session-token>
→ {apiKey: string, theme: string, profileName: string}

POST /api/v1/headplane/settings
Authorization: Bearer <session-token>
Body: {apiKey?: string, theme?: string, profileName?: string}
→ {success: true}

POST /api/v1/headplane/change-password
Authorization: Bearer <session-token>
Body: {currentPassword: string, newPassword: string}
→ {success: true}
```

**Frontend can mock these initially** using stub responses.

### Theme System (Backend ↔ Frontend)

- **Backend:** Stores theme as string (`'light'` or `'dark'`) in settings table
- **Frontend:** Applies via `data-theme` attribute on `<html>` element with CSS variables

### Session Enhancement (Backend → Frontend)

- **Backend:** Loads stored API key on session validation, injects into context
- **Frontend:** API calls automatically use stored key (transparent to UI)

---

## Development Workflow

### Phase 1: Independent Development (Day 1-3)

**Agent 1 (Backend):**
```bash
cd /home/denny/Project/headscale-settings-backend
cat AGENT_INSTRUCTIONS.md
git branch --show-current  # verify: feature/settings-backend
# Implement backend per instructions
go test ./hscontrol -v -run TestHeadplane
git add .
git commit -m "backend: implement settings storage and endpoints"
git push -u origin feature/settings-backend
```

**Agent 2 (Frontend):**
```bash
cd /home/denny/Project/headscale-settings-frontend
cat AGENT_INSTRUCTIONS.md
git branch --show-current  # verify: feature/settings-frontend
# Implement frontend per instructions (mock backend initially)
cd headplane
pnpm test:unit && pnpm test:e2e
cd ..
git add .
git commit -m "frontend: implement settings UI and theme system"
git push -u origin feature/settings-frontend
```

**Agent 3 (Docs):**
```bash
cd /home/denny/Project/headscale-settings-docs
cat AGENT_INSTRUCTIONS.md
git branch --show-current  # verify: feature/settings-docs
# Draft documentation per instructions
git add docs/
git commit -m "docs: draft Phase 13b settings documentation"
# Wait for agents 1 & 2 before integration testing
```

---

### Phase 2: Integration Testing (Day 4-5)

**After agents 1 & 2 complete:**

1. **Agent 3 sets up integration environment:**
   ```bash
   cd /home/denny/Project/headscale-settings-docs
   
   # Merge backend branch locally
   git merge feature/settings-backend --no-commit
   
   # Merge frontend branch locally
   git merge feature/settings-frontend --no-commit
   
   # Build and test
   go build -o headscale cmd/headscale/main.go
   cd headplane && pnpm install && pnpm dev
   ```

2. **Run integration test scenarios** (see `docs/phase13b-integration-testing.md`)

3. **Document results:**
   - All test scenarios pass/fail
   - Screenshots for documentation
   - Known issues
   - Recommendations

4. **Update documentation** with accurate details from implementations

5. **Commit final docs:**
   ```bash
   git add docs/
   git commit -m "docs: complete Phase 13b documentation with integration results"
   git push -u origin feature/settings-docs
   ```

---

### Phase 3: Review & Merge (Day 5-6)

**Merge order:** backend → frontend → docs

1. **Review backend branch:**
   ```bash
   cd /home/denny/Project/headscale
   git checkout dev
   git pull origin feature/settings-backend
   # Review, test
   git merge feature/settings-backend --no-ff
   git push origin dev
   ```

2. **Review frontend branch:**
   ```bash
   git pull origin feature/settings-frontend
   # Review, test
   git merge feature/settings-frontend --no-ff
   git push origin dev
   ```

3. **Review docs branch:**
   ```bash
   git pull origin feature/settings-docs
   # Review, test
   git merge feature/settings-docs --no-ff
   git push origin dev
   ```

4. **Final integration test on dev branch**

5. **Clean up worktrees:**
   ```bash
   git worktree remove /home/denny/Project/headscale-settings-backend
   git worktree remove /home/denny/Project/headscale-settings-frontend
   git worktree remove /home/denny/Project/headscale-settings-docs
   git branch -d feature/settings-backend
   git branch -d feature/settings-frontend
   git branch -d feature/settings-docs
   ```

---

## Communication Protocol

### Daily Sync (Optional)

- **Time:** End of each day
- **Format:** Status update from each agent
- **Content:**
  - Progress summary
  - Blockers (if any)
  - Questions for other agents
  - Expected completion date

### Blockers

If blocked:
1. Document the blocker clearly
2. Notify other agents (if applicable)
3. Propose solutions or workarounds
4. Continue on non-blocked tasks

### API Changes

If backend API needs to change:
1. Backend agent proposes change with rationale
2. Frontend agent confirms feasibility
3. Update coordination document with new contract
4. Both agents implement updated contract

---

## Quality Checklist

### Before Pushing Branch

**Backend:**
- [ ] All tests pass (`go test ./...`)
- [ ] No linting errors (`make lint`)
- [ ] No regressions in Phase 13a
- [ ] Code formatted (`make fmt`)
- [ ] Commit messages follow convention

**Frontend:**
- [ ] All tests pass (`pnpm test:unit && pnpm test:e2e`)
- [ ] No type errors (`pnpm typecheck`)
- [ ] No linting errors (`pnpm lint`)
- [ ] UI works in dev mode (`pnpm dev`)
- [ ] Responsive design verified
- [ ] Accessibility checked

**Docs:**
- [ ] All links work
- [ ] Code examples are accurate
- [ ] Screenshots are current (if any)
- [ ] Spelling/grammar checked
- [ ] Integration tests documented

---

## Troubleshooting

### Merge Conflicts

Unlikely but possible if agents modify same files.

**If conflict occurs:**
1. Identify conflicting files
2. Coordinate resolution between affected agents
3. Test merged result
4. Document resolution in commit message

### Integration Failures

If integration tests fail:
1. Identify root cause (backend, frontend, or both)
2. Create issue with details
3. Assign to responsible agent(s)
4. Fix and retest
5. Update documentation with resolution

### Worktree Issues

If worktree becomes corrupted:
```bash
git worktree remove --force <path>
git worktree add -b <branch> <path>
```

---

## Success Criteria

Phase 13b is complete when:

- [ ] All three branches merged to `dev`
- [ ] All acceptance criteria met (see implementation plan)
- [ ] Integration tests pass
- [ ] Documentation complete
- [ ] No regressions in Phase 13a
- [ ] CHANGELOG updated
- [ ] Ready for release

---

## References

- **Implementation Plan:** `docs/phase13b-implementation-plan.md`
- **ROADMAP:** `ROADMAP.md` Phase 13b
- **Agent Instructions:**
  - Backend: `/home/denny/Project/headscale-settings-backend/AGENT_INSTRUCTIONS.md`
  - Frontend: `/home/denny/Project/headscale-settings-frontend/AGENT_INSTRUCTIONS.md`
  - Docs: `/home/denny/Project/headscale-settings-docs/AGENT_INSTRUCTIONS.md`

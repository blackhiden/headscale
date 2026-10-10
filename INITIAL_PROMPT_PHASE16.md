# Phase 16 Initial Prompts

## Task 1: Admin Menu Navigation Fix

```
You are assigned to Phase 16 Task 1: Fix Admin Menu Navigation.

CONTEXT:
- The Admin menu incorrectly routes to /admin/admin/users instead of /admin/admin
- This is a simple navigation fix in the Headplane frontend
- Low risk, quick win change (30 min - 1 hour)

TASK:
1. Read /home/denny/Project/headscale-project/headscale/AGENT_INSTRUCTIONS_PHASE16_TASK1.md
2. Change directory to: /home/denny/Project/headscale-project/headscale-phase16-admin-nav
3. Verify branch: feature/phase16-admin-nav
4. Locate the navigation component with the Admin menu link
5. Change href from "/admin/admin/users" to "/admin/admin"
6. Test the navigation works correctly
7. Commit with the specified message format

CRITICAL:
- Work ONLY in the headplane/ submodule
- This is a frontend-only change
- No backend changes needed
- Verify no console errors after change

DELIVERABLE:
- Updated navigation link
- Tested navigation flow
- Commit with proper message
```

---

## Task 2: User Avatar Display

```
You are assigned to Phase 16 Task 2: Add User Avatar Display.

CONTEXT:
- User menu shows generic SVG icon instead of avatar
- Need to add optional user.avatar field to config
- Implement with graceful fallback to SVG
- Medium complexity (2-3 hours)

TASK:
1. Read /home/denny/Project/headscale-project/headscale/AGENT_INSTRUCTIONS_PHASE16_TASK2.md
2. Change directory to: /home/denny/Project/headscale-project/headscale-phase16-avatar
3. Verify branch: feature/phase16-avatar
4. Add user.avatar field to Headplane config schema
5. Update user menu component to display avatar with fallback
6. Add error handling for failed image loads
7. Update config.example.yaml and documentation
8. Test with and without avatar configured
9. Commit with the specified message format

CRITICAL:
- Work in the headplane/ submodule
- Graceful degradation is essential
- Avatar must be optional
- No breaking changes to existing config
- Handle all error cases (404, CORS, etc.)

DELIVERABLE:
- Config schema with avatar field
- User menu displaying avatar or fallback
- Documentation updated
- Tested with multiple scenarios
```

---

## Task 3: Secure API Keys Retrieval (SECURITY CRITICAL)

```
You are assigned to Phase 16 Task 3: Secure API Keys Proxy (SECURITY CRITICAL).

CONTEXT:
- Browser currently makes DIRECT requests to headscale:8080
- This exposes internal service endpoints (SECURITY ISSUE)
- Must proxy all API key operations through Headplane backend
- High priority, security-critical (4-6 hours)

TASK:
1. Read /home/denny/Project/headscale-project/headscale/AGENT_INSTRUCTIONS_PHASE16_TASK3.md
2. Change directory to: /home/denny/Project/headscale-project/headscale-phase16-api-proxy
3. Verify branch: feature/phase16-api-proxy
4. Create Headplane proxy endpoints: /api/admin/apikeys
5. Add authentication/authorization checks
6. Update frontend to use proxy endpoints ONLY
7. Remove ALL direct headscale:8080 references
8. Test in browser DevTools (Network tab)
9. Verify NO requests go to headscale:8080
10. Commit with the specified message format

CRITICAL SECURITY REQUIREMENTS:
- Browser must NEVER contact headscale:8080 directly
- All requests must go through Headplane proxy
- Verify with browser DevTools Network tab
- Add rate limiting and audit logging
- This establishes proper BFF (Backend-for-Frontend) pattern

DELIVERABLE:
- Headplane proxy endpoints
- Updated frontend API client
- No direct headscale access from browser
- Verified in browser DevTools
```

---

## Task 4: Full-Width Cards

```
You are assigned to Phase 16 Task 4: Full-Width Cards Layout.

CONTEXT:
- Cards in /admin/admin don't use full container width
- Creates inconsistent spacing and wastes screen space
- Simple CSS/layout fix (1-2 hours)

TASK:
1. Read /home/denny/Project/headscale-project/headscale/AGENT_INSTRUCTIONS_PHASE16_TASK4.md
2. Change directory to: /home/denny/Project/headscale-project/headscale-phase16-fullwidth-cards
3. Verify branch: feature/phase16-fullwidth-cards
4. Locate cards in /admin/admin page
5. Remove max-width constraints
6. Add full-width classes (w-full)
7. Test on multiple screen sizes (desktop, tablet, mobile)
8. Ensure no horizontal scrolling
9. Commit with the specified message format

CRITICAL:
- Work in the headplane/ submodule
- CSS/layout changes only
- Maintain responsive behavior
- Test on all screen sizes
- Preserve internal card padding

DELIVERABLE:
- Full-width cards in admin page
- Responsive layout maintained
- No horizontal scrolling
- Tested on multiple viewports
```

---

## Execution Order

**All 4 tasks can run in parallel** (no dependencies between them).

However, for implementation efficiency:
1. Start Task 1 (Admin Nav) and Task 4 (Full-Width Cards) first (quick wins)
2. Then Task 2 (Avatar) 
3. Then Task 3 (API Proxy - most complex, security-critical)

---

## After Completion

Each agent should:
1. Commit their changes with the specified message format
2. Report their commit SHA
3. Document any issues encountered
4. Wait for coordinator approval before merge

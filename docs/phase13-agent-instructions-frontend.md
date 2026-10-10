# Phase 13 Agent Instructions: Frontend Password Login UI

## Task Assignment

**Agent:** Frontend Login UI  
**Worktree:** `/home/denny/Project/headscale-password-frontend`  
**Branch:** `feature/password-frontend`  
**Base Commit:** `b1495ddf5d26ec4089cd6459a491e33cffcf60ed`

---

## Objective

Update Headplane login form to support password-based authentication alongside the existing API key method.

---

## Roadmap Reference

See `ROADMAP.md` Phase 13 for complete requirements.

---

## Technical Requirements

### 1. Login Form Updates
- Add password input field to login form
- Add toggle or tab to switch between API key and password login
- Preserve existing API key login functionality
- Add "Remember me" option (optional)
- Add "Show/hide password" toggle

### 2. Authentication Flow
- Call backend `/api/v1/headplane/login` endpoint with password
- Handle successful login (store session token)
- Handle failed login (show error message)
- Handle rate limiting errors (show appropriate message)
- Redirect to dashboard on success

### 3. Session Management
- Store session token securely (localStorage or sessionStorage)
- Add token to subsequent API requests
- Handle token expiration
- Clear token on logout

### 4. Error Handling
- Display clear error messages for invalid password
- Display rate limiting messages
- Display network errors
- Provide helpful feedback to users

### 5. UI/UX
- Follow existing Headplane design system
- Reuse existing form components
- Maintain visual consistency
- Ensure accessibility (WCAG compliance)

---

## Expected Files to Change

- `headplane/app/components/LoginForm.tsx` or similar
- `headplane/app/pages/Login.tsx` or similar
- `headplane/app/services/auth.ts` (add password login method)
- `headplane/app/stores/authStore.ts` (update session handling)
- `headplane/tests/unit/auth/*.test.ts` (add tests)
- `headplane/tests/e2e/login.test.ts` (update E2E tests)

---

## Testing Requirements

1. Unit tests for login form component
2. Unit tests for auth service
3. Integration tests for login flow
4. E2E tests for password login
5. Test error handling
6. Test session management
7. Test accessibility

---

## Acceptance Criteria

- [ ] Login form supports both API key and password
- [ ] Password login works end-to-end
- [ ] Error messages are clear and helpful
- [ ] Session management works correctly
- [ ] All tests pass
- [ ] Accessibility requirements met
- [ ] Visual design matches existing Headplane UI
- [ ] API key login remains unchanged and functional

---

## Implementation Steps

1. Read existing login code in `headplane/app/`
2. Review design system and component library
3. Design password login UI (wireframe/mockup optional)
4. Implement form updates
5. Implement API integration
6. Implement session management
7. Add error handling
8. Write tests
9. Test manually in browser
10. Verify all tests pass
11. Commit to feature branch

---

## Development Workflow

```bash
# Navigate to worktree
cd /home/denny/Project/headscale-password-frontend

# Verify branch
git branch --show-current  # should show: feature/password-frontend

# Install dependencies (if needed)
cd headplane
npm install

# Run dev server
npm run dev

# Run tests
npm test

# Make changes, test, commit
git add -A
git commit -m "ui: add password-based login to Headplane"

# Push when ready
git push -u origin feature/password-frontend
```

---

## Dependencies

- Backend password authentication endpoint (can mock until ready)

---

## Mock Backend for Development

If backend is not ready, create a mock:

```typescript
// headplane/app/services/auth.mock.ts
export async function mockPasswordLogin(password: string) {
  await new Promise(resolve => setTimeout(resolve, 500)); // simulate network delay
  
  if (password === "test123") {
    return { token: "mock-session-token", user: "admin" };
  }
  
  throw new Error("Invalid password");
}
```

---

## Notes

- Do NOT break existing API key login
- Follow Headplane design system and patterns
- Reuse existing components where possible
- Ensure proper TypeScript types
- Follow accessibility guidelines
- Test in multiple browsers

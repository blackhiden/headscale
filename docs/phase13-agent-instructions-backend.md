# Phase 13 Agent Instructions: Backend Password Authentication

## Task Assignment

**Agent:** Backend Auth Logic  
**Worktree:** `/home/denny/Project/headscale-password-backend`  
**Branch:** `feature/password-backend`  
**Base Commit:** `b1495ddf5d26ec4089cd6459a491e33cffcf60ed`

---

## Objective

Implement password-based authentication for Headplane web UI that:
- Does NOT require database schema changes
- Stores password in configuration file or environment variable
- Provides an alternative to API key authentication
- Is scoped specifically to Headplane (does not grant general API access)

---

## Roadmap Reference

See `ROADMAP.md` Phase 13 for complete requirements.

---

## Technical Requirements

### 1. Password Storage
- Store password in `config.yaml` as `headplane.password` field
- Support environment variable `HEADSCALE_HEADPLANE_PASSWORD`
- Default to a secure randomly-generated password if not configured
- Log a warning if using the default password

### 2. Authentication Endpoint
- Implement `POST /api/v1/headplane/login` or similar
- Accept `{ "password": "..." }` in request body
- Return session token/JWT on success
- Return 401 on invalid password
- Rate limit to prevent brute-force (e.g., 5 attempts per minute per IP)

### 3. Authentication Middleware
- Create middleware to validate Headplane password authentication
- Distinguish between API key auth and password auth
- Password auth only grants Headplane web UI access
- Do NOT grant broader API access

### 4. Configuration Schema
- Update config schema to include `headplane.password` field
- Provide clear documentation in config-example.yaml
- Validate password is set in production deployments

### 5. Security
- Use bcrypt or similar for password hashing if stored
- Transmit over HTTPS only
- Implement rate limiting
- Add security warnings to logs

---

## Expected Files to Change

- `hscontrol/auth.go` or `hscontrol/headplane_auth.go` (new file)
- `hscontrol/types/config.go` (add Headplane config struct)
- `config-example.yaml` (add headplane.password example)
- `hscontrol/app.go` (register new routes)
- `hscontrol/*_test.go` (unit and integration tests)

---

## Testing Requirements

1. Unit tests for password validation
2. Integration tests for login endpoint
3. Test rate limiting
4. Test invalid password handling
5. Test configuration loading
6. Test environment variable override

---

## Acceptance Criteria

- [ ] Password can be configured via config.yaml or environment variable
- [ ] Login endpoint works and returns valid session token
- [ ] Invalid passwords return 401
- [ ] Rate limiting prevents brute-force
- [ ] All tests pass
- [ ] No database migrations required
- [ ] API key authentication remains unchanged and functional
- [ ] Security warnings logged appropriately

---

## Implementation Steps

1. Read existing authentication code in `hscontrol/auth.go`, `hscontrol/oidc.go`, `hscontrol/api_key.go`
2. Review config structure in `hscontrol/types/config.go`
3. Design password auth endpoint and middleware
4. Implement configuration loading
5. Implement authentication endpoint
6. Add rate limiting
7. Write tests
8. Update config-example.yaml
9. Verify all tests pass
10. Commit to feature branch

---

## Development Workflow

```bash
# Navigate to worktree
cd /home/denny/Project/headscale-password-backend

# Verify branch
git branch --show-current  # should show: feature/password-backend

# Make changes, test, commit
make test
make lint
git add -A
git commit -m "auth: implement password-based authentication for Headplane"

# Push when ready
git push -u origin feature/password-backend
```

---

## Dependencies

- None (can be implemented independently)

---

## Notes

- Do NOT modify existing API key authentication
- Do NOT add database migrations
- Password auth is Headplane-specific, not for general API access
- Follow existing code patterns in `hscontrol/`
- Use existing logging, error handling, and testing patterns

# Initial Prompt: QR Backend Hardening

**Copy and paste this into a new Cline agent instance:**

---

I need you to work on hardening the backend QR code generation and registration logic for non-OIDC interactive registration in the Headscale project.

**Your working directory is:**
```
/home/denny/Project/headscale-project/headscale-qr-backend-harden
```

**Your branch is:** `feature/qr-backend-harden`

**Your task:**
Harden the QR code generation, expiry, and auth cache logic. Add comprehensive tests for edge cases including expired QR codes, reused/invalid auth IDs, malformed payloads, and cache expiry enforcement.

**Instructions:**
Read the complete agent instructions at:
```
/home/denny/Project/headscale-project/headscale/AGENT_INSTRUCTIONS_QR_BACKEND_HARDEN.md
```

**Start by:**
1. Changing your working directory to `/home/denny/Project/headscale-project/headscale-qr-backend-harden`
2. Reading the agent instructions file
3. Reading the current implementation:
   - `hscontrol/qr/qr.go`
   - `hscontrol/state/auth_cache.go`
   - `hscontrol/handlers.go` (register endpoints)
4. Identifying edge cases not covered by tests
5. Adding tests and implementing hardening

**Important project rules:**
- Read `AGENTS.md` for project interaction rules
- Follow existing code style and patterns
- Add tests before implementing changes (TDD)
- Do not modify database schema
- Preserve existing API contract
- Commit incrementally with clear messages

**When complete:**
Report which files were changed, which tests were added, and the results of running:
- `go test ./hscontrol/qr -v -count=1`
- `go test ./hscontrol/state -run AuthCache -v -count=1`
- `go test ./hscontrol -run Register -v -count=1`
- `make lint`
- `go build ./cmd/headscale`

Begin by reading the agent instructions and confirming your understanding of the task.

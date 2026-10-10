# Initial Prompt: Headscale Frontend Template Hardening

**Copy and paste this into a new Cline agent instance:**

---

I need you to work on hardening the Headscale built-in registration page templates for non-OIDC interactive registration.

**IMPORTANT:** This is about Headscale's own Go HTML templates (`/register/{auth_id}` pages), NOT Headplane (the React UI). Headplane is already complete.

**Your working directory is:**
```
/home/denny/Project/headscale-project/headscale-headscale-frontend-harden
```

**Your branch is:** `feature/headscale-frontend-harden`

**Your task:**
Harden the registration page templates in `hscontrol/templates/`. Improve error handling, QR code display, and edge-case handling in the Go HTML templates served by Headscale.

**Instructions:**
Read the complete agent instructions at:
```
/home/denny/Project/headscale-project/headscale/AGENT_INSTRUCTIONS_HEADSCALE_FRONTEND_HARDEN.md
```

**Start by:**
1. Changing your working directory to `/home/denny/Project/headscale-project/headscale-headscale-frontend-harden`
2. Reading the agent instructions file
3. Reading the current templates:
   - `hscontrol/templates/auth_web.go` (non-OIDC registration page)
   - `hscontrol/templates/register_confirm.go` (OIDC registration page)
   - `hscontrol/handlers.go` (template rendering logic)
4. Identifying edge cases not covered by tests
5. Adding tests and improving template robustness

**Important project rules:**
- Read `AGENTS.md` for project interaction rules
- Follow existing template patterns
- Add tests before implementing changes (TDD)
- Do not modify backend logic
- This is NOT about Headplane (React UI)
- Commit incrementally with clear messages

**When complete:**
Report which files were changed, which tests were added, and the results of running:
- `go test ./hscontrol -run Template -v -count=1`
- `go test ./hscontrol -run Register -v -count=1`
- `make lint`
- `go build ./cmd/headscale`

Begin by reading the agent instructions and confirming your understanding of the task.

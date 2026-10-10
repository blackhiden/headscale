# Initial Prompt: QR Integration/E2E Testing

**Copy and paste this into a new Cline agent instance:**

---

I need you to add comprehensive integration and end-to-end tests for the non-OIDC interactive registration flow with QR code support in the Headscale project.

**Your working directory is:**
```
/home/denny/Project/headscale-project/headscale-qr-e2e-tests
```

**Your branch is:** `feature/qr-e2e-tests`

**Your task:**
Add integration/E2E tests covering the full registration flow, including CLI + QR, expiry enforcement, single-use validation, admin approval, and error scenarios.

**Instructions:**
Read the complete agent instructions at:
```
/home/denny/Project/headscale-project/headscale/AGENT_INSTRUCTIONS_QR_E2E_TESTS.md
```

**Start by:**
1. Changing your working directory to `/home/denny/Project/headscale-project/headscale-qr-e2e-tests`
2. Reading the agent instructions file
3. Reading the integration test documentation:
   - `cmd/hi/README.md`
   - `integration/README.md`
4. Reviewing existing registration tests
5. Planning new test scenarios

**Important project rules:**
- Read `AGENTS.md` for project interaction rules
- Read `cmd/hi/README.md` and `integration/README.md` BEFORE running tests
- Use `IntegrationSkip(t)` for integration tests
- Use `EventuallyWithT` for external calls
- Do not guess at `hi` flags
- Clean up test artifacts

**Test scenarios to cover:**
- Happy path: node register → QR display → scan → approval
- Expiry: QR expired before scan
- Single-use: auth ID consumed
- Invalid: malformed QR payload

**When complete:**
Report which tests were added and the results of running:
- `go run ./cmd/hi doctor`
- `go test ./integration -v -run QR`

**Note:** This task depends on backend/frontend hardening. If those tasks are not complete, focus on planning tests and reviewing existing patterns.

Begin by reading the agent instructions and confirming your understanding of the task.

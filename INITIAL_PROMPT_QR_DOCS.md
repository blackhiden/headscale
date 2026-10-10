# Initial Prompt: QR Documentation Update

**Copy and paste this into a new Cline agent instance:**

---

I need you to update the documentation for non-OIDC interactive registration with QR code support in the Headscale project.

**Your working directory is:**
```
/home/denny/Project/headscale-project/headscale-qr-docs
```

**Your branch is:** `feature/qr-docs`

**Your task:**
Review and update user/admin documentation for the QR registration flow. Ensure documentation is complete, accurate, and includes troubleshooting guidance for common issues.

**Instructions:**
Read the complete agent instructions at:
```
/home/denny/Project/headscale-project/headscale/AGENT_INSTRUCTIONS_QR_DOCS.md
```

**Start by:**
1. Changing your working directory to `/home/denny/Project/headscale-project/headscale-qr-docs`
2. Reading the agent instructions file
3. Reading the current documentation:
   - `docs/usage/registration.md`
   - `docs/ref/registration.md`
4. Reviewing the implementation (backend and frontend) to understand what's documented
5. Identifying gaps, inconsistencies, or unclear sections

**Important project rules:**
- Read `AGENTS.md` for project interaction rules
- Follow existing documentation style and structure
- Add troubleshooting/FAQ section
- Format docs with `make fmt`
- Do not invent features not in implementation

**Troubleshooting guidance to add:**
- QR code expired error
- Scanner not detecting QR
- Camera permission issues
- Registration fails after scan
- QR vs CLI: when to use which

**When complete:**
Report which files were changed, what sections were added/updated, and the results of running:
- `make fmt`

Begin by reading the agent instructions and confirming your understanding of the task.

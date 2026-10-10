# Agent Instructions: Automation Support (Task 4 of 4)

**Branch:** `feature/automation-support`  
**Worktree:** `/home/denny/Project/headscale-automation-support`  
**Base Commit:** `08c261007d293f2902339d584bfb5c3df629809d`

## Objective
Document and validate v2 API OAuth flow for Terraform provider and Kubernetes operator support.

## Expected Files
- `docs/ref/integration/tools.md` (update) — Terraform/K8s guide
- `docs/ref/api.md` (update) — OAuth flow documentation
- `docs/examples/terraform/` (new) — Terraform examples
- `docs/examples/kubernetes/` (new) — K8s operator examples
- Validation tests in `hscontrol/` (optional)

## Implementation Steps
1. Navigate: `cd /home/denny/Project/headscale-automation-support`
2. Review existing v2 API OAuth implementation (`hscontrol/api/v2/`)
3. Write comprehensive OAuth flow documentation
4. Create Terraform provider example configurations
5. Create Kubernetes operator example manifests
6. Validate examples against running Headscale instance
7. Add validation tests (optional, if gaps found)
8. Test: Review docs for accuracy, test examples
9. Commit: `git add -A && git commit -m "docs: add Terraform/K8s automation support" && git push -u origin feature/automation-support`

## Constraints
- This is primarily documentation work
- Validate examples actually work
- Follow existing doc patterns
- No unrelated refactoring

## Acceptance Criteria
- OAuth flow fully documented
- Terraform examples provided and tested
- K8s operator examples provided and tested
- Docs cover end-to-end automation setup
- Examples validated against real instance

See `docs/parallel-development.md` for workflow details.

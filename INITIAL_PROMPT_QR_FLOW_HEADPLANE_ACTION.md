# Initial Prompt: QR Flow Headplane Action Verification

You are assigned to verify Headplane's QR registration scanner action at the
exact submodule revision recorded by the approved Headscale baseline.

1. Change to
   `/home/denny/Project/headscale-project/headscale-qr-flow-headplane-action`.
2. Verify branch `feature/qr-flow-headplane-action`, superproject baseline
   `49ebce52d3608803573cb2a07ba21c461e4c18da`, and `headplane` commit
   `b470d2f87e52f3ab654897fff00732bd021a8495`.
3. Read `AGENTS.md`, `headplane/AGENTS.md`, and then
   `AGENT_INSTRUCTIONS_QR_FLOW_HEADPLANE_ACTION.md` in full.
4. Implement only the assigned focused regression coverage or a defect fix
   proven by that coverage.
5. Run the validation commands, commit scoped product changes in the submodule
   and then its gitlink only if needed, and report exact results.

Do not change Headscale Go code, documentation, API/schema/dependency contracts,
other worktrees, or the original checkout. Do not merge, rebase, force-push, or
rewrite history.
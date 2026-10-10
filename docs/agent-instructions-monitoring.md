# Agent Instructions: Monitoring Integrations (Task 3 of 4)

**Branch:** `feature/monitoring-integrations`  
**Worktree:** `/home/denny/Project/headscale-monitoring-integrations`  
**Base Commit:** `08c261007d293f2902339d584bfb5c3df629809d`

## Objective
Add support for Prometheus/Grafana webhooks and improve observability integrations.

## Expected Files
- `hscontrol/webhooks/` (new) — webhook implementation
- `hscontrol/api/v1/webhooks.go` (new) — webhook configuration API
- `config-example.yaml` (update) — webhook config examples
- `docs/usage/observability.md` (update) — webhook setup guide
- `docs/ref/integration/tools.md` (update) — alerting integrations
- Tests in `hscontrol/webhooks/` and `hscontrol/`

## Implementation Steps
1. Navigate: `cd /home/denny/Project/headscale-monitoring-integrations`
2. Design webhook configuration schema
3. Implement webhook dispatcher (events → webhook calls)
4. Add webhook configuration API endpoints
5. Document webhook events and payloads
6. Add Grafana/Prometheus integration examples
7. Add tests (unit + integration)
8. Test: `make test && go test ./...`
9. Commit: `git add -A && git commit -m "feat: add monitoring/alerting webhooks" && git push -u origin feature/monitoring-integrations`

## Constraints
- Work only in this worktree
- Don't break existing metrics
- Follow existing patterns
- No unrelated refactoring

## Acceptance Criteria
- Webhooks configurable via config file
- Webhook events documented
- Grafana/Prometheus examples provided
- Docs cover integration setup
- All tests pass

See `docs/parallel-development.md` for workflow details.

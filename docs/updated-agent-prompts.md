# Updated Agent Prompts for Phase 11

Due to disk space constraints, agents need to install frontend dependencies themselves.

## Task 1: Plugin System (UPDATED)

Context:
- Working directory: /home/denny/Project/headscale-plugin-system
- Branch: feature/plugin-system

IMPORTANT - Frontend Setup Required:
1. cd /home/denny/Project/headscale-plugin-system/headplane
2. Run: pnpm install
3. Verify: pnpm typecheck

Task: Implement plugin/extension system for Headplane.

Read: docs/agent-instructions-plugin-system.md

---

## Task 2: Multi-Instance Dashboard (UPDATED)

Context:
- Working directory: /home/denny/Project/headscale-multi-instance-dashboard
- Branch: feature/multi-instance-dashboard

IMPORTANT - Frontend Setup Required:
1. cd /home/denny/Project/headscale-multi-instance-dashboard/headplane
2. Run: pnpm install
3. Verify: pnpm typecheck

Task: Add multi-instance dashboard support.

Read: docs/agent-instructions-multi-instance.md

---

## Task 3: Monitoring Integrations (NO CHANGES)

Context:
- Working directory: /home/denny/Project/headscale-monitoring-integrations
- Branch: feature/monitoring-integrations

Task: Add Prometheus/Grafana webhooks (backend only).

Read: docs/agent-instructions-monitoring.md

---

## Task 4: Automation Support (NO CHANGES)

Context:
- Working directory: /home/denny/Project/headscale-automation-support
- Branch: feature/automation-support

Task: Document Terraform/K8s automation (docs only).

Read: docs/agent-instructions-automation.md

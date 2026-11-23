# Merge Preservation Contract — Doni Arsy Fork Features

> **ACTIVE MERGE CONTRACT — DO NOT DELETE OR REPLACE THIS FILE DURING AN
> UPSTREAM MERGE.** This file is the authority for what must survive a merge
> from `upstream/main`. If a conflict touches this file, keep the fork version
> and append only new information.

## 0. Read this before merging

1. Never accept an upstream deletion of a `FP-*` path or symbol without an
   explicit retirement entry in the **Merge Decision Log** at the end.
2. For fork-owned files, keep the fork version. For shared files, merge by hand
   and preserve the fork block described in the feature card.
3. Run the **One-Command Preservation Audit** before committing the merge.
4. Run `git -C headplane pnpm run check:fork` for submodule features.
5. Record the upstream SHA, merge commit, audit result, and any restores in the
   **Merge Decision Log**.
6. Update `CHANGELOG.md` and `session.md` in the same merge commit.

If any audit check fails, stop. Do not push the merge.

## 1. Feature manifest

Stable IDs are used in conflict notes, commit messages, and the decision log.

| ID | Feature | Author | Status | Critical paths / tokens |
| --- | --- | --- | --- | --- |
| `FP-001` | Fork identity, module path, build metadata | Doni Arsy | Active | `go.mod`, `Dockerfile`, `Makefile`, `APP_VERSION`, `__COMMIT_HASH__` |
| `FP-002` | Update checker | Doni Arsy | Active / DO NOT REMOVE | `hscontrol/updatecheck/`, `hscontrol/api/v1/updatecheck.go`, `HEADSCALE_UPDATE_CHECK_REPO` |
| `FP-003` | DERP status endpoint | Doni Arsy | Active | `hscontrol/api/v1/derp.go`, `GET /api/v1/derp` |
| `FP-004` | Webhooks | Doni Arsy | Active | `hscontrol/api/v1/webhooks.go`, `hscontrol/webhook.go`, `202609301000-create-webhooks-table` |
| `FP-005` | Machine approval API | Doni Arsy | Active | `hscontrol/api/v1/machines.go`, `machine_approve`, `devices:core` |
| `FP-006` | Headplane local admin, settings, admin UI | Doni Arsy (gitlinks) | Active | `headplane`, `hscontrol/db/migrate_headplane.go`, `202610050900-repair-headplane-schema` |
| `FP-007` | OIDC groups in ACLs | Doni Arsy | Active | `OIDCGroups`, `oidc_groups`, `Groups.ContainsOrOIDC` |
| `FP-008` | Phase 10/11 expansion | Doni Arsy | Active | `/routes`, `/instances`, `headplane/app/plugins/`, `docs/automation.md` |
| `FP-009` | QR registration | Doni Arsy | Active | `hscontrol/qr/`, `headscale-registration`, `/machines/scan-qr` |
| `FP-010` | Server info / admin technical info | Denny Lastpaste (adjacent) | Active | `hscontrol/api/v1/serverinfo.go`, `/api/admin/server-info` |
| `FP-011` | Observability and readiness | Doni Arsy | Active | `GET /ready`, `docs/usage/observability.md` |
| `FP-012` | UI/UX and accessibility | Doni Arsy (gitlinks) | Active | `EmptyState`, `LoadingSpinner`, `docs/usage/ui-guide.md` |
| `FP-013` | Documentation and development process | Doni Arsy | Active | `ARCHITECTURE.md`, `CODE_GRAPH.md`, `SPECIFICATION.md`, `ROADMAP.md`, `SCHEMA.md`, this file |
| `FP-014` | Small load-bearing fixes | Doni Arsy | Active | `Credential.Clone`, `testify/assert`, `capver`, `zerolog` |
| `FP-015` | Deliberate removals | Doni Arsy | Negative list | `*.md` ignore, removed agent docs — do **not** restore |
| `FP-016` | Headplane gitlink selection | Doni Arsy | Active | `headplane` gitlink, `pnpm run check:fork` |
| `FP-017` | Legacy fork config key aliases | Denny Lastpaste | Active | `hscontrol/config/keys.go`, `applyLegacyAliases`, `derp.server.auto_update_enabled` |

## 2. Scope and attribution

The scan was performed on 2026-11-22 against all `origin/*` refs. At scan time
`origin/main` was `9ef033c7`; the local working tree was at `3c788b7a`.

- Doni-authored commits reachable from `origin`: **297**
- Non-merge commits: **172**
- Merge commits: **124**

Pure merge commits, Cline checkpoint/index commits, and release/changelog-only
commits are not standalone feature code. They may contain conflict resolutions,
so review them when a feature disappears after a merge.

Attribution is by commit author `Doni Arsy <arsydoni4326@gmail.com>`.
`FP-010` and the `hscontrol/qr_api.go` part of `FP-009` were authored by
**Denny Lastpaste**, not Doni; they are included because they are load-bearing
fork features and easy to lose in an upstream merge.

The `headplane/` submodule contains **zero** Doni-authored commits in this
clone. Doni's parent-repository contribution there is gitlink selection and
integration (`FP-016`). Preserve the intended gitlink and the features it
represents.

Regenerate the raw inventory with:

```bash
git log --remotes=origin --author='Doni Arsy' --no-merges --oneline
git log --remotes=origin --author='Doni Arsy' --merges --oneline
git diff --name-status upstream/main...origin/main
```

## 3. Feature preservation cards

Each card is a self-contained merge instruction. Copy the card ID into conflict
notes and the decision log.

### FP-001 — Fork identity, module path, and build metadata

- **Author:** Doni Arsy
- **Last known good:** `7337a8a9`, `573ac29b`, `a6c6563d`, `9d4a0bba`,
  `5b7f06ce`, `190cdd75`, `dcc7d641`
- **Protected paths:** `go.mod`, Go imports across `cmd/`, `hscontrol/`,
  `integration/`, `gen/`, `Dockerfile`, `Makefile`,
  `.github/workflows/deploy.yaml`, `.goreleaser.yml`, `nix/`, docs, and
  `headplane/vite.config.ts`.
- **Required tokens:** `module github.com/blackhiden/headscale`, `APP_VERSION`,
  `APP_COMMIT`, `BUILD_DATE`, `__COMMIT_HASH__`, `__BUILD_TIME__`.
- **Contract:** Fork version suffix is documented as `-arsydoni4326-alt`.
  Build metadata is injected into `hscontrol/types.{Version,Commit,BuildDate}`.
  Some upstream URLs remain intentionally in user docs; review before blanket
  replacement.
- **Conflict rule:** Keep the fork module path in Go code and configuration.
  Do not accept an upstream reversion to `github.com/juanfont/headscale` in
  imports or build flags.
- **Restore:** `git checkout <last-good> -- go.mod Dockerfile Makefile .github/workflows/deploy.yaml .goreleaser.yml`
- **Verify:** `go build ./cmd/headscale` and
  `git grep -n 'github.com/juanfont/headscale' -- '*.go' '*.yaml' '*.yml'`
  (review remaining hits).

### FP-002 — Update checker

- **Author:** Doni Arsy
- **Last known good:** `99245073`, `af2d315d`, `1df134a8`, `e985a037`,
  `08c26100`, `3ba121c7`
- **Protected paths:** `hscontrol/updatecheck/`,
  `hscontrol/api/v1/updatecheck.go`, `hscontrol/apiv1_updatecheck_test.go`,
  `headplane/app/update-check/`, `headplane/app/routes/util/update-check.ts`.
- **Required tokens:** `DO NOT REMOVE`, `UpdateCheckResult`,
  `HEADSCALE_UPDATE_CHECK_REPO`, `DisableUpdateCheck`,
  `headscale_updatecheck_requests_total`.
- **Contract:** `GET /api/v1/update-check` is unauthenticated. Without
  `?check=true` it returns current version info; with it, fetch remote
  release/commit, compare, and cache for 15 minutes. Metrics must remain
  registered.
- **Conflict rule:** If upstream deletes `hscontrol/updatecheck/`, restore the
  whole directory. If upstream rewrites `updatecheck.go`, keep
  `registerUpdateCheck` with no `Security` requirement and the
  `BuildResponse(check, disabled)` call.
- **Restore:** `git checkout <last-good> -- hscontrol/updatecheck hscontrol/api/v1/updatecheck.go headplane/app/update-check`
- **Verify:** `go test ./hscontrol/updatecheck/...` and
  `go test ./hscontrol -run TestAPIV1UpdateCheck -count=1`.

### FP-003 — DERP status endpoint

- **Author:** Doni Arsy
- **Last known good:** `b0764611`
- **Protected paths:** `hscontrol/api/v1/derp.go`,
  `hscontrol/apiv1_derp_test.go`, `hscontrol/servertest/derp_test.go`,
  `headplane/app/routes/derp/`.
- **Required tokens:** `getDerp`, `/api/v1/derp`, `DERPResponseBody`,
  `totalRegions`.
- **Contract:** Authenticated read-only endpoint returning `configured`,
  `totalRegions`, and sorted `regions` with DERP/STUN ports and IP addresses.
- **Conflict rule:** Never accept deletion of `derp.go`. If upstream adds a
  similar endpoint, merge it but keep this response shape and OpenAPI
  operation ID.
- **Restore:** `git checkout b0764611 -- hscontrol/api/v1/derp.go hscontrol/apiv1_derp_test.go`
- **Verify:** `go test ./hscontrol -run TestAPIV1Derp -count=1`.

### FP-004 — Webhooks

- **Author:** Doni Arsy
- **Last known good:** `d3b29def`, `2421cbdb`
- **Protected paths:** `hscontrol/api/v1/webhooks.go`, `hscontrol/webhook.go`,
  `hscontrol/webhook_test.go`, `hscontrol/types/webhook.go`,
  `hscontrol/db/webhook.go`, `hscontrol/db/migrate_webhooks.go`,
  `hscontrol/db/schema.sql`, `docs/usage/webhooks.md`.
- **Required tokens:** `202609301000-create-webhooks-table`,
  `X-Headscale-Signature`, `node_up`, `alert_firing`,
  `headscale_webhook_dispatch_total`.
- **Contract:** CRUD at `/api/v1/webhooks`; async JSON POST dispatch; optional
  HMAC-SHA256; explicit per-dialect DDL. Do not replace the migration DDL with
  `AutoMigrate`.
- **Conflict rule:** Keep the migration ID and its position in the migration
  array. Never reorder migrations.
- **Restore:** `git checkout <last-good> -- hscontrol/api/v1/webhooks.go hscontrol/webhook.go hscontrol/db/webhook.go hscontrol/db/migrate_webhooks.go hscontrol/db/schema.sql`
- **Verify:** `go test ./hscontrol -run TestWebhook -count=1` and
  `go test ./hscontrol/db -count=1`.

### FP-005 — Machine approval API

- **Author:** Doni Arsy
- **Last known good:** `2421cbdb`
- **Protected paths:** `hscontrol/api/v1/machines.go`,
  `hscontrol/apiv1_machines_test.go`, `gen/client/v1/client.gen.go`,
  `docs/ref/api.md`, `docs/ref/registration.md`.
- **Required tokens:** `approveMachine`, `approveMachines`, `devices:core`,
  `machine_approve_bulk`, `nodeIds`.
- **Contract:** Single and bulk approval endpoints. Approval clears expiry.
  Bulk is non-atomic with per-node results. Admin API key or `devices:core`
  OAuth scope required; read-only tokens get `403`.
- **Conflict rule:** If upstream adds a different approval API, keep both or
  explicitly retire this one in the decision log. Do not drop the generated
  client operations.
- **Restore:** `git checkout 2421cbdb -- hscontrol/api/v1/machines.go hscontrol/apiv1_machines_test.go gen/client/v1/client.gen.go`
- **Verify:** `go test ./hscontrol -run TestAPIV1Machines -count=1`.

### FP-006 — Headplane local admin, settings, and admin UI

- **Author:** Doni Arsy (parent-repo gitlink integration)
- **Last known good:** `8e2131bf`, `0370cd5d`, `04f5d145`, `a64257cc`,
  `e82f2fbd`, `1e3a0b50`, `9fbba83b`, `6975b340`, `a13af9cb`, `01c8474c`,
  `254587a4`, `28222aa6`, `c5076575`, `c09ade7d`
- **Protected paths:** `headplane` gitlink,
  `hscontrol/db/migrate_headplane.go`, `hscontrol/db/schema.sql`,
  `hscontrol/db/db.go`, `docs/usage/authentication.md`,
  `docs/usage/settings.md`, `docs/ref/api/headplane-settings.md`,
  `docs/ref/api/headplane-users.md`.
- **Required tokens:** `202610050900-repair-headplane-schema`,
  `202610050930-normalize-headplane-indexes`, `headplane_users`,
  `headplane_settings`, `localAdminContext`, `createPasswordSession`,
  `check:fork`.
- **Contract:** One local administrator in Headplane config with bcrypt
  password; API-key login is secondary. Password sessions, rate limiting,
  atomic config writes, session invalidation, optional `user.name`/`user.avatar`.
  API-key operations are proxied through the Headplane backend. Legacy tables
  are rollback data only.
- **Conflict rule:** Do **not** restore retired Go files
  `hscontrol/headplane_auth.go`, `hscontrol/headplane_settings.go`, or
  `hscontrol/headplane_users.go`. Preserve the Headplane implementation and the
  schema repair migrations.
- **Restore:** `git checkout <last-good> -- hscontrol/db/migrate_headplane.go hscontrol/db/schema.sql hscontrol/db/db.go` then restore the intended `headplane` gitlink (see `FP-016`).
- **Verify:** `git -C headplane pnpm run check:fork`,
  `git -C headplane pnpm run test:unit`, and
  `go test ./hscontrol/db -count=1`.

### FP-007 — OIDC groups in ACLs

- **Author:** Doni Arsy
- **Last known good:** `08c26100`
- **Protected paths:** `hscontrol/types/users.go`,
  `hscontrol/types/types_view.go`, `hscontrol/types/types_clone.go`,
  `hscontrol/policy/v2/types.go`, `hscontrol/policy/v2/policy.go`,
  `hscontrol/api/v1/types.go`, `hscontrol/api/v2/users.go`,
  `gen/client/v1/client.gen.go`, `gen/client/v2/client.gen.go`,
  `hscontrol/db/schema.sql`, `hscontrol/db/db.go`,
  `hscontrol/policy/v2/oidc_groups_test.go`, `docs/ref/policy.md`,
  `docs/ref/oidc.md`.
- **Required tokens:** `OIDCGroups`, `oidc_groups`,
  `202609291402-add-oidc-groups-to-users`, `Groups.ContainsOrOIDC`,
  `oidcGroups`.
- **Contract:** OIDC `groups` claim is stored on users and refreshed at login.
  Groups resolve in ACLs, grants, `tagOwners`, and `autoApprovers`.
  Policy-defined groups take precedence.
- **Conflict rule:** Keep the migration ID and column. If upstream adds a
  different group model, merge it but keep `oidcGroups` in v1/v2 responses.
- **Restore:** `git checkout 08c26100 -- hscontrol/types/users.go hscontrol/policy/v2/types.go hscontrol/api/v1/types.go hscontrol/api/v2/users.go hscontrol/db/schema.sql`
- **Verify:** `go test ./hscontrol/policy/v2 -run OIDC -count=1` and
  `go test ./hscontrol -run 'TestAPIV1.*User' -count=1`.

### FP-008 — Phase 10/11 feature expansion

- **Author:** Doni Arsy
- **Last known good:** `08c26100`, `3edf7ee8`, `b3dc40c0`, `3f0df94e`,
  `d3b29def`, `8fafea34`
- **Protected paths:** `headplane/app/routes/routes/`,
  `headplane/app/plugins/`, `headplane/app/routes/instances/`,
  `docs/ref/extending.md`, `docs/automation.md`,
  `docs/ref/integration/tools.md`, `docs/usage/observability.md`.
- **Required tokens:** `machine.rotate_key`, `/routes`, `/instances`,
  `INSTANCE_ENCRYPTION_SECRET`, `docs/ref/extending.md`.
- **Contract:** Key rotation, route overview, DNS `extra_records_path` notice,
  CNAME/IP validation, plugin system, multi-instance dashboard with encrypted
  credentials, and automation docs.
- **Conflict rule:** These are mostly submodule features; preserve the gitlink
  and the docs. Do not let upstream docs replace `docs/ref/extending.md` with a
  version that drops fork extension points.
- **Restore:** `git checkout <last-good> -- docs/ref/extending.md docs/automation.md docs/ref/integration/tools.md` and restore the gitlink.
- **Verify:** `git -C headplane pnpm run test:unit` and inspect `/routes` and
  `/instances` routes.

### FP-009 — QR registration

- **Author:** Doni Arsy; adjacent `qr_api.go` by Denny Lastpaste
- **Last known good:** `f5521675`, `b55c6e18`, `3190ef73`, `a3911b1e`,
  `6cfae2e3`, `70fb8f45`, `681d5675`, `22d85755`, `9ca1b306`, `94aac78d`,
  `4cf43205`, `67cb981b`, `77fc2040`, `526f3a4f`, `c5bfb7ba`, `6312bb7f`
- **Protected paths:** `hscontrol/qr/`, `hscontrol/qr_api.go`,
  `hscontrol/handlers.go`, `hscontrol/templates/register_confirm.go`,
  `hscontrol/state/auth_cache_test.go`, `headplane/app/routes/machines/scan-qr.tsx`,
  `headplane/app/components/qr-scanner/`, `docs/api/qr-payload.md`,
  `docs/usage/registration.md`, `docs/ref/registration.md`.
- **Required tokens:** `headscale-registration`, `NewRegistrationPayload`,
  `GenerateRegistrationQR`, `/api/v1/registration/qr/{auth_id}`,
  `/machines/scan-qr`, `html5-qrcode`.
- **Contract:** QR payload contains `type`, `version`, `auth_id`, `server_url`,
  `expires_at`. Expiry is bound to the pending registration cache deadline.
  Headscale remains authoritative for single-use auth IDs. The open JSON API
  returns 200/400/404/410.
- **Conflict rule:** Never let upstream remove `hscontrol/qr/` or the QR section
  of the registration templates. Keep the expiry binding; do not refresh the
  deadline on page reload.
- **Restore:** `git checkout <last-good> -- hscontrol/qr hscontrol/qr_api.go hscontrol/templates/register_confirm.go docs/api/qr-payload.md` and restore the gitlink.
- **Verify:** `go test ./hscontrol/qr -count=1`,
  `go test ./hscontrol/state -run 'Test(AuthCache|SetAuthCacheEntry)' -count=1`,
  `git -C headplane pnpm run test:component`.

### FP-010 — Server info / admin technical info

- **Author:** Denny Lastpaste (adjacent, not Doni)
- **Last known good:** current `origin/main`
- **Protected paths:** `hscontrol/api/v1/serverinfo.go`,
  `headplane/app/routes/api/admin/server-info.ts`,
  `headplane/app/routes/admin/components/technical-info-section.tsx`,
  `docs/api/server-info.md`, `docs/usage/admin-technical-info.md`.
- **Required tokens:** `getServerInfo`, `/api/v1/server-info`,
  `/api/admin/server-info`, `ServerInfoResponseBody`.
- **Contract:** Authenticated, read-only, non-sensitive system info. Headplane
  merges Headscale and Headplane info and degrades gracefully.
- **Conflict rule:** Keep both endpoints and the admin-only access. Do not
  expose secrets or detailed hardware identifiers.
- **Restore:** `git checkout <last-good> -- hscontrol/api/v1/serverinfo.go docs/api/server-info.md docs/usage/admin-technical-info.md` and restore the gitlink.
- **Verify:** `go test ./hscontrol -run TestAPIV1ServerInfo -count=1` and
  `git -C headplane pnpm run test:unit`.

### FP-011 — Observability and readiness

- **Author:** Doni Arsy
- **Last known good:** `bf271929`
- **Protected paths:** `hscontrol/handlers.go`, `hscontrol/app.go`,
  `hscontrol/handlers_test.go`, `docs/usage/observability.md`.
- **Required tokens:** `ReadyHandler`, `/ready`, `TestReadyHandler`,
  `headscale_updatecheck_requests_total`.
- **Contract:** `/ready` returns 200 when the database is reachable and 503
  otherwise. Observability docs catalog endpoints, metrics, logging, debug
  endpoints, alerts, and the Grafana dashboard.
- **Conflict rule:** Do not drop `/ready` or the fork-specific metrics section.
- **Restore:** `git checkout bf271929 -- hscontrol/handlers.go hscontrol/app.go hscontrol/handlers_test.go docs/usage/observability.md`
- **Verify:** `go test ./hscontrol -run 'Test(Health|Ready|Version)Handler' -count=1`.

### FP-012 — UI/UX and accessibility

- **Author:** Doni Arsy (gitlink integration)
- **Last known good:** `250afc5c`, `9929dd07`, `2d3c18af`, `411bebf4`
- **Protected paths:** `headplane/app/components/empty-state.tsx`,
  `headplane/app/components/loading-spinner.tsx`,
  `headplane/app/components/error-banner.tsx`, `docs/usage/ui-guide.md`.
- **Required tokens:** `EmptyState`, `LoadingSpinner`, `ErrorBanner`,
  `test:a11y`.
- **Contract:** Standardized empty/loading/error states, accessibility docs,
  axe-core/e2e coverage, design tokens, responsive navigation, full-width
  admin cards, avatar display, and admin menu fixes.
- **Conflict rule:** Preserve the gitlink and `docs/usage/ui-guide.md`. Do not
  treat screenshot placeholders as missing files.
- **Restore:** `git checkout <last-good> -- docs/usage/ui-guide.md` and restore
  the gitlink.
- **Verify:** `git -C headplane pnpm run test:component` and
  `git -C headplane pnpm run test:a11y`.

### FP-013 — Documentation and development process

- **Author:** Doni Arsy
- **Last known good:** `af2d315d`, `a06d0a9e`, `1df134a8`, `68920d38`,
  `3ba121c7`, `08c26100`, `23a88d24`, `3edf7ee8`
- **Protected paths:** `ARCHITECTURE.md`, `CODE_GRAPH.md`, `SPECIFICATION.md`,
  `ROADMAP.md`, `CHANGELOG.md`, `CONTRIBUTING.md`, `session.md`, `SCHEMA.md`,
  `docs/MERGE-PRESERVATION.md`, `docs/about/fork.md`, `docs/about/testing.md`,
  `docs/troubleshooting.md`, `docs/migration-v0.36.3-to-v1.0.0.md`,
  `docs/parallel-development.md`, `docs/parallel-development-workflow.md`,
  `scripts/parallel-feature-worktrees.sh`, `scripts/phase13c-setup.sh`,
  `.github/workflows/fork-features.yml`.
- **Required tokens:** `Feature Preservation`, `DO NOT REMOVE`, `check:fork`,
  `fork-features.yml`.
- **Contract:** Root and user docs describe fork features. The CI gate checks
  four paths only; this contract is the broader checklist.
- **Conflict rule:** If upstream deletes a required doc, restore it. If upstream
  edits it, merge the upstream text but keep the fork-specific sections.
- **Restore:** `git checkout <last-good> -- ARCHITECTURE.md CODE_GRAPH.md SPECIFICATION.md ROADMAP.md CONTRIBUTING.md SCHEMA.md docs/MERGE-PRESERVATION.md`
- **Verify:** `mkdocs build --strict` (when the Material theme is installed).

### FP-014 — Small but load-bearing fixes

- **Author:** Doni Arsy
- **Last known good:** `ab651692`, `5c9c91eb`, plus lint-fix commits
- **Protected paths:** `hscontrol/types/types_clone.go`,
  `hscontrol/servertest/sshcheck_test.go`, `hscontrol/noise_test.go`,
  `hscontrol/db/db.go`, `hscontrol/webhook_test.go`.
- **Required tokens:** `src.User.Clone()`, `testify/assert`, `capver`,
  `zerolog`, `202610031721-create-headplane-users` must appear **once**.
- **Contract:** `Credential.Clone()` deep-copies `User`; SSH check tests import
  `assert`; noise tests import `capver` and `zerolog`; duplicate migration
  entry removed; webhook test isolated; context-aware HTTP requests in tests.
- **Conflict rule:** These are single-line fixes that upstream conflict
  resolution often drops. Check them explicitly after every merge.
- **Restore:** `git checkout <last-good> -- hscontrol/types/types_clone.go hscontrol/servertest/sshcheck_test.go hscontrol/noise_test.go`
- **Verify:** `go test ./hscontrol/types ./hscontrol/servertest ./hscontrol -run 'Test(SSHCheck|Noise|Credential)' -count=1`.

### FP-015 — Deliberate removals (negative list)

- **Author:** Doni Arsy
- **Last known good:** `25350a6d`, `d290acea`
- **Protected paths:** `.gitignore` (`*.md`), removed root agent docs,
  `CODE_OF_CONDUCT.md`, `CLAUDE.md`, some upstream Dockerfiles/workflows.
- **Required tokens:** `*.md` in `.gitignore`.
- **Contract:** Generated agent-instruction and implementation-summary Markdown
  files were intentionally removed from tracking. New Markdown files must be
  added with `git add -f`. Required project docs remain tracked.
- **Conflict rule:** Do **not** restore removed files unless explicitly
  requested. Do **not** remove `*.md` from `.gitignore` without a decision.
- **Restore:** Not applicable. This is a do-not-restore list.
- **Verify:** `git check-ignore -v docs/MERGE-PRESERVATION.md` should show the
  ignore rule, while `git ls-files docs/MERGE-PRESERVATION.md` still lists it.

### FP-016 — Headplane gitlink selection

- **Author:** Doni Arsy (parent-repo gitlink updates)
- **Last known good:** intended gitlink recorded in the merge decision log
- **Protected paths:** `headplane` gitlink, `headplane/docs/MERGE-PRESERVATION.md`,
  `headplane/scripts/check-fork-features.sh`, `headplane/lefthook.yml`,
  `headplane/.github/workflows/build.yaml`.
- **Required tokens:** `check:fork`, `localAdminContext`, `qr-scanner`,
  `update-check`.
- **Contract:** The parent repo pins a Headplane commit that contains all
  submodule-side fork features. Upstream merges must not reset this gitlink to
  an upstream-only commit.
- **Conflict rule:** If the gitlink changes during a merge, compare the old and
  new submodule commits. Accept the new one only if `pnpm run check:fork`
  passes and all submodule features are present.
- **Restore:** `git checkout <last-good> -- headplane` then
  `git submodule update --init --recursive`.
- **Verify:** `git ls-tree HEAD headplane`, `git submodule status`, and
  `git -C headplane pnpm run check:fork`.

### FP-017 — Legacy fork config key aliases

- **Author:** Denny Lastpaste
- **Last known good:** `7aa662f2`
- **Protected paths:** `hscontrol/config/keys.go`, `hscontrol/config/load.go`,
  `hscontrol/config/load_test.go`, `CHANGELOG.md`, `docs/troubleshooting.md`.
- **Required tokens:** `legacyAliases`, `applyLegacyAliases`,
  `derp.server.auto_update_enabled`, `derp.server.update_frequency`,
  `node.routes.probe_interval`, `node.routes.probe_timeout`.
- **Contract:** Existing fork configs using the pre-refactor key locations must
  keep loading. The loader rewrites them to `derp.auto_update_enabled`,
  `derp.update_frequency`, `node.routes.ha.probe_interval`, and
  `node.routes.ha.probe_timeout`, emitting one deprecation warning per key. If
  both legacy and current keys are present, the current key wins.
- **Conflict rule:** Never accept an upstream config-loader change that removes
  `legacyAliases` or `applyLegacyAliases` without an explicit retirement entry
  in the decision log. Keep the four alias mappings and the regression test.
- **Restore:** `git checkout 7aa662f2 -- hscontrol/config/keys.go hscontrol/config/load.go hscontrol/config/load_test.go`
- **Verify:** `go test ./hscontrol/config -run '^TestConfigSource$/fork_legacy_key_aliases$' -count=1`
  and `headscale check config -f <config-with-legacy-keys>`.

## 4. One-Command Preservation Audit

Run this immediately after resolving merge conflicts and before committing the
merge. It checks paths, required tokens, migration IDs, and submodule
integration points. A failure means **stop and restore**.

```bash
FAILED=0
need_path() { [ -e "$1" ] || { echo "❌ missing path: $1"; FAILED=1; }; }
need_token() { grep -qF -- "$2" "$1" || { echo "❌ missing token '$2' in $1"; FAILED=1; }; }
need_count() {
  n=$(grep -cF -- "$2" "$1" || true)
  [ "$n" = "$3" ] || { echo "❌ expected $3 occurrence(s) of '$2' in $1, got $n"; FAILED=1; }
}

for p in \
  go.mod Dockerfile Makefile .github/workflows/deploy.yaml .goreleaser.yml \
  hscontrol/updatecheck hscontrol/api/v1/updatecheck.go \
  hscontrol/api/v1/derp.go hscontrol/api/v1/webhooks.go hscontrol/webhook.go \
  hscontrol/api/v1/machines.go hscontrol/qr hscontrol/qr_api.go \
  hscontrol/api/v1/serverinfo.go hscontrol/handlers.go \
  hscontrol/db/migrate_headplane.go hscontrol/db/migrate_webhooks.go \
  hscontrol/db/schema.sql hscontrol/types/types_clone.go \
  hscontrol/config/keys.go hscontrol/config/load.go \
  hscontrol/servertest/sshcheck_test.go hscontrol/noise_test.go \
  docs/MERGE-PRESERVATION.md SCHEMA.md docs/ref/extending.md \
  docs/usage/observability.md docs/usage/ui-guide.md \
  headplane/app/update-check headplane/app/components/qr-scanner \
  headplane/scripts/check-fork-features.sh \
  headplane/docs/MERGE-PRESERVATION.md; do
  need_path "$p"
done

need_token go.mod 'module github.com/blackhiden/headscale'
need_token Dockerfile 'APP_VERSION'
need_token headplane/vite.config.ts '__COMMIT_HASH__'
need_token hscontrol/updatecheck/updatecheck.go 'DO NOT REMOVE'
need_token hscontrol/api/v1/updatecheck.go 'update-check'
need_token hscontrol/api/v1/derp.go 'getDerp'
need_token hscontrol/webhook.go 'X-Headscale-Signature'
need_token hscontrol/api/v1/machines.go 'machine_approve_bulk'
need_token hscontrol/types/users.go 'OIDCGroups'
need_token hscontrol/policy/v2/types.go 'Groups.ContainsOrOIDC'
need_token hscontrol/qr/payload.go 'headscale-registration'
need_token hscontrol/api/v1/serverinfo.go 'getServerInfo'
need_token hscontrol/handlers.go 'ReadyHandler'
need_token hscontrol/types/types_clone.go 'src.User.Clone()'
need_token hscontrol/servertest/sshcheck_test.go 'testify/assert'
need_token hscontrol/noise_test.go 'capver'
need_token hscontrol/config/keys.go 'derp.server.auto_update_enabled'
need_token hscontrol/config/load.go 'applyLegacyAliases'
need_token docs/ref/extending.md 'Extension points'
need_token docs/usage/observability.md 'headscale_updatecheck_requests_total'
need_token docs/usage/ui-guide.md 'EmptyState'
need_token .gitignore '*.md'

need_count hscontrol/db/db.go '202609291402-add-oidc-groups-to-users' 1
need_count hscontrol/db/db.go '202609301000-create-webhooks-table' 1
need_count hscontrol/db/db.go '202610050900-repair-headplane-schema' 1
need_count hscontrol/db/db.go '202610050930-normalize-headplane-indexes' 1
need_count hscontrol/db/db.go '202610031721-create-headplane-users' 1

need_token headplane/app/routes/auth/login/action.ts 'localAdminContext'
need_token headplane/app/routes/auth/login/page.tsx 'name="password"'
need_token headplane/app/server/context.ts 'buildLocalAdmin'
need_token headplane/app/server/db/schema.ts "'password'"
need_token headplane/app/server/web/auth.ts 'createPasswordSession'
need_token headplane/app/routes/machines/scan-qr.tsx 'scan-qr'
need_token headplane/package.json 'check:fork'

echo "HEAD gitlink:  $(git ls-tree HEAD headplane)"
echo "Submodule:     $(git submodule status)"
echo "Upstream drift:"
git diff --name-status upstream/main...HEAD | head -80

if [ "$FAILED" -ne 0 ]; then
  echo "🚨 PRESERVATION AUDIT FAILED — do not commit the merge."
  exit 1
fi
echo "✅ Preservation audit passed."
```

Then run the feature-specific verification commands from the cards and:

```bash
git -C headplane pnpm run check:fork
go build ./cmd/headscale
go test ./hscontrol/...
```

## 5. Merge protocol

### Before

- [ ] Record `git rev-parse origin/main`, `git rev-parse HEAD`, and
      `git ls-tree HEAD headplane` in `session.md`.
- [ ] Create a backup branch: `git checkout -b backup/pre-merge-$(date +%Y%m%d)`.
- [ ] Run the audit in section 4 on the pre-merge tree; it must pass.
- [ ] Review `git diff --name-status upstream/main...origin/main`.
- [ ] List every upstream deletion or rename that touches a `FP-*` path.

### During

- [ ] For fork-owned files, keep the fork version.
- [ ] For shared files, merge manually and preserve the fork block named in the
      card's **Contract**.
- [ ] Never run `git checkout --theirs` on a protected path.
- [ ] Use `git checkout --ours` only for fork-owned files, never for shared
      files such as `hscontrol/app.go`, `hscontrol/db/db.go`, or
      `headplane/app/server/context.ts`.
- [ ] Restore deleted files with the **Restore** command in the relevant card.
- [ ] Keep migration order immutable; append new migrations only.

### After

- [ ] Run the section 4 audit. Any failure blocks the merge.
- [ ] Run `git -C headplane pnpm run check:fork`.
- [ ] Run the **Verify** command from every affected card.
- [ ] Run `go build ./cmd/headscale` and `go test ./hscontrol/...`.
- [ ] Update `CHANGELOG.md`, `session.md`, and the decision log below.
- [ ] Commit with a message that names the reviewed `FP-*` IDs.

## 6. Red flags and immediate action

| Symptom after merge | Feature | Immediate action |
| --- | --- | --- |
| `hscontrol/updatecheck/` missing | `FP-002` | Restore the whole directory and API handler |
| `go.mod` says `github.com/juanfont/headscale` | `FP-001` | Restore fork module path and imports |
| `hscontrol/api/v1/derp.go` missing | `FP-003` | Restore from `b0764611` |
| Webhook migration ID missing or duplicated | `FP-004` | Restore migration and `schema.sql` |
| `machine_approve_bulk` absent | `FP-005` | Restore `machines.go` and generated client |
| `headplane` gitlink changed | `FP-016` | Compare commits; run `check:fork` before accepting |
| `oidc_groups` or `oidcGroups` absent | `FP-007` | Restore type, migration, policy, and API fields |
| `hscontrol/qr/` missing | `FP-009` | Restore QR package and templates |
| `/ready` or `ReadyHandler` absent | `FP-011` | Restore handler and route |
| `docs/MERGE-PRESERVATION.md` deleted | `FP-013` | Restore from backup branch or git history |
| `*.md` removed from `.gitignore` | `FP-015` | Restore ignore rule unless explicitly retired |

## 7. Merge decision log

Append one row for every upstream merge. Never delete old rows.

| Date | Upstream SHA | Merge commit | FP IDs reviewed | Deletions accepted | Restores | Audit | Operator |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 2026-11-22 | `9ef033c7` | `3c788b7a` | `FP-001`–`FP-016` | none | none | baseline | Cline |
| 2026-11-23 | `0c304b06` | `9a56c617` | `FP-001`–`FP-016` | `.github/workflows/docs-deploy.yml`, `.github/workflows/version-bump.yml`, `nix/example-configuration.nix` | `headplane` gitlink restored to `9b8bfc2c` | PASS | Cline |

Notes for the 2026-11-23 merge:

- Upstream changes accepted: config loader and generated options, auth-cache
  admission and tagged-PAK re-registration, policy group pre-resolution,
  deleted-node 404 behavior, new migrations, and dependency updates.
- Upstream changes rejected to preserve fork behavior: DockerHub login in the
  main container workflow, replacement of the Headplane gitlink, and removal of
  the deprecated `headplane` config key.
- Features explicitly retired: none. The three deleted files above were already
  removed by the fork before this merge; upstream's modified versions were not
  restored.
- Follow-up work: environment-only test limitations remain (`TestTailSQLAccess`,
  `TestAPIv2` missing `tscli`/`tofu`, and the upstream servertest timeout).
- `FP-017` was added after release `v1.1.6` in commit `7aa662f2` when deployment
  configs using the legacy key locations were found to be rejected by the new
  strict config loader. No feature was retired.

Copy this template for the next merge:

```markdown
| YYYY-MM-DD | <upstream-sha> | <merge-sha> | <FP-ids> | <none or list> | <none or restore commits> | PASS/FAIL | <name> |
```

Notes to append below the table:

- Upstream changes accepted:
- Upstream changes rejected to preserve fork behavior:
- Features explicitly retired (requires approval):
- Follow-up work:

## 8. CI gate coverage

`.github/workflows/fork-features.yml` currently checks only four paths:

| Checked path | Feature |
| --- | --- |
| `hscontrol/updatecheck/` | `FP-002` |
| `hscontrol/api/v1/updatecheck.go` | `FP-002` |
| `hscontrol/api/v1/derp.go` | `FP-003` |
| `headplane/app/update-check/` | `FP-002` |

The following are **not** covered by CI and must be checked with section 4:

- `FP-001` module path and build metadata
- `FP-004` webhook migration and dispatcher
- `FP-005` machine approval API
- `FP-006` Headplane schema repair migrations
- `FP-007` OIDC groups
- `FP-009` QR registration
- `FP-010` server info
- `FP-011` readiness
- `FP-013` required documentation
- `FP-014` single-line fixes
- `FP-016` gitlink selection
- `FP-017` legacy config key aliases

Submodule-side features are covered by
`git -C headplane pnpm run check:fork`, which checks local admin files, login
integration, build defines, and the QR scanner.

## 9. Known limitations and current discrepancies

- The `headplane` gitlink recorded in `HEAD` is
  `9b8bfc2c15744e33cb9f8d9d2c998260bcff5aad`, but the checked-out submodule
  worktree is at `f9b2a3ed30c7e429f6ac4367dadf57fd05f8b7f0`
  (`git submodule status` shows a leading `+`). Do not reset the submodule
  without an explicit decision about which commit is authoritative.
- The submodule has no Doni-authored commits; preservation there means keeping
  the selected gitlink and the features it represents.
- Some docs still describe the superseded multi-user local-auth design and the
  `-arsydoni4326-alt` suffix. Treat this contract and the current code as the
  source of truth until those docs are reconciled.
- `FP-010` and the `hscontrol/qr_api.go` part of `FP-009` were authored by
  Denny Lastpaste, not Doni Arsy.
- The audit in section 4 is intentionally conservative. It detects deletion
  and token loss, not semantic regressions; run the card verification commands.

## 10. Maintenance

Update this contract when:

- a new fork-specific feature is added — assign the next `FP-*` ID and add a
  card;
- an upstream merge removes, renames, or replaces a listed path;
- the `headplane` gitlink changes;
- a documented convention is intentionally retired — record it in `FP-015`
  and the decision log;
- a verification command becomes stale.

If this file itself is lost during a merge, restore it from the backup branch
or from git history before doing anything else. It is the merge contract.

Record every decision in `session.md` and add a user-visible note to
`CHANGELOG.md` when behavior changes.

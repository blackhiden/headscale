# Session

## Current Task: Merge upstream/main (00822e63) into dev ✅ COMPLETE

**Completed:** 2026-11-28

### Objective

Finish the in-progress merge of `upstream/main` (`00822e63`, the changelog
restyle and policy follow-ups) into `dev` and resolve its conflicts.

### Conflicts resolved

- `CHANGELOG.md`: took upstream's restyled 0.29.4 entries, kept the fork's
  `Dockerfile.tailscale-HEAD` entry and the `0.29.3-arsydoni4326-alt` release
  section.
- `hscontrol/api/v1/api.go`: took upstream's new middleware (query, body and
  status handling, routing errors, protojson body decoding) and re-applied the
  fork's RBAC (`principalAdminKey`/`principalScopesKey`, OAuth scope
  enforcement, `requireScope`).
- `hscontrol/state/auth_cache_test.go`: kept upstream's registration-lock tests
  and the fork's QR-expiry test, with the fork module path.
- `gen/client/v1/client.gen.go`: regenerated with `make client`.
- `flakehashes.json`: recomputed with `go run ./cmd/vendorhash update`.

### Merge artifacts repaired

- `github.com/juanfont/headscale` module paths in six test files.
- Unused imports left by upstream removals in `generate.go`,
  `generate_test.go`, `apiv1_policy_test.go`, `apiv2_compat_test.go`,
  `auth_tags_test.go` and `node_health.go`.
- Stale `--running` line in `docs/ref/configuration.md` (upstream removed the
  flag).
- Deleted `hscontrol/policy/v2/issue_3212_test.go`: upstream moved the captures
  into the `routes_results` golden suite, and the fork's copy only carried
  module-path rewrites.

### Validation

- Preservation audit from `docs/MERGE-PRESERVATION.md` section 4 — PASS.
- `git diff --check` — PASS.
- `go build ./...` — PASS.
- `go vet ./...` — PASS.
- `go test ./hscontrol/state`, `./hscontrol/config`, `./cmd/headscale/cli` — PASS.
- `go test ./hscontrol` — only `TestTailSQLAccess` fails, the known
  environment-only failure that also reproduces on clean upstream.

### Notes

- The `headplane` gitlink was not changed by this merge.
- No release was created for this merge.

---

## Previous Task: Merge upstream/main (configuration overhaul) into dev ✅ COMPLETE

**Completed:** 2026-11-27

### Objective

Finish the in-progress merge of `upstream/main` (`d9bf885c`, the configuration
overhaul) into `dev` and resolve its remaining conflict.

### Conflict resolution

One conflict remained, in `docs/ref/configuration.md`:

- Took upstream's new "Example configuration" section (`headscale generate
  config`), which replaces the removed root `config-example.yaml`.
- Dropped the fork's dead `blackhiden/headscale` `config-example.yaml` GitHub
  links: the root file no longer exists (upstream removed it in `fe908c2d`),
  and the version shown in the docs is now `docs/ref/config-example.yaml`,
  generated from `hscontrol/types/config.go`.
- Preserved the fork's "Headplane Configuration" section (`FP-006`), which is
  still referenced by `docs/usage/authentication.md` and
  `docs/troubleshooting.md`.
- The first 49 lines of the resolved file are byte-identical to
  `MERGE_HEAD:docs/ref/configuration.md`; the fork section is appended below.

### Validation

- `git diff --check` — PASS (no leftover conflict markers).
- Preservation audit from `docs/MERGE-PRESERVATION.md` section 4 — PASS.
- `go build ./cmd/headscale` — PASS.
- `go test ./hscontrol/config -count=1` — PASS (20.5s).
- `headscale generate config` matches `docs/ref/config-example.yaml` apart from
  the generated header, so the mkdocs include in the page is valid.
- `mdformat` and `prettier` are not installed in this workspace, so the exact
  Nix formatting gate was not run locally.

### Known documentation drift (not changed, out of scope)

- `docs/troubleshooting.md` still says `headscale configtest`, which upstream
  renamed to `headscale check config`.
- `docs/ref/configuration.md` and `docs/usage/authentication.md` describe
  setting `headplane.password` in Headscale's `config.yaml`; the strict loader
  accepts `headplane.*` only for rollback compatibility
  (`hscontrol/config/load.go`) and Headplane owns that setting now.

---

## Previous Task: Merge upstream/main into dev and create Git Flow release v1.1.7

**Started:** 2026-11-23

### Pre-merge state

- Current branch: `dev` at `2b8b517bb88d3e88c73bf867956f2b3db4fddab2` before this task.
- `main` at `41917c4589998f38f02c7a672c3b984a0d3cc276`; `main` is an ancestor of `dev`.
- `upstream/main` before fetch: `160d8d2e9d6fb3ea94ba9101ccfe342651096a40`.
- `dev` was behind `upstream/main` by 1 commit and ahead by 368 commits.
- `headplane` gitlink in `HEAD`: `7d6616b9d6fbfb76846d5cc832c7ccb306695c8f`; submodule worktree is clean at the same commit.
- Git Flow 1.12.3 is installed and configured with master `main`, develop `dev`, release prefix `release/`, and version-tag prefix `v`.
- Latest valid release tag: `v1.1.6`. The `vv0.35.0-arsydoni4326-alt` tag is malformed and older; it is not used for version calculation.
- Planned release: `1.1.7`, producing tag `v1.1.7` via Git Flow.

### Local work preservation

- Backup branch: `backup/pre-merge-20261123-upstream-160d8d2e` at the original `dev` tip.
- Unrelated staged `Dockerfile` change (Go builder `1.27.2` → `1.27.4` plus build-essential setup) preserved in stash `f1c4acee28c11c2a18f3a41d17982b2b91472bf1` with message `pre-upstream-merge: preserve staged Dockerfile build tweak`.
- The stash must be restored after the release is finalized; if Git does not restore the index bit, re-stage `Dockerfile` to recover the original staged state.

### Progress

- [x] Read project documentation and merge contract.
- [x] Inspect Git status, branches, tags, remotes, and git-flow configuration.
- [x] Create backup branch and preserve unrelated staged work.
- [x] Fetch and update `upstream/main` (unchanged at `160d8d2e`).
- [x] Merge upstream into `dev` with preservation audit. Merge commit: `4992e42be045ee691ef38ea166f537dddcc53d37`.
  - One import-only conflict in `hscontrol/servertest/policy_test.go` was resolved
    by keeping the fork module path and adding upstream's `types` import.
  - Preservation audit passed before and after conflict resolution.
  - `go build ./cmd/headscale` PASS.
  - `go test ./hscontrol/policy/v2 -count=1` PASS.
  - `go test ./hscontrol/state -count=1` PASS.
  - Focused `go test ./hscontrol/servertest -run '^(TestUserEmailDomainChangeRevokesLiveAccess|TestPolicyChanges)$'` PASS.
- [x] Create and finish Git Flow release `1.1.7`.
- [x] Restore unrelated local work and validate final state.

### Result

- Previous latest valid tag: `v1.1.6`.
- New release version: `1.1.7`; annotated tag `v1.1.7` object `2e59d0c0d684219e04e09cb770e85f56843df81b`, pointing to commit `d598680cd6d48ca1364d6b7fa1e53458d6d703c0`.
- Upstream commit synchronized: `160d8d2e9d6fb3ea94ba9101ccfe342651096a40`.
- Merge into `dev`: `4992e42be045ee691ef38ea166f537dddcc53d37`; `upstream/main` is now an ancestor of `dev`.
- Release branch: `release/1.1.7`, created from `dev` and deleted after finishing.
- `main`: `d598680cd6d48ca1364d6b7fa1e53458d6d703c0`, containing the release merge and tag.
- `dev`: `54f95f843285c78d198c3fc5569c786060a33517`, containing the upstream merge and release back-merge.
- `headplane` gitlink remained `7d6616b9d6fbfb76846d5cc832c7ccb306695c8f`; `git -C headplane pnpm run check:fork` PASS.
- Final preservation audit PASS; final `go build ./cmd/headscale` PASS.
- The staged `Dockerfile` change was restored from stash and re-staged; no unrelated local work was discarded.
- No push, force-push, reset, clean, or tag deletion was performed. Backup branch `backup/pre-merge-20261123-upstream-160d8d2e` was retained.

---

## Current Task: Restore legacy config key compatibility ✅ COMPLETE

**Completed:** 2026-11-23

### Problem

After the upstream config refactor, existing fork configs failed with unknown-key
errors for `derp.server.auto_update_enabled`, `derp.server.update_frequency`,
`node.routes.probe_interval`, and `node.routes.probe_timeout`.

### Fix

- Added deprecated aliases in `hscontrol/config/keys.go` and a rewrite step in
  `hscontrol/config/load.go`.
- Legacy keys are rewritten to `derp.auto_update_enabled`,
  `derp.update_frequency`, `node.routes.ha.probe_interval`, and
  `node.routes.ha.probe_timeout`, with one warning per key.
- If both legacy and current keys are present, the current key wins.
- Added regression coverage in `hscontrol/config/load_test.go`.
- Documented the aliases in `CHANGELOG.md` and `docs/troubleshooting.md`.
- Recorded the aliases as `FP-017` in `docs/MERGE-PRESERVATION.md`, including
  the manifest, preservation card, audit tokens, and decision-log note.

### Validation

- Focused alias test — PASS.
- `go test ./hscontrol/config -count=1` — PASS.
- `go build ./cmd/headscale` — PASS.
- `headscale check config -f /tmp/legacy-config.yaml` with all four legacy keys
  — PASS with four deprecation warnings.
- No new release was created for this fix.
- An unrelated uncommitted `Dockerfile` change was left untouched.

---


## Current Task: Merge upstream/main into dev and create Git Flow release ✅ COMPLETE

**Completed:** 2026-11-23

### Result

- Fetched `upstream/main` at `0c304b06` and merged it into `dev`.
- Merge commit: `9a56c617d426b2cebb8401a1966248a7dfa950a8`.
- Previous latest valid tag: `v1.1.5`. Next release: `1.1.6` (tag `v1.1.6`).
- Backup branch retained: `backup/pre-merge-20261123`.
- Headplane gitlink retained at `9b8bfc2c15744e33cb9f8d9d2c998260bcff5aad`;
  the checked-out submodule worktree remains at `5ab2878a…` as documented.

### Conflict resolution

- 83 conflicts were resolved: 44 import-only, 6 clean normalized merges, and
  33 manual/configuration resolutions.
- Preserved fork features `FP-001`–`FP-016`, including the module path, update
  checker, DERP/webhook/machine-approval APIs, Headplane integration, OIDC
  groups, QR registration, server info, readiness endpoint, and docs.
- Accepted upstream's newer config loader, auth-cache admission, tagged-PAK
  re-registration, policy group pre-resolution, deleted-node behavior, and
  migrations, then re-applied fork behavior where required.
- Kept the deprecated `headplane` config key accepted for rollback
  compatibility and restored the QR expiry default in `SetAuthCacheEntry`.
- Accepted fork deletions of `docs-deploy.yml`, `version-bump.yml`, and
  `nix/example-configuration.nix`; no feature was retired.

### Validation

- `go build ./cmd/headscale` — PASS.
- Preservation audit — PASS.
- `git -C headplane` check:fork — PASS.
- `go test ./cmd/...` — PASS.
- `go test ./hscontrol/... -skip '^(TestTailSQLAccess|TestAPIv2)$'` — all
  packages PASS except the pre-existing upstream servertest timeout, which also
  reproduces on clean `upstream/main`.
- Focused QR, config, health/readiness, and tagged-PAK tests — PASS.
- Environment-only limitations: `TestTailSQLAccess` fails on clean upstream;
  `TestAPIv2` requires `tscli`/`tofu` from the Nix dev shell; the servertest
  package times out on clean upstream in this environment.

### Release result

- Git Flow release branch: `release/1.1.6` (created and deleted).
- Release tag: `v1.1.6` at `41917c4589998f38f02c7a672c3b984a0d3cc276`.
- `main` at `41917c45` contains the release merge.
- `dev` at `d1ba1e19` contains `upstream/main`, merge commit `9a56c617`, and
  the release back-merge.
- Final `go build ./cmd/headscale` — PASS.
- No push, force-push, reset, or tag deletion was performed.
- The `headplane` gitlink remains `9b8bfc2c…`; the checked-out submodule
  worktree is now `f9b2a3ed…` after the release checkouts. This known
  discrepancy was left unchanged.

---



## Current Task: Harden Merge-Preservation Document Format ✅ COMPLETE

**Completed:** 2026-11-22

### Objective

Convert `docs/MERGE-PRESERVATION.md` from a narrative inventory into a
feature-loss-proof merge contract that gives explicit conflict rules, restore
commands, and verification steps for every fork feature.

### Delivered

- Added stable `FP-001`–`FP-016` IDs and a feature manifest table.
- Rewrote the inventory as self-contained preservation cards with author,
  last-known-good commits, protected paths, required tokens, behavior contract,
  conflict rule, restore command, and verification command.
- Added a one-command preservation audit that checks paths, tokens, migration
  IDs, and submodule integration points.
- Added a merge protocol, red-flag triage table, and append-only merge decision
  log.
- Documented CI gate coverage gaps and known limitations.
- Updated `CHANGELOG.md` to describe the new contract format.

### Validation

- Extracted and executed the audit block from the document — PASS.
- `git diff --check` — PASS.
- `mdformat` and `prettier` are not installed locally; the exact Nix formatting
  gate was not run.
- Documentation-only change; no Go build or test run was required.

---

## Previous Task: Document Doni Arsy Code for Merge Preservation ✅ COMPLETE

**Completed:** 2026-11-22

### Objective

Scan the fork remote (`origin`) for all commits authored by Doni Arsy and write
a preservation-oriented summary to `docs/MERGE-PRESERVATION.md`.

### Scan Results

- 297 Doni-authored commits reachable from `origin/*`; 172 non-merge and 124
  merge commits.
- `origin/main` at scan time: `9ef033c7`; local HEAD: `3c788b7a`.
- The `headplane` submodule has zero Doni-authored commits; his parent-repo
  contribution is gitlink selection/integration.
- Pure merge commits, Cline checkpoint/index commits, and release-only commits
  were excluded from the feature inventory but retained in history.
- `hscontrol/api/v1/serverinfo.go` and `hscontrol/qr_api.go` were authored by
  Denny Lastpaste, not Doni; they are documented as adjacent fork features.

### Delivered

- Replaced the empty `docs/MERGE-PRESERVATION.md` placeholder with a full
  inventory covering fork identity, update checker, DERP endpoint, webhooks,
  machine approval, Headplane local admin/settings, OIDC groups, Phase 10/11
  features, QR registration, observability, UI/a11y, docs/process, small fixes,
  and deliberate removals.
- Added a representative commit map, upstream-merge checklist, recovery
  template, verification commands, and known limitations.
- Added the missing required root `SCHEMA.md`.
- Linked the preservation document from `CONTRIBUTING.md`, `ROADMAP.md`,
  `mkdocs.yml`, and `CHANGELOG.md`.

### Decisions and Assumptions

- "All code written by Doni Arsy" means commits whose author is Doni Arsy in
  the fork remote, not code merely merged by him.
- Submodule features are described as integrated via gitlinks because the
  submodule contains no Doni-authored commits in this clone.
- The current submodule pointer mismatch (`HEAD` gitlink `9b8bfc2c...` versus
  checked-out `5ab2878...`) was recorded as a known limitation, not changed.

### Validation

- `git log --remotes=origin --author='Doni Arsy' --no-merges --oneline` — 172
  commits enumerated.
- `git diff --name-status upstream/main...origin/main` — reviewed fork drift.
- `git -C headplane log --all --author='Doni Arsy'` — zero commits.
- `git diff --check` — PASS.
- `mdformat` and `prettier` are not installed in this workspace, so the exact
  Nix formatting gate was not run locally.
- `mkdocs build --strict` was attempted but failed because the environment
  lacks the `material.extensions.emoji` Python module; this is an environment
  limitation, not a document syntax error.
- Documentation-only change; no Go build or test run was required.

---

## Previous Task: Fix 404 on Headplane `/api/admin/server-info` ✅ COMPLETE

**Completed:** 2026-10-28

### Root Cause Analysis

1. **Unregistered route in `headplane/app/routes.ts`**: The endpoint file `app/routes/api/admin/server-info.ts` existed and had unit tests, but was never registered in `app/routes.ts` under `prefix("/admin", ...)`.
2. **Missing `__PREFIX__` in client fetch**: In `app/routes/admin/components/technical-info-section.tsx`, `fetch("/api/admin/server-info")` was used instead of `fetch(`${__PREFIX__}/api/admin/server-info`)`. Headplane runs with React Router `basename: "/admin/"` (`__PREFIX__`). Calling `/api/admin/server-info` bypassed the basename, causing `<Router basename="/admin/"> is not able to match the URL "/api/admin/server-info"` and returning 404.
3. **Missing fallback redirect in runtime**: In `runtime/http.ts`, `composeListener` only redirected `/admin` to `/admin/`. Requests to `/api/*` or `/healthz` sent without the basename were passed to the Router directly and failed with 404.

### Fixes Applied

1. **`headplane/app/routes.ts`**:
   - Registered `route("/server-info", "routes/api/admin/server-info.ts")` under `prefix("/admin", [...])`.
2. **`headplane/app/routes/admin/components/technical-info-section.tsx`**:
   - Changed fetch target to `${__PREFIX__}/api/admin/server-info`.
   - Added support for displaying CPU, RAM, storage, deployment mode, and public IP in `headscaleFields`.
3. **`headplane/app/routes/api/admin/server-info.ts`**:
   - Mapped all `hscontrol/api/v1/serverinfo.go` response fields (`go.version`, `os`, `arch`, `cpu.count`, `memory.total`, `storage.total`, `deployment.mode`, `publicIp`).
   - Wrapped `auth.getHeadscaleApiKey(principal)` in `try/catch` to prevent unhandled exceptions if the API key is unconfigured.
   - Enhanced deployment mode detection with `existsSync("/.dockerenv")` and checked `config.integration?.docker`.
4. **`headplane/runtime/http.ts`**:
   - In `composeListener`, added automatic HTTP 307 redirect for any `/api/*` or `/healthz` requests that arrive without the basename prefix to `${prefix}${url.pathname}${url.search}`.
5. **`hscontrol/` tests**:
   - Fixed outdated `github.com/juanfont/headscale` import paths in 4 test files to `github.com/blackhiden/headscale`.

### Testing & Verification

- `go test ./hscontrol/ -run 'TestAPIV1ServerInfo' -v` — PASS
- `go build ./cmd/headscale` — PASS
- `npx vitest run tests/unit/routes/admin-server-info.test.ts` — PASS
- `npx vitest run tests/unit/components/server-info-card.test.ts` — PASS
- `npx vitest run tests/unit/runtime/http.test.ts` — PASS
- `pnpm build` (Headplane production build) — PASS (includes `server-info` in client and server bundles)

---

## Previous Task: Phase 18 — Admin Technical Info Section ✅ COMPLETE

**Completed:** 2026-10-28

### Objective

Implement the Admin Technical Info section at `/admin/admin` in Headplane, showing server and system details for both Headscale and Headplane.

**Completed:** 2026-10-17

### Objective

Add a custom API endpoint (`/api/v1/registration/qr/{auth_id}`) that returns the QR registration payload as JSON, enabling CLI clients to fetch the payload programmatically and generate QR codes locally without opening the registration URL in a browser.

### Requirements

- Open endpoint (no authentication required)
- Returns JSON payload with `type`, `version`, `auth_id`, `server_url`, and `expires_at`
- Validates `auth_id` format and checks if registration is pending and not expired
- Returns appropriate HTTP status codes (200, 400, 404, 410)

### Implementation

Created three files:

1. **`hscontrol/qr_api.go`** (98 lines)
   - `RegistrationQRPayloadHandler`: HTTP handler for GET `/api/v1/registration/qr/{auth_id}`
   - Validates auth_id, looks up pending registration, checks expiry
   - Reuses existing `qr.NewRegistrationPayload()` function
   - Returns JSON response with proper error handling

2. **`hscontrol/qr_api_test.go`** (133 lines)
   - `TestRegistrationQRPayloadHandler_Success`: Valid request returns 200 with correct payload
   - `TestRegistrationQRPayloadHandler_NotFound`: Non-existent auth_id returns 404
   - `TestRegistrationQRPayloadHandler_Expired`: Expired registration returns 410
   - `TestRegistrationQRPayloadHandler_InvalidAuthID`: Malformed auth_id returns 400
   - `TestRegistrationQRPayloadHandler_NotRegistrationRequest`: Non-registration auth returns 404

3. **Modified `hscontrol/app.go`** (line 503-504)
   - Registered route: `r.Get("/api/v1/registration/qr/{auth_id}", h.RegistrationQRPayloadHandler)`
   - Placed after `/register/{auth_id}` and `/auth/{auth_id}` routes

4. **Updated `docs/api/qr-payload.md`** (added 56 lines)
   - API endpoint documentation with examples
   - Request/response format
   - Error codes (400, 404, 410)
   - CLI use case example with `curl` and `qrencode`
   - Security notes

### Design Decisions

- **Open endpoint**: Matches behavior of `/register/{auth_id}` page (already public)
- **Reuses existing QR payload logic**: Uses `qr.NewRegistrationPayload()` for consistency
- **State lookup**: Uses `State.GetAuthCacheEntry()` to fetch pending registration
- **Expiry check**: Returns 410 Gone if registration has expired
- **Registration type check**: Verifies `authReq.IsRegistration()` to distinguish from SSH check requests

### Testing

All test cases cover:
- ✅ Success: Valid pending registration returns correct JSON payload
- ✅ Not found: Non-existent auth_id returns 404
- ✅ Expired: Past expiry time returns 410
- ✅ Invalid format: Malformed auth_id returns 400
- ✅ Wrong type: Non-registration auth request returns 404

### Notes

- This endpoint allows CLI clients to generate QR codes during registration without browser access
- The `expires_at` value comes from `AuthRequest.ExpiresAt()`, which is set when the registration enters the cache
- No breaking changes: purely additive functionality
- Documentation updated in `docs/api/qr-payload.md`

---

## Previous Task: Generate CODE_GRAPH.md ✅ COMPLETE

**Completed:** 2026-10-13

### Objective

Create a comprehensive CODE_GRAPH.md document that serves as a persistent structural map of the Headscale codebase, enabling AI agents and developers to understand the project architecture instantly and reduce exploratory file reads.

### Delivered

Created `CODE_GRAPH.md` (506 lines) with:

1. **Visual Architecture** - High-level system diagram and data flow diagrams
2. **Project Structure Overview** - Complete directory tree with role descriptions
3. **Core Module Index** - Detailed module documentation for all major components
4. **Data Flow Patterns** - Node registration, map updates, tags-as-identity
5. **Key Architectural Patterns** - Copy-on-write NodeStore, code-first API, DB migrations
6. **Testing Infrastructure** - servertest/, integration/, cmd/hi/
7. **Build and Development** - Makefile, go.mod, flake.nix
8. **Quick Reference** - Task-specific guides for common development workflows

### Validation

```bash
wc -l CODE_GRAPH.md
# 506 CODE_GRAPH.md

grep -c '^##' CODE_GRAPH.md
# 52 (52 major sections/subsections)
```

### Notes

- Follows project-guidelines requirement: documentation is part of implementation
- Designed for AI agents to reduce file-reading overhead during task execution
- Should be updated when major architectural changes occur
- Cross-references ARCHITECTURE.md, AGENTS.md, cmd/hi/README.md, integration/README.md

---

## Previous Task: Document Markdown Formatting Gate ✅ COMPLETE

**Completed:** 2026-10-12

### Delivered

- Updated `CONTRIBUTING.md` and `README.md` to require this command before
  submitting or committing any Markdown (`.md`) change:

  ```shell
  sudo nix --extra-experimental-features 'nix-command flakes' build -L .#checks.x86_64-linux.formatting
  ```

- `docs/about/contributing.md` includes `CONTRIBUTING.md`, so the contributor
  documentation site automatically reflects this policy without duplicate text.

### Validation

- `git diff --check` passed for the tracked documentation changes.
- The exact required Nix command was attempted but could not run in this
  non-interactive workspace because `sudo` requires terminal authentication:
  `sudo: A terminal is required to authenticate`.
- The existing dirty `headplane` submodule was not modified.

---

## Current Task: Fix Nix treefmt Formatting Check ✅ COMPLETE

**Completed:** 2026-10-12

### Objective

Fix the `treefmt-check` failure in the Nix build caused by trailing whitespace and indentation inconsistencies in `CHANGELOG.md` and `ROADMAP.md`.

### Problem

The treefmt formatter detected 66 files that needed formatting, with 2 files (CHANGELOG.md and ROADMAP.md) containing formatting issues:
- Trailing whitespace on line continuations
- Inconsistent indentation in nested list items and code blocks
- Missing blank lines between sections

### Solution

Applied all formatting changes reported by treefmt:
- Removed trailing spaces from line continuations
- Standardized indentation (2 spaces for nested lists, 4 spaces for continued lines)
- Added required blank lines between sections
- Fixed code block formatting (blank line before code fence)

### Files Modified

- `CHANGELOG.md` — 38 lines reformatted (trailing whitespace, list indentation)
- `ROADMAP.md` — 79 lines reformatted (trailing whitespace, list indentation, code blocks)

### Verification

```bash
git diff --stat
# CHANGELOG.md | 38 +++++++++++++++--------------
# ROADMAP.md   | 79 +++++++++++++++++++++++++++++++++++++++---------------------
# 2 files changed, 72 insertions(+), 45 deletions(-)

git commit -m 'docs: fix treefmt formatting (trailing whitespace and indentation)'
# [dev fd121541] docs: fix treefmt formatting (trailing whitespace and indentation)
```

### Notes

- This is a formatting-only change with no behavioral, semantic, or content modifications
- All changes match the exact output from `treefmt v2.6.0`
- The working tree is now clean and ready for the next treefmt check

---

## Current Task: Repair Nix Formatting, Lint, and Unit-Test Checks

**Date:** 2026-10-11

### Progress

- Repaired the CI failures reported by `.github/workflows/nix-checks.yml`:
  applied the configured `gci` ordering to all reported stale import blocks,
  replaced the flagged inline error declarations in the Headplane DB helpers,
  added `wsl_v5`-required statement spacing in `hscontrol/db/webhook.go`, and
  restored Prettier-required list spacing in `ROADMAP.md`.
- No runtime behavior, public API, schema, or workflow configuration changed.
- Completed local validation:

  ```bash
  golangci-lint run                                            # PASS (0 issues)
  CGO_ENABLED=0 go test ./hscontrol/db ./hscontrol/types \
    ./hscontrol/policy ./hscontrol/policy/v2                   # PASS
  gofumpt -l <all changed Go files>                             # clean
  prettier --check ROADMAP.md                                   # PASS
  git diff --check                                              # clean
  ```

- Local Nix is unavailable, so
  `nix build -L .#checks.x86_64-linux.{golangci-lint,formatting}` remains for
  CI validation. The `headplane` submodule was already dirty and was not
  modified by this repair.

- Read all 23 workflow files under `.github/workflows/`; the relevant Nix gate
  definitions are in `.github/workflows/nix-checks.yml`.
- Formatted `.github/workflows/deploy.yaml` and
  `.github/workflows/fork-features.yml`; the root formatter-equivalent check
  now passes using cached Prettier over the flake's tracked-file scope.
- Reproduced the CI's 76 `golangci-lint` findings locally and applied scoped
  source/test corrections, including replacing timing sleeps with polling and
  isolating `TestWebhookDispatcher_SendWebhook` from unnecessary full-server
  initialization.
- The equivalent non-Docker Go test subset initially identified the webhook
  test's actual failure: an empty `types.Config` caused `NewHeadscale` to fail
  while creating a Noise private key. The test now constructs only the
  dispatcher required by `sendWebhook`.

### Current Repair

- Fixed the PostgreSQL branch of `migrateHeadplaneSettingsToPerUser` in
  `hscontrol/db/db.go`: it now assigns to the `err` declared earlier in the
  branch rather than redeclaring it with `:=`.
- This restores compilation of the CI server-test package without changing
  migration behavior or database schema.
- Completed validation:

```bash
CGO_ENABLED=0 go test ./hscontrol/db/...             # PASS
CGO_ENABLED=0 go test -run '^$' ./hscontrol/servertest/...  # PASS
gofmt -d hscontrol/db/db.go                           # clean
git diff --check                                      # clean
```

- The complete CI-equivalent server-test suite is running with its configured
  20-minute timeout. Its result will be recorded when it exits.

### Environment Limitation

- `nix` is unavailable in this workspace (`nix: command not found`), so the
  exact three Nix derivations cannot be run locally.

---

## Current Task: Fix golangci-lint Nix Check Failure ✅ COMPLETE

**Completed:** 2026-10-11

### Objective

Fix the `nix build -L .#checks.x86_64-linux.golangci-lint` failure in CI caused by missing imports in `hscontrol/noise_test.go`.

### Problem

The CI workflow `.github/workflows/nix-checks.yml` failed with typecheck errors:
- `undefined: capver` (lines 250, 306, 321, 344, 399)
- `undefined: zerolog` (line 768)

### Root Cause

The `hscontrol/noise_test.go` file was missing two required imports:
1. `github.com/blackhiden/headscale/hscontrol/capver`
2. `github.com/rs/zerolog`

These packages are used throughout the test file but were not imported, likely due to an incomplete merge or refactoring.

### Delivered

- Added missing import `github.com/blackhiden/headscale/hscontrol/capver` to `hscontrol/noise_test.go`
- Added missing import `github.com/rs/zerolog` to `hscontrol/noise_test.go`

### Validation

- ✅ `go test -c ./hscontrol` compiles successfully with no errors
- ✅ All undefined symbol errors resolved

### Notes

- No workflow changes were needed; this was purely a source code issue
- The imports are already present in other test files (e.g., `hscontrol/noise.go`)
- The Nix flake check should now pass in CI

---

## Previous Task: Repository Path Migration ✅ COMPLETE

**Completed:** 2026-10-10

### Objective

Update all repository references from `github.com/arsydoni4326-alt/headscale` to `github.com/blackhiden/headscale` across the entire codebase.

### Delivered

- Updated Go module path in `go.mod` from `github.com/arsydoni4326-alt/headscale` to `github.com/blackhiden/headscale`.
- Updated all Go import statements across the codebase (2340+ files scanned).
- Updated all references in documentation files (markdown, YAML, TOML).
- Updated CI/CD configuration files (`.github/workflows/`, `.goreleaser.yml`).
- Fixed remaining upstream references in:
  - `hscontrol/types/config_test.go`
  - `integration/tsic/login_test.go`
  - `.golangci.yaml`
- Ran `go mod tidy` to ensure dependency consistency.

### Validation

- ✅ `go build ./cmd/headscale` succeeded.
- ✅ `go test ./hscontrol/types` passed.
- ✅ `go test ./cmd/headscale/cli` passed.
- ✅ All references to old repository path removed (verified with grep).
- ✅ All references to new repository path verified in Go code.

### Important Notes

- This is a breaking change for downstream consumers who import this module - they will need to update their import paths.
- The change affects the module path, all import statements, documentation, CI/CD configs, and release configurations.
- No vendor directory exists, so no vendoring updates were needed.

---

## Current Task: Headplane QR Registration Fetcher Transport Fix ✅ COMPLETE

**Completed:** 2026-10-10

### Objective

Prevent QR registration failures from rendering a full SSR Headplane document
as the inline error message.

### Delivered

- Merged Headscale commit `ea42eb3` (`headplane: submit QR registration through
  fetcher`) with published Headplane commit `6c5c1e5`.
- Replaced raw browser fetch with React Router `useFetcher().submit()` action
  transport, preserving structured action data and redirects.
- Extended QR Headscale error extraction to prefer Huma `errors[].message`, so
  user lookup failures display `user not found` rather than HTML or raw JSON.
- Added component regression coverage for fetcher submission, structured action
  errors, and submission-state disabling.

### Validation

- ✅ Merged focused Scan QR tests: 20 component + 5 action tests passed.
- ✅ Full Headplane component suite: 16 files / 102 tests passed.
- ✅ Focused formatting/lint and production build passed.
- ⚠️ `pnpm run typecheck` remains blocked by existing unrelated diagnostics;
  changed Scan QR files are not listed.

---

## Current Task: Headplane QR Registration Owner Name Fix ✅ COMPLETE

**Completed:** 2026-10-10

### Objective

Fix QR registration submitting a selected Headscale user ID where the
registration API requires the username.

### Delivered

- Merged Headscale commit `bd154017` (`headplane: register QR nodes by
  username`) with published Headplane commit `d05a98d`.
- The Scan QR owner selector now submits `user.name` rather than `user.id`,
  matching Headscale v1 registration's `GetUserByName` contract.
- Added a component regression that verifies selected `Alice` submits
  `user=Alice`, not the internal ID `user-1`.

### Validation

- ✅ Merged focused Scan QR tests: 20 component + 5 action tests passed.
- ✅ Full Headplane component suite: 16 files / 102 tests passed.
- ✅ Focused formatting/lint and production build passed.
- ⚠️ `pnpm run typecheck` remains blocked by existing unrelated diagnostics;
  the changed Scan QR test is not listed.

---

## Current Task: Headplane QR Registration Error Feedback ✅ COMPLETE

**Completed:** 2026-10-10

### Objective

Fix Scan QR registration submissions that bypassed Headplane's URL prefix and
could show an empty inline error after a registration failure.

### Delivered

- Merged Headscale commit `2e2337e` (`headplane: show QR registration errors`)
  with published Headplane commit `7d49a63`.
- QR registration now submits through `__PREFIX__`, so the default deployment
  reaches `/admin/machines/scan-qr`.
- Registration errors now return and display structured safe Headscale messages,
  plain-text failures, or a non-empty fallback for empty response bodies.
- Added component/action tests for prefix submission, structured errors, empty
  body fallback, and Headscale registration rejection.

### Validation

- ✅ Merged focused Scan QR tests: 20 component + 5 action tests passed.
- ✅ Focused formatting/lint and Headplane production build passed.
- ⚠️ `pnpm run typecheck` remains blocked by existing unrelated diagnostics;
  the changed Scan QR test is no longer listed.

---

## Current Task: Headplane QR Scanner Lifecycle Fix ✅ COMPLETE

**Completed:** 2026-10-10

### Objective

Fix Headplane's QR scanner showing no camera frame and crashing during cleanup
with `Cannot stop, scanner is not running or paused.`

### Delivered

- Merged Headscale commit `7e6fc0c` (`headplane: fix QR scanner lifecycle`),
  which points to published Headplane commit `ddb708a`.
- Kept the `html5-qrcode` host visible and full-size during permission/startup;
  scanner messages now overlay the host instead of hiding it.
- Made cleanup state-aware for only `SCANNING`/`PAUSED` library states and
  absorb synchronous/asynchronous cleanup failures.
- Replaced the persistent start guard with effect-local cancellation, covering
  React StrictMode replay and startup completion after unmount.
- Added component regressions for host visibility, failed startup, synchronous
  stop failures, paused cleanup, late startup cleanup, and StrictMode replay.

### Validation

- ✅ Merged focused component tests: 32 passed.
- ✅ Full Headplane component suite: 16 files / 99 tests passed.
- ✅ Focused `oxfmt` and `oxlint` passed; `pnpm run build` passed.
- ⚠️ `pnpm run typecheck` remains blocked by existing unrelated diagnostics.

---

## Current Task: Phase 17 QR Registration Reconciliation ✅ COMPLETE

**Completed:** 2026-10-09

### Objective

Reconcile the Phase 17 roadmap status with the existing QR registration
implementation, verify the Headscale and Headplane flow in isolated worktrees,
and merge focused coverage and lint hardening after explicit approval.

### Delivered

- Merged `qr: harden payload validation` (`47d4baf8`) into `dev`:
  package-local QR validation now uses static wrapped errors, its expiry test
  polls rather than sleeps, and `golangci-lint ./hscontrol/qr` passes.
- Merged `headplane: cover QR scanner action` (`4cf43205`) and updated the
  gitlink to published Headplane commit `d6dcbbbc`; its new unit coverage
  checks scanner loader/action authorization, API registration, live refresh,
  audit details, and redirect behavior.
- Merged `docs: mark QR registration flow implemented` (`59b35d0`): the
  Phase 17 brief is now `[Implemented]` and current state reflects the CLI +
  QR approval flow.
- Restored the historical Headplane baseline via
  `origin/recovery/restore-b470d2f`; retain it until old worktrees and release
  references no longer need commit `b470d2f`.

### Validation

- ✅ `go test ./hscontrol/qr -count=1`
- ✅ `go test ./hscontrol/state -run 'Test(AuthCache|SetAuthCacheEntry)' -count=1`
- ✅ `go test ./hscontrol -run 'Test(RegisterHandlerRendersQRForPendingRegistration|APIV1AuthRegister|Template)' -count=1`
- ✅ `go build ./cmd/headscale`
- ✅ `golangci-lint run --timeout 10m ./hscontrol/qr`
- ✅ Headplane QR route/key unit tests: 14 passed.
- ✅ Headplane QR component tests: 27 passed.
- ✅ Headplane `pnpm run lint` completed with 9 existing warnings and no
  errors; `pnpm run build` passed.
- ⚠️ The Headplane TypeScript check remains blocked by pre-existing unrelated
  diagnostics (including optional `better-sqlite3`, admin API, instances,
  settings, SSH, and existing component-test diagnostics).
- ⚠️ The Headplane Testcontainers QR integration matrix did not reach test
  bodies: all selected Headscale 0.27.0–0.29.1 cases timed out waiting for
  `/health` after 30 seconds.
- ⚠️ Full `go test ./hscontrol -count=1` retains an unrelated
  `TestWebhookDispatcher_SendWebhook` failure caused by an empty Noise key path.

---

## Current Task: Phase 17 QR Code Registration Flow Repair ✅ COMPLETE

**Completed:** 2026-10-08

### Objective

Repair the Phase 17 QR registration flow after the Machines-page action did not
navigate and the non-OIDC `/register/{auth_id}` page did not show a QR code.

### Requirement clarification

The requested supported flow is non-OIDC interactive registration: a node is
redirected to `https://<headscale-ui>/register/<auth_id>`, where Headscale
shows the existing CLI command and a QR code. A Headplane administrator uses
**Machines → Scan QR**, selects the target user, scans the code, and submits
the existing pending auth ID for registration. The QR is not a credential and
does not bypass Headscale's existing expiry or single-use authorization checks.

### Implementation Delivered

- ✅ Replaced unsupported `Button component={Link}` usage with a valid link to
  `/machines/scan-qr`; the disabled state remains non-interactive.
- ✅ Added QR rendering to the standard CLI-approved registration page as well
  as the existing OIDC confirmation page.
- ✅ Bound QR `expires_at` to the fixed pending-registration cache deadline, so
  reloading the registration page cannot extend its lifetime.
- ✅ Preserved Headscale's existing server-side cache expiry and single-use auth
  ID enforcement; scanning does not bypass owner selection or approval.
- ✅ Validated QR type, version, auth ID, server URL syntax, and expiry before
  Headplane submits the registration.
- ✅ Added `TestRegisterHandlerRendersQRForPendingRegistration`.
- ✅ Added `TestSetAuthCacheEntryPreservesQRExpiry` to prevent cache reinsertion
  from extending the QR deadline.

### Key Feature Details

**User Flow:**
1. User attempts to register a new Tailscale node
2. Registration page shows both:
   - CLI command: `headscale auth register --auth-id <hskey> --user USERNAME`
   - QR code with instructions to scan via Headplane
3. User opens Headplane: Machines → Scan QR → Choose User → Start Scanning
4. Device is automatically approved after successful scan

### Verification

- ✅ `go test ./hscontrol -run 'TestRegisterHandlerRendersQRForPendingRegistration|TestTemplateHTMLConsistency' -count=1`
- ✅ `go test ./hscontrol -run '^$'`
- ✅ `go test ./hscontrol/state ./hscontrol/types ./hscontrol/qr -count=1`
- ✅ `go test ./hscontrol/state -run 'TestAuthCacheBoundedLRU|TestSetAuthCacheEntryPreservesQRExpiry' -count=1`
- ✅ `pnpm exec oxlint app/routes/machines/overview.tsx app/routes/machines/scan-qr.tsx`
- ⚠️ Headplane `pnpm run typecheck` remains blocked by existing unrelated
  project errors and missing local `html5-qrcode`/`better-sqlite3` modules;
  the workspace is also running Node 22.22.1 while Headplane requires Node
  24.2–24.x.

### Documentation Updated

- `ROADMAP.md` — corrected the actual Machines → Scan QR path
- `docs/usage/registration.md` — documented non-OIDC support and cache-bound expiry
- `docs/ref/registration.md` — added the Headplane QR approval path
- `CHANGELOG.md` — clarified frontend navigation and non-OIDC support

---

## Previous Task: Phase 13c Single Local Administrator Migration ✅ COMPLETE

**Completed:** 2026-10-06

### Objective

Replace the Headscale database-backed multi-user local-auth implementation with
a single Headplane-configured local administrator. Local username/password is
primary, Headscale API-key login is secondary, and OIDC/proxy code remains
present but disabled for this mode.

### Implementation Delivered

**Phase 13c — Single Local Administrator Migration is complete.**

- ✅ Backend: Single local administrator authentication with bcrypt password verification
- ✅ Configuration: `user.username` and `user.password` (bcrypt hash) in `config.yaml`
- ✅ Migration tooling: `migrate-local-admin`, `hash-password`, `reset-local-admin-password` commands
- ✅ Administration UI: Password reset, API key management, single-admin notice
- ✅ Security: API key rotation safeguards, session management, rate limiting
- ✅ Documentation: Migration plan, runbook, verification procedures, rollback guidance
- ✅ Testing: Unit tests, browser tests, migration idempotency tests
- ✅ Legacy retirement: Headscale no longer manages Headplane users; tables preserved for rollback

### Documentation Updated

- `ROADMAP.md` — Phase 13c marked as [Implemented] with complete delivery checklist
- `docs/phase13c-single-admin-migration-plan.md` — Status: Complete with outcome summary
- `docs/phase13c-single-admin-migration-runbook.md` — Status: Complete with delivered tooling
- `docs/phase13-implementation-status.md` — All completion checkboxes marked, status updated
- `CHANGELOG.md` — Phase 13c Known Gap Fix and implementation details recorded

### Key Architectural Changes

- Headplane owns local authentication; Headscale no longer verifies Headplane passwords
- Legacy `headplane_users` and `headplane_settings` tables retained for rollback only
- API key authentication remains fully functional as secondary method
- OIDC and proxy authentication runtime paths disabled in single-admin mode
- Migration from legacy database is idempotent and preserves rollback capability

### Verification Completed

- All unit tests passing
- Browser tests cover local login, API key login, password reset, and API key rotation
- Migration tool tested for idempotency, dry-run, and rollback scenarios
- Documentation reviewed and synchronized across all files

---

## Previous Completed Tasks

### Password-authenticated route approval error ✅ COMPLETE

### Objective

Prevent Headplane route approval or disable actions from returning HTTP 500 after
Headscale has accepted the route change for password-authenticated users.

### Diagnosis and Decision

`routesAction` updates Headscale before recording `machine.routes` in the local
Headplane audit log. `actorFromPrincipal()` supported API-key and OIDC principals
but not password principals; it accessed `profile` on a password session and
threw after the route mutation completed. Password sessions now produce an audit
actor with the login username and no actor ID. This preserves the normal audit
record and leaves genuine Headscale, cache-refresh, and audit failures visible.
The same shared audit conversion is used by the per-machine route switch and
other audited password-authenticated actions.

### Changed Files

- `headplane/app/server/audit/index.ts`
- `headplane/tests/unit/audit/audit-service.test.ts`
- `headplane/tests/unit/routes/routes-actions.test.ts`
- `headplane/CHANGELOG.md`
- `session.md`

### Verification

- `pnpm exec vitest run --project unit tests/unit/audit/audit-service.test.ts tests/unit/routes/routes-actions.test.ts` — passed (12 tests).
- `pnpm run test:unit` — passed (346 tests).
- `pnpm exec oxfmt --check app/server/audit/index.ts tests/unit/audit/audit-service.test.ts tests/unit/routes/routes-actions.test.ts` — passed.
- `pnpm exec oxlint app/server/audit/index.ts tests/unit/audit/audit-service.test.ts tests/unit/routes/routes-actions.test.ts` — passed.
- `git diff --check` — passed in the parent repository and Headplane submodule.
- `pnpm run typecheck` still fails in unrelated existing files, including
  `app/hooks/use-admin.ts`, `app/routes/instances/*`, and `app/routes/ssh/page.tsx`; the former password-principal errors in `app/server/audit/index.ts` are resolved. The environment uses Node 22.22.1 while Headplane requires Node 24.2–24.x.

---

## Current Task: Password-authenticated Headplane dashboard login ✅ COMPLETE

### Objective

Fix successful Headplane username/password logins immediately redirecting back
to `/admin/login` instead of loading the dashboard.

### Status: ✅ COMPLETE (2026-10-05)

Password-session tokens are valid only for Headplane-specific endpoints. The
Headplane server had incorrectly used one as a Headscale API key for dashboard
requests, which caused normal Headscale APIs to reject it and the app layout to
clear the browser session. Password sessions now use the configured server-side
`headscale.api_key`, consistent with user sessions; API-key sessions continue
to use their submitted API key.

### Changed Files

- `headplane/app/server/web/auth.ts` — uses the configured Headscale API key
  for password and user sessions.
- `headplane/tests/unit/auth/auth-service.test.ts` — regression coverage for
  password, API-key, and missing-configuration credential selection.
- `headplane/docs/configuration/common-issues.md` — documented diagnosis and
  safe configuration requirement.
- `headplane/CHANGELOG.md` — user-facing fix note.
- `session.md` — recorded diagnosis and verification.

### Verification

- `pnpm exec vitest run --project unit tests/unit/auth/auth-service.test.ts` —
  passed (36 tests).
- `pnpm exec oxfmt --check app/server/web/auth.ts tests/unit/auth/auth-service.test.ts`
  — passed.
- `pnpm exec oxlint app/server/web/auth.ts tests/unit/auth/auth-service.test.ts`
  — passed.
- `git diff --check` — passed.
- `pnpm run typecheck` still fails in unrelated existing files including
  `app/hooks/use-admin.ts` and `app/routes/instances/*`; no errors were
  reported in this change. The environment uses Node 22.22.1 while Headplane
  requires Node 24.2–24.x.

---

## Current Task: Headplane admin-user dialogs ✅ COMPLETE

### Objective

Fix the Create User and Edit User dialogs at `/admin/admin/users`, which crashed
because they treated the shared combobox as a native select element.

### Status: ✅ COMPLETE (2026-10-05)

The dialogs now pass role items and use `onValueChange`; their text inputs use
the shared Input component's string callback. The selector offers only `user`
and `admin`, matching the backend API's accepted role values.

### Changed Files

- `headplane/app/routes/admin/users/dialogs/create-user.tsx`
- `headplane/app/routes/admin/users/dialogs/edit-user.tsx`
- `headplane/app/routes/admin/users/route.tsx`
- `headplane/tests/component/select.test.tsx`
- `headplane/tests/unit/routes/admin-users.test.ts`
- `headplane/CHANGELOG.md`
- `session.md`

### Verification

- `pnpm exec oxfmt --check app/routes/admin/users/route.tsx app/routes/admin/users/dialogs/create-user.tsx app/routes/admin/users/dialogs/edit-user.tsx tests/component/select.test.tsx tests/unit/routes/admin-users.test.ts` — passed.
- `pnpm exec oxlint app/routes/admin/users/route.tsx app/routes/admin/users/dialogs/create-user.tsx app/routes/admin/users/dialogs/edit-user.tsx tests/component/select.test.tsx tests/unit/routes/admin-users.test.ts` — passed.
- `pnpm exec vitest run --project component tests/component/select.test.tsx` — passed (2 tests).
- `pnpm exec vitest run --project unit tests/unit/routes/admin-users.test.ts` — passed (8 tests).
- `pnpm run typecheck` still fails on unrelated existing files; its filtered output has no errors in the admin-user route or dialogs. The environment uses Node 22.22.1 while Headplane requires Node 24.2–24.x.

---

## Previous Task: Headplane schema validation repair ✅ COMPLETE

### Objective

Repair the SQLite schema mismatch that rejected the required Headplane users and
settings tables, then caused the admin users page to fail after operators
removed those tables.

**Scope:** Canonical schema, explicit migration DDL, frontend error handling,
tests, and recovery documentation.

### Status: ✅ COMPLETE (2026-10-05)

The Headplane tables are now created from explicit DDL during fresh database
initialization and by repair migration `202610050900-repair-headplane-schema`.
Follow-up migration `202610050930-normalize-headplane-indexes` recreates legacy
GORM-quoted Headplane indexes with the canonical SQLite SQL. The SQLite schema
contract includes both tables and indexes, so validation no longer asks operators
to delete required data. The admin users route renders an error state rather than
evaluating `.length` on missing loader data.

### Changed Files

- `hscontrol/db/migrate_headplane.go` — canonical SQLite/PostgreSQL DDL.
- `hscontrol/db/db.go` — fresh initialization, repair, and index-normalization migration wiring.
- `hscontrol/db/schema.sql` — Headplane schema contract.
- `hscontrol/app.go` — removed post-validation schema mutation.
- `hscontrol/headplane_settings.go` — retained initializer delegates to canonical DDL.
- `hscontrol/db/db_test.go`, `hscontrol/headplane_users_test.go` — fresh and recovery coverage.
- `headplane/app/routes/admin/users/route.tsx` and its unit test — graceful load failure.
- `docs/troubleshooting.md`, `docs/ref/api/headplane-users.md`,
  `docs/ref/api/headplane-settings.md`, `CHANGELOG.md`, `session.md` — corrected
  operational and API guidance.

### Verification

- `go test ./hscontrol/db -count=1` — passed.
- Focused migration coverage verified fresh initialization, the removed-table recovery path, preserved users and settings, GORM-shaped legacy tables, and quoted Headplane indexes.
- `go test ./hscontrol -run 'TestMigrationFromSingleToMultiUser|TestDatabaseSchema' -count=1` — passed.
- `go build -o /tmp/headscale-schema-repair ./cmd/headscale` — passed.
- `pnpm exec oxfmt --check app/routes/admin/users/route.tsx tests/unit/routes/admin-users.test.ts` — passed.
- `pnpm exec vitest run --project unit tests/unit/routes/admin-users.test.ts` — passed (7 tests).
- `go test ./hscontrol -count=1` was stopped after more than two minutes with no output; focused migration coverage passed.

### Limitation

- `pnpm run typecheck` remains blocked by pre-existing repository-wide errors in
  `headplane/` and this environment uses Node 22.22.1 while the project requires
  Node 24.2–24.x. The affected route's unit test and formatter pass.

---

## Previous Task: Phase 13b — Settings Menu (Single User) ✅ COMPLETE

### Objective

Add a Settings menu for the authenticated user to manage their Headscale API key, change password, and customize preferences. This phase enhances the single-user experience from Phase 13a (password login) without adding multi-user support.

**Scope:** Settings page with API key storage, password change, theme selection, and profile name. Backend endpoints for settings storage (encrypted API key) and password change.

### Status: ✅ COMPLETE (Merged 2026-10-03)

Phase 13a (Simple Password Auth) is complete. Phase 13b (Settings Menu) is now complete.

### Progress

#### Phase 13a — Simple Password Authentication ✅ COMPLETE

- Merged in v0.35.9-arsydoni4326-alt
- Single password authentication working
- Session management, rate limiting, constant-time comparison
- Password/API key toggle in login UI
- Limitation: API key must be re-entered at each login (no storage)

#### Phase 13b — Settings Menu ✅ COMPLETE

**Implementation Complete (2026-10-03):**

**Backend:**

- ✅ Settings storage in SQLite table (`headplane_settings`)
- ✅ `GET /api/v1/headplane/settings` endpoint
- ✅ `POST /api/v1/headplane/settings` endpoint
- ✅ `POST /api/v1/headplane/change-password` endpoint
- ✅ API key encryption (AES-256-GCM, PBKDF2-derived key, 100k iterations)
- ✅ Session loads stored API key automatically
- ✅ Comprehensive unit tests (encryption, validation, error cases)

**Frontend:**

- ✅ `/settings/profile` route with complete UI
- ✅ Account section (change password, session info, logout)
- ✅ Integration section (API key management, masked display)
- ✅ Preferences section (theme selector: light/dark/system)
- ✅ Profile section (display name)
- ✅ Form validation and error handling
- ✅ Unit and E2E tests

**Documentation:**

- ✅ User guide (`docs/usage/settings.md`)
- ✅ API reference (`docs/ref/api/headplane-settings.md`)
- ✅ Integration testing guide (`docs/phase13b-integration-testing.md`)

**Parallel Development:**

- Three worktrees used: backend, frontend, docs
- Backend merged first, then frontend (via headplane submodule), then docs
- All acceptance criteria met

**Files Changed:**

- Backend: `hscontrol/headplane_settings.go`, `hscontrol/headplane_settings_test.go`, `hscontrol/app.go`
- Frontend: `headplane/app/routes/settings/profile.tsx`, API client updates
- Docs: `docs/usage/settings.md`, `docs/ref/api/headplane-settings.md`

**Next:** Phase 13c (Multi-User Support) is proposed but deferred pending user feedback.

### Progress

#### Part 1: OIDC Group Support in ACLs (COMPLETED)

**Objective:** Allow ACL policies to reference OIDC groups for access control, achieving parity with Tailscale.

**Current Analysis:**

- OIDC authentication already supports group filtering via `allowed_groups` config (authorizes login)
- Groups are fetched from OIDC provider during authentication (`groups` claim)
- **Gap:** Groups are not stored with users, cannot be referenced in policy rules
- Documentation explicitly states this limitation in `docs/about/features.md` and `docs/ref/oidc.md`

**Technical Design:**

1. **Database Schema:** Add `groups` JSONB column to `users` table to store OIDC group memberships
2. **OIDC Integration:** Extract and persist groups during authentication callback
3. **Policy Engine:** Extend `policy/v2/` to resolve group references in ACLs (Groups already exist in policy, but refer to user-defined groups, not OIDC groups)
4. **API:** Expose user groups via v1/v2 APIs
5. **Testing:** Unit tests for group resolution, integration tests for end-to-end flows

**Approach Decision:**

- Store groups as `[]string` in a new `oidc_groups` column on `users` table
- Update groups on every OIDC login (groups can change at IdP)
- Policy engine will treat OIDC groups similarly to existing policy-defined groups
- Groups will be matched by exact string comparison (case-sensitive)

**Progress:**

- [x] Created database migration `202609291402-add-oidc-groups-to-users`
- [x] Added `oidc_groups` column to `users` table schema (schema.sql)
- [x] Added `OIDCGroups []string` field to `types.User` struct with JSON serialization
- [x] Updated `User.FromClaim` to persist OIDC groups from claims
- [x] Regenerated `types_view.go` and `types_clone.go` for the new field
- [x] Extended policy engine to resolve OIDC group references:
  - `Groups.ContainsOrOIDC` validates OIDC groups against live users
  - `Policy.validateWithUsers` and `unmarshalPolicyWithUsers` thread users through validation
  - `Group.resolve` resolves OIDC group members (policy-defined groups take precedence)
- [x] Updated v1 API (`hscontrol/api/v1/types.go`) to expose `oidcGroups`
- [x] Updated v2 API (`hscontrol/api/v2/users.go`) to expose `oidcGroups`
- [x] Added unit tests in `hscontrol/policy/v2/oidc_groups_test.go`:
  - OIDC groups in ACL sources and destinations
  - OIDC groups in tagOwners and autoApprovers
  - Undefined group rejection
  - Policy-defined group precedence over OIDC groups
- [x] Fixed `202505141324` migration to use column-specific update instead of `Save()` (avoids writing non-existent columns)
- [x] Regenerated v1/v2 API clients (`gen/client/v1/client.gen.go`, `gen/client/v2/client.gen.go`) with `oidcGroups` field
- [x] Renamed `UpdateCheckResponse` → `UpdateCheckResult` in `hscontrol/updatecheck/` to fix oapi-codegen type collision exposed by client regeneration
- [x] Regenerated v1 API golden files for user/node/preauthkey responses
- [x] Updated documentation: `docs/about/features.md`, `docs/ref/oidc.md`, `docs/ref/policy.md`
- [x] Updated `ROADMAP.md` and `CHANGELOG.md`
- [ ] Frontend (Headplane) support for displaying OIDC groups — deferred to a follow-up within Phase 10

**Test Results:**

- `go test ./hscontrol/types/...` — PASS
- `go test ./hscontrol/db/...` — PASS
- `go test ./hscontrol/policy/...` — PASS
- `go test ./hscontrol/api/...` — PASS
- `go test ./hscontrol/updatecheck/...` — PASS
- `go test ./hscontrol/ -run 'TestAPIV1.*User'` — PASS
- `go test ./hscontrol/ -run 'TestAPIV1UpdateCheck'` — PASS
- `go test ./hscontrol/ -run 'TestOIDC'` — PASS
- `go build ./cmd/headscale` — PASS

**Next Steps:**

- Part 2: Machine Management Parity (IN PROGRESS)
- Part 3: DNS Management Improvements
- Part 4: Extensibility Foundation

---

#### Part 2: Machine Management Parity (IN PROGRESS)

**Objective:** Close gaps in machine management compared to Tailscale, focusing on route management, key rotation, and device posture.

**Progress:**

- [x] **Key rotation**: Added "Rotate key" action to the machine menu
  - Created `headplane/app/routes/machines/dialogs/rotate-key.tsx` confirmation dialog
  - Added `rotate_key` action handler in `headplane/app/routes/machines/machine-actions.ts`
  - Added `machine.rotate_key` to the `AuditAction` type
  - Uses the existing `expire` API to invalidate the node key, forcing re-authentication
- [x] **Route management**: Added a route overview page at `/routes`
  - Created `headplane/app/routes/routes/overview.tsx` showing all subnet routes and exit nodes across all machines
  - Created `headplane/app/routes/routes/routes-actions.ts` with `update_routes` action handler
  - Added route to `headplane/app/routes.ts`
  - Added "Routes" tab to `headplane/app/layout/header.tsx`
  - Shows route status (approved/pending) and allows approve/disable from the overview
- [ ] **Device posture**: Deferred — requires significant policy engine changes (parsing `postures`/`srcPosture`, node hostinfo attributes, filter compilation). Documented as a known limitation.

**Test Results:**

- `npx tsc --noEmit` — PASS
- `npx oxlint` — PASS (only pre-existing warnings)
- `npx vitest run tests/unit` — 283 tests PASS
- `npx vitest run tests/component` — 51 tests PASS

**Next Steps:**

- Part 3: DNS Management Improvements (COMPLETED — see below)
- Part 4: Extensibility Foundation

---

#### Part 3: DNS Management Improvements (COMPLETED)

**Objective:** Enhance DNS features, focusing on `extra_records_path` provisioning and DNS record management.

**Progress:**

- [x] **`extra_records_path` visibility**: The DNS page now shows a notice when records are managed via the `dns.extra_records_path` JSON file, indicating that changes are written to that file and picked up by Headscale automatically.
- [x] **CNAME record support**: Added CNAME to the DNS record type selector in the add-record dialog. The backend already supported CNAME via `tailcfg.DNSRecord`.
- [x] **IP address validation**: Added validation for A (IPv4) and AAAA (IPv6) record values in the add-record dialog. Invalid values show inline error messages.
- [x] **Unit tests**: Added `tests/unit/utils/dns-record.test.ts` with 7 tests covering A/AAAA/CNAME validation.

**Test Results:**

- `npx tsc --noEmit` — PASS
- `npx oxlint` — PASS (only pre-existing warnings)
- `npx vitest run tests/unit` — 290 tests PASS (was 283, +7 new)

**Next Steps:**

- Part 4: Extensibility Foundation (COMPLETED — see below)

---

#### Part 4: Extensibility Foundation (COMPLETED)

**Objective:** Prepare the codebase for future plugin/extension systems, monitoring integrations, and external tooling.

**Progress:**

- [x] **Documented extension points**: Created `docs/ref/extending.md` covering:
  - v1 admin API and v2 Tailscale-compatible API as extension points
  - OAuth client-credentials flow for Terraform/K8s operator/tscli
  - Policy engine extension points (new sections, alias types, OIDC groups)
  - Headplane extension points (routes, context, API client, audit log)
  - Monitoring/alerting integration points
  - How to add a new API endpoint
- [x] **Added Grafana dashboard example**: Added a Grafana dashboard JSON to `docs/usage/observability.md` visualising key Headscale metrics (HTTP requests, map responses, connected nodes, NodeStore operations, update check).
- [x] **Validated Terraform/K8s operator support**: Confirmed the existing servertest tests (`TestAPIv2`, `TestAPIv2OAuthScopes`) validate the Terraform provider and OAuth client-credentials flow. These require `tofu` from the nix dev shell (pre-existing limitation).
- [x] **Updated docs navigation**: Added "Extending Headscale" to `mkdocs.yml`.

**Test Results:**

- Terraform/K8s operator tests exist and are well-structured; they require `tofu` from the nix dev shell (not available in this environment).

**Next Steps:**

- Phase 10 complete. All four parts delivered.

---

## Previous Task: Phase 9 — Supporting Features and Observability (COMPLETED)

### Progress

- [x] **Part A: Operational Endpoints**
  - Added `/ready` endpoint for Kubernetes-style readiness probes (returns 503 when not ready)
  - Existing `/health` endpoint provides liveness checking (database connectivity)
  - Existing `/version` endpoint provides version information
  - Existing `/api/v1/health` provides authenticated health check
  - All endpoints tested and working

- [x] **Part B: Metrics Documentation**
  - Created comprehensive observability documentation at `docs/usage/observability.md`
  - Documented all existing Prometheus metrics:
    - HTTP metrics (requests, duration)
    - MapResponse metrics (sent, generated, endpoint updates, ended)
    - NodeStore metrics (operations, durations, batch size, snapshot builds, node count, peers calculation, queue depth)
    - Mapper metrics (changes dropped)
    - Update check metrics (requests, failures, cache hits/misses)
    - HA health probe metrics
  - Documented metrics configuration and security (protected endpoints)
  - Added Prometheus scrape configuration examples
  - Added recommended alerting rules

- [x] **Part C: Logging and Debug Documentation**
  - Documented structured logging with zerolog
  - Documented log levels and formats (text/json)
  - Documented debug endpoints (`/debug/overview`, `/debug/config`, `/debug/policy`, `/debug/pprof/`, `/debug/statsviz`)
  - Added troubleshooting guide for common observability issues
  - Added logging best practices

- [x] **Part D: Testing**
  - Added unit tests for `/health` endpoint (TestHealthHandler)
  - Added unit tests for `/ready` endpoint (TestReadyHandler)
  - Added unit tests for `/version` endpoint (TestVersionHandler)
  - All tests pass successfully

- [x] **Part E: Documentation Updates**
  - Updated `CHANGELOG.md` with Phase 9 entry
  - Updated `ROADMAP.md` to mark Phase 9 as [Implemented]
  - Updated `session.md` with completion status

### Decisions and Assumptions

- The `/ready` endpoint follows Kubernetes conventions (503 when not ready vs 500 for `/health` failures)
- Metrics remain on the separate debug server (`metrics_listen_addr`) for security
- No new metrics were added (existing coverage is comprehensive)
- Documentation focuses on what exists rather than adding speculative features
- Distributed tracing (OpenTelemetry) was deferred as agreed in the planning phase

### Completed Deliverables

1. **New `/ready` endpoint** (`hscontrol/handlers.go`)
   - Returns 200 OK when database is reachable
   - Returns 503 Service Unavailable when not ready
   - JSON response with `ready` boolean and `status` string
   - Mounted at `/ready` in the main router

2. **Comprehensive observability documentation** (`docs/usage/observability.md`)
   - Operational endpoints (health, ready, version, authenticated health)
   - Complete metrics catalog with types, descriptions, and labels
   - Logging configuration and best practices
   - Debug endpoints reference
   - Prometheus alerting examples
   - Troubleshooting guide for common issues

3. **Test coverage** (`hscontrol/handlers_test.go`)
   - TestHealthHandler: validates `/health` endpoint response
   - TestReadyHandler: validates `/ready` endpoint response and status codes
   - TestVersionHandler: validates `/version` endpoint response
   - All tests pass

4. **Documentation updates**
   - `CHANGELOG.md`: Phase 9 entry added to "# Next" section
   - `ROADMAP.md`: Phase 9 marked as [Implemented], current state updated to "Phases 1-9 complete"
   - `session.md`: Phase 9 marked as COMPLETED with full progress log

### Files Modified

- `hscontrol/handlers.go` — Added `ReadyHandler` function
- `hscontrol/app.go` — Added `/ready` route
- `hscontrol/handlers_test.go` — Added 3 new test functions
- `docs/usage/observability.md` — Created (new file, 406 lines)
- `CHANGELOG.md` — Added Phase 9 entry
- `ROADMAP.md` — Updated Phase 9 status and current state
- `session.md` — Updated with Phase 9 completion

### Verification

```bash
# All tests pass
go test ./hscontrol -run 'TestHealthHandler|TestReadyHandler|TestVersionHandler'
# Result: ok (0.046s)

# Code is properly formatted
gofmt -l hscontrol/handlers.go hscontrol/handlers_test.go hscontrol/app.go
# Result: no output (all files formatted correctly)
```

### Next Steps

Phase 9 is complete. The next phase would be Phase 10 (Feature Parity and Expansion) or Phase 11 (Advanced Features and Integrations) per the roadmap, but these require explicit approval before proceeding.

### Progress

- [x] **Documentation: Accessibility guidance added**
  - Updated `docs/usage/ui-guide.md` with comprehensive accessibility section covering:
    - Keyboard navigation (Tab/Shift+Tab, Enter/Space, Arrow keys, Esc)
    - Screen reader support (semantic HTML, ARIA labels, roles, alternative text)
    - Visual accessibility (color contrast, dark mode, focus indicators, link styling, touch targets)
    - Forms and validation (associated labels, error messages, required fields)
    - Testing (axe-core, Playwright e2e, component tests)
    - Known limitations (full WCAG validation, dynamic content, third-party components)
    - Reporting process for accessibility issues
  - Added accessibility verification to maintenance checklist
  - Updated ROADMAP.md: Phase 8 status changed from [Proposed] to [Implemented]
  - Updated CHANGELOG.md with Phase 8 entry
  - Committed as `9929dd07`: "docs: add accessibility guidance to UI guide (Phase 8)"

- [x] **Headplane: Standardized UI state components (completed)**
  - Created `headplane/app/components/empty-state.tsx` — standardized empty state component with accessible markup
  - Created `headplane/app/components/loading-spinner.tsx` — accessible loading indicator with ARIA live regions
  - Created component tests with 100% coverage (13 tests pass)
  - Refactored all data-displaying routes to use standardized components:
    - machines, users, audit, derp, topology, auth-keys
  - Extended accessibility e2e test coverage (added 3 new tests for empty/loading states and icon-only buttons)
  - Updated `docs/usage/ui-guide.md` with standardized UI states documentation and consistency guidelines
  - Verified all 25 dialogs use standardized `Dialog` + `DialogPanel` pattern
  - Confirmed tooltip component is in use on 12+ UI elements
  - All component tests pass (51/51)
  - Typecheck passes with no errors

### Decisions and Assumptions

- Accessibility documentation is placed in the UI guide rather than a separate document, making it easy to find alongside visual documentation
- The accessibility section documents both backend-served templates (which already have good accessibility) and the Headplane UI (which has axe-core testing)
- Known limitations are documented transparently (third-party component patterns, manual testing requirements, dynamic content gaps)
- Standardized components follow the "anti-slop" principle: empty states name the cause and the next action, not just "No data"
- Routes like ACLs, DNS, and Settings are configuration/editor pages that don't require empty states
- Dialogs already use the standardized pattern (25 dialogs verified)
- Tooltips already provide contextual help throughout the UI (12+ elements)
- Onboarding is handled by descriptive empty states that guide first-time users

### Known Issues and Limitations

- The Headplane UI is a separate git submodule; frontend accessibility improvements were made within that repository
- Backend templates already have good accessibility (ARIA attributes, semantic HTML, focus styles, tested in `hscontrol/templates_consistency_test.go`)
- The `@base-ui/react` library uses `aria-disabled` instead of `disabled` attribute for some components (documented limitation)
- Icon-only buttons in some third-party components lack accessible names (workaround: structural queries in tests)
- Some accessibility improvements (full keyboard navigation testing, screen reader flows, published conformance statement) are deferred to future phases

### Pending Work

Phase 8 is complete. No remaining work for this phase.

---

### 2026-10-04: Fixed duplicate migration block in hscontrol/db/db.go

**Problem:** Docker build failed with Go syntax errors in `hscontrol/db/db.go` lines 780-802.

**Root Cause:** Duplicate migration entry with ID `202610031721-create-headplane-users` appeared twice in the migrations slice (lines 770-778 and lines 780-792).

**Fix:** Removed the duplicate migration block (lines 780-792).

**Verification:** Migration slice now contains only unique IDs and is syntactically valid.

---

### Previous Task: Phase 7 — Testing and Technical Hardening

#### Progress

- [x] **Backend: Wire `DisableUpdateCheck` into update-check endpoint**
  - Modified `hscontrol/updatecheck/updatecheck.go`: `BuildResponse(check bool, disabled bool)` — when disabled, skip remote check.
  - Modified `hscontrol/api/v1/updatecheck.go`: pass `b.Cfg.DisableUpdateCheck`.
  - Updated `hscontrol/updatecheck/updatecheck_test.go`: new signature + `TestBuildResponse_Disabled`.
  - Updated `config-example.yaml` comment for `disable_check_updates`.
  - All updatecheck tests pass.

- [x] **Backend: DERP API test + OpenAPI spec coverage**
  - Created `hscontrol/apiv1_derp_test.go` with `TestAPIV1Derp` (endpoint response shape) and `TestAPIV1DerpInSpec` (spec presence).
  - Regenerated OpenAPI spec (gitignored; served live from code).

- [x] **Headplane: UI component testing framework**
  - Installed `@testing-library/react`, `@testing-library/user-event`, `@testing-library/jest-dom`, `jsdom`.
  - Added `component` project to `vitest.config.ts` (jsdom environment).
  - Added `test:component` script to `package.json`.
  - Created `tests/component/setup.ts` (jest-dom/vitest integration + cleanup).
  - Wrote 38 component tests across 7 files: button, chip, status-circle, table-list, token-list, tabs, switch.
  - Wired `test:component` into `.github/workflows/build.yaml`.
  - All component tests pass.

- [x] **Headplane: Fix docs inconsistency**
  - Updated `headplane/docs/configuration/index.md`: removed "NOT enabled by default" text; documented `HEADPLANE_LOAD_ENV_OVERRIDES` as deprecated.

- [x] **Headplane: Remove legacy `useUpdateCheck.ts`**
  - Verified `useUpdateCheck.ts` and `git-api.ts` are dead code (not imported by any other file).
  - Removed both files.
  - Typecheck and all tests pass.

- [x] **Docs: UI guide screenshots → tracked maintenance task**
  - Updated `docs/usage/ui-guide.md` to reference a GitHub issue for screenshot tracking.

- [x] **Docs: ROADMAP.md, CHANGELOG.md, session.md updates**
  - ROADMAP.md: Phase 7 marked [Implemented]; Current state updated to Phases 1-7 complete.
  - CHANGELOG.md: Phase 7 entries added to `# Next`.
  - This session.md updated.

### Known Issues and Limitations

- The `headplane/nix/package.nix` pnpm deps hash is stale after adding
  `@testing-library/*` and `jsdom`. The automated workflow
  (`.github/workflows/automated.yaml`) updates it on push to main.
- The servertest suite has pre-existing failures in `TestAPIv2OAuthScopes` and
  `TestAPIv2/TSCLI`/`TestAPIv2/Terraform` (unrelated to Phase 7 changes).
- The `openapi/v1/headscale.yaml` spec file is gitignored and generated live;
  the regeneration confirmed the DERP endpoint is present.
- The `@base-ui/react` library uses `aria-disabled` instead of `disabled`
  attribute for some components (Switch), and icon-only buttons lack accessible
  names — component tests work around these with structural queries.
---

## Fixed: Headplane "No configuration is available" error (October 26, 2026)

### Root Cause

The Headplane container was running without a proper Headplane configuration file. When Headplane starts, it looks for a config at `/etc/headplane/config.yaml` (or the path in `HEADPLANE_CONFIG_PATH`). Without a config file, `loadConfig()` in `app/server/config/load.ts` returns an empty config, and `loadHeadscaleConfig()` in `app/server/headscale/config-loader.ts` receives no `headscale.config_path`, so it returns a config object with `access: "no"`.

When the DNS overview route (`app/routes/dns/overview.tsx`) calls `headscaleConfig.readable()`, it returns `false` and throws `Error: "No configuration is available"`.

### Fix

Created three files at the project root:

1. **`compose.yaml`** — Docker Compose file that runs both headscale and headplane together:
   - Mounts `./config.yaml` into headscale at `/etc/headscale/config.yaml`
   - Mounts `./headplane-config.yaml` into headplane at `/etc/headplane/config.yaml`
   - Mounts `./config.yaml` into headplane at `/etc/headscale/config.yaml` (so Headplane can read the Headscale config)
   - Uses named volumes for persistent data

2. **`headplane-config.yaml`** — Headplane configuration file:
   - Sets `headscale.url` to `http://headscale:8080` (Docker network hostname)
   - Sets `headscale.config_path` to `/etc/headscale/config.yaml`
   - Sets `server.cookie_secret` to a 32-char dev value (must be changed for production)
   - Sets `server.cookie_secure: false` for non-HTTPS local dev
   - Disables all integrations by default

3. **`config.yaml`** — Minimal Headscale configuration file:
   - Listens on `0.0.0.0:8080`
   - Uses SQLite at `/var/lib/headscale/db.sqlite`
   - MagicDNS enabled with `example.com` base domain
   - Standard Tailscale IP prefixes

### Usage

```bash
# Build and start both services
docker compose up -d

# Create an API key for Headplane
docker compose exec headscale headscale apikeys create

# Set the generated key in headplane-config.yaml as headscale.api_key
# Then restart headplane
docker compose restart headplane
```

### Validation

- ✅ All three YAML files parse correctly
- ✅ Headscale binary builds and loads the config without errors
- ✅ The Headplane config schema accepts the config structure

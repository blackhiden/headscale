# Testing and CI

This page describes how the fork is tested and what the continuous integration
(CI) pipeline enforces. It is aimed at contributors and operators who want to
understand the quality gates or run the test suites locally.

## Test suites

### Backend (Go)

The Headscale control server is tested with the standard Go toolchain:

```bash
# Unit tests (the pure subset CI runs)
CGO_ENABLED=0 go test ./...

# With the race detector (slower)
go test -race ./...

# Server-level tests (hscontrol/servertest) — slow, timing-sensitive
go test -timeout=20m ./hscontrol/servertest/...
```

The `hscontrol/updatecheck/` package (the fork's update-check feature) has
deterministic edge-case tests that point the GitHub API base URL at a local
`httptest` server, covering non-200 responses, empty/malformed payloads, cache
behavior, repo-config environment handling, and release/commit comparison
fallbacks.

### Frontend (Headplane)

Headplane has three layers of tests:

| Suite | Command | What it covers |
| --- | --- | --- |
| Unit | `pnpm run test:unit` | Service logic, config loading, utilities |
| Integration | `pnpm run test:integration` | API flows against real Headscale containers |
| E2E + a11y | `pnpm run test:e2e` | Browser flows (login, machines, ACL, DNS) and axe-core accessibility |

The Playwright e2e suite (`headplane/tests/e2e/`) starts a real Headscale
container via testcontainers, writes a throwaway Headplane config, and drives
the dev server in a browser. The accessibility tests (`a11y.spec.ts`) fail on
critical or serious axe-core violations across the main routes.

## Performance tooling

- **Bundle-size analysis** — `pnpm run analyze:bundle` (in `headplane/`) reads
  the Vite build output and fails when any JS chunk exceeds 1 MB (gzipped) or
  the total exceeds 2.5 MB. The per-chunk budget is generous because the lazy
  SSH route chunk includes the `restty` terminal emulator; tighten both after
  capturing a baseline. Budgets are overridable via
  `HEADPLANE_BUNDLE_CHUNK_BUDGET` / `HEADPLANE_BUNDLE_TOTAL_BUDGET`.
- **Lighthouse CI** — `pnpm run lighthouse:ci` (in `headplane/`) measures the
  login page against performance, accessibility, best-practices, and SEO
  budgets. The budgets are initial values; tighten them after capturing a
  baseline.
- **WASM SSH payload** — the `hp_ssh.wasm` module is excluded from the client
  bundle and fetched at runtime only when the SSH page is opened.

## CI gates

The fork's CI runs on GitHub Actions. The load-bearing gates are:

| Workflow | Enforces |
| --- | --- |
| `nix-checks.yml` | Build, Go unit tests, golangci-lint, formatting (via `nix build .#checks.*`) |
| `servertest.yml` | `hscontrol/servertest` (race, stress, HA property tests) |
| `build.yml` | Nix build + cross-compilation, vendor hash freshness |
| `check-generated.yml` | `make generate` output is committed |
| `check-tests.yaml` | Integration test workflow is generated and up to date |
| `docs-test.yml` | `mkdocs build --strict` succeeds |
| `fork-features.yml` | Fork-specific features survive upstream merges |
| `test-integration.yaml` | Docker-based end-to-end integration tests |
| Headplane `build.yaml` | Frontend build, unit/integration tests, bundle-size analysis, Playwright e2e + a11y, Lighthouse CI |

### Fork-feature gate

`fork-features.yml` fails the build if any of the fork's load-bearing features
are removed by an upstream merge:

- `hscontrol/updatecheck/` (backend update check)
- `hscontrol/api/v1/updatecheck.go` (update-check API registration)
- `hscontrol/api/v1/derp.go` (DERP status endpoint)
- `headplane/app/update-check/` (frontend update check)

## Local development

The full local workflow is `make dev` (format + lint + test + build) in the
repository root, or `nix develop` for the pinned toolchain. For Headplane,
`pnpm run test:unit` and `pnpm run test:e2e` cover the frontend; the e2e suite
requires Docker (for the Headscale container) and Playwright browsers
(`pnpm run test:e2e:install`).
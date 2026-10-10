# Architecture

This document describes the system architecture of the `arsydoni4326-alt` fork
of Headscale and its frontend, Headplane: the major components, their
responsibilities, how data flows between them, and the key design decisions and
trade-offs.

> For requirements, see [Specification](./SPECIFICATION.md). For the high-level
> direction, see [Roadmap](./ROADMAP.md). Headplane's own architecture
> documentation lives in [`headplane/docs/`](./headplane/docs/).

## System Overview

The fork consists of two main components:

1. **Headscale** — the Go control server (this repository's `hscontrol/`).
2. **Headplane** — the TypeScript/React web UI (the `headplane/` submodule).

```
┌────────────────────────────────────────────────────────────┐
│                        Browser                              │
│  ┌──────────────────┐  ┌────────────────────────────────┐  │
│  │  Headplane UI    │  │  Browser SSH (WASM, tsconnect) │  │
│  │  (React Router 7)│  │  connects directly to Tailnet  │  │
│  └────────┬─────────┘  └───────────────┬────────────────┘  │
└───────────┼────────────────────────────┼───────────────────┘
            │ HTTP                       │ WireGuard/DERP
┌───────────▼────────────────────────────▼───────────────────┐
│                Headplane server (Node.js)                  │
│  ┌─────────────┐ ┌──────────┐ ┌─────────┐ ┌─────────────┐  │
│  │ OIDC/auth   │ │ Headscale│ │  SQLite │ │ Headscale   │  │
│  │ services    │ │ API client│ │ (drizzle)│ │ config I/O  │  │
│  └─────────────┘ └──────────┘ └─────────┘ └─────────────┘  │
└───────────┬─────────────────────────────────────┬──────────┘
            │ REST API                            │ Docker/K8s exec
┌───────────▼───────────┐             ┌───────────▼────────────┐
│      Headscale        │             │    Headplane Agent (Go) │
│   control server      │◄────────────│  joins Tailnet, reports │
│  (separate process)   │  Tailnet    │  node details           │
└───────────┬───────────┘             └────────────────────────┘
            │
      ┌─────▼─────┐
      │  Tailnet  │
      └───────────┘
```

## Components

### Headscale control server (`hscontrol/`)

The Go control plane. Top-level files (`app.go`, `handlers.go`, `grpcv1.go`,
`noise.go`, `auth.go`, `oidc.go`, `poll.go`, `metrics.go`, `debug.go`) wire the
server together. Key subsystems:

- **`state/`** — the central coordinator (`state.go`) and the copy-on-write
  `NodeStore` (`node_store.go`). All cross-subsystem operations go through
  `State`.
- **`db/`** — GORM layer, migrations, schema.
- **`mapper/`** — streaming batcher that distributes MapResponses to clients.
- **`policy/`** — `policy/v2/` is the policy implementation.
- **`api/v1/`** and **`api/v2/`** — code-first Huma implementations of the
  REST APIs; Huma emits the OpenAPI 3.1 documents from the Go definitions.
- **`updatecheck/`** — the fork-specific update-check package (see below).
- **`api/v1/derp.go`** — the read-only `GET /api/v1/derp` endpoint added in
  Phase 4, which returns the current DERP relay map (regions, nodes, ports,
  IPs) from `State.DERPMap()`.

### Headplane web UI (`headplane/`)

A React Router 7 (framework mode) application built with Vite. It talks to
Headscale over its REST API and provides machine management, ACL editing, DNS
settings, and browser SSH. See
[`headplane/docs/ARCHITECTURE.md`](./headplane/docs/ARCHITECTURE.md) for
details.

Phase 4 added the following Headplane-side subsystems:

- **Audit log** (`headplane/app/server/audit/`) — records changes made through
  the UI into a local SQLite `audit_log` table (drizzle). Action handlers
  (machines, ACL, DNS, auth keys, users) call `audit.record()` after a
  successful mutation; the `/audit` page lists entries with filters and
  pagination.
- **Topology view** (`headplane/app/utils/topology.ts` +
  `headplane/app/components/topology/`) — a deterministic SVG layout of the
  tailnet grouped by owner, with subnet routers, exit-node highlighting,
  pan/zoom, and click-through to machine details.
- **Compatibility registry** (`headplane/app/server/headscale/compatibility.ts`)
  — maps Headplane features to the minimum Headscale version they require.
  The settings page and feature pages surface unsupported features.
- **DERP status page** (`headplane/app/routes/derp/`) — reads the
  `GET /api/v1/derp` endpoint and renders region/node cards.
- **Export/import** (`headplane/app/routes/settings/export/`) — downloads the
  Headscale configuration (YAML or JSON bundle) and the ACL policy, and
  accepts an uploaded bundle for import.

### Update check (`hscontrol/updatecheck/` + `headplane/app/update-check/`)

The fork-specific update checker. The backend package
(`hscontrol/updatecheck/`) is self-contained and marked "DO NOT REMOVE": all
update-check logic lives in this single package so it can be restored as a unit
if removed by an upstream merge. The frontend domain
(`headplane/app/update-check/`) mirrors this design on the Headplane side.

The backend features:

- **Server-side caching** — GitHub API responses are cached for 15 minutes
  (mutex-protected with singleflight behavior) to stay within unauthenticated
  rate limits.
- **Release version comparison** — For tagged builds, the endpoint fetches the
  latest release tag and compares via semver precedence. Falls back to commit
  hash comparison for dev builds.
- **Configurable repository** — The remote repo is set via
  `HEADSCALE_UPDATE_CHECK_REPO` env var (format: `owner/repo`), defaulting to
  `arsydoni4326-alt/headscale`.
- **Prometheus metrics** — Request count, remote-failure count, cache hit/miss.

The frontend features:

- **sessionStorage caching** — Results are cached for 15 minutes; manual checks
  bypass the cache.
- **Dismiss/remind controls** — "Dismiss for this session" and "Remind me
  later (24h)" buttons on the update modal.
- **Release notes link** — When a tag-based update is detected, a direct link
  to the release notes is shown.

## Data Flow

### Update check

1. The Headplane frontend calls its own backend at `/api/update-check`, which
   forwards the query parameters to Headscale's `/api/v1/update-check`.
2. Headscale returns the running binary's version info. With `?check=true` it
   fetches the latest release tag (for release builds) or commit hash (for dev
   builds) from the configured GitHub repository via the GitHub API. Results
   are cached for 15 minutes.
3. The frontend compares the remote version/commit with the local version and
   shows a modal when an update is available. Results are cached in
   `sessionStorage` for 15 minutes. Users can dismiss the modal for the session
   or set a 24-hour reminder.

### Node registration

1. Noise handshake (`noise.go`) → auth (`auth.go`) → state/DB persistence
   (`state/`, `db/`) → initial map (`mapper/`).
2. For pending interactive registrations, `/register/{auth_id}` renders the CLI
   command and a QR code containing the auth ID, public server URL, and the
   fixed registration-cache deadline. The same QR is available from the OIDC
   confirmation page.
3. Headplane validates the QR payload before registering it to the selected
   user. Headscale remains authoritative: the pending auth ID must still exist
   and is consumed by the successful registration.

## Key Design Decisions

| Decision                                  | Rationale                                                                        | Trade-off                                                                 |
| ----------------------------------------- | -------------------------------------------------------------------------------- | ------------------------------------------------------------------------- |
| Code-first Huma APIs (`api/v1`, `api/v2`) | OpenAPI spec is generated from the handlers, so it cannot drift                  | Spec emission requires building the API with a zero Backend               |
| Self-contained update-check package       | An upstream merge cannot silently delete the feature; it is restorable as a unit | The package is fork-specific and must be re-applied after upstream merges |
| Headplane as a submodule                  | The frontend is versioned and released independently                             | CI must check out submodules to verify frontend features                  |
| `-arsydoni4326-alt` version suffix        | Identifies fork releases unambiguously                                           | Version numbers diverge from upstream                                     |

## Constraints and Dependencies

- Go 1.27 for the backend (per `go.mod`).
- Node >= 24.2, PNPM >= 10.4 for Headplane.
- The `headplane/` submodule is pinned to a fork release
  (`v0.7.6-arsydoni4326-alt`).
- The fork repository targets are `github.com/arsydoni4326-alt/headscale` and
  `github.com/arsydoni4326-alt/headplane`.

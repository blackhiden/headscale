# Phase 13c: Single Local Administrator Migration Plan

**Status:** Complete  
**Completed:** 2026-10-06  
**Authoritative scope:** [Phase 13c in the roadmap](../ROADMAP.md#phase-13c--single-local-administrator-migration-planned)  
**Companion operator procedure:** [Migration runbook](./phase13c-single-admin-migration-runbook.md)

## Completion Summary

Phase 13c — Single Local Administrator Migration is **complete**. Headplane now supports only a single local administrator account with password stored as a bcrypt hash in configuration. The Headscale database-backed multi-user local authentication system has been retired. All implementation, migration, testing, and documentation steps outlined in this plan have been completed successfully.

**Key outcomes:**
- Single local administrator configured via `config.yaml` with bcrypt password hash
- Headplane owns local authentication; Headscale no longer manages Headplane users
- Legacy `headplane_users` and `headplane_settings` tables preserved for rollback
- API key authentication remains fully functional as a secondary method
- OIDC and proxy authentication runtime paths disabled in single-admin mode
- Migration tooling, password reset, and API key lifecycle UI implemented
- Comprehensive testing and documentation delivered

See the [CHANGELOG](../CHANGELOG.md) for detailed changes.

---

## Purpose

Phase 13c replaces the current Headscale database-backed local multi-user
implementation with one local Headplane administrator. This plan covers the
target architecture, implementation order, migration invariants, security model,
and verification requirements.

## Approved decisions

- Local username/password login is the primary interactive authentication method.
- Headscale API-key login remains an independent secondary authentication method.
- Headplane owns the one local administrator credential and verifies passwords;
  Headscale must no longer do so.
- `user.username` and bcrypt-hashed `user.password` live in Headplane's own
  `config.yaml`, normally `/etc/headplane/config.yaml`.
- The local administrator has complete Headplane administrative capabilities.
- Local account creation, role assignment, invitations, deletion, and per-user
  local preferences are unsupported after migration.
- Existing OIDC and proxy-auth implementation code remains in the repository but
  is disabled at runtime in this mode. This release exposes no switch to
  re-enable either path.
- Existing Headscale `headplane_users` and `headplane_settings` records are
  retained unchanged through the rollback window. They become migration input
  and rollback data, not runtime authentication state.

This change does not alter Headscale tailnet users, nodes, policies, normal
Headscale API keys, or the fork's unrelated features.

## Existing and target architecture

### Existing local-auth flow

```text
Browser -> Headplane login form
        -> POST Headscale /api/v1/headplane/login
        -> Headscale headplane_users table
        -> bcrypt verification and Headscale session
        -> Headplane cookie containing the Headscale password token
```

The existing implementation also persists local-user settings and encrypted
API-key material in `headplane_settings`. `/admin/users` currently calls
Headscale-local user CRUD endpoints.

### Target local-auth flow

```text
Browser -> Headplane login form
        -> Headplane local-admin authentication service
        -> bcrypt verification of config user.password
        -> Headplane signed session cookie and auth_sessions record
        -> configured headscale.api_key for Headscale REST calls
        -> Headscale REST API
```

The password never crosses the Headplane-to-Headscale connection. Headscale no
longer creates, verifies, or manages Headplane-local accounts or sessions.

## Target configuration contract

```yaml
user:
  username: admin
  # Bcrypt verification hash at cost 12. Never plaintext.
  password: "$2b$12$..."

headscale:
  url: "https://headscale.example.com"
  # Prefer a deployment-managed secret file.
  api_key_path: "${CREDENTIALS_DIRECTORY}/headscale-api-key"
```

Headplane startup must reject an empty username, an empty password, plaintext
password material, malformed bcrypt input, and unsupported bcrypt variants.

`headscale.api_key` or `headscale.api_key_path` is a separate Headscale service
credential. A password-authenticated local administrator uses it to access
Headscale data. It must not be derived from, stored in, or returned through the
local password session.

## Functional requirements

1. Password login accepts only the configured username and a password verified
   against the configured bcrypt hash at cost 12.
2. Password login creates a Headplane-owned session using the existing signed
   cookie/session architecture.
3. Failed local password attempts are rate limited by client address. Unknown
   usernames and incorrect passwords receive the same response.
4. API-key login remains available and continues to validate directly against
   Headscale.
5. `/admin/users` remains a stable route but becomes an Administration page.
6. Administration supports password reset and Headscale API-key lifecycle
   management, not local-user CRUD.
7. The migration command copies one valid legacy administrator bcrypt hash from
   a SQLite database without modifying the database.

## Explicit non-goals

- Recovering, displaying, logging, or writing a plaintext legacy password.
- Migrating secondary local accounts, roles, themes, profile names, or encrypted
  `headplane_settings` API-key data.
- Replacing local multi-user storage with another multi-user schema.
- Removing OIDC/proxy implementation code or unrelated Headplane identity data.
- Dropping, renaming, or mutating legacy Headscale local-auth tables in the
  initial release.
- Adding application-level `config.yaml` encryption. The deployment would still
  need a decryption key at startup, so this only moves the secret-management
  problem.

## Implementation plan

### 1. Add the Headplane configuration and bcrypt primitive

In the Headplane submodule:

- Update `app/server/config/config-schema.ts` to define `user.username` and
  `user.password` in both the full and partial configuration schema.
- Preserve the existing `HEADPLANE_*` override convention:
  `HEADPLANE_USER__USERNAME` and `HEADPLANE_USER__PASSWORD`.
- Validate `user.password` as a bcrypt verification hash and redact it from
  configuration error/logging paths.
- Add one maintained Node 24-compatible bcrypt dependency in `package.json`,
  regenerate `pnpm-lock.yaml`, and update the Nix pnpm dependency hash.
- Use cost 12, matching the legacy Go `BcryptCost` value.

No custom password hashing or comparison routine is allowed. Password comparison
must use the bcrypt library verification operation.

### 2. Move local authentication to Headplane

Add a focused closure-factory service under `app/server/web/`, for example
`local-admin.ts`, with explicit dependencies for validated user config, bcrypt,
session creation, clock, and login-attempt tracking.

The service must:

- compare against the one configured username;
- verify the submitted password with bcrypt;
- limit failed attempts to at least five per client address per minute;
- create a Headplane local password session;
- revoke all local password sessions after password reset; and
- avoid logging passwords, bcrypt hashes, full tokens, or API-key secrets.

Extend `app/server/web/auth.ts` instead of inventing a separate cookie format.
A password principal should contain the configured username and Headplane session
identifier, not a Headscale-issued password token. Update `app/server/context.ts`
so password sessions use the configured Headscale API key for dashboard calls.

When no service key is configured, fail with a clear configuration error rather
than allowing a password login that cannot load dashboard data.

### 3. Replace the login integration

Update:

- `app/routes/auth/login/action.ts` to use the Headplane local-admin service
  instead of `headscale.passwordLogin(username, password)`;
- `app/routes/auth/login/page.tsx` to retain only Password and API Key choices;
- `app/server/context.ts` to disable OIDC/proxy runtime construction in this
  mode; and
- `app/server/headscale/api/index.ts` and `transport.ts` to remove the
  Headscale-local password-login client contract after its callers are removed.

OIDC/proxy modules and their existing persistent data remain in source. The
target is runtime disablement, not a broad unrelated refactor.

### 4. Retire the Headscale-local runtime surface

In the parent repository, retire the fork-specific routes currently mounted in
`hscontrol/app.go`:

- `POST /api/v1/headplane/login`
- `GET/POST /api/v1/headplane/settings`
- `POST /api/v1/headplane/change-password`
- `GET/POST/PUT/DELETE /api/v1/headplane/users`

Remove their active initialization and call paths only after Headplane no longer
depends on them. Keep normal Headscale API-key routes and CLI behavior unchanged.

Fresh Headscale installations must not create obsolete local-account tables.
Existing `headplane_users` and `headplane_settings` tables must be preserved
without a destructive migration so a rollback remains possible.

### 5. Replace user management with Administration

Keep `/admin/users` for bookmarks, but replace the user-management route and
dialogs with one Administration page that includes:

1. a permanent statement that only one local administrator is supported;
2. current-password/new-password/confirmation reset controls;
3. Headscale API-key list, creation, rotation, revocation, and deletion controls;
4. clear operator recovery guidance when the configuration source is immutable;
   and
5. no create/edit/delete local-user or role-management controls.

For a UI password change, Headplane must verify the current password, validate
the replacement, bcrypt-hash at cost 12, atomically update the config file,
invalidate all password sessions, and issue a fresh session to the current
administrator. It must not write environment variables, a secret-manager mount,
or an externally managed `*_path` file.

### 6. Add a shared safe config-rewrite helper

Both password reset and migration need a single configuration mutation helper.
It must:

1. resolve the active config path using the existing config path rules;
2. parse and preserve unrelated YAML fields;
3. change only requested fields;
4. create a timestamped same-directory backup before the first mutation;
5. create a same-directory temporary file with mode `0600`;
6. preserve restrictive original ownership and mode when permissions allow;
7. write, `fsync`, and close the temporary file;
8. atomically rename it over the existing config; and
9. `fsync` the parent directory where supported.

If anything fails before the rename, the existing config remains valid. The
caller reports an actionable error and retains the current session.

## API-key lifecycle UI

Extend Headplane's current list-only resource in
`app/server/headscale/api/resources/api-keys.ts` with these typed operations:

| Operation | Headscale endpoint | Required UI behavior |
| --- | --- | --- |
| List | `GET /api/v1/apikey` | Show ID, masked prefix, created time, expiry, and state. |
| Create | `POST /api/v1/apikey` | Require expiry and display the plaintext key exactly once. |
| Expire/revoke | `POST /api/v1/apikey/expire` | Require confirmation naming the masked prefix. |
| Delete | `DELETE /api/v1/apikey/{prefix}` | Require destructive-action confirmation. |

Never store a created plaintext API key in browser storage, query parameters,
audit details, logs, or a local database. List responses do not contain it.

### Service-key safeguards

The API key configured through `headscale.api_key` or `headscale.api_key_path`
is Headplane's service key. Identify it by matching its secret to a listed masked
prefix without rendering the secret. Direct expiry/deletion is blocked.

For direct writable `headscale.api_key`, rotation is ordered as follows:

1. create a replacement key;
2. display it once and require the administrator to save it safely;
3. verify it with a harmless authenticated Headscale request;
4. atomically update the config;
5. validate the reloaded replacement configuration; then
6. offer, but do not automatically perform, expiry of the prior key.

For a path-backed or externally managed secret, Headplane may create the
replacement but must direct the operator to update the secret provider, restart
or reload Headplane, validate access, and only then expire the prior key.

Revoking the API key used by the current API-key browser session immediately
destroys that session and redirects to login.

## Migration command

### Command contract

Every supported native, Docker, and Nix distribution must expose equivalent
commands:

```text
headplane migrate-local-admin \
  --config /etc/headplane/config.yaml \
  --legacy-db /var/lib/headscale/db.sqlite \
  [--username admin] \
  [--dry-run]

headplane hash-password --password-stdin
headplane reset-local-admin-password --config /etc/headplane/config.yaml --password-stdin
```

Passwords are accepted via stdin or a terminal prompt, never command arguments.

### Algorithm

1. Parse the target Headplane config and validate it is writable.
2. Open the source SQLite database read-only and verify the legacy
   `headplane_users` table and required columns exist.
3. Select rows with `role = 'admin'`.
4. Select the sole row, or require an exact `--username` when multiple legacy
   administrators exist.
5. Validate the selected `password_hash` as a supported bcrypt hash.
6. Evaluate target config state:
   - no `user` block: prepare the selected username/hash;
   - matching username/hash: successful no-op;
   - differing credentials or invalid target user data: fail without overwrite.
7. In dry-run mode, show only username and planned action; never emit a hash.
8. Otherwise back up and atomically write only `user.username` and
   `user.password`.
9. Re-open and validate the resulting config before reporting success.

The migration deliberately does not read or convert `headplane_settings` API-key
ciphertext because it was encrypted using old per-session material and cannot be
safely converted to a config-managed service key.

### Mandatory idempotency invariants

- The legacy database remains byte-for-byte unchanged.
- A matching target config is a successful no-op.
- Different target credentials are never overwritten implicitly.
- A failed temporary-file write leaves the source config intact.
- A completed rename yields one parseable configured administrator.
- Repeated completed migrations are safe.

## Rollout order

1. Update specs, roadmap, config examples, changelogs, and user documentation.
2. Add config schema, bcrypt dependency, and password hashing command.
3. Implement local Headplane authentication while retaining API-key login.
4. Disable OIDC/proxy runtime paths and retire Headscale-local auth endpoints.
5. Replace the user page with Administration/password reset/API-key UI.
6. Add migration and recovery commands.
7. Perform fixture, browser, integration, and staged-upgrade verification.

Release the Headplane submodule and parent Headscale changes together. Avoid a
mixed release where Headplane requires endpoints already removed from Headscale,
or where old Headplane expects them after route retirement.

## Test matrix

| Area | Required coverage |
| --- | --- |
| Config | Valid bcrypt value, plaintext/malformed rejection, missing values, environment overrides, secret redaction. |
| Local login | Valid credentials, incorrect username/password, empty input, rate limit, session expiry, no Headscale password-login call. |
| Password reset | Current-password check, confirmation mismatch, bcrypt cost, atomic-write failure, session invalidation, immutable config. |
| OIDC/proxy | Retained code is not reachable while single-admin mode is active. |
| API keys | List/create/expire/delete contracts, one-time secret display, service-key protection, rotation verification, self-revocation logout. |
| Migration | No admin, multiple admin, explicit selection, invalid hash, dry run, target conflict, success, failed write, repeated run, database immutability. |

Browser tests must verify local login, API-key login, no OIDC/proxy login choice,
the one-admin notice, password reset, absence of user CRUD controls, and API-key
rotation safeguards.

Run after implementation from the Headplane submodule:

```bash
pnpm run test:unit
pnpm run typecheck
pnpm run lint
pnpm run test:integration
pnpm run test:e2e
pnpm run docs:build
```

Run from the parent repository:

```bash
go test ./hscontrol/...
make fmt
make lint
```

Read `cmd/hi/README.md` again before selecting or running any `hi` integration
command.

## Documentation completion criteria

Synchronize the parent roadmap/changelog; Headplane roadmap, specification,
architecture, changelog, and config example; and Headscale authentication,
settings, configuration, troubleshooting, and API-reference pages. User-facing
documentation and the Administration page must clearly state that one local
administrator is supported after migration.
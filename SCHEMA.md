# Schema

This document is the index for database schema definitions in this fork.

## Authoritative sources

- [`hscontrol/db/schema.sql`](./hscontrol/db/schema.sql) is the canonical SQLite
  schema contract used by schema validation and tests.
- [`hscontrol/db/migrate_*.go`](./hscontrol/db/) contains the explicit migration
  DDL. Migration order is immutable; new migrations are appended only.
- [`hscontrol/types/`](./hscontrol/types/) defines the Go records and generated
  view/clone accessors that map to the database rows.
- [`hscontrol/api/v1/`](./hscontrol/api/v1/) and
  [`hscontrol/api/v2/`](./hscontrol/api/v2/) define the HTTP request and response
  shapes exposed to clients.

## Fork-specific schema areas

The fork adds or depends on the following schema areas:

- `oidc_groups` on `users`: OIDC `groups` claim values refreshed at login and
  resolved by the policy engine.
- `webhooks`: webhook configuration persisted for monitoring/alerting events.
- `headplane_users` and `headplane_settings`: legacy Headplane local-account
  data retained for rollback after the Phase 13c single-administrator migration.
  They are not the runtime source of local authentication.
- `credentials`: shared credential storage used by API keys, pre-auth keys,
  OAuth clients, and OAuth access tokens in current upstream-compatible code.

Any schema change must update the canonical schema, append an explicit migration,
regenerate affected view/clone/API code, add tests, and update the relevant
documentation. See [`AGENTS.md`](./AGENTS.md) for the migration rules.

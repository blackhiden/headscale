# Phase 13c: Single Local Administrator Migration Runbook

**Status:** Complete — Implementation and migration finished  
**Completed:** 2026-10-06  
**Technical design:** [Single Local Administrator Migration Plan](./phase13c-single-admin-migration-plan.md)  
**Authoritative scope:** [Phase 13c in the roadmap](../ROADMAP.md#phase-13c--single-local-administrator-migration-planned)

## Completion Summary

Phase 13c — Single Local Administrator Migration is **complete**. This runbook documents the migration procedure that has been implemented and verified. Operators can now use the delivered migration tooling and follow the procedures outlined below to migrate their deployments from the legacy multi-user system to the new single-admin configuration.

**What was delivered:**
- Migration command for converting legacy database users to config-based admin
- Password hashing and reset commands
- Comprehensive validation and idempotency checks
- Rollback procedures and backup guidance
- API key rotation and management UI
- Full verification and testing coverage

See the [CHANGELOG](../CHANGELOG.md) for implementation details.

---

## What changes

This migration changes Headplane local authentication from Headscale-backed
multi-user accounts to one Headplane-configured local administrator.

After migration:

- You sign in locally with one configured username and password.
- The stored local password is a bcrypt verification hash, never plaintext.
- Headscale API keys remain a separate secondary sign-in method.
- The Administration page replaces local user management with password reset and
  API-key lifecycle controls.
- OIDC and proxy authentication are disabled by this single-admin mode.
- Extra legacy local accounts cannot sign in.

The migration copies one legacy administrator's existing bcrypt hash. The
selected administrator continues using the same password. The tool cannot
recover or display plaintext passwords.

## Prerequisites and limitations

### Supported legacy source

The migration supports the fork's current **SQLite** legacy source containing
the `headplane_users` table. It opens that database read-only.

It does not migrate secondary accounts, roles, themes, profile names,
`headplane_settings` API-key ciphertext, OIDC/proxy users, Headscale tailnet
users, nodes, policies, or API keys.

### Required access

Before you begin, ensure you can:

- read the legacy Headscale SQLite database;
- write and back up Headplane's config, normally `/etc/headplane/config.yaml`;
- stop and start the Headplane service; and
- identify the configured Headscale service API key.

Your target Headplane config must already contain a valid `headscale.api_key` or
`headscale.api_key_path`. The migration does not create or copy this service key.

### Protect the configuration

The config file is sensitive. Use a restricted owner and mode, replacing the
service account to match your deployment:

```bash
sudo chown headplane:headplane /etc/headplane/config.yaml
sudo chmod 0600 /etc/headplane/config.yaml
sudo chmod 0700 /etc/headplane
```

Never commit config, bcrypt hashes, API keys, or cookie secrets to source
control. Prefer a managed secret file for the Headscale API key:

```yaml
headscale:
  api_key_path: "${CREDENTIALS_DIRECTORY}/headscale-api-key"
```

Systemd credentials, Docker secrets, Kubernetes Secrets, and equivalent managed
secret files are appropriate. Encrypting config alone is insufficient because
Headplane would still require a decryption key at startup.

## Preflight and backup

### Record paths and deployed versions

Record the current Headplane and Headscale versions/image tags plus the paths to:

- Headplane config;
- Headplane persistent data;
- Headscale SQLite database; and
- any SQLite `-wal` and `-shm` files.

Keep the existing artifacts available throughout the rollback window.

### Stop writes

Stop Headplane before migration. For the most conservative SQLite backup, stop
Headscale too before copying database files.

```bash
sudo systemctl stop headplane
sudo systemctl stop headscale
```

For Docker Compose, stop services without deleting persistent volumes:

```bash
docker compose stop headplane headscale
```

### Create a consistent backup

```bash
backup_dir="/root/headplane-single-admin-backup-$(date -u +%Y%m%dT%H%M%SZ)"
sudo install -d -m 0700 "$backup_dir"

sudo cp -a /etc/headplane/config.yaml "$backup_dir/config.yaml.before"

# Replace with your real Headscale database path.
sudo cp -a /var/lib/headscale/db.sqlite "$backup_dir/"
sudo cp -a /var/lib/headscale/db.sqlite-wal "$backup_dir/" 2>/dev/null || true
sudo cp -a /var/lib/headscale/db.sqlite-shm "$backup_dir/" 2>/dev/null || true
```

Do not copy only the primary SQLite file while Headscale is writing in WAL mode.
Stop the service first or use a SQLite-consistent backup method for your
deployment.

You may start Headscale again after this backup if required. Keep Headplane
stopped until migration validation completes.

## Inspect legacy administrators

The migration selects rows with `role = 'admin'` from legacy
`headplane_users`. List only non-secret metadata:

```bash
sqlite3 /var/lib/headscale/db.sqlite \
  "SELECT id, username, role, created_at FROM headplane_users WHERE role = 'admin' ORDER BY id;"
```

| Result | Required action |
| --- | --- |
| One administrator | Continue; the tool selects it automatically. |
| Multiple administrators | Choose the retained account and pass its exact username with `--username`. |
| No administrator | Stop; repair or create a known administrator before migration. |
| Missing table/columns | Stop; verify database path and supported legacy release. |

Do not query, print, copy, or share `password_hash` values. The migration reads
the selected hash directly.

## Dry run

Use the packaged Headplane command. Native, Docker, and Nix installations may
invoke it differently, but its command contract is:

```bash
headplane migrate-local-admin \
  --config /etc/headplane/config.yaml \
  --legacy-db /var/lib/headscale/db.sqlite \
  --dry-run
```

When multiple legacy administrators exist:

```bash
headplane migrate-local-admin \
  --config /etc/headplane/config.yaml \
  --legacy-db /var/lib/headscale/db.sqlite \
  --username admin \
  --dry-run
```

The dry run may report the selected username and whether the target config is
already migrated. It must not print password hashes, plaintext passwords, API
keys, cookie secrets, or session tokens.

Do not proceed unless the selected username is the account you intend to retain.

## Perform the migration

Repeat the successful dry-run command without `--dry-run`:

```bash
headplane migrate-local-admin \
  --config /etc/headplane/config.yaml \
  --legacy-db /var/lib/headscale/db.sqlite
```

Add `--username` exactly as used in a dry run when the source has more than one
administrator.

The migration must:

1. open SQLite read-only;
2. select exactly one legacy administrator;
3. validate that the selected stored value is a supported bcrypt hash;
4. create a timestamped `0600` config backup beside the target config;
5. write through a same-directory temporary file;
6. atomically rename the new config into place; and
7. re-read and validate the result before returning success.

The resulting config includes:

```yaml
user:
  username: admin
  password: "$2b$12$..."
```

The command must not alter `headplane_users`, `headplane_settings`, or any other
legacy database table.

## Idempotency and expected failures

You can safely repeat the migration command.

| Condition | Expected result |
| --- | --- |
| No target `user` block | The selected username and copied bcrypt hash are written. |
| Target username/hash already match | Successful no-op. |
| Target config contains different credentials | Safe failure without overwrite. |
| Multiple legacy admins without `--username` | Safe failure requesting explicit selection. |
| Malformed legacy hash | Safe failure without config change. |
| Failure before atomic rename | Existing config remains valid; correct permissions, disk space, or input and retry. |
| Interruption after completed rename | Rerun validates the completed result and reports a no-op. |

Do not manually copy a legacy password hash unless you have verified backups and
understand the resulting config state. Prefer to correct inputs and rerun the
tool.

## Start and validate

Start the single-admin Headplane release:

```bash
sudo systemctl start headplane
# Or: docker compose up -d headplane
```

Validate in a private browser session:

1. Open `/admin/login`.
2. Confirm Password and API Key are visible login methods.
3. Confirm OIDC/proxy sign-in controls are absent.
4. Log in with the retained username and existing password.
5. Confirm dashboard data loads. If it does not, validate
   `headscale.api_key` or `headscale.api_key_path`.
6. Open `/admin/users`; it should be an Administration page, not a local user
   management page.
7. Confirm the page says only one local administrator is supported.
8. Log out and confirm valid API-key login still works.
9. Confirm another legacy local account can no longer complete password login.

Do not delete backups or legacy tables after the first validation.

## Password reset and recovery

### Administration page reset

For a regular writable config, the Administration page verifies the current
password, writes a fresh bcrypt hash, and invalidates existing local password
sessions.

If configuration is immutable or comes from an externally managed source, the
page must be disabled. Update the managed config/secret and restart Headplane
instead.

### Host-level recovery

If you lose the local password, use the host-level command without putting a
password in shell history or process arguments:

```bash
read -rs -p 'New Headplane password: ' new_password
printf '\n'
printf '%s' "$new_password" | \
  headplane reset-local-admin-password \
    --config /etc/headplane/config.yaml \
    --password-stdin
unset new_password
```

For fresh setup, generate a bcrypt value by reading a password from stdin:

```bash
headplane hash-password --password-stdin
```

Store only the output hash in `user.password`, restart Headplane, and verify in
a new browser session.

## API-key rotation and revocation

The Administration page lists API-key metadata and can create, expire, revoke,
or delete keys where the connected Headscale version supports each operation.

### Rotate a normal API key

1. Create a replacement key with the desired expiry.
2. Copy the plaintext key while it is shown; it is displayed once only.
3. Update any scripts or services using the previous key.
4. Verify those clients work with the replacement.
5. Expire or delete the old key.

### Rotate Headplane's configured service key

The configured `headscale.api_key` or `headscale.api_key_path` is protected. The
UI must not allow direct expiry or deletion of it.

For direct writable `headscale.api_key`:

1. Create replacement key.
2. Save it while displayed.
3. Allow Headplane to validate it.
4. Allow Headplane to atomically update its direct configuration and validate
   dashboard access.
5. Expire the prior service key only after this succeeds.

For an externally managed `headscale.api_key_path`:

1. Create and save the replacement key.
2. Update the external secret provider yourself.
3. Restart or reload Headplane.
4. Validate dashboard access.
5. Expire the prior service key.

If you revoke the API key used by your current API-key browser session, Headplane
signs you out immediately. This is expected.

## Rollback

### Migration fails before atomic rename

1. Keep Headplane stopped.
2. Correct the reported issue, such as path, permissions, source selection,
   malformed hash, or disk capacity.
3. Run the dry run again.
4. Rerun migration.

The target config remains unchanged in this state. Database restore is not
required.

### Migration completed but must be abandoned

1. Stop Headplane.
2. Restore the original config from the backup created before migration or from
   the migration-created timestamped backup.
3. Redeploy the prior compatible Headplane and Headscale releases.
4. Clear the `_hp_auth` cookie or use a private browser session.
5. Start the prior services.
6. Validate the legacy login and administration behavior.

```bash
sudo systemctl stop headplane
sudo cp -a "$backup_dir/config.yaml.before" /etc/headplane/config.yaml
# Re-deploy previous compatible Headplane and Headscale versions here.
sudo systemctl start headscale
sudo systemctl start headplane
```

The migration reads the legacy SQLite database only. Database restoration is
normally unnecessary; restore the complete SQLite backup only after an unrelated
database change, corruption, or manual modification.

## Cleanup policy

Do not remove legacy `headplane_users` or `headplane_settings` data immediately.
Keep database and config backups through the documented rollback window. Any
future destructive cleanup requires a separately announced migration, fresh
backups, release notes, and tested downgrade/recovery instructions.
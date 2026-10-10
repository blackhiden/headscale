# Upgrade an existing installation

!!! tip "Required update path"

    It's required to update from one stable version to the next (e.g. 0.26.0 → 0.27.1 → 0.28.0) without skipping minor
    versions in between. You should always pick the latest available patch release.

Update an existing Headscale installation to a new version:

- Read the announcement on the [GitHub releases](https://github.com/arsydoni4326-alt/headscale/releases) page for the new
  version. It lists the changes of the release along with possible breaking changes and version-specific upgrade
  instructions.
- Stop Headscale
- **[Create a backup of your installation](#backup)**
- Update Headscale to the new version, preferably by following the same installation method.
- Compare and update the [configuration](../ref/configuration.md) file.
- Start Headscale

## Backup

Headscale applies database migrations during upgrades and we highly recommend to create a backup of your database before
upgrading. A full backup of Headscale depends on your individual setup, but below are some typical setup scenarios.

=== "Standard installation"

    An installation that follows our [official releases](install/official.md) setup guide uses the following paths:

    - [Configuration file](../ref/configuration.md): `/etc/headscale/config.yaml`
    - Data directory: `/var/lib/headscale`
    - SQLite as database: `/var/lib/headscale/db.sqlite`

    ```console
    TIMESTAMP=$(date +%Y%m%d%H%M%S)
    cp -aR /etc/headscale /etc/headscale.backup-$TIMESTAMP
    cp -aR /var/lib/headscale /var/lib/headscale.backup-$TIMESTAMP
    ```

=== "Container"

    An installation that follows our [container](install/container.md) setup guide uses a single source volume directory
    that contains the configuration file, data directory and the SQLite database.

    ```console
    cp -aR /path/to/headscale /path/to/headscale.backup-$(date +%Y%m%d%H%M%S)
    ```

=== "PostgreSQL"

    Please follow PostgreSQL's [Backup and Restore](https://www.postgresql.org/docs/current/backup.html) documentation
    to create a backup of your PostgreSQL database.

## Manual Migration Scenarios

### Upgrading from v0.35.4 to v0.35.5+ with Unversioned Binary

If you deployed a binary built without version information (e.g., using `make build` or `go build` directly
instead of `make build VERSION=...`), automatic database migrations will be skipped. You may encounter this error:

```
SQLite schema failed to validate:
>> Add table "webhooks"
Error: initializing: creating new headscale: init state: initializing database: validating schema
```

**Root Cause:** The binary lacks embedded version information, causing headscale to skip database version checks
and migrations. The code expects the `webhooks` table (added in v0.35.5), but it doesn't exist in the database.

**Manual Fix:**

1. **Stop headscale and backup your database:**

    ```bash
    # Stop headscale service
    systemctl stop headscale  # or docker compose down, etc.
    
    # Find database location (check your config.yaml)
    grep -i "db_path\|database_path" /etc/headscale/config.yaml
    
    # Backup database
    cp /var/lib/headscale/db.sqlite /var/lib/headscale/db.sqlite.backup-$(date +%Y%m%d-%H%M%S)
    ```

2. **Apply missing schema manually:**

    ```bash
    sqlite3 /var/lib/headscale/db.sqlite << 'EOF'
    -- Create the webhooks table
    CREATE TABLE webhooks(
      id integer PRIMARY KEY AUTOINCREMENT,
      created_at datetime,
      updated_at datetime,
      deleted_at datetime,
      name text NOT NULL,
      url text NOT NULL,
      events text NOT NULL,
      headers text,
      secret text,
      enabled numeric NOT NULL DEFAULT true,
      timeout_seconds integer NOT NULL DEFAULT 10
    );
    
    -- Create indexes
    CREATE INDEX idx_webhooks_deleted_at ON webhooks(deleted_at);
    CREATE UNIQUE INDEX idx_webhooks_name ON webhooks(name);
    EOF
    ```

3. **Verify schema was applied:**

    ```bash
    sqlite3 /var/lib/headscale/db.sqlite "SELECT name FROM sqlite_master WHERE type='table' AND name='webhooks';"
    # Should output: webhooks
    ```

4. **Restart headscale:**

    ```bash
    systemctl start headscale  # or docker compose up -d, etc.
    ```

**Prevention:** Always build binaries with version information to enable automatic migrations:

```bash
make build VERSION=v0.35.5-arsydoni4326-alt
# or
make release VERSION=v0.35.5-arsydoni4326-alt
```

Alternatively, use pre-built binaries from [GitHub releases](https://github.com/arsydoni4326-alt/headscale/releases),
which always include proper version information.

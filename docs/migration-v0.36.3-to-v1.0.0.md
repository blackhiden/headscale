# Migration Guide: v0.36.3-arsydoni4326-alt → v1.0.0-arsydoni4326

**Date:** 2026-10-06  
**Status:** Active Migration Guide  
**Phase:** Phase 15 - Database Version Migration

## Overview

This guide covers migrating from Headscale v0.36.3-arsydoni4326-alt to v1.0.0-arsydoni4326.

**Important:** This migration updates **version metadata only**. No database schema changes exist between these versions. The v1.0.0 release includes Phase 13c code retirement (Headplane local-auth runtime removal) but preserves all database structures.

## What's Changed in v1.0.0

### Code Changes Only
- ✅ **Removed:** Automatic default admin user creation from `cfg.Headplane.Password`
- ✅ **Removed:** Tests for database-backed Headplane user creation
- ✅ **Preserved:** All Headscale core functionality (nodes, users, routes, policies, API keys)
- ✅ **Preserved:** Legacy `headplane_users` and `headplane_settings` tables (for rollback)
- ✅ **Required:** Headplane Phase 13c configuration (config-based authentication)

### No Schema Changes
- Database structure is **identical** between v0.36.3 and v1.0.0
- All migrations from v0.36.3 are present in v1.0.0
- No new migrations added in v1.0.0

## Prerequisites

### 1. Current Version Check

Verify you're running v0.36.3-arsydoni4326-alt:

```bash
headscale version
# or
docker compose exec headscale headscale version
```

Expected output: `v0.36.3-arsydoni4326-alt` or similar v0.36.x

### 2. System Requirements

- SQLite3 command-line tool installed
- Access to Headscale database file (typically `/var/lib/headscale/headscale.db`)
- Backup storage space (database + WAL + SHM files)
- If using Headplane: access to `/etc/headplane/config.yaml` or equivalent

### 3. Headplane Phase 13c Configuration

If you use Headplane, verify Phase 13c configuration exists:

```yaml
# /etc/headplane/config.yaml
user:
  username: admin
  password: "$2b$12$..."  # bcrypt hash, NOT plaintext

headscale:
  url: "http://headscale:8080"
  api_key: "your-admin-api-key"
  # or: api_key_path: "${CREDENTIALS_DIRECTORY}/headscale-api-key"
```

If not configured, you'll need to run the Headplane migration tool (see Step 4 below).

## Migration Procedure

### Step 1: Backup Everything

**Critical:** Always backup before migration.

```bash
# Stop services
docker compose down
# or
systemctl stop headscale headplane

# Create backup directory
mkdir -p /tmp/headscale-backup-$(date +%Y%m%d-%H%M%S)
cd /tmp/headscale-backup-*

# Backup database files
cp /var/lib/headscale/headscale.db ./headscale.db.v0.36.3.backup
cp /var/lib/headscale/headscale.db-wal ./headscale.db-wal.backup 2>/dev/null || true
cp /var/lib/headscale/headscale.db-shm ./headscale.db-shm.backup 2>/dev/null || true

# Backup entire directory (safest)
tar -czf headscale-full-backup.tar.gz /var/lib/headscale/

# Backup Headplane config if using Headplane
cp /etc/headplane/config.yaml ./headplane-config.yaml.backup 2>/dev/null || true

# Verify backups exist
ls -lh
```

### Step 2: Update Database Version Metadata

Create the migration SQL script:

```bash
cat > /tmp/migrate-v0-to-v1.sql << 'EOF'
-- Phase 15: Update database version metadata for v1.0.0 migration
-- NO SCHEMA CHANGES - version metadata only

BEGIN TRANSACTION;

-- Update the version in database_versions table
UPDATE database_versions 
SET version = 'v1.0.0-arsydoni4326',
    updated_at = datetime('now')
WHERE id = 1;

-- If no version record exists, insert it
INSERT OR IGNORE INTO database_versions (id, version, updated_at) 
VALUES (1, 'v1.0.0-arsydoni4326', datetime('now'));

-- Verify the migration
SELECT 'Version updated to: ' || version as result 
FROM database_versions 
WHERE id = 1;

COMMIT;
EOF
```

Apply the migration:

```bash
# Apply SQL migration
sqlite3 /var/lib/headscale/headscale.db < /tmp/migrate-v0-to-v1.sql

# Expected output: Version updated to: 1.0.0-arsydoni4326

# Verify version was updated
sqlite3 /var/lib/headscale/headscale.db \
  "SELECT key, value FROM kv WHERE key = 'last_seen_version';"

# Expected output: last_seen_version|1.0.0-arsydoni4326
```

### Step 3: Verify Data Integrity

Check that your data is intact:

```bash
# Count nodes (should match your known count)
sqlite3 /var/lib/headscale/headscale.db \
  "SELECT COUNT(*) FROM nodes;"

# Count users (should match your known count)
sqlite3 /var/lib/headscale/headscale.db \
  "SELECT COUNT(*) FROM users;"

# Count routes
sqlite3 /var/lib/headscale/headscale.db \
  "SELECT COUNT(*) FROM routes;"

# Verify Headplane tables still exist (for rollback)
sqlite3 /var/lib/headscale/headscale.db \
  "SELECT name FROM sqlite_master WHERE type='table' AND name LIKE 'headplane%';"

# Expected: headplane_users, headplane_settings
```

### Step 4: Upgrade to v1.0.0

```bash
# Switch to v1.0.0
cd /home/denny/Project/headscale-project/headscale
git fetch --tags
git checkout v1.0.0-arsydoni4326

# Rebuild
make clean
make build

# Or for Docker
docker compose build --no-cache

# Start services
docker compose up -d

# Or for systemd
systemctl start headscale
systemctl start headplane  # if applicable
```

### Step 5: Verify Migration Success

```bash
# Check Headscale version
docker compose exec headscale headscale version
# Expected: v1.0.0-arsydoni4326

# Check logs for successful startup
docker compose logs headscale | head -50
# Should NOT see "version check" errors

# Verify all nodes are present
docker compose exec headscale headscale nodes list

# Verify all users are present
docker compose exec headscale headscale users list

# Verify all routes are present
docker compose exec headscale headscale routes list
```

## Verification Checklist

After migration, verify:

- [ ] Headscale starts without "version check: major version change not supported" error
- [ ] `headscale version` reports `v1.0.0-arsydoni4326`
- [ ] `headscale nodes list` shows all existing nodes (count matches backup)
- [ ] `headscale users list` shows all existing users (count matches backup)
- [ ] `headscale routes list` shows all existing routes (count matches backup)
- [ ] `headscale apikeys list` shows existing API keys
- [ ] API keys work for Headscale API access
- [ ] No data loss compared to v0.36.3 backup

## Rollback Procedure

If v1.0.0 fails or has issues:

```bash
# Stop services immediately
docker compose down

# Restore database from backup
cp /tmp/headscale-backup-*/headscale.db.v0.36.3.backup \
   /var/lib/headscale/headscale.db

# Remove WAL and SHM files to force clean state
rm -f /var/lib/headscale/headscale.db-wal
rm -f /var/lib/headscale/headscale.db-shm

# Revert to v0.36.3
cd /home/denny/Project/headscale-project/headscale
git checkout v0.36.3-arsydoni4326-alt

# Rebuild
make clean && make build

# Restart services
docker compose up -d

# Verify rollback success
headscale version  # should show v0.36.3
headscale nodes list  # should show all nodes
```

## Troubleshooting

### Error: "version check: major version change not supported"

**Cause:** Database version metadata was not updated  
**Solution:** Rerun Step 2

```bash
sqlite3 /var/lib/headscale/headscale.db < /tmp/migrate-v0-to-v1.sql
```

### Missing nodes/users/routes after migration

**Cause:** Database corruption or incomplete backup restore  
**Solution:** Immediately rollback

---

**Migration Guide Version:** 1.0  
**Last Updated:** 2026-10-06

# Troubleshooting

This guide covers common issues you may encounter with Headscale and Headplane, along with their solutions.

## Authentication Issues

### Password Authentication (Headplane Web UI)

#### Invalid password error

**Symptoms:**

- Login form shows "Invalid password" error
- HTTP 401 Unauthorized response

**Possible causes and solutions:**

1. **Password mismatch:**
   - Verify the password matches your `config.yaml` or environment variable exactly
   - Check for leading/trailing whitespace in the password
   - Passwords are case-sensitive

2. **Configuration not loaded:**
   - Restart Headscale after changing the password
   - Verify Headscale is reading the correct configuration file
   - Check logs for configuration parsing errors: `journalctl -u headscale -n 50`

3. **Environment variable precedence:**
   - If `HEADSCALE_HEADPLANE_PASSWORD` is set, it overrides `config.yaml`
   - Check environment variables: `env | grep HEADSCALE_HEADPLANE`
   - Unset the environment variable or update it to match your intended password

4. **YAML syntax error:**
   - Verify the `headplane.password` field is correctly indented in `config.yaml`
   - Run `headscale configtest` to validate your configuration
   - Example correct syntax:
     ```yaml
     headplane:
       password: "your-password-here"
     ```

**Diagnostic commands:**

```bash
# Verify configuration is valid
headscale configtest

# Check which config file is being used
headscale version

# Check environment variables
env | grep HEADSCALE_HEADPLANE

# View recent logs
journalctl -u headscale -n 50
```

#### Rate limiting (429 Too Many Requests)

**Symptoms:**

- Login attempts return HTTP 429
- Error message: "Rate limit exceeded. Please try again later."

**Causes:**

- More than 5 failed login attempts within 1 minute from the same IP address

**Solutions:**

1. **Wait and retry:**
   - Wait 60 seconds for the rate limit window to reset
   - Try logging in again with the correct password

2. **Check for automated scripts:**
   - Verify no scripts are making excessive login attempts
   - Review cron jobs or monitoring tools that might be hitting the login endpoint

3. **Shared IP address:**
   - If behind a proxy or NAT, multiple users may share the same IP
   - Consider increasing the rate limit in a future Headscale update (not currently configurable)

4. **Verify password before retrying:**
   - Double-check the password before attempting again
   - Each failed attempt counts toward the rate limit

#### Session expired

**Symptoms:**

- Redirected to login page after being logged in
- HTTP 401 Unauthorized on API requests
- "Session expired" message in UI

**Causes:**

- Sessions expire after 24 hours
- System clock skew between client and server

**Solutions:**

1. **Log in again:**
   - Sessions are designed to expire after 24 hours for security
   - Simply log in again to get a new session token

2. **Check system clock:**
   - Verify both client and server clocks are synchronized
   - Use NTP to keep clocks in sync:
     ```bash
     timedatectl status
     systemctl status systemd-timesyncd
     ```

3. **Browser issues:**
   - Clear browser cache and cookies for the Headscale domain
   - Try logging in with an incognito/private window

#### Password change not taking effect

**Symptoms:**

- Changed password in config but old password still works
- Changed password but can't log in with new one

**Solutions:**

1. **Restart Headscale:**
   - Configuration changes require a restart
   ```bash
   systemctl restart headscale
   ```

2. **Verify configuration precedence:**
   - Environment variable takes precedence over config file
   - If `HEADSCALE_HEADPLANE_PASSWORD` is set, update it or unset it
   ```bash
   # Check current env var
   env | grep HEADSCALE_HEADPLANE_PASSWORD
   
   # Unset if needed
   unset HEADSCALE_HEADPLANE_PASSWORD
   ```

3. **Validate configuration syntax:**
   ```bash
   headscale configtest
   ```

4. **Check which config file is loaded:**
   - Headscale searches multiple paths in order:
     - `/etc/headscale/config.yaml`
     - `~/.headscale/config.yaml`
     - `./config.yaml` (current directory)
   - Verify you're editing the correct file
   - Use `headscale --config /path/to/config.yaml` to specify explicitly

5. **Existing sessions remain valid:**
   - Changing the password does not invalidate existing sessions
   - Users with active sessions can continue using Headplane for up to 24 hours
   - To force logout, restart Headscale or wait for session expiry

### API Key Authentication

#### API key not working

**Symptoms:**

- HTTP 401 Unauthorized when using API key
- "Invalid API key" error

**Solutions:**

1. **Verify API key format:**
   - API keys should be used in the `Authorization: Bearer <KEY>` header
   - Example:
     ```bash
     curl -H "Authorization: Bearer YOUR_API_KEY" \
         https://headscale.example.com/api/v1/user
     ```

2. **Check API key expiry:**
   ```bash
   headscale apikeys list
   ```
   - Look for your API key and check its expiration date
   - Create a new key if expired

3. **API key vs. session token:**
   - Session tokens (from password login) start with `hp_` and expire after 24 hours
   - API keys are longer-lived (default 90 days)
   - Don't confuse the two token types

4. **Verify API key wasn't expired:**
   ```bash
   headscale apikeys list
   ```
   - If the key is no longer in the list, it was expired
   - Create a new API key

## Configuration Issues

### Configuration file not found

**Symptoms:**

- "Configuration file not found" error on startup
- Headscale fails to start

**Solutions:**

1. **Create configuration file:**
   - Download example config:
     ```bash
     wget -O /etc/headscale/config.yaml \
         https://raw.githubusercontent.com/juanfont/headscale/main/config-example.yaml
     ```

2. **Specify config path explicitly:**
   ```bash
   headscale --config /path/to/config.yaml serve
   ```

3. **Use environment variable:**
   ```bash
   export HEADSCALE_CONFIG=/path/to/config.yaml
   headscale serve
   ```

### YAML syntax errors

**Symptoms:**

- "Error parsing configuration file"
- Headscale fails to start with parsing error

**Solutions:**

1. **Validate configuration:**
   ```bash
   headscale configtest
   ```

2. **Common YAML errors:**
   - Incorrect indentation (use 2 spaces, not tabs)
   - Missing quotes around special characters
   - Incorrect nesting of fields

3. **Check headplane section syntax:**
   ```yaml
   # Correct:
   headplane:
     password: "your-password"
   
   # Wrong (missing colon):
   headplane
     password: "your-password"
   
   # Wrong (incorrect indentation):
   headplane:
   password: "your-password"
   ```

### Environment variable not taking effect

**Symptoms:**

- Set `HEADSCALE_HEADPLANE_PASSWORD` but password doesn't work
- Configuration seems to be ignored

**Solutions:**

1. **Verify environment variable is set:**
   ```bash
   env | grep HEADSCALE_HEADPLANE_PASSWORD
   ```

2. **Restart Headscale after setting:**
   ```bash
   export HEADSCALE_HEADPLANE_PASSWORD="your-password"
   systemctl restart headscale
   ```

3. **Check systemd service file:**
   - If running as a systemd service, env vars must be set in the service file
   - Edit `/etc/systemd/system/headscale.service`:
     ```ini
     [Service]
     Environment="HEADSCALE_HEADPLANE_PASSWORD=your-password"
     ```
   - Reload and restart:
     ```bash
     systemctl daemon-reload
     systemctl restart headscale
     ```

4. **Use systemd environment file:**
   ```bash
   # Create /etc/headscale/environment
   echo 'HEADSCALE_HEADPLANE_PASSWORD=your-password' > /etc/headscale/environment
   chmod 600 /etc/headscale/environment
   ```
   
   Update service file:
   ```ini
   [Service]
   EnvironmentFile=/etc/headscale/environment
   ```

## Database Migration Issues

### Headplane schema validation rejects required tables

**Symptoms:**

- Startup fails with schema validation output that requests removal of
  `headplane_users`, `headplane_settings`, or their indexes.
- Removing the objects lets Headscale start, but `GET /api/v1/headplane/users`
  returns HTTP 500 and the Headplane users page cannot load.

**Cause:**

Affected builds created the `headplane_users` and `headplane_settings` tables
without listing them in Headscale's canonical SQLite schema. The validator then
considered the required tables unexpected. Deleting them only hides the schema
error; it breaks Headplane authentication, user management, and settings.

**Resolution:**

Upgrade to a build containing the Headplane schema repair migrations. They define
both tables in the canonical schema, recreate either table if it was removed by
the earlier workaround, and normalize indexes previously created by GORM with
quoted identifiers.

1. Stop Headscale and back up the database:

   ```bash
   systemctl stop headscale
   cp /var/lib/headscale/db.sqlite \
     /var/lib/headscale/db.sqlite.backup-$(date +%Y%m%d-%H%M%S)
   ```

2. Deploy the fixed Headscale image or binary, then start it normally:

   ```bash
   systemctl start headscale
   # or: docker compose up -d headscale
   ```

3. Confirm the repair migration and required objects:

   ```bash
   sqlite3 /var/lib/headscale/db.sqlite <<'EOF'
   SELECT id FROM migrations
   WHERE id IN (
     '202610050900-repair-headplane-schema',
     '202610050930-normalize-headplane-indexes'
   )
   ORDER BY id;
   SELECT name FROM sqlite_master
   WHERE type IN ('table', 'index')
     AND name IN (
       'headplane_users',
       'headplane_settings',
       'idx_headplane_users_username',
       'idx_headplane_settings_user_id'
     )
   ORDER BY name;
   EOF
   ```

4. Reload `/admin/admin/users` and verify that
   `/api/v1/headplane/users` returns an object containing a `users` array.

Do not manually create or drop only one of these tables, and do not mark the
repair migration as applied by hand. The migration is idempotent and preserves
existing Headplane data. If the repaired build still fails, restore the backup
and include the complete validation diff and `sqlite_master` output in a bug
report.

### Migration failed during upgrade

**Symptoms:**

- Headscale won't start after upgrade
- Error: `migration failed: ...`
- Database is in an inconsistent state

**Solution:**

1. **Restore from backup:**
   
   ```bash
   systemctl stop headscale
   cp /var/lib/headscale/db.sqlite.backup /var/lib/headscale/db.sqlite
   systemctl start headscale
   ```

2. **If no backup exists, try to repair:**
   
   Check the specific migration error in logs and consult the [Headplane schema validation section](#headplane-schema-validation-rejects-required-tables) above.

3. **Report the issue:**
   
   If the repair migration does not complete, file a bug report with:
   - The exact error message from logs
   - Output of `.schema` from sqlite3
   - Output of `SELECT * FROM migrations;`

## Network and Connectivity Issues

### Cannot connect to Headplane web interface

**Symptoms:**

- Browser shows "Connection refused" or "Unable to connect"
- Headplane UI doesn't load

**Solutions:**

1. **Verify Headscale is running:**
   ```bash
   systemctl status headscale
   ```

2. **Check listen address:**
   - Verify `listen_addr` in config.yaml
   - For remote access, use `0.0.0.0:8080` not `127.0.0.1:8080`

3. **Check firewall rules:**
   ```bash
   # Allow port 8080
   ufw allow 8080/tcp
   
   # Or for firewalld
   firewall-cmd --add-port=8080/tcp --permanent
   firewall-cmd --reload
   ```

4. **Verify port is listening:**
   ```bash
   ss -tlnp | grep 8080
   ```

### HTTPS/TLS issues

**Symptoms:**

- "Your connection is not private" browser warning
- TLS certificate errors

**Solutions:**

See the [TLS Configuration Guide](ref/tls.md) for detailed HTTPS setup instructions.

Quick checks:

1. **Verify TLS is configured:**
   - Check `tls_cert_path` and `tls_key_path` in config.yaml
   - Or verify Let's Encrypt config

2. **Test certificate:**
   ```bash
   openssl s_client -connect headscale.example.com:443
   ```

3. **Check certificate expiry:**
   ```bash
   echo | openssl s_client -connect headscale.example.com:443 2>/dev/null | \
       openssl x509 -noout -dates
   ```

## General Troubleshooting Steps

### Enable debug logging

Add to your `config.yaml`:

```yaml
log:
  level: debug
  format: text
```

Restart Headscale and check logs:

```bash
systemctl restart headscale
journalctl -u headscale -f
```

### Check Headscale version

```bash
headscale version
```

Ensure you're running a version that supports password authentication (v0.30.0+).

### Review logs

```bash
# View recent logs
journalctl -u headscale -n 100

# Follow logs in real-time
journalctl -u headscale -f

# View logs for specific time range
journalctl -u headscale --since "2026-10-02 10:00:00"

# Search logs for errors
journalctl -u headscale | grep -i error
```

## Getting Help

If you've tried the solutions above and still have issues:

1. **Check the documentation:**
   - [Authentication Guide](usage/authentication.md)
   - [Configuration Reference](ref/configuration.md)
   - [API Reference](ref/api.md)

2. **Search existing issues:**
   - [Headscale GitHub Issues](https://github.com/arsydoni4326-alt/headscale/issues)

3. **Join the community:**
   - [Discord server](https://discord.gg/c84AZQhmpx) for real-time help

4. **File a bug report:**
   - Include Headscale version (`headscale version`)
   - Include relevant logs (with sensitive info redacted)
   - Include configuration (with passwords/secrets redacted)
   - Describe steps to reproduce the issue



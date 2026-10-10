# Settings

The Headplane Settings page allows you to customize your experience and manage your Headscale integration after logging in with your password.

## Overview

After authenticating to Headplane with your password (see [Authentication](authentication.md)), you can access the Settings page to:

- **Store your Headscale API key** — Save your API key securely so you don't need to re-enter it at each login
- **Change your password** — Update your Headplane login password with proper verification
- **Select a theme** — Choose between light and dark themes for the Headplane interface
- **Set a profile name** — Add an optional display name for personalization

All settings are stored securely and persist across login sessions.

---

## Accessing Settings

1. Log in to Headplane with your password
2. Click **Settings** in the navigation bar (top-right after login)
3. You'll see the Settings page with multiple sections

---

## Settings Sections

### Account

The Account section manages your Headplane authentication and session.

#### Change Password

Update your Headplane login password:

1. Navigate to Settings → **Account**
2. In the "Change Password" section:
   - Enter your **current password** for verification
   - Enter your **new password**
   - Confirm your **new password** by entering it again
3. Click **Change Password**
4. Upon success, you'll see a confirmation message
5. Your current session remains valid — you'll use the new password on your next login

**Security Notes:**

- The current password is required to prevent unauthorized password changes
- Password validation happens server-side with constant-time comparison to prevent timing attacks
- All existing sessions remain valid after password change (24-hour expiry unchanged)
- The new password takes effect immediately for future logins

**Troubleshooting:**

- **"Current password is incorrect"**: Verify you're entering your current password correctly
- **"New passwords do not match"**: Ensure both new password fields are identical
- **"Password change failed"**: Check server logs for configuration issues (see [Troubleshooting](#troubleshooting))

#### Session Information

View your current session details:

- **Session expires at**: Your session token expires 24 hours after login
- **Logged in since**: Timestamp of your current login
- **Logout**: Click to immediately terminate your session and return to the login page

---

### Integration

The Integration section manages your connection to the Headscale control server.

#### Headscale API Key

Store your Headscale API key to avoid re-entering it at each login.

**Why store an API key?**

Headplane needs a Headscale API key to display your tailnet information (nodes, users, routes, etc.). By default, you must enter this key at every login. Storing it in Settings eliminates this step.

**How to store your API key:**

1. **Obtain an API key from Headscale:**
   ```bash
   headscale apikeys create
   ```
   This generates an API key like: `hs_1234567890abcdef...`

2. **Navigate to Settings → Integration**

3. **Paste the API key** into the "Headscale API Key" field
   - The key is masked by default (shows `••••••••`)
   - Click the "eye" icon to reveal the key if needed

4. **Click "Save"**

5. **Verify success**: You'll see a confirmation message

6. **Your API key is now stored securely**
   - Encrypted using AES-256-GCM before storage
   - Automatically loaded at future logins
   - You won't need to enter it again unless you clear it or it expires

**Updating or removing an API key:**

- To **update**: Paste a new API key and click "Save"
- To **remove**: Clear the field and click "Save"

**Security Notes:**

- The API key is encrypted at rest using AES-256-GCM with a session-derived encryption key
- The key is never logged or exposed in server responses
- The stored key is decrypted only when you log in with your correct password
- If you change your Headplane password, re-save your API key to ensure it remains accessible

**Troubleshooting:**

- **"API key is invalid"**: Test the key directly with Headscale CLI to verify it's valid
- **"Failed to save API key"**: Check server logs for encryption errors
- **"Stored API key not loading"**: Re-save the API key after logging in

---

### Preferences

The Preferences section customizes your Headplane experience.

#### Theme

Choose between light and dark themes for the Headplane interface.

**Available themes:**

- **Light** — Clean, bright interface (default)
- **Dark** — Reduced eye strain for low-light environments

**How to change your theme:**

1. Navigate to Settings → **Preferences**
2. Select your preferred theme from the dropdown or toggle
3. The theme applies **immediately** (no save button required)
4. Your preference is saved automatically and persists across login sessions

**Technical notes:**

- Theme is stored server-side in your settings
- Theme is applied before page render to prevent flash of unstyled content
- Survives logout/login, browser cache clear, and device switches (as long as you log in with the same password)

---

### Profile

The Profile section manages your personal information and display preferences.

#### Display Name

Set an optional display name for your Headplane profile.

**How to set a display name:**

1. Navigate to Settings → **Profile**
2. Enter your desired name in the "Display Name" field (e.g., "John Doe", "Admin", "Network Team")
3. Click **Save**
4. Your display name will appear in:
   - The navigation bar (top-right)
   - Audit logs (if applicable)
   - Future multi-user features (Phase 13c)

**Notes:**

- Display name is **optional** — if not set, your session ID or username is shown
- Maximum length: 100 characters
- Can be changed at any time
- Stored server-side and persists across sessions

---

## Security Considerations

The Settings system is designed with security in mind:

### API Key Encryption

- **Encryption algorithm**: AES-256-GCM (Galois/Counter Mode)
- **Key derivation**: PBKDF2 with salt from session token
- **Storage**: Encrypted ciphertext + nonce stored in SQLite
- **Decryption**: Only occurs when you successfully authenticate with your password
- **Key rotation**: If you change your password, re-save your API key to maintain access

### Password Change Security

- **Current password verification**: Required to prevent unauthorized changes
- **Constant-time comparison**: Prevents timing attacks
- **Session preservation**: All active sessions remain valid after password change
- **Audit logging**: Password changes are logged for security monitoring

### CSRF Protection

- All settings endpoints are protected against Cross-Site Request Forgery attacks
- Session tokens are bound to your browser session
- Secure, HTTP-only cookies prevent JavaScript access

### Data Storage

- Settings are stored in Headplane's SQLite database (separate from Headscale)
- Single-row constraint enforces single-user model (Phase 13b)
- Automatic cleanup of expired sessions every 5 minutes

### Best Practices

1. **Use HTTPS in production**: Session tokens and API keys should never be transmitted over unencrypted connections
2. **Rotate API keys regularly**: Generate new Headscale API keys periodically and update Settings
3. **Use strong passwords**: Choose a strong, unique password for Headplane access
4. **Log out when done**: Explicitly log out on shared machines to invalidate your session
5. **Monitor audit logs**: Review logs for unexpected settings changes

---

## Troubleshooting

### Settings Page Not Loading

**Symptoms**: Settings page shows error or redirects to login

**Causes and solutions:**

- **Session expired**: Your 24-hour session has expired. Log in again.
- **Not logged in**: Settings requires password authentication. Log in first.
- **Server error**: Check Headscale logs for database or configuration errors.

### API Key Not Saving

**Symptoms**: "Failed to save API key" error or key doesn't persist

**Causes and solutions:**

1. **Database error**:
   - Check server logs for SQLite errors
   - Verify database file permissions: `headplane.db` should be writable by Headscale process
   - Ensure sufficient disk space

2. **Encryption error**:
   - Check server logs for "encryption failed" messages
   - Verify session token is valid (re-login if needed)

3. **Invalid API key format**:
   - Headscale API keys start with `hs_` prefix
   - Verify the key with: `headscale apikeys list`

### API Key Not Loading After Login

**Symptoms**: Must re-enter API key despite saving it

**Causes and solutions:**

1. **Password changed**:
   - If you recently changed your password, re-save your API key
   - Encryption key is derived from session token, which changes on password reset

2. **Database corruption**:
   - Check database integrity: `sqlite3 headplane.db "PRAGMA integrity_check;"`
   - Restore from backup if corruption detected

3. **API key expired or deleted**:
   - Verify the key still exists in Headscale: `headscale apikeys list`
   - Generate a new key if expired and update Settings

### Password Change Fails

**Symptoms**: "Current password is incorrect" or "Password change failed"

**Causes and solutions:**

1. **Wrong current password**:
   - Verify you're entering your current password correctly
   - Check for caps lock, leading/trailing spaces

2. **Configuration file not writable**:
   - Password is stored in `config.yaml` or environment variable
   - If using config file, verify file permissions allow writes by Headscale process
   - Consider using `HEADSCALE_HEADPLANE_PASSWORD` environment variable instead

3. **Environment variable override**:
   - If `HEADSCALE_HEADPLANE_PASSWORD` is set, config file changes are ignored
   - Unset environment variable or update it instead

### Theme Not Persisting

**Symptoms**: Theme resets to light after logout or page refresh

**Causes and solutions:**

1. **Not saved properly**:
   - Verify you see a success message after selecting theme
   - Check browser console for errors

2. **Database error**:
   - Check server logs for SQLite write errors
   - Verify database is writable

3. **Browser cache**:
   - Hard refresh: Ctrl+Shift+R (Windows/Linux) or Cmd+Shift+R (Mac)
   - Clear browser cache if theme styles seem corrupted

### Session Expired Unexpectedly

**Symptoms**: Logged out before 24 hours elapsed

**Causes and solutions:**

1. **Server restart**:
   - In-memory sessions are lost on server restart
   - This is expected behavior; log in again

2. **Clock skew**:
   - Verify server and client clocks are synchronized
   - Use NTP for time synchronization

3. **Session cleanup**:
   - Sessions are automatically cleaned up every 5 minutes
   - Expired sessions (>24 hours old) are purged

---

## Known Issues

### Phase 13c Updates

1. **Admin user management now available via API key** (Known Gap resolved):
   - Admin API keys can now access the full user management UI
   - The `/admin/users` page shows "Add Headplane User" and full CRUD operations when authenticated via admin API key
   - See [Authentication Guide](authentication.md#admin-user-management-with-api-keys) for details

2. **Multi-user support** (Phase 13c):
   - Headplane now supports multiple dashboard users with role-based access control
   - Each user has their own settings (API key, theme, profile)
   - Admin users can manage other users via the Settings interface or API

### Phase 13b Known Limitations

1. **Password stored in config file**:
   - Changing password requires file write access or environment variable update
   - On some deployments, this may require admin intervention
   - Consider using environment variable for easier password rotation

3. **API key re-save after password change**:
   - If you change your Headplane password, you may need to re-save your API key
   - This is due to the session-derived encryption key changing
   - Future versions may support automatic re-encryption

4. **Theme flash on slow connections**:
   - On very slow connections, you may briefly see the default theme before your saved theme loads
   - This is mitigated by inline theme application, but not entirely eliminated

5. **No API key rotation notification**:
   - If your Headscale API key expires, Headplane doesn't proactively notify you
   - You'll see API errors when trying to fetch tailnet data
   - Solution: Manually rotate keys and update Settings periodically

### Reporting Issues

If you encounter a bug or unexpected behavior:

1. Check the [Headscale logs](../../troubleshooting.md#checking-logs) for errors
2. Review this troubleshooting section for known solutions
3. Report issues on GitHub: [arsydoni4326-alt/headscale](https://github.com/arsydoni4326-alt/headscale/issues)
4. Include:
   - Headscale version (`headscale version`)
   - Steps to reproduce
   - Relevant log excerpts (redact secrets!)
   - Browser and OS information (for UI issues)

---

## Related Documentation

- [Authentication Guide](authentication.md) — Password authentication setup and login
- [API Reference: Settings](../ref/api/headplane-settings.md) — Technical API details for settings endpoints
- [Configuration Reference](../ref/configuration.md) — Headplane configuration options
- [Troubleshooting Guide](../../troubleshooting.md) — General troubleshooting for Headscale and Headplane
- [ROADMAP: Phase 13b](../../ROADMAP.md#phase-13b--settings-menu-single-user-in-progress) — Feature overview and implementation details




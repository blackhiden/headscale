# Authentication

Headscale provides multiple authentication methods for different use cases. This guide covers authentication for the Headplane web interface and programmatic API access.

## Authentication Methods

Headscale supports two primary authentication methods:

1. **Password Authentication** — Simple password-based login for Headplane web UI access
2. **API Keys** — Token-based authentication for programmatic API access and automation

### When to Use Each Method

- **Use Password Authentication** when:
    - Logging into the Headplane web interface
    - You need a simple, memorable credential for web UI access
    - You want to quickly grant access to administrators without managing API keys
    - You need to change your password or manage settings

- **Use API Keys** when:
    - Writing scripts or automation that calls the Headscale API
    - Integrating third-party tools with Headscale
    - You need fine-grained access control or token rotation
    - You want to access the admin UI without password login (admin API keys only)

!!! note "Admin API Keys"
    As of Phase 13c, **admin API keys can now access the full admin UI**, including user management. This means you can manage Headplane dashboard users directly from the UI when authenticated via an admin API key, without needing to log in with a Headplane password.

---

## Password Authentication (Headplane Web UI)

Headplane supports password-based authentication as a simplified login method for web interface access.

### How It Works

1. An administrator configures a password in the Headscale configuration file or via environment variable
2. Users navigate to the Headplane web interface
3. Users enter the configured password to log in
4. Upon successful authentication, Headplane issues a session token valid for 24 hours
5. The session token is stored securely in the browser and used for subsequent requests

### Configuration

The Headplane password can be configured in two ways:

#### Option 1: Configuration File

Add the `headplane.password` field to your `config.yaml`:

```yaml
# Headplane web UI configuration
headplane:
  # Password for Headplane web UI login.
  # ⚠️  SECURITY: Change this from the default in production!
  # This password grants access to the Headplane web interface only,
  # not the Headscale API. For API access, use API keys.
  password: "your-secure-password-here"
```

#### Option 2: Environment Variable

Set the `HEADSCALE_HEADPLANE_PASSWORD` environment variable:

```bash
export HEADSCALE_HEADPLANE_PASSWORD="your-secure-password-here"
```

!!! note "Precedence"
    If both the configuration file and environment variable are set, the **environment variable takes precedence**. This allows you to override the config file password without modifying the file itself.

### Logging In

1. Navigate to your Headplane URL (e.g., `https://headscale.example.com`)
2. On the login page, select **Password** as the authentication method
3. Enter the configured password
4. Click **Log In**
5. Upon successful authentication, you'll be redirected to the Headplane dashboard

### Session Management

- **Session Duration:** 24 hours from login
- **Session Storage:** Session tokens are stored securely in your browser
- **Session Expiry:** After 24 hours, you'll need to log in again
- **Automatic Cleanup:** Expired sessions are automatically cleaned up by the server

### Changing the Password

To change the Headplane password:

1. **Via Configuration File:**
    - Edit `config.yaml` and update the `headplane.password` field
    - Restart Headscale for the change to take effect
    - All existing sessions remain valid until expiry

2. **Via Environment Variable:**
    - Update the `HEADSCALE_HEADPLANE_PASSWORD` environment variable
    - Restart Headscale for the change to take effect
    - All existing sessions remain valid until expiry

!!! warning "Session Invalidation"
    Changing the password does **not** invalidate existing sessions. Users with active sessions can continue to use Headplane until their sessions expire (24 hours). To immediately revoke access, restart Headscale or wait for session expiry.

---

## Security Considerations

### Password Security

!!! danger "Production Deployment"
    **Always change the default password before deploying to production!**
    
    - Use a strong, randomly generated password (minimum 16 characters)
    - Never commit passwords to version control
    - Use environment variables or secure secret management for production deployments
    - Rotate passwords periodically according to your security policy

### HTTPS Requirement

!!! warning "Use HTTPS in Production"
    Password authentication transmits credentials over HTTP. **Always use HTTPS in production** to prevent credentials from being intercepted in transit.
    
    See the [TLS configuration guide](../ref/tls.md) for details on setting up HTTPS.

### Rate Limiting

Headscale implements rate limiting on the login endpoint to prevent brute-force attacks:

- **Limit:** 5 login attempts per minute per IP address
- **Behavior:** After exceeding the limit, subsequent requests receive HTTP 429 (Too Many Requests)
- **Reset:** The rate limit window resets after 1 minute

If you're experiencing rate limiting issues, see the [Troubleshooting](#troubleshooting) section below.

### Password Scope

The Headplane password **only grants access to the Headplane web interface**. It does not:

- Grant access to the Headscale REST API (use API keys for that)
- Allow command-line access to Headscale
- Provide node registration capabilities
- Allow programmatic automation

For API access and automation, use [API Keys](#api-keys-programmatic-access) instead.

---

## API Keys (Programmatic Access)

API keys provide token-based authentication for programmatic access to the Headscale REST API.

### Creating an API Key

To create an API key, log into your Headscale server and run:

```shell
headscale apikeys create
```

The command outputs an API key with a default expiration of 90 days. **Copy and save this key immediately** — you cannot retrieve it later.

### Using an API Key

Include the API key in the `Authorization` header of your HTTP requests:

```bash
curl -H "Authorization: Bearer YOUR_API_KEY" \
    https://headscale.example.com/api/v1/user
```

### Managing API Keys

List all API keys:

```shell
headscale apikeys list
```

Expire (revoke) an API key:

```shell
headscale apikeys expire --prefix <PREFIX>
```

!!! tip "API Key Best Practices"
    - Create separate API keys for each application or script
    - Rotate API keys regularly
    - Expire API keys immediately when they're no longer needed
    - Store API keys securely (environment variables, secret managers, not in code)

### API Documentation

For complete API reference, see:

- [API Reference](../ref/api.md) — Full endpoint documentation
- [OpenAPI Documentation](https://headscale.example.com/api/v1/docs) — Interactive API docs at `/api/v1/docs`

---

## Headplane Login API

The Headplane login endpoint is used internally by the Headplane web interface. You generally don't need to call this endpoint directly, but it's documented here for completeness.

### Endpoint

`POST /api/v1/headplane/login`

### Authentication

This endpoint **does not require** prior authentication (it's the authentication endpoint itself).

### Request

**Headers:**

```
Content-Type: application/json
```

**Body:**

```json
{
  "password": "your-headplane-password"
}
```

**Example:**

```bash
curl -X POST https://headscale.example.com/api/v1/headplane/login \
  -H "Content-Type: application/json" \
  -d '{"password": "your-secure-password"}'
```

### Response

**Success (200 OK):**

```json
{
  "token": "hp_AbCdEfGhIjKlMnOpQrStUvWxYz0123456789AbCdEfGhIjKlMnOpQrStUvWxYz",
  "expires_at": 1735689600
}
```

- `token` — Session token to use in subsequent requests (prefix: `hp_`)
- `expires_at` — Unix timestamp when the token expires (24 hours from now)

**Error (401 Unauthorized):**

```json
{
  "error": "Invalid password"
}
```

**Error (429 Too Many Requests):**

```json
{
  "error": "Rate limit exceeded. Please try again later."
}
```

### Using the Session Token

Include the session token in the `Authorization` header for Headplane web requests:

```bash
curl -H "Authorization: Bearer hp_AbCdEfGhIjKlMnOpQrStUvWxYz..." \
    https://headscale.example.com/api/v1/user
```

!!! note "Token Prefix"
    Session tokens begin with `hp_` to distinguish them from API keys. Both can be used in the `Authorization: Bearer` header, but session tokens are only valid for 24 hours.

---

## Troubleshooting

For common authentication issues and solutions, see the [Troubleshooting Guide](../../troubleshooting.md#authentication-issues).

### Quick Troubleshooting

**"Invalid password" error:**

- Verify the password matches your configuration file or environment variable exactly
- Check for leading/trailing whitespace in the password
- Ensure Headscale was restarted after changing the password
- Verify the configuration file path is correct

**Rate limiting (429 error):**

- Wait 1 minute and try again
- Check for automated scripts making excessive login attempts
- Verify your IP address is not behind a shared proxy that may trigger rate limits

**Session expired:**

- Sessions last 24 hours — log in again to get a new session
- Check your system clock is synchronized (sessions use timestamps)

**Password not working after configuration change:**

- Restart Headscale to load the new configuration
- Verify the `headplane.password` field is at the correct level in the YAML hierarchy
- Check Headscale logs for configuration parsing errors

For more detailed troubleshooting, see the [full troubleshooting guide](../../troubleshooting.md).

---

## Migration from API Key-Only Authentication

If you were previously using API keys for Headplane web access:

1. **Configure a password** using one of the methods above
2. **Restart Headscale** to enable password authentication
3. **Log in with the password** — API keys remain valid for API access
4. **Keep your API keys** for programmatic access and automation

Both authentication methods work simultaneously. Password authentication does not replace or invalidate API keys.

---

## Admin User Management with API Keys

As of Phase 13c, admin API keys can access the full Headplane admin interface, including user management capabilities. This resolves the "Known Gap" where admin features were only available to password-authenticated sessions.

### What You Can Do with Admin API Keys

When logged in to Headplane with an admin API key, you have full access to:

- **User Management UI** — The `/admin/users` page displays "Add Headplane User" button and full CRUD operations for dashboard users
- **Create Users** — Add new Headplane dashboard users with admin or regular user roles
- **Edit Users** — Update usernames, passwords, and roles for existing users
- **Delete Users** — Remove users from the system (with protection against deleting the last admin)
- **View All Users** — List all Headplane dashboard users with their roles and creation dates

### How to Access Admin UI with API Key

1. **Obtain an admin API key from Headscale:**
   ```bash
   headscale apikeys create
   ```

2. **Navigate to your Headplane URL** (e.g., `https://headscale.example.com`)

3. **On the login page, select "API Key" as the authentication method**

4. **Enter your API key** and click **Log In**

5. **Access admin features:**
   - Navigate to **Admin → Users** in the sidebar
   - You'll see the "Add Headplane User" button and full user management interface
   - Create, edit, or delete Headplane dashboard users as needed

### Admin API Key Requirements

To access admin features, your API key must:
- Be a valid Headscale API key (not expired or revoked)
- Have admin-level permissions in Headscale
- Be properly configured in your Headscale deployment

**Note**: Regular (non-admin) API keys will not have access to the admin user management interface.

### Use Cases

**Scenario 1: Automated Admin Workflows**

You can now script admin operations using API keys:

```bash
# List all Headplane users
curl -X GET https://headscale.example.com/api/v1/headplane/users \
  -H "Authorization: Bearer hs_your_admin_api_key"

# Create a new user
curl -X POST https://headscale.example.com/api/v1/headplane/users \
  -H "Authorization: Bearer hs_your_admin_api_key" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newoperator",
    "password": "secure-password",
    "role": "user"
  }'
```

**Scenario 2: Admin Without Password**

You can manage Headplane users entirely through API key authentication, without needing to configure or remember a Headplane password.

**Scenario 3: Consistent Auth Method**

If your team already uses API keys for automation, you can now use the same auth method for interactive admin tasks in the UI.

### Screenshots

![Admin Users Page - API Key Session](../assets/screenshots/admin-users-api-key.png)
*Admin users page showing "Add Headplane User" button when authenticated via API key*

**Note**: This screenshot is a placeholder reference. Actual screenshots should be captured from the implementation showing the unified admin UI for API key sessions.

---

## Settings and API Key Storage

After logging in with your password, you can save your Headscale API key in Settings to avoid re-entering it at each login:

1. Log in to Headplane with your password
2. Navigate to **Settings** in the navigation bar (top-right)
3. Go to **Settings → Integration**
4. Enter your Headscale API key (generate one with `headscale apikeys create`)
5. Click **Save**

Your API key will be stored securely (encrypted with AES-256-GCM) and reused across sessions. You can also change your password, select a theme preference, and set a display name in Settings.

See the [Settings Guide](settings.md) for complete documentation on managing your Headplane settings.

---

## Related Documentation


- [Configuration Reference](../ref/configuration.md) — Complete configuration options
- [API Reference](../ref/api.md) — REST API documentation
- [TLS Configuration](../ref/tls.md) — Setting up HTTPS
- [Troubleshooting Guide](../../troubleshooting.md) — Common issues and solutions


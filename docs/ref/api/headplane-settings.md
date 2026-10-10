# Headplane Settings API Reference

This document provides technical API documentation for the Headplane Settings endpoints introduced in Phase 13b. These endpoints allow authenticated users to manage their settings, including API key storage, password changes, theme preferences, and profile information.

## Overview

All settings endpoints require authentication via a valid Headplane session token obtained through password login (see [Authentication API](../authentication.md#login-endpoint)).

**Base URL**: `/api/v1/headplane`

**Authentication**: All endpoints require a valid session token in the `Authorization` header:
```
Authorization: Bearer hp_<session_token>
```

**Content-Type**: All requests and responses use `application/json`.

---

## Endpoints

### GET /api/v1/headplane/settings

Retrieve the current user's settings, including stored API key (decrypted), theme preference, and profile information.

#### Authentication

Required. Must include valid session token.

#### Request

**Method**: `GET`

**Headers**:
```
Authorization: Bearer hp_AbCdEfGhIjKlMnOpQrStUvWxYz0123456789AbCdEfGhIjKlMnOpQrStUvWxYz
```

**Parameters**: None

#### Response

**Success (200 OK)**:

Returns the user's current settings with the API key decrypted.

**Schema**:

```yaml
type: object
properties:
  api_key:
    type: string
    nullable: true
    description: Decrypted Headscale API key, or null if not set
    example: "hs_1234567890abcdefghijklmnopqrstuvwxyz"
  theme:
    type: string
    enum: [light, dark]
    description: Current theme preference
    default: light
    example: "dark"
  profile_name:
    type: string
    nullable: true
    maxLength: 100
    description: Optional display name for the user profile
    example: "John Doe"
  updated_at:
    type: string
    format: date-time
    description: ISO 8601 timestamp of last settings update
    example: "2026-10-03T04:52:04.978Z"
required:
  - theme
  - updated_at
```

**Example Response**:

```json
{
  "api_key": "hs_1234567890abcdefghijklmnopqrstuvwxyz",
  "theme": "dark",
  "profile_name": "Network Administrator",
  "updated_at": "2026-10-03T04:52:04.978Z"
}
```

**Example Response (No API Key Set)**:

```json
{
  "api_key": null,
  "theme": "light",
  "profile_name": null,
  "updated_at": "2026-10-01T12:00:00.000Z"
}
```

#### Error Responses

**401 Unauthorized**:

Session token is invalid, expired, or missing.

```yaml
type: object
properties:
  error:
    type: string
    example: "Unauthorized"
required:
  - error
```

```json
{
  "error": "Unauthorized"
}
```

**500 Internal Server Error**:

Server error (database failure, decryption error, etc.).

```yaml
type: object
properties:
  error:
    type: string
    example: "Failed to retrieve settings"
required:
  - error
```

```json
{
  "error": "Failed to retrieve settings"
}
```

#### Example Usage

```bash
curl -X GET https://headscale.example.com/api/v1/headplane/settings \
  -H "Authorization: Bearer hp_AbCdEfGhIjKlMnOpQrStUvWxYz0123456789..."
```

---

### POST /api/v1/headplane/settings

Update the user's settings. All fields are optional; only provided fields are updated.

#### Authentication

Required. Must include valid session token.

#### Request

**Method**: `POST`

**Headers**:
```
Authorization: Bearer hp_AbCdEfGhIjKlMnOpQrStUvWxYz0123456789AbCdEfGhIjKlMnOpQrStUvWxYz
Content-Type: application/json
```

**Body Schema**:

```yaml
type: object
properties:
  api_key:
    type: string
    nullable: true
    description: |
      Headscale API key to store (will be encrypted).
      Set to null or empty string to remove stored key.
    example: "hs_1234567890abcdefghijklmnopqrstuvwxyz"
  theme:
    type: string
    enum: [light, dark]
    description: Theme preference to save
    example: "dark"
  profile_name:
    type: string
    nullable: true
    maxLength: 100
    description: |
      Display name for user profile.
      Set to null or empty string to clear.
    example: "John Doe"
```

**Notes**:
- All fields are optional; send only the fields you want to update
- Empty string or `null` for `api_key` or `profile_name` clears the stored value
- API key is encrypted using AES-256-GCM before storage
- Theme changes apply immediately on next page load

**Example Request (Update API Key Only)**:

```json
{
  "api_key": "hs_1234567890abcdefghijklmnopqrstuvwxyz"
}
```

**Example Request (Update Theme and Profile Name)**:

```json
{
  "theme": "dark",
  "profile_name": "Network Administrator"
}
```

**Example Request (Clear API Key)**:

```json

#### Response

**Success (200 OK)**:

Settings updated successfully. Returns the updated settings (same format as GET).

**Schema**: Same as GET response schema (see above).

**Example Response**:

```json
{
  "api_key": "hs_1234567890abcdefghijklmnopqrstuvwxyz",
  "theme": "dark",
  "profile_name": "John Doe",
  "updated_at": "2026-10-03T04:52:04.978Z"
}
```

#### Error Responses

**400 Bad Request**:

Invalid input (e.g., invalid theme value, profile name too long).

```yaml
type: object
properties:
  error:
    type: string
    example: "Invalid theme value. Must be 'light' or 'dark'"
required:
  - error
```

```json
{
  "error": "Invalid theme value. Must be 'light' or 'dark'"
}
```

**401 Unauthorized**:

Session token is invalid, expired, or missing.

```yaml
type: object
properties:
  error:
    type: string
    example: "Unauthorized"
required:
  - error
```

```json
{
  "error": "Unauthorized"
}
```

**500 Internal Server Error**:

Server error (database failure, encryption error, etc.).

```yaml
type: object
properties:
  error:
    type: string
    example: "Failed to save settings"
required:
  - error
```

```json
{
  "error": "Failed to save settings"
}
```

#### Example Usage

```bash
# Update API key and theme
curl -X POST https://headscale.example.com/api/v1/headplane/settings \
  -H "Authorization: Bearer hp_AbCdEfGhIjKlMnOpQrStUvWxYz..." \
  -H "Content-Type: application/json" \
  -d '{
    "api_key": "hs_1234567890abcdefghijklmnopqrstuvwxyz",
    "theme": "dark"
  }'

# Clear API key
curl -X POST https://headscale.example.com/api/v1/headplane/settings \
  -H "Authorization: Bearer hp_AbCdEfGhIjKlMnOpQrStUvWxYz..." \
  -H "Content-Type: application/json" \
  -d '{
    "api_key": null
  }'
```

---

### POST /api/v1/headplane/change-password

Change the Headplane login password. Requires current password for verification.

#### Authentication

Required. Must include valid session token.

#### Request

**Method**: `POST`

**Headers**:
```
Authorization: Bearer hp_AbCdEfGhIjKlMnOpQrStUvWxYz0123456789AbCdEfGhIjKlMnOpQrStUvWxYz
Content-Type: application/json
```

**Body Schema**:

```yaml
type: object
properties:
  current_password:
    type: string
    description: Current Headplane password for verification
    minLength: 1
    example: "old-secure-password"
  new_password:
    type: string
    description: New password to set
    minLength: 1
    example: "new-secure-password"
required:
  - current_password
  - new_password
```

**Security Notes**:
- `current_password` is validated using constant-time comparison to prevent timing attacks
- Existing sessions remain valid after password change (24-hour expiry unchanged)
- Password is updated in the configuration source (file or environment variable)
- If using config file, the Headscale process must have write permissions

**Example Request**:

```json
{
  "current_password": "old-secure-password",
  "new_password": "new-secure-password"
}
```

#### Response

**Success (200 OK)**:

Password changed successfully. Current session remains valid.


#### Error Responses

**400 Bad Request**:

Missing required fields or invalid input.

```yaml
type: object
properties:
  error:
    type: string
    example: "Both current_password and new_password are required"
required:
  - error
```

```json
{
  "error": "Both current_password and new_password are required"
}
```

**401 Unauthorized**:

Either session token is invalid or current password is incorrect.

```yaml
type: object
properties:
  error:
    type: string
    example: "Current password is incorrect"
required:
  - error
```

**Example (Invalid Session)**:
```json
{
  "error": "Unauthorized"
}
```

**Example (Wrong Current Password)**:
```json
{
  "error": "Current password is incorrect"
}
```

**500 Internal Server Error**:

Server error (failed to write config file, permission denied, etc.).

```yaml
type: object
properties:
  error:
    type: string
    example: "Failed to update password"
required:
  - error
```

```json
{
  "error": "Failed to update password"
}
```

#### Example Usage

```bash
curl -X POST https://headscale.example.com/api/v1/headplane/change-password \
  -H "Authorization: Bearer hp_AbCdEfGhIjKlMnOpQrStUvWxYz..." \
  -H "Content-Type: application/json" \
  -d '{
    "current_password": "old-secure-password",
    "new_password": "new-secure-password"
  }'
```

---

## Security Considerations

### API Key Encryption

- **Algorithm**: AES-256-GCM (Galois/Counter Mode)
- **Key Derivation**: PBKDF2-HMAC-SHA256 with salt derived from session token
- **Storage**: Encrypted ciphertext and nonce stored in SQLite database
- **Decryption**: Occurs server-side only when valid session token is presented
- **Key Rotation**: If password changes, re-save API key to ensure continued access

### Password Change Security

- **Current Password Verification**: Required to prevent unauthorized changes
- **Constant-Time Comparison**: Prevents timing attacks during verification
- **Session Preservation**: All active sessions remain valid after password change
- **Audit Logging**: Password changes should be logged for security monitoring
- **Configuration Persistence**: Password is written to config file or environment variable

### CSRF Protection

All settings endpoints are protected against Cross-Site Request Forgery attacks through:
- Session token validation
- Secure, HTTP-only cookies
- Origin/referer header validation (if configured)

### Transport Security

- **HTTPS Required**: Use HTTPS in production to protect session tokens and API keys in transit
- **Secure Cookies**: Session cookies should have `Secure` flag set (automatic with HTTPS)
- **No Caching**: Settings responses include `Cache-Control: no-store` headers

### Rate Limiting

Settings endpoints may be rate-limited to prevent abuse:
- **Recommended Limit**: 10 requests per minute per session
- **Password Change**: 5 attempts per hour per session (stricter for security)

---

## Database Schema

Settings are stored per Headplane user in the Headscale database:

```sql
CREATE TABLE headplane_settings (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  api_key_encrypted TEXT,
  api_key_nonce TEXT,
  api_key_salt TEXT,
  theme TEXT DEFAULT "light",
  profile_name TEXT,
  updated_at DATETIME
);
CREATE UNIQUE INDEX idx_headplane_settings_user_id ON headplane_settings(user_id);
```

**Notes**:
- The unique index allows one settings record per Headplane user.
- `api_key_encrypted`, `api_key_nonce`, and `api_key_salt` are stored together when an API key is set.
- `theme` defaults to `'light'` if not set.
- `profile_name` is optional.

---

## Testing

### Unit Tests

Backend unit tests should cover:
- API key encryption/decryption
- Settings CRUD operations
- Password change validation
- Error handling (invalid input, missing fields, etc.)
- Session authentication

### Integration Tests

End-to-end tests should verify:
1. **API Key Persistence**: Save API key → logout → login → verify key is loaded
2. **Password Change Flow**: Change password → logout → login with new password
3. **Theme Persistence**: Set theme → logout → login → verify theme is applied
4. **Profile Name**: Set profile name → verify it appears in UI
5. **Security**: Attempt unauthorized access → verify 401 response
6. **Encryption**: Verify API key is encrypted in database (not plaintext)

### Example Test Scenarios

See `docs/phase13b-integration-testing.md` for detailed integration test scenarios and results (created after backend/frontend implementation is complete).

---

## Related Documentation

- [Settings User Guide](../../usage/settings.md) — End-user documentation for the Settings page
- [Authentication API](../authentication.md) — Password login endpoint documentation
- [Headplane Users API](./headplane-users.md) — User management endpoints (admin only)
- [Configuration Reference](../configuration.md) — Headplane configuration options
- [Phase 13b Implementation Plan](../../../docs/phase13b-implementation-plan.md) — Technical implementation details

---

## Changelog

- **2026-10-03**: Phase 13b — Initial settings API documentation created
- **[To be updated]**: After backend implementation, update with actual endpoint paths and behavior

**Schema**:

```yaml
type: object
properties:
  message:
    type: string
    example: "Password changed successfully"
required:
  - message
```

**Example Response**:

```json
{
  "message": "Password changed successfully"
}
```

{
  "api_key": null
}
```

**Example Request (Update All Settings)**:

```json
{
  "api_key": "hs_1234567890abcdefghijklmnopqrstuvwxyz",
  "theme": "dark",
  "profile_name": "John Doe"
}
```


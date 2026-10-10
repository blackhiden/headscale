# QR Code Payload Format

Registration QR codes encode a JSON payload:

## Structure

```json
{
  "type": "headscale-registration",
  "version": "1",
  "auth_id": "hskey-...",
  "server_url": "https://headscale.example.com",
  "expires_at": "2026-10-08T07:00:00Z"
}
```

## Fields

- `type`: Always "headscale-registration"
- `version`: Payload format version (currently "1")
- `auth_id`: Registration authentication ID (hskey)
- `server_url`: Base URL of Headscale instance
- `expires_at`: Expiration timestamp (ISO 8601)

## Security

- No secrets or API keys included
- Single-use: invalidated after successful registration
- Time-limited: expires with registration session

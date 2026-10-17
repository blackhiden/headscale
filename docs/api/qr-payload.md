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

## API Endpoint

You can retrieve the QR registration payload programmatically using the API endpoint:

```bash
GET /api/v1/registration/qr/{auth_id}
```

### Example Request

```bash
curl https://headscale.example.com/api/v1/registration/qr/hskey-authreq-abc123
```

### Example Response (200 OK)

```json
{
  "type": "headscale-registration",
  "version": "1",
  "auth_id": "hskey-authreq-abc123",
  "server_url": "https://headscale.example.com",
  "expires_at": "2026-10-17T13:53:34.387885+07:00"
}
```

### Error Responses

- **400 Bad Request**: Invalid `auth_id` format
- **404 Not Found**: Pending registration not found or not a registration request
- **410 Gone**: Registration has expired

### Notes

- This endpoint is **open** (no authentication required)
- The `expires_at` value matches the pending registration cache deadline
- The payload is only available while the registration is pending and not expired
- Once the registration is completed or expires, the endpoint returns 404 or 410

### Use Case: CLI Client QR Generation

CLI clients can use this endpoint to fetch the payload and generate QR codes locally:

```bash
# Fetch payload
PAYLOAD=$(curl -s https://headscale.example.com/api/v1/registration/qr/hskey-authreq-abc123)

# Generate QR code (using qrencode or similar tool)
echo "$PAYLOAD" | qrencode -t ansiutf8
```

## Security

- No secrets or API keys included
- Single-use: invalidated after successful registration
- Time-limited: expires with registration session
- Open endpoint: anyone with an `auth_id` can retrieve the payload (same as `/register/{auth_id}` page)

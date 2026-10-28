# Server Info API

The Server Info API provides system and runtime specifications for the Headscale control server and the Headplane web interface.

## Headscale Server Info

### Endpoint

```
GET /api/v1/server-info
```

### Authentication

Requires API key authentication via HTTP Bearer token:

```
Authorization: Bearer <API_KEY>
```

### Response

Returns system information for the Headscale server.

#### Success Response (200 OK)

```json
{
  "os": "linux",
  "arch": "amd64",
  "goVersion": "go1.27.2",
  "cpuCores": 8,
  "ramBytes": 16777216000,
  "storageBytes": 107374182400,
  "dockerDetected": true,
  "dockerVersion": "24.0.7",
  "publicIP": "203.0.113.42",
  "tailnetIPv4": "100.64.0.0/10",
  "tailnetIPv6": "fd7a:115c:a1e0::/48"
}
```

#### Response Fields

| Field | Type | Description |
| --- | --- | --- |
| `os` | string | Operating system (e.g., `linux`, `darwin`, `windows`) |
| `arch` | string | CPU architecture (e.g., `amd64`, `arm64`) |
| `goVersion` | string | Go runtime version |
| `cpuCores` | integer | Number of logical CPU cores |
| `ramBytes` | integer | Total system memory in bytes |
| `storageBytes` | integer | Total disk space in bytes |
| `dockerDetected` | boolean | Whether Docker is detected |
| `dockerVersion` | string \| null | Docker version if detected, otherwise null |
| `publicIP` | string \| null | Public IP address, or null if unavailable |
| `tailnetIPv4` | string | IPv4 CIDR range allocated to the tailnet |
| `tailnetIPv6` | string \| null | IPv6 CIDR range allocated to the tailnet, or null if IPv6 is not configured |

### Example Request

```bash
curl -H "Authorization: Bearer <API_KEY>" \
    https://headscale.example.com/api/v1/server-info
```

### Error Responses

- **401 Unauthorized**: Missing or invalid API key
- **500 Internal Server Error**: Server error while gathering information

### Notes

- This endpoint is **read-only** and does not modify server state.
- Fields that cannot be determined return `null` rather than failing the entire request.
- The public IP lookup is non-blocking and may return `null` if the IP echo service is unreachable.

---

## Headplane Server Info

### Endpoint

```
GET /api/admin/server-info
```

### Authentication

Requires administrator session authentication (logged in as admin).

### Response

Returns combined system information for both Headscale and Headplane.

#### Success Response (200 OK)

```json
{
  "headscale": {
    "os": "linux",
    "arch": "amd64",
    "goVersion": "go1.27.2",
    "cpuCores": 8,
    "ramBytes": 16777216000,
    "storageBytes": 107374182400,
    "dockerDetected": true,
    "dockerVersion": "24.0.7",
    "publicIP": "203.0.113.42",
    "tailnetIPv4": "100.64.0.0/10",
    "tailnetIPv6": "fd7a:115c:a1e0::/48"
  },
  "headplane": {
    "os": "linux",
    "arch": "x64",
    "nodeVersion": "v20.10.0",
    "cpuCores": 8,
    "ramBytes": 16777216000,
    "storageBytes": 107374182400,
    "dockerDetected": true,
    "dockerVersion": "24.0.7"
  }
}
```

#### Response Fields

**`headscale` object:** Same structure as the Headscale Server Info endpoint above.

**`headplane` object:**

| Field | Type | Description |
| --- | --- | --- |
| `os` | string | Operating system |
| `arch` | string | CPU architecture |
| `nodeVersion` | string | Node.js runtime version |
| `cpuCores` | integer | Number of logical CPU cores |
| `ramBytes` | integer | Total system memory in bytes |
| `storageBytes` | integer | Total disk space in bytes |
| `dockerDetected` | boolean | Whether Docker is detected |
| `dockerVersion` | string \| null | Docker version if detected, otherwise null |

### Example Request

```bash
curl -H "Cookie: session=<SESSION_COOKIE>" \
    https://headscale.example.com/api/admin/server-info
```

### Error Responses

- **401 Unauthorized**: Not logged in or not an administrator
- **500 Internal Server Error**: Server error while gathering information
- **502 Bad Gateway**: Headscale server is unreachable (Headplane info still returned if available)

### Graceful Degradation

If the Headscale server is unreachable, the endpoint still returns a 200 response with the `headplane` object populated and the `headscale` object set to `null` or an error indicator. This allows the UI to display partial information rather than failing entirely.

```json
{
  "headscale": null,
  "headplane": {
    "os": "linux",
    "arch": "x64",
    "nodeVersion": "v20.10.0",
    "cpuCores": 8,
    "ramBytes": 16777216000,
    "storageBytes": 107374182400,
    "dockerDetected": false,
    "dockerVersion": null
  },
  "error": "Headscale server unreachable"
}
```

### Notes

- This endpoint proxies the Headscale `/api/v1/server-info` endpoint and merges the result with Headplane's own server info.
- Admin authentication is required; regular users cannot access this endpoint.
- Network failures to Headscale are handled gracefully to ensure the UI remains functional.

---

## Security

- **No sensitive information is exposed.** API keys, passwords, tokens, private keys, and detailed hardware identifiers are never included in the response.
- **Admin-only access.** Both endpoints require administrator privileges (API key for Headscale, session authentication for Headplane).
- **Read-only operations.** These endpoints do not modify server state or configuration.

## Use Cases

- **Troubleshooting:** Programmatically retrieve server specs for diagnostics and support requests.
- **Monitoring:** Integrate server info into monitoring dashboards or alerting systems.
- **Capacity planning:** Automate collection of resource utilization data for infrastructure planning.

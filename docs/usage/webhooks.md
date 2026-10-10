# Webhooks

Headscale supports webhooks for monitoring and alerting integrations. Webhooks allow you to receive HTTP POST notifications when specific events occur in your Headscale instance.

## Overview

Webhooks enable real-time notifications for events like:

- **Node up/down** - When nodes connect or disconnect
- **Health check failures** - When health checks fail
- **Alert firing** - When external alerts are triggered

Each webhook can be configured with:

- Custom HTTP headers for authentication
- HMAC-SHA256 payload signing for security
- Configurable timeout and enabled/disabled state
- Event filtering to receive only specific event types

## Configuration

Webhooks are managed via the REST API at `/api/v1/webhooks`. All webhook operations require admin API key authentication.

### Event Types

| Event Type | Description |
|------------|-------------|
| `node_up` | Fired when a node comes online |
| `node_down` | Fired when a node goes offline |
| `health_check_fail` | Fired when a health check fails |
| `alert_firing` | Fired when an external alert is received |

## API Operations

### Create Webhook

**Endpoint:** `POST /api/v1/webhooks`

**Request Body:**

```json
{
  "name": "prometheus-alerts",
  "url": "https://alertmanager.example.com/webhook",
  "events": ["node_up", "node_down", "health_check_fail"],
  "headers": {
    "Authorization": "Bearer your-token"
  },
  "secret": "your-signing-secret",
  "enabled": true,
  "timeoutSeconds": 10
}
```

**Response:**

```json
{
  "webhook": {
    "id": 1,
    "name": "prometheus-alerts",
    "url": "https://alertmanager.example.com/webhook",
    "events": ["node_up", "node_down", "health_check_fail"],
    "enabled": true,
    "timeoutSeconds": 10,
    "createdAt": "2026-09-30T03:43:34Z",
    "updatedAt": "2026-09-30T03:43:34Z"
  }
}
```

**Note:** The `secret` is not returned in responses for security.

### List Webhooks

**Endpoint:** `GET /api/v1/webhooks`

### Get Webhook

**Endpoint:** `GET /api/v1/webhooks/{id}`

### Update Webhook

**Endpoint:** `PUT /api/v1/webhooks/{id}`

**Request Body:** (all fields optional)

```json
{
  "name": "updated-name",
  "enabled": false
}
```

### Delete Webhook

**Endpoint:** `DELETE /api/v1/webhooks/{id}`

## Webhook Payload

All webhooks receive a JSON payload with the following structure:

```json
{
  "event_type": "node_up",
  "timestamp": "2026-09-30T03:43:34.873Z",
  "data": {
    "node_id": 123,
    "node_name": "my-laptop",
    "user_name": "alice"
  }
}
```

### Event Data Schemas

#### Node Events (`node_up`, `node_down`)

```json
{
  "node_id": 123,
  "node_name": "my-laptop",
  "user_name": "alice"
}
```

#### Health Check Events (`health_check_fail`)

```json
{
  "endpoint": "/health",
  "error": "database connection failed"
}
```

#### Alert Events (`alert_firing`)

```json
{
  "alert_name": "HighErrorRate",
  "severity": "critical",
  "labels": {
    "instance": "headscale-1"
  },
  "details": {
    "error_rate": 0.15
  }
}
```

## Security

### HMAC Signature

When a webhook has a `secret` configured, all payloads are signed with HMAC-SHA256. The signature is sent in the `X-Headscale-Signature` header.

**Verifying signatures in Python:**

```python
import hmac
import hashlib

def verify_webhook(payload, signature, secret):
    expected_signature = hmac.new(
        secret.encode('utf-8'),
        payload.encode('utf-8'),
        hashlib.sha256
    ).hexdigest()
    
    return hmac.compare_digest(signature, expected_signature)
```

**Verifying signatures in Go:**

```go
func verifyWebhook(payload []byte, signature, secret string) bool {
    h := hmac.New(sha256.New, []byte(secret))
    h.Write(payload)
    expectedSignature := hex.EncodeToString(h.Sum(nil))
    
    return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
```

### Custom Headers

Use the `headers` field to include authentication tokens:

```json
{
  "headers": {
    "Authorization": "Bearer your-token",
    "X-Custom-Header": "custom-value"
  }
}
```

## Integration Examples

### Prometheus Alertmanager

**alertmanager.yml:**

```yaml
receivers:
  - name: headscale
    webhook_configs:
      - url: 'https://your-handler.example.com/alerts'
        http_config:
          authorization:
            credentials: 'your-api-key'
```

### Grafana

Configure Grafana alerts with a webhook contact point pointing to your webhook handler.

## Metrics

Webhook operations expose Prometheus metrics:

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `headscale_webhook_dispatch_total` | Counter | Total webhook dispatches | `event_type`, `status` |
| `headscale_webhook_dispatch_duration_seconds` | Histogram | Webhook dispatch duration | `event_type` |

**Status labels:** `success`, `http_error`, `send_error`, `request_error`, `marshal_error`, `list_error`

## Best Practices

1. **Use HMAC signatures** - Always configure a `secret` to verify webhook authenticity
2. **Set appropriate timeouts** - Default is 10 seconds
3. **Monitor webhook metrics** - Track success/failure rates via Prometheus
4. **Use HTTPS** - Configure webhook URLs with HTTPS
5. **Handle retries externally** - Headscale does not retry failed webhooks

## See Also

- [Observability](observability.md) - Prometheus metrics and monitoring
- [API Reference](../ref/api.md) - Complete API documentation
- [Extending Headscale](../ref/extending.md) - Integration points

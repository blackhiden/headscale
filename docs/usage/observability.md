# Observability

This document describes the observability features available in Headscale, including metrics, logging, and operational endpoints.

## Operational Endpoints

Headscale exposes several endpoints for health checking, monitoring, and operational purposes.

### Health Endpoints

#### `/health` - Liveness Check

**Purpose:** Basic health check that verifies the server is running and the database is reachable.

**Authentication:** None (publicly accessible)

**Response Format:** `application/health+json`

**Success Response (200 OK):**
```json
{
  "status": "pass"
}
```

**Failure Response (500 Internal Server Error):**
```json
{
  "status": "fail"
}
```

**Use Case:** Use this endpoint for liveness probes in container orchestration platforms (Kubernetes `livenessProbe`, Docker healthchecks).

**Example:**
```bash
curl http://localhost:8080/health
```

#### `/ready` - Readiness Check

**Purpose:** Indicates whether the server is ready to accept traffic. Checks database connectivity and returns 503 if not ready.

**Authentication:** None (publicly accessible)

**Response Format:** `application/json`

**Success Response (200 OK):**
```json
{
  "ready": true,
  "status": "ready"
}
```

**Not Ready Response (503 Service Unavailable):**
```json
{
  "ready": false,
  "status": "not ready"
}
```

**Use Case:** Use this endpoint for readiness probes in container orchestration platforms (Kubernetes `readinessProbe`). The server returns 503 during startup or when the database is unreachable, signaling load balancers and orchestrators to temporarily stop routing traffic.

**Example:**
```bash
curl http://localhost:8080/ready
```

**Kubernetes Example:**
```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
```

#### `/version` - Version Information

**Purpose:** Returns version, commit, and build information about the running Headscale server.

**Authentication:** None (publicly accessible)

**Response Format:** `application/json`

**Example Response:**
```json
{
  "version": "0.34.0-arsydoni4326-alt",
  "commit": "abc123def456",
  "buildDate": "2026-09-29T06:00:00Z",
  "dirty": false
}
```

**Example:**
```bash
curl http://localhost:8080/version
```

#### `/api/v1/health` - Authenticated Health Check

**Purpose:** Detailed health check for authenticated monitoring systems.

**Authentication:** Required (API key via Bearer token)

**Response Format:** `application/json`

**Example Response:**
```json
{
  "databaseConnectivity": true
}
```

**Example:**
```bash
curl -H "Authorization: Bearer YOUR_API_KEY" http://localhost:8080/api/v1/health
```

## Metrics

Headscale exposes Prometheus metrics on a separate listener configured via `metrics_listen_addr` in the configuration file.

### Configuration

```yaml
# Address to listen to /metrics and /debug
# Use an empty value to disable the metrics listener
metrics_listen_addr: 127.0.0.1:9090
```

**Security Note:** The metrics endpoint is protected and only accessible from:
- Loopback addresses (127.0.0.1, ::1)
- Tailscale CGNAT IPs (100.64.0.0/10)
- Private network addresses (RFC 1918 / RFC 4193)

### Available Metrics

Headscale exports metrics in the Prometheus exposition format at `http://<metrics_listen_addr>/metrics`.

#### HTTP Metrics

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `http_requests_total` | Counter | Total number of HTTP requests | `method`, `status`, `proto` |
| `http_request_duration_seconds` | Histogram | HTTP request duration in seconds | `method`, `status`, `proto` |

**Note:** `OPTIONS` requests are excluded from HTTP metrics to reduce noise.

#### MapResponse Metrics

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `headscale_mapresponse_sent_total` | Counter | Total count of MapResponses sent to clients | `status`, `type` |
| `headscale_mapresponse_generated_total` | Counter | Total count of MapResponses generated | `reason` |
| `headscale_mapresponse_endpoint_updates_total` | Counter | Total count of endpoint updates received | `status` |
| `headscale_mapresponse_ended_total` | Counter | Total count of map sessions ended | `reason` |
| `headscale_mapresponse_last_sent_seconds` | Gauge | Last sent metric to node (high cardinality, requires `HEADSCALE_DEBUG_HIGH_CARDINALITY_METRICS=1`) | `type`, `id` |

#### NodeStore Metrics

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `headscale_nodestore_operations_total` | Counter | Total number of NodeStore operations | `operation` |
| `headscale_nodestore_operation_duration_seconds` | Histogram | NodeStore operation duration | `operation` |
| `headscale_nodestore_batch_size` | Histogram | Number of nodes processed in batch operations | - |
| `headscale_nodestore_batch_duration_seconds` | Histogram | Batch operation duration | - |
| `headscale_nodestore_snapshot_build_duration_seconds` | Histogram | Snapshot build duration | - |
| `headscale_nodestore_snapshot_builds_total` | Counter | Total snapshot builds | `trigger` |
| `headscale_nodestore_nodes` | Gauge | Number of nodes in the NodeStore | - |
| `headscale_nodestore_peers_calculation_duration_seconds` | Histogram | Peer visibility calculation duration | - |
| `headscale_nodestore_queue_depth` | Gauge | Current queue depth | - |

#### Mapper Metrics

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `headscale_mapper_changes_dropped_total` | Counter | Changes the batcher refused to fan out | `reason` |

#### Update Check Metrics (Fork-Specific)

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `headscale_updatecheck_requests_total` | Counter | Total update check requests | `check` |
| `headscale_updatecheck_remote_failures_total` | Counter | Remote GitHub API failures | `reason` |
| `headscale_updatecheck_cache_hits_total` | Counter | Cache hits | - |
| `headscale_updatecheck_cache_misses_total` | Counter | Cache misses | - |

#### HA Health Probe Metrics

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `headscale_ha_health_updates_total` | Counter | Health probe outcomes for HA subnet routers | `outcome` |

### Querying Metrics

**Example: Scrape all metrics**
```bash
curl http://127.0.0.1:9090/metrics
```

**Example: Filter specific metrics**
```bash
curl http://127.0.0.1:9090/metrics | grep nodestore
```

**Prometheus Configuration:**
```yaml
scrape_configs:
  - job_name: 'headscale'
    static_configs:
      - targets: ['localhost:9090']
```

## Logging

Headscale uses structured logging via [zerolog](https://github.com/rs/zerolog).

### Log Levels

Log output can be controlled via the `log.level` configuration option:

```yaml
log:
  level: info  # trace, debug, info, warn, error, fatal, panic
  format: text  # text or json
```

### Structured Fields

Logs include structured fields for filtering and analysis:

- `level`: Log level (info, warn, error, etc.)
- `time`: Timestamp (RFC3339)
- `message`: Log message
- `error`: Error details (when applicable)
- Context-specific fields: `node.id`, `user`, `version`, `commit`, etc.

### Example Log Output (text format)

```
2026-09-29T06:27:53Z INF starting headscale version=0.34.0-arsydoni4326-alt commit=abc123
2026-09-29T06:27:53Z INF Clients with a lower minimum version will be rejected minimum_version=v1.82.0
2026-09-29T06:27:53Z INF HTTP server listening on 127.0.0.1:8080
2026-09-29T06:27:53Z INF Metrics server listening on 127.0.0.1:9090
```

### Example Log Output (json format)

```json
{"level":"info","time":"2026-09-29T06:27:53Z","message":"starting headscale","version":"0.34.0-arsydoni4326-alt","commit":"abc123"}
{"level":"info","time":"2026-09-29T06:27:53Z","message":"HTTP server listening on 127.0.0.1:8080"}
```

### Logging Best Practices

1. **Use structured fields** instead of string concatenation for easier filtering
2. **Set appropriate log levels** (use `debug` for development, `info` for production)
3. **Rotate logs** using your OS log rotation mechanism (logrotate, journald)
4. **Monitor error logs** for operational issues

## Debug Endpoints

Additional debug endpoints are available on the metrics listener (`metrics_listen_addr`).

### `/debug/`

Interactive debug interface with links to all available debug endpoints.

**Access:** Protected (loopback, Tailscale IPs, private networks only)

### `/debug/overview`

Current state overview including node count, user count, and policy summary.

**Example:**
```bash
curl http://127.0.0.1:9090/debug/overview
```

### `/debug/config`

Current server configuration (sensitive values redacted).

**Example:**
```bash
curl http://127.0.0.1:9090/debug/config
```

### `/debug/policy`

Current policy (ACLs) loaded in the server.

**Example:**
```bash
curl http://127.0.0.1:9090/debug/policy
```

### `/debug/pprof/`

Go runtime profiling endpoints for performance analysis.

**Example:**
```bash
# CPU profile
curl http://127.0.0.1:9090/debug/pprof/profile?seconds=30 > cpu.pprof

# Heap profile
curl http://127.0.0.1:9090/debug/pprof/heap > heap.pprof

# Analyze with pprof
go tool pprof cpu.pprof
```

### `/debug/statsviz`

Live visualization of Go runtime statistics (memory, GC, goroutines).

**Access:** Visit `http://127.0.0.1:9090/debug/statsviz` in your browser.

## Alerting

### Recommended Prometheus Alerts

```yaml
groups:
  - name: headscale
    rules:
      - alert: HeadscaleDown
        expr: up{job="headscale"} == 0
        for: 5m
        annotations:
          summary: "Headscale is down"

      - alert: HeadscaleDatabaseUnreachable
        expr: headscale_nodestore_operation_duration_seconds{operation="ping"} > 5
        for: 5m
        annotations:
          summary: "Headscale database is slow or unreachable"

      - alert: HeadscaleHighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
        for: 5m
        annotations:
          summary: "Headscale is returning a high rate of 5xx errors"

      - alert: HeadscaleMapResponseFailures
        expr: rate(headscale_mapresponse_sent_total{status="error"}[5m]) > 0.05
        for: 5m
        annotations:
          summary: "Headscale is failing to send map responses"
```

### Grafana Dashboard

The following Grafana dashboard JSON visualises the key Headscale metrics.
Import it in Grafana via **Dashboards → Import** and paste the JSON, or save it
as a file and use the Grafana provisioning API.

```json
{
  "title": "Headscale",
  "uid": "headscale-overview",
  "tags": ["headscale"],
  "timezone": "browser",
  "panels": [
    {
      "title": "HTTP Requests",
      "type": "timeseries",
      "datasource": { "type": "prometheus", "uid": "prometheus" },
      "targets": [
        {
          "expr": "sum(rate(http_requests_total[5m])) by (code)",
          "legendFormat": "{{code}}"
        }
      ],
      "fieldConfig": {
        "defaults": { "unit": "reqps" }
      }
    },
    {
      "title": "HTTP Request Duration (p95)",
      "type": "timeseries",
      "datasource": { "type": "prometheus", "uid": "prometheus" },
      "targets": [
        {
          "expr": "histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))",
          "legendFormat": "p95"
        }
      ],
      "fieldConfig": {
        "defaults": { "unit": "s" }
      }
    },
    {
      "title": "Map Responses Sent",
      "type": "timeseries",
      "datasource": { "type": "prometheus", "uid": "prometheus" },
      "targets": [
        {
          "expr": "sum(rate(headscale_mapresponse_sent_total[5m])) by (status)",
          "legendFormat": "{{status}}"
        }
      ],
      "fieldConfig": {
        "defaults": { "unit": "ops" }
      }
    },
    {
      "title": "Connected Nodes",
      "type": "stat",
      "datasource": { "type": "prometheus", "uid": "prometheus" },
      "targets": [
        {
          "expr": "headscale_nodestore_nodes_total",
          "legendFormat": "nodes"
        }
      ],
      "fieldConfig": {
        "defaults": { "unit": "short" }
      }
    },
    {
      "title": "NodeStore Operation Duration (p95)",
      "type": "timeseries",
      "datasource": { "type": "prometheus", "uid": "prometheus" },
      "targets": [
        {
          "expr": "histogram_quantile(0.95, sum(rate(headscale_nodestore_operation_duration_seconds_bucket[5m])) by (le, operation))",
          "legendFormat": "{{operation}}"
        }
      ],
      "fieldConfig": {
        "defaults": { "unit": "s" }
      }
    },
    {
      "title": "Update Check Requests",
      "type": "timeseries",
      "datasource": { "type": "prometheus", "uid": "prometheus" },
      "targets": [
        {
          "expr": "sum(rate(headscale_updatecheck_requests_total[5m])) by (check)",
          "legendFormat": "check={{check}}"
        }
      ],
      "fieldConfig": {
        "defaults": { "unit": "reqps" }
      }
    }
  ],
  "refresh": "30s",
  "schemaVersion": 39,
  "version": 1
}
```

The dashboard assumes a Prometheus datasource with UID `prometheus`. Adjust the
datasource UID to match your Grafana setup.

## Troubleshooting

### Metrics Not Available

**Symptom:** `/metrics` endpoint returns 404 or connection refused.

**Solutions:**
1. Check `metrics_listen_addr` is configured in `config.yaml`
2. Verify the metrics server is running: check logs for "Metrics server listening"
3. Ensure you're connecting from an allowed address (loopback, Tailscale IP, or private network)

### Health Check Fails

**Symptom:** `/health` or `/ready` returns `fail` or `not ready`.

**Solutions:**
1. Check database connectivity: `headscale nodes list` (should succeed if DB is reachable)
2. Check database configuration in `config.yaml`
3. Check database logs for connection errors
4. Verify database is running and accessible

### High Memory Usage

**Solutions:**
1. Check `/debug/statsviz` for memory trends
2. Capture a heap profile: `curl http://127.0.0.1:9090/debug/pprof/heap > heap.pprof`
3. Analyze with `go tool pprof heap.pprof`
4. Review `headscale_nodestore_nodes` metric for node count

### Slow Responses

**Solutions:**
1. Check `http_request_duration_seconds` histogram for latency distribution
2. Check `headscale_nodestore_operation_duration_seconds` for database performance
3. Check `headscale_nodestore_snapshot_build_duration_seconds` for NodeStore performance
4. Capture a CPU profile: `curl 'http://127.0.0.1:9090/debug/pprof/profile?seconds=30' > cpu.pprof`

## See Also

- [Configuration Reference](../ref/configuration.md)
- [Debug Endpoints](../ref/debug.md)
- [Webhooks](webhooks.md) - Webhook notifications for monitoring and alerting
- [Getting Started](getting-started.md)

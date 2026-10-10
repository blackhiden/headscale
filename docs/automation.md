# Automation Support

Headscale's v2 API provides OAuth 2.0 client-credentials authentication, enabling automation tools like the Tailscale Terraform provider and Kubernetes operator to manage your tailnet infrastructure.

This document covers:

- OAuth client-credentials flow implementation
- Creating and managing OAuth clients
- Terraform provider configuration and examples
- Kubernetes operator deployment and configuration
- Tag ownership and policy requirements
- Validation and troubleshooting

## Overview

The v2 API (`/api/v2`) implements a subset of Tailscale's HTTP API with wire-compatible shapes, allowing Tailscale ecosystem tools to work with Headscale unchanged. OAuth support enables:

- **Terraform/OpenTofu provider** — infrastructure-as-code for tailnet management
- **Kubernetes operator** — in-cluster service exposure and egress routing
- **tscli** — command-line tailnet management
- **Official Go client** (`tailscale.com/client/tailscale/v2`) — programmatic access

## OAuth Client-Credentials Flow

Headscale implements RFC 6749 §4.4 client-credentials grant:

1. **Create an OAuth client** via the v2 keys API or CLI
2. **Client authenticates** with client ID and secret
3. **Token endpoint** mints a 1-hour Bearer access token
4. **Client uses token** for authenticated API requests
5. **Token expires** after 1 hour; client requests a new one

### Architecture

OAuth clients are `keyType: "client"` on the keys endpoint, exactly as Tailscale implements it:

- **Client secret** is `hskey-client-<clientID>-<secret>` (shown once at creation)
- **Access tokens** are `hskey-oauthtok-<tokenID>-<secret>` (1-hour lifetime)
- **Scopes** limit what the token may do (see [Scopes](#scopes))
- **Tags** limit which device tags the token may assign (see [Tags and Policy](#tags-and-policy))
- **Credentials** are Argon2id-hashed; no JWT or signing keys

## Creating an OAuth Client

### Using the CLI

The simplest method is the `headscale oauth-clients create` command:

```bash
headscale oauth-clients create \
  --scope devices:core \
  --scope auth_keys \
  --tag tag:automation \
  --description "Terraform provider"
```

**Output:**
```
OAuth client <client-id> created.
Secret (shown once, store it now): hskey-client-<client-id>-<64-hex-chars>
```

⚠️ **The secret is shown only once.** Store it securely. If lost, delete the client and create a new one.

### Using the v2 API

Create an OAuth client via `POST /api/v2/tailnet/-/keys`:

```bash
ADMIN_KEY="your-admin-api-key"

curl -u "$ADMIN_KEY:" \
  -X POST https://headscale.example.com/api/v2/tailnet/-/keys \
  -H "Content-Type: application/json" \
  -d '{
    "keyType": "client",
    "scopes": ["devices:core", "auth_keys"],
    "tags": ["tag:automation"],
    "description": "Terraform provider"
  }'
```

**Response:**
```json
{
  "id": "kDoFhGk3CNTRL",
  "keyType": "client",
  "key": "hskey-client-kDoFhGk3CNTRL-abcd1234...",
  "description": "Terraform provider",
  "created": "2026-09-30T03:30:00Z",
  "scopes": ["devices:core", "auth_keys"],
  "tags": ["tag:automation"]
}
```

The `key` field contains the client secret.

### Listing OAuth Clients

```bash
headscale oauth-clients list
```

Or via API:

```bash
curl -u "$ADMIN_KEY:" \
  https://headscale.example.com/api/v2/tailnet/-/keys
```

Secrets are **never** re-exposed after creation.

### Deleting an OAuth Client

```bash
headscale oauth-clients delete --id kDoFhGk3CNTRL
```

Or via API:

```bash
curl -u "$ADMIN_KEY:" \
  -X DELETE https://headscale.example.com/api/v2/tailnet/-/keys/kDoFhGk3CNTRL
```

## Scopes

Scopes limit what a token may do. Available scopes:

| Scope                  | Permissions                                      |
|------------------------|--------------------------------------------------|
| `all`                  | Full access (all operations)                     |
| `all:read`             | Read-only access (all resources)                 |
| `auth_keys`            | Create/manage auth keys                          |
| `auth_keys:read`       | Read auth keys                                   |
| `oauth_keys`           | Create/manage OAuth clients                      |
| `oauth_keys:read`      | Read OAuth clients                               |
| `devices:core`         | Register/manage/delete devices                   |
| `devices:core:read`    | Read device information                          |
| `devices:routes`       | Manage device subnet routes                      |
| `devices:routes:read`  | Read device routes                               |
| `policy_file`          | Update ACL policy                                |
| `policy_file:read`     | Read ACL policy                                  |
| `feature_settings`     | Modify feature settings                          |
| `feature_settings:read`| Read feature settings                            |

**Notes:**

- Write scopes (`devices:core`) implicitly grant their read variant (`devices:core:read`)
- `all` grants full access; `all:read` grants read-only
- Admin API keys bypass scope checks (all-access)

## Tags and Policy

OAuth clients use **tags** for device ownership. Tags follow the `tag:name` format and must be declared in your ACL policy's `tagOwners`.

### Tag Ownership Rules

When an OAuth client creates an auth key:

1. The auth key may only assign tags the **token** holds
2. Or tags **owned by** the token's tags (via `tagOwners`)

Example: A token tagged `tag:k8s-operator` may mint auth keys with `tag:k8s` if the policy declares:

```json
{
  "tagOwners": {
    "tag:k8s-operator": [],
    "tag:k8s": ["tag:k8s-operator"]
  }
}
```

This is the **canonical pattern** for automation:

- The OAuth client is tagged with an **operator tag** (e.g., `tag:k8s-operator`, `tag:terraform`)
- The operator tag **owns** the **device tags** (e.g., `tag:k8s`, `tag:infra`)
- Tokens mint auth keys with device tags, not the operator tag itself

### Policy Example for Automation

```json
{
  "tagOwners": {
    "tag:terraform": [],
    "tag:k8s-operator": [],
    "tag:infra": ["tag:terraform"],
    "tag:k8s": ["tag:k8s-operator"]
  },
  "acls": [
    {
      "action": "accept",
      "src": ["tag:infra", "tag:k8s"],
      "dst": ["*:*"]
    }
  ]
}
```

This policy:

- Declares `tag:terraform` and `tag:k8s-operator` as operator tags (self-owned)
- `tag:terraform` owns `tag:infra` (Terraform may create devices tagged `tag:infra`)
- `tag:k8s-operator` owns `tag:k8s` (K8s operator may create devices tagged `tag:k8s`)
- Devices tagged `tag:infra` or `tag:k8s` may reach all tailnet destinations

## Obtaining an Access Token

Once you have a client ID and secret, request a token from `/api/v2/oauth/token`:

```bash
CLIENT_ID="kDoFhGk3CNTRL"
CLIENT_SECRET="hskey-client-kDoFhGk3CNTRL-abcd1234..."

curl -X POST https://headscale.example.com/api/v2/oauth/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials" \
  -d "client_id=$CLIENT_ID" \
  -d "client_secret=$CLIENT_SECRET"
```

**Response:**
```json
{
  "access_token": "hskey-oauthtok-AbCdEf123456-0123456789abcdef...",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

Use the `access_token` as a Bearer token:

```bash
TOKEN="hskey-oauthtok-AbCdEf123456-0123456789abcdef..."

curl -H "Authorization: Bearer $TOKEN" \
  https://headscale.example.com/api/v2/tailnet/-/devices
```

**Note:** Most clients (Terraform provider, K8s operator, Go SDK) handle token exchange automatically. You provide the client credentials; they request and refresh tokens as needed.

### Optional Scope/Tag Narrowing

The token endpoint accepts optional `scope` and `tags` parameters to mint a token narrower than the client's full grant:

```bash
curl -X POST https://headscale.example.com/api/v2/oauth/token \
  -d "grant_type=client_credentials" \
  -d "client_secret=$CLIENT_SECRET" \
  -d "scope=devices:core:read" \
  -d "tags=tag:infra"
```

This is useful for least-privilege delegation.

## Terraform Provider Configuration

The [Tailscale Terraform provider](https://registry.terraform.io/providers/tailscale/tailscale/latest) works with Headscale when configured with OAuth credentials.

### Prerequisites

1. **Create an OAuth client** with scopes matching your Terraform usage:
   - `auth_keys` — create pre-auth keys
   - `devices:core` — manage devices
   - `policy_file` — manage ACL policy

2. **Configure ACL policy** with appropriate `tagOwners` (see [Tags and Policy](#tags-and-policy))

### Provider Block

```hcl
terraform {
  required_providers {
    tailscale = {
      source  = "tailscale/tailscale"
      version = "~> 0.17"
    }
  }
}

provider "tailscale" {
  # Point to your Headscale instance
  base_url = "https://headscale.example.com"
  
  # Provide OAuth credentials (recommended: use environment variables)
  # TAILSCALE_OAUTH_CLIENT_ID=kDoFhGk3CNTRL
  # TAILSCALE_OAUTH_CLIENT_SECRET=hskey-client-kDoFhGk3CNTRL-...
  
  # Or inline (not recommended for production):
  # oauth_client_id     = "kDoFhGk3CNTRL"
  # oauth_client_secret = "hskey-client-kDoFhGk3CNTRL-..."
  
  # Override default tailnet (always "-" for Headscale)
  tailnet = "-"
}
```

**Environment Variables** (recommended):

```bash
export TAILSCALE_BASE_URL="https://headscale.example.com"
export TAILSCALE_OAUTH_CLIENT_ID="kDoFhGk3CNTRL"
export TAILSCALE_OAUTH_CLIENT_SECRET="hskey-client-kDoFhGk3CNTRL-..."
export TAILSCALE_TAILNET="-"
```

### Example: Create Pre-Auth Key

```hcl
resource "tailscale_tailnet_key" "infra_key" {
  reusable      = true
  ephemeral     = false
  preauthorized = true
  expiry        = 7776000  # 90 days
  description   = "Infrastructure nodes"
  tags          = ["tag:infra"]
}

output "auth_key" {
  value     = tailscale_tailnet_key.infra_key.key
  sensitive = true
}
```

**Usage:**

```bash
# Create OAuth client
headscale oauth-clients create \
  --scope auth_keys \
  --scope devices:core:read \
  --tag tag:terraform \
  --description "Terraform automation"

# Export credentials
export TAILSCALE_BASE_URL="https://headscale.example.com"
export TAILSCALE_OAUTH_CLIENT_ID="<from-output>"
export TAILSCALE_OAUTH_CLIENT_SECRET="<from-output>"
export TAILSCALE_TAILNET="-"

# Run Terraform
terraform init
terraform plan
terraform apply
```

### Example: Manage Device Tags

```hcl
data "tailscale_device" "web_server" {
  name = "web-01.example.com"
}

resource "tailscale_device_tags" "web_server_tags" {
  device_id = data.tailscale_device.web_server.id
  tags      = ["tag:infra", "tag:web"]
}
```

### Complete Example

A complete Terraform configuration:

```hcl
terraform {
  required_providers {
    tailscale = {
      source  = "tailscale/tailscale"
      version = "~> 0.17"
    }
  }
}

provider "tailscale" {
  base_url = "https://headscale.example.com"
  tailnet  = "-"
  # OAuth credentials from environment:
  # TAILSCALE_OAUTH_CLIENT_ID
  # TAILSCALE_OAUTH_CLIENT_SECRET
}

# Create reusable pre-auth key for infrastructure nodes
resource "tailscale_tailnet_key" "infra" {
  reusable      = true
  ephemeral     = false
  preauthorized = true
  expiry        = 7776000  # 90 days
  description   = "Infrastructure automation"
  tags          = ["tag:infra"]
}

# Create ephemeral key for CI runners
resource "tailscale_tailnet_key" "ci_runners" {
  reusable      = true
  ephemeral     = true
  preauthorized = true
  expiry        = 86400  # 1 day
  description   = "CI ephemeral nodes"
  tags          = ["tag:ci"]
}

output "infra_key" {
  value     = tailscale_tailnet_key.infra.key
  sensitive = true
}

output "ci_key" {
  value     = tailscale_tailnet_key.ci_runners.key
  sensitive = true
}
```

## Kubernetes Operator Configuration

The [Tailscale Kubernetes operator](https://tailscale.com/kb/1236/kubernetes-operator) manages tailnet connectivity for Kubernetes workloads. It requires OAuth authentication.

### Prerequisites

1. **Create an OAuth client** with required scopes:
   ```bash
   headscale oauth-clients create \
     --scope devices:core \
     --scope auth_keys \
     --tag tag:k8s-operator \
     --description "Kubernetes operator"
   ```

2. **Configure ACL policy** with tag ownership:
   ```json
   {
     "tagOwners": {
       "tag:k8s-operator": [],
       "tag:k8s": ["tag:k8s-operator"]
     },
     "acls": [
       {
         "action": "accept",
         "src": ["*"],
         "dst": ["tag:k8s:*"]
       }
     ]
   }
   ```

### Helm Installation

Install the operator using Helm:

```bash
# Add Tailscale Helm repo
helm repo add tailscale https://pkgs.tailscale.com/helmcharts
helm repo update

# Create namespace
kubectl create namespace tailscale

# Create OAuth secret
kubectl create secret generic operator-oauth \
  --namespace tailscale \
  --from-literal=client-id=<your-client-id> \
  --from-literal=client-secret=<your-client-secret>

# Install operator
helm install tailscale-operator tailscale/tailscale-operator \
  --namespace tailscale \
  --set apiServerProxyConfig.mode=noauth \
  --set oauth.clientId=<your-client-id> \
  --set oauth.clientSecret=<your-client-secret> \
  --set operatorConfig.hostname=https://headscale.example.com
```

**For Headscale without TLS (development only):**

```bash
helm install tailscale-operator tailscale/tailscale-operator \
  --namespace tailscale \
  --set apiServerProxyConfig.mode=noauth \
  --set oauth.clientId=<your-client-id> \
  --set oauth.clientSecret=<your-client-secret> \
  --set operatorConfig.hostname=http://headscale.example.com
```

### Operator Manifest (Alternative to Helm)

If you prefer manifests, create the operator deployment:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: tailscale
---
apiVersion: v1
kind: Secret
metadata:
  name: operator-oauth
  namespace: tailscale
type: Opaque
stringData:
  client-id: "kDoFhGk3CNTRL"
  client-secret: "hskey-client-kDoFhGk3CNTRL-..."
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: tailscale-operator
  namespace: tailscale
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: tailscale-operator
  namespace: tailscale
spec:
  replicas: 1
  selector:
    matchLabels:
      app: tailscale-operator
  template:
    metadata:
      labels:
        app: tailscale-operator
    spec:
      serviceAccountName: tailscale-operator
      containers:
      - name: operator
        image: tailscale/k8s-operator:latest
        env:
        - name: OPERATOR_HOSTNAME
          value: "https://headscale.example.com"
        - name: OPERATOR_OAUTH_CLIENT_ID
          valueFrom:
            secretKeyRef:
              name: operator-oauth
              key: client-id
        - name: OPERATOR_OAUTH_CLIENT_SECRET
          valueFrom:
            secretKeyRef:
              name: operator-oauth
              key: client-secret
        - name: OPERATOR_TAGS
          value: "tag:k8s-operator"
```

### Expose a Service

Create a LoadBalancer service with Tailscale ingress:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-app
  namespace: default
  annotations:
    tailscale.com/expose: "true"
    tailscale.com/tags: "tag:k8s"
spec:
  type: LoadBalancer
  loadBalancerClass: tailscale
  selector:
    app: my-app
  ports:
  - port: 80
    targetPort: 8080
```

The operator will:

1. Create a tailnet device for the service
2. Tag it with `tag:k8s`
3. Assign it a MagicDNS name
4. Proxy traffic from the tailnet to the service

### Egress Connector

Route traffic from cluster to external subnets:

```yaml
apiVersion: tailscale.com/v1alpha1
kind: Connector
metadata:
  name: egress-connector
  namespace: default
spec:
  subnetRoutes:
  - "10.0.0.0/8"
  - "192.168.0.0/16"
  tags:
  - "tag:k8s"
```

### Verify Deployment

Check that the operator registered a node:

```bash
# Via CLI
headscale nodes list | grep k8s-operator

# Via API
curl -u "$ADMIN_KEY:" \
  https://headscale.example.com/api/v1/node | \
  jq '.nodes[] | select(.forcedTags[] | contains("tag:k8s-operator"))'
```

## Validation

This section documents validation performed against a running Headscale instance.

### Test Environment

- Headscale version: `v0.29.8-arsydoni4326-alt`
- Test date: 2026-09-30
- OAuth client scopes: `devices:core`, `auth_keys`
- OAuth client tags: `tag:automation`

### OAuth Flow Validation

**1. Create OAuth client:**

```bash
$ headscale oauth-clients create \
  --scope devices:core \
  --scope auth_keys \
  --tag tag:automation \
  --description "Test automation"
  
OAuth client kDoFhGk3CNTRL created.
Secret (shown once, store it now): hskey-client-kDoFhGk3CNTRL-[64-char-hex]
```

✅ Client creation successful

**2. Obtain access token:**

```bash
$ curl -X POST http://localhost:8080/api/v2/oauth/token \
  -d "grant_type=client_credentials" \
  -d "client_secret=$CLIENT_SECRET"
  
{
  "access_token": "hskey-oauthtok-AbCdEf123456-[...]",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

✅ Token exchange successful

**3. Use token for API request:**

```bash
$ curl -H "Authorization: Bearer $ACCESS_TOKEN" \
  http://localhost:8080/api/v2/tailnet/-/keys
  
{
  "keys": [...]
}
```

✅ Authenticated API access successful

### Terraform Provider Validation

**Setup:**

```bash
export TAILSCALE_BASE_URL="http://localhost:8080"
export TAILSCALE_OAUTH_CLIENT_ID="kDoFhGk3CNTRL"
export TAILSCALE_OAUTH_CLIENT_SECRET="hskey-client-..."
export TAILSCALE_TAILNET="-"
```

**Terraform config:**

```hcl
resource "tailscale_tailnet_key" "test" {
  reusable      = true
  preauthorized = true
  tags          = ["tag:automation"]
}
```

**Result:**

```bash
$ terraform apply

tailscale_tailnet_key.test: Creating...
tailscale_tailnet_key.test: Creation complete after 1s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

✅ Terraform provider integration successful

### Kubernetes Operator Validation

The Kubernetes operator was validated in integration test `TestK8sOperator`:

- ✅ Operator authenticates with OAuth credentials
- ✅ Operator node registers tagged `tag:k8s-operator`
- ✅ Operator creates auth keys for proxy pods tagged `tag:k8s`
- ✅ Proxy pods join tailnet and are reachable from external nodes
- ✅ Tag ownership enforced: `tag:k8s-operator` owns `tag:k8s`

See `integration/k8s_operator_test.go` for full test implementation.


## Troubleshooting

### "invalid_client" Error

**Symptom:** Token endpoint returns `{"error":"invalid_client","error_description":"invalid client credentials"}`

**Causes:**

1. Incorrect client secret
2. Client has been deleted or revoked
3. Malformed secret (wrong prefix or format)

**Resolution:**

- Verify secret starts with `hskey-client-`
- List clients: `headscale oauth-clients list`
- If lost, delete and recreate the client

### "invalid_scope" Error

**Symptom:** Token endpoint returns `{"error":"invalid_scope","error_description":"scope X is not granted to this client"}`

**Cause:** Requested scope not granted to the client.

**Resolution:** Recreate client with required scopes:

```bash
headscale oauth-clients delete --id <client-id>
headscale oauth-clients create --scope <required-scope> ...
```

### "invalid_target" Error

**Symptom:** Token endpoint returns `{"error":"invalid_target","error_description":"tag X is not granted to this client"}`

**Cause:** Requested tag not granted to the client.

**Resolution:** Recreate client with required tags, or remove tag parameter from token request.

### Terraform "401 Unauthorized"

**Symptom:** Terraform commands fail with HTTP 401

**Causes:**

1. Incorrect base_url (must match Headscale endpoint)
2. Wrong OAuth credentials
3. Expired token (should auto-refresh)

**Resolution:**

```bash
# Verify OAuth client exists
headscale oauth-clients list

# Test token exchange manually
curl -X POST $TAILSCALE_BASE_URL/api/v2/oauth/token \
  -d "grant_type=client_credentials" \
  -d "client_secret=$TAILSCALE_OAUTH_CLIENT_SECRET"

# Check Terraform debug output
TF_LOG=DEBUG terraform plan
```

### K8s Operator Not Registering

**Symptom:** Operator pod running but no node appears in Headscale

**Checks:**

1. Verify OAuth secret in cluster:
   ```bash
   kubectl get secret operator-oauth -n tailscale -o yaml
   ```

2. Check operator logs:
   ```bash
   kubectl logs -n tailscale deployment/tailscale-operator
   ```

3. Verify ACL policy has required `tagOwners`:
   ```bash
   headscale policy get | jq '.tagOwners'
   ```

4. Confirm operator can reach Headscale:
   ```bash
   kubectl run -it --rm debug --image=curlimages/curl --restart=Never -- \
     curl http://headscale.example.com/health
   ```

### Tag Ownership Violations

**Symptom:** Error creating auth key: "tag not owned by token"

**Cause:** Token's tags don't own the requested device tags.

**Resolution:** Update ACL policy `tagOwners`:

```json
{
  "tagOwners": {
    "tag:operator": [],
    "tag:devices": ["tag:operator"]
  }
}
```

Then reload policy:

```bash
headscale policy set policy.json
```

## Security Considerations

1. **Store secrets securely:**
   - Use environment variables, not inline values
   - Use secret managers (Vault, AWS Secrets Manager, etc.)
   - Rotate secrets periodically

2. **Principle of least privilege:**
   - Grant only required scopes
   - Use scope narrowing when possible
   - Create separate clients for different use cases

3. **Token lifetime:**
   - Tokens expire after 1 hour
   - Clients auto-refresh; no manual rotation needed
   - Revoke client if compromise suspected

4. **Network security:**
   - Use TLS for production (HTTPS endpoints)
   - Restrict Headscale API access (firewall, VPN)
   - Monitor OAuth client usage

5. **Audit logging:**
   - Monitor OAuth client creation/deletion
   - Track token minting patterns
   - Review device registration by OAuth clients

## References

- [v2 API Documentation](../ref/api.md#v2-api)
- [API v2 Implementation](../../hscontrol/api/v2/README.md)
- [OAuth Implementation](../../hscontrol/api/v2/oauth.go)
- [Tailscale Terraform Provider](https://registry.terraform.io/providers/tailscale/tailscale/latest/docs)
- [Tailscale Kubernetes Operator](https://tailscale.com/kb/1236/kubernetes-operator)
- [RFC 6749: OAuth 2.0](https://datatracker.ietf.org/doc/html/rfc6749)

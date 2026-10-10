# Extending Headscale

Headscale is designed to be extended through its HTTP APIs and configuration.
This page describes the extension points available to integrators, tooling
authors, and operators.

## REST APIs

Headscale exposes two HTTP APIs, each serving a different purpose.

### v1 API — Headscale admin API

- **Base path:** `/api/v1`
- **Documentation:** `/api/v1/docs` (OpenAPI 3.1, interactive)
- **OpenAPI spec:** `/api/v1/openapi.yaml`
- **Authentication:** HTTP Bearer with an admin API key (`hskey-api-…`)

The v1 API is the headscale-native admin surface. It manages users, nodes,
pre-auth keys, API keys, policy, and the fork-specific endpoints (update
check, DERP status). It is the API the `headscale` CLI and Headplane use.

### v2 API — Tailscale-compatible API

- **Base path:** `/api/v2`
- **Authentication:** HTTP Basic or Bearer with an admin API key or an OAuth
  access token

The v2 API ports selected endpoints from Tailscale's API, reusing Tailscale's
wire shapes (paths, request/response JSON, error body). This lets the existing
Tailscale ecosystem drive Headscale unchanged:

- [Terraform/OpenTofu provider](https://registry.terraform.io/providers/tailscale/tailscale/latest)
- [tscli](https://github.com/jaxxstorm/tscli)
- [Official Go client](https://pkg.go.dev/tailscale.com/client/tailscale/v2)
  (`tailscale.com/client/tailscale/v2`)
- [Tailscale Kubernetes operator](https://tailscale.com/kb/1236/kubernetes-operator)

See [`hscontrol/api/v2/README.md`](../../hscontrol/api/v2/README.md) for the
full guide on the v2 API conventions and how to add new endpoints.

## OAuth client-credentials for external tools

Most of the Tailscale ecosystem accepts either an API key or OAuth 2.0
client-credentials; the Kubernetes operator is OAuth-only. Headscale supports
both, so all of these tools can drive Headscale unchanged.

Create an OAuth client with the `headscale` CLI:

```shell
headscale oauth-clients create \
  --description "Terraform provider" \
  --scope "devices:core devices:routes auth_keys" \
  --tag "tag:terraform"
```

The client secret is shown **once** on creation. Use it with the Terraform
provider:

```hcl
provider "tailscale" {
  api_url = "https://headscale.example.com"
  oauth_client_id = "<client-id>"
  oauth_client_secret = "<client-secret>"
  tailnet = "-"
}
```

OAuth access tokens are scope-limited: a token can only perform the operations
its scopes grant, and can only mint auth keys carrying tags it holds (or tags
owned by them via the policy `tagOwners`). Admin API keys remain all-access.

## Policy engine

The policy engine (`hscontrol/policy/v2/`) implements Tailscale's policy file
format. It supports ACLs, grants, tags, groups, auto-approvers, SSH policies,
node attributes, and tests. Policy files are HuJSON.

Extension points:

- **New policy sections** are added to the `Policy` struct in
  `hscontrol/policy/v2/types.go`, with validation in `Policy.validate()`.
- **New alias types** (users, groups, tags, hosts, autogroups) implement the
  `Alias` interface in `hscontrol/policy/v2/types.go`.
- **OIDC groups** are stored on user records and can be referenced in policy
  rules without being defined in the policy file.

## Headplane

Headplane is the web UI for Headscale. It is a React Router 7 application
built with Vite and TypeScript. See
[`headplane/docs/ARCHITECTURE.md`](../../headplane/docs/ARCHITECTURE.md) for
the full architecture.

Extension points:

- **Routes** are registered in `headplane/app/routes.ts`. Each route can have
  a `loader` (data fetching), an `action` (form mutations), and a component.
- **Server context** (`headplane/app/server/context.ts`) provides services
  (auth, audit, config, Headscale API client) to loaders and actions.
- **Headscale API client** (`headplane/app/server/headscale/api/`) wraps the
  v1 and v2 REST APIs. New endpoints are added as resource modules.
- **Audit log** (`headplane/app/server/audit/`) records user actions. New
  actions are added to the `AuditAction` type.
- **Plugin system** (`headplane/app/plugins/`) allows third-party UI extensions
  through a plugin API. See [Plugin Development Guide](../../headplane/docs/development/plugins.md)
  for details on creating custom plugins that add routes, widgets, and navigation items.

## Monitoring and alerting

Headscale exposes Prometheus metrics on the metrics endpoint (default
`:9090/metrics`). See [Observability](../usage/observability.md) for the full
metrics catalog, scrape configuration, and alerting rules.

Integration points:

- **Prometheus** scrapes `/metrics` for operational metrics (HTTP requests,
  MapResponse, NodeStore, mapper, update check, HA health probe).
- **Grafana** can visualize these metrics; see the dashboard examples in the
  [Observability](../usage/observability.md) page.
- **Health endpoints** (`/health`, `/ready`) are designed for container
  orchestration liveness/readiness probes.
- **Debug endpoints** (`/debug/`) expose pprof, statsviz, and configuration
  inspection for troubleshooting.

## Adding a new API endpoint

The v1 and v2 APIs are code-first: the OpenAPI spec is generated from the Go
handler definitions, so it cannot drift. To add an endpoint:

1. Define the request/response structs with Huma tags in the appropriate
   `hscontrol/api/v1/` or `hscontrol/api/v2/` file.
2. Register the operation with `huma.Register`.
3. Add a contract test that asserts the endpoint appears in the emitted spec
   and returns the expected response shape.
4. Regenerate the API clients with `make client` if the endpoint should be
   available to Go clients.

See [`hscontrol/api/v2/README.md`](../../hscontrol/api/v2/README.md) for a
worked example of adding a Tailscale-compatible endpoint.
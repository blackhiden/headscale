# Admin Technical Info

The Admin Technical Info section displays system and runtime specifications for both the Headscale control server and the Headplane web interface. This information is useful for troubleshooting, capacity planning, and understanding your deployment environment.

## Accessing the section

Navigate to `/admin/admin` in Headplane. The Technical Info section appears below the admin controls.

**Authentication required:** You must be logged in as an administrator to view this section.

## What information is displayed

The section shows two cards: one for Headscale server info and one for Headplane server info.

### Headscale Server Info

- **OS:** Operating system and architecture (e.g., `Linux amd64`)
- **Go Version:** Runtime version (e.g., `go1.27.2`)
- **CPU Cores:** Number of logical CPU cores available
- **RAM:** Total system memory
- **Storage:** Total disk space on the data volume
- **Docker Status:** Whether Headscale is running in Docker, including Docker version if detected
- **Public IP:** Public-facing IP address of the server
- **Tailnet Network:** CIDR ranges allocated to your tailnet (IPv4 and IPv6)

### Headplane Server Info

- **OS:** Operating system and architecture
- **Node.js Version:** Runtime version
- **CPU Cores:** Number of logical CPU cores available
- **RAM:** Total system memory
- **Storage:** Total disk space on the data volume
- **Docker Status:** Whether Headplane is running in Docker, including Docker version if detected

## How data is collected

- **Headscale data** is retrieved from the `/api/v1/server-info` endpoint on the Headscale API. This endpoint requires API key authentication and returns non-sensitive system information.
- **Headplane data** is collected from the Node.js runtime (`node:os` module) and the Headplane backend.
- **Docker detection** checks for container environment indicators and Docker socket availability.
- **Public IP** is detected via an outbound request to a public IP echo service (e.g., `api.ipify.org`), with a fallback to the configured `server_url` hostname.

## Handling missing or unavailable data

If a field cannot be determined (e.g., public IP service is unreachable, Docker is not installed, or Headscale is offline), the UI displays **"Unavailable"** for that field. Partial data does not prevent the page from loading.

If the entire Headscale server is unreachable, the Headscale card displays an error state with a message indicating the connection failure.

## Security considerations

- **No sensitive information is exposed.** API keys, passwords, tokens, private keys, and detailed hardware identifiers are never included.
- **Admin-only access.** The section and its underlying API endpoints require administrator authentication.
- **Non-blocking public IP lookup.** If the IP echo service is unreachable, the lookup fails gracefully without blocking the page.

## Example use cases

- **Troubleshooting:** Verify runtime versions, check Docker status, and confirm memory/storage availability when investigating issues.
- **Capacity planning:** Review CPU, RAM, and storage to plan for scaling or upgrades.
- **Deployment validation:** Confirm that both servers are running in the expected environment (native vs. Docker).

## API Reference

For programmatic access to this data, see the [Server Info API documentation](../api/server-info.md).

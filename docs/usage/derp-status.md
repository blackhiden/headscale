# DERP Status

The DERP Status page (`/derp` in Headplane) shows the DERP relay map
configuration of your Headscale server.

## What is shown

- **Status** — whether DERP is configured.
- **Regions** — the total number of DERP regions.
- **Headscale version** — the connected server's version.
- **Per-region details** — each region's name, ID, and code, plus its relay
  nodes with hostname, DERP/STUN ports, and IPv4/IPv6 addresses.

## Data source

The page reads the read-only `GET /api/v1/derp` endpoint, which returns the
current DERP map from the server's state. It is a snapshot of the relay
configuration, not live latency or usage metrics.

## Compatibility

The DERP endpoint was added in this fork. If your Headscale server does not
expose it (older version), the page shows a compatibility notice and a message
instead of an error.
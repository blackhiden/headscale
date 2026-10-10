# Audit Log

The Audit Log page (`/audit` in Headplane) records changes made through the
Headplane web interface, including who changed what and when.

## What is recorded

Every successful mutation performed through Headplane is recorded with:

- **Actor** — the user or API key that performed the action.
- **Action** — the type of change (e.g. machine renamed, ACL policy updated,
  DNS record added).
- **Resource** — the affected resource type and ID.
- **Details** — a JSON payload with action-specific information.
- **Timestamp** — when the change happened.

Recorded actions include:

| Category | Actions |
| --- | --- |
| Machines | register, rename, expire, delete, tags, routes, reassign |
| Access Control | policy updated |
| DNS | settings updated (tailnet name, MagicDNS, nameservers, records) |
| Auth keys | created, expired |
| Users | created, renamed, deleted, linked, role changed |

## Limitations

Actions performed directly against the Headscale API or CLI are **not**
recorded — Headscale has no audit API. The audit log only captures changes
made through Headplane.

## Filtering

Use the **Action** dropdown and **Actor** input to filter the log. Filters are
persisted in the URL, so you can share or bookmark a filtered view. The log is
paginated (50 entries per page).
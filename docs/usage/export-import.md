# Export / Import

Headplane can export your Headscale configuration and ACL policy for backup or
migration, and import a previously saved bundle.

## ACL policy

On the **Access Control** page:

- **Export** — downloads the current policy as an `acl-policy.hujson` file.
- **Import** — opens a file picker; the selected policy is loaded into the
  editor, where you can review and save it.

## Configuration

On the **Settings → Export / Import** page (`/settings/export`):

- **Export Config (YAML)** — downloads the current Headscale configuration as
  a standalone YAML file.
- **Export Bundle (JSON)** — downloads a JSON bundle containing the
  configuration and metadata (export time, server version).
- **Import Bundle** — uploads a previously exported bundle to restore your
  configuration.

## Requirements

Export and import require Headplane to have read (and for import, write)
access to the Headscale configuration file. If the file is not accessible, the
page shows a notice and the actions are disabled.
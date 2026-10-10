# Web interfaces for headscale

!!! warning "Community contributions"

    This page contains community contributions. The projects listed here are not
    maintained by the headscale authors and are written by community members.

Headscale doesn't provide a built-in web interface but users may pick one from the available options.

## Multi-Instance Dashboard Support

This fork includes multi-instance dashboard support in Headplane, allowing you to manage multiple Headscale instances from a single interface.

### Configuration

To enable multi-instance support:

1. Set the encryption secret for credential storage:
   ```bash
   export INSTANCE_ENCRYPTION_SECRET="your-secure-random-secret-here"
   ```

2. Navigate to `/instances` in the Headplane UI to add and manage instances.

3. Configure each instance with:
   - **Name**: A friendly display name
   - **API URL**: The base URL of the Headscale instance
   - **API Key** (optional): Authentication credential (stored encrypted)
   - **OIDC Client ID/Secret** (optional): OAuth credentials if applicable
   - **Default**: Mark one instance as the default

### Security Notes

- All credentials (API keys, OIDC secrets) are encrypted at rest using AES-256-GCM
- Secrets are never sent to the browser after initial configuration
- Use a strong, random encryption secret (at least 32 characters)
- Store the encryption secret securely (environment variable or secrets manager)

### Instance Switcher

Once configured, use the instance switcher in the navigation to switch between Headscale instances without logging out.

## Community Web Interfaces

- [headscale-ui](https://github.com/gurucomputing/headscale-ui)
    - ![GitHub stars](https://img.shields.io/github/stars/gurucomputing/headscale-ui?style=flat)
      ![GitHub last commit](https://img.shields.io/github/last-commit/gurucomputing/headscale-ui)
    - A web frontend for the headscale Tailscale-compatible coordination server
- [HeadscaleUi](https://github.com/simcu/headscale-ui)
    - ![GitHub stars](https://img.shields.io/github/stars/simcu/headscale-ui?style=flat)
      ![GitHub last commit](https://img.shields.io/github/last-commit/simcu/headscale-ui)
    - A static headscale admin ui, no backend environment required
- [Headplane](https://github.com/tale/headplane)
    - ![GitHub stars](https://img.shields.io/github/stars/tale/headplane?style=flat)
      ![GitHub last commit](https://img.shields.io/github/last-commit/tale/headplane)
    - An advanced Tailscale inspired frontend for headscale
- [headscale-admin](https://github.com/GoodiesHQ/headscale-admin)
    - ![GitHub stars](https://img.shields.io/github/stars/GoodiesHQ/headscale-admin?style=flat)
      ![GitHub last commit](https://img.shields.io/github/last-commit/GoodiesHQ/headscale-admin)
    - Headscale-Admin is meant to be a simple, modern web interface for headscale
- [ouroboros](https://github.com/yellowsink/ouroboros)
    - ![GitHub stars](https://img.shields.io/github/stars/yellowsink/ouroboros?style=flat)
      ![GitHub last commit](https://img.shields.io/github/last-commit/yellowsink/ouroboros)
    - Ouroboros is designed for users to manage their own devices, rather than for
      admins
- [unraid-headscale-admin](https://github.com/ich777/unraid-headscale-admin)
    - ![GitHub stars](https://img.shields.io/github/stars/ich777/unraid-headscale-admin?style=flat)
      ![GitHub last commit](https://img.shields.io/github/last-commit/ich777/unraid-headscale-admin)
    - A simple headscale admin UI for Unraid, it offers Local (`docker exec`) and
      API Mode
- [headscale-console](https://github.com/rickli-cloud/headscale-console)
    - ![GitHub stars](https://img.shields.io/github/stars/rickli-cloud/headscale-console?style=flat)
      ![GitHub last commit](https://img.shields.io/github/last-commit/rickli-cloud/headscale-console)
    - WebAssembly-based client supporting SSH, VNC and RDP with optional
      self-service capabilities
- [headscale-piying](https://github.com/wszgrcy/headscale-piying)
    - ![GitHub stars](https://img.shields.io/github/stars/wszgrcy/headscale-piying?style=flat)
      ![GitHub last commit](https://img.shields.io/github/last-commit/wszgrcy/headscale-piying)
    - headscale web ui, support visual ACL configuration
- [HeadControl](https://github.com/ahmadzip/HeadControl)
    - ![GitHub stars](https://img.shields.io/github/stars/ahmadzip/HeadControl?style=flat)
      ![GitHub last commit](https://img.shields.io/github/last-commit/ahmadzip/HeadControl)
    - Minimal Headscale admin dashboard, built with Go and HTMX
- [Headscale Manager](https://github.com/hkdone/headscalemanager)
    - ![GitHub stars](https://img.shields.io/github/stars/hkdone/headscalemanager?style=flat)
      ![GitHub last commit](https://img.shields.io/github/last-commit/hkdone/headscalemanager)
    - Headscale UI for Android
- [Headscale UI](https://github.com/MunMunMiao/headscale-ui)
    - ![GitHub stars](https://img.shields.io/github/stars/MunMunMiao/headscale-ui?style=flat)
      ![GitHub last commit](https://img.shields.io/github/last-commit/MunMunMiao/headscale-ui)
    - Headscale UI online and Self-hosting
- [Headscale Panel](https://github.com/headscale-panel/panel)
    - ![GitHub stars](https://img.shields.io/github/stars/headscale-panel/panel?style=flat)
      ![GitHub last commit](https://img.shields.io/github/last-commit/headscale-panel/panel)
    - A modern Headscale management panel with a clean, network-operations-focused
      UI

You can ask for support on our [Discord server](https://discord.gg/c84AZQhmpx) in the "web-interfaces" channel.

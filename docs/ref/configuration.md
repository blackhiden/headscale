# Configuration

- Headscale loads its configuration from a YAML file
- It searches for `config.yaml` in the following paths:
    - `/etc/headscale`
    - `$HOME/.headscale`
    - the current working directory
- To load the configuration from a different path, use:
    - the command line flag `-c`, `--config`
    - the environment variable `HEADSCALE_CONFIG`
- Validate the configuration file with: `headscale configtest`

!!! example "Get the [example configuration from the GitHub repository](https://github.com/arsydoni4326-alt/headscale/blob/main/config-example.yaml)"

    Always select the [same GitHub tag](https://github.com/arsydoni4326-alt/headscale/tags) as the released version you use to
    ensure you have the correct example configuration. The `main` branch might contain unreleased changes.

    === "View on GitHub"

        - Development version: <https://github.com/arsydoni4326-alt/headscale/blob/main/config-example.yaml>
        - Version {{ headscale.version }}: https://github.com/arsydoni4326-alt/headscale/blob/v{{ headscale.version }}/config-example.yaml

    === "Download with `wget`"

        ```shell
        # Development version
        wget -O config.yaml https://raw.githubusercontent.com/juanfont/headscale/main/config-example.yaml

        # Version {{ headscale.version }}
        wget -O config.yaml https://raw.githubusercontent.com/juanfont/headscale/v{{ headscale.version }}/config-example.yaml
        ```

    === "Download with `curl`"

        ```shell
        # Development version
        curl -o config.yaml https://raw.githubusercontent.com/juanfont/headscale/main/config-example.yaml

        # Version {{ headscale.version }}
        curl -o config.yaml https://raw.githubusercontent.com/juanfont/headscale/v{{ headscale.version }}/config-example.yaml
        ```

## Headplane Configuration

### Password Authentication

The `headplane.password` configuration option enables password-based authentication for the Headplane web interface.

#### Configuration File

Add the `headplane` section to your `config.yaml`:

```yaml
# Headplane web UI configuration
headplane:
  # Password for Headplane web UI login.
  # This password grants access to the Headplane web interface only,
  # not the Headscale API. For API access, use API keys.
  #
  # ⚠️  SECURITY WARNING:
  # - Change this password before deploying to production
  # - Use a strong, randomly generated password (minimum 16 characters)
  # - Never commit passwords to version control
  # - Use HTTPS in production to protect credentials in transit
  #
  # Can also be set via the HEADSCALE_HEADPLANE_PASSWORD environment variable,
  # which takes precedence over this config file value.
  password: "your-secure-password-here"
```

#### Environment Variable

Alternatively, set the password via the `HEADSCALE_HEADPLANE_PASSWORD` environment variable:

```bash
export HEADSCALE_HEADPLANE_PASSWORD="your-secure-password-here"
```

**Precedence:** If both the configuration file and environment variable are set, the **environment variable takes precedence**.

#### Security Considerations

!!! danger "Production Deployment"
    - **Always change the default password** before deploying to production
    - Use a strong, randomly generated password (minimum 16 characters)
    - Never commit passwords to version control
    - Use environment variables or secret management systems for production
    - **Always use HTTPS** in production to prevent credential interception

#### Authentication Details

- **Scope:** Password authentication is for Headplane web UI access only
- **Rate Limiting:** 5 login attempts per minute per IP address
- **Session Duration:** 24 hours from login
- **Token Format:** Session tokens start with `hp_` prefix
- **API Endpoint:** `POST /api/v1/headplane/login`

For complete authentication documentation, see the [Authentication Guide](../usage/authentication.md).

#### Changing the Password

To change the Headplane password:

1. Update the `headplane.password` field in `config.yaml` or the `HEADSCALE_HEADPLANE_PASSWORD` environment variable
2. Restart Headscale: `systemctl restart headscale`
3. Existing sessions remain valid until they expire (24 hours)

#### Troubleshooting

Common issues:

- **"Invalid password" error:** Verify password matches configuration exactly; check for whitespace
- **Rate limiting (429):** Wait 1 minute between login attempts
- **Password change not working:** Ensure Headscale was restarted after configuration change

See the [Troubleshooting Guide](../../troubleshooting.md#authentication-issues) for more details.


# Configuration

- Headscale loads its configuration from a YAML, JSON or HuJSON file
- It searches for `config.json`, `config.yaml`, `config.yml` or `config.hujson`, in that order, in the following paths:
    - `/etc/headscale`
    - `$XDG_CONFIG_HOME/headscale` (default `$HOME/.config/headscale`)
    - `headscale` under each directory in `$XDG_CONFIG_DIRS` (default `/etc/xdg`)
    - `$HOME/.headscale`
    - the current working directory
- Configuration keys are always lower case. An unknown or removed configuration entry prevents Headscale from starting
- An unknown or removed `HEADSCALE_*` environment variable logs a warning
- A [secret](configuration-options.md#secrets) can be either a string value or a file path starting with `file:`:
    - an empty file is an error
    - file is only read where the secret is used
- To load the configuration from a different path, use:
    - the command line flag `-c`, `--config`
    - the environment variable `HEADSCALE_CONFIG`
- The configuration can be checked with: `headscale check config`. Use `-f <path>` to load from path.
- Additional configuration can be checked with: `headscale check all`:
    - [Policy](configuration-options.md#policy.path)
    - [DERP Maps](configuration-options.md#derp.paths)
    - [TLS Certs](configuration-options.md#tls.cert_path)
    - [Extra records](configuration-options.md#dns.extra_records_path)

## Example configuration

The full configuration for your version can be generated using `headscale generate config` and will print a
configuration file with all configuration options, their documentation and the default values.

Write it to a file to start a new configuration:

```shell
headscale generate config > /etc/headscale/config.yaml
```

Use `--no-docs` to print a dense configuration file without the comments/documentation.

Use `--minimal` to print the most essential configuration options.

Use `--running` to print the current running configuration of the live Headscale server, with secrets redacted.

[Configuration options](configuration-options.md) describes every key, with examples for the configuration file, the
environment and NixOS.

??? example "Example configuration of this version"

    ```yaml
    {% include "./config-example.yaml" start="do not edit.\n" %}
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

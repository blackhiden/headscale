# Phase 13 Agent Instructions: Documentation

## Task Assignment

**Agent:** Documentation  
**Worktree:** `/home/denny/Project/headscale-password-docs`  
**Branch:** `feature/password-docs`  
**Base Commit:** `b1495ddf5d26ec4089cd6459a491e33cffcf60ed`

---

## Objective

Update all relevant documentation to cover password-based login for Headplane, including setup, configuration, security, and usage.

---

## Roadmap Reference

See `ROADMAP.md` Phase 13 for complete requirements.

---

## Documentation Requirements

### 1. User Documentation
- Add setup guide for password authentication
- Document configuration options (config.yaml and environment variables)
- Document login process for end users
- Document password management (changing password)
- Add troubleshooting section

### 2. Security Documentation
- Document security considerations
- Warn about default passwords
- Document HTTPS requirement
- Document rate limiting behavior
- Best practices for production deployment

### 3. API Documentation
- Document new `/api/v1/headplane/login` endpoint (or equivalent)
- Document request/response format
- Document error codes
- Clarify difference between API key and password auth
- Add examples

### 4. Configuration Documentation
- Update config reference to include `headplane.password`
- Document environment variable `HEADSCALE_HEADPLANE_PASSWORD`
- Update config-example.yaml comments
- Document default behavior

### 5. README Updates
- Update main README if needed
- Mention password login as an authentication option
- Link to detailed documentation

---

## Expected Files to Change

- `docs/setup/headplane.md` or similar (new section on password auth)
- `docs/ref/configuration.md` (add headplane.password field)
- `docs/usage/authentication.md` or similar (new section)
- `docs/troubleshooting.md` (add password auth issues)
- `README.md` (mention password login)
- `config-example.yaml` (comments for password field)
- OpenAPI spec if applicable (add new endpoint)

---

## Documentation Structure

Create a new documentation file if needed:

```markdown
# Headplane Password Authentication

## Overview

Headplane supports password-based authentication as an alternative to API keys...

## Configuration

### Using config.yaml

```yaml
headplane:
  password: "your-secure-password"
```

### Using Environment Variables

```bash
export HEADSCALE_HEADPLANE_PASSWORD="your-secure-password"
```

## Usage

### Logging In

1. Navigate to Headplane URL
2. Select "Password" authentication method
3. Enter your password
4. Click "Log In"

## Security

⚠️ **Important:** Change the default password in production...

## Troubleshooting

### "Invalid password" error
...

### Rate limiting
...
```

---

## Acceptance Criteria

- [ ] All password authentication features are documented
- [ ] Configuration options are clear and complete
- [ ] Security warnings are prominent
- [ ] API endpoint is documented
- [ ] Examples are provided
- [ ] Troubleshooting covers common issues
- [ ] Links are correct and working
- [ ] Documentation follows existing style and structure
- [ ] No broken links or references

---

## Implementation Steps

1. Read backend implementation to understand exact behavior
2. Read frontend implementation to understand user flow
3. Review existing Headscale/Headplane documentation structure
4. Create/update documentation files
5. Add code examples
6. Add configuration examples
7. Add troubleshooting section
8. Update API reference
9. Update README
10. Review for completeness and accuracy
11. Test all examples
12. Commit to feature branch

---

## Development Workflow

```bash
# Navigate to worktree
cd /home/denny/Project/headscale-password-docs

# Verify branch
git branch --show-current  # should show: feature/password-docs

# Make changes
# Edit docs/ files

# Preview documentation (if using mkdocs)
mkdocs serve

# Commit changes
git add -A
git commit -m "docs: add password authentication documentation for Headplane"

# Push when ready
git push -u origin feature/password-docs
```

---

## Dependencies

- Backend implementation (for accurate API documentation)
- Frontend implementation (for accurate user flow documentation)

---

## Style Guide

- Follow existing Headscale documentation style
- Use clear, concise language
- Include code examples with syntax highlighting
- Use admonitions for warnings and tips
- Cross-reference related documentation
- Ensure examples are copy-pasteable

---

## Notes

- Wait for backend and frontend implementation to be complete before finalizing docs
- Verify all examples work
- Check for consistency with existing documentation
- Ensure security warnings are prominent
- Test all configuration examples

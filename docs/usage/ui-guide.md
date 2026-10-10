---
hide:
  - toc
---

# Headscale UI Guide

This guide provides visual references and flow diagrams for the Headscale
user interface. It is intended for contributors and maintainers who need
to understand, update, or test the UI.

> **Note:** Screenshots in this guide are tracked as a [maintenance task](https://github.com/arsydoni4326-alt/headscale/issues). When the UI changes, contributors should update the corresponding images to keep this guide accurate. See [How to take and update screenshots](#how-to-take-and-update-screenshots) below.

---

## UI Screenshots

Screenshots are stored under
[`docs/assets/screenshots/`](../assets/screenshots/). Each page or feature
has its own subdirectory when needed.

### How to take and update screenshots

1. Use a consistent browser (Chrome/Firefox latest) and OS (Linux/macOS).
2. Capture both **desktop** and **mobile** views where relevant.
   - Desktop: viewport at least 1280×800 px.
   - Mobile: viewport 375×812 px (iPhone X size).
3. Name files clearly, e.g.:
   - `ping-default-desktop.png`
   - `ping-error-mobile.png`
   - `register-confirm-default-desktop.png`
4. When you make a UI change, update the affected screenshots in the same
   pull request.

### Pages and states

<!--
  Contributors: replace each placeholder with the actual screenshot.
  Use the pattern:
    ![Description](../assets/screenshots/<filename>)
-->

#### Apple Configuration (`/apple`)

| State     | Desktop | Mobile |
|-----------|---------|--------|
| Default   | _Screenshot placeholder_ | _Screenshot placeholder_ |

#### Windows Configuration (`/windows`)

| State     | Desktop | Mobile |
|-----------|---------|--------|
| Default   | _Screenshot placeholder_ | _Screenshot placeholder_ |

#### Authentication Success

| State     | Desktop | Mobile |
|-----------|---------|--------|
| Default   | _Screenshot placeholder_ | _Screenshot placeholder_ |

#### Authentication Error

| State     | Desktop | Mobile |
|-----------|---------|--------|
| Default   | _Screenshot placeholder_ | _Screenshot placeholder_ |

#### Ping / Debug (`/debug/ping`)

| State     | Desktop | Mobile |
|-----------|---------|--------|
| Default (no result) | _Screenshot placeholder_ | _Screenshot placeholder_ |
| Success (pong) | _Screenshot placeholder_ | _Screenshot placeholder_ |
| Error / Timeout | _Screenshot placeholder_ | _Screenshot placeholder_ |

#### Registration Confirmation (`/register/confirm`)

| State     | Desktop | Mobile |
|-----------|---------|--------|
| Default   | _Screenshot placeholder_ | _Screenshot placeholder_ |

#### Machine Approval

!!! note "Coordinate with frontend agent"
    The exact UI workflow, terminology, permissions, and component names need to be
    confirmed with the frontend UI agent. Placeholders below represent expected
    behavior based on the task specification.

| State     | Desktop | Mobile |
|-----------|---------|--------|
| Pending Machines List | _Screenshot placeholder_ | _Screenshot placeholder_ |
| Single Machine Approval | _Screenshot placeholder_ | _Screenshot placeholder_ |
| Bulk Approval | _Screenshot placeholder_ | _Screenshot placeholder_ |
| Approval Success | _Screenshot placeholder_ | _Screenshot placeholder_ |
| Approval Error | _Screenshot placeholder_ | _Screenshot placeholder_ |

---

## UI Flow Diagrams

Diagrams are stored under
[`docs/assets/diagrams/`](../assets/diagrams/). They illustrate the main
user-facing flows.

### How to create and update diagrams

1. Use **Mermaid** (native mkdocs support) for inline diagrams, or
   draw.io / similar for complex flows.
2. Save source files alongside the rendered PNG/SVG when possible.
3. Update diagrams when the corresponding user flow changes.

### Device registration flow

```mermaid
flowchart LR
    A[Register request] --> B[Confirm page]
    B -->|Accept| C[Success page]
    B -->|Reject / close| D[Error / Expire]
```

### Machine approval flow

!!! note "Coordinate with frontend agent"
    This workflow diagram represents the expected behavior. Final details should be
    confirmed with the frontend UI agent.

```mermaid
flowchart LR
    A[Machine requests registration] --> B[Pending state]
    B --> C{Approval method}
    C -->|Web UI| D[View pending machines]
    C -->|CLI| E[headscale auth register]
    C -->|API| F[POST /api/v1/machines/.../approve]
    D --> G[Single or bulk approve]
    G --> H[Machine approved]
    E --> H
    F --> H
    H --> I[Machine joins network]
```

### Authentication flow

```mermaid
flowchart LR
    A[Login] --> B{OIDC / CLI}
    B -->|OIDC| C[Auth page]
    C -->|Success| D[Success page]
    C -->|Failure| E[Error page]
    B -->|CLI| F[Token / command page]
```

### Ping / Debug flow

```mermaid
flowchart LR
    A[Open /debug/ping] --> B[Enter node query]
    B --> C[POST form]
    C --> D{Result}
    D -->|Success| E[Success box]
    D -->|Timeout| F[Warning box]
    D -->|Error| G[Error box]
```

---

## Accessibility Features

Headscale's web interfaces (both backend-served templates and the Headplane
UI) are designed to be accessible to all users. This section documents the
accessibility features and provides guidance for contributors.

### Keyboard Navigation

All interactive elements can be accessed and operated using only a keyboard:

- **Tab / Shift+Tab** — Navigate between interactive elements
- **Enter / Space** — Activate buttons, links, and form controls
- **Arrow keys** — Navigate within select menus, radio groups, and custom controls
- **Esc** — Close modals and dialogs

All interactive elements have visible focus indicators (rings) to show the
current keyboard focus position.

### Screen Reader Support

All UI elements provide appropriate semantic information for screen readers:

- **Semantic HTML** — Proper use of `<button>`, `<form>`, `<nav>`, heading hierarchy
- **ARIA labels** — All icon-only buttons and controls have accessible names via `aria-label`
- **ARIA roles** — Dynamic content uses `role="alert"` and `aria-live` for announcements
- **Alternative text** — Images and icons have descriptive text or are marked decorative with `aria-hidden="true"`

### Visual Accessibility

- **Color contrast** — All text and interactive elements meet WCAG 2.1 AA contrast requirements (4.5:1 for normal text, 3:1 for large text)
- **Dark mode** — Full dark mode support with appropriate contrast in both themes
- **Focus indicators** — High-visibility focus rings on all interactive elements
- **Link styling** — Links are underlined to avoid relying on color alone
- **Touch targets** — Interactive elements meet minimum size requirements (44×44px on mobile)

### Forms and Validation

- **Associated labels** — All form inputs have visible labels or `aria-label` attributes
- **Error messages** — Validation errors are announced to screen readers via `aria-live` regions
- **Required fields** — Required inputs are marked and announced appropriately
- **Placeholder text** — Used only for hints, never as the sole label

### Testing

The Headplane UI includes automated accessibility testing:

- **axe-core** — Automated scans for critical and serious WCAG violations in CI/CD
- **Playwright e2e tests** — Include accessibility checks on all major pages
- **Component tests** — Verify keyboard navigation, focus management, and ARIA attributes

### Known Limitations

- **Full WCAG validation** — Automated testing catches many issues, but full WCAG 2.1 AA conformance requires manual testing with assistive technologies
- **Dynamic content** — Some complex interactions (drag-and-drop in topology view) may have reduced screen reader support
- **Third-party components** — Some upstream components (`@base-ui/react`) use non-standard patterns (e.g., `aria-disabled` instead of `disabled` attribute)

### Reporting Accessibility Issues

If you encounter an accessibility barrier, please report it via the
[GitHub issue tracker](https://github.com/arsydoni4326-alt/headscale/issues).
Include:

- Description of the issue
- The page/component affected
- Your assistive technology (screen reader, keyboard-only navigation, etc.)
- Steps to reproduce
- Expected behavior

---

## Standardized UI States

The Headplane UI uses standardized components for empty, loading, and error
states across all routes. These components ensure consistency and
accessibility throughout the application.

### EmptyState Component

The `EmptyState` component displays when no data is available. It provides
clear messaging about why the state is empty and what the user can do next.

**Features:**

- **Clear messaging** — Title and description explain the empty state
- **Actionable** — Optional primary and secondary action buttons
- **Variants** — Three visual variants: `default`, `filtered`, `error`
- **Accessible** — Icons are decorative (`aria-hidden`), text is semantic

**Usage examples:**

- **No machines yet** — First-time user sees instructions on adding their first device
- **No results** — User's filters don't match any items; "Clear Filters" button is provided
- **API error** — Connection or permission error with appropriate visual treatment

### LoadingSpinner Component

The `LoadingSpinner` component displays during data fetching operations.

**Features:**

- **Accessible loading indicator** — Uses `role="status"` and `aria-live="polite"`
- **Optional label** — Describes what is being loaded (e.g., "Loading machines...")
- **Size variants** — `sm`, `md`, `lg` to fit different contexts
- **Semantic markup** — Icon is decorative, state is announced to screen readers

**Accessibility:**

Loading states are announced to screen readers via ARIA live regions without
interrupting the user's current task (`aria-live="polite"`).

### ErrorBanner Component

The `ErrorBanner` component displays error messages with appropriate context
and recovery actions. It provides detailed information about API errors,
connection failures, and unexpected exceptions.

**Features:**

- **Contextual error messages** — Tailored to the error type (API, connection, validation)
- **Recovery guidance** — Clear instructions on what to do next
- **Accessible** — Error content is semantic HTML, announced to screen readers
- **Detailed debugging** — Technical details in code blocks for troubleshooting

### Consistency Guidelines

When implementing new routes or features:

1. **Empty states** — Use `EmptyState` component instead of plain text
2. **Loading states** — Use `LoadingSpinner` with a descriptive label
3. **Error states** — Use `ErrorBanner` for errors and `PageError` for route-level errors
4. **Filtered empty states** — Use `EmptyState` variant `filtered` with a "Clear Filters" action
5. **Icon-only buttons** — Always provide `aria-label` or visible text via `.sr-only`

---

## Machine Approval Workflow

!!! note "Coordinate with frontend agent"
    The specific UI pages, routes, components, terminology, and interaction patterns
    below are placeholders. Final details must be confirmed with the frontend UI agent
    to ensure accuracy.

### Overview

The machine approval feature allows administrators to review and approve machines
that have initiated web authentication. Pending machines can be approved individually
or in bulk through the Headplane web interface.

### Viewing Pending Machines

**Expected route:** `/machines` or `/machines?status=pending`

**Steps:**

1. Navigate to the Machines page in Headplane
2. Filter or view machines with "Pending" status
3. Review machine details: hostname, user, IP address, last seen timestamp

**UI elements:**

- Status badge showing "Pending" state
- Machine metadata (hostname, requested user, tags if applicable)
- Timestamp of registration request
- Actions: Approve button, Reject button (if supported)

### Approving a Single Machine

**Steps:**

1. From the pending machines list, locate the machine to approve
2. Click the "Approve" button or action menu
3. Confirm the user assignment (if required)
4. The machine transitions to "Approved" status
5. Success notification appears confirming approval

**Expected behavior:**

- Optimistic UI update (machine moves from pending to approved immediately)
- API call to `POST /api/v1/machines/{id}/approve`
- Error handling if approval fails (rollback UI state, show error message)
- Success feedback (toast notification or inline message)

### Bulk Approval

**Steps:**

1. From the pending machines list, select multiple machines using checkboxes
2. Click "Approve Selected" or bulk action button
3. Confirm bulk approval (if required)
4. All selected machines transition to "Approved" status
5. Success notification shows count of approved machines

**Expected behavior:**

- Bulk selection UI with "Select All" and individual checkboxes
- Disabled bulk actions when no machines are selected
- API call to `POST /api/v1/machines/approve` with array of machine IDs
- Partial success handling (some machines approved, others failed)
- Detailed error reporting for any failed approvals

### Permissions

**Required permissions:** (TBD - confirm with frontend agent)

- Admin or machine management role
- API key with appropriate scopes
- User must have permission to assign machines to the specified user

**Behavior when unauthorized:**

- Approve buttons are hidden or disabled
- Attempting approval shows permission error
- User is directed to contact an administrator

### Troubleshooting

**Common issues:**

| Issue | Cause | Solution |
|-------|-------|----------|
| "Permission denied" error | Insufficient API key permissions | Contact administrator to verify API key scopes |
| Machine not appearing in pending list | Already approved or expired | Check all machines view or registration logs |
| Bulk approval partially failed | Some machines in invalid state | Review error details, retry failed machines individually |
| "Machine not found" error | Machine removed or auth session expired | Machine must re-initiate registration |

**Debug steps:**

1. Check browser console for API errors
2. Verify API key is valid and not expired
3. Confirm machine is in pending state via CLI: `headscale nodes list`
4. Review Headscale server logs for detailed error messages
5. Test API endpoint directly with curl (see [API docs](../ref/api.md#machine-approval))

---

## Contributing to this guide

- **Add a new page:** Follow the table format in [Pages and states](#pages-and-states).
  Add the screenshot files and update the table.
- **Update a screenshot:** Replace the file in
  `docs/assets/screenshots/` and keep the same filename if the content
  hasn't changed structurally.
- **Update a diagram:** Edit the Mermaid block or replace the file in
  `docs/assets/diagrams/`.
- **Reference screenshots in other docs:** Use relative links:
  ```markdown
  ![Ping page](../assets/screenshots/ping-default-desktop.png)
  ```

### Maintenance checklist

Before each release, verify that:

- [ ] All screenshots reflect the current UI.
- [ ] No broken image links exist.
- [ ] Flow diagrams match the current user experience.
- [ ] Accessibility tests pass (run `pnpm run test:a11y` in the Headplane directory).
- [ ] Empty/loading/error states use standardized components.
- [ ] Any new UI pages or states are documented.
- [ ] Accessibility features remain functional and documented.

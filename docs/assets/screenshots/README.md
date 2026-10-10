# Screenshot Placeholders

This directory contains placeholder references for screenshots that document the Known Gap fix in Phase 13c.

## Required Screenshots

### 1. admin-users-api-key.png
**Description**: Admin users page showing full user management UI when authenticated via API key

**What to capture**:
- Navigate to Headplane and log in with an admin API key
- Go to `/admin/users` page
- Capture the full page showing:
  - "Add Headplane User" button visible
  - List of Headplane dashboard users with their roles
  - Edit and delete buttons for each user
  - Sidebar showing admin is logged in via API key

**Dimensions**: Approximately 1920x1080 or similar standard resolution

---

### 2. add-headplane-user-dialog.png
**Description**: Dialog for creating a new Headplane dashboard user

**What to capture**:
- From the `/admin/users` page (authenticated via API key)
- Click "Add Headplane User" button
- Capture the modal/dialog showing:
  - Username field
  - Password field
  - Role selector (admin/user dropdown)
  - Save and Cancel buttons

**Dimensions**: Crop to show just the dialog and some background context

---

### 3. edit-headplane-user-dialog.png
**Description**: Dialog for editing an existing Headplane dashboard user

**What to capture**:
- From the `/admin/users` page
- Click edit button on an existing user
- Capture the modal/dialog showing:
  - Username field (pre-filled)
  - Password field (optional update)
  - Role selector showing current role
  - Save and Cancel buttons

**Dimensions**: Crop to show just the dialog and some background context

---

## Usage in Documentation

These screenshots are referenced in:
- `docs/ref/api/headplane-users.md`
- `docs/usage/authentication.md`

## Notes

- Screenshots should be taken from the actual implementation after the Known Gap fix is deployed
- Use a consistent browser and theme (preferably light theme for better visibility)
- Ensure no sensitive information (real API keys, passwords, etc.) is visible
- PNG format with reasonable compression
- Consider adding annotations (arrows, highlights) if needed for clarity

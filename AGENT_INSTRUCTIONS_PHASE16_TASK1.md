# Phase 16 Task 1: Admin Menu Navigation Fix

**Branch:** `feature/phase16-admin-nav`  
**Worktree:** `/home/denny/Project/headscale-project/headscale-phase16-admin-nav`  
**Base Commit:** `813c46df` (dev branch)

---

## Objective

Fix the Admin menu navigation to route to `/admin/admin` instead of `/admin/admin/users`.

---

## Problem Statement

The Admin menu item in the main navigation bar incorrectly routes to `/admin/admin/users` instead of `/admin/admin`. This causes confusion and breaks the expected navigation hierarchy.

**Current behavior:**
```html
<a href="/admin/admin/users">
  <svg>...</svg>Admin
</a>
```

**Expected behavior:**
```html
<a href="/admin/admin">
  <svg>...</svg>Admin
</a>
```

---

## Requirements

1. Update the Admin menu item `href` from `/admin/admin/users` to `/admin/admin`
2. Verify the Admin page (`/admin/admin`) has proper default view
3. Test navigation flow: clicking Admin should land on the Admin overview page
4. Ensure backward compatibility if any bookmarks or external links reference the old path

---

## Implementation Steps

### 1. Locate Navigation Component

Find the navigation component in the Headplane submodule:
```bash
cd headplane
find app -name "*.tsx" -o -name "*.ts" | xargs grep -l "admin/admin/users"
```

Likely locations:
- `headplane/app/components/navigation/*.tsx`
- `headplane/app/layout/*.tsx`
- `headplane/app/routes/admin/layout.tsx`

### 2. Update the Navigation Link

Change the `href` attribute:
```typescript
// Before:
<a href="/admin/admin/users">
  <svg>...</svg>Admin
</a>

// After:
<a href="/admin/admin">
  <svg>...</svg>Admin
</a>
```

### 3. Verify Admin Page Route

Check that `/admin/admin` route exists and has proper content:
```bash
cd headplane
ls -la app/routes/admin/admin/
```

If the route file is named `users.tsx`, it may need to be renamed to `index.tsx` or `page.tsx` depending on React Router 7 conventions.

### 4. Test the Fix

**Manual testing:**
1. Start Headplane: `cd headplane && pnpm run dev`
2. Login to Headplane
3. Click the Admin menu item
4. Verify URL is `/admin/admin`
5. Verify page loads correctly

**Automated testing:**
- Check if existing E2E tests cover admin navigation
- Update or add tests if needed

---

## Files Expected to Change

- `headplane/app/layout/*.tsx` or `headplane/app/components/navigation/*.tsx` (navigation component)
- Possibly `headplane/app/routes/admin/admin/page.tsx` or routing config

---

## Acceptance Criteria

- [x] Admin menu item routes to `/admin/admin`
- [x] Admin page loads correctly at `/admin/admin`
- [x] No console errors or warnings
- [x] Existing admin functionality remains intact
- [x] Navigation breadcrumb (if present) reflects correct hierarchy

---

## Testing Checklist

- [ ] Click Admin menu → lands on `/admin/admin`
- [ ] Page displays expected content
- [ ] No JavaScript errors in console
- [ ] Navigation state properly reflects active menu item
- [ ] Bookmarks to `/admin/admin/users` still work (redirect or display content)

---

## Commit Message Format

```
headplane: fix admin menu navigation route

Change Admin menu href from /admin/admin/users to /admin/admin
to match expected navigation hierarchy.

Fixes Phase 16 Issue 1.

Refs: ROADMAP.md Phase 16
```

---

## Notes

- This is a **low-risk, quick-win** change
- Estimated time: 30 minutes - 1 hour
- No backend changes required
- Focus on Headplane submodule only

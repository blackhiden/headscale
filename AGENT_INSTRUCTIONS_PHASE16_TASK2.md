# Phase 16 Task 2: User Avatar Display

**Branch:** `feature/phase16-avatar`  
**Worktree:** `/home/denny/Project/headscale-project/headscale-phase16-avatar`  
**Base Commit:** `813c46df` (dev branch)

---

## Objective

Add user avatar display support from `config.yaml` with fallback to SVG icon.

---

## Problem Statement

The user menu button displays a generic SVG icon instead of the user's avatar picture. The avatar URL should be retrieved from `config.yaml` and displayed when available.

**Current behavior:**
```html
<button>
  <svg class="lucide-circle-user">...</svg>
</button>
```

**Expected behavior:**
```html
<button>
  <img src="{avatar_url_from_config}" alt="User avatar" />
  <!-- Fallback to SVG if avatar not configured -->
</button>
```

---

## Requirements

1. Add optional `user.avatar` field to Headplane config schema
2. Expose avatar URL via API or authentication response
3. Update user menu component to display avatar image when available
4. Fallback to SVG icon when avatar is not configured or fails to load
5. Apply proper styling (circular avatar, appropriate sizing)
6. Handle image loading errors gracefully

---

## Implementation Steps

### 1. Update Config Schema

Add avatar field to Headplane config:
```typescript
// headplane/app/server/config/config-schema.ts
user?: {
  username: string;
  password: string;
  name?: string;
  avatar?: string; // URL or path to avatar image
}
```

### 2. Expose Avatar in API

Option A: Include in settings endpoint response:
```typescript
// headplane/app/server/api/settings.ts or similar
export async function getSettings() {
  return {
    username: config.user.username,
    name: config.user.name,
    avatar: config.user.avatar, // Add this
    // ... other settings
  };
}
```

Option B: Include in auth response during login.

### 3. Update User Menu Component

Locate the user menu component:
```bash
cd headplane
find app -name "*.tsx" | xargs grep -l "lucide-circle-user"
```

Update to conditionally render avatar:
```typescript
// Pseudocode
{avatar ? (
  <img 
    src={avatar} 
    alt="User avatar"
    className="w-8 h-8 rounded-full"
    onError={(e) => {
      // Fallback to SVG on error
      e.currentTarget.style.display = 'none';
      setShowFallback(true);
    }}
  />
) : (
  <svg className="lucide-circle-user">...</svg>
)}
```

### 4. Add Error Handling

Handle image load failures:
- Network errors
- Invalid URLs
- CORS issues
- Missing images

Always fallback to SVG icon.

### 5. Update Example Config

Add avatar example to `headplane/config.example.yaml`:
```yaml
user:
  username: admin
  password: "$2b$12$..."
  name: "Administrator"  # Optional
  avatar: "https://example.com/avatar.jpg"  # Optional
```

### 6. Update Documentation

Update `headplane/docs/CONFIGURATION.md`:
- Document `user.avatar` field
- Explain supported formats (URLs)
- Note that avatar is optional

---

## Files Expected to Change

- `headplane/app/server/config/config-schema.ts` (config schema)
- `headplane/app/components/user-menu/*.tsx` or `headplane/app/layout/header.tsx` (user menu)
- `headplane/app/server/api/settings.ts` or auth handler (if avatar in API response)
- `headplane/config.example.yaml` (documentation)
- `headplane/docs/CONFIGURATION.md` (documentation)

---

## Acceptance Criteria

- [x] Config with avatar URL displays image
- [x] Config without avatar URL displays SVG fallback
- [x] Invalid image URL displays SVG fallback
- [x] Image has proper alt text and accessibility
- [x] Avatar is circular and properly sized
- [x] No layout shift during image load
- [x] Example config documents avatar field

---

## Testing Checklist

### With Avatar URL
- [ ] Image displays correctly
- [ ] Circular styling applied
- [ ] Proper sizing (matches design)
- [ ] Alt text present

### Without Avatar URL
- [ ] SVG icon displays
- [ ] No console errors
- [ ] No broken image placeholder

### Error Scenarios
- [ ] Invalid URL → SVG fallback
- [ ] 404 image → SVG fallback
- [ ] CORS blocked image → SVG fallback
- [ ] Slow-loading image → loading state or immediate fallback

### Accessibility
- [ ] Alt text describes avatar
- [ ] Keyboard navigation works
- [ ] Screen reader announces correctly

---

## Commit Message Format

```
headplane: add user avatar display from config

Add optional user.avatar field to config schema and display
avatar image in user menu when configured. Falls back to
SVG icon when avatar is not configured or fails to load.

Features:
- Avatar URL from config.yaml (user.avatar)
- Circular avatar styling
- Graceful error handling with SVG fallback
- Accessibility support

Fixes Phase 16 Issue 2.

Refs: ROADMAP.md Phase 16
```

---

## Notes

- **Medium complexity** change
- Estimated time: 2-3 hours
- Requires config schema change and UI update
- No backend API changes required (uses existing settings endpoint)
- Focus on graceful degradation

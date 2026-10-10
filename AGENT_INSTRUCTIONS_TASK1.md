# Task 1: Core Module Path Rewrite

**Worktree:** `/home/denny/Project/headscale-project/headscale-module-path-core`  
**Branch:** `feature/module-path-core`  
**Roadmap Item:** Phase 14 — Module Path Rewrite  
**Sequence:** 1 (must complete first)  
**Dependencies:** None

---

## Objective

Update the root Go module path and all core Go source imports from 
`github.com/juanfont/headscale` to `github.com/arsydoni4326-alt/headscale`.

---

## Scope — What to Change

1. **`go.mod`**: Update the `module` directive
2. **All `.go` files** in the repository (except `headplane/` submodule):
   - Update all import statements
   - Search for: `github.com/juanfont/headscale`
   - Replace with: `github.com/arsydoni4326-alt/headscale`

---

## Scope — What NOT to Change

- Do NOT modify `headplane/` submodule (separate task)
- Do NOT modify `Dockerfile*` files (separate task)
- Do NOT modify `Makefile` (separate task)
- Do NOT modify documentation files (separate task)
- Do NOT modify Nix files (separate task)
- Do NOT run `go mod tidy` yet (will be done later)
- Do NOT update `go.sum` manually

---

## Implementation Steps

### 1. Change to worktree directory

```bash
cd /home/denny/Project/headscale-project/headscale-module-path-core
```

### 2. Verify clean state

```bash
git status
git branch --show-current  # Should be feature/module-path-core
```

### 3. Update go.mod

```bash
sed -i 's|module github.com/juanfont/headscale|module github.com/arsydoni4326-alt/headscale|' go.mod
```

Verify:
```bash
head -3 go.mod
# Should show: module github.com/arsydoni4326-alt/headscale
```

### 4. Update all Go imports (excluding headplane)

```bash
find . -name "*.go" -type f -not -path "./headplane/*" -exec sed -i \
  's|github.com/juanfont/headscale|github.com/arsydoni4326-alt/headscale|g' {} +
```

### 5. Verify no old references remain

```bash
grep -r "github.com/juanfont/headscale" --include="*.go" --exclude-dir=headplane
```

**Expected result:** No matches (empty output)

If any matches found, review manually and update if needed.

### 6. Quick syntax check

```bash
go list ./... 2>&1 | head -20
```

**Expected:** Module path errors are OK at this stage (dependencies not updated yet).  
Just verify Go can parse the files without syntax errors.

---

## Validation Checklist

- [ ] `go.mod` module directive shows new path
- [ ] All `.go` files (except `headplane/`) updated
- [ ] No references to old path in core Go files
- [ ] Go files parse without syntax errors
- [ ] Only Go files and go.mod changed (check `git status`)

---

## Commit

### Before committing

```bash
git status
# Review changed files - should only be .go files and go.mod
```

### Commit message

```
module: rewrite path to github.com/arsydoni4326-alt/headscale

Update the root go.mod module directive and all Go import statements
from github.com/juanfont/headscale to github.com/arsydoni4326-alt/headscale.

This is part of Phase 14 to properly reflect the fork's independent
identity and fix version reporting in built binaries.

Refs: ROADMAP.md Phase 14
```

### Execute commit

```bash
git add go.mod
git add $(find . -name "*.go" -type f -not -path "./headplane/*")
git commit -m "module: rewrite path to github.com/arsydoni4326-alt/headscale

Update the root go.mod module directive and all Go import statements
from github.com/juanfont/headscale to github.com/arsydoni4326-alt/headscale.

This is part of Phase 14 to properly reflect the fork's independent
identity and fix version reporting in built binaries.

Refs: ROADMAP.md Phase 14"
```

### Record commit SHA

```bash
git rev-parse HEAD
# Save this SHA - other tasks will merge from this commit
```

---

## After Completion

1. **Record commit SHA**: `_________________`
2. **Changed files count**: Check `git show --stat`
3. **Any issues encountered**: Document below

---

## Notes

- This is the foundational change that other tasks depend on
- Other worktrees will merge this commit before proceeding
- Do not attempt a full build yet; other files still need updating
- The build will fail until Makefile, Dockerfiles, and other files are updated

---

## Troubleshooting

**Issue:** `grep` still finds old references  
**Solution:** Review those files manually - they may be in comments or strings that need updating

**Issue:** `go list` reports errors  
**Solution:** This is expected - module dependencies haven't been updated yet. Just ensure no syntax errors.

**Issue:** `git status` shows unexpected files  
**Solution:** Only .go files and go.mod should be changed. If others appear, do not commit them.

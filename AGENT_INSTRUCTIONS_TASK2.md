# Task 2: Headplane Submodule Rewrite

**Worktree:** `/home/denny/Project/headscale-project/headscale-module-path-headplane`  
**Branch:** `feature/module-path-headplane`  
**Roadmap Item:** Phase 14 — Module Path Rewrite  
**Sequence:** 2 (after core rewrite)  
**Dependencies:** Task 1 (Core Module Path Rewrite)

---

## Objective

Update the `headplane/` Git submodule's Go module path and any references to the 
parent Headscale module.

---

## Scope — What to Change

1. **`headplane/go.mod`**: Update module directive if it references headscale
2. **All `.go` files in `headplane/`**: Update any imports referencing the parent module
3. Any cross-references between headplane and headscale

---

## Scope — What NOT to Change

- Do NOT modify files outside `headplane/`
- Do NOT modify Dockerfiles, Makefiles, or docs
- Do NOT run `go mod tidy` yet

---

## Implementation Steps

### 1. Change to worktree directory

```bash
cd /home/denny/Project/headscale-project/headscale-module-path-headplane
```

### 2. Merge core rewrite

```bash
git merge feature/module-path-core --no-edit
```

**Verify merge succeeded:**
```bash
git log --oneline -5
# Should show the core rewrite commit
```

### 3. Check headplane for references

```bash
cd headplane
grep -r "github.com/juanfont/headscale" --include="*.go" .
```

**Expected:** May find references if headplane imports headscale packages

### 4. Update headplane Go files if needed

```bash
cd headplane
find . -name "*.go" -type f -exec sed -i \
  's|github.com/juanfont/headscale|github.com/arsydoni4326-alt/headscale|g' {} +
```

### 5. Check headplane's go.mod

```bash
cd headplane
cat go.mod | head -5
```

If it contains the old module path, update it:
```bash
sed -i 's|github.com/juanfont/headscale|github.com/arsydoni4326-alt/headscale|g' go.mod
```

### 6. Verify no old references

```bash
cd headplane
grep -r "github.com/juanfont/headscale" --include="*.go" --include="go.mod" .
```

**Expected:** No matches

### 7. Syntax check

```bash
cd headplane
go list ./... 2>&1 | head -20
```

---

## Validation Checklist

- [ ] Core rewrite merged successfully
- [ ] All `.go` files in `headplane/` updated
- [ ] `headplane/go.mod` updated if needed
- [ ] No references to old path in headplane
- [ ] Go syntax valid

---

## Commit

```bash
cd /home/denny/Project/headscale-project/headscale-module-path-headplane
git add headplane/
git commit -m "headplane: update module path references

Update any references to the parent Headscale module from
github.com/juanfont/headscale to github.com/arsydoni4326-alt/headscale
in the headplane submodule.

Refs: ROADMAP.md Phase 14"
```

### Record commit SHA

```bash
git rev-parse HEAD
```

---

## After Completion

1. **Record commit SHA**: `_________________`
2. **Files changed in headplane**: Check `git show --stat`
3. **Any issues**: Document below

# Task 3: Build System Update

**Worktree:** `/home/denny/Project/headscale-project/headscale-module-path-build`  
**Branch:** `feature/module-path-build`  
**Roadmap Item:** Phase 14 — Module Path Rewrite  
**Sequence:** 2 (parallel with tasks 2, 4, 5)  
**Dependencies:** Task 1 (Core Module Path Rewrite)

---

## Objective

Update `Makefile` and all `Dockerfile*` files to use the new module path in 
ldflags and build commands.

---

## Scope — What to Change

1. **`Makefile`**: Update ldflags with new module path
2. **All `Dockerfile*` files**: Update ldflags in build commands
3. **Build scripts**: Any shell scripts with hardcoded module paths

---

## Scope — What NOT to Change

- Do NOT modify Go source files
- Do NOT modify documentation
- Do NOT modify Nix files
- Do NOT modify CI/CD workflows

---

## Implementation Steps

### 1. Change to worktree directory

```bash
cd /home/denny/Project/headscale-project/headscale-module-path-build
```

### 2. Merge core rewrite

```bash
git merge feature/module-path-core --no-edit
```

### 3. Update Makefile

```bash
sed -i 's|github.com/juanfont/headscale|github.com/arsydoni4326-alt/headscale|g' Makefile
```

**Verify:**
```bash
grep "github.com/arsydoni4326-alt/headscale" Makefile
# Should show updated ldflags
```

### 4. Update all Dockerfiles

```bash
find . -name "Dockerfile*" -type f -exec sed -i \
  's|github.com/juanfont/headscale|github.com/arsydoni4326-alt/headscale|g' {} +
```

**List all Dockerfiles to verify:**
```bash
find . -name "Dockerfile*" -type f
```

### 5. Check shell scripts

```bash
grep -r "github.com/juanfont/headscale" --include="*.sh" scripts/ 2>/dev/null || echo "No matches"
```

If matches found, update:
```bash
find scripts/ -name "*.sh" -type f -exec sed -i \
  's|github.com/juanfont/headscale|github.com/arsydoni4326-alt/headscale|g' {} +
```

### 6. Check docker-build.sh if exists

```bash
if [ -f docker-build.sh ]; then
  sed -i 's|github.com/juanfont/headscale|github.com/arsydoni4326-alt/headscale|g' docker-build.sh
fi
```

### 7. Verify no old references

```bash
grep -r "github.com/juanfont/headscale" Makefile Dockerfile* scripts/ 2>/dev/null || echo "All updated"
```

**Expected:** "All updated" or no matches

---

## Validation Checklist

- [ ] Makefile ldflags updated
- [ ] All Dockerfile* files updated
- [ ] Build scripts updated
- [ ] No old path references in build files

---

## Commit

```bash
git add Makefile Dockerfile* scripts/ docker-build.sh 2>/dev/null || true
git commit -m "build: update module path in Makefile and Dockerfiles

Update ldflags and build commands from github.com/juanfont/headscale
to github.com/arsydoni4326-alt/headscale in:
- Makefile
- All Dockerfile* files
- Build scripts

This ensures version information is correctly injected at build time.

Refs: ROADMAP.md Phase 14"
```

### Record commit SHA

```bash
git rev-parse HEAD
```

---

## After Completion

1. **Record commit SHA**: `_________________`
2. **Files updated**: List them
3. **Any issues**: Document below

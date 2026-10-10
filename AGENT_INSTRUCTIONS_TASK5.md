# Task 5: Nix and Tooling Update

**Worktree:** `/home/denny/Project/headscale-project/headscale-module-path-nix`  
**Branch:** `feature/module-path-nix`  
**Roadmap Item:** Phase 14 — Module Path Rewrite  
**Sequence:** 2 (parallel with tasks 2, 3, 4)  
**Dependencies:** Task 1 (Core Module Path Rewrite)

---

## Objective

Update Nix files, flake configuration, and development tooling to reference the 
new module path.

---

## Scope — What to Change

1. **`flake.nix`**: Update any module path references
2. **`nix/` directory**: Update all Nix files
3. **`tools/` directory**: Update tooling scripts (e.g., `tools/bump`)
4. **OpenAPI configs**: `openapi/` if they reference the module

---

## Scope — What NOT to Change

- Do NOT modify Go source files
- Do NOT modify Dockerfiles or Makefile
- Do NOT modify documentation
- Do NOT run `go run ./cmd/vendorhash update` yet (later phase)
- Do NOT update `flakehashes.json` yet

---

## Implementation Steps

### 1. Change to worktree directory

```bash
cd /home/denny/Project/headscale-project/headscale-module-path-nix
```

### 2. Merge core rewrite

```bash
git merge feature/module-path-core --no-edit
```

### 3. Check flake.nix

```bash
grep "github.com/juanfont/headscale" flake.nix || echo "None found"
```

If found, update:
```bash
sed -i 's|github.com/juanfont/headscale|github.com/arsydoni4326-alt/headscale|g' flake.nix
```

### 4. Check nix/ directory

```bash
grep -r "github.com/juanfont/headscale" nix/ --include="*.nix" || echo "None found"
```

If found, update:
```bash
find nix/ -name "*.nix" -type f -exec sed -i \
  's|github.com/juanfont/headscale|github.com/arsydoni4326-alt/headscale|g' {} +
```

### 5. Check tools/ directory

```bash
grep -r "github.com/juanfont/headscale" tools/ || echo "None found"
```

If found, update:
```bash
find tools/ -type f -exec sed -i \
  's|github.com/juanfont/headscale|github.com/arsydoni4326-alt/headscale|g' {} +
```

### 6. Check openapi/ directory

```bash
grep -r "github.com/juanfont/headscale" openapi/ 2>/dev/null || echo "None found"
```

If found, update:
```bash
find openapi/ -type f -exec sed -i \
  's|github.com/juanfont/headscale|github.com/arsydoni4326-alt/headscale|g' {} +
```

### 7. Check mkdocs.yml

```bash
grep "github.com/juanfont/headscale" mkdocs.yml 2>/dev/null || echo "None found"
```

If found, update:
```bash
sed -i 's|github.com/juanfont/headscale|github.com/arsydoni4326-alt/headscale|g' mkdocs.yml
```

### 8. Verify no old references in Nix/tooling

```bash
grep -r "github.com/juanfont/headscale" \
  flake.nix nix/ tools/ openapi/ mkdocs.yml \
  2>/dev/null || echo "All updated"
```

**Expected:** "All updated" or no matches

---

## Validation Checklist

- [ ] flake.nix updated (if needed)
- [ ] nix/*.nix files updated
- [ ] tools/ scripts updated
- [ ] openapi/ configs updated (if needed)
- [ ] No old references remain

---

## Commit

```bash
git add flake.nix nix/ tools/ openapi/ mkdocs.yml 2>/dev/null || true
git commit -m "nix: update module path references

Update Nix flake, nix/ configuration files, and development tooling
from github.com/juanfont/headscale to github.com/arsydoni4326-alt/headscale.

Files updated:
- flake.nix (if applicable)
- nix/*.nix
- tools/ scripts
- openapi/ configs (if applicable)

Refs: ROADMAP.md Phase 14"
```

### Record commit SHA

```bash
git rev-parse HEAD
```

---

## After Completion

1. **Record commit SHA**: `_________________`
2. **Files updated**: Check `git show --stat`
3. **Any issues**: Document below

---

## Note

Do NOT run `go run ./cmd/vendorhash update` yet. This will be done in the 
final integration phase after all branches are merged.

# Task 4: Documentation and Scripts Update

**Worktree:** `/home/denny/Project/headscale-project/headscale-module-path-docs`  
**Branch:** `feature/module-path-docs`  
**Roadmap Item:** Phase 14 — Module Path Rewrite  
**Sequence:** 2 (parallel with tasks 2, 3, 5)  
**Dependencies:** Task 1 (Core Module Path Rewrite)

---

## Objective

Update all documentation files, CI/CD configs, and non-build scripts to reference 
the new module path.

---

## Scope — What to Change

1. **Documentation**: `README.md`, `CONTRIBUTING.md`, `ARCHITECTURE.md`, `docs/`
2. **CI/CD**: `.github/workflows/*.yml`
3. **Test configs**: `integration/`, `cmd/hi/` configs
4. **Release config**: `.goreleaser.yml`
5. **Package configs**: `packaging/`

---

## Scope — What NOT to Change

- Do NOT modify Go source files
- Do NOT modify Dockerfiles or Makefile
- Do NOT modify Nix files
- Do NOT modify `go.mod` or `go.sum`

---

## Implementation Steps

### 1. Change to worktree directory

```bash
cd /home/denny/Project/headscale-project/headscale-module-path-docs
```

### 2. Merge core rewrite

```bash
git merge feature/module-path-core --no-edit
```

### 3. Update README.md

```bash
sed -i 's|github.com/juanfont/headscale|github.com/arsydoni4326-alt/headscale|g' README.md
```

### 4. Update CONTRIBUTING.md

```bash
sed -i 's|github.com/juanfont/headscale|github.com/arsydoni4326-alt/headscale|g' CONTRIBUTING.md
```

### 5. Update ARCHITECTURE.md if exists

```bash
if [ -f ARCHITECTURE.md ]; then
  sed -i 's|github.com/juanfont/headscale|github.com/arsydoni4326-alt/headscale|g' ARCHITECTURE.md
fi
```

### 6. Update all docs/ markdown files

```bash
find docs/ -name "*.md" -type f -exec sed -i \
  's|github.com/juanfont/headscale|github.com/arsydoni4326-alt/headscale|g' {} +
```

### 7. Update GitHub workflows

```bash
find .github/workflows/ -name "*.yml" -o -name "*.yaml" -type f -exec sed -i \
  's|github.com/juanfont/headscale|github.com/arsydoni4326-alt/headscale|g' {} +
```

### 8. Update .goreleaser.yml

```bash
sed -i 's|github.com/juanfont/headscale|github.com/arsydoni4326-alt/headscale|g' .goreleaser.yml
```

### 9. Check integration configs

```bash
grep -r "github.com/juanfont/headscale" integration/ cmd/hi/ --include="*.md" --include="*.yml" --include="*.yaml" 2>/dev/null || echo "None found"
```

If found, update:
```bash
find integration/ cmd/hi/ -type f \( -name "*.md" -o -name "*.yml" -o -name "*.yaml" \) -exec sed -i \
  's|github.com/juanfont/headscale|github.com/arsydoni4326-alt/headscale|g' {} +
```

### 10. Check packaging/

```bash
grep -r "github.com/juanfont/headscale" packaging/ 2>/dev/null || echo "None found"
```

If found, update:
```bash
find packaging/ -type f -exec sed -i \
  's|github.com/juanfont/headscale|github.com/arsydoni4326-alt/headscale|g' {} +
```

### 11. Verify no old references in docs/configs

```bash
grep -r "github.com/juanfont/headscale" \
  --include="*.md" \
  --include="*.yml" \
  --include="*.yaml" \
  --exclude-dir=.git \
  --exclude-dir=headplane \
  . 2>/dev/null || echo "All updated"
```

---

## Validation Checklist

- [ ] README.md updated
- [ ] CONTRIBUTING.md updated
- [ ] docs/ files updated
- [ ] GitHub workflows updated
- [ ] .goreleaser.yml updated
- [ ] No old references in documentation

---

## Commit

```bash
git add README.md CONTRIBUTING.md ARCHITECTURE.md docs/ .github/ .goreleaser.yml integration/ cmd/hi/ packaging/ 2>/dev/null || true
git commit -m "docs: update module path references

Update all documentation, CI/CD workflows, and configuration files
from github.com/juanfont/headscale to github.com/arsydoni4326-alt/headscale.

Files updated:
- README.md, CONTRIBUTING.md, ARCHITECTURE.md
- docs/*.md
- .github/workflows/*.yml
- integration and test configs
- .goreleaser.yml
- packaging configs

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

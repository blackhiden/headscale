# Parallel Feature Development with Git Worktrees

This document describes the workflow for developing multiple independent features in parallel using Git Worktrees. This approach allows multiple developers or AI agents to work on separate features simultaneously without conflicts.

---

## Overview

Git Worktrees allow you to have multiple working trees (directories) attached to the same repository, each on a different branch. This enables:

- **Parallel development**: Multiple features can be developed simultaneously
- **Isolation**: Each feature has its own working directory and branch
- **Independence**: Features can be built, tested, and reviewed independently
- **Clean merges**: All features originate from the same base commit

---

## Directory Structure

```
parent/
├── headscale/                              # main repository (dev or main branch)
├── headscale-feature-a/                    # worktree for feature/feature-a
├── headscale-feature-b/                    # worktree for feature/feature-b
├── headscale-feature-c/                    # worktree for feature/feature-c
└── headscale-feature-d/                    # worktree for feature/feature-d
```

---

## Workflow

### 1. Planning Phase

Before creating worktrees:

1. **Analyze the roadmap** and identify the current phase
2. **Break down the phase** into logical, independent tasks
3. **Identify parallelizable tasks** (low coupling, different files/modules)
4. **Create an implementation plan** with:
   - Task descriptions
   - Expected files/modules to change
   - Dependencies between tasks
   - Testing requirements
   - Acceptance criteria
5. **Get explicit approval** before proceeding

### 2. Setup Phase

#### a. Prepare the Base Branch

```bash
# Ensure you're on a clean state
git checkout main  # or dev
git pull --ff-only
git status  # should be clean

# Record the base commit SHA
BASE_COMMIT=$(git rev-parse HEAD)
echo "Base commit: $BASE_COMMIT"
```

#### b. Create Feature Branches

```bash
# Create all feature branches from the same base commit
git branch feature/feature-a
git branch feature/feature-b
git branch feature/feature-c
git branch feature/feature-d

# Verify branches
git branch --list 'feature/*'
```

#### c. Create Worktrees

```bash
# Create worktrees in parent directory
git worktree add ../headscale-feature-a feature/feature-a
git worktree add ../headscale-feature-b feature/feature-b
git worktree add ../headscale-feature-c feature/feature-c
git worktree add ../headscale-feature-d feature/feature-d

# Verify worktrees
git worktree list
```

### 3. Implementation Phase

Each developer/agent works in their assigned worktree:

```bash
# Navigate to your worktree
cd ../headscale-feature-a

# Verify you're on the correct branch
git branch --show-current  # should show feature/feature-a

# Implement the feature
# ... make changes, add tests, update docs ...

# Commit your work
git add -A
git commit -m "feat: implement feature A"

# Push to remote (optional, for collaboration)
git push -u origin feature/feature-a
```

**Important rules for implementation:**

- **Work only in your assigned worktree** — never modify another worktree
- **Stay on your feature branch** — do not checkout other branches
- **Follow existing conventions** — match coding style, architecture, patterns
- **Add tests** — unit, integration, E2E as appropriate
- **Update documentation** — keep docs synchronized with changes
- **Keep scope focused** — implement only the assigned task
- **Avoid refactoring** unrelated code unless required for the task

### 4. Testing and Validation

Each feature must be independently tested:

```bash
# Build the project
make build

# Run tests
make test
go test ./...

# Run integration tests (if applicable)
go run ./cmd/hi run "TestFeatureName"

# Format and lint
make fmt
make lint

# Verify the build
./headscale version
```

Document the exact commands used for validation.

### 5. Integration Analysis

After all features are complete:

1. **Inspect all feature branches** for potential conflicts
2. **Identify shared files** modified by multiple branches
3. **Check for**:
   - Database migration conflicts
   - API contract conflicts
   - Configuration conflicts
   - Shared component changes
   - Documentation conflicts
4. **Plan merge order** if dependencies exist
5. **Request approval** before merging

### 6. Merge Phase

Merge features one at a time after approval:

```bash
# Switch to main repository
cd /home/denny/Project/headscale

# Ensure main/dev is up to date
git checkout main  # or dev
git pull --ff-only

# Merge the feature branch
git merge --no-ff feature/feature-a -m "Merge feature A"

# Push to remote
git push origin main

# Repeat for other features
```

### 7. Cleanup Phase

After a feature is successfully merged and verified:

```bash
# Remove the worktree
git worktree remove ../headscale-feature-a

# Delete the feature branch (optional, only if safe)
git branch -d feature/feature-a

# Delete remote branch (optional)
git push origin --delete feature/feature-a
```

**Never delete a worktree with uncommitted changes.**

---

## Best Practices

1. **All worktrees from the same commit** — ensures clean merge base
2. **Maximum 5 parallel tasks** — more increases merge complexity
3. **Choose low-coupling tasks** — avoid tasks that touch the same core code
4. **Document everything** — commands, decisions, assumptions
5. **Test independently** — each feature must work standalone
6. **Commit frequently** — small, focused commits are easier to review
7. **Push to remote** — enables collaboration and backup
8. **Keep worktrees short-lived** — merge and clean up promptly

---

## Additional Resources

- [Git Worktree Documentation](https://git-scm.com/docs/git-worktree)
- [Project Guidelines](../CONTRIBUTING.md)
- [Architecture Documentation](../ARCHITECTURE.md)
- [Roadmap](../ROADMAP.md)
- Automation script: `scripts/parallel-feature-worktrees.sh`



# Parallel Development Workflow with Git Worktrees

This document describes the reusable workflow for executing parallel development tasks using Git worktrees and isolated feature branches.

## Overview

Execute multiple independent tasks simultaneously in isolated Git worktrees, each with its own feature branch, all originating from the same base commit.

## When to Use This Workflow

**Use when:**
- A phase contains multiple independent tasks
- Minimal file overlap between tasks
- Tasks can be validated independently
- 1-10 tasks need to run in parallel
- Each task can complete without waiting for others

**Do NOT use when:**
- Tasks modify the same core files extensively
- Sequential implementation is required for correctness
- Tasks are tightly coupled or have complex dependencies
- Database migrations require coordination
- The merge strategy is not yet clear

## Prerequisites

- Clean working tree (no uncommitted changes)
- Explicit base ref or commit SHA identified
- Task breakdown and dependency analysis complete
- Approval to proceed with parallel implementation

## Workflow Steps

### 1. Planning Phase (Required)

Before creating any worktrees:

1. Read `ROADMAP.md`, `SPECIFICATION.md`, `ARCHITECTURE.md`, and relevant documentation
2. Identify the current phase and break it into logical, independent tasks
3. Analyze file overlap and dependencies between tasks
4. Document each task's objective, scope, dependencies, testing requirements, and commit message format
5. **GET EXPLICIT APPROVAL** before proceeding

### 2. Worktree Creation

Use the provided script for safe, atomic worktree creation:

```bash
# Verify clean state
git status --porcelain  # Must be empty
BASE_REF="main"  # or 'dev', or a specific commit SHA

# Record the base commit
BASE_COMMIT=$(git rev-parse "$BASE_REF")
echo "Base commit: $BASE_COMMIT"

# Create worktrees (requires explicit base ref)
./scripts/parallel-feature-worktrees.sh "$BASE_REF" task1 task2 task3

# Verify creation
git worktree list
```

The script guarantees:
- All branches created from the exact same base commit
- Non-interactive operation (no prompts)
- Collision detection (fails if branch/path exists)
- Atomic rollback on any error

### 3. Task Assignment and Agent Instructions

For each task, create clear instructions with:
- Worktree path and branch name
- Objective and scope (what to change and NOT change)
- Step-by-step implementation guidance
- Validation checklist (build, test, lint commands)
- Exact commit message format
- Dependencies on other tasks (if any)
- Restrictions (files not to touch, worktrees not to access)

### 4. Execution Rules

**For independent tasks:** Start all agents in parallel, each working only in its assigned worktree.

**For dependent tasks:**
1. Primary/foundation task completes first
2. Agent reports commit SHA and affected files
3. Dependent agents merge the foundation commit
4. Dependent agents proceed in parallel

**Isolation rules (mandatory):**
- Work only in assigned worktree
- Never modify other worktrees
- Never merge to main/dev yourself
- Stay strictly within assigned scope

### 5. Integration Analysis

After all agents complete, analyze for conflicts:

```bash
# Compare changed files between two branches
git diff feature/task1..feature/task2 --name-only

# List files changed by each task
git log feature/task1 --name-only --pretty=format: | sort -u > task1.txt
git log feature/task2 --name-only --pretty=format: | sort -u > task2.txt
comm -12 task1.txt task2.txt
```

### 6. Validation (Per-Branch)

Each branch must pass validation independently:

```bash
cd /path/to/worktree
make build
make test
make fmt
make lint
git diff --check
```

### 7. Merge Strategy

Merge in dependency order after approval:

```bash
cd /original/repo/path
git checkout main
git merge --no-ff feature/task1 -m "Merge feature/task1: <description>"
make clean && make build && make test
git push origin main
```

### 8. Cleanup (Only After Approval)

```bash
git worktree remove /path/to/worktree
git branch -d feature/task-name  # Optional
```

## Script Reference

### `scripts/parallel-feature-worktrees.sh`

Non-interactive script for safe worktree creation.

**Usage:**
```bash
./scripts/parallel-feature-worktrees.sh <base-ref> <name1> [name2] ... [name10]
```

**Requirements:**
- Clean working tree
- Explicit base ref (branch name or commit SHA)
- 1-10 feature names
- No existing branches/paths with those names

**Safety guarantees:**
- Atomic operation (all succeed or all rolled back)
- Collision detection before any changes
- All branches from exact same base commit

## Best Practices

1. **Single base commit** — All worktrees from exactly the same commit
2. **Maximum 10 parallel tasks** — More tasks increase merge complexity
3. **Choose low-coupling tasks** — Avoid tasks that touch the same core code
4. **Document everything** — Commands run, decisions made, assumptions
5. **Test independently** — Each feature must work standalone
6. **Commit frequently** — Small, focused commits
7. **Never skip validation** — Build + test + lint before merge
8. **Keep worktrees short-lived** — Merge and clean up promptly

## See Also

- [Git Worktree Documentation](https://git-scm.com/docs/git-worktree)
- [Project Guidelines](../CONTRIBUTING.md)
- [Architecture Documentation](../ARCHITECTURE.md)
- [Roadmap](../ROADMAP.md)
- Detailed reference: `docs/parallel-development.md`

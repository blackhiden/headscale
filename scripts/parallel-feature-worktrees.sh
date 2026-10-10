#!/usr/bin/env bash
# Parallel Feature Worktrees Setup Script
# Creates Git worktrees for parallel feature development
#
# Usage: ./scripts/parallel-feature-worktrees.sh <base-ref> <name1> [name2] ... [name10]
#
# Requirements:
# - Clean working tree (no uncommitted changes)
# - Explicit base ref or SHA (branch name or commit hash)
# - 1-10 feature names
# - No existing branches or paths with the same names
#
# Safety guarantees:
# - Non-interactive (no prompts)
# - All branches created from the exact same base commit
# - Collision detection prevents accidental overwrites
# - Atomic operation: either all succeed or none are created

set -euo pipefail

# Configuration
REPO_DIR="$(pwd)"
PARENT_DIR="$(dirname "$REPO_DIR")"
REPO_NAME="$(basename "$REPO_DIR")"

# Usage
usage() {
    cat <<EOF
Usage: $0 <base-ref> <name1> [name2] ... [name10]

Creates Git worktrees for parallel feature development.

Arguments:
  base-ref    Base branch or commit SHA (e.g., 'main', 'dev', or a commit hash)
  name1-10    Feature names (1-10 names; will create feature/<name> branches)

Example:
  $0 main plugin-system multi-instance monitoring

This creates:
  - Branches: feature/plugin-system, feature/multi-instance, feature/monitoring
  - Worktrees: ../${REPO_NAME}-plugin-system, ../${REPO_NAME}-multi-instance, ../${REPO_NAME}-monitoring
  - All from the exact commit at 'main'

Requirements:
  - Clean working tree (no uncommitted changes)
  - No existing branches or paths with the same names
  - Base ref must resolve to a valid commit

Safety:
  - Non-interactive (no prompts)
  - All branches created from the exact same base commit
  - Collision detection prevents overwrites
  - Fails atomically if any step cannot complete

See docs/parallel-development-workflow.md for the full workflow.
EOF
    exit 1
}

# Check arguments
if [ $# -lt 2 ]; then
    echo "Error: At least a base ref and one feature name are required." >&2
    echo "" >&2
    usage
fi

if [ $# -gt 11 ]; then
    echo "Error: Maximum 10 feature names allowed (got $(($# - 1)))." >&2
    echo "" >&2
    usage
fi

BASE_REF="$1"
shift
FEATURE_NAMES=("$@")

# Verify we're in a git repository
if [ ! -d .git ]; then
    echo "Error: Not in a git repository root. Run this script from the repository root." >&2
    exit 1
fi

# Check for uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo "Error: Working tree has uncommitted changes. Commit or stash them first." >&2
    git status --short >&2
    exit 1
fi

# Resolve base ref to commit SHA
if ! BASE_COMMIT=$(git rev-parse --verify "$BASE_REF^{commit}" 2>/dev/null); then
    echo "Error: Base ref '$BASE_REF' does not resolve to a valid commit." >&2
    exit 1
fi

echo "Base ref: $BASE_REF"
echo "Base commit: $BASE_COMMIT"
echo "Repository: $REPO_NAME"
echo "Features: ${FEATURE_NAMES[*]}"
echo ""

# Pre-flight collision check
COLLISION_DETECTED=0
for name in "${FEATURE_NAMES[@]}"; do
    BRANCH_NAME="feature/$name"
    WORKTREE_PATH="${PARENT_DIR}/${REPO_NAME}-${name}"

    if git show-ref --verify --quiet "refs/heads/$BRANCH_NAME"; then
        echo "Error: Branch '$BRANCH_NAME' already exists." >&2
        COLLISION_DETECTED=1
    fi

    if [ -e "$WORKTREE_PATH" ]; then
        echo "Error: Path '$WORKTREE_PATH' already exists." >&2
        COLLISION_DETECTED=1
    fi
done

if [ $COLLISION_DETECTED -eq 1 ]; then
    echo "" >&2
    echo "Collision detected. No branches or worktrees were created." >&2
    exit 1
fi

# Record branches and worktrees created for potential cleanup
CREATED_BRANCHES=()
CREATED_WORKTREES=()

# Cleanup function for error handling
cleanup_on_error() {
    echo "" >&2
    echo "Error occurred. Rolling back partially created branches and worktrees..." >&2

    for worktree in "${CREATED_WORKTREES[@]}"; do
        if [ -n "$worktree" ]; then
            git worktree remove "$worktree" --force 2>/dev/null || true
            echo "  Removed worktree: $worktree" >&2
        fi
    done

    for branch in "${CREATED_BRANCHES[@]}"; do
        if [ -n "$branch" ]; then
            git branch -D "$branch" 2>/dev/null || true
            echo "  Removed branch: $branch" >&2
        fi
    done

    echo "Rollback complete. No changes were made to the repository." >&2
    exit 1
}

trap cleanup_on_error ERR

# Create branches and worktrees
echo "Creating branches and worktrees from $BASE_COMMIT..."
echo ""

for name in "${FEATURE_NAMES[@]}"; do
    BRANCH_NAME="feature/$name"
    WORKTREE_PATH="${PARENT_DIR}/${REPO_NAME}-${name}"

    echo "Creating: $name"

    # Create branch from base commit
    git branch "$BRANCH_NAME" "$BASE_COMMIT"
    CREATED_BRANCHES+=("$BRANCH_NAME")
    echo "  ✓ Branch: $BRANCH_NAME"

    # Create worktree
    git worktree add "$WORKTREE_PATH" "$BRANCH_NAME" >/dev/null 2>&1
    CREATED_WORKTREES+=("$WORKTREE_PATH")
    echo "  ✓ Worktree: $WORKTREE_PATH"
    echo ""
done

# Summary
echo "========================================"
echo "Worktree setup complete!"
echo "========================================"
echo ""
echo "All branches created from: $BASE_COMMIT"
echo ""
git worktree list
echo ""
echo "Next steps:"
echo "1. Navigate to each worktree: cd $PARENT_DIR/${REPO_NAME}-<name>"
echo "2. Implement the feature in isolation"
echo "3. Test independently: make build && make test"
echo "4. Commit changes: git add -A && git commit"
echo "5. Request merge approval (do not merge yourself)"
echo ""
echo "See docs/parallel-development-workflow.md for the complete workflow."

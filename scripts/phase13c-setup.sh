#!/bin/bash
# Phase 13c Parallel Development Setup Script
# Creates worktrees and branches for multi-user support implementation

set -e

echo "=== Phase 13c Multi-User Support: Parallel Worktree Setup ==="
echo ""

# Configuration
BASE_BRANCH="dev"
BASE_COMMIT="880a8a7d9abb0cfb77fcb20e767aa67cee4cc02a"
PARENT_DIR="../"

# Task definitions
declare -a TASKS=(
    "multiuser-backend-auth:Backend User Model & Auth"
    "multiuser-per-user-settings:Backend Per-User Settings"
    "multiuser-user-mgmt-api:Backend User Management API"
    "multiuser-frontend-login:Frontend Multi-User Login"
    "multiuser-frontend-user-mgmt:Frontend User Management UI"
    "multiuser-frontend-settings:Frontend Per-User Settings"
    "multiuser-docs:Documentation"
    "multiuser-testing:Testing & Validation"
)

# Check we're in the right place
if [[ ! -d ".git" ]]; then
    echo "Error: Must run from repository root"
    exit 1
fi

# Check for clean state
if [[ -n $(git status --porcelain) ]]; then
    echo "Warning: Working directory has uncommitted changes"
    read -p "Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Verify base commit
CURRENT_COMMIT=$(git rev-parse HEAD)
if [[ "$CURRENT_COMMIT" != "$BASE_COMMIT" ]]; then
    echo "Warning: Current commit ($CURRENT_COMMIT) != expected base ($BASE_COMMIT)"
    read -p "Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

echo "Base branch: $BASE_BRANCH"
echo "Base commit: $BASE_COMMIT"
echo "Creating ${#TASKS[@]} worktrees..."
echo ""

# Create branches and worktrees
for task in "${TASKS[@]}"; do
    IFS=':' read -r name description <<< "$task"
    branch="feature/$name"
    worktree="${PARENT_DIR}headscale-${name}"
    
    echo "Creating: $branch -> $worktree"
    
    # Check if branch already exists
    if git show-ref --verify --quiet "refs/heads/$branch"; then
        echo "  Branch $branch already exists, skipping..."
    else
        git branch "$branch"
        echo "  Branch created"
    fi
    
    # Check if worktree already exists
    if [[ -d "$worktree" ]]; then
        echo "  Worktree $worktree already exists, skipping..."
    else
        git worktree add "$worktree" "$branch"
        echo "  Worktree created"
    fi
    
    echo ""
done

echo "=== Setup Complete ==="
echo ""
echo "Worktrees created:"
git worktree list | grep multiuser
echo ""
echo "Next steps:"
echo "1. Assign one agent per worktree"
echo "2. Each agent reads AGENT_INSTRUCTIONS.md in their worktree"
echo "3. Task 1 (backend-auth) must complete and merge first"
echo "4. Other tasks can proceed in parallel after Task 1"
echo ""
echo "See docs/parallel-development.md for workflow details"

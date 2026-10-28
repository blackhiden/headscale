#!/bin/bash
# Merge and cleanup worktrees script

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

BASE_BRANCH="dev"

echo -e "${BLUE}=== Worktree Merge Tool ===${NC}"
echo ""

# Get worktrees
declare -a BRANCHES
declare -a PATHS
index=0

while IFS= read -r line; do
    path=$(echo "$line" | awk '{print $1}')
    branch=$(echo "$line" | grep -oP '\[\K[^\]]+' || echo "")
    
    if [[ $branch == feature/* ]]; then
        BRANCHES[$index]="$branch"
        PATHS[$index]="$path"
        ((index++))
    fi
done < <(git worktree list)

if [ $index -eq 0 ]; then
    echo -e "${YELLOW}No feature branches found.${NC}"
    exit 0
fi

echo -e "${GREEN}Available Feature Branches:${NC}"
echo "============================================================"
printf "%-4s %-35s %s\n" "No." "Branch" "Path"
echo "============================================================"

for i in "${!BRANCHES[@]}"; do
    printf "%-4s %-35s %s\n" "$((i+1))" "${BRANCHES[$i]}" "$(basename "${PATHS[$i]}")"
done

echo "============================================================"
echo ""
echo -e "${YELLOW}Options:${NC}"
echo "  1-$index : Merge specific feature"
echo "  a       : Merge all features"
echo "  q       : Quit"
echo ""
read -p "Your choice: " choice

if [[ "$choice" == "q" ]] || [[ "$choice" == "Q" ]]; then
    echo "Cancelled."
    exit 0
fi

merge_feature() {
    local branch="$1"
    local path="$2"
    
    echo ""
    echo -e "${BLUE}>>> Processing: ${GREEN}$branch${NC}"
    
    echo "  1. Checkout $BASE_BRANCH"
    git checkout "$BASE_BRANCH" || {
        echo -e "${RED}  X Failed to checkout $BASE_BRANCH${NC}"
        return 1
    }
    
    echo "  2. Merge $branch"
    git merge --no-ff "$branch" -m "Merge branch '$branch' into $BASE_BRANCH" || {
        echo -e "${RED}  X Merge failed${NC}"
        return 1
    }
    echo -e "${GREEN}     √ Merged${NC}"
    
    echo "  3. Remove worktree"
    git worktree remove --force "$path" || {
        echo -e "${YELLOW}     ! Manual: rm -rf $path${NC}"
    }
    echo -e "${GREEN}     √ Removed${NC}"
    
    echo "  4. Delete branch"
    git branch -D "$branch" || echo -e "${YELLOW}     ! Not deleted${NC}"
    echo -e "${GREEN}     √ Done${NC}"
}

if [[ "$choice" == "a" ]] || [[ "$choice" == "A" ]]; then
    for i in "${!BRANCHES[@]}"; do
        merge_feature "${BRANCHES[$i]}" "${PATHS[$i]}"
    done
elif [[ "$choice" =~ ^[0-9]+$ ]] && [ "$choice" -ge 1 ] && [ "$choice" -le "$index" ]; then
    idx=$((choice-1))
    merge_feature "${BRANCHES[$idx]}" "${PATHS[$idx]}"
else
    echo -e "${RED}Invalid choice${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}============================================================${NC}"
echo -e "${GREEN}√ Complete!${NC}"
echo ""
echo "Current branch: $(git branch --show-current)"
echo ""
echo "Remaining worktrees:"
git worktree list

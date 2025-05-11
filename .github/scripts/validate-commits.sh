#!/bin/bash
set -e

# Improve debugging
echo "Current directory: $(pwd)"
echo "Git version: $(git --version)"

# Determine the commit range to check
if [ -n "$GITHUB_BASE_REF" ]; then
  # GitHub Actions PR
  echo "Running in GitHub Actions PR mode"
  git fetch origin $GITHUB_BASE_REF --depth=1000
  COMMIT_RANGE="origin/$GITHUB_BASE_REF..HEAD"
elif [ -n "$GITHUB_EVENT_NAME" ] && [ "$GITHUB_EVENT_NAME" = "push" ]; then
  # GitHub Actions push
  echo "Running in GitHub Actions push mode"
  COMMIT_RANGE="$GITHUB_BEFORE..$GITHUB_SHA"
elif [ -n "$CI_MERGE_REQUEST_TARGET_BRANCH_NAME" ]; then
  # GitLab CI
  echo "Running in GitLab CI mode"
  git fetch origin $CI_MERGE_REQUEST_TARGET_BRANCH_NAME --depth=1000
  COMMIT_RANGE="origin/$CI_MERGE_REQUEST_TARGET_BRANCH_NAME..HEAD"
else
  # Local or other CI - try to determine from git directly
  echo "Running in generic mode, attempting to determine commits"
  # Check for merge base with main/master
  if git rev-parse --verify origin/main >/dev/null 2>&1; then
    BASE_BRANCH="origin/main"
  elif git rev-parse --verify origin/master >/dev/null 2>&1; then
    BASE_BRANCH="origin/master"
  else
    echo "Could not determine base branch, defaulting to last 10 commits"
    COMMIT_RANGE="HEAD~10..HEAD"
  fi
  
  if [ -z "$COMMIT_RANGE" ]; then
    MERGE_BASE=$(git merge-base HEAD $BASE_BRANCH)
    COMMIT_RANGE="$MERGE_BASE..HEAD"
  fi
fi

echo "Checking commits in range: $COMMIT_RANGE"
ERROR=0

# Get all commit SHAs in the range
COMMIT_LIST=$(git rev-list $COMMIT_RANGE)
if [ -z "$COMMIT_LIST" ]; then
  echo "No commits found in range: $COMMIT_RANGE"
  exit 0
fi

for commit_hash in $COMMIT_LIST; do
  echo "Checking commit: $commit_hash"
  
  # Get the commit message
  commit_msg=$(git log --format=%B -n 1 $commit_hash)
  
  # Skip GitHub's automatic merge commits for PRs
  if [[ "$commit_msg" =~ ^Merge\ [0-9a-f]+\ into\ [0-9a-f]+ ]]; then
    echo "⏩ Skipping GitHub automatic merge commit: $commit_hash"
    continue
  fi
  
  # Also skip regular merge commits if needed
  if [[ "$commit_msg" =~ ^Merge\ (pull\ request|branch) ]]; then
    echo "⏩ Skipping merge commit: $commit_hash"
    continue
  fi
  
  # Get the first line of the commit message (subject line)
  subject_line=$(echo "$commit_msg" | head -n 1)
  
  # Display the subject for debugging
  echo "Subject: $subject_line"
  
  # Regular expression for Conventional Commits pattern
  pattern='^(build|chore|ci|docs|feat|fix|perf|refactor|revert|style|test)(\([a-z0-9 -]+\))?: .+'
  
  if ! [[ "$subject_line" =~ $pattern ]]; then
    echo "❌ ERROR: Commit $commit_hash does not follow the Conventional Commits standard."
    echo "Example: feat(auth): add login functionality"
    echo "Types: build, chore, ci, docs, feat, fix, perf, refactor, revert, style, test"
    ERROR=1
    continue
  fi
  
  # Extract the part after the colon and space
  message_part=$(echo "$subject_line" | sed -E 's/^[^:]+: (.*)$/\1/')
  
  # Check first letter capitalization
  first_letter=$(echo "$message_part" | cut -c1)
  if [[ "$first_letter" =~ [A-Z] ]]; then
    echo "❌ ERROR: Commit $commit_hash has a capitalized first letter after the colon."
    ERROR=1
    continue
  fi
  
  # Check for period at the end
  if [[ "$subject_line" =~ \.$  ]]; then
    echo "❌ ERROR: Commit $commit_hash has a period at the end of the subject line."
    ERROR=1
    continue
  fi
  
  echo "✅ Commit $commit_hash is valid"
done

if [ $ERROR -ne 0 ]; then
  echo "❌ One or more commits have invalid messages. See above for details."
  exit 1
else
  echo "✅ All commits are valid!"
  exit 0
fi

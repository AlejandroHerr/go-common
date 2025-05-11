#!/bin/bash

# Get all commits in the range (from base branch to HEAD)
if [ "$GITHUB_BASE_REF" != "" ]; then
  # For pull requests
  COMMIT_RANGE="origin/${GITHUB_BASE_REF}...HEAD"
else
  # For direct pushes (you can adjust the depth as needed)
  COMMIT_RANGE="HEAD~10..HEAD"
fi

echo "Checking commits in range: $COMMIT_RANGE"
ERROR=0

validate_commit_message() {
  local commit_msg="$1"
  # Regular expression for Conventional Commits pattern
  local pattern='^(build|chore|ci|docs|feat|fix|perf|refactor|revert|style|test)(\([a-z0-9 -]+\))?: .+'
  
  if ! [[ "$commit_msg" =~ $pattern ]]; then
    echo "❌ Commit $(echo $commit_msg | cut -c1-10) does not follow the Conventional Commits standard."
    ERROR=1
    return
  fi
  
  # Ensure first letter after colon is not capitalized
  local first_letter=$(echo "$commit_msg" | sed -n 's/^[^:]*: \(.\).*/\1/p')
  if [[ "$first_letter" =~ [A-Z] ]]; then
    echo "❌ Commit $(echo $commit_msg | cut -c1-10) has capitalized first letter after the colon."
    ERROR=1
    return
  fi
  
  # Ensure no period at the end of subject line
  local subject_line=$(echo "$commit_msg" | head -n 1)
  if [[ "$subject_line" =~ \.$  ]]; then
    echo "❌ Commit $(echo $commit_msg | cut -c1-10) has a period at the end of the subject line."
    ERROR=1
    return
  fi
  
  echo "✅ Commit $(echo $commit_msg | cut -c1-10) is valid"
}

# Get all commit messages in the range
for commit_hash in $(git rev-list $COMMIT_RANGE); do
  commit_msg=$(git log --format=%B -n 1 $commit_hash)
  validate_commit_message "$commit_msg"
done

if [ $ERROR -ne 0 ]; then
  echo "❌ One or more commits have invalid messages. See above for details."
  echo "Valid format: type(scope): message"
  echo "Types: build, chore, ci, docs, feat, fix, perf, refactor, revert, style, test"
  exit 1
else
  echo "✅ All commits are valid!"
  exit 0
fi

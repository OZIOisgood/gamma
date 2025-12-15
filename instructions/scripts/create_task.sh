#!/bin/bash

# Check if a task name was provided
if [ -z "$1" ]; then
  echo "Usage: $0 \"task name\""
  exit 1
fi

# Get current timestamp
TIMESTAMP=$(date +"%Y%m%d%H%M%S")

# Convert task name to snake_case
TASK_NAME=$(echo "$1" | tr '[:upper:]' '[:lower:]' | tr ' ' '_')

# Create directory name
DIR_NAME="${TIMESTAMP}_${TASK_NAME}"
TASK_DIR="instructions/tasks/${DIR_NAME}"

# Create directory
mkdir -p "$TASK_DIR"

# Create README.md (Task Definition)
cat <<EOL > "$TASK_DIR/README.md"
# Task: $1

## Status
- [ ] Defined
- [ ] In Progress
- [ ] Completed

## Description
<!-- Describe the task, user story, or bug report here -->

## Context
<!-- Relevant files, links, or background info -->

## Acceptance Criteria
- [ ] 
EOL

# Create RESOLUTION.md
cat <<EOL > "$TASK_DIR/RESOLUTION.md"
# Resolution

## Summary
<!-- What was done? -->

## Technical Details
<!-- Key changes, decisions made -->

## Verification
- [ ] Build passed
- [ ] Verified in UI/API
EOL

echo "Created task: $TASK_DIR"
echo "Files created:"
echo "  - $TASK_DIR/README.md"
echo "  - $TASK_DIR/RESOLUTION.md"

#!/bin/bash

# Check if an issue name was provided
if [ -z "$1" ]; then
  echo "Usage: $0 \"issue name\""
  exit 1
fi

# Get current timestamp
TIMESTAMP=$(date +"%Y%m%d%H%M%S")
DATE_PRETTY=$(date +"%Y-%m-%d %H:%M:%S")

# Convert issue name to snake_case
ISSUE_NAME=$(echo "$1" | tr '[:upper:]' '[:lower:]' | tr ' ' '_')

# Create file path
FILE_PATH="instructions/common_issues/${TIMESTAMP}_${ISSUE_NAME}.md"

# Create Issue File
cat <<EOL > "$FILE_PATH"
# $1

**Date:** $DATE_PRETTY
**Status:** Resolved

## Description
<!-- Describe the issue, error message, or unexpected behavior -->

## Solution
<!-- Describe how it was fixed. Include code snippets if relevant. -->

### Code Change
<!-- Optional: Show the diff or the code that fixed it -->

## Prevention / Notes
<!-- How to avoid this in the future -->
EOL

echo "Created common issue file: $FILE_PATH"

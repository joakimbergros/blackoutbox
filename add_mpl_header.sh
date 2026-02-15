#!/usr/bin/env bash
set -euo pipefail

# Function to add header to a single file
add_header() {
    local file="$1"
    local comment_char="$2"

    # Check if header already exists
    if head -3 "$file" | grep -q "Mozilla Public"; then
        return 0
    fi

    # Create temp file
    tmpfile=$(mktemp)

    # Write header
    {
        echo "${comment_char} This Source Code Form is subject to the terms of the Mozilla Public"
        echo "${comment_char} License, v. 2.0. If a copy of the MPL was not distributed with this"
        echo "${comment_char} file, You can obtain one at https://mozilla.org/MPL/2.0/."
        echo
    } > "$tmpfile"

    # Append original file content
    cat "$file" >> "$tmpfile"

    # Replace original file
    mv "$tmpfile" "$file"
}

export -f add_header

# Process Go files
find . -type f -name "*.go" -exec bash -c 'add_header "$0" "//"' {} \;

# Process SQL files
find . -type f -name "*.sql" -exec bash -c 'add_header "$0" "--"' {} \;

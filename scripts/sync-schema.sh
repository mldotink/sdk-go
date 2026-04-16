#!/usr/bin/env bash
# Sync schema.graphql from backend source files.
# Usage: ./scripts/sync-schema.sh [path-to-backend-graph-dir]
#
# Defaults to ../backend/go-backend/internal/graph relative to the SDK root.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SDK_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

GRAPH_DIR="${1:-$SDK_ROOT/../backend/go-backend/internal/graph}"

if [[ ! -d "$GRAPH_DIR" ]]; then
  echo "error: graph directory not found: $GRAPH_DIR" >&2
  echo "usage: $0 [path-to-backend/internal/graph]" >&2
  exit 1
fi

echo "Syncing schema from $GRAPH_DIR ..."

# Write the filter script to a temp file (can't use heredoc + pipe to same process stdin).
TMPPY=$(mktemp /tmp/schema_filter_XXXXXX.py)
trap 'rm -f "$TMPPY"' EXIT

cat > "$TMPPY" <<'PYEOF'
import re, sys

s = sys.stdin.read()

# ── Step 1: Remove complete directive definitions (handles multi-line) ────────
# Track parenthesis depth so we skip the full arg list before "on LOCATION".
lines = s.split('\n')
result = []
in_directive = False
paren_depth = 0

for line in lines:
    stripped = line.strip()
    if not in_directive:
        if stripped.startswith('directive @'):
            in_directive = True
            paren_depth = stripped.count('(') - stripped.count(')')
            # Single-line directive already has "on"
            if paren_depth <= 0 and ' on ' in stripped:
                in_directive = False
            continue  # skip directive line
        result.append(line)
    else:
        paren_depth += line.count('(') - line.count(')')
        if paren_depth <= 0 and ' on ' in line:
            in_directive = False
        continue  # skip inside / end of directive

s = '\n'.join(result)

# ── Step 2: Remove inline field-level directive annotations ───────────────────
s = re.sub(r'\s*@(goField|isAuthenticated|agent|hasRole)\b(\([^)]*\))?', '', s)

# ── Step 3: Tidy up extra blank lines ─────────────────────────────────────────
s = re.sub(r'\n{3,}', '\n\n', s)

print(s.strip())
PYEOF

cat "$GRAPH_DIR"/*.graphqls | python3 "$TMPPY" > "$SDK_ROOT/schema.graphql"

echo "Written to $SDK_ROOT/schema.graphql"
echo "Lines: $(wc -l < "$SDK_ROOT/schema.graphql")"

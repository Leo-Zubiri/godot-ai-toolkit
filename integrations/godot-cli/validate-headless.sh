#!/usr/bin/env bash
# validate-headless.sh — headless Godot validation for CI / local checks.
#
# Runs `--headless --import` on a Godot project and reports failure reliably.
# IMPORTANT (verified): headless Godot can exit 0 even when a script fails to
# parse or a resource fails to import. We therefore ALSO scan output for error
# markers and fail on them, instead of trusting the exit code alone.
# See: https://docs.godotengine.org/en/stable/tutorials/editor/command_line_tutorial.html
#
# Usage:
#   validate-headless.sh <project-dir> [godot-binary]
#
# Env:
#   GODOT_BIN   Path to the Godot executable (overrides auto-detection / arg 2).
#
# Exit codes:
#   0  success
#   2  Godot binary not found      (reported as SKIPPED — not a pass)
#   3  usage error
#   1  validation failed (nonzero exit OR error markers in output)

set -uo pipefail

PROJECT_DIR="${1:-}"
GODOT_BIN="${GODOT_BIN:-${2:-}}"

if [ -z "$PROJECT_DIR" ]; then
  echo "usage: validate-headless.sh <project-dir> [godot-binary]" >&2
  exit 3
fi
if [ ! -f "$PROJECT_DIR/project.godot" ]; then
  echo "error: '$PROJECT_DIR' does not contain project.godot" >&2
  exit 3
fi

# Resolve a Godot binary if not provided.
if [ -z "$GODOT_BIN" ]; then
  for cand in godot godot4 Godot; do
    if command -v "$cand" >/dev/null 2>&1; then GODOT_BIN="$cand"; break; fi
  done
fi
if [ -z "$GODOT_BIN" ] || ! command -v "$GODOT_BIN" >/dev/null 2>&1; then
  echo "SKIPPED: no Godot binary found (set GODOT_BIN or pass it as arg 2)." >&2
  echo "         This is a skip, NOT a pass." >&2
  exit 2
fi

echo "Validating '$PROJECT_DIR' with '$GODOT_BIN' (headless import)..."
OUTPUT="$("$GODOT_BIN" --headless --path "$PROJECT_DIR" --import 2>&1)"
STATUS=$?
echo "$OUTPUT"

# Scan for error markers regardless of exit status.
if echo "$OUTPUT" | grep -Eiq 'SCRIPT ERROR|Parse Error|ERROR:|Failed to load'; then
  echo "FAILED: error markers found in Godot output." >&2
  exit 1
fi
if [ "$STATUS" -ne 0 ]; then
  echo "FAILED: Godot exited with status $STATUS." >&2
  exit 1
fi

echo "OK: no errors detected."
exit 0

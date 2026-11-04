#!/usr/bin/env bash
# parse-gdscript.sh — parse-check a single GDScript file with headless Godot.
#
# Mirrors validate-headless.sh but targets one script instead of a whole project.
# As with import, a parse check can pass yet say nothing about runtime behavior —
# this verifies SYNTAX ONLY. Output is also scanned for error markers because
# headless Godot can exit 0 on failure.
# Docs: https://docs.godotengine.org/en/stable/tutorials/editor/command_line_tutorial.html
#
# Usage:
#   parse-gdscript.sh <script.gd> [godot-binary]
# Env:
#   GODOT_BIN   Path to the Godot executable (overrides auto-detection / arg 2).
#
# Exit codes:
#   0  parsed OK (syntax only — NOT a runtime guarantee)
#   1  parse failed (nonzero exit OR error markers in output)
#   2  Godot binary not found      (reported as SKIPPED — not a pass)
#   3  usage error

set -uo pipefail

SCRIPT="${1:-}"
GODOT_BIN="${GODOT_BIN:-${2:-}}"

if [ -z "$SCRIPT" ]; then
  echo "usage: parse-gdscript.sh <script.gd> [godot-binary]" >&2
  exit 3
fi
if [ ! -f "$SCRIPT" ]; then
  echo "error: '$SCRIPT' not found" >&2
  exit 3
fi

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

echo "Parse-checking '$SCRIPT' with '$GODOT_BIN' (syntax only)..."
OUTPUT="$("$GODOT_BIN" --headless --check-only --script "$SCRIPT" 2>&1)"
STATUS=$?
echo "$OUTPUT"

if echo "$OUTPUT" | grep -Eiq 'SCRIPT ERROR|Parse Error|ERROR:'; then
  echo "FAILED: error markers found in Godot output." >&2
  exit 1
fi
if [ "$STATUS" -ne 0 ]; then
  echo "FAILED: Godot exited with status $STATUS." >&2
  exit 1
fi

echo "OK: parsed (syntax only; not a runtime guarantee)."
exit 0

# Godot CLI helpers

Dependency-free helpers for validating Godot projects from the command line (CI and local).

## `validate-headless.sh`

Runs `godot --headless --path <project> --import` and reports failure **reliably**:

```sh
integrations/godot-cli/validate-headless.sh tests/fixtures/gdscript-project
GODOT_BIN=/opt/godot/godot integrations/godot-cli/validate-headless.sh tests/fixtures/mixed-project
```

Why it exists: **headless Godot can exit `0` even on failure** (verified — e.g. a scene whose script
fails to parse). The script therefore scans output for `SCRIPT ERROR` / `Parse Error` / `ERROR:` /
`Failed to load` and fails on them, not just on the exit code.

Exit codes: `0` ok · `1` validation failed · `2` **skipped** (no Godot binary — not a pass) · `3` usage.

## Notes

- For running test suites, prefer each framework's own runner (GUT `-gexit`, GdUnit4 runner/Action),
  which exits non-zero on test failure. See [`../../docs/compatibility.md`](../../docs/compatibility.md).
- Useful flags: `--check-only --script <s>` (parse-check a single script), `--export-release <preset>
  <path>` (export; `--headless` required on machines without a GPU).

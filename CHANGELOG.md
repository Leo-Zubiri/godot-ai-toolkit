# Changelog

All notable changes to `godot-ai-toolkit` are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

The toolkit version is what a consumer pins in `.godot-ai/toolkit.lock.json`. SemVer here
describes the **toolkit contract** (CLI behavior, lock-file schema, adapter schema, skill
format expectations), not any single game.

## [Unreleased]

### Added
- Repository skeleton and architecture (Phase B initial structure).
- Rector document (`README.md`) and operational agent rules (`AGENTS.md`).
- Architecture, compatibility matrix, skill-authoring and distribution docs under `docs/`.
- Agent adapters as data descriptors (`adapters/claude.json`, `opencode.json`, `codex.json`).
- Go CLI (`cmd/godot-ai-toolkit`) with `install`, `update`, `status`, `validate`:
  - Agent Skills (`SKILL.md`) structure + frontmatter validation.
  - Copy-based installation into per-agent discovery paths.
  - `.godot-ai/toolkit.lock.json` generation with per-file SHA-256.
  - Local-modification/conflict detection, timestamped backups, explicit updates.
  - Generic release-tarball fetch (HTTPS + checksum) wired behind `--url`/`--sha256`.
- Test fixtures (valid + invalid skills; gdscript / csharp / mixed Godot projects).
- Headless Godot validation helpers (`integrations/godot-cli`) and .NET notes (`integrations/dotnet`).
- CI workflow building and testing the CLI on Windows, Linux and macOS.
- Research doc `docs/research/existing-skills-ecosystem.md`: comparative study of four existing Godot AI
  skills repos, with an evidence-based taxonomy/quality/source/version/testing policy and a prioritized roadmap.
- P0/P1 design proposals under `docs/proposals/` (authoring & quality standard; verification workflow).
- Skill manifest (`internal/manifest`, `skills/<name>/manifest.json`): schema for compatibility
  (minimum/tested/reference), languages, categorized examples and sources; validated by `validate`.
- Deterministic catalog (`internal/catalog`) and `toolkit catalog` subcommand (`--write` / `--check` staleness);
  initial empty `skills/catalog.json`.
- `toolkit verify` subcommand (`internal/verify`): engine-level verification of skill example files —
  GDScript parse (`godot --check-only`) and C# build (`dotnet build`, reported BUILD_VERIFIED, never runtime);
  PASS/FAIL/SKIPPED/NOT_APPLICABLE with per-example and global aggregation (a skipped mandatory check never
  yields a global PASS); text and JSON reports.
- Per-file GDScript parse helper `integrations/godot-cli/parse-gdscript.sh`.
- CI: a Linux engine tier provisioning Godot 4.7.2 (+ .NET) to run `verify`, plus a `catalog --check` gate.

### Changed
- **Frontmatter contract finalized to `name` + `description` only.** The validator now rejects any other
  top-level key; compatibility/metadata move to the skill manifest. `docs/authoring-skills.md` updated and the
  `example-skill` fixture migrated accordingly.

### Fixed
- `adapter.Validate` now rejects absolute/rooted `skillsPath` cross-platform (a POSIX-style `/abs` slipped
  through `filepath.IsAbs` on Windows, failing CI); added a Windows drive-letter test case.

### Notes
- No official Godot knowledge skills are authored yet. `skills/` holds only the taxonomy
  (`skills/README.md`); the example skill used to exercise the pipeline lives in
  `tests/fixtures/`, never in `skills/`.

<!--
Release template:

## [0.1.0] - YYYY-MM-DD
### Added / Changed / Deprecated / Removed / Fixed / Security
-->

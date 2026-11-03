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

### Notes
- No official Godot knowledge skills are authored yet. `skills/` holds only the taxonomy
  (`skills/README.md`); the example skill used to exercise the pipeline lives in
  `tests/fixtures/`, never in `skills/`.

<!--
Release template:

## [0.1.0] - YYYY-MM-DD
### Added / Changed / Deprecated / Removed / Fixed / Security
-->

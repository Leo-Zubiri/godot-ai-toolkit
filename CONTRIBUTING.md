# Contributing to godot-ai-toolkit

Thank you for helping build a durable, centralized knowledge base for Godot + AI development.
Before contributing, read the [`README.md`](README.md) (the project's rector document) and
[`AGENTS.md`](AGENTS.md). This file covers the *mechanics*; the README owns the *vision and
non-negotiable principles*. Do not duplicate README content here.

## Golden rules

1. **Respect the non-negotiable principles** in the README. If a change conflicts with one,
   say so explicitly and justify why it is necessary — do not work around it silently.
2. **Verify Godot / .NET facts against official documentation** before writing them. Never
   invent APIs, methods, properties, version requirements or platform support. Cite the source.
3. **Single source of truth.** Skills are authored once under `skills/`. Never duplicate skill
   content per agent — per-agent placement is the installer's job via `adapters/`.
4. **Do not over-engineer.** No speculative features, no empty directories, no dependencies that
   the Go standard library can replace.

## Repository layout

See [`docs/architecture.md`](docs/architecture.md) for the full model. In short:

| Path | Responsibility |
|---|---|
| `skills/` | Official Godot knowledge catalog (Agent Skills `SKILL.md`). **Single source of truth.** |
| `adapters/` | Per-agent **data** descriptors (discovery paths, instruction file). No content, no code. |
| `cmd/`, `internal/` | The Go CLI (`install`/`update`/`status`/`validate`). Tooling only. |
| `integrations/` | Optional capabilities: MCP notes, Godot CLI helpers, .NET helpers. |
| `tests/fixtures/` | Deterministic fixtures to exercise the pipeline. The example skill lives **here**, never in `skills/`. |
| `docs/` | Architecture, compatibility matrix, skill authoring, distribution. |

## Adding a skill

> **First** complete the [skill research & technical-grounding protocol](docs/skill-research-protocol.md)
> and present its research report. Then follow [`docs/authoring-skills.md`](docs/authoring-skills.md).
> Authored examples: [`signals`](skills/signals/SKILL.md), [`scenes-and-nodes`](skills/scenes-and-nodes/SKILL.md).
> Summary:

1. Create `skills/<skill-name>/SKILL.md`. The folder name **must equal** the `name` frontmatter
   field (lowercase `[a-z0-9-]`, 1–64 chars).
2. Frontmatter is **`name` + `description` only** (no other key is allowed). Compatibility, languages,
   examples and sources go in the skill's `manifest.json`.
3. Keep `SKILL.md` under ~500 lines (a warning, not a hard cap). Put large material in `references/`.
4. State version/build requirements (e.g. "C# requires the .NET editor build; no web export")
   in `compatibility` and in the body where relevant.
5. Validate before committing (see below).

## Working on the Go CLI

- **Standard library first.** Add a third-party module only when stdlib genuinely cannot do the
  job, and justify it in the PR. The current CLI has **no third-party dependencies**.
- Keep **agent adapters decoupled** from installer core — adapters are loaded as data; never
  hard-code agent names in the installer.
- Module path in `go.mod` (`github.com/godot-ai-toolkit/godot-ai-toolkit`) is a placeholder —
  set it to the real repository path when the remote is created.

## Validating locally

```sh
# Build (requires a Go toolchain; end users do NOT need Go — they get a binary).
go build ./...

# Run all Go tests (validator, installer, lock, fetch).
go test ./...

# Validate the repository's adapters and any skills.
go run ./cmd/godot-ai-toolkit validate
```

Environment-gated checks (run when the tool is available; otherwise reported as skipped — never
as passed):

```sh
# C# fixture build (needs .NET 8+ SDK):
dotnet build tests/fixtures/csharp-project

# Headless Godot parse/import (needs a Godot binary; scan stdout for errors — see the helper):
integrations/godot-cli/validate-headless.sh tests/fixtures/gdscript-project
```

## Commits & pull requests

- Use clear, imperative commit subjects. Reference the principle or deliverable when relevant.
- Update `CHANGELOG.md` under `[Unreleased]`.
- CI (`.github/workflows/validate.yml`) must pass: it builds and tests the CLI on Windows, Linux
  and macOS and runs the always-on validation tier.
- Keep PRs focused. Do not mix a new skill with CLI changes unless they are inseparable.

## Releases (maintainers)

See [`docs/distribution.md`](docs/distribution.md). Summary: bump `CHANGELOG.md`, tag `vX.Y.Z`
(SemVer), build per-OS/arch binaries, publish a content tarball with a `.sha256` checksum.

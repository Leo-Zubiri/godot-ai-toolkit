# Distribution & updates

How the toolkit is versioned, delivered, installed into a game, and updated — always on the
developer's explicit decision (README principle 11).

## Versioning

- **SemVer** (`vMAJOR.MINOR.PATCH`) describes the **toolkit contract**: CLI behavior, the lock-file
  schema, the adapter schema, and skill-format expectations.
- Each release is a **Git tag** + a GitHub **release**, with a `CHANGELOG.md` entry.
- A release publishes: per-OS/arch **CLI binaries**, and a **content tarball** (`skills/`, `adapters/`)
  with a `.sha256` checksum for integrity verification.

## What the end user needs

**Nothing but the binary.** No Go, Python, Node or uv. The CLI is a single native executable per
OS/arch. Git is optional — the install/update happy path fetches the content tarball over HTTPS and
verifies its checksum.

Install methods (in order of simplicity):

1. Download the binary for your OS/arch from the GitHub release; verify its `.sha256`.
2. `curl | sh` / PowerShell bootstrap script (planned).
3. Homebrew tap / Scoop manifest (planned).
4. `go install github.com/godot-ai-toolkit/godot-ai-toolkit/cmd/godot-ai-toolkit@vX.Y.Z` (for Go users).

> The module path `github.com/godot-ai-toolkit/godot-ai-toolkit` is a placeholder until the real
> repository exists; update `go.mod` accordingly.

## The consumer control file: `.godot-ai/toolkit.lock.json`

Written into each game project. It makes the installation **reproducible, inspectable and
independent**. Schema (v1):

```json
{
  "schemaVersion": 1,
  "toolkit": {
    "repository": "https://github.com/godot-ai-toolkit/godot-ai-toolkit",
    "version": "v0.1.0",
    "commit": "<git commit sha, if known>"
  },
  "agents": ["claude", "codex"],
  "skills": [
    {
      "name": "example-skill",
      "source": "skills/example-skill",
      "files": [
        { "path": ".claude/skills/example-skill/SKILL.md", "sha256": "…" },
        { "path": ".codex/skills/example-skill/SKILL.md", "sha256": "…" }
      ]
    }
  ],
  "installedAt": "2026-11-02T00:00:00Z",
  "generatedBy": "godot-ai-toolkit v0.1.0"
}
```

- `toolkit.version` + `toolkit.commit` — the exact pinned source.
- `agents` — which agents were configured.
- `skills[].files[].sha256` — per managed file, so the CLI can detect local edits and know exactly
  what it placed (and what to remove on update).

Commit `toolkit.lock.json` to the game's repository. Do **not** commit `.godot-ai/backups/`.

## Commands

```sh
# Install (first time): pin a version, choose agents and skills.
godot-ai-toolkit install --tag v0.1.0 --agents claude,codex --skills example-skill
#   or install from a local toolkit checkout (used in development/CI):
godot-ai-toolkit install --source /path/to/godot-ai-toolkit --agents claude --skills example-skill

# Update (later, deliberately): change the pinned version.
godot-ai-toolkit update --tag v0.2.0

# Inspect what is installed and whether managed files were changed locally.
godot-ai-toolkit status

# Validate skills/adapters (in the toolkit repo) or a target skills directory.
godot-ai-toolkit validate
```

Common flags: `--project <dir>` (consumer project, default `.`), `--force` (override conflict
protection), `--url <tarball> --sha256 <sum>` (explicit fetch source).

## Safety model

- **Copy, not symlink**, by default — each project is self-contained. (Symlink is a local-dev-only
  option, never the default.)
- **Conflict detection:** before an update overwrites a managed file, the CLI compares the on-disk
  SHA-256 with the lock. If you edited it, the update **stops and lists the conflicts**; proceed only
  with `--force`.
- **Backups:** anything about to be overwritten or removed is first copied to
  `.godot-ai/backups/<timestamp>/`, preserving relative paths — so an update is recoverable.
- **Explicit updates only:** the CLI never updates on its own.
- **Integrity:** remote fetches verify the tarball against its SHA-256 before extraction.

## Releasing (maintainers)

1. Update `CHANGELOG.md` (move `[Unreleased]` → `vX.Y.Z`).
2. Tag `vX.Y.Z` and push the tag.
3. Build per-OS/arch binaries (CI / goreleaser-style) and the content tarball + `.sha256`.
4. Create the GitHub release and attach the artifacts.

No remotes are configured and no releases are published by the toolkit itself — these are deliberate
maintainer actions.

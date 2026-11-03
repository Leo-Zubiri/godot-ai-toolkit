# Architecture

This document explains the model and the responsibility of each component. For *why* the project
exists and its principles, see the [rector README](../README.md).

## The model in one picture

```
            ┌─────────────────────────────┐
            │  skills/  (SINGLE SOURCE)    │  canonical SKILL.md (Agent Skills spec)
            │  references live inside each │
            └──────────────┬──────────────┘
                           │ consumed by
        ┌──────────────────┼────────────────────┐
        ▼                  ▼                     ▼
   adapters/           cmd/ + internal/     integrations/
  (per-agent DATA:     (Go CLI:             (MCP notes, Godot CLI
   discovery path,      install/update/      helpers, .NET helpers —
   instruction file)    status/validate)     optional capabilities)
                           │
                           ▼
              consumer game project
      .claude/skills/… .codex/skills/…  AGENTS.md
      .godot-ai/toolkit.lock.json   (version + commit + selection + hashes)
```

## Four responsibilities, kept separate

| Component | Is | Is not |
|---|---|---|
| **Skills** (`skills/`) | The knowledge: procedures/instructions for agents, authored once in the portable `SKILL.md` format. | Agent-specific copies. Code. |
| **Adapters** (`adapters/*.json`) | *Data* describing where each agent discovers skills/instructions. | Content. Code. A plugin system. |
| **CLI** (`cmd/`, `internal/`) | Tooling: install, update, status, validate; lock file; conflict/backup; fetch. | Knowledge. Coupled to any one agent. |
| **Integrations** (`integrations/`) | Optional capabilities an agent can *use* (MCP, Godot CLI, .NET). | Required dependencies. |

## Why skills are authored once

The open Agent Skills `SKILL.md` format is consumed by Claude Code, OpenCode and Codex. Crucially,
**OpenCode reads `.claude/skills/` directly**, so Claude Code and OpenCode share one layout. Only
Codex expects a different root (`.codex/skills/`). Therefore the only real cross-agent difference is a
**path** (and the instruction filename). We capture those as data in `adapters/` and let the CLI place
the one canonical skill where each selected agent looks. No content is duplicated; nothing drifts.

```
adapters/claude.json    → skillsPath ".claude/skills",  instructionsFile "CLAUDE.md"
adapters/opencode.json  → skillsPath ".claude/skills",  instructionsFile "AGENTS.md"   (reads Claude layout)
adapters/codex.json     → skillsPath ".codex/skills",   instructionsFile "AGENTS.md"
```

Adding a future agent that adopts Agent Skills = adding one JSON descriptor. No code change.

## Installation data flow

```
source (pinned toolkit version)        consumer project
┌───────────────┐                      ┌────────────────────────────┐
│ skills/<name> │ --select--> CLI ---> │ <skillsPath>/<name>/...     │  (per selected agent)
│ adapters/*    │           copy +     │ .godot-ai/toolkit.lock.json │  (version, commit, hashes)
└───────────────┘           hash       │ .godot-ai/backups/<ts>/     │  (only when overwriting)
                                       └────────────────────────────┘
```

- **Copy, not symlink** by default: each project is self-contained and independent of the central
  repo. (Symlink is offered only as a local-development convenience, never the default.)
- **Per-file SHA-256** is recorded in the lock so the CLI can detect local modifications before an
  update overwrites them, and back up anything it is about to replace.

## The CLI internals (packages)

| Package | Responsibility |
|---|---|
| `internal/skillspec` | Parse `SKILL.md` frontmatter (constrained YAML subset) and validate a skill directory. |
| `internal/adapter` | Load and validate `adapters/*.json` descriptors (decoupled from install logic). |
| `internal/lockfile` | Read/write `.godot-ai/toolkit.lock.json`. |
| `internal/fsutil` | Copy files/dirs, SHA-256, atomic writes. |
| `internal/fetch` | Download + checksum-verify + extract a content tarball (stdlib only). |
| `internal/installer` | Orchestrate install/update/status: selection, conflict detection, backup, lock. |
| `internal/cli` | Parse arguments, dispatch subcommands, format output. |

Standard library only — no third-party Go modules. See [distribution](distribution.md) for how the
binary reaches end users without imposing any runtime.

## What is deliberately NOT here (yet)

- A central `references/`/knowledge layer — references start inside each skill; extract only when
  reuse is demonstrated.
- A bundled MCP server — the Godot MCP ecosystem is community/experimental (see
  [`../integrations/README.md`](../integrations/README.md)).
- A plugin system, specialized sub-agents, or a complex config system — none are needed for the MVP.

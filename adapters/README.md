# adapters/ — per-agent data descriptors

Each `*.json` here tells the CLI **where a given agent discovers skills and instructions**. They are
**data, not code** — the installer core never hard-codes agent names. Supporting a new agent that
adopts the Agent Skills format is a matter of adding one descriptor here (README principle 2, 10).

## Schema

```json
{
  "agent": "claude",
  "skillsPath": ".claude/skills",
  "instructionsFile": "CLAUDE.md",
  "notes": "Free-text note about this agent's discovery behavior."
}
```

| Field | Required | Meaning |
|---|---|---|
| `agent` | ✅ | Identifier used on the CLI (`--agents claude,codex`). Must match the file's base name. |
| `skillsPath` | ✅ | Project-relative directory where the agent looks for skills. Must be relative, no `..`. |
| `instructionsFile` | ✅ | The agent's project instruction file (for documentation/wiring). |
| `notes` | optional | Human note; ignored by the installer logic. |

## Current adapters (verified discovery paths)

| Agent | `skillsPath` | `instructionsFile` | Note |
|---|---|---|---|
| `claude` | `.claude/skills` | `CLAUDE.md` | Claude Code; reads `AGENTS.md` only when no `CLAUDE.md` exists. |
| `opencode` | `.claude/skills` | `AGENTS.md` | OpenCode reads the Claude skills layout directly. |
| `codex` | `.codex/skills` | `AGENTS.md` | Codex uses its own skills root. |

Because Claude Code and OpenCode share `.claude/skills`, selecting both installs each skill **once**
there (the installer de-duplicates identical target paths). Codex gets a copy under `.codex/skills`.

See [`../docs/architecture.md`](../docs/architecture.md) for why this keeps a single source of truth.

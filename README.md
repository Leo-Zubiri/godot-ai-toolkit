# godot-ai-toolkit

> A centralized, versioned, reusable source of knowledge, instructions, skills and integrations
> for building **Godot** games with **AI agents** — shared across many projects and across
> different agents (Claude Code, OpenCode, Codex), on the developer's own schedule.

**This README is the rector document of the project.** It defines *why* the toolkit exists and the
principles it must not betray. Agents and contributors must consult it before any architectural,
skill, or distribution change. `AGENTS.md` is the short operational companion; it does not repeat
this document.

---

## 1. Purpose

Build a **single, central, reusable and evolving** base of Godot development knowledge — skills,
instructions, references and integrations — that any game project can incorporate and that any
compatible AI agent can use.

## 2. The problem it solves

Without a central source, every new game re-creates the same thing:

- The same Godot instructions and conventions, copy-pasted and then edited per project.
- Divergent, drifting, outdated versions of the *same* knowledge across games.
- Per-agent setup repeated by hand (Claude Code, OpenCode, Codex each configured separately).
- No way to know *which* version of the guidance a given game is actually using.

The toolkit replaces that with **one upstream** that games pull from deliberately.

## 3. Philosophy

- **One source of truth, pulled on demand.** Knowledge is authored once, versioned centrally, and
  incorporated into a game only when the developer chooses to install or update it.
- **Never auto-update a consumer.** Updates are always an explicit, deliberate action.
- **Each game is independent.** Different games can pin different toolkit versions and select
  different skills, forever, without affecting each other.

## 4. Multi-agent by design

The toolkit targets a **common format** (the open [Agent Skills](https://agentskills.io) `SKILL.md`
spec) so one body of knowledge serves multiple agents:

| Agent | How it discovers skills | Instruction file |
|---|---|---|
| **Claude Code** | `.claude/skills/<name>/` (project), `~/.claude/skills/` (global) | `CLAUDE.md` (reads `AGENTS.md` when no `CLAUDE.md`) |
| **OpenCode** | reads `.claude/skills/` directly, plus `.opencode/skills/`, `.agents/skills/` | `AGENTS.md` |
| **Codex** | `.codex/skills/<name>/` | `AGENTS.md` |

Because the format is shared, we **author skills once** and let the installer place them where each
agent looks. We do **not** duplicate skill content per vendor, and we do **not** couple the toolkit
to any single AI provider. New agents that adopt the Agent Skills format are supported by adding a
small **data** descriptor in `adapters/` — no new code.

## 5. Languages: GDScript and C#/.NET are both first-class

- **Shared Godot concepts** (nodes, scenes, signals, resources, lifecycle, physics, export, …) are
  expressed once, with language-specific examples shown side by side only where the engine API
  genuinely differs.
- **Language-specific concerns** stay separate (project setup, typing, async, NuGet, build/debug).
- **Their capabilities are not assumed identical.** The toolkit records real, verified differences —
  for example, **C# cannot export to Web/HTML5**, Android/iOS C# support is experimental, C# requires
  the **.NET editor build**, and the .NET SDK floor is **.NET 8** (**.NET 9** for Android export). See
  [`docs/compatibility.md`](docs/compatibility.md).

## 6. Project independence

Each consuming game:

- Chooses which agents it configures.
- Selects only the skills it needs.
- Pins an exact toolkit **version + commit**.
- Records what it has in [`.godot-ai/toolkit.lock.json`](docs/distribution.md).
- Updates **only when the developer decides**, never automatically.
- Keeps a version independent from every other game.

## 7. Evolution

The toolkit is meant to grow for years — new skills, new integrations, new Godot/.NET versions —
**without breaking existing consumers**. Because consumers pin versions and update explicitly, new
releases cannot destabilize a shipped game. SemVer communicates the nature of each change.

## 8. Non-negotiable principles

1. A single, centralized source of knowledge.
2. Compatibility with different AI agents and providers.
3. GDScript and C#/.NET are both first-class languages.
4. Share Godot concepts where it makes sense, without needless duplication.
5. Keep language-specific particularities separate where necessary.
6. Prefer simplicity, clarity and maintainability over over-engineering.
7. Prioritize official documentation and verifiable practice.
8. Never invent Godot APIs, methods, properties or behaviors.
9. Never assume GDScript and C# have identical capabilities or limitations.
10. Do not couple the toolkit to a single AI provider.
11. Never auto-update consumer projects.
12. Let each project select which skills it installs.
13. Let different projects pin different toolkit versions.
14. Prioritize working examples and real validation over voluminous but impractical docs.

## 9. Responsibility of the agents

Before proposing or making a significant change — architecture, new skills, or distribution — an
agent **must**:

- Read this README and confirm the change respects the purpose and the principles above.
- Verify any Godot/.NET technical claim against official documentation.
- If a change contradicts a principle, **flag it explicitly and justify** why it is necessary rather
  than quietly working around it.

See `AGENTS.md` for the operational checklist.

## 10. What it is / what it is not

**It is:**
- A central, versioned library of Godot knowledge as Agent Skills.
- A small, dependency-free CLI to install, update, inspect and validate that knowledge per project.
- A compatibility layer (data descriptors) so one source serves multiple agents.

**It is not:**
- A Godot plugin or engine fork.
- A game, a game template, or a runtime dependency of your game.
- An auto-updater, a package registry, or an AI provider.
- A dumping ground that restates the official Godot manual.

## 11. How it is organized

```
skills/        Official Godot knowledge catalog (SKILL.md) — single source of truth
adapters/      Per-agent data descriptors (discovery paths, instruction file)
cmd/ internal/ The Go CLI (install / update / status / validate) — tooling only
integrations/  Optional capabilities: MCP notes, Godot CLI helpers, .NET helpers
tests/         Fixtures (valid/invalid skills; gdscript/csharp/mixed projects) + checks
docs/          Architecture, compatibility matrix, skill authoring, distribution
```

Full detail in [`docs/architecture.md`](docs/architecture.md).

## 12. How to add a skill

See [`docs/authoring-skills.md`](docs/authoring-skills.md) and [`CONTRIBUTING.md`](CONTRIBUTING.md).
One skill is one folder `skills/<name>/SKILL.md`; the folder name must equal the `name` frontmatter
field; `name` + `description` are required. Validate with `godot-ai-toolkit validate`.

## 13. How it is distributed

Git + GitHub, **SemVer** tags and releases, a `CHANGELOG.md`, and per-OS/arch CLI binaries plus a
content tarball with a SHA-256 checksum. End users need **no runtime** (no Go, Python, Node or uv) —
only the single native binary. Details in [`docs/distribution.md`](docs/distribution.md).

## 14. How a consumer updates

```sh
# First time — pin a version, pick agents and skills:
godot-ai-toolkit install --tag v0.1.0 --agents claude,codex --skills <names...>

# Later, deliberately, when YOU decide:
godot-ai-toolkit update --tag v0.2.0

# See what is installed and whether any managed file was changed locally:
godot-ai-toolkit status
```

Updates are explicit. The CLI detects local modifications to managed files and backs them up before
overwriting. See [`docs/distribution.md`](docs/distribution.md).

## 15. How to contribute

See [`CONTRIBUTING.md`](CONTRIBUTING.md). In short: respect the principles, verify facts against
official docs, keep the single source of truth, avoid over-engineering, and make `go test ./...` and
`godot-ai-toolkit validate` pass.

## 16. Status: implemented vs planned

**Implemented (this phase — structure & tooling):**
- Repository architecture, rector README, `AGENTS.md`, governance and hygiene files.
- Docs: architecture, compatibility matrix, skill authoring, distribution.
- Go CLI: `install` / `update` / `status` / `validate`, lock file, conflict detection, backups,
  checksum-verified fetch — standard-library only.
- Agent adapters (Claude Code, OpenCode, Codex) as data descriptors.
- Test fixtures (valid + invalid skills; gdscript/csharp/mixed projects) and CI.

**Planned (deliberately deferred):**
- The actual Godot knowledge skills (the `skills/` catalog) — authored in later versions.
- MCP integration shipped as configuration; bundling/auto-wiring an MCP server (all community/
  experimental today — see [`integrations/README.md`](integrations/README.md)).
- A shared/central `references/` layer (kept inside skills until reuse is proven).
- Additional distribution channels (Claude plugin/marketplace, Homebrew/Scoop), binary signing.

> `skills/` is intentionally empty of skills right now; it holds only the taxonomy in
> [`skills/README.md`](skills/README.md). The example skill used to test the pipeline lives under
> `tests/fixtures/`, never in `skills/`.

---

## License

[MIT](LICENSE).

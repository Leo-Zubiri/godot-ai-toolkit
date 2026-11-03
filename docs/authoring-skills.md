# Authoring skills

Skills are the heart of the toolkit. They follow the open
[Agent Skills specification](https://agentskills.io/specification) so that Claude Code, OpenCode and
Codex can all consume them. **We author to the strict spec** for maximum portability — even though
some agents (e.g. Claude Code) are more lenient.

> No official Godot skills are authored yet. This document defines *how* they will be written. The
> example used to test the pipeline lives in `tests/fixtures/`, not in `skills/`.

## One skill = one folder

```
skills/<skill-name>/
├── SKILL.md            # required: YAML frontmatter + Markdown body
├── references/         # optional: detailed material loaded on demand
├── scripts/            # optional: executable helpers
└── assets/             # optional: templates/data
```

The folder name **must equal** the `name` frontmatter field.

## Frontmatter

```yaml
---
name: scenes-and-nodes
description: >-
  Godot scene/node architecture and the node lifecycle. Use when structuring
  scenes, instancing, or ordering _ready/_process logic in GDScript or C#.
license: MIT
compatibility: "Godot 4.x. C# requires the .NET editor build; no Web export for C#."
metadata:
  godot: "4.x"
  languages: "gdscript, csharp"
---
```

Rules enforced by `godot-ai-toolkit validate`:

| Field | Required | Constraint |
|---|---|---|
| `name` | ✅ | 1–64 chars, `^[a-z0-9]+(-[a-z0-9]+)*$`, **equals the folder name**. |
| `description` | ✅ | 1–1024 characters; state *what* the skill does **and when** to use it. |
| `license` | optional | Short identifier (e.g. `MIT`) or a bundled file reference. |
| `compatibility` | optional | ≤ 500 chars; environment/version/platform requirements. |
| `metadata` | optional | Map of string → string. |
| `allowed-tools` | optional | Space-separated tool list. **Experimental**; support varies by agent. |

### Supported frontmatter format (validator subset)

The validator parses a constrained YAML subset and **fails loudly** on anything it does not support,
rather than guessing:

- Top-level `key: value` pairs (values may be quoted with `"` or `'`).
- A single-line folded/plain scalar, or `>-` / `|`-style block scalars for `description`.
- A nested `metadata:` block of indented `key: value` string pairs.

Avoid YAML features outside this subset (anchors, top-level sequences, nested maps other than
`metadata`). Keep frontmatter simple — it is read by every agent.

## Body

- Keep `SKILL.md` **under ~500 lines / ~5k tokens**. Move depth into `references/` (progressive
  disclosure: agents load the body on activation, references only when needed).
- Write **practical, specific** instructions — not a restatement of the official manual.
- Include: when to use it, correct examples, common mistakes/pitfalls, and references.
- **Separate GDScript and C#** where the API differs; show both where it helps. Never assume parity.
- **Verify every Godot/.NET fact** against official docs before writing it (README principle 7–8).
  State version/build/platform constraints explicitly (see [compatibility](compatibility.md)).

## Referencing bundled files

Link to bundled material with relative paths:

```markdown
See [the signal patterns reference](references/signal-patterns.md).
```

The validator checks that such relative links resolve (reference integrity). Keep references one level
deep.

## Validate before committing

```sh
go run ./cmd/godot-ai-toolkit validate            # validates skills/ + adapters/
go run ./cmd/godot-ai-toolkit validate --skills-dir tests/fixtures/skills/valid
```

A skill with a missing/invalid `name`, an over-long `description`, a folder/name mismatch, or a broken
relative reference is **rejected**.

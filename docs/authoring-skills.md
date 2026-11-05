# Authoring skills

Skills are the heart of the toolkit. They follow the open
[Agent Skills specification](https://agentskills.io/specification) so that Claude Code, OpenCode and
Codex can all consume them. **We author to the strict spec** for maximum portability — even though
some agents (e.g. Claude Code) are more lenient.

> The first production skill — [`signals`](../skills/signals/SKILL.md) — is authored. This document defines
> *how* skills are written. Throwaway pipeline/example fixtures live in `tests/fixtures/`, never in `skills/`.

## One skill = one folder

```
skills/<skill-name>/
├── SKILL.md            # required: YAML frontmatter (name + description) + Markdown body
├── manifest.json       # required: toolkit-owned metadata (compatibility, examples, sources)
├── examples/           # optional: canonical, verifiable example files/projects
├── references/         # optional: detailed material loaded on demand
├── scripts/            # optional: executable helpers
└── assets/             # optional: templates/data
```

The folder name **must equal** the `name` frontmatter field (and the `name` in `manifest.json`).

## Frontmatter — the contract is `name` + `description` only

```yaml
---
name: scenes-and-nodes
description: >-
  Godot scene/node architecture and the node lifecycle. Use when structuring
  scenes, instancing, or ordering _ready/_process logic in GDScript or C#.
---
```

Rules enforced by `godot-ai-toolkit validate`:

| Field | Required | Constraint |
|---|---|---|
| `name` | ✅ | 1–64 chars, `^[a-z0-9]+(-[a-z0-9]+)*$`, **equals the folder name**. |
| `description` | ✅ | 1–1024 characters; state *what* the skill does **and when** to use it. |

**No other frontmatter key is allowed.** `compatibility`, `metadata`, `license`, `allowed-tools` and any
other field are **rejected** — that information lives in `manifest.json` (below). This keeps the frontmatter
maximally portable across Claude Code / OpenCode / Codex.

### Supported frontmatter format (validator subset)

The validator parses a constrained YAML subset and **fails loudly** rather than guessing:

- The two `key: value` pairs (`name`, `description`); values may be quoted with `"` or `'`.
- A folded/plain scalar, or a `>-` / `|`-style block scalar for `description`.
- Any disallowed key (including nested blocks such as `metadata:`) is recognized and **rejected**.

Avoid YAML features outside this subset (anchors, top-level sequences, nested maps).

## Skill manifest (`manifest.json`)

Richer metadata lives here, **not** in frontmatter. It is the per-skill **source of truth**; the aggregated
`skills/catalog.json` is **generated** from all manifests and is never hand-edited. Shape:

```json
{
  "schemaVersion": 1,
  "name": "scenes-and-nodes",
  "languages": ["gdscript", "csharp"],
  "godot": { "minimum": "4.2", "tested": ["4.7.2"], "reference": "4.7.2" },
  "dotnet": { "minimum": "8.0", "notes": "Android export needs .NET 9+; no C# web export." },
  "examples": [
    { "path": "examples/scene_root.gd",     "language": "gdscript", "tier": "complete",   "verify": "parse" },
    { "path": "examples/csharp/SceneDemo",   "language": "csharp",   "tier": "executable", "verify": "build" }
  ],
  "sources": [
    { "claim": "C# cannot export to the web", "url": "https://docs.godotengine.org/en/4.7/tutorials/scripting/c_sharp/index.html", "version": "4.7" }
  ]
}
```

- **Compatibility model:** `godot.minimum` (lowest supported, justified), `godot.tested` (actually verified),
  `godot.reference` (canonical target; baseline **4.7.2**). Mirror for `.NET` where C# applies. Must not
  contradict [compatibility.md](compatibility.md).
- **Examples** are categorized by `tier` — *illustrative* (inline, review-only, not listed here),
  *complete* (self-contained file: parse/compile), *executable* (runnable project). `verify` ∈
  `parse | build | run | build+run | none`.
- **`sources[]`** mirrors the body's citations for audit/tooling.
- `verification` status is written by `verify` — never hand-set.
- A size exception may be recorded with `"sizeExceptionReason": "<why>"` (see Body).

## Body

- Keep `SKILL.md` **under ~500 lines / ~5k tokens** (a **warning**, not a hard limit; justify exceptions via
  `manifest.sizeExceptionReason`). Move depth into `references/` (progressive disclosure).
- Write **practical, specific** instructions — not a restatement of the official manual.
- Include: when to use it, correct examples, common mistakes/pitfalls, references, and a **`## Sources`** section.
- **Cite every non-obvious technical claim** (APIs, engine behavior, version differences, language
  restrictions, lifecycle, performance) with a URL + Godot version. **Distinguish documented facts from the
  toolkit's own recommendations** (label recommendations as such).
- **Separate GDScript and C#** where the API differs; show both where it helps. Never assume parity.
- Verify every Godot/.NET fact against official docs before writing it (README principle 7–8); state
  version/build/platform constraints explicitly (see [compatibility](compatibility.md)).

## Referencing bundled files

Link to bundled material with relative paths:

```markdown
See [the signal patterns reference](references/signal-patterns.md).
```

The validator checks that such relative links resolve (reference integrity). Keep references one level
deep.

## Validate and verify before committing

```sh
# Structural + metadata: SKILL.md contract, manifest schema, reference integrity, size warnings.
go run ./cmd/godot-ai-toolkit validate --skills-dir tests/fixtures/skills/valid

# Engine-level: parse GDScript / build C# for the manifest's complete & executable examples.
go run ./cmd/godot-ai-toolkit verify   --skills-dir tests/fixtures/skills/valid
```

- `validate` **rejects** a missing/invalid `name`, an over-long `description`, a **disallowed frontmatter key**,
  a folder/name mismatch, a broken relative reference, or an invalid `manifest.json`.
- `verify` reports **PASS / FAIL / SKIPPED / NOT_APPLICABLE** per example and globally; a **skipped mandatory
  check never yields a global PASS**, and C# is reported as **BUILD_VERIFIED** (never runtime) for now. See
  [`proposals/verification-workflow.md`](proposals/verification-workflow.md).

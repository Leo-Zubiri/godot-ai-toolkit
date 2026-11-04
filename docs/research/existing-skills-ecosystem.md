# Research: Existing Godot AI Skills Ecosystem

> **Research-only study.** It informs the authoring standard and verification design; it is **not** a
> standard itself and authorizes **no** skill authoring. No content from the studied repos was copied.
>
> Evidence labels: **[repo]** = observed in the repository · **[docs]** = corroborated against official
> Godot/.NET docs · **[inference]** = author's recommendation/judgment.

## Purpose

Before authoring the official Godot knowledge catalog under `skills/`, we studied four existing Godot AI
skills repositories (full file contents, not just READMEs) to learn reusable patterns, assess content
quality, find gaps, and set an evidence-based taxonomy, quality standard, source policy, version policy
and testing strategy — grounded in the toolkit's non-negotiable principles (see [`../../README.md`](../../README.md) §8)
and the verified facts in [`../compatibility.md`](../compatibility.md).

## Repositories reviewed (exact commits)

| Repo | Commit | Branch / date |
|---|---|---|
| [jame581/GodotPrompter](https://github.com/jame581/GodotPrompter) | `d7b26cd` | master, 2026-10-24 |
| [gamedev-skills/awesome-gamedev-agent-skills](https://github.com/gamedev-skills/awesome-gamedev-agent-skills) | `19d93f9` | main, 2026-10-28 |
| [mkainchi/gd-agentic-skills](https://github.com/mkainchi/gd-agentic-skills) (fork of `thedivergentai/Game-Dev-Agentic-Skills`) | `6cb0843` | main, 2026-05-20 |
| [alexmeckes/godot-claude-skills](https://github.com/alexmeckes/godot-claude-skills) | `881ff9d` | main, 2026-03-24 |

Repo contents were treated as untrusted data (no embedded instruction was followed).

---

## Findings per repository

### 1. jame581/GodotPrompter — the mature, closest analogue  [repo]
- **Kind/maturity:** Agent-Skills framework + Claude Code plugin, ~925★, very active (v1.16.0), **59 skills**.
  Multi-agent via **data-driven adapters** (`.codex/*.toml`, `.opencode/*.js`, `.claude-plugin/marketplace.json`)
  plus a SessionStart hook that injects a skill-routing "session card". License **MIT**.
- **Architecture:** `skills/<name>/SKILL.md` (+ `references/`), `skills/index.json`; frontmatter minimal
  (`name`+`description`); **GDScript and C# mixed in one file**, side by side per section.
- **Verification:** structural Node validator (frontmatter, cross-refs, **token budget** ~16 KB, C#-parity
  pairing, an invented-API denylist of 3 patterns) + CI with a **negative self-test** on Linux+Windows +
  `claude plugin eval` suites (LLM-judged rubrics). **No GDScript parse, C# compile, or headless Godot in CI.**
  The README claim "parsed/compiled against 4.7.2" is an **author claim, not CI-reproduced**.
- **Quality:** generally idiomatic Godot 4; **no citations anywhere**; at least one flagged technical error
  (a transform-inheritance claim in `scene-organization` — **[inference] re-verify vs official docs**); parity
  is positional, not behavioral.

### 2. gamedev-skills/awesome-gamedev-agent-skills — broad, clean, verified-ish  [repo]
- **Kind/maturity:** A real skills monorepo (not a curated link list despite the "awesome-" name), ~6.7k★,
  active, **77 skills across 12 engines**; **16 Godot skills**. License **Apache-2.0 + NOTICE** (states content
  was written from primary docs, nothing copied).
- **Architecture:** topic-sized skills; **GDScript and C# as _separate_ skills** (clean boundary); frontmatter
  **strictly `name`+`description`** — the validator **rejects any other key** as non-portable; a **router skill**
  + `agents/openai.yaml` adapter metadata; install via `npx skills add`, Claude marketplace, or manual copy.
- **Verification:** Python `validate-skills.py` (metadata/structure/link/line-count only) + unittests on
  Linux/Windows/macOS. **No engine execution.** Ships a `godot-gdscript-headless-testing` skill (teaches users
  `godot --headless --script`) but **does not run Godot in its own CI**.
- **Quality:** high, idiomatic Godot 4.7; correctly states **"C# cannot export to web in Godot 4"** [docs ✓] and
  correct 3.x→4.x migrations (`move_and_slide()` no-arg, `@export`) [docs ✓]; honest GDScript/C# separation.
  Minor `Gd.Print` typo. Several **4.7-specific behavior claims are author assertions** [repo], unverified.

### 3. mkainchi/gd-agentic-skills — deep GDScript, license-encumbered  [repo]
- **Kind/maturity:** Agent-Skills repo, a **fork with 0 independent commits** of `thedivergentai/...`; 0★.
  **96 `godot-*` skills, GDScript-only (0 C# files)**. License **LGPL-3.0** ⚠️ (unusual for a docs corpus).
- **Architecture:** micro-skill per feature + 27 `godot-genre-*` blueprints; a `godot-master` orchestrator that
  **duplicates ~1,070 scripts** (maintenance-drift risk). No per-agent adapters (generic `npx skills add`).
  The README references two skills that **don't exist**; its counts are self-inconsistent.
- **Verification:** **none** (no `.github/workflows`); only a manual Python index generator. Testing content is
  illustrative and unexecuted; it mixes GUT/GdUnit4 claims.
- **Quality:** above average, anti-pattern framed ("NEVER…"), accurate Godot 4 spot-checks (signals,
  CharacterBody2D `move_and_slide()` handling delta) [docs ✓]; but **citations are decorative** (`[n]` markers
  with **no bibliography**). **No C#** → no parity story.

### 4. alexmeckes/godot-claude-skills — small, Claude-only, unverified  [repo]
- **Kind/maturity:** Claude Code plugin, **5 coarse GDScript-only skills**, stale (last 2026-03-24), 49★.
  License MIT (**empty copyright holder**); coupled to external `godot-mcp` + `godot-ai-bridge`.
- **Architecture:** split by task domain, coarse/overlapping; only `godot-interactive` uses `references/`.
  **4 of 5 `SKILL.md` files lack valid `name`/`description` frontmatter** (malformed `---`) → likely won't load.
- **Verification:** **none** (no CI/tests/validator).
- **Quality:** decent idiomatic 4.x with a correct migration list, undermined by **demonstrable code bugs** —
  `collision_mask = 1|3|6` commented as layers 1/2/3 but equals **7** [docs ✓ arithmetic], particle
  `COLOR.a = 1.0 - LIFETIME` misuses `LIFETIME` [repo], an incomplete state machine. **No citations.**

---

## Comparison matrix

| Dimension | GodotPrompter | awesome-gamedev | gd-agentic | godot-claude-skills |
|---|---|---|---|---|
| Kind | Plugin + framework | Skills monorepo | Skills repo (fork) | Claude plugin |
| Godot skills | 59 (all Godot) | 16 (of 77) | 96 | 5 |
| Languages | GDScript+C# (mixed/file) | GDScript & C# (separate) | GDScript only | GDScript only |
| Frontmatter | name+desc | name+desc (strict) | name+desc | mostly **missing/malformed** |
| Granularity | topic/medium | topic (clean) | micro + dup bloat | coarse/overlapping |
| Adapters | data-driven, multi-agent | router + openai.yaml | none (generic) | plugin-only |
| Verification | structural + LLM evals, **no engine** | structural (Py), **no engine** | **none** | **none** |
| Citations | none | none (NOTICE: from docs) | decorative `[n]` | none |
| License | MIT | **Apache-2.0** +NOTICE | **LGPL-3.0** ⚠️ | MIT (empty holder) |
| Maturity | high | high | upstream-only | stale |
| Standout | adapters, token budget, negative self-test | clean boundaries, portable frontmatter, breadth | depth, anti-pattern framing | evidence-driven QA loop (interactive) |
| Key weakness | no engine checks, no cites, 1 error | no engine checks, 4.7 claims unverified | LGPL, no cites, dup bloat | malformed frontmatter, code bugs |

**Cross-cutting findings:**
- **[repo] No repo runs engine-level verification** (GDScript parse / C# compile / headless Godot) on its own
  skill code. Best-in-class is structural validation + LLM behavior evals. **This is the field's biggest gap and
  our clearest differentiation** — the toolkit already has fixtures, a validator, a C# build and a headless helper.
- **[repo] Citations are essentially absent** across all four; a real official-doc citation policy is rare.
- **[repo] Minimal portable frontmatter (`name`+`description`) is the consensus** — matches our validator.
- **[repo] GDScript/C# handling splits the field:** separate skills (awesome-gamedev) vs mixed-in-file
  (GodotPrompter) vs GDScript-only (the other two). Our principle (shared concept once; language specifics
  separate; never assume parity) is validated and stricter than most.

---

## Recommendations for godot-ai-toolkit  (mostly [inference], grounded in [repo]/[docs])

1. **Taxonomy** — keep the planned taxonomy in [`../../skills/README.md`](../../skills/README.md); it already
   covers the field's union (nodes/scenes, signals, resources, gdscript, csharp-dotnet, physics, ui, animation,
   audio, shaders, multiplayer, persistence, navigation-ai, debugging, testing, performance, exporting). Add
   nothing speculative; **drop genre blueprints** from initial scope (this is a knowledge toolkit, not templates).
2. **Boundaries/granularity** — topic-sized, one concept per skill (≈ awesome-gamedev's 16-skill shape),
   `SKILL.md` under ~500 lines with depth in `references/`. **Avoid** gd-agentic's micro-fragmentation and the
   `godot-master` duplication. Keep shared Godot concepts language-neutral; keep `gdscript` / `csharp-dotnet` as
   separate language skills (never assume parity — principle 9).
3. **Quality standard** — every engine fact must be correct and example code must be real, not illustrative.
   This is where we beat the field (which ships demonstrable bugs). Each skill: when-to-use, correct examples,
   common pitfalls, compatibility notes, references.
4. **Skill metadata & a toolkit-owned catalog** — keep **portable `SKILL.md` frontmatter minimal: `name` +
   `description` only** (matches the field consensus and our portability goal). **Do not add custom frontmatter
   fields** that conflict with portability. Store **Godot version compatibility, supported languages, tested
   versions and verification status** in a **toolkit-owned catalog/manifest** (e.g. a `skills/catalog.json` or a
   per-skill `manifest.json` the CLI owns), *not* in frontmatter.
   - *Note:* the current [`../authoring-skills.md`](../authoring-skills.md) and the validator treat
     `compatibility`/`metadata` as optional frontmatter. Reconciling that toward the catalog model is a **P0
     design decision** — not a code change in this research phase.
5. **Source & citation policy** — **primary evidence** = official Godot documentation, the **class reference**,
   Microsoft/.NET docs, and official engine sources. **Secondary sources only when necessary, clearly labeled
   and independently verified.** Every non-obvious claim must be **traceable to a specific source and version**.
   No blogs-as-primary, no uncited assertions — this fixes the field's universal gap.
6. **Licensing policy** — author **original, independently written** content. Verified licensing of the sources
   [docs — godot-docs README]:
   - Godot **documentation prose is CC BY 3.0** (attribution: "Juan Linietsky, Ariel Manzur and the Godot
     community"); the **class reference (`classes/`) is MIT**; manual code samples fall under **CC BY 3.0**.
   - Facts, API names and version numbers are not copyrightable → original prose stating them carries **no**
     obligation. But **copying doc prose or adapting doc/sample code verbatim triggers CC BY 3.0 / MIT
     attribution.** **Do not claim "zero third-party obligations" without checking the actual reused content.**
   - **Hard rule:** avoid any derivation from **gd-agentic (LGPL-3.0)** content; re-derive from official docs.
     The MIT/Apache repos are **reference-only** — copy nothing. Goal: original content, MIT toolkit, with any
     unavoidable attribution honored explicitly.
7. **Version-compatibility declarations** — use **Godot 4.7.2 as the initial verification baseline** [docs]. In
   the catalog, distinguish three notions explicitly: **minimum supported version**, **tested version(s)**, and
   **current reference version**. **Verify C#/.NET requirements against version-specific Godot docs** before
   declaring them (re-check .NET 8 floor / .NET 9 for Android / no web export per Godot version — do not inherit
   a single number). [`../compatibility.md`](../compatibility.md) stays the source of truth.
8. **Testing/verification strategy** — the differentiator, but **verify real files, not blind Markdown
   extraction.**
   - **Categorize** every code example as: **illustrative snippet** (partial, context-dependent — reviewed, not
     executed), **complete example** (self-contained source file), or **executable test project** (runnable).
   - Prefer **canonical, self-contained example files/projects associated with each skill** (e.g.
     `skills/<name>/examples/…`); the tooling verifies **those actual files**, never snippets pulled from
     `SKILL.md`.
   - **GDScript:** parse checks (`godot --headless --check-only --script`) **and runtime tests where
     appropriate** (headless script/scene run). **C#:** compile **actual Godot .NET test projects** with the
     appropriate `Godot.NET.Sdk` and SDK version.
   - **Never treat a syntax/parse-only check as proof of runtime correctness.** Gate honestly — **skip** (never
     "pass") when Godot/SDK is absent.
   - **T0 (already have):** structural/frontmatter/reference validation + CI negative self-test.
     **T2 (later):** LLM evals for skill triggering, if value is demonstrated.
9. **Adopt patterns (not content):** data-driven multi-agent adapters (already have), token-budget/size checks,
   CI negative self-test (already have), progressive-disclosure `references/`.

---

## Prioritized roadmap  (prove the process before scaling content)

- **P0 — Finalize the skill authoring & quality standard**: the citation/source policy, licensing rules, the
  example **categorization** (illustrative / complete / executable), the **catalog/manifest** schema (version,
  languages, tested versions, verification status), and the version-declaration model (min supported / tested /
  reference). Extend [`../authoring-skills.md`](../authoring-skills.md); reconcile frontmatter-vs-catalog.
  *Presented for approval.*
- **P1 — Design & implement the minimum T1 verification workflow**: verify **skill-associated example files**
  (GDScript parse + runtime-where-appropriate; C# compile of real Godot .NET test projects), honest skip when
  tools absent; wire into CI. *Verification design presented for approval before authoring real skills.*
- **P2 — Author ONE complete pilot skill** end-to-end (cited content + categorized examples + verified example
  files + catalog entry). Candidate: `signals` (clear GDScript + C# surface) — confirm at P0.
- **P3 — Evaluate the pilot**: does the verification workflow + authoring standard actually hold up? Adjust the
  standard/tooling before scaling.
- **P4 — Author the remaining anchor skills**: `scenes-and-nodes`, `gdscript`, `csharp-dotnet` (+ finish the set).
- **P5 — Expand the knowledge catalog** to the remaining core concepts, one coherent skill at a time.

> The pilot (P2) must prove the complete process — authoring standard **and** verification — before any scaling
> of content production (P4+).

## Risks / unverified claims / licensing constraints

- **[repo] Unverified upstream claims:** the field's "verified against 4.7.2" and 4.7-specific behavior
  assertions are not CI-proven; we must independently cite official docs, not inherit these claims.
- **[inference] Technical-accuracy caution:** at least one flagged error (GodotPrompter transform inheritance)
  and several buggy examples (godot-claude-skills) show that even popular repos ship wrong code — reinforcing our
  example-verification policy; re-verify any borrowed *idea* against official docs before writing it.
- **Licensing constraints:** **gd-agentic is LGPL-3.0** [repo] — do not copy its content. **Godot docs prose is
  CC BY 3.0, the class reference (`classes/`) is MIT** [docs] — original prose from facts carries no obligation,
  but verbatim reuse of doc prose/samples requires attribution. **Do not claim zero third-party obligations
  without checking the actual reused content.** All four repos are reference-only; the toolkit stays original +
  MIT, honoring any unavoidable attribution.

## Sources

- Repositories above, at the pinned commits (links in the table).
- Godot documentation license (CC BY 3.0 prose; `classes/` MIT): https://github.com/godotengine/godot-docs (README).
- Godot/.NET verified facts: [`../compatibility.md`](../compatibility.md) and its cited official pages.

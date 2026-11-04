# Proposal P0 — Skill Authoring & Quality Standard

> **Status: APPROVED (design) — not yet implemented.** Design/documentation only. The open decisions have been
> **reviewed and approved**; they are recorded in *Resolved decisions* below and folded into the text. This does
> **not** authorize authoring real skills or changing the installer/CLI/adapters/fixtures. It builds on
> [`../research/existing-skills-ecosystem.md`](../research/existing-skills-ecosystem.md) and must respect the
> non-negotiable principles in [`../../README.md`](../../README.md) §8.

## Scope

This standard governs **how a skill is authored and when it is accepted** into `skills/`. It covers: source
& citation policy, licensing, example categorization, the skill metadata catalog, the version-compatibility
model, and the quality/maintenance/acceptance criteria. The verification mechanics live in the companion
proposal [`verification-workflow.md`](verification-workflow.md) (P1).

---

## 1. Source, citation & attribution policy

**Primary sources (required basis for every engine claim):**
- Official Godot documentation — `docs.godotengine.org`, **version-pinned URLs** (e.g. `/en/4.7/…`, never `/stable/`).
- The Godot **class reference** (per-class pages / the engine's XML, emitted via `--doctool`).
- **Microsoft / .NET** official docs for C#/.NET facts.
- Official Godot engine sources: release notes, migration guides, the engine repo.

**Secondary sources:** allowed only when primary sources are silent, **clearly labeled as secondary** and
**independently verified** (by a doc cross-check or by an executable example). A blog/third-party repo is
never the sole basis for an API claim.

**Citation rule (resolved):** **back every non-obvious technical claim** with a source and the Godot version it
was checked against — in particular: **APIs/signatures, engine behavior, differences between versions, language
restrictions (GDScript vs C#), node/object lifecycle, and performance claims.** Trivially-true prose needs no
citation. **Clearly distinguish documented facts from the toolkit's own architectural recommendations** — facts
carry a citation; a recommendation is labeled as such (e.g. "Recommendation:") and does not masquerade as a
documented guarantee.

**Citation storage (resolved — both are required):**
- **In-body:** a `## Sources` section at the end of each `SKILL.md` listing `claim → URL → version` — optimized
  for the agent reading the skill.
- **Machine-tracked:** a `sources[]` array in the skill manifest (§4) — optimized for audit and tooling, and it
  survives independent of prose. Both are mandatory; they serve different consumers.

## 2. Licensing policy

Author **original, independently written** content. Verified license facts of the sources
([docs — godot-docs README](https://github.com/godotengine/godot-docs)):

- Godot **documentation prose → CC BY 3.0** (attribution: "Juan Linietsky, Ariel Manzur and the Godot community").
- Godot **class reference (`classes/`) → MIT**.
- Manual **code samples → CC BY 3.0**.

Rules:
- Facts, API names, signatures and version numbers are **not copyrightable** → original prose stating them
  carries **no obligation**. Prefer writing original examples.
- **Do not** copy doc prose or adapt doc/sample code verbatim. If a doc sample is ever reused ~verbatim, it
  **must** carry the required CC BY 3.0 / MIT attribution, recorded in a `THIRD-PARTY-NOTICES.md` (created only
  if/when first needed).
- **Hard rule:** no derivation from **gd-agentic-skills (LGPL-3.0)** content; re-derive from official docs.
  The MIT/Apache external repos are **reference-only** — copy nothing.
- Never claim "zero third-party obligations" without checking the actual reused content.

The toolkit stays original + MIT.

## 3. Example categorization

Every code example belongs to exactly one tier; the tier determines how (and whether) it is verified (P1):

| Tier | What it is | Lives in | Verified by |
|---|---|---|---|
| **Illustrative snippet** | Partial, context-dependent code inside `SKILL.md` fences. | `SKILL.md` body | **Review only** — never executed; must be correct-by-review and clearly fragmentary. |
| **Complete example** | A self-contained source file that stands on its own. | `skills/<name>/examples/…` | **Parse (GDScript)** / **compile (C#)**. |
| **Executable verification project** | A minimal, runnable Godot project (optionally with a tiny test). | `skills/<name>/examples/<proj>/` | **Headless run / test** (highest assurance). |

- **Default for inline code is "illustrative."** A snippet is never treated as proof of anything.
- A skill SHOULD back its key claims with at least one **complete** example; **executable** projects are
  encouraged for behaviorally-subtle topics (signals, physics, multiplayer).
- Each example's tier and language are declared in the manifest (§4) so the verifier knows what to do.

## 4. Skill metadata & the toolkit-owned catalog

**Frontmatter stays minimal and portable: `name` + `description` only.** No custom frontmatter fields that
would break portability across Claude Code / OpenCode / Codex.

All richer metadata lives in a **toolkit-owned manifest**, co-located with each skill.

**Manifest shape (resolved):** a per-skill **`skills/<name>/manifest.json` is the single source of truth**
(co-located, scales, easy to review per PR). An aggregated **`skills/catalog.json` is generated automatically by
the CLI** from the manifests and is **never hand-edited** — the two are not maintained in parallel.

**Proposed per-skill `manifest.json` (example — illustrative):**

```json
{
  "schemaVersion": 1,
  "name": "signals",
  "languages": ["gdscript", "csharp"],
  "godot": { "minimum": "4.2", "tested": ["4.7.2"], "reference": "4.7.2" },
  "dotnet": { "minimum": "8.0", "notes": "Android export requires .NET 9+; no C# web export." },
  "examples": [
    { "path": "examples/signal_bus.gd",        "language": "gdscript", "tier": "complete",   "verify": "parse" },
    { "path": "examples/csharp/SignalsDemo",    "language": "csharp",   "tier": "executable", "verify": "build+run" }
  ],
  "sources": [
    { "claim": "signal.emit() is the 4.x form of emit_signal()", "url": "https://docs.godotengine.org/en/4.7/tutorials/scripting/gdscript/gdscript_basics.html", "version": "4.7" }
  ],
  "verification": { "status": "unverified", "lastRun": null, "lastVersion": null }
}
```

- `verification.status` is **written by the tool**, never hand-set (see P1 for the status model).
- Keeping this out of frontmatter preserves portability while giving the toolkit a rich, checkable record.

## 5. Compatibility & version model

Three distinct notions, recorded per skill in `manifest.godot` (and mirrored for `.NET`):

- **minimum supported** — the lowest Godot version the content is valid for (stated **with justification**,
  cited to version-specific docs). Default conservative; raise only with evidence.
- **tested** — versions actually run through verification (initially `["4.7.2"]`).
- **reference** — the canonical version examples target. **Baseline = Godot 4.7.2.**

Rules:
- **Verify C#/.NET requirements against version-specific Godot docs** before declaring them; do not inherit a
  single number (e.g. .NET 8 floor / .NET 9 for Android / **no C# web export** must be re-checked per version).
- [`../compatibility.md`](../compatibility.md) remains the **single source of truth** for the matrix; a skill's
  manifest declares only its own support window, which must not contradict the matrix.

## 6. Quality, maintenance & acceptance criteria

**Acceptance checklist — a skill is "accepted" only when all hold:**

1. **Structure:** spec-valid `SKILL.md` (folder == `name`, `description` ≤1024, minimal frontmatter) — the
   existing `toolkit validate` passes.
2. **Required body sections:** *When to use*, *Correct usage*, *Common pitfalls / anti-patterns*,
   *Compatibility notes*, *Sources*.
3. **Language handling:** shared Godot concept written once; GDScript vs C# specifics kept separate; **no false
   parity claims** (principle 9). If only one language applies, say so explicitly.
4. **Citations:** every non-obvious claim cited (URL + version) per §1.
5. **Examples:** each categorized (§3); all **complete/executable** examples verify **PASS** (or **SKIPPED**
   with the exact missing dependency — **never FAIL**, never a false PASS; see P1).
6. **Manifest:** present and schema-valid; compatibility declared (§4/§5).
7. **Licensing:** original content; no external copying; §2 satisfied.
8. **Size budget (resolved):** target **≤ 500 lines / ~5,000 tokens** per `SKILL.md`, with depth pushed to
   `references/`. Tooling emits an **automatic warning** when exceeded — **not a hard failure** that would force
   sacrificing quality. **Justified exceptions are allowed** and recorded (e.g. a `sizeExceptionReason` note in
   the manifest) so the warning is acknowledged rather than silently ignored.

**Maintenance policy:**
- On a new Godot release, re-run verification; update `manifest.godot.tested` and `verification`.
- Record changes in `CHANGELOG.md`.
- Stale policy: if a skill's `tested` set lags the current reference by more than one minor version, flag it for
  re-verification (reported, not auto-changed).

---

## Resolved decisions (approved)

1. **Citations storage:** `## Sources` in each `SKILL.md` **and** `sources[]` in its `manifest.json` (both; §1).
2. **Catalog:** per-skill `manifest.json` is the source of truth; `catalog.json` is auto-generated, not
   hand-maintained (§4).
3. **Size:** target ≤500 lines / ~5k tokens per `SKILL.md` as **automatic warnings**, not hard limits; justified
   exceptions allowed (§6.8).
4. **Pilot skill:** `signals`.
5. **Citation strictness:** back every non-obvious technical claim (APIs, engine behavior, version differences,
   language restrictions, lifecycle, performance); distinguish documented facts from own recommendations (§1).
6. **Versions:** keep the `minimum / tested / reference` model; baseline **Godot 4.7.2** (§5).

The P1 verification decisions are recorded in [`verification-workflow.md`](verification-workflow.md).

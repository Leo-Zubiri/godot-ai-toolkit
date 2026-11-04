# Proposal P1 — Minimal Verification Workflow

> **Status: APPROVED (design) — not yet implemented.** Design/documentation only. The decisions are **reviewed
> and approved** (see *Resolved decisions*). Where tooling is needed it remains a **gated change**: no
> installer/CLI/adapter/fixture code is modified in this consolidation phase. Companion to
> [`authoring-and-quality-standard.md`](authoring-and-quality-standard.md) (P0).

## Principles

1. **Verify real files, not Markdown.** The workflow verifies **skill-associated example files/projects** (P0 §3),
   never snippets extracted from `SKILL.md`.
2. **A parse/syntax check is not proof of runtime correctness.** Parsing and compiling are necessary, not
   sufficient; runtime behavior needs an actual run.
3. **Honest statuses, no false positives.** Every check resolves to **PASS / FAIL / SKIPPED / NOT_APPLICABLE**.
   A missing tool is **SKIPPED** (with the exact missing dependency), never PASS.
4. **Reuse before building.** Prefer existing components; new tooling must be justified and approved.

## Reuse of existing components (reuse-first)

| Need | Existing component | Reuse? |
|---|---|---|
| SKILL.md structural + frontmatter + reference validation | `internal/skillspec` + `toolkit validate` | **Reuse as-is.** |
| Negative self-test (validator must reject broken input) | invalid fixtures + CI | **Reuse as-is.** |
| Headless Godot project check + stdout error scan | `integrations/godot-cli/validate-headless.sh` | **Reuse**; extend for per-script parse (small). |
| C# compilation | `dotnet build` on Godot.NET.Sdk projects (proven on fixtures) | **Reuse as-is.** |
| Example Godot projects (gdscript/csharp/mixed) | `tests/fixtures/` | **Reuse as templates** for `examples/`. |
| CI matrix (Win/Linux/macOS) | `.github/workflows/validate.yml` | **Reuse**; add a Linux Godot 4.7.2 engine tier (resolved). |

## Verification layers

Each layer declares, per example, a status. Layers are ordered; a FAIL in an earlier layer short-circuits later
ones for that example.

1. **Structural validation** *(have it)* — `SKILL.md` spec + reference integrity, **+ manifest schema** (new,
   small). Runs always, no engine needed. PASS/FAIL.
   - *Gated tooling:* manifest-schema validation would extend the validator → **justify & approve before coding.**
2. **Parse / syntax**
   - **GDScript:** `godot --headless --check-only --script <file>` per *complete/executable* GDScript example.
   - **C#:** `dotnet build` of the example project (restores `Godot.NET.Sdk`).
   - SKIPPED if the required tool (Godot binary / .NET SDK) is absent.
3. **Compilation**
   - **C#:** the `dotnet build` above is the compile gate (success = compiles).
   - **GDScript:** no separate compile step exists; parse (`--check-only`) is the closest static gate. *Do not*
     label a passed parse as "compiled."
4. **Runtime** *(executable tier only)*
   - **GDScript:** run the example project headless (`godot --headless --path <proj> --quit`, or a `--script`
     test runner), **scanning stdout for `SCRIPT ERROR`/`ERROR`/`Parse Error`** (headless Godot can exit 0 on
     failure — verified caveat), or a framework runner (GUT/GdUnit4) that exits non-zero on failure.
   - **C#:** requires the **.NET-enabled Godot editor build**; run via gdUnit4Net `[RequireGodotRuntime]` or a
     headless run. Pure-logic C# can be tested with xUnit/NUnit without Godot.
   - SKIPPED if no (suitable) Godot binary.

### GDScript vs C#/.NET — what runs

| | GDScript | C#/.NET |
|---|---|---|
| Parse/syntax | `godot --headless --check-only --script` | `dotnet build` (Godot.NET.Sdk) |
| Compile | n/a (parse is the static gate) | `dotnet build` |
| Runtime | headless script/scene run + stdout scan; GUT/GdUnit4 | needs .NET editor build; gdUnit4Net / headless run; xUnit/NUnit for pure logic |

### When to run Godot headless

Only when **(a)** a *complete* or *executable* example exists for the skill, **and (b)** a Godot binary is
available. Parse checks accept any Godot build; **C# runtime needs the .NET editor build**. Otherwise the check
is **SKIPPED** (binary absent) or **NOT_APPLICABLE** (no example of that kind).

## Status model (no false positives)

Per check:
- **PASS** — the check ran and succeeded.
- **FAIL** — the check ran and failed (syntax error, compile error, runtime error markers, non-zero framework exit).
- **SKIPPED** — the check *could* apply but a required tool/binary is absent; the report names the exact missing
  dependency (e.g. "no Godot binary", "no .NET SDK", "Godot build lacks C# support").
- **NOT_APPLICABLE** — the check does not apply by design (e.g. runtime test for an *illustrative* snippet; C#
  runtime where the skill has no C# example; C# **web export**, which is unsupported).

**Mandatory vs optional checks.** *Mandatory* = structural; parse of every *complete* GDScript example; build of
every *complete/executable* C# example; runtime of every *executable* project. *Optional/NA* = runtime of an
*illustrative* snippet, C# runtime where there is no C# example, C# **web export** (unsupported), etc.

Aggregation (per example/project, then global):
- **Per example/project:** the status is the worst of its checks (FAIL > SKIPPED > PASS; NOT_APPLICABLE ignored).
- **Global (per skill):** **FAIL** if any applicable check FAILed. **PASS (verified)** only if **every mandatory
  applicable check actually ran and PASSed**. If any **mandatory** check is **SKIPPED** (tool/binary absent), the
  global result is **SKIPPED / unverified — never a global PASS.** SKIPPED and NOT_APPLICABLE never count as PASS.
- `manifest.verification` records the global status, the timestamp, and the Godot version it ran against.

This is the core anti-false-positive rule: **a skipped mandatory test can never produce a green skill.** In CI,
the Linux engine tier (below) provides Godot 4.7.2 so these mandatory checks actually run rather than skip.

## Example verification report (deliverable — illustrative)

Human-readable:

```
Skill: signals            reference: Godot 4.7.2
  structural                        PASS
  examples/signal_bus.gd
    gdscript:parse                  PASS
    gdscript:runtime                SKIPPED  (no Godot binary; set GODOT_BIN)
  examples/csharp/SignalsDemo
    csharp:build                    PASS     (.NET 10.0.113, Godot.NET.Sdk 4.4.0)
    csharp:runtime                  SKIPPED  (needs .NET-enabled Godot build)
  web-export (C#)                   NOT_APPLICABLE  (C# cannot export to web)
  => status: unverified (skipped)   [no complete/executable example fully ran]
```

Machine-readable (proposed shape):

```json
{
  "skill": "signals",
  "reference": "4.7.2",
  "checks": [
    { "id": "structural", "status": "PASS" },
    { "example": "examples/signal_bus.gd", "id": "gdscript:parse", "status": "PASS" },
    { "example": "examples/signal_bus.gd", "id": "gdscript:runtime", "status": "SKIPPED", "reason": "no Godot binary" },
    { "example": "examples/csharp/SignalsDemo", "id": "csharp:build", "status": "PASS" },
    { "example": "examples/csharp/SignalsDemo", "id": "csharp:runtime", "status": "SKIPPED", "reason": "needs .NET-enabled Godot build" },
    { "id": "csharp:web-export", "status": "NOT_APPLICABLE", "reason": "unsupported" }
  ],
  "status": "unverified-skipped"
}
```

## Driver (resolved — gated, not implemented here)

**Resolved:** a **new, independent `toolkit verify` subcommand** reads each skill's manifest, runs the applicable
layers per example, and emits the report. **`toolkit validate` stays focused on structural + metadata validation**
(fast, always-on, no engine); `verify` owns the engine-level tiers. Implementing `verify` is a **CLI change →
gated behind approval**.

- Reuses: `skillspec` (structural), the `validate-headless.sh` patterns (GDScript parse/run + stdout scan),
  `dotnet build` (C#). New code is only the manifest reader, the per-example dispatcher, and the report emitter.
- **Reuse-first scope (resolved):** no new frameworks, dependencies or abstractions — Go stdlib + shelling to
  `godot`/`dotnet`. **No new third-party dependencies.**

## Gaps vs reuse (deliverable)

**Reusable now:** `skillspec` validator, `validate-headless.sh`, the proven `dotnet build` path, fixtures, CI
negative self-test, CI OS matrix.

**Real gaps (new work, each gated):**
1. **Manifest schema + validation** (small; extends the validator).
2. **A per-GDScript-file parse wrapper** (`--check-only --script`); current helper imports a whole project
   (small addition to `integrations/godot-cli`).
3. **The verify driver + report** (manifest reader, dispatcher, PASS/FAIL/SKIPPED/NA, JSON+text).
4. **C# runtime harness** (gdUnit4Net template) — later; build-only is enough for the pilot.
5. **CI engine tier (resolved):** add a **Linux job that installs Godot 4.7.2** and runs the engine-level tier
   (GDScript parse/run + C# build, and C# runtime where a .NET-enabled build is available). Windows/macOS engine
   tiers are a **possible later expansion**; until then they keep running structural + C# build only.

---

## Resolved decisions (approved)

7. **Driver:** a new independent **`toolkit verify`** subcommand; **`toolkit validate` stays structural +
   metadata** only. **Manifest-schema validation is metadata → it lives in `validate`**; `verify` consumes
   already-valid manifests and runs the engine tiers.
8. **CI:** add a **Linux engine tier with Godot 4.7.2**; other platforms are a later expansion.
9. **Results:** `PASS / FAIL / SKIPPED / NOT_APPLICABLE`, aggregated per example/project then global; **a skipped
   mandatory check never yields a global PASS** (see *Status model*).
10. **Scope:** reuse existing tools first; no unnecessary frameworks, dependencies or abstractions.

---

## Proposed P1 implementation scope (for approval — not done in this phase)

Minimal, reuse-first. In rough order:

1. **Manifest schema + validation in `validate`** — define the `manifest.json` schema (P0 §4) and validate it
   alongside `SKILL.md` (metadata check). Add the size-budget **warning** (≤500 lines / ~5k tokens, with the
   justified-exception note).
2. **`catalog.json` generation** — a `validate`/`verify` step that regenerates the aggregated catalog from
   manifests (never hand-edited).
3. **`toolkit verify` subcommand** — manifest reader → per-example dispatcher → report emitter
   (text + JSON), implementing the status model and aggregation. Go stdlib only.
4. **GDScript per-file parse helper** — extend `integrations/godot-cli` with a `--check-only --script` wrapper
   (the current helper imports whole projects); reuse its stdout error-scan.
5. **CI Linux engine tier** — a job that installs Godot 4.7.2 and runs `verify`'s engine tier; keep the existing
   matrix for structural + C# build.
6. *(Later, not in the first cut)* C# runtime harness (gdUnit4Net) — build-only is enough for the `signals` pilot.

### Files that would be modified/created during P1 implementation (for approval)
- `internal/skillspec/` (or a new `internal/manifest/`) — manifest schema + validation. *(tooling change)*
- `internal/cli/cli.go` — register the `verify` subcommand; `validate` gains manifest/size checks. *(tooling change)*
- `internal/verify/` (new package) — dispatcher, runners (GDScript parse/run, C# build/run), report model.
- `integrations/godot-cli/` — add the per-script parse helper.
- `.github/workflows/validate.yml` — add the Linux Godot 4.7.2 engine tier.
- `docs/authoring-skills.md` — fold in the approved standard (frontmatter, manifest, examples, citations).
- `CHANGELOG.md` — `[Unreleased]` entries.
- New fixtures under `tests/fixtures/` for `verify` (valid/invalid manifests; a tiny complete GDScript example and
  a minimal C# example project) — **additive test fixtures**, to be confirmed at approval.

> All of the above are **gated**: each is presented and approved before coding. None are touched in this
> consolidation phase.

### Tests worth running to validate the P1 implementation
- `go test ./...` — unit tests for the manifest parser/schema, the status aggregator (PASS/FAIL/SKIPPED/NA, and
  the "mandatory-skipped ⇒ not global PASS" rule), and the report emitter.
- `go run ./cmd/godot-ai-toolkit validate` — structural + manifest + size-warning checks on fixtures.
- `go run ./cmd/godot-ai-toolkit verify` — against valid/invalid manifest fixtures: assert correct statuses and
  that a broken example yields FAIL, a missing tool yields SKIPPED, and an N/A check is reported N/A.
- Engine tier (Linux CI / local with Godot 4.7.2): GDScript parse/run on the example; `dotnet build` on the C#
  example; confirm SKIPPED (not PASS) when the Godot binary/SDK is absent.
- Negative self-test: `verify` must **fail** on a fixture whose example has a deliberate syntax/compile error.

**No skills are authored and no tooling is implemented until this scope is approved.**

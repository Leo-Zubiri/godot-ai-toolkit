# Proposal P1 — Minimal Verification Workflow

> **Status: PROPOSAL — pending approval.** Design/documentation only. It proposes a workflow and, where tooling
> would be needed, **flags it as a gated change** (no installer/CLI/adapter/fixture code is modified in this
> phase). Companion to [`authoring-and-quality-standard.md`](authoring-and-quality-standard.md) (P0). Open
> decisions marked **[DECISION]**.

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
| CI matrix (Win/Linux/macOS) | `.github/workflows/validate.yml` | **Reuse**; add an engine tier — **[DECISION]**. |

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

Aggregation:
- A skill's `verification.status` = **FAIL** if any applicable check FAILed; **verified** only if **all
  applicable checks PASS and at least one complete/executable example actually ran** (parsed/compiled, and ran
  where an executable example exists). If applicable checks were only SKIPPED, status = **unverified (skipped)** —
  **never** "verified". SKIPPED/NOT_APPLICABLE never count as PASS.

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

## Proposed driver (gated — not implemented here)

A single entry point reads each skill's manifest, runs the applicable layers per example, and emits the report.

- **[DECISION 6]** Implement as a **new `toolkit verify` subcommand** (recommended — keeps `validate` fast and
  always-on) vs extending `toolkit validate`. Either is a **CLI change → gated behind approval** (P1 build step).
- Reuses: `skillspec` (structural), `validate-headless.sh` patterns (GDScript parse/run + stdout scan),
  `dotnet build` (C#). New code is the manifest reader, the per-example dispatcher, and the report emitter.
- **No new third-party dependencies** anticipated (Go stdlib + shelling to `godot`/`dotnet`).

## Gaps vs reuse (deliverable)

**Reusable now:** `skillspec` validator, `validate-headless.sh`, the proven `dotnet build` path, fixtures, CI
negative self-test, CI OS matrix.

**Real gaps (new work, each gated):**
1. **Manifest schema + validation** (small; extends the validator).
2. **A per-GDScript-file parse wrapper** (`--check-only --script`); current helper imports a whole project
   (small addition to `integrations/godot-cli`).
3. **The verify driver + report** (manifest reader, dispatcher, PASS/FAIL/SKIPPED/NA, JSON+text).
4. **C# runtime harness** (gdUnit4Net template) — later; build-only is enough for the pilot.
5. **CI engine tier** — provisioning Godot in CI has a cost; **[DECISION 7]** run engine checks in CI (install
   Godot) vs keep them local/opt-in and have CI run structural + C# build only.

---

## Open decisions (consolidated)

- **[DECISION 6]** `toolkit verify` (new subcommand, recommended) vs extend `validate`.
- **[DECISION 7]** Engine-level checks in CI (install Godot) vs local/opt-in; CI keeps structural + C# build.
- Manifest-schema validation location (validator vs verify driver).
- Example project layout for C#: per-skill project (cleaner, heavier restore) vs a shared examples project.

Presented together with the P0 decisions. **No skills are authored and no tooling is implemented until these are
approved.**

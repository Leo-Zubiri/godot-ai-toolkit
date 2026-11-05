# Skill research & technical-grounding protocol

**Mandatory before authoring any skill.** This protocol governs *how a skill is grounded and validated*;
[`authoring-skills.md`](authoring-skills.md) governs the *format/contract* and
[`proposals/verification-workflow.md`](proposals/verification-workflow.md) the *verification mechanics*.
A skill is not complete merely because its documentation is accurate — it must also be **discoverable,
understandable and useful to an AI coding agent**, and every non-obvious claim must be traceable to an official
source.

> Do not fabricate research results. Do not claim to have consulted a source you did not actually inspect.

## 1. Official Godot documentation — primary technical authority

Use the official, **version-pinned** docs (`docs.godotengine.org/en/<version>/`), never `stable`, and never
search snippets as technical evidence:

- Godot 4.7 documentation and **class reference** (reference version).
- Version-specific docs for the declared **minimum supported** version when behavior differs.
- Official migration guides / release notes when version differences matter.

For every relevant topic: (1) identify the doc page, (2) identify the classes/methods/properties/concepts,
(3) verify signatures and behavior against the class reference, (4) check version compatibility, (5) note
documented limitations and lifecycle considerations, (6) record the **source URL + applicable Godot version**.
Do not invent API behavior from familiarity with another version. If the docs are ambiguous or silent, say so
explicitly instead of presenting an assumption as fact.

## 2. GodotPrompter — secondary skill-engineering reference

Study the repo's existing research ([`research/existing-skills-ecosystem.md`](research/existing-skills-ecosystem.md),
pinned revision) for engineering *patterns* only — organization, boundaries, descriptions/discoverability,
instruction writing, validation, negative tests, evaluation, agent-oriented examples.

- **Do not** copy its `SKILL.md` content, examples or prose; **do not** treat it as authoritative for Godot APIs;
  **do not** assume its architectural choices fit this toolkit. Respect its license and the P0 licensing rules.
- Use the **pinned revision**; do not silently substitute another.
- **Record** which patterns were considered, adopted, adapted or rejected, and why (in the research report).

## 3. Technical evidence & traceability

Every skill must clearly distinguish: **documented behavior** · **version-specific behavior** · **toolkit
recommendations** · **architectural opinions/trade-offs** · **known limitations / open questions**.

- Put a source reference near every non-obvious technical claim.
- `## Sources` in `SKILL.md` lists the **relevant official sources**, not a generic link dump.
- The manifest's `sources[]` carries structured source metadata (claim → URL → version).

## 4. Practical validation

Design representative tasks an AI agent could realistically receive, with examples covering correct
implementation, common mistakes, relevant edge cases, and version compatibility where applicable. Use the
existing verification system and **distinguish the tiers**:

- **Structural** (frontmatter contract, manifest schema, reference integrity).
- **Parse/syntax** (GDScript `--check-only`).
- **Build** (C# `dotnet build` → `BUILD_VERIFIED`).
- **Runtime** (headless execution that asserts behavior → `RUNTIME_VERIFIED`).

**Never claim runtime correctness from parsing or compilation alone.** A skipped mandatory check never yields
a global PASS.

## 5. Agent-usefulness evaluation

For each skill define at least: one task that **should trigger** it; one that **should not**; the expected
technical decisions; the common **failure modes** it prevents; and **observable acceptance criteria**.

## 6. Research deliverable before implementation

Before writing the skill, present a short **research report** containing:

1. Official Godot sources consulted (URLs + versions).
2. Relevant APIs and version constraints (with facts-vs-flags).
3. GodotPrompter patterns examined.
4. Patterns adopted / adapted / rejected and why.
5. Proposed skill boundaries.
6. Technical risks and ambiguities.
7. Verification plan.
8. Agent evaluation scenarios.

Only after this research is complete and approved does implementation begin.

---

*Reference example:* the `signals` and `scenes-and-nodes` skills were authored under this protocol; their
research reports and the ecosystem study live under [`research/`](research/) and the plan history.

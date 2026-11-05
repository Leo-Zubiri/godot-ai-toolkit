# skills/ — the official Godot knowledge catalog

This directory is the **single source of truth** for the toolkit's knowledge. Each skill is one folder
with a `SKILL.md` following the [Agent Skills spec](https://agentskills.io/specification). See
[`../docs/authoring-skills.md`](../docs/authoring-skills.md).

## Status

**Authored:** [`signals`](signals/SKILL.md) — the first production skill (Godot 4.7, GDScript + C#).
The rest of the taxonomy below is planned. `catalog.json` is generated from the skills' manifests.

> The throwaway example skill used to exercise the installer/validator lives in
> [`../tests/fixtures/skills/`](../tests/fixtures/), **never here**. `skills/` is reserved for the
> real, published catalog.

## Grouping philosophy

Do **not** create one skill per bullet below. Over-fragmentation is an anti-goal (README principle 4,
6). Group related concepts into a coherent skill; split only when a topic is large enough to stand on
its own and is used independently. Shared Godot concepts are written once; language specifics stay
separate (README principle 3–5).

## Planned taxonomy (authored in later versions)

**Core Godot concepts (language-neutral, with GDScript + C# examples where the API differs):**

| Planned skill | Covers |
|---|---|
| `godot-fundamentals` | Engine model, project structure, the editor, first principles. |
| `scenes-and-nodes` | Scene/node architecture, instancing, node lifecycle. |
| `signals` | Signal declaration, connection, patterns; GDScript vs C# (`[Signal]`, `ToSignal`). |
| `resources` | `Resource`/custom resources, loading, `.tres`/`.res`. |
| `ui-ux` | Control nodes, anchors/containers, themes, input focus. |
| `animation` | AnimationPlayer/Tree, tweens, state. |
| `physics` | Bodies, areas, collision layers, `_physics_process`. |
| `shaders` | Shading language, canvas/spatial shaders, uniforms. |
| `audio` | Buses, streams, playback, positional audio. |
| `navigation-ai` | Navigation meshes/agents, pathfinding, basic game AI. |
| `multiplayer` | High-level multiplayer, RPCs, synchronization. |
| `persistence` | Save/load, `user://`, serialization. |
| `debugging` | Debugger, logging, profiler basics. |
| `testing` | GUT and GdUnit4; headless CI; what each supports. |
| `performance` | Profiling, common costs, draw calls, GC/alloc notes. |
| `exporting` | Export presets/templates; platform/language constraints. |

**Language-specific skills (only what is truly language-bound):**

| Planned skill | Covers |
|---|---|
| `gdscript` | Style/typing, script organization, patterns, engine integration, debug/test. |
| `csharp-dotnet` | Project/`.csproj` setup, .NET editor build, API mapping, async, NuGet, build/debug/test, export limits. |

Each authored skill must include: when to use it, practical instructions, correct examples, common
pitfalls, references, and explicit compatibility notes (verified against official docs).

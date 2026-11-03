# Compatibility matrix

Verified against official documentation. **Every version/platform claim here must be sourced.** Do not
add a requirement you have not confirmed in official docs, and do not promise a platform Godot does not
support for a given language. Re-verify point-release numbers before pinning — they move fast.

_Last verified: 2026-11-02. Current stable Godot: **4.7.2**._

## Godot ↔ language ↔ .NET ↔ platforms

| Godot | Language | .NET SDK required | Desktop (Win/mac/Linux) | Web / HTML5 | Android | iOS |
|---|---|---|---|---|---|---|
| 4.7 (stable) | GDScript | n/a | ✅ stable | ✅ stable | ✅ | ✅ |
| 4.7 (stable) | C#/.NET | **.NET 8+** (**.NET 9+** for Android) | ✅ stable (since 4.2) | ❌ **not supported** | ⚠️ experimental | ⚠️ experimental (macOS host only) |
| 4.5 | C#/.NET | .NET 8+ (.NET 9+ for Android) | ✅ | ❌ | ⚠️ | ⚠️ |
| 4.4 | C#/.NET | .NET 8 | ✅ | ❌ | ⚠️ | ⚠️ |
| 4.3 | C#/.NET | .NET 6 | ✅ | ❌ | ⚠️ | ⚠️ |

Legend: ✅ supported/stable · ⚠️ experimental/limited · ❌ not supported.

### Key facts (with sources)

- **.NET SDK floor** is **.NET 8** for current stable Godot; **Android export requires .NET 9+**.
  Progression across releases: 4.3 → .NET 6, 4.4 → .NET 8, 4.5+ → .NET 8 (9 for Android).
  — [C# Basics](https://docs.godotengine.org/en/stable/tutorials/scripting/c_sharp/c_sharp_basics.html)
- **C# cannot be exported to the Web platform.** If web is a hard requirement, that path must stay
  GDScript (or use Godot 3 for C# on web).
  — [C# index](https://docs.godotengine.org/en/stable/tutorials/scripting/c_sharp/index.html)
- **Android and iOS C# support are experimental.** iOS can only be exported from macOS; simulator
  templates are x64-only.
  — [C# index](https://docs.godotengine.org/en/stable/tutorials/scripting/c_sharp/index.html),
  [Platform state in C# (4.2)](https://godotengine.org/article/platform-state-in-csharp-for-godot-4-2/)
- **C# desktop export is stable** on Windows, Linux and macOS (since Godot 4.2).
- **C# requires the `.NET` editor build.** The standard Godot editor cannot build or run C#.
  — [C# index](https://docs.godotengine.org/en/stable/tutorials/scripting/c_sharp/index.html)
- **Godot does not bundle the C# build tools** (MSBuild, C# compiler); they come with the .NET SDK,
  installed separately. Godot does bundle the runtime parts to *run* compiled games.
  — [C# Basics](https://docs.godotengine.org/en/stable/tutorials/scripting/c_sharp/c_sharp_basics.html)

## GDScript vs C# — documented differences that affect the toolkit

| Topic | GDScript | C#/.NET | Source |
|---|---|---|---|
| Editor support | Full built-in editor + IntelliSense | **No built-in IntelliSense**; needs an external IDE | [C# differences](https://docs.godotengine.org/en/stable/tutorials/scripting/c_sharp/c_sharp_differences.html) |
| Web export | Supported | **Not supported** | C# index |
| Editor build | Standard build | **.NET build required** | C# index |
| Naming/API | snake_case, `@export`, `signal` | PascalCase, `[Export]`, `[Signal]`; globals under `GD`, `Mathf`, `GodotObject` | C# differences |
| Lifecycle helpers | `@onready`, `preload` | No `@onready` (assign in `_Ready()`); no `preload` (use `GD.Load`/`ResourceLoader.Load`) | C# differences |
| Async signals | coroutines/`await` | `await ToSignal(...)` | C# differences |
| Strings/paths | UTF-32; understands `res://`/`user://` | UTF-16; `System.IO.Path` does **not** understand `res://`/`user://` | C# differences |

> Do not assume parity (README principle 9). When a skill states a capability, it must also state the
> language/build/platform constraints that apply.

## Testing frameworks (verified)

| Framework | Language | Headless CI | Notes |
|---|---|---|---|
| [GUT](https://github.com/bitwes/Gut) | GDScript only | ✅ (`-gexit`, JUnit XML) | 9.x line = Godot 4; 7.x = Godot 3. |
| [GdUnit4](https://github.com/MikeSchulze/gdUnit4) | **GDScript + C#** | ✅ (runner + GitHub Action, JUnit XML) | C# via `gdUnit4.api` + `gdUnit4.test.adapter` (VSTest). |
| xUnit / NUnit (`dotnet test`) | C# (pure logic) | ✅ | For code that does **not** touch engine types; no Godot needed. |

- **Engine-dependent C# tests require launching the Godot runtime** (.NET editor build). `gdUnit4Net`
  v5 runs without the runtime by default; opt in with `[RequireGodotRuntime]` — those are *skipped* if
  the Godot binary lacks C# support.

## Godot CLI for automation (verified)

[Command-line tutorial](https://docs.godotengine.org/en/stable/tutorials/editor/command_line_tutorial.html)

| Flag | Meaning |
|---|---|
| `--headless` | `--display-driver headless --audio-driver Dummy`; required on machines without GPU (e.g. CI). |
| `--import` | Import resources, then quit (implies `--editor --quit`). |
| `--export-release <preset> <path>` | Export with a preset. `--export-debug` uses the debug template. |
| `--script <res://…>` | Run a script (path relative to the project). |
| `--check-only --script <s>` | Parse-check a script and quit. |
| `--path <dir>` | Project directory (must contain `project.godot`). |
| `--quit` / `--quit-after <n>` | Quit after one / N iterations. |

> **Exit-code caveat (verified):** headless Godot can **exit 0 even on failure** (e.g. import/run with
> a script that fails to parse). CI must **also scan stdout/stderr** for `SCRIPT ERROR`, `ERROR`,
> `Parse Error`, and should gate on the test framework's own runner, which exits non-zero on failure.
> See [`../integrations/godot-cli/`](../integrations/godot-cli/).

## Toolkit tooling requirements (separate from consumer requirements)

| Layer | Requirement |
|---|---|
| Toolkit **development** | Go toolchain (build/test the CLI); optionally Godot + .NET SDK for fixture checks. |
| Installer **implementation** | Go (standard library only). |
| Installer **distribution** | Native binaries per OS/arch + checksums. |
| **End-user runtime** | **None.** No Go/Python/Node/uv required to use the CLI. |

# integrations/

Optional **capabilities** an agent can *use* while working on a Godot project — as opposed to skills
(knowledge) or adapters (placement data). Integrations are optional by design: the toolkit never
requires them, and we do not add a tool just because it exists (README principle 6).

Distinction (see [`../docs/architecture.md`](../docs/architecture.md)):

- **Skills** = instructions/procedures for agents.
- **MCP** = tools/capabilities for interacting with Godot.
- **CLI** = install/update/validate (the Go tooling).
- **Adapters** = per-agent discovery descriptors.

## MCP for Godot — status: experimental, **deferred**

As of 2026-11, there is **no official** (Godot Foundation or Anthropic) MCP server for Godot. The
ecosystem is community-driven and uneven in release discipline. We therefore **document** it and let
developers opt in manually; we do **not** bundle or auto-configure a server.

| Server | What it does | Maturity |
|---|---|---|
| [Coding-Solo/godot-mcp](https://github.com/Coding-Solo/godot-mcp) | CLI-driven: launch editor, run/stop projects, capture debug output, create scenes/nodes. Standalone. | Most adopted; still experimental (no tagged releases). |
| [hybridindie/godot-mcp](https://github.com/hybridindie/godot-mcp) | Editor plugin + bridge; deep *live* editor control. | Feature-rich; less battle-tested. |
| [godot-mcp-docs](https://pypi.org/project/godot-mcp-docs/) | Version-aware access to the Godot class reference for agents. | Docs-only; useful and low-risk. |

**Source-of-truth alternative without any MCP server:** Godot ships its class reference as **XML** in
the engine source and can emit/merge it via `godot --doctool`. Feeding an agent that XML is a robust,
official way to give it accurate API facts.

When a demonstrable need appears, a dedicated `integrations/mcp/` folder can hold example client
configuration (e.g. for `godot-mcp-docs`). It is intentionally absent until then.

## Godot CLI helpers

See [`godot-cli/`](godot-cli/). Small, dependency-free scripts for headless validation (parse/import)
suitable for CI. They implement the verified caveat that **headless Godot can exit 0 on failure**, so
they also scan output for error markers.

## .NET helpers

See [`dotnet/`](dotnet/). Notes and an example `.runsettings` for testing C# Godot code, plus the
split between pure-logic tests (`dotnet test` with xUnit/NUnit) and engine-bound tests (`gdUnit4Net`,
which needs the .NET editor build).

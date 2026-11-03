# tests/

Fixtures and the validation strategy. **Go unit tests live beside the code** under `internal/`
(run with `go test ./...`); this directory holds the shared, deterministic *fixtures* those tests and
CI consume.

## Fixtures

```
fixtures/
├── skills/
│   ├── valid/example-skill/      # spec-valid SKILL.md (+ a reference) — must pass
│   └── invalid/                  # each must be REJECTED by the validator:
│       ├── name-mismatch/        #   name != folder
│       ├── missing-description/  #   required field absent
│       ├── bad-name/             #   name violates [a-z0-9-] and folder match
│       ├── no-frontmatter/       #   no YAML fence
│       └── broken-reference/     #   body links to a missing bundled file
├── gdscript-project/             # minimal Godot 4.x GDScript project
├── csharp-project/               # minimal Godot 4.x C# project (Godot.NET.Sdk)
└── mixed-project/                # GDScript + C# in one project
```

> The example skill lives here **only** as test infrastructure. It is not a real skill and never
> appears under `skills/` (reserved for the official catalog).

## Validation tiers (honest about what runs where)

1. **Always runnable (no Godot/.NET).** `go test ./...` and `godot-ai-toolkit validate`:
   - `SKILL.md` structure + frontmatter (folder == name, description length, required fields).
   - Reference integrity (relative links resolve).
   - Adapter descriptor and lock-file schema.
   - End-to-end pipeline against `fixtures/skills/` — the **valid** skill installs to each agent path,
     produces a correct lock, and survives update/conflict/backup flows; **invalid** skills are rejected.

2. **SDK-gated.** `dotnet build tests/fixtures/csharp-project` when a .NET 8+ SDK is present (restores
   `Godot.NET.Sdk` from NuGet). If the SDK or packages are unavailable, CI reports **skipped**.

3. **Godot-gated.** `integrations/godot-cli/validate-headless.sh tests/fixtures/<project>` when a Godot
   binary is available. It scans stdout for error markers (headless Godot can exit 0 on failure). If no
   Godot binary is present, it reports **skipped**.

A skipped tier is never reported as passed.

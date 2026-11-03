# .NET helpers

Notes and an example configuration for building and testing C# Godot code. See the verified facts in
[`../../docs/compatibility.md`](../../docs/compatibility.md).

## Requirements (verified)

- **.NET 8 SDK** floor for current stable Godot; **.NET 9 SDK** if you export to **Android**.
- C# requires the **.NET editor build** of Godot; the standard build cannot build/run C#.
- Godot does not bundle MSBuild / the C# compiler — they come with the .NET SDK.

## Two testing tiers

1. **Pure logic (no engine types)** — test with `dotnet test` using **xUnit** or **NUnit**. No Godot
   required; fast and CI-friendly. Use for anything that does not touch `Node`, `Resource`, scenes, etc.

2. **Engine-dependent code** — needs the Godot runtime (.NET editor build). Use
   [`gdUnit4Net`](https://github.com/MikeSchulze/gdUnit4): NuGet `gdUnit4.api` (+ `gdUnit4.test.adapter`
   for VSTest/IDE). In v5, tests run **without** the runtime by default; opt in per test with
   `[RequireGodotRuntime]`. Those tests are **skipped** if the Godot binary lacks C# support.

## Example `.runsettings`

[`gdunit4.runsettings`](gdunit4.runsettings) shows the knobs gdUnit4Net reads (e.g. the Godot
connection timeout). Point your runner at it:

```sh
dotnet test --settings integrations/dotnet/gdunit4.runsettings
```

## Build check (used by CI when the SDK is present)

```sh
dotnet build tests/fixtures/csharp-project
```

If no suitable SDK is installed, CI reports this as **skipped**, never as passed.

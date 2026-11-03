# AGENTS.md — operating rules for AI agents in this repository

This is the short, operational companion to [`README.md`](README.md). The **README is the source of
truth** for the project's vision and its non-negotiable principles; this file tells you how to act.
Do not restate the README here.

## Before you change anything significant

Significant = architecture, adding/removing skills, changing the CLI contract, the lock-file or
adapter schema, or the distribution model.

1. **Read [`README.md`](README.md)** and confirm your change respects its purpose and the 14
   non-negotiable principles (README §8).
2. If your change conflicts with a principle, **stop and flag it explicitly**, with justification.
   Do not quietly work around a principle.
3. **Verify every Godot / .NET fact against official documentation** before writing it. Never invent
   APIs, methods, properties, version requirements or platform support. Prefer the official docs over
   blogs, third-party repos, or memory. Cite the source in the PR or the file.
4. Check [`docs/compatibility.md`](docs/compatibility.md) for verified version/platform facts, and
   do not contradict it without updating it (with sources).

## How this repository is organized (where things go)

- **Skills are authored once** under `skills/`, following the Agent Skills `SKILL.md` spec. The folder
  name must equal the `name` frontmatter field. Never duplicate skill content per agent.
- **Per-agent differences live in `adapters/*.json`** as *data* (discovery path, instruction file).
  Never hard-code agent-specific logic into the installer core.
- **The CLI (`cmd/`, `internal/`) is tooling only** and is written in Go, standard-library first.
  The toolkit's *content* (skills, docs, metadata) stays independent of Go.
- **Fixtures live in `tests/fixtures/`.** The throwaway example skill belongs there, never in
  `skills/` (which is reserved for the official catalog).

## Do

- Keep it simple. Prefer the smallest change that works.
- Keep GDScript and C# concerns correctly separated; never assume parity between them.
- Update `CHANGELOG.md` (`[Unreleased]`) and run `go test ./...` and
  `go run ./cmd/godot-ai-toolkit validate` before proposing changes.
- Report honestly when a validation could not run (e.g. no Godot binary, no .NET SDK). Never claim a
  check passed if it did not execute.

## Don't

- Don't auto-update consumer projects or introduce background update behavior.
- Don't add dependencies the standard library can cover, speculative features, or empty directories.
- Don't couple the toolkit to a single AI provider.
- Don't invent engine behavior or promise unsupported platforms (e.g. C# → Web export is **not**
  supported).

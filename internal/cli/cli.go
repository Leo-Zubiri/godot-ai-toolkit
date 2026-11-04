// Package cli parses arguments and dispatches the toolkit subcommands.
package cli

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/godot-ai-toolkit/godot-ai-toolkit/internal/adapter"
	"github.com/godot-ai-toolkit/godot-ai-toolkit/internal/catalog"
	"github.com/godot-ai-toolkit/godot-ai-toolkit/internal/fetch"
	"github.com/godot-ai-toolkit/godot-ai-toolkit/internal/fsutil"
	"github.com/godot-ai-toolkit/godot-ai-toolkit/internal/installer"
	"github.com/godot-ai-toolkit/godot-ai-toolkit/internal/manifest"
	"github.com/godot-ai-toolkit/godot-ai-toolkit/internal/skillspec"
	"github.com/godot-ai-toolkit/godot-ai-toolkit/internal/verify"
)

// Version is the CLI version (override at build time with -ldflags).
var Version = "0.1.0-dev"

// Run dispatches a subcommand and returns a process exit code.
func Run(args []string) int {
	if len(args) == 0 {
		printUsage(os.Stderr)
		return 2
	}
	switch args[0] {
	case "install":
		return cmdInstall(args[1:])
	case "update":
		return cmdUpdate(args[1:])
	case "status":
		return cmdStatus(args[1:])
	case "validate":
		return cmdValidate(args[1:])
	case "catalog":
		return cmdCatalog(args[1:])
	case "verify":
		return cmdVerify(args[1:])
	case "version", "--version", "-v":
		fmt.Println("godot-ai-toolkit " + Version)
		return 0
	case "help", "-h", "--help":
		printUsage(os.Stdout)
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", args[0])
		printUsage(os.Stderr)
		return 2
	}
}

func cmdInstall(args []string) int {
	fs := flag.NewFlagSet("install", flag.ContinueOnError)
	source := fs.String("source", "", "local toolkit checkout to install from")
	project := fs.String("project", ".", "consumer project directory")
	agents := fs.String("agents", "", "comma-separated agents (e.g. claude,codex)")
	skills := fs.String("skills", "", "comma-separated skill names")
	tag := fs.String("tag", "", "toolkit version tag to record (and fetch, with --repo)")
	repo := fs.String("repo", "", "toolkit repository URL (for fetch and lock)")
	commit := fs.String("commit", "", "source commit to record")
	url := fs.String("url", "", "explicit content tarball URL to fetch")
	sha := fs.String("sha256", "", "expected SHA-256 of the fetched tarball")
	force := fs.Bool("force", false, "override conflict protection (backs up first)")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	src, cleanup, err := resolveSource(*source, *url, *sha, *repo, *tag)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}
	if cleanup != nil {
		defer cleanup()
	}
	rep, err := installer.Install(installer.Options{
		SourceDir:  src,
		ProjectDir: *project,
		Agents:     splitCSV(*agents),
		Skills:     splitCSV(*skills),
		Version:    *tag,
		Commit:     *commit,
		Repository: *repo,
		Force:      *force,
	})
	return reportApply("installed", rep, err)
}

func cmdUpdate(args []string) int {
	fs := flag.NewFlagSet("update", flag.ContinueOnError)
	source := fs.String("source", "", "local toolkit checkout to update from")
	project := fs.String("project", ".", "consumer project directory")
	agents := fs.String("agents", "", "comma-separated agents (default: keep current)")
	skills := fs.String("skills", "", "comma-separated skill names (default: keep current)")
	tag := fs.String("tag", "", "toolkit version tag to record (and fetch, with --repo)")
	repo := fs.String("repo", "", "toolkit repository URL")
	commit := fs.String("commit", "", "source commit to record")
	url := fs.String("url", "", "explicit content tarball URL to fetch")
	sha := fs.String("sha256", "", "expected SHA-256 of the fetched tarball")
	force := fs.Bool("force", false, "override conflict protection (backs up first)")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	src, cleanup, err := resolveSource(*source, *url, *sha, *repo, *tag)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}
	if cleanup != nil {
		defer cleanup()
	}
	rep, err := installer.Update(installer.Options{
		SourceDir:  src,
		ProjectDir: *project,
		Agents:     splitCSV(*agents),
		Skills:     splitCSV(*skills),
		Version:    *tag,
		Commit:     *commit,
		Repository: *repo,
		Force:      *force,
	})
	return reportApply("updated", rep, err)
}

func cmdStatus(args []string) int {
	fs := flag.NewFlagSet("status", flag.ContinueOnError)
	project := fs.String("project", ".", "consumer project directory")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	rep, err := installer.Status(*project)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}
	l := rep.Lock
	fmt.Printf("toolkit:    %s %s\n", orDash(l.Toolkit.Repository), orDash(l.Toolkit.Version))
	if l.Toolkit.Commit != "" {
		fmt.Printf("commit:     %s\n", l.Toolkit.Commit)
	}
	fmt.Printf("agents:     %s\n", strings.Join(l.Agents, ", "))
	fmt.Printf("skills:     %d\n", len(l.Skills))
	fmt.Printf("installed:  %s\n", orDash(l.InstalledAt))
	fmt.Println("files:")
	modified := 0
	for _, f := range rep.Files {
		if f.State != "ok" {
			modified++
		}
		fmt.Printf("  [%-8s] %s\n", f.State, f.Path)
	}
	if modified > 0 {
		fmt.Fprintf(os.Stderr, "\n%d managed file(s) differ from the lock (modified/missing).\n", modified)
		return 1
	}
	fmt.Println("\nall managed files match the lock.")
	return 0
}

func cmdValidate(args []string) int {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	repoDir := fs.String("repo-dir", ".", "toolkit repository root")
	skillsDir := fs.String("skills-dir", "", "skills directory to validate (default <repo-dir>/skills)")
	adaptersDir := fs.String("adapters-dir", "", "adapters directory to validate (default <repo-dir>/adapters)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	sd := *skillsDir
	if sd == "" {
		sd = filepath.Join(*repoDir, "skills")
	}
	ad := *adaptersDir
	if ad == "" {
		ad = filepath.Join(*repoDir, "adapters")
	}

	ok := true
	if fsutil.Exists(ad) {
		if _, err := adapter.LoadAll(ad); err != nil {
			fmt.Fprintln(os.Stderr, "FAIL adapters:", err)
			ok = false
		} else {
			fmt.Println("ok   adapters")
		}
	}

	count := 0
	entries, err := os.ReadDir(sd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read skills directory %s: %v\n", sd, err)
		return 1
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		dir := filepath.Join(sd, e.Name())
		skillMD := filepath.Join(dir, "SKILL.md")
		if !fsutil.Exists(skillMD) {
			continue
		}
		count++
		skillOK := true

		// 1. SKILL.md frontmatter contract, structure and reference integrity.
		for _, is := range skillspec.ValidateDir(dir) {
			skillOK = false
			fmt.Fprintf(os.Stderr, "FAIL %s\n", is)
		}

		// 2. manifest.json (required; the source of truth for metadata).
		mpath := filepath.Join(dir, manifest.FileName)
		if !fsutil.Exists(mpath) {
			skillOK = false
			fmt.Fprintf(os.Stderr, "FAIL %s: %s is required\n", e.Name(), manifest.FileName)
		} else if m, err := manifest.Load(mpath); err != nil {
			skillOK = false
			fmt.Fprintf(os.Stderr, "FAIL %s/%s: %v\n", e.Name(), manifest.FileName, err)
		} else {
			for _, me := range m.Validate(e.Name(), dir) {
				skillOK = false
				fmt.Fprintf(os.Stderr, "FAIL %s/%s: %s\n", e.Name(), manifest.FileName, me)
			}
			// 3. Size budget (warning, never fatal).
			if w := sizeWarning(skillMD, m.SizeExceptionReason); w != "" {
				fmt.Fprintf(os.Stderr, "WARN %s: %s\n", e.Name(), w)
			}
		}

		if skillOK {
			fmt.Printf("ok   %s\n", e.Name())
		} else {
			ok = false
		}
	}
	if count == 0 {
		fmt.Printf("no skills found in %s (nothing to validate)\n", sd)
	}
	if !ok {
		fmt.Fprintln(os.Stderr, "\nvalidation failed")
		return 1
	}
	fmt.Println("\nvalidation passed")
	return 0
}

func cmdCatalog(args []string) int {
	fs := flag.NewFlagSet("catalog", flag.ContinueOnError)
	skillsDir := fs.String("skills-dir", "skills", "skills directory to index")
	write := fs.Bool("write", false, "regenerate and write catalog.json")
	_ = fs.Bool("check", false, "check that catalog.json is up to date (the default when --write is absent)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	cat, err := catalog.Build(*skillsDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error building catalog:", err)
		return 1
	}
	data, err := cat.Marshal()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}
	path := catalog.Path(*skillsDir)
	if *write {
		if err := fsutil.WriteFileAtomic(path, data, 0o644); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return 1
		}
		fmt.Printf("wrote %s (%d skills)\n", path, len(cat.Skills))
		return 0
	}
	existing, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s is missing or unreadable — run 'catalog --write'\n", path)
		return 1
	}
	if !bytes.Equal(existing, data) {
		fmt.Fprintf(os.Stderr, "%s is out of date — run 'catalog --write'\n", path)
		return 1
	}
	fmt.Printf("%s up to date (%d skills)\n", path, len(cat.Skills))
	return 0
}

func cmdVerify(args []string) int {
	fs := flag.NewFlagSet("verify", flag.ContinueOnError)
	skillsDir := fs.String("skills-dir", "skills", "skills directory to verify")
	asJSON := fs.Bool("json", false, "emit the report as JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	entries, err := os.ReadDir(*skillsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read skills directory %s: %v\n", *skillsDir, err)
		return 1
	}
	runners := verify.DefaultRunners()
	var reports []verify.SkillReport
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		dir := filepath.Join(*skillsDir, e.Name())
		mpath := filepath.Join(dir, manifest.FileName)
		if !fsutil.Exists(mpath) {
			continue
		}
		m, err := manifest.Load(mpath)
		if err != nil {
			reports = append(reports, verify.SkillReport{
				Skill:  e.Name(),
				Checks: []verify.Check{{ID: "structural", Status: verify.FAIL, Mandatory: true, Reason: err.Error()}},
				Status: verify.FAIL,
			})
			continue
		}
		reports = append(reports, verify.VerifySkill(dir, m, runners))
	}

	overall := verify.Aggregate(reports)
	report := verify.Report{Skills: reports, Status: overall}

	if *asJSON {
		data, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(data))
	} else {
		printVerifyText(report, len(reports) == 0, *skillsDir)
	}

	switch overall {
	case verify.FAIL:
		return 1
	case verify.SKIP:
		return 2 // incomplete: a mandatory check could not run (never a PASS)
	default:
		return 0
	}
}

func printVerifyText(report verify.Report, empty bool, skillsDir string) {
	if empty {
		fmt.Printf("no skills with a manifest found in %s (nothing to verify)\n", skillsDir)
		return
	}
	for _, s := range report.Skills {
		fmt.Printf("Skill: %-24s reference: %s\n", s.Skill, orDash(s.Reference))
		for _, c := range s.Checks {
			id := c.ID
			if c.Example != "" {
				id = c.Example + "  " + c.ID
			}
			line := fmt.Sprintf("  %-40s %s", id, c.Status)
			if c.Verified != "" {
				line += "  (" + c.Verified + "-verified)"
			}
			if c.Reason != "" {
				line += "  — " + c.Reason
			}
			fmt.Println(line)
		}
		fmt.Printf("  => %s\n\n", s.Status)
	}
	fmt.Printf("overall: %s\n", report.Status)
}

// resolveSource returns a local source directory. If source is empty, it fetches
// a content tarball (from --url, or --repo + --tag) into a temp dir.
func resolveSource(source, url, sha, repo, tag string) (string, func(), error) {
	if source != "" {
		return source, nil, nil
	}
	fetchURL := url
	if fetchURL == "" && repo != "" && tag != "" {
		fetchURL = fetch.ReleaseTarballURL(repo, tag)
	}
	if fetchURL == "" {
		return "", nil, fmt.Errorf("provide --source <dir>, or --url (with optional --sha256), or --repo with --tag")
	}
	tmp, err := os.MkdirTemp("", "gait-src-")
	if err != nil {
		return "", nil, err
	}
	cleanup := func() { os.RemoveAll(tmp) }
	if err := fetch.DownloadAndExtract(fetchURL, sha, tmp); err != nil {
		cleanup()
		return "", nil, err
	}
	return descendSingleDir(tmp), cleanup, nil
}

// sizeWarning returns a non-empty warning when a SKILL.md exceeds the size
// target (≤500 lines / ~5k tokens). It is advisory, never fatal; a recorded
// exception reason is acknowledged rather than suppressing the notice.
func sizeWarning(skillMD, exceptionReason string) string {
	data, err := os.ReadFile(skillMD)
	if err != nil {
		return ""
	}
	lines := strings.Count(string(data), "\n") + 1
	approxTokens := len(data) / 4 // rough heuristic
	if lines <= 500 && approxTokens <= 5000 {
		return ""
	}
	if exceptionReason != "" {
		return fmt.Sprintf("large SKILL.md (%d lines, ~%d tokens) — exception noted: %s",
			lines, approxTokens, exceptionReason)
	}
	return fmt.Sprintf("SKILL.md over target (%d lines, ~%d tokens; target ≤500 lines / ~5k tokens)",
		lines, approxTokens)
}

// descendSingleDir returns the single child directory of dir if that is the only
// entry (release tarballs often wrap content in one folder); otherwise dir.
func descendSingleDir(dir string) string {
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != 1 || !entries[0].IsDir() {
		return dir
	}
	return filepath.Join(dir, entries[0].Name())
}

func reportApply(verb string, rep *installer.Report, err error) int {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}
	if rep.BackupDir != "" {
		fmt.Printf("backed up %d file(s) to %s\n", len(rep.BackedUp), rep.BackupDir)
	}
	fmt.Printf("%s %d file(s)", verb, len(rep.Installed))
	if len(rep.Removed) > 0 {
		fmt.Printf(", removed %d", len(rep.Removed))
	}
	fmt.Println(".")
	if rep.Lock != nil {
		fmt.Printf("lock: %s (agents: %s, skills: %d)\n",
			orDash(rep.Lock.Toolkit.Version), strings.Join(rep.Lock.Agents, ","), len(rep.Lock.Skills))
	}
	return 0
}

func splitCSV(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func printUsage(w *os.File) {
	fmt.Fprint(w, `godot-ai-toolkit — install, update, inspect and validate Godot AI skills.

Usage:
  godot-ai-toolkit <command> [flags]

Commands:
  install    Install selected skills for selected agents into a project.
  update     Update an existing installation (explicit; detects local changes).
  status     Show the installed version and whether managed files changed.
  validate   Validate skills and adapters (SKILL.md contract, manifest, structure).
  catalog    Generate (--write) or check skills/catalog.json from manifests.
  verify     Engine-level verification of skill example files (GDScript parse, C# build).
  version    Print the CLI version.
  help       Show this help.

Examples:
  godot-ai-toolkit install --source . --agents claude,codex --skills example-skill
  godot-ai-toolkit install --repo https://github.com/you/godot-ai-toolkit --tag v0.1.0 \
      --agents claude --skills example-skill
  godot-ai-toolkit update --tag v0.2.0
  godot-ai-toolkit status
  godot-ai-toolkit validate

End users need no Go/Python/Node/uv runtime — only this binary.
`)
}

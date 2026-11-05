// Package verify runs engine-level verification of a skill's example files, as
// declared in its manifest, and reports honest per-check and global statuses.
//
// It is deliberately separate from `validate` (structural + metadata). Runners
// are injectable so the status logic is testable without Godot or the .NET SDK.
// A parse/build check is never reported as runtime verification, and a skipped
// mandatory check never yields a global PASS.
package verify

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/godot-ai-toolkit/godot-ai-toolkit/internal/manifest"
)

// Status is a check/skill/report result.
type Status string

const (
	PASS Status = "PASS"
	FAIL Status = "FAIL"
	SKIP Status = "SKIPPED"
	NA   Status = "NOT_APPLICABLE"
)

// Check is a single verification step's result.
type Check struct {
	Example   string `json:"example,omitempty"`
	ID        string `json:"id"`
	Status    Status `json:"status"`
	Verified  string `json:"verified,omitempty"` // what was actually proven: parse | build | runtime
	Reason    string `json:"reason,omitempty"`
	Mandatory bool   `json:"mandatory"`
}

// SkillReport aggregates a skill's checks.
type SkillReport struct {
	Skill     string  `json:"skill"`
	Reference string  `json:"reference"`
	Checks    []Check `json:"checks"`
	Status    Status  `json:"status"`
}

// Report is the overall result across skills.
type Report struct {
	Skills []SkillReport `json:"skills"`
	Status Status        `json:"status"`
}

// Outcome is what a runner returns for an executed (or skippable) check.
type Outcome struct {
	Status   Status
	Verified string // set on PASS: parse | build | runtime
	Reason   string
}

// Runners are the injectable engine-level operations.
type Runners struct {
	GodotAvailable  func() bool
	DotnetAvailable func() bool
	GDScriptParse   func(absFile string) Outcome
	GDScriptRun     func(absProjectDir string) Outcome
	CSharpBuild     func(absProjectDir string) Outcome
}

// VerifySkill verifies one skill against its manifest.
func VerifySkill(skillDir string, m *manifest.Manifest, r Runners) SkillReport {
	rep := SkillReport{Skill: m.Name, Reference: m.Godot.Reference}

	// Structural gate: the manifest must be valid. (verify assumes validate ran,
	// but stays self-contained.)
	if errs := m.Validate(filepath.Base(skillDir), skillDir); len(errs) > 0 {
		rep.Checks = append(rep.Checks, Check{
			ID: "structural", Status: FAIL, Mandatory: true,
			Reason: strings.Join(errs, "; "),
		})
		rep.Status = aggregate(rep.Checks)
		return rep
	}
	rep.Checks = append(rep.Checks, Check{ID: "structural", Status: PASS, Mandatory: true})

	for _, ex := range m.Examples {
		rep.Checks = append(rep.Checks, checksForExample(skillDir, ex, r)...)
	}
	rep.Status = aggregate(rep.Checks)
	return rep
}

func checksForExample(skillDir string, ex manifest.Example, r Runners) []Check {
	abs := filepath.Join(skillDir, filepath.FromSlash(ex.Path))

	// Illustrative snippets are review-only: not executed.
	if ex.Tier == "illustrative" || ex.Verify == "none" {
		return []Check{{
			Example: ex.Path, ID: ex.Language + ":review", Status: NA, Mandatory: false,
			Reason: "illustrative example (review only)",
		}}
	}

	var checks []Check
	for _, action := range strings.Split(ex.Verify, "+") {
		switch {
		case ex.Language == "gdscript" && action == "parse":
			c := Check{Example: ex.Path, ID: "gdscript:parse", Mandatory: true}
			if !r.GodotAvailable() {
				c.Status, c.Reason = SKIP, "no Godot binary (set GODOT_BIN)"
			} else {
				o := r.GDScriptParse(abs)
				c.Status, c.Verified, c.Reason = o.Status, o.Verified, o.Reason
			}
			checks = append(checks, c)
		case ex.Language == "csharp" && action == "build":
			c := Check{Example: ex.Path, ID: "csharp:build", Mandatory: true}
			if !r.DotnetAvailable() {
				c.Status, c.Reason = SKIP, "no .NET SDK (dotnet not found)"
			} else {
				o := r.CSharpBuild(abs)
				c.Status, c.Verified, c.Reason = o.Status, o.Verified, o.Reason
			}
			checks = append(checks, c)
		case ex.Language == "gdscript" && action == "run":
			// Runtime verification of a GDScript project: it must boot headless
			// and exit 0 (its res://verify.gd asserts behavior and quits non-zero
			// on failure). This proves real behavior, not just that it parses.
			c := Check{Example: ex.Path, ID: "gdscript:runtime", Mandatory: true}
			if !r.GodotAvailable() {
				c.Status, c.Reason = SKIP, "no Godot binary (set GODOT_BIN)"
			} else {
				o := r.GDScriptRun(abs)
				c.Status, c.Verified, c.Reason = o.Status, o.Verified, o.Reason
			}
			checks = append(checks, c)
		case action == "run":
			// C# (and any other) runtime is still deferred: it needs a .NET-enabled
			// Godot build, and gdUnit4Net is intentionally not integrated. Reported
			// SKIPPED so it can never masquerade as runtime-verified.
			checks = append(checks, Check{
				Example: ex.Path, ID: ex.Language + ":runtime", Status: SKIP, Mandatory: true,
				Reason: "C# runtime verification deferred (needs .NET-enabled Godot; gdUnit4Net excluded)",
			})
		default:
			checks = append(checks, Check{
				Example: ex.Path, ID: ex.Language + ":" + action, Status: NA, Mandatory: false,
				Reason: "no verifier for this language/action combination",
			})
		}
	}
	return checks
}

// aggregate computes the global status: FAIL if any check failed; else SKIPPED if
// any mandatory check was skipped (so a skipped mandatory check never yields a
// global PASS); else PASS.
func aggregate(checks []Check) Status {
	hasFail := false
	mandatory, mandatoryPass, mandatorySkip := 0, 0, 0
	for _, c := range checks {
		if c.Status == FAIL {
			hasFail = true
		}
		if c.Mandatory {
			mandatory++
			switch c.Status {
			case PASS:
				mandatoryPass++
			case SKIP:
				mandatorySkip++
			}
		}
	}
	if hasFail {
		return FAIL
	}
	if mandatorySkip > 0 {
		return SKIP
	}
	if mandatory > 0 && mandatoryPass == mandatory {
		return PASS
	}
	return PASS // only non-mandatory/NA checks: nothing to fail
}

// Aggregate combines skill reports into an overall status (FAIL > SKIPPED > PASS).
func Aggregate(reports []SkillReport) Status {
	overall := PASS
	for _, s := range reports {
		switch s.Status {
		case FAIL:
			return FAIL
		case SKIP:
			overall = SKIP
		}
	}
	return overall
}

// DefaultRunners shells out to the real godot/dotnet tools.
func DefaultRunners() Runners {
	return Runners{
		GodotAvailable:  func() bool { return godotBin() != "" },
		DotnetAvailable: func() bool { _, err := exec.LookPath("dotnet"); return err == nil },
		GDScriptParse: func(abs string) Outcome {
			bin := godotBin()
			if bin == "" {
				return Outcome{Status: SKIP, Reason: "no Godot binary (set GODOT_BIN)"}
			}
			out, err := exec.Command(bin, "--headless", "--check-only", "--script", abs).CombinedOutput()
			if hasGodotError(string(out)) {
				return Outcome{Status: FAIL, Reason: "parse errors in output"}
			}
			if err != nil {
				return Outcome{Status: FAIL, Reason: err.Error()}
			}
			return Outcome{Status: PASS, Verified: "parse"}
		},
		GDScriptRun: func(absDir string) Outcome {
			bin := godotBin()
			if bin == "" {
				return Outcome{Status: SKIP, Reason: "no Godot binary (set GODOT_BIN)"}
			}
			// Convention: the project's entry point is res://verify.gd, which
			// asserts behavior and quits non-zero on failure. quit(code) sets a
			// reliable exit code at runtime, so we gate on exit code + explicit
			// failure markers (not a generic "ERROR:" which the engine may print
			// benignly).
			out, err := exec.Command(bin, "--headless", "--path", absDir, "--script", "res://verify.gd").CombinedOutput()
			s := string(out)
			if strings.Contains(s, "SCRIPT ERROR") || strings.Contains(s, "FAIL") {
				return Outcome{Status: FAIL, Reason: "runtime script error or assertion failure in output"}
			}
			if err != nil {
				return Outcome{Status: FAIL, Reason: "non-zero exit: " + err.Error()}
			}
			return Outcome{Status: PASS, Verified: "runtime"}
		},
		CSharpBuild: func(absDir string) Outcome {
			if _, err := exec.LookPath("dotnet"); err != nil {
				return Outcome{Status: SKIP, Reason: "no .NET SDK (dotnet not found)"}
			}
			out, err := exec.Command("dotnet", "build", absDir, "-v", "minimal", "--nologo").CombinedOutput()
			if err != nil {
				return Outcome{Status: FAIL, Reason: buildErrorSummary(string(out))}
			}
			return Outcome{Status: PASS, Verified: "build"}
		},
	}
}

func godotBin() string {
	if b := os.Getenv("GODOT_BIN"); b != "" {
		if _, err := exec.LookPath(b); err == nil {
			return b
		}
		if _, err := os.Stat(b); err == nil {
			return b
		}
	}
	for _, c := range []string{"godot", "godot4", "Godot"} {
		if _, err := exec.LookPath(c); err == nil {
			return c
		}
	}
	return ""
}

func hasGodotError(s string) bool {
	for _, marker := range []string{"SCRIPT ERROR", "Parse Error", "ERROR:", "Failed to load"} {
		if strings.Contains(s, marker) {
			return true
		}
	}
	return false
}

// buildErrorSummary extracts the useful compiler diagnostics (file, line, code,
// message) from `dotnet build` output, rather than its last line. Kept simple:
// it collects the first few lines that look like MSBuild/C# errors.
func buildErrorSummary(out string) string {
	var errs []string
	seen := map[string]bool{}
	for _, line := range strings.Split(out, "\n") {
		l := strings.TrimSpace(line)
		if strings.Contains(l, "): error ") || strings.Contains(l, " error CS") ||
			strings.Contains(l, " error MSB") || strings.Contains(l, " error NU") {
			if seen[l] {
				continue
			}
			seen[l] = true
			errs = append(errs, l)
			if len(errs) >= 3 {
				break
			}
		}
	}
	if len(errs) == 0 {
		return "dotnet build failed (see full output)"
	}
	return strings.Join(errs, " | ")
}

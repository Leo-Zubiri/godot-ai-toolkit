package verify_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/godot-ai-toolkit/godot-ai-toolkit/internal/manifest"
	"github.com/godot-ai-toolkit/godot-ai-toolkit/internal/verify"
)

func demoSkill(t *testing.T, examples []manifest.Example) (string, *manifest.Manifest) {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "demo")
	if err := os.MkdirAll(filepath.Join(dir, "examples", "csharp"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "examples", "hello.gd"), []byte("extends Node\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	m := &manifest.Manifest{
		SchemaVersion: 1, Name: "demo", Languages: []string{"gdscript", "csharp"},
		Godot:    manifest.Godot{Minimum: "4.2", Tested: []string{"4.7.2"}, Reference: "4.7.2"},
		DotNet:   &manifest.DotNet{Minimum: "8.0"},
		Examples: examples,
	}
	return dir, m
}

func fakeRunners(godot, dotnet bool, gd, cs verify.Status) verify.Runners {
	return verify.Runners{
		GodotAvailable:  func() bool { return godot },
		DotnetAvailable: func() bool { return dotnet },
		GDScriptParse:   func(string) verify.Outcome { return verify.Outcome{Status: gd, Verified: "parse"} },
		CSharpBuild:     func(string) verify.Outcome { return verify.Outcome{Status: cs, Verified: "build"} },
	}
}

var bothExamples = []manifest.Example{
	{Path: "examples/hello.gd", Language: "gdscript", Tier: "complete", Verify: "parse"},
	{Path: "examples/csharp", Language: "csharp", Tier: "executable", Verify: "build"},
}

func find(rep verify.SkillReport, id string) (verify.Check, bool) {
	for _, c := range rep.Checks {
		if c.ID == id {
			return c, true
		}
	}
	return verify.Check{}, false
}

func TestAllPass(t *testing.T) {
	dir, m := demoSkill(t, bothExamples)
	rep := verify.VerifySkill(dir, m, fakeRunners(true, true, verify.PASS, verify.PASS))
	if rep.Status != verify.PASS {
		t.Fatalf("global = %s; want PASS (%+v)", rep.Status, rep.Checks)
	}
	if c, ok := find(rep, "gdscript:parse"); !ok || c.Verified != "parse" {
		t.Errorf("gdscript parse not verified-as-parse: %+v", c)
	}
	if c, ok := find(rep, "csharp:build"); !ok || c.Verified != "build" {
		t.Errorf("csharp build not verified-as-build: %+v", c)
	}
}

func TestSkippedMandatoryNeverGlobalPass(t *testing.T) {
	dir, m := demoSkill(t, bothExamples)
	// Godot absent -> gdscript parse is a SKIPPED mandatory check.
	rep := verify.VerifySkill(dir, m, fakeRunners(false, true, verify.PASS, verify.PASS))
	if c, ok := find(rep, "gdscript:parse"); !ok || c.Status != verify.SKIP {
		t.Errorf("expected gdscript parse SKIPPED, got %+v", c)
	}
	if rep.Status == verify.PASS {
		t.Fatalf("a skipped mandatory check must not yield global PASS; got %s", rep.Status)
	}
	if rep.Status != verify.SKIP {
		t.Errorf("global = %s; want SKIPPED", rep.Status)
	}
}

func TestFailPropagates(t *testing.T) {
	dir, m := demoSkill(t, bothExamples)
	rep := verify.VerifySkill(dir, m, fakeRunners(true, true, verify.FAIL, verify.PASS))
	if rep.Status != verify.FAIL {
		t.Fatalf("global = %s; want FAIL", rep.Status)
	}
}

func TestIllustrativeIsNotApplicable(t *testing.T) {
	dir, m := demoSkill(t, []manifest.Example{
		{Path: "examples/hello.gd", Language: "gdscript", Tier: "illustrative", Verify: "none"},
	})
	rep := verify.VerifySkill(dir, m, fakeRunners(true, true, verify.PASS, verify.PASS))
	if c, ok := find(rep, "gdscript:review"); !ok || c.Status != verify.NA {
		t.Errorf("expected NOT_APPLICABLE review check, got %+v", c)
	}
	if rep.Status != verify.PASS {
		t.Errorf("global with only illustrative = %s; want PASS", rep.Status)
	}
}

func TestAggregate(t *testing.T) {
	if verify.Aggregate([]verify.SkillReport{{Status: verify.PASS}, {Status: verify.SKIP}}) != verify.SKIP {
		t.Error("PASS+SKIP should aggregate to SKIP")
	}
	if verify.Aggregate([]verify.SkillReport{{Status: verify.SKIP}, {Status: verify.FAIL}}) != verify.FAIL {
		t.Error("any FAIL should aggregate to FAIL")
	}
	if verify.Aggregate([]verify.SkillReport{{Status: verify.PASS}, {Status: verify.PASS}}) != verify.PASS {
		t.Error("all PASS should aggregate to PASS")
	}
}

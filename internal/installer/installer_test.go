package installer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/godot-ai-toolkit/godot-ai-toolkit/internal/lockfile"
)

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

func makeSource(t *testing.T) string {
	t.Helper()
	src := t.TempDir()
	write(t, filepath.Join(src, "adapters", "claude.json"),
		`{"agent":"claude","skillsPath":".claude/skills","instructionsFile":"CLAUDE.md"}`)
	write(t, filepath.Join(src, "adapters", "opencode.json"),
		`{"agent":"opencode","skillsPath":".claude/skills","instructionsFile":"AGENTS.md"}`)
	write(t, filepath.Join(src, "adapters", "codex.json"),
		`{"agent":"codex","skillsPath":".codex/skills","instructionsFile":"AGENTS.md"}`)
	write(t, filepath.Join(src, "skills", "example-skill", "SKILL.md"),
		"---\nname: example-skill\ndescription: A test skill.\n---\nbody\n")
	write(t, filepath.Join(src, "skills", "second-skill", "SKILL.md"),
		"---\nname: second-skill\ndescription: Another test skill.\n---\nbody\n")
	return src
}

func TestInstallDedupAndLock(t *testing.T) {
	src := makeSource(t)
	proj := t.TempDir()
	if _, err := Install(Options{
		SourceDir: src, ProjectDir: proj,
		Agents: []string{"claude", "opencode", "codex"},
		Skills: []string{"example-skill"}, Version: "v0.0.1", Repository: "r",
	}); err != nil {
		t.Fatalf("install: %v", err)
	}
	if !fileExists(filepath.Join(proj, ".claude/skills/example-skill/SKILL.md")) {
		t.Error("missing claude/opencode skill path")
	}
	if !fileExists(filepath.Join(proj, ".codex/skills/example-skill/SKILL.md")) {
		t.Error("missing codex skill path")
	}
	l, err := lockfile.Load(proj)
	if err != nil {
		t.Fatal(err)
	}
	if len(l.Agents) != 3 {
		t.Errorf("agents = %v", l.Agents)
	}
	if len(l.Skills) != 1 {
		t.Fatalf("skills = %d", len(l.Skills))
	}
	// claude + opencode dedupe to one .claude path; codex adds one .codex path.
	if n := len(l.Skills[0].Files); n != 2 {
		t.Errorf("managed files = %d (%v)", n, l.Skills[0].Files)
	}
}

func TestInstallRefusesOverExisting(t *testing.T) {
	src := makeSource(t)
	proj := t.TempDir()
	o := Options{SourceDir: src, ProjectDir: proj, Agents: []string{"claude"}, Skills: []string{"example-skill"}}
	if _, err := Install(o); err != nil {
		t.Fatal(err)
	}
	if _, err := Install(o); err == nil {
		t.Error("expected install to refuse over an existing lock")
	}
}

func TestConflictDetectionAndForceBackup(t *testing.T) {
	src := makeSource(t)
	proj := t.TempDir()
	if _, err := Install(Options{SourceDir: src, ProjectDir: proj,
		Agents: []string{"claude"}, Skills: []string{"example-skill"}, Version: "v1"}); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(proj, ".claude/skills/example-skill/SKILL.md")
	write(t, target, "LOCALLY EDITED")

	st, err := Status(proj)
	if err != nil {
		t.Fatal(err)
	}
	modified := false
	for _, f := range st.Files {
		if f.Path == ".claude/skills/example-skill/SKILL.md" && f.State == "modified" {
			modified = true
		}
	}
	if !modified {
		t.Error("status did not report the modified file")
	}

	if _, err := Update(Options{SourceDir: src, ProjectDir: proj, Version: "v2"}); err == nil {
		t.Error("expected a conflict error on update without --force")
	}

	rep, err := Update(Options{SourceDir: src, ProjectDir: proj, Version: "v2", Force: true})
	if err != nil {
		t.Fatalf("forced update: %v", err)
	}
	if rep.BackupDir == "" {
		t.Error("expected a backup directory from the forced update")
	}
	got, _ := os.ReadFile(target)
	if string(got) == "LOCALLY EDITED" {
		t.Error("forced update did not restore the managed file")
	}
	st2, _ := Status(proj)
	for _, f := range st2.Files {
		if f.State != "ok" {
			t.Errorf("file %s state %q after forced update", f.Path, f.State)
		}
	}
}

func TestUpdateRemovesDeselectedSkill(t *testing.T) {
	src := makeSource(t)
	proj := t.TempDir()
	if _, err := Install(Options{SourceDir: src, ProjectDir: proj,
		Agents: []string{"claude"}, Skills: []string{"example-skill", "second-skill"}, Version: "v1"}); err != nil {
		t.Fatal(err)
	}
	second := filepath.Join(proj, ".claude/skills/second-skill/SKILL.md")
	if !fileExists(second) {
		t.Fatal("second skill was not installed")
	}
	rep, err := Update(Options{SourceDir: src, ProjectDir: proj,
		Skills: []string{"example-skill"}, Version: "v2"})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if fileExists(second) {
		t.Error("deselected skill should have been removed")
	}
	if len(rep.Removed) == 0 {
		t.Error("expected removals to be reported")
	}
	if len(rep.BackedUp) == 0 {
		t.Error("expected the removed file to be backed up")
	}
	l, _ := lockfile.Load(proj)
	if len(l.Skills) != 1 {
		t.Errorf("skills after update = %d", len(l.Skills))
	}
}

func TestInvalidSkillRejected(t *testing.T) {
	src := t.TempDir()
	write(t, filepath.Join(src, "adapters", "claude.json"),
		`{"agent":"claude","skillsPath":".claude/skills","instructionsFile":"CLAUDE.md"}`)
	// name does not match the folder -> invalid.
	write(t, filepath.Join(src, "skills", "bad", "SKILL.md"),
		"---\nname: not-bad\ndescription: x\n---\n")
	proj := t.TempDir()
	if _, err := Install(Options{SourceDir: src, ProjectDir: proj,
		Agents: []string{"claude"}, Skills: []string{"bad"}}); err == nil {
		t.Error("expected install to reject an invalid skill")
	}
}

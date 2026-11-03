package adapter

import (
	"path/filepath"
	"testing"

	"github.com/godot-ai-toolkit/godot-ai-toolkit/internal/testutil"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name    string
		a       Adapter
		base    string
		wantErr bool
	}{
		{"ok", Adapter{Agent: "claude", SkillsPath: ".claude/skills", InstructionsFile: "CLAUDE.md"}, "claude", false},
		{"missing agent", Adapter{SkillsPath: "x", InstructionsFile: "y"}, "", true},
		{"base mismatch", Adapter{Agent: "claude", SkillsPath: "x", InstructionsFile: "y"}, "codex", true},
		{"abs path", Adapter{Agent: "a", SkillsPath: "/abs", InstructionsFile: "y"}, "a", true},
		{"escape path", Adapter{Agent: "a", SkillsPath: "../x", InstructionsFile: "y"}, "a", true},
		{"missing instructions", Adapter{Agent: "a", SkillsPath: "x"}, "a", true},
	}
	for _, c := range cases {
		if err := c.a.Validate(c.base); (err != nil) != c.wantErr {
			t.Errorf("%s: err=%v wantErr=%v", c.name, err, c.wantErr)
		}
	}
}

func TestLoadAllRepoAdapters(t *testing.T) {
	r, err := testutil.RepoRoot()
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	m, err := LoadAll(filepath.Join(r, "adapters"))
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	for _, want := range []string{"claude", "opencode", "codex"} {
		if _, ok := m[want]; !ok {
			t.Errorf("missing adapter %q", want)
		}
	}
}

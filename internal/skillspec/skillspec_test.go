package skillspec_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/godot-ai-toolkit/godot-ai-toolkit/internal/skillspec"
	"github.com/godot-ai-toolkit/godot-ai-toolkit/internal/testutil"
)

func root(t *testing.T) string {
	t.Helper()
	r, err := testutil.RepoRoot()
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	return r
}

func TestParseValidFixture(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(root(t), "tests/fixtures/skills/valid/example-skill/SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	fm, body, err := skillspec.Parse(content)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if fm.Name != "example-skill" {
		t.Errorf("name = %q", fm.Name)
	}
	if !strings.Contains(fm.Description, "pipeline") {
		t.Errorf("description = %q", fm.Description)
	}
	if len(fm.Keys) != 2 || fm.Keys[0] != "name" || fm.Keys[1] != "description" {
		t.Errorf("frontmatter keys = %v; want [name description]", fm.Keys)
	}
	if !strings.Contains(body, "Example skill") {
		t.Errorf("body missing expected content")
	}
}

func TestValidateValidFixture(t *testing.T) {
	dir := filepath.Join(root(t), "tests/fixtures/skills/valid/example-skill")
	if issues := skillspec.ValidateDir(dir); len(issues) != 0 {
		t.Fatalf("expected no issues, got: %v", issues)
	}
}

func TestValidateInvalidFixtures(t *testing.T) {
	base := filepath.Join(root(t), "tests/fixtures/skills/invalid")
	entries, err := os.ReadDir(base)
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		n++
		dir := filepath.Join(base, e.Name())
		if issues := skillspec.ValidateDir(dir); len(issues) == 0 {
			t.Errorf("%s: expected validation failure, got none", e.Name())
		}
	}
	if n == 0 {
		t.Fatal("no invalid fixtures found")
	}
}

func TestDescriptionTooLong(t *testing.T) {
	skillDir := filepath.Join(t.TempDir(), "too-long")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	desc := strings.Repeat("a", skillspec.MaxDescriptionLen+1)
	content := "---\nname: too-long\ndescription: " + desc + "\n---\nbody\n"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, is := range skillspec.ValidateDir(skillDir) {
		if strings.Contains(is.Message, "description") {
			found = true
		}
	}
	if !found {
		t.Fatal("expected a description-length issue")
	}
}

func TestParseRejectsMissingFrontmatter(t *testing.T) {
	if _, _, err := skillspec.Parse([]byte("# no frontmatter\n")); err == nil {
		t.Fatal("expected error for missing frontmatter")
	}
}

func TestValidateRejectsDisallowedFrontmatterKey(t *testing.T) {
	skillDir := filepath.Join(t.TempDir(), "x")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "---\nname: x\ndescription: ok\nlicense: MIT\n---\nbody\n"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, is := range skillspec.ValidateDir(skillDir) {
		if strings.Contains(is.Message, "not allowed") {
			found = true
		}
	}
	if !found {
		t.Fatal("expected a 'not allowed' frontmatter-key issue for 'license'")
	}
}

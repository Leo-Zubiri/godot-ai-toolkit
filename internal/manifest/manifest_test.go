package manifest_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/godot-ai-toolkit/godot-ai-toolkit/internal/manifest"
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

func TestValidManifest(t *testing.T) {
	m, err := manifest.Load(filepath.Join(root(t), "tests/fixtures/manifests/valid.json"))
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if errs := m.Validate("sample", ""); len(errs) != 0 {
		t.Fatalf("expected no errors, got: %v", errs)
	}
	if !m.HasLanguage("csharp") || !m.HasLanguage("gdscript") {
		t.Errorf("languages = %v", m.Languages)
	}
}

func TestInvalidManifests(t *testing.T) {
	base := filepath.Join(root(t), "tests/fixtures/manifests/invalid")
	entries, err := os.ReadDir(base)
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		n++
		path := filepath.Join(base, e.Name())
		m, loadErr := manifest.Load(path)
		if loadErr != nil {
			continue // rejected at decode time (e.g. unknown field) — correct
		}
		if errs := m.Validate("", ""); len(errs) == 0 {
			t.Errorf("%s: expected a load error or validation errors, got none", e.Name())
		}
	}
	if n == 0 {
		t.Fatal("no invalid manifest fixtures found")
	}
}

func TestUnknownFieldRejectedAtLoad(t *testing.T) {
	_, err := manifest.Load(filepath.Join(root(t), "tests/fixtures/manifests/invalid/unknown-field.json"))
	if err == nil {
		t.Fatal("expected load to reject an unknown field")
	}
}

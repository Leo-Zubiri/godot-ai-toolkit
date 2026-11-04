package catalog_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/godot-ai-toolkit/godot-ai-toolkit/internal/catalog"
)

func writeSkill(t *testing.T, skillsDir, name, manifestJSON string) {
	t.Helper()
	dir := filepath.Join(skillsDir, name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "manifest.json"), []byte(manifestJSON), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestBuildDeterministicAndSorted(t *testing.T) {
	dir := t.TempDir()
	// Intentionally create out of order.
	writeSkill(t, dir, "b-skill", "{\"schemaVersion\":1,\"name\":\"b-skill\",\"languages\":[\"gdscript\"],\"godot\":{\"minimum\":\"4.2\",\"tested\":[\"4.7.2\"],\"reference\":\"4.7.2\"},\"examples\":[],\"sources\":[]}")
	writeSkill(t, dir, "a-skill", "{\"schemaVersion\":1,\"name\":\"a-skill\",\"languages\":[\"gdscript\"],\"godot\":{\"minimum\":\"4.2\",\"tested\":[\"4.7.2\"],\"reference\":\"4.7.2\"},\"examples\":[],\"sources\":[]}")

	cat, err := catalog.Build(dir)
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if len(cat.Skills) != 2 || cat.Skills[0].Name != "a-skill" || cat.Skills[1].Name != "b-skill" {
		t.Fatalf("skills not sorted deterministically: %+v", cat.Skills)
	}
	b1, _ := cat.Marshal()
	cat2, _ := catalog.Build(dir)
	b2, _ := cat2.Marshal()
	if string(b1) != string(b2) {
		t.Error("catalog marshaling is not deterministic")
	}
}

func TestBuildRejectsInvalidManifest(t *testing.T) {
	dir := t.TempDir()
	// name mismatch -> invalid
	writeSkill(t, dir, "x", "{\"schemaVersion\":1,\"name\":\"not-x\",\"languages\":[\"gdscript\"],\"godot\":{\"minimum\":\"4.2\",\"tested\":[\"4.7.2\"],\"reference\":\"4.7.2\"},\"examples\":[],\"sources\":[]}")
	if _, err := catalog.Build(dir); err == nil {
		t.Fatal("expected Build to reject an invalid manifest")
	}
}

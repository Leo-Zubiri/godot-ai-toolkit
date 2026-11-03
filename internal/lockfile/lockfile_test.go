package lockfile

import "testing"

func TestSaveLoadRoundtrip(t *testing.T) {
	dir := t.TempDir()
	l := &Lock{
		Toolkit: ToolkitRef{Repository: "r", Version: "v0.1.0", Commit: "c"},
		Agents:  []string{"claude"},
		Skills:  []Skill{{Name: "s", Files: []ManagedFile{{Path: "p", SHA256: "h"}}}},
	}
	if err := l.Save(dir); err != nil {
		t.Fatal(err)
	}
	if !Exists(dir) {
		t.Fatal("Exists returned false after Save")
	}
	got, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got.Toolkit.Version != "v0.1.0" {
		t.Errorf("version = %q", got.Toolkit.Version)
	}
	if got.SchemaVersion != SchemaVersion {
		t.Errorf("schemaVersion = %d", got.SchemaVersion)
	}
	if got.InstalledAt == "" {
		t.Error("installedAt not set")
	}
	if n := len(got.AllFiles()); n != 1 {
		t.Errorf("AllFiles = %d", n)
	}
}

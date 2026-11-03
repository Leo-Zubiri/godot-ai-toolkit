// Package lockfile reads and writes .godot-ai/toolkit.lock.json, the per-project
// control file that records the pinned toolkit version and installed skills.
package lockfile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const (
	SchemaVersion = 1
	RelDir        = ".godot-ai"
	FileName      = "toolkit.lock.json"
)

// ToolkitRef identifies the source the installation was taken from.
type ToolkitRef struct {
	Repository string `json:"repository,omitempty"`
	Version    string `json:"version,omitempty"`
	Commit     string `json:"commit,omitempty"`
}

// ManagedFile is one installed file and its recorded hash.
type ManagedFile struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

// Skill is one installed skill and the files it placed.
type Skill struct {
	Name   string        `json:"name"`
	Source string        `json:"source,omitempty"`
	Files  []ManagedFile `json:"files"`
}

// Lock is the full control-file document.
type Lock struct {
	SchemaVersion int        `json:"schemaVersion"`
	Toolkit       ToolkitRef `json:"toolkit"`
	Agents        []string   `json:"agents"`
	Skills        []Skill    `json:"skills"`
	InstalledAt   string     `json:"installedAt,omitempty"`
	GeneratedBy   string     `json:"generatedBy,omitempty"`
}

// Path returns the lock file path within projectDir.
func Path(projectDir string) string {
	return filepath.Join(projectDir, RelDir, FileName)
}

// Exists reports whether a lock file is present in projectDir.
func Exists(projectDir string) bool {
	_, err := os.Stat(Path(projectDir))
	return err == nil
}

// Load reads and parses the lock file from projectDir.
func Load(projectDir string) (*Lock, error) {
	data, err := os.ReadFile(Path(projectDir))
	if err != nil {
		return nil, err
	}
	var l Lock
	if err := json.Unmarshal(data, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// Save writes the lock file (pretty-printed) into projectDir.
func (l *Lock) Save(projectDir string) error {
	l.SchemaVersion = SchemaVersion
	l.InstalledAt = time.Now().UTC().Format(time.RFC3339)
	data, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	p := Path(projectDir)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o644)
}

// AllFiles returns every managed file across all skills.
func (l *Lock) AllFiles() []ManagedFile {
	var out []ManagedFile
	for _, s := range l.Skills {
		out = append(out, s.Files...)
	}
	return out
}

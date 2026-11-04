// Package catalog builds a deterministic aggregated index (skills/catalog.json)
// from the per-skill manifests. The catalog is always generated, never hand-
// edited: it is a derived view, not a second source of truth. Determinism (stable
// ordering, no timestamps) lets `catalog --check` detect when it is out of date.
package catalog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/godot-ai-toolkit/godot-ai-toolkit/internal/fsutil"
	"github.com/godot-ai-toolkit/godot-ai-toolkit/internal/manifest"
)

const (
	SchemaVersion = 1
	FileName      = "catalog.json"
)

// Entry is the derived per-skill summary (never authoritative).
type Entry struct {
	Name      string           `json:"name"`
	Languages []string         `json:"languages"`
	Godot     manifest.Godot   `json:"godot"`
	DotNet    *manifest.DotNet `json:"dotnet,omitempty"`
	Examples  int              `json:"examples"`
}

// Catalog is the aggregated index.
type Catalog struct {
	SchemaVersion int     `json:"schemaVersion"`
	Skills        []Entry `json:"skills"`
}

// Path returns the catalog path within a skills directory.
func Path(skillsDir string) string { return filepath.Join(skillsDir, FileName) }

// Build scans skillsDir for skills (directories containing a manifest.json),
// validates each manifest, and returns a deterministic catalog. It returns an
// error naming the first invalid/unloadable manifest.
func Build(skillsDir string) (*Catalog, error) {
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if fsutil.Exists(filepath.Join(skillsDir, e.Name(), manifest.FileName)) {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names) // stable ordering

	cat := &Catalog{SchemaVersion: SchemaVersion, Skills: []Entry{}}
	for _, name := range names {
		dir := filepath.Join(skillsDir, name)
		m, err := manifest.Load(filepath.Join(dir, manifest.FileName))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		if errs := m.Validate(name, dir); len(errs) > 0 {
			return nil, fmt.Errorf("%s: invalid manifest: %s", name, strings.Join(errs, "; "))
		}
		cat.Skills = append(cat.Skills, Entry{
			Name:      m.Name,
			Languages: m.Languages,
			Godot:     m.Godot,
			DotNet:    m.DotNet,
			Examples:  len(m.Examples),
		})
	}
	return cat, nil
}

// Marshal returns the canonical (deterministic) JSON encoding with a trailing
// newline, suitable for writing and for byte-exact staleness comparison.
func (c *Catalog) Marshal() ([]byte, error) {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(b, '\n'), nil
}

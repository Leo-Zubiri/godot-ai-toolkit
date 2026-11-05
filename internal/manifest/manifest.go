// Package manifest defines and validates the toolkit-owned per-skill manifest
// (skills/<name>/manifest.json). The manifest is the source of truth for a
// skill's compatibility, languages, examples and sources; frontmatter stays
// minimal (name + description). See docs/authoring-skills.md.
package manifest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	SchemaVersion = 1
	FileName      = "manifest.json"
)

var (
	validLanguages = map[string]bool{"gdscript": true, "csharp": true}
	validTiers     = map[string]bool{"illustrative": true, "complete": true, "executable": true}
	validVerify    = map[string]bool{"none": true, "parse": true, "build": true, "run": true, "build+run": true}
)

// Godot captures the compatibility model: minimum supported, tested, reference.
type Godot struct {
	Minimum   string   `json:"minimum"`
	Tested    []string `json:"tested"`
	Reference string   `json:"reference"`
}

// DotNet captures C#/.NET requirements (required when a skill lists csharp).
// Minimum is the historical technical floor; Recommended is the current advice;
// Tested lists .NET target frameworks actually built.
type DotNet struct {
	Minimum     string   `json:"minimum,omitempty"`
	Recommended string   `json:"recommended,omitempty"`
	Tested      []string `json:"tested,omitempty"`
	Notes       string   `json:"notes,omitempty"`
}

// Example is one categorized code example associated with the skill.
type Example struct {
	Path     string `json:"path"`     // relative to the skill directory
	Language string `json:"language"` // gdscript | csharp
	Tier     string `json:"tier"`     // illustrative | complete | executable
	Verify   string `json:"verify"`   // none | parse | build | run | build+run
}

// Source is a traceable citation backing a claim in the skill.
type Source struct {
	Claim   string `json:"claim"`
	URL     string `json:"url"`
	Version string `json:"version"`
}

// Verification records the last verification outcome. Written by `verify`,
// never hand-set.
type Verification struct {
	Status      string `json:"status,omitempty"`
	LastRun     string `json:"lastRun,omitempty"`
	LastVersion string `json:"lastVersion,omitempty"`
}

// Manifest is the full per-skill document.
type Manifest struct {
	SchemaVersion       int           `json:"schemaVersion"`
	Name                string        `json:"name"`
	Languages           []string      `json:"languages"`
	Godot               Godot         `json:"godot"`
	DotNet              *DotNet       `json:"dotnet,omitempty"`
	Examples            []Example     `json:"examples"`
	Sources             []Source      `json:"sources"`
	Verification        *Verification `json:"verification,omitempty"`
	SizeExceptionReason string        `json:"sizeExceptionReason,omitempty"`
}

// Load reads and parses a manifest.json, rejecting unknown fields.
func Load(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var m Manifest
	if err := dec.Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

// HasLanguage reports whether the manifest lists a given language.
func (m *Manifest) HasLanguage(lang string) bool {
	for _, l := range m.Languages {
		if l == lang {
			return true
		}
	}
	return false
}

// Validate checks the manifest against the schema and contract. folder is the
// expected skill name (its directory); dir, if non-empty, is the skill directory
// used to check that example paths exist. Returns human-readable error strings.
func (m *Manifest) Validate(folder, dir string) []string {
	var errs []string
	add := func(f string, a ...any) { errs = append(errs, fmt.Sprintf(f, a...)) }

	if m.SchemaVersion != SchemaVersion {
		add("schemaVersion must be %d, got %d", SchemaVersion, m.SchemaVersion)
	}
	if m.Name == "" {
		add("name is required")
	} else if folder != "" && m.Name != folder {
		add("name %q must equal the folder name %q", m.Name, folder)
	}
	if len(m.Languages) == 0 {
		add("languages must list at least one of gdscript/csharp")
	}
	for _, l := range m.Languages {
		if !validLanguages[l] {
			add("unknown language %q (allowed: gdscript, csharp)", l)
		}
	}
	if m.Godot.Reference == "" {
		add("godot.reference is required")
	}
	if m.Godot.Minimum == "" {
		add("godot.minimum is required")
	}
	if len(m.Godot.Tested) == 0 {
		add("godot.tested must list at least one version")
	}
	if m.HasLanguage("csharp") && m.DotNet == nil {
		add("dotnet is required when 'csharp' is a listed language")
	}

	for i, ex := range m.Examples {
		if ex.Path == "" {
			add("examples[%d].path is required", i)
		} else if filepath.IsAbs(ex.Path) || strings.HasPrefix(filepath.ToSlash(ex.Path), "/") ||
			strings.Contains(filepath.ToSlash(ex.Path), "..") {
			add("examples[%d].path %q must be relative and inside the skill", i, ex.Path)
		}
		if !validLanguages[ex.Language] {
			add("examples[%d].language %q is invalid", i, ex.Language)
		}
		if !validTiers[ex.Tier] {
			add("examples[%d].tier %q is invalid", i, ex.Tier)
		}
		if !validVerify[ex.Verify] {
			add("examples[%d].verify %q is invalid", i, ex.Verify)
		}
		// Tier/verify/language consistency.
		if ex.Tier == "illustrative" && ex.Verify != "none" {
			add("examples[%d]: illustrative examples must have verify=none", i)
		}
		if ex.Tier != "illustrative" && ex.Verify == "none" {
			add("examples[%d]: %q examples must declare a verify action", i, ex.Tier)
		}
		if ex.Language == "csharp" && ex.Verify == "parse" {
			add("examples[%d]: C# cannot be parse-verified; use build or build+run", i)
		}
		if ex.Language == "gdscript" && (ex.Verify == "build" || ex.Verify == "build+run") {
			add("examples[%d]: GDScript cannot be build-verified; use parse or run", i)
		}
		if dir != "" && ex.Path != "" {
			if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(ex.Path))); err != nil {
				add("examples[%d]: path %q does not exist", i, ex.Path)
			}
		}
	}

	for i, s := range m.Sources {
		if s.Claim == "" || s.URL == "" || s.Version == "" {
			add("sources[%d] requires claim, url and version", i)
		}
	}
	return errs
}

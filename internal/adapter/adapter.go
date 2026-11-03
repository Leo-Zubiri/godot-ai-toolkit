// Package adapter loads per-agent discovery descriptors (adapters/*.json).
//
// Adapters are DATA, decoupled from the installer core: supporting a new agent
// that adopts the Agent Skills format is a matter of adding one JSON file.
package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Adapter describes where one agent discovers skills and instructions.
type Adapter struct {
	Agent            string `json:"agent"`
	SkillsPath       string `json:"skillsPath"`
	InstructionsFile string `json:"instructionsFile"`
	Notes            string `json:"notes,omitempty"`
}

// Validate checks the descriptor. If fileBase is non-empty, Agent must match it.
func (a Adapter) Validate(fileBase string) error {
	if a.Agent == "" {
		return fmt.Errorf("'agent' is required")
	}
	if fileBase != "" && a.Agent != fileBase {
		return fmt.Errorf("'agent' (%q) must match the file base name (%q)", a.Agent, fileBase)
	}
	if a.SkillsPath == "" {
		return fmt.Errorf("'skillsPath' is required")
	}
	if filepath.IsAbs(a.SkillsPath) {
		return fmt.Errorf("'skillsPath' must be relative, got %q", a.SkillsPath)
	}
	clean := filepath.ToSlash(filepath.Clean(a.SkillsPath))
	if clean == ".." || strings.HasPrefix(clean, "../") {
		return fmt.Errorf("'skillsPath' must not escape the project (no '..'): %q", a.SkillsPath)
	}
	if a.InstructionsFile == "" {
		return fmt.Errorf("'instructionsFile' is required")
	}
	return nil
}

// LoadAll loads and validates every *.json in dir, keyed by agent name.
func LoadAll(dir string) (map[string]Adapter, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	out := map[string]Adapter{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		p := filepath.Join(dir, e.Name())
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", p, err)
		}
		var a Adapter
		dec := json.NewDecoder(bytes.NewReader(data))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&a); err != nil {
			return nil, fmt.Errorf("%s: %w", p, err)
		}
		base := strings.TrimSuffix(e.Name(), ".json")
		if err := a.Validate(base); err != nil {
			return nil, fmt.Errorf("%s: %w", p, err)
		}
		if _, dup := out[a.Agent]; dup {
			return nil, fmt.Errorf("%s: duplicate agent %q", p, a.Agent)
		}
		out[a.Agent] = a
	}
	return out, nil
}

// Names returns the sorted agent names present in m.
func Names(m map[string]Adapter) []string {
	names := make([]string, 0, len(m))
	for k := range m {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

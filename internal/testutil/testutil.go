// Package testutil holds small helpers shared by tests.
package testutil

import (
	"os"
	"path/filepath"
)

// RepoRoot walks up from the current working directory to the module root
// (the directory containing go.mod).
func RepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

// Package installer orchestrates install/update/status: it plans which files to
// place from a toolkit source into a consumer project's per-agent paths, detects
// conflicts, backs up what it overwrites, and records a lock file.
//
// Agent-specific behavior comes entirely from adapter descriptors (data); this
// core never hard-codes agent names.
package installer

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/godot-ai-toolkit/godot-ai-toolkit/internal/adapter"
	"github.com/godot-ai-toolkit/godot-ai-toolkit/internal/fsutil"
	"github.com/godot-ai-toolkit/godot-ai-toolkit/internal/lockfile"
	"github.com/godot-ai-toolkit/godot-ai-toolkit/internal/skillspec"
)

// Options configures an install/update.
type Options struct {
	SourceDir  string   // local toolkit checkout (must contain skills/ and adapters/)
	ProjectDir string   // consumer project directory (default ".")
	Agents     []string // selected agent names (must have matching adapters)
	Skills     []string // selected skill names
	Version    string   // toolkit version to record in the lock
	Commit     string   // source commit to record
	Repository string   // toolkit repository URL to record
	Force      bool     // override conflict protection (still backs up)
}

// Report summarizes what an apply did.
type Report struct {
	Installed []string
	Removed   []string
	BackedUp  []string
	BackupDir string
	Conflicts []string
	Lock      *lockfile.Lock
}

type plannedFile struct {
	Source string // absolute source path
	Skill  string
}

// Install performs a first-time installation. It refuses to run over an existing
// lock unless Force is set (use Update instead).
func Install(o Options) (*Report, error) {
	o.ProjectDir = orDot(o.ProjectDir)
	var prev *lockfile.Lock
	if lockfile.Exists(o.ProjectDir) {
		if !o.Force {
			return nil, fmt.Errorf("a lock file already exists at %s; use 'update' (or --force)", lockfile.Path(o.ProjectDir))
		}
		prev, _ = lockfile.Load(o.ProjectDir)
	}
	return apply(o, prev)
}

// Update reconciles an existing installation to a (possibly new) version or
// selection. Missing selection/repository are inherited from the current lock.
func Update(o Options) (*Report, error) {
	o.ProjectDir = orDot(o.ProjectDir)
	prev, err := lockfile.Load(o.ProjectDir)
	if err != nil {
		return nil, fmt.Errorf("no installation found at %s (run 'install' first): %w", lockfile.Path(o.ProjectDir), err)
	}
	if len(o.Agents) == 0 {
		o.Agents = prev.Agents
	}
	if len(o.Skills) == 0 {
		for _, s := range prev.Skills {
			o.Skills = append(o.Skills, s.Name)
		}
	}
	if o.Repository == "" {
		o.Repository = prev.Toolkit.Repository
	}
	return apply(o, prev)
}

func apply(o Options, prev *lockfile.Lock) (*Report, error) {
	rep := &Report{}
	if o.SourceDir == "" {
		return nil, fmt.Errorf("no source: provide --source <dir> (or a fetched tarball)")
	}
	adapters, err := adapter.LoadAll(filepath.Join(o.SourceDir, "adapters"))
	if err != nil {
		return nil, fmt.Errorf("loading adapters: %w", err)
	}

	agents := uniqueSorted(o.Agents)
	skills := uniqueSorted(o.Skills)
	if len(agents) == 0 {
		return nil, fmt.Errorf("no agents selected (use --agents)")
	}
	if len(skills) == 0 {
		return nil, fmt.Errorf("no skills selected (use --skills)")
	}
	for _, a := range agents {
		if _, ok := adapters[a]; !ok {
			return nil, fmt.Errorf("unknown agent %q (no adapters/%s.json in source)", a, a)
		}
	}

	// Plan target files. Identical target paths (e.g. claude + opencode both use
	// .claude/skills) are placed once.
	planned := map[string]plannedFile{}
	skillFiles := map[string][]string{}
	for _, sk := range skills {
		srcDir := filepath.Join(o.SourceDir, "skills", sk)
		if !fsutil.Exists(filepath.Join(srcDir, "SKILL.md")) {
			return nil, fmt.Errorf("skill %q not found at %s", sk, srcDir)
		}
		if issues := skillspec.ValidateDir(srcDir); len(issues) > 0 {
			return nil, fmt.Errorf("skill %q is invalid:\n  %s", sk, formatIssues(issues))
		}
		rels, err := fsutil.ListFiles(srcDir)
		if err != nil {
			return nil, err
		}
		for _, a := range agents {
			base := filepath.ToSlash(filepath.Clean(adapters[a].SkillsPath))
			for _, rf := range rels {
				targetRel := path.Join(base, sk, rf)
				if _, exists := planned[targetRel]; !exists {
					planned[targetRel] = plannedFile{Source: filepath.Join(srcDir, filepath.FromSlash(rf)), Skill: sk}
					skillFiles[sk] = append(skillFiles[sk], targetRel)
				}
			}
		}
	}
	for k := range skillFiles {
		sort.Strings(skillFiles[k])
	}

	prevFiles := map[string]string{}
	if prev != nil {
		for _, mf := range prev.AllFiles() {
			prevFiles[mf.Path] = mf.SHA256
		}
	}

	// Conflict detection.
	var conflicts []string
	for rel := range planned {
		abs := filepath.Join(o.ProjectDir, filepath.FromSlash(rel))
		if !fsutil.Exists(abs) {
			continue
		}
		cur, err := fsutil.SHA256File(abs)
		if err != nil {
			return nil, err
		}
		if prevSha, tracked := prevFiles[rel]; tracked {
			if cur != prevSha {
				conflicts = append(conflicts, rel+" (locally modified)")
			}
		} else {
			conflicts = append(conflicts, rel+" (exists, not managed by the toolkit)")
		}
	}
	var removals []string
	for rel, prevSha := range prevFiles {
		if _, keep := planned[rel]; keep {
			continue
		}
		removals = append(removals, rel)
		abs := filepath.Join(o.ProjectDir, filepath.FromSlash(rel))
		if fsutil.Exists(abs) {
			cur, err := fsutil.SHA256File(abs)
			if err != nil {
				return nil, err
			}
			if cur != prevSha {
				conflicts = append(conflicts, rel+" (locally modified, would be removed)")
			}
		}
	}
	sort.Strings(conflicts)
	sort.Strings(removals)
	if len(conflicts) > 0 && !o.Force {
		rep.Conflicts = conflicts
		return rep, fmt.Errorf("conflicts detected (use --force to override; files are backed up first):\n  %s",
			strings.Join(conflicts, "\n  "))
	}

	// Back up anything existing we will overwrite or remove.
	ts := time.Now().UTC().Format("20060102T150405Z")
	backupDir := filepath.Join(o.ProjectDir, lockfile.RelDir, "backups", ts)
	toBackup := map[string]bool{}
	for rel := range planned {
		toBackup[rel] = true
	}
	for _, rel := range removals {
		toBackup[rel] = true
	}
	var backupList []string
	for rel := range toBackup {
		abs := filepath.Join(o.ProjectDir, filepath.FromSlash(rel))
		if fsutil.Exists(abs) {
			if err := fsutil.CopyFile(abs, filepath.Join(backupDir, filepath.FromSlash(rel))); err != nil {
				return nil, fmt.Errorf("backup %s: %w", rel, err)
			}
			backupList = append(backupList, rel)
		}
	}
	sort.Strings(backupList)
	rep.BackedUp = backupList
	if len(backupList) > 0 {
		rep.BackupDir = backupDir
	}

	// Write planned files.
	for rel, pf := range planned {
		abs := filepath.Join(o.ProjectDir, filepath.FromSlash(rel))
		if err := fsutil.CopyFile(pf.Source, abs); err != nil {
			return nil, fmt.Errorf("install %s: %w", rel, err)
		}
		rep.Installed = append(rep.Installed, rel)
	}
	// Remove files no longer managed, then prune empty dirs.
	for _, rel := range removals {
		abs := filepath.Join(o.ProjectDir, filepath.FromSlash(rel))
		if fsutil.Exists(abs) {
			if err := os.Remove(abs); err != nil {
				return nil, fmt.Errorf("remove %s: %w", rel, err)
			}
		}
		rep.Removed = append(rep.Removed, rel)
	}
	removeEmptyDirs(o.ProjectDir, removals)
	sort.Strings(rep.Installed)

	// Build and save the new lock.
	newLock := &lockfile.Lock{
		Toolkit:     lockfile.ToolkitRef{Repository: o.Repository, Version: o.Version, Commit: o.Commit},
		Agents:      agents,
		GeneratedBy: strings.TrimSpace("godot-ai-toolkit " + o.Version),
	}
	for _, sk := range skills {
		s := lockfile.Skill{Name: sk, Source: path.Join("skills", sk)}
		for _, rel := range skillFiles[sk] {
			abs := filepath.Join(o.ProjectDir, filepath.FromSlash(rel))
			sha, err := fsutil.SHA256File(abs)
			if err != nil {
				return nil, err
			}
			s.Files = append(s.Files, lockfile.ManagedFile{Path: rel, SHA256: sha})
		}
		newLock.Skills = append(newLock.Skills, s)
	}
	if err := newLock.Save(o.ProjectDir); err != nil {
		return nil, err
	}
	rep.Lock = newLock
	return rep, nil
}

// FileStatus reports one managed file's on-disk state.
type FileStatus struct {
	Path  string
	State string // "ok", "modified", "missing"
}

// StatusReport is the result of Status.
type StatusReport struct {
	Lock  *lockfile.Lock
	Files []FileStatus
}

// Status reads the lock and compares each managed file to its recorded hash.
func Status(projectDir string) (*StatusReport, error) {
	projectDir = orDot(projectDir)
	l, err := lockfile.Load(projectDir)
	if err != nil {
		return nil, fmt.Errorf("no installation found at %s: %w", lockfile.Path(projectDir), err)
	}
	rep := &StatusReport{Lock: l}
	for _, mf := range l.AllFiles() {
		abs := filepath.Join(projectDir, filepath.FromSlash(mf.Path))
		st := FileStatus{Path: mf.Path}
		switch {
		case !fsutil.Exists(abs):
			st.State = "missing"
		default:
			cur, err := fsutil.SHA256File(abs)
			if err != nil {
				return nil, err
			}
			if cur == mf.SHA256 {
				st.State = "ok"
			} else {
				st.State = "modified"
			}
		}
		rep.Files = append(rep.Files, st)
	}
	sort.Slice(rep.Files, func(i, j int) bool { return rep.Files[i].Path < rep.Files[j].Path })
	return rep, nil
}

func removeEmptyDirs(projectDir string, rels []string) {
	stop := filepath.Clean(projectDir)
	for _, rel := range rels {
		dir := filepath.Dir(filepath.Join(projectDir, filepath.FromSlash(rel)))
		for {
			if filepath.Clean(dir) == stop {
				break
			}
			entries, err := os.ReadDir(dir)
			if err != nil || len(entries) > 0 {
				break
			}
			if err := os.Remove(dir); err != nil {
				break
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
}

func orDot(s string) string {
	if s == "" {
		return "."
	}
	return s
}

func uniqueSorted(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s == "" || seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

func formatIssues(issues []skillspec.Issue) string {
	var parts []string
	for _, i := range issues {
		parts = append(parts, i.String())
	}
	return strings.Join(parts, "\n  ")
}
